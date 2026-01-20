# Idea Hamster - Database Integration & Submission System
## Implementation Plan & Design Document

**Project**: Idea Hamster
**Features**: Database Integration (Supabase), Idea Submission Form, Content Moderation, Email Verification
**Status**: Planning Phase
**Reviewed**: Pending

---

## Table of Contents
1. [Executive Summary](#executive-summary)
2. [Architecture Overview](#architecture-overview)
3. [Security Considerations](#security-considerations)
4. [Phased Implementation Plan](#phased-implementation-plan)
5. [Integration Testing Strategy](#integration-testing-strategy)
6. [Best Practices & Patterns](#best-practices--patterns)
7. [Rollback & Recovery Plan](#rollback--recovery-plan)

---

## Executive Summary

### Objectives
- **Database Integration**: Replace mock data with Supabase PostgreSQL backend
- **Idea Submission**: Allow users to submit new ideas with validation
- **Content Moderation**: Block profanity/abuse using Go library (TwiN/go-away)
- **Email Verification**: Require verified email for submissions (reuse voting pattern)

### Key Decisions
- **Database**: Supabase (cloud PostgreSQL with RLS policies)
- **Profanity Filter**: Go library (github.com/TwiN/go-away) - offline, fast
- **Moderation Action**: Block submission with error message
- **Email Requirement**: Yes, require verified email before submission
- **Database Library**: pgx/v5 (already configured)

### Success Criteria
- ✅ All ideas load from database instead of mock data
- ✅ Users can submit new ideas through web form
- ✅ Profanity/abuse is detected and blocked
- ✅ Email verification prevents spam submissions
- ✅ Integration tests validate end-to-end flows
- ✅ Security best practices followed (input validation, SQL injection prevention, XSS protection)

---

## Architecture Overview

### System Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                          CLIENT BROWSER                              │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │  Leaderboard Page (leaderboard.templ)                          │ │
│  │  - Theme Switcher (4 themes)                                   │ │
│  │  - Idea Cards (IdeaCard component)                             │ │
│  │  - Submit Button (SubmitIdeaButton)                            │ │
│  │  - Filter Tabs (TOP VOTED / RECENT / ELIGIBLE)                 │ │
│  └────────────────────────────────────────────────────────────────┘ │
│                              ↕ HTMX (hx-get/hx-post)                 │
└─────────────────────────────────────────────────────────────────────┘
                                      ↕
┌─────────────────────────────────────────────────────────────────────┐
│                       GO WEB SERVER (Chi Router)                     │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │  Routes (cmd/server/main.go)                                   │ │
│  │  - GET  /leaderboard          → HandleLeaderboard              │ │
│  │  - GET  /submit-form          → HandleSubmitForm [NEW]         │ │
│  │  - POST /api/submit-idea      → HandleSubmitIdea [NEW]         │ │
│  │  - POST /api/vote/{ideaID}    → HandleVote                     │ │
│  │  - POST /api/verify-email     → HandleVerifyEmail              │ │
│  │  - GET  /idea/{ideaID}/expand → HandleExpandIdea               │ │
│  └────────────────────────────────────────────────────────────────┘ │
│                              ↕                                       │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │  Handlers Layer (internal/handlers/)                           │ │
│  │  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────┐ │ │
│  │  │ leaderboard.go   │  │  submit.go [NEW] │  │  vote.go     │ │ │
│  │  │ - HandleLeader   │  │  - HandleSubmit  │  │  - HandleVote│ │ │
│  │  │   board()        │  │    Form()        │  │  - HandleExp │ │ │
│  │  │                  │  │  - HandleSubmit  │  │    andIdea() │ │ │
│  │  │                  │  │    Idea()        │  │              │ │ │
│  │  └──────────────────┘  └──────────────────┘  └──────────────┘ │ │
│  │                                                                 │ │
│  │  ┌──────────────────────────────────────────────────────────┐ │ │
│  │  │  moderation.go [NEW]                                      │ │ │
│  │  │  - ValidateContent(title, desc) → (bool, errorMsg)        │ │ │
│  │  │  - SanitizeInput(title, desc)   → (bool, errorMsg)        │ │ │
│  │  │  Uses: github.com/TwiN/go-away (profanity detection)      │ │ │
│  │  └──────────────────────────────────────────────────────────┘ │ │
│  └────────────────────────────────────────────────────────────────┘ │
│                              ↕                                       │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │  Database Layer (internal/database/)                           │ │
│  │  ┌──────────────────┐  ┌──────────────────────────────────┐   │ │
│  │  │ database.go      │  │  queries.go [NEW]                 │   │ │
│  │  │ - DB (pool)      │  │  - GetAllIdeas(sortBy)           │   │ │
│  │  │ - Connect()      │  │  - GetIdeaByID(id)               │   │ │
│  │  │ - Close()        │  │  - CreateIdea(idea)              │   │ │
│  │  │                  │  │  - CreateVote(ideaID, email)     │   │ │
│  │  │                  │  │  - HasUserVoted(ideaID, email)   │   │ │
│  │  └──────────────────┘  └──────────────────────────────────┘   │ │
│  └────────────────────────────────────────────────────────────────┘ │
│                              ↕                                       │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │  Models (internal/models/idea.go)                              │ │
│  │  - type Idea struct (ID, Title, Description, Category, Tags,  │ │
│  │    Submitter, VoteCount, Status, Timestamps)                   │ │
│  │  - type Vote struct (ID, IdeaID, Email, Verified, Timestamp)  │ │
│  │  - Constants: Categories, Statuses                             │ │
│  └────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
                                      ↕
                    pgx/v5 Connection Pool (max 10, min 2)
                                      ↕
┌─────────────────────────────────────────────────────────────────────┐
│                    SUPABASE POSTGRESQL DATABASE                      │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │  ideas Table                                                    │ │
│  │  - id (UUID, PK)                                                │ │
│  │  - title (VARCHAR 80, NOT NULL)                                │ │
│  │  - description (VARCHAR 500, NOT NULL)                         │ │
│  │  - category (VARCHAR 20, CHECK: Frontend/Backend/Full Stack)  │ │
│  │  - tags (TEXT[])                                                │ │
│  │  - submitter (VARCHAR 100)                                      │ │
│  │  - vote_count (INTEGER, DEFAULT 0)                             │ │
│  │  - status (VARCHAR 20, DEFAULT 'pending')                      │ │
│  │  - created_at, updated_at (TIMESTAMPTZ)                        │ │
│  │                                                                  │ │
│  │  Indexes:                                                        │ │
│  │  - idx_ideas_vote_count (vote_count DESC)                      │ │
│  │  - idx_ideas_created_at (created_at DESC)                      │ │
│  │  - idx_ideas_status (status)                                   │ │
│  └────────────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │  votes Table                                                    │ │
│  │  - id (UUID, PK)                                                │ │
│  │  - idea_id (UUID, FK → ideas.id ON DELETE CASCADE)            │ │
│  │  - email (VARCHAR 255, NOT NULL)                               │ │
│  │  - verified (BOOLEAN, DEFAULT FALSE)                           │ │
│  │  - created_at (TIMESTAMPTZ)                                     │ │
│  │  - UNIQUE (idea_id, email) ← Prevents duplicate votes          │ │
│  │                                                                  │ │
│  │  Indexes:                                                        │ │
│  │  - idx_votes_idea_id (idea_id)                                 │ │
│  │  - idx_votes_email (email)                                     │ │
│  └────────────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │  Triggers (Auto-Updates)                                        │ │
│  │  1. update_idea_vote_count()                                    │ │
│  │     - Increments vote_count on INSERT into votes               │ │
│  │     - Decrements vote_count on DELETE from votes               │ │
│  │     - Auto-promotes status to 'eligible' when count ≥ 50       │ │
│  │  2. update_updated_at_column()                                  │ │
│  │     - Sets updated_at = NOW() on any idea UPDATE               │ │
│  └────────────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │  Row Level Security (RLS) Policies                             │ │
│  │  - ideas: SELECT/INSERT allowed for all (anon + auth)          │ │
│  │  - votes: SELECT/INSERT allowed for all, DELETE admin-only     │ │
│  └────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
```

### Data Flow Diagrams

#### Flow 1: Leaderboard Load
```
User → Browser
         ↓ GET /leaderboard
     Chi Router → HandleLeaderboard()
                    ↓ database.GetAllIdeas(sortBy)
                  pgx/v5 → SELECT * FROM ideas ORDER BY [sort]
                    ↓ rows.Scan() into []models.Idea
                  templates.Leaderboard(ideas)
                    ↓ Render HTML
     Browser ← HTMX swap innerHTML into #ideas-list
```

#### Flow 2: Idea Submission (Happy Path)
```
User → Click "SUBMIT IDEA"
         ↓ hx-get="/submit-form"
     HandleSubmitForm()
         ↓ Check email cookie
         ✓ Email exists
         ↓ templates.SubmitForm()
     Browser ← Form loads in modal
         ↓ User fills form
         ↓ POST /api/submit-idea
     HandleSubmitIdea()
         ↓ Parse form (title, desc, category, tags, submitter)
         ↓ SanitizeInput(title, desc)
             ↓ Check length constraints
             ↓ ValidateContent(title, desc)
                 ↓ go-away profanity detection
                 ✓ Clean
         ↓ Validate category
         ↓ Parse tags (comma-separated)
         ↓ database.CreateIdea(idea)
             ↓ INSERT INTO ideas ... RETURNING id, timestamps
             ✓ Success
         ↓ templates.SubmitSuccess()
     Browser ← Success modal
         ↓ Click "VIEW LEADERBOARD"
         ↓ window.location.reload()
     Leaderboard refreshes with new idea
```

#### Flow 3: Content Moderation (Blocked)
```
User → Submit form with profanity
         ↓ POST /api/submit-idea
     HandleSubmitIdea()
         ↓ SanitizeInput(title, desc)
             ↓ ValidateContent(title, desc)
                 ↓ go-away.IsProfane(title) → TRUE
                 ✗ Return (false, "inappropriate language...")
         ↓ templates.SubmitError(msg)
     Browser ← Error modal
         ↓ Click "TRY AGAIN"
         ↓ hx-get="/submit-form"
     Form reloads (user can retry)
```

#### Flow 4: Email Verification Flow
```
User → Click "SUBMIT IDEA" (no email cookie)
         ↓ hx-get="/submit-form"
     HandleSubmitForm()
         ↓ Check email cookie
         ✗ Missing
         ↓ templates.EmailVerificationModal()
     Browser ← Email modal
         ↓ User enters email
         ↓ POST /api/verify-email
     HandleVerifyEmail()
         ↓ Validate email format
         ↓ Set cookie (expires 1 year)
         ✓ Success
     Browser ← Script closes modal, reloads page
         ↓ User clicks "SUBMIT IDEA" again
         ↓ hx-get="/submit-form"
     HandleSubmitForm()
         ↓ Check email cookie
         ✓ Email exists
         ↓ templates.SubmitForm()
     Browser ← Form loads (submission allowed)
```

#### Flow 5: Vote Recording
```
User → Click "VOTE" button
         ↓ hx-post="/api/vote/{ideaID}"
     HandleVote()
         ↓ Check email cookie
         ✓ Email exists
         ↓ database.HasUserVoted(ideaID, email)
             ↓ SELECT EXISTS(SELECT 1 FROM votes WHERE ...)
         ✗ Has not voted yet
         ↓ database.CreateVote(ideaID, email)
             ↓ INSERT INTO votes (idea_id, email, verified) VALUES (...)
             ↓ Trigger: update_idea_vote_count() fires
                 ↓ UPDATE ideas SET vote_count = vote_count + 1
                 ↓ IF vote_count >= 50 THEN status = 'eligible'
             ✓ Vote recorded
         ↓ templates.VoteButton(ideaID, voted=true)
     Browser ← Button updates to "✓ VOTED" (disabled)
```

---

## Security Considerations

### 1. SQL Injection Prevention

**Threat**: Malicious SQL in user input
**Mitigation**: pgx/v5 uses prepared statements automatically

```go
// ✅ SAFE - pgx automatically escapes parameters
query := `INSERT INTO ideas (title, description) VALUES ($1, $2)`
DB.QueryRow(ctx, query, userTitle, userDescription)

// ❌ UNSAFE - Never do this
query := fmt.Sprintf("INSERT INTO ideas (title) VALUES ('%s')", userTitle)
```

**Test**:
```bash
# Try SQL injection in submission
Title: '; DROP TABLE ideas;--
Expected: Title saved as literal string, no SQL executed
```

### 2. Cross-Site Scripting (XSS) Prevention

**Threat**: JavaScript injection in user content
**Mitigation**: Templ automatically escapes HTML

```templ
// ✅ SAFE - templ escapes HTML entities
<h2>{ idea.Title }</h2>

// If user submits: <script>alert('xss')</script>
// Rendered as: &lt;script&gt;alert('xss')&lt;/script&gt;
```

**Test**:
```bash
# Try XSS in submission
Title: <script>alert('XSS')</script>
Expected: Displayed as text, not executed as code
```

### 3. Email Verification Bypass Prevention

**Threat**: Users submitting without verified email
**Mitigation**: Server-side cookie validation

```go
// ✅ Check cookie on every submission
cookie, err := r.Cookie(emailCookieName)
if err != nil || cookie.Value == "" {
    // Reject submission
}
```

**Cookie Security**:
- `HttpOnly: true` - Prevents JavaScript access
- `SameSite: http.SameSiteStrictMode` - CSRF protection
- `Path: "/"` - Available site-wide
- `MaxAge: 365 days` - Long-lived for convenience

**Test**:
```bash
# Try submission without cookie
curl -X POST http://localhost:3001/api/submit-idea -d "title=Test"
Expected: Error "Email verification required"
```

### 4. Content Security Policy (CSP)

**Recommendation**: Add CSP headers to prevent inline script execution

```go
// In main.go middleware
func securityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Security-Policy",
            "default-src 'self'; script-src 'self' https://unpkg.com; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com")
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        next.ServeHTTP(w, r)
    })
}

// Apply to router
r.Use(securityHeaders)
```

### 5. Rate Limiting

**Threat**: Spam submissions or DDoS attacks
**Mitigation**: Implement rate limiting (future enhancement)

```go
// Example using go-chi/httprate
import "github.com/go-chi/httprate"

// Limit to 5 submissions per minute per IP
r.With(httprate.LimitByIP(5, 1*time.Minute)).Post("/api/submit-idea", handlers.HandleSubmitIdea)
```

### 6. Input Validation Defense-in-Depth

**Layer 1 - Client Side**: HTML5 validation
```html
<input required maxlength="80" />
```

**Layer 2 - Server Side**: Length + format checks
```go
if len(title) == 0 || len(title) > 80 {
    return error
}
```

**Layer 3 - Database**: Column constraints
```sql
title VARCHAR(80) NOT NULL
```

**Layer 4 - Content Moderation**: Profanity filter
```go
if detector.IsProfane(title) {
    return error
}
```

### 7. Database Security

**Row Level Security (RLS)**: Already configured in schema
```sql
-- Only allow INSERT and SELECT
ALTER TABLE ideas ENABLE ROW LEVEL SECURITY;
CREATE POLICY "Allow public read" ON ideas FOR SELECT USING (true);
CREATE POLICY "Allow public insert" ON ideas FOR INSERT WITH CHECK (true);
```

**Connection Security**:
- Use SSL/TLS for database connections
- Store credentials in environment variables (never in code)
- Use Supabase connection pooler for security

### 8. OWASP Top 10 Checklist

| Vulnerability | Mitigation |
|---|---|
| A01: Broken Access Control | Email verification, cookie validation |
| A02: Cryptographic Failures | HTTPS enforced, secure cookies |
| A03: Injection | pgx prepared statements, templ escaping |
| A04: Insecure Design | Security-first architecture, validation layers |
| A05: Security Misconfiguration | Security headers, RLS policies |
| A06: Vulnerable Components | Regular dependency updates (go mod tidy) |
| A07: Authentication Failures | Email verification with secure cookies |
| A08: Data Integrity Failures | UNIQUE constraints, foreign keys |
| A09: Logging Failures | Log all errors, security events |
| A10: SSRF | Not applicable (no external fetches) |

---

## Phased Implementation Plan

### Phase 1: Database Integration (Foundation) 🔨

**Duration**: 1-2 hours
**Goal**: Replace mock data with Supabase backend

#### Tasks

**1.1 Initialize Database Connection**
- File: `cmd/server/main.go`
- Add: `database.Connect()` call in main() before router setup
- Add: `defer database.Close()` for cleanup
- Add: Startup logging for connection success/failure
- Test: Server starts without errors, log shows "✓ Database connected"

**1.2 Create Query Functions**
- File: `internal/database/queries.go` (NEW)
- Implement:
  - `GetAllIdeas(ctx, sortBy string) ([]models.Idea, error)`
  - `GetIdeaByID(ctx, id string) (*models.Idea, error)`
  - `CreateIdea(ctx, idea *models.Idea) error`
  - `CreateVote(ctx, ideaID, email string) error`
  - `HasUserVoted(ctx, ideaID, email string) (bool, error)`
- Patterns: Use pgx/v5 Query/QueryRow, proper error wrapping
- Test: Unit test each function with test database

**1.3 Update Handlers to Use Database**
- Files: `internal/handlers/leaderboard.go`, `internal/handlers/vote.go`
- Changes:
  - `HandleLeaderboard()`: Call `database.GetAllIdeas()` instead of `getMockIdeas()`
  - `HandleVote()`: Call `database.CreateVote()` to persist votes
  - `HandleExpandIdea()`: Call `database.GetIdeaByID()` for expansion
- Remove: `getMockIdeas()` function
- Test: Leaderboard loads from database, sorting works, voting persists

**1.4 Environment Configuration**
- Files: `.env`, `.env.example`
- Add: `DATABASE_URL=postgresql://postgres:[PASSWORD]@db.[PROJECT].supabase.co:5432/postgres`
- Document: Setup instructions in PROJECT_README.md
- Test: Connection works with provided credentials

#### Integration Tests - Phase 1

**Test P1.1: Database Connection**
```bash
# Start server
make run

# Expected output:
✓ Database connected successfully
🐹 Idea Hamster starting on http://localhost:3001

# Verify connection
curl http://localhost:3001/leaderboard | grep "IDEA LEADERBOARD"
```

**Test P1.2: Leaderboard Load**
```bash
# Navigate to leaderboard
open http://localhost:3001/leaderboard

# Verify:
- Ideas load from database (not mock data)
- Vote counts match database
- Click "TOP VOTED" → sorted by vote_count DESC
- Click "RECENT" → sorted by created_at DESC
- Click "ELIGIBLE (50+)" → filtered to vote_count >= 50
```

**Test P1.3: Vote Persistence**
```bash
# 1. Get current vote count from database
psql $DATABASE_URL -c "SELECT title, vote_count FROM ideas LIMIT 1;"

# 2. Click VOTE button in browser
# 3. Query database again
psql $DATABASE_URL -c "SELECT title, vote_count FROM ideas LIMIT 1;"

# Verify: vote_count incremented by 1
```

**Test P1.4: Error Handling**
```bash
# Simulate database down: Stop Supabase temporarily
# Restart server
make run

# Expected: Server fails to start with error message
Failed to connect to database: connection refused
```

#### Commit Message - Phase 1
```
feat: integrate Supabase database backend

- Initialize database connection pool in main.go
- Implement query layer with GetAllIdeas, GetIdeaByID, CreateVote
- Update handlers to use database instead of mock data
- Add error handling for database operations
- Configure DATABASE_URL environment variable

Breaking: getMockIdeas() removed, database required to run

Closes: #[issue-number]
Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>
```

---

### Phase 2: Content Moderation (Security) 🛡️

**Duration**: 30 minutes
**Goal**: Block profanity and abuse in submissions

#### Tasks

**2.1 Add Profanity Detection Library**
- Run: `go get github.com/TwiN/go-away`
- Run: `go mod tidy`
- Test: Library imports successfully

**2.2 Create Moderation Module**
- File: `internal/handlers/moderation.go` (NEW)
- Implement:
  - `ValidateContent(fields ...string) (bool, string)` - Checks for profanity
  - `SanitizeInput(title, description string) (bool, string)` - Full validation
- Logic:
  - Check title length (1-80 chars)
  - Check description length (1-500 chars)
  - Check for profanity using go-away
  - Return (false, errorMsg) on any violation
- Test: Unit tests for edge cases

#### Integration Tests - Phase 2

**Test P2.1: Profanity Detection**
```go
// Unit test
func TestValidateContent_Profanity(t *testing.T) {
    valid, msg := ValidateContent("This is damn inappropriate")
    assert.False(t, valid)
    assert.Contains(t, msg, "inappropriate language")
}
```

**Test P2.2: Length Validation**
```go
func TestSanitizeInput_TooLong(t *testing.T) {
    longTitle := strings.Repeat("a", 81)
    valid, msg := SanitizeInput(longTitle, "Description")
    assert.False(t, valid)
    assert.Contains(t, msg, "80 characters")
}
```

**Test P2.3: Clean Content**
```go
func TestSanitizeInput_Valid(t *testing.T) {
    valid, _ := SanitizeInput("Clean Title", "Clean description")
    assert.True(t, valid)
}
```

#### Commit Message - Phase 2
```
feat: add content moderation with profanity filtering

- Integrate TwiN/go-away library for profanity detection
- Create moderation module with ValidateContent and SanitizeInput
- Implement multi-layer validation (length + profanity + format)
- Add comprehensive unit tests for edge cases

Security: Blocks inappropriate content before database insertion

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>
```

---

### Phase 3: Submission Form (User Feature) 📝

**Duration**: 2-3 hours
**Goal**: Allow users to submit ideas through web form

#### Tasks

**3.1 Create Submission Form Template**
- File: `web/templates/submit_form.templ` (NEW)
- Components:
  - `SubmitForm()` - Main form with fields (title, description, category, tags, submitter)
  - `SubmitSuccess()` - Success modal with confetti effect
  - `SubmitError(errorMsg)` - Error modal with retry button
- Styling: Use retro-input, retro-button, theme-primary classes
- HTMX: `hx-post="/api/submit-idea"` with `hx-target="#submit-modal"`
- Test: Templ compiles without errors

**3.2 Create Submission Handlers**
- File: `internal/handlers/submit.go` (NEW)
- Implement:
  - `HandleSubmitForm(w, r)` - Shows form or email modal
  - `HandleSubmitIdea(w, r)` - Processes submission
- Logic for HandleSubmitIdea:
  1. Check email cookie (reject if missing)
  2. Parse form data (r.ParseForm, FormValue)
  3. Sanitize input (call moderation.SanitizeInput)
  4. Validate category (must be in enum)
  5. Parse tags (split by comma, trim)
  6. Create models.Idea object
  7. Call database.CreateIdea()
  8. Return success or error component
- Error handling: Return specific error messages via SubmitError template
- Test: Handler logic tested with mock request

**3.3 Add Routes**
- File: `cmd/server/main.go`
- Add:
  - `r.Get("/submit-form", handlers.HandleSubmitForm)`
  - `r.Post("/api/submit-idea", handlers.HandleSubmitIdea)`
- Test: Routes registered, server starts

**3.4 Update Submit Button**
- File: `web/templates/components.templ`
- Verify: `SubmitIdeaButton()` has correct hx-get="/submit-form"
- Already exists, no changes needed
- Test: Button click triggers form load

#### Integration Tests - Phase 3

**Test P3.1: Full Submission Flow (Happy Path)**
```
1. Open http://localhost:3001/leaderboard
2. Click "SUBMIT IDEA" button
3. If email modal appears:
   - Enter: test@example.com
   - Click VERIFY
   - Wait for page reload
4. Click "SUBMIT IDEA" again
5. Fill form:
   - Title: "Integration Test Idea"
   - Description: "This is a test submission to verify the flow works end-to-end."
   - Category: Full Stack
   - Tags: Test, Demo, Integration
   - Submitter: Tester
6. Click "SUBMIT IDEA"
7. Verify success modal appears with 🎉
8. Click "VIEW LEADERBOARD"
9. Verify idea appears in list
10. Query database:
    psql $DATABASE_URL -c "SELECT title, submitter FROM ideas WHERE title = 'Integration Test Idea';"
11. Verify row exists
```

**Test P3.2: Profanity Rejection**
```
1. Click "SUBMIT IDEA"
2. Fill form:
   - Title: "This damn idea"
   - Description: "Testing profanity filter"
3. Click "SUBMIT IDEA"
4. Verify error modal appears:
   - Message: "Your submission contains inappropriate language..."
5. Click "TRY AGAIN"
6. Verify form reloads
7. Submit clean content
8. Verify success
```

**Test P3.3: Validation Errors**
```
Test A: Empty title
- Leave title blank
- Verify HTML5 validation prevents submission

Test B: Title too long
- Enter 81 characters in title
- Click submit
- Verify error: "Title must be between 1 and 80 characters"

Test C: Description too long
- Enter 501 characters in description
- Click submit
- Verify error: "Description must be between 1 and 500 characters"

Test D: Invalid category
- Manually POST with category="Invalid"
- Verify error: "Invalid category selected"
```

**Test P3.4: Email Requirement**
```
1. Clear cookies (Dev Tools → Application → Cookies → Delete All)
2. Click "SUBMIT IDEA"
3. Verify EmailVerificationModal appears (not form)
4. Cancel modal
5. Reload page
6. Click "SUBMIT IDEA" again
7. Verify modal still appears (email not verified)
8. Enter email, verify
9. Click "SUBMIT IDEA"
10. Verify form now appears
```

**Test P3.5: Anonymous Submission**
```
1. Submit idea with:
   - Submitter field: EMPTY
   - Email cookie: test@example.com
2. Click submit
3. Query database:
   psql $DATABASE_URL -c "SELECT submitter FROM ideas WHERE title = '...';"
4. Verify submitter = 'test@example.com' (defaulted to email)
```

**Test P3.6: Tag Parsing**
```
1. Submit idea with tags: "AI, Mobile,Productivity  , Design"
2. Query database:
   psql $DATABASE_URL -c "SELECT tags FROM ideas WHERE title = '...';"
3. Verify tags = ['AI', 'Mobile', 'Productivity', 'Design'] (trimmed)
```

#### Commit Message - Phase 3
```
feat: add idea submission form with validation

- Create submission form template with retro theme styling
- Implement HandleSubmitForm and HandleSubmitIdea handlers
- Add routes for form display and submission processing
- Integrate content moderation for profanity filtering
- Add comprehensive validation (length, format, category)
- Implement email verification requirement
- Add success/error modals with HTMX swap

Features:
- Fields: title, description, category (dropdown), tags (comma-separated), submitter (optional)
- Validation: HTML5 + server-side + profanity filter
- Email: Required verification before submission
- Anonymous: Defaults to email if submitter field empty

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>
```

---

### Phase 4: Integration & Polish (Quality) ✨

**Duration**: 1-2 hours
**Goal**: End-to-end testing, bug fixes, documentation

#### Tasks

**4.1 End-to-End Testing**
- Test: All flows documented in Integration Tests sections above
- Fix: Any bugs discovered during testing
- Verify: Error messages are user-friendly
- Verify: HTMX swaps work smoothly
- Verify: Theme switcher still works with new modals

**4.2 Security Audit**
- Review: All user inputs are validated
- Review: SQL queries use prepared statements
- Review: Cookies have secure flags
- Test: SQL injection attempts (see Security section)
- Test: XSS attempts (see Security section)
- Add: Security headers (CSP, X-Frame-Options, etc.)

**4.3 Performance Testing**
- Test: Leaderboard load time with 100+ ideas
- Test: Submission response time
- Verify: Database indexes are used (EXPLAIN ANALYZE)
- Optimize: Any slow queries

**4.4 Documentation**
- Update: PROJECT_README.md with setup instructions
- Update: .env.example with all required variables
- Create: TESTING.md with test scenarios
- Add: Code comments for complex logic
- Document: Database schema in migrations/README.md

**4.5 Error Logging**
- Add: Structured logging for errors
- Add: Log submission events (idea created, vote recorded)
- Add: Log security events (profanity blocked, email verification)
- Configure: Log rotation (if using file logging)

#### Integration Tests - Phase 4

**Test P4.1: Concurrent Voting**
```bash
# Open 2 browser tabs
# Tab 1: Click VOTE on idea #1
# Tab 2: Immediately click VOTE on idea #1
# Verify: Only 1 vote recorded (UNIQUE constraint works)
```

**Test P4.2: Database Performance**
```sql
-- Create 100 test ideas
INSERT INTO ideas (title, description, category)
SELECT
    'Test Idea ' || generate_series,
    'Description for test idea ' || generate_series,
    CASE (generate_series % 3)
        WHEN 0 THEN 'Frontend'
        WHEN 1 THEN 'Backend'
        ELSE 'Full Stack'
    END
FROM generate_series(1, 100);

-- Test query performance
EXPLAIN ANALYZE SELECT * FROM ideas ORDER BY vote_count DESC LIMIT 20;
-- Verify: Uses idx_ideas_vote_count index
```

**Test P4.3: Stress Testing**
```bash
# Use Apache Bench (ab) for load testing
ab -n 1000 -c 10 http://localhost:3001/leaderboard

# Verify:
- No errors (0 failed requests)
- Response time < 200ms for 95th percentile
- Server doesn't crash
```

**Test P4.4: Theme Compatibility**
```
1. Load leaderboard
2. Switch to each theme (Original Pink, Arctic Vapor, Rivian Earth, Rivian Cyber)
3. Click "SUBMIT IDEA" in each theme
4. Verify modal colors match theme
5. Submit idea in each theme
6. Verify success modal colors match theme
```

**Test P4.5: Email Persistence**
```
1. Verify email
2. Close browser (not just tab)
3. Reopen browser
4. Go to leaderboard
5. Click "SUBMIT IDEA"
6. Verify form loads immediately (cookie persisted)
```

#### Commit Message - Phase 4
```
chore: integration testing, security hardening, documentation

- Add comprehensive integration test suite
- Implement security headers (CSP, X-Frame-Options, etc.)
- Add structured logging for errors and events
- Optimize database queries with index verification
- Update documentation (README, .env.example, TESTING.md)
- Fix bugs discovered during testing

Security:
- SQL injection protection verified
- XSS protection verified
- Rate limiting considerations documented

Performance:
- Leaderboard loads in <200ms with 100+ ideas
- Database indexes verified with EXPLAIN ANALYZE

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>
```

---

## Integration Testing Strategy

### Testing Levels

```
┌─────────────────────────────────────────────────────────┐
│  Level 1: Unit Tests                                    │
│  - database/queries_test.go (query functions)           │
│  - handlers/moderation_test.go (validation logic)       │
│  - Run: go test ./internal/...                          │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│  Level 2: Integration Tests (Handler + DB)              │
│  - handlers/submit_test.go (with test database)         │
│  - Test database transactions, rollback after test      │
│  - Run: go test -tags=integration ./internal/handlers   │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│  Level 3: End-to-End Tests (Browser)                    │
│  - Manual testing with real browser                     │
│  - Playwright/Selenium automation (future)              │
│  - Test all HTMX interactions                           │
└─────────────────────────────────────────────────────────┘
```

### Test Database Setup

**File**: `internal/database/testing.go` (NEW)

```go
package database

import (
    "context"
    "testing"
)

// SetupTestDB creates a test database connection
func SetupTestDB(t *testing.T) {
    testDBURL := os.Getenv("TEST_DATABASE_URL")
    if testDBURL == "" {
        t.Skip("TEST_DATABASE_URL not set")
    }

    // Connect to test database
    os.Setenv("DATABASE_URL", testDBURL)
    err := Connect()
    if err != nil {
        t.Fatalf("Failed to connect to test database: %v", err)
    }

    // Clean tables before test
    t.Cleanup(func() {
        DB.Exec(context.Background(), "TRUNCATE ideas, votes CASCADE")
        Close()
    })
}
```

### Example Integration Test

**File**: `internal/handlers/submit_test.go` (NEW)

```go
//go:build integration
// +build integration

package handlers

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/abhi10/idea-hamster/internal/database"
)

func TestHandleSubmitIdea_Integration(t *testing.T) {
    database.SetupTestDB(t)

    // Create test request
    form := url.Values{}
    form.Add("title", "Integration Test Idea")
    form.Add("description", "Testing the full stack")
    form.Add("category", "Full Stack")
    form.Add("tags", "Test, Integration")

    req := httptest.NewRequest("POST", "/api/submit-idea", strings.NewReader(form.Encode()))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.AddCookie(&http.Cookie{Name: "voter_email", Value: "test@example.com"})

    w := httptest.NewRecorder()

    // Execute handler
    HandleSubmitIdea(w, req)

    // Assertions
    assert.Equal(t, http.StatusOK, w.Code)

    // Verify database
    ideas, err := database.GetAllIdeas(context.Background(), "votes")
    assert.NoError(t, err)
    assert.Len(t, ideas, 1)
    assert.Equal(t, "Integration Test Idea", ideas[0].Title)
}
```

### Test Execution Plan

**Development Workflow**:
```bash
# 1. Unit tests (fast, run often)
go test ./internal/database -v
go test ./internal/handlers -v

# 2. Integration tests (slower, run before commit)
TEST_DATABASE_URL=$DATABASE_URL go test -tags=integration ./internal/handlers -v

# 3. E2E tests (manual, run before PR)
make run
# Open browser, follow test scenarios
```

**CI/CD Pipeline** (future):
```yaml
# .github/workflows/test.yml
jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: testpass
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
      - run: go test ./...
      - run: go test -tags=integration ./...
        env:
          TEST_DATABASE_URL: postgresql://postgres:testpass@localhost/postgres
```

---

## Best Practices & Patterns

### 1. Repository Pattern (Data Access Layer)

**Current**: Queries in `database/queries.go`
**Future**: Consider repository pattern for testability

```go
// internal/repository/idea_repository.go
type IdeaRepository interface {
    GetAll(ctx context.Context, sortBy string) ([]models.Idea, error)
    GetByID(ctx context.Context, id string) (*models.Idea, error)
    Create(ctx context.Context, idea *models.Idea) error
}

type PostgresIdeaRepository struct {
    db *pgxpool.Pool
}

func (r *PostgresIdeaRepository) GetAll(ctx context.Context, sortBy string) ([]models.Idea, error) {
    // Implementation
}

// Benefits:
// - Easy to mock for testing
// - Swap databases without changing handlers
// - Clear separation of concerns
```

### 2. Context Propagation

**Pattern**: Pass context.Context through all layers

```go
// ✅ GOOD - Context passed from request to database
func HandleLeaderboard(w http.ResponseWriter, r *http.Request) {
    ideas, err := database.GetAllIdeas(r.Context(), sortBy)
    // ...
}

func GetAllIdeas(ctx context.Context, sortBy string) ([]models.Idea, error) {
    rows, err := DB.Query(ctx, query)  // Context passed to pgx
    // ...
}

// Benefits:
// - Request cancellation propagates
// - Timeouts work correctly
// - Distributed tracing support
```

### 3. Error Wrapping

**Pattern**: Add context to errors using fmt.Errorf with %w

```go
// ✅ GOOD - Error wrapped with context
func CreateIdea(ctx context.Context, idea *models.Idea) error {
    err := DB.QueryRow(ctx, query, ...).Scan(...)
    if err != nil {
        return fmt.Errorf("failed to create idea: %w", err)
    }
    return nil
}

// Handler can unwrap and log original error
if err := database.CreateIdea(ctx, idea); err != nil {
    log.Printf("Error creating idea: %v", err)  // Logs full error chain
}

// Benefits:
// - Original error preserved
// - Error messages have context
// - Can use errors.Is() and errors.As()
```

### 4. Structured Logging

**Pattern**: Use structured logs instead of printf

```go
// ❌ Current
log.Printf("Idea created: id=%s, title=%s", idea.ID, idea.Title)

// ✅ Better (with structured logger like zerolog)
log.Info().
    Str("idea_id", idea.ID).
    Str("title", idea.Title).
    Str("submitter", idea.Submitter).
    Msg("Idea created")

// Benefits:
// - Easy to parse and query
// - Can filter by field
// - Machine-readable
```

### 5. Database Transactions

**Pattern**: Use transactions for multi-step operations

```go
// Example: Creating idea and auto-voting in one transaction
func CreateIdeaAndVote(ctx context.Context, idea *models.Idea, email string) error {
    tx, err := DB.Begin(ctx)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback(ctx)  // Rollback if not committed

    // Create idea
    err = tx.QueryRow(ctx, insertIdeaQuery, ...).Scan(&idea.ID)
    if err != nil {
        return fmt.Errorf("failed to create idea: %w", err)
    }

    // Auto-vote
    _, err = tx.Exec(ctx, insertVoteQuery, idea.ID, email)
    if err != nil {
        return fmt.Errorf("failed to create vote: %w", err)
    }

    // Commit both operations
    if err := tx.Commit(ctx); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    return nil
}

// Benefits:
// - Atomic operations
// - No partial states
// - Automatic rollback on error
```

### 6. Dependency Injection

**Pattern**: Inject dependencies instead of using globals

```go
// ❌ Current
var DB *pgxpool.Pool  // Global variable

// ✅ Better
type Server struct {
    db     *pgxpool.Pool
    router *chi.Mux
}

func NewServer(db *pgxpool.Pool) *Server {
    s := &Server{db: db, router: chi.NewRouter()}
    s.routes()
    return s
}

func (s *Server) HandleLeaderboard(w http.ResponseWriter, r *http.Request) {
    ideas, err := GetAllIdeas(r.Context(), s.db, sortBy)
    // ...
}

// Benefits:
// - Easy to test with mock DB
// - Clear dependencies
// - No hidden globals
```

### 7. Configuration Management

**Pattern**: Centralize configuration

```go
// internal/config/config.go
type Config struct {
    Port        string
    DatabaseURL string
    Environment string
    LogLevel    string
}

func Load() (*Config, error) {
    if err := godotenv.Load(); err != nil {
        // OK if .env doesn't exist (use system env)
    }

    cfg := &Config{
        Port:        getEnv("PORT", "3001"),
        DatabaseURL: getEnv("DATABASE_URL", ""),
        Environment: getEnv("ENV", "development"),
        LogLevel:    getEnv("LOG_LEVEL", "info"),
    }

    if cfg.DatabaseURL == "" {
        return nil, fmt.Errorf("DATABASE_URL is required")
    }

    return cfg, nil
}

// main.go
cfg, err := config.Load()
if err != nil {
    log.Fatal(err)
}

// Benefits:
// - Validation in one place
// - Default values
// - Type safety
```

### 8. Graceful Shutdown

**Pattern**: Handle shutdown signals properly

```go
func main() {
    // ... setup ...

    srv := &http.Server{
        Addr:    ":" + port,
        Handler: r,
    }

    // Start server in goroutine
    go func() {
        log.Printf("🐹 Idea Hamster starting on http://localhost:%s\n", port)
        if err := srv.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatalf("Server failed: %v", err)
        }
    }()

    // Wait for interrupt signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    <-sigChan

    log.Println("Shutting down gracefully...")

    // Graceful shutdown with 5 second timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Printf("Server shutdown error: %v", err)
    }

    database.Close()
    log.Println("Server stopped")
}

// Benefits:
// - Finish in-flight requests
// - Close database connections
// - Clean shutdown
```

---

## Rollback & Recovery Plan

### Pre-Implementation Backup

```bash
# 1. Create feature branch
git checkout -b feature/database-submission

# 2. Tag current state
git tag pre-database-integration

# 3. Backup database (if has data)
pg_dump $DATABASE_URL > backup_$(date +%Y%m%d).sql
```

### Phase Rollback Procedures

**If Phase 1 Fails (Database Integration)**:
```bash
# Revert database connection
git checkout main -- cmd/server/main.go

# Revert handlers to use mock data
git checkout main -- internal/handlers/leaderboard.go
git checkout main -- internal/handlers/vote.go

# Remove queries.go
rm internal/database/queries.go

# Restart server
make run
```

**If Phase 2 Fails (Content Moderation)**:
```bash
# Remove moderation module
rm internal/handlers/moderation.go

# Remove library
go mod edit -droprequire github.com/TwiN/go-away
go mod tidy

# Update handlers to skip validation
# (comment out SanitizeInput calls)
```

**If Phase 3 Fails (Submission Form)**:
```bash
# Remove submission files
rm web/templates/submit_form.templ
rm internal/handlers/submit.go

# Remove routes from main.go
git checkout main -- cmd/server/main.go

# Hide submit button (CSS)
echo ".submit-button { display: none; }" >> web/static/css/custom.css
```

### Database Rollback

**If need to revert database changes**:
```bash
# Option 1: Drop all data (development only)
psql $DATABASE_URL -c "TRUNCATE ideas, votes CASCADE;"

# Option 2: Restore from backup
psql $DATABASE_URL < backup_20260120.sql

# Option 3: Rerun migration
psql $DATABASE_URL -f migrations/001_initial_schema.sql
```

### Quick Recovery Commands

```bash
# Emergency: Revert everything to main branch
git reset --hard main
git clean -fd
go mod tidy
make run

# Verify: Server runs with mock data
curl http://localhost:3001/leaderboard
```

---

## Critical Files Reference

### Files to Create (NEW)

| File | Purpose | Lines |
|------|---------|-------|
| `internal/database/queries.go` | Database query functions | ~200 |
| `internal/handlers/moderation.go` | Content validation & profanity filter | ~50 |
| `internal/handlers/submit.go` | Submission form handlers | ~150 |
| `web/templates/submit_form.templ` | Submission form UI components | ~200 |
| `internal/database/testing.go` | Test database setup helpers | ~30 |
| `internal/handlers/submit_test.go` | Integration tests | ~100 |

### Files to Modify

| File | Changes | Impact |
|------|---------|--------|
| `cmd/server/main.go` | Add database.Connect(), submission routes | Medium |
| `internal/handlers/leaderboard.go` | Replace getMockIdeas() with DB query | High |
| `internal/handlers/vote.go` | Add database.CreateVote() call | High |
| `.env` | Add DATABASE_URL | Required |
| `.env.example` | Document DATABASE_URL | Documentation |
| `go.mod` | Add go-away library | Dependencies |

### Files to Review (Existing)

| File | Why |
|------|-----|
| `migrations/001_initial_schema.sql` | Understand database schema |
| `internal/models/idea.go` | Review Idea/Vote structs |
| `internal/database/database.go` | Review connection setup |
| `web/templates/components.templ` | Review SubmitIdeaButton |

---

## Verification Checklist

Before considering implementation complete, verify:

### Database
- [ ] Server starts with database connection
- [ ] Leaderboard loads ideas from database
- [ ] Sorting works (votes, recent, eligible)
- [ ] Voting persists to database
- [ ] Vote count updates automatically (trigger)
- [ ] Duplicate votes prevented (UNIQUE constraint)

### Submission
- [ ] Submit button opens modal
- [ ] Email verification required
- [ ] Form validates input client-side (HTML5)
- [ ] Form validates input server-side (length, profanity)
- [ ] Category dropdown has 3 options
- [ ] Tags parsed correctly (comma-separated)
- [ ] Submitter defaults to email if empty
- [ ] Success modal appears after submission
- [ ] New idea appears in leaderboard

### Security
- [ ] SQL injection attempts blocked (tested)
- [ ] XSS attempts escaped (tested)
- [ ] Profanity blocked (tested)
- [ ] Email cookie required
- [ ] Cookie has HttpOnly + SameSite flags
- [ ] Security headers added (CSP, X-Frame-Options)

### Performance
- [ ] Leaderboard loads < 200ms with 100 ideas
- [ ] Database indexes used (EXPLAIN ANALYZE)
- [ ] No N+1 queries

### Testing
- [ ] All integration tests pass
- [ ] Manual E2E testing complete
- [ ] Error scenarios tested
- [ ] Edge cases covered

### Documentation
- [ ] README updated with setup instructions
- [ ] .env.example has all variables
- [ ] Code comments added for complex logic
- [ ] Commit messages follow conventional commits

---

## Next Steps After Implementation

### Immediate (Day 1-2)
1. Monitor error logs for issues
2. Test with real users (small group)
3. Gather feedback on UX
4. Fix critical bugs

### Short-term (Week 1)
1. Add rate limiting (prevent spam)
2. Implement admin panel (moderate ideas)
3. Add email notifications (idea reaches 50 votes)
4. Improve error messages based on user feedback

### Medium-term (Month 1)
1. Add user profiles (track submissions/votes)
2. Implement idea comments/feedback
3. Add search functionality
4. Build analytics dashboard

### Long-term (Quarter 1)
1. Move to production Supabase instance
2. Set up CI/CD pipeline
3. Implement monitoring (Sentry, DataDog)
4. Add automated E2E tests (Playwright)

---

## Success Metrics

**Technical Metrics**:
- Database uptime: >99.9%
- API response time: p95 < 200ms
- Error rate: <0.1%
- Test coverage: >80%

**Business Metrics**:
- Ideas submitted per day
- Voting engagement rate
- Ideas reaching 50 votes
- User retention (return visitors)

**Security Metrics**:
- Profanity blocks per day
- SQL injection attempts (should be 0)
- XSS attempts (should be 0)
- Authentication failures

---

## Summary

This implementation plan provides a comprehensive, phased approach to adding database integration, idea submission, and content moderation to Idea Hamster. The plan emphasizes:

1. **Security-first**: Multiple validation layers, input sanitization, OWASP compliance
2. **Test-driven**: Integration tests at each phase, clear test scenarios
3. **Best practices**: Repository pattern, context propagation, error wrapping
4. **Maintainability**: Clear separation of concerns, documentation, rollback plan

**Estimated Total Time**: 5-8 hours across 4 phases
**Risk Level**: Medium (database integration always has risks)
**Rollback Plan**: Well-defined, can revert to main branch at any point

**Ready to proceed with implementation once approved.**
