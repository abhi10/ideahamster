# Idea Hamster - System Architecture 🏗️

**Visual Overview of the Complete System**

---

## 🎯 Phase 1: Community Validation (Current Focus)

```
┌─────────────────────────────────────────────────────────────┐
│                         USER BROWSER                         │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  Next.js 15 Frontend (React 18 + Tailwind)            │ │
│  │  • Leaderboard Component                               │ │
│  │  • Idea Submission Form                                │ │
│  │  • Vote Buttons                                        │ │
│  │  • 90's Retro UI Theme                                 │ │
│  └──────────────────┬─────────────────────────────────────┘ │
└─────────────────────┼───────────────────────────────────────┘
                      │
                      │ HTTPS/REST API
                      ▼
┌─────────────────────────────────────────────────────────────┐
│              VERCEL (Hosting + Serverless)                   │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  Next.js API Routes                                    │ │
│  │  • POST /api/ideas (submit new idea)                   │ │
│  │  • GET  /api/ideas (fetch leaderboard)                 │ │
│  │  • POST /api/votes (cast vote)                         │ │
│  │  • GET  /api/ideas/[id] (idea details)                 │ │
│  └──────────────────┬─────────────────────────────────────┘ │
└─────────────────────┼───────────────────────────────────────┘
                      │
                      │ PostgreSQL Protocol
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                  SUPABASE (Backend as a Service)             │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  PostgreSQL Database                                   │ │
│  │  ┌──────────────┐  ┌──────────────┐                   │ │
│  │  │ ideas table  │  │ votes table  │                   │ │
│  │  │ - id         │  │ - id         │                   │ │
│  │  │ - title      │  │ - idea_id    │                   │ │
│  │  │ - description│  │ - voter_email│                   │ │
│  │  │ - category   │  │ - voted_at   │                   │ │
│  │  │ - tags[]     │  └──────────────┘                   │ │
│  │  │ - vote_count │                                      │ │
│  │  │ - status     │                                      │ │
│  │  └──────────────┘                                      │ │
│  │                                                         │ │
│  │  • Row Level Security (RLS)                            │ │
│  │  • Auto-updating vote counts (triggers)                │ │
│  │  • Indexes for performance                             │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘

Cost: ~$1/month
Users: 500-1000/month
Features: Submit, Vote, Leaderboard
```

---

## 🚀 Phase 2: Semi-Automated Building

```
┌─────────────────────────────────────────────────────────────┐
│                    ADMIN DASHBOARD                           │
│  • Review eligible ideas (50+ votes)                         │
│  • Trigger build process                                     │
│  • Monitor build progress                                    │
│  • Approve PRD and deployments                               │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        │ Admin clicks "Start Build"
                        ▼
┌─────────────────────────────────────────────────────────────┐
│                   BUILD ORCHESTRATION                        │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  BullMQ Job Queue (Redis)                              │ │
│  │  ┌──────────────────────────────────────────────────┐ │ │
│  │  │  Job: build-idea-12345                           │ │ │
│  │  │  Status: Processing                               │ │ │
│  │  │  Progress: 45%                                    │ │ │
│  │  │  Phase: "Generating Code"                         │ │ │
│  │  └──────────────────────────────────────────────────┘ │ │
│  └────────────────────────────────────────────────────────┘ │
└───────────────────────┬─────────────────────────────────────┘
                        │
        ┌───────────────┼───────────────┐
        │               │               │
        ▼               ▼               ▼
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│  PHASE 1:    │  │  PHASE 2:    │  │  PHASE 3:    │
│  Generate    │  │  Generate    │  │  Wait for    │
│  PRD         │  │  Tech Stack  │  │  Approval    │
└──────┬───────┘  └──────┬───────┘  └──────┬───────┘
       │                 │                 │
       └─────────────────┴─────────────────┘
                        │
                        │ Approved by Admin
                        ▼
┌─────────────────────────────────────────────────────────────┐
│              ANTHROPIC CLAUDE API (Sonnet 4)                 │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  Prompt: "Generate code files for this PRD..."        │ │
│  │  Output: {                                             │ │
│  │    "package.json": {...},                              │ │
│  │    "src/app/page.tsx": "...",                          │ │
│  │    "src/components/Header.tsx": "...",                 │ │
│  │    ...                                                 │ │
│  │  }                                                     │ │
│  └────────────────────────────────────────────────────────┘ │
│                                                              │
│  Cost per build: ~$10-20                                     │
│  Tokens: 20k-50k tokens                                      │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        │ Files generated
                        ▼
┌─────────────────────────────────────────────────────────────┐
│                ISOLATED BUILD CONTAINER                      │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  Docker Container (node:20-alpine)                     │ │
│  │  ┌──────────────────────────────────────────────────┐ │ │
│  │  │  Security Config:                                 │ │ │
│  │  │  • Memory: 512MB limit                            │ │ │
│  │  │  • CPU: 50% quota                                 │ │ │
│  │  │  • Network: ISOLATED (no internet)                │ │ │
│  │  │  • Filesystem: Read-only root                     │ │ │
│  │  │  • Time limit: 15 minutes                         │ │ │
│  │  └──────────────────────────────────────────────────┘ │ │
│  │                                                         │ │
│  │  Build Steps:                                           │ │
│  │  1. Write files to /app                                 │ │
│  │  2. npm install                                         │ │
│  │  3. npm run build                                       │ │
│  │  4. npm test                                            │ │
│  │  5. Extract output                                      │ │
│  └────────────────────────────────────────────────────────┘ │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        │ Build successful
                        ▼
┌─────────────────────────────────────────────────────────────┐
│                    SECURITY SCANNER                          │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  Checks:                                               │ │
│  │  ✓ npm audit (dependency vulnerabilities)             │ │
│  │  ✓ Secrets detection (API keys, tokens)               │ │
│  │  ✓ Code patterns (eval, dangerouslySetInnerHTML)      │ │
│  │  ✓ OWASP scan                                          │ │
│  │                                                         │ │
│  │  Result: ✅ PASSED (0 critical issues)                 │ │
│  └────────────────────────────────────────────────────────┘ │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        │ Security approved
                        ▼
┌─────────────────────────────────────────────────────────────┐
│                  DEPLOYMENT ROUTER                           │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  if (category === 'frontend') → Vercel                │ │
│  │  if (category === 'fullstack') → Railway              │ │
│  │  if (category === 'backend') → Fly.io                 │ │
│  └────────────────────────────────────────────────────────┘ │
└───────────────────────┬─────────────────────────────────────┘
                        │
        ┌───────────────┼───────────────┐
        │               │               │
        ▼               ▼               ▼
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│   VERCEL     │  │   RAILWAY    │  │    FLY.IO    │
│  (Frontend)  │  │ (Full Stack) │  │  (Backend)   │
└──────┬───────┘  └──────┬───────┘  └──────┬───────┘
       │                 │                 │
       └─────────────────┴─────────────────┘
                        │
                        │ Deploy complete
                        ▼
┌─────────────────────────────────────────────────────────────┐
│                    DEPLOYED APP                              │
│  URL: {slug}.ideahamster.dev                                 │
│  Status: LIVE ✅                                             │
│                                                              │
│  Post-Deploy:                                                │
│  • Update database (status = "built", deploy_url)           │
│  • Email submitter                                           │
│  • Tweet announcement                                        │
│  • Add to archive                                            │
└─────────────────────────────────────────────────────────────┘

Cost: ~$116/month (4 builds)
Build Time: 10-30 minutes
Success Rate Target: 80%+
```

---

## 🔄 Data Flow Diagram

```
USER SUBMITS IDEA
       │
       ▼
┌─────────────────┐
│  Validation     │ → Profanity filter, length checks
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Save to DB     │ → Supabase ideas table
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Display on     │ → Leaderboard (sorted by votes)
│  Leaderboard    │
└────────┬────────┘
         │
         │ Users vote...
         ▼
┌─────────────────┐
│  Vote Count     │ → Increment vote_count
│  Increases      │
└────────┬────────┘
         │
         │ Reaches 50 votes
         ▼
┌─────────────────┐
│  Eligible for   │ → Shows "ELIGIBLE!" badge
│  Build          │
└────────┬────────┘
         │
         │ Admin triggers build
         ▼
┌─────────────────┐
│  Job Queued     │ → BullMQ adds to queue
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  AI Generates   │ → Claude API creates PRD
│  PRD            │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  AI Generates   │ → Claude API recommends stack
│  Tech Stack     │
└────────┬────────┘
         │
         │ Admin approves
         ▼
┌─────────────────┐
│  AI Generates   │ → Claude API writes code
│  Code Files     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Build in       │ → Docker container
│  Container      │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Security Scan  │ → Check for vulnerabilities
└────────┬────────┘
         │
         │ Passed
         ▼
┌─────────────────┐
│  Deploy to      │ → Vercel/Railway/Fly.io
│  Platform       │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  APP LIVE! 🎉  │ → {slug}.ideahamster.dev
└─────────────────┘
```

---

## 🗄️ Database Schema

```
┌─────────────────────────────────────────────────────────────┐
│                         IDEAS TABLE                          │
├──────────────────┬──────────────────────┬───────────────────┤
│ Column           │ Type                 │ Description       │
├──────────────────┼──────────────────────┼───────────────────┤
│ id               │ UUID (PK)            │ Unique ID         │
│ title            │ TEXT (≤80 chars)     │ Idea title        │
│ description      │ TEXT (≤500 chars)    │ Full description  │
│ category         │ ENUM                 │ frontend/backend  │
│ tags             │ TEXT[]               │ Max 5 tags        │
│ submitted_by     │ TEXT                 │ Username/email    │
│ vote_count       │ INTEGER              │ Auto-updated      │
│ status           │ ENUM                 │ active/building   │
│ created_at       │ TIMESTAMPTZ          │ Submission time   │
│ updated_at       │ TIMESTAMPTZ          │ Last modified     │
└──────────────────┴──────────────────────┴───────────────────┘
         │
         │ 1:N relationship
         ▼
┌─────────────────────────────────────────────────────────────┐
│                         VOTES TABLE                          │
├──────────────────┬──────────────────────┬───────────────────┤
│ Column           │ Type                 │ Description       │
├──────────────────┼──────────────────────┼───────────────────┤
│ id               │ UUID (PK)            │ Unique ID         │
│ idea_id          │ UUID (FK)            │ References ideas  │
│ voter_email      │ TEXT                 │ Verified email    │
│ voted_at         │ TIMESTAMPTZ          │ Vote timestamp    │
│ CONSTRAINT       │ UNIQUE(idea_id,      │ One vote per      │
│                  │ voter_email)         │ email per idea    │
└──────────────────┴──────────────────────┴───────────────────┘
         │
         │ Trigger updates vote_count
         ▼
┌─────────────────────────────────────────────────────────────┐
│                        BUILDS TABLE                          │
│                      (Added in Phase 2)                      │
├──────────────────┬──────────────────────┬───────────────────┤
│ Column           │ Type                 │ Description       │
├──────────────────┼──────────────────────┼───────────────────┤
│ id               │ UUID (PK)            │ Unique ID         │
│ idea_id          │ UUID (FK)            │ References ideas  │
│ status           │ ENUM                 │ queued/building   │
│ phase            │ TEXT                 │ Current phase     │
│ prd_content      │ TEXT                 │ Generated PRD     │
│ tech_stack       │ JSONB                │ Tech choices      │
│ deploy_url       │ TEXT                 │ Live app URL      │
│ github_repo      │ TEXT                 │ Source code       │
│ started_at       │ TIMESTAMPTZ          │ Build start       │
│ completed_at     │ TIMESTAMPTZ          │ Build end         │
│ error_message    │ TEXT                 │ If failed         │
└──────────────────┴──────────────────────┴───────────────────┘
```

---

## 🔐 Security Layers

```
┌─────────────────────────────────────────────────────────────┐
│  LAYER 1: Input Validation                                   │
│  • Profanity filter                                          │
│  • Length limits (title, description)                        │
│  • Tag whitelist                                             │
│  • Email format verification                                 │
└─────────────────────────────────────────────────────────────┘
                        ▼
┌─────────────────────────────────────────────────────────────┐
│  LAYER 2: Rate Limiting                                      │
│  • 3 submissions per user per week                           │
│  • 10 votes per IP per hour                                  │
│  • API rate limits (100 req/min)                             │
└─────────────────────────────────────────────────────────────┘
                        ▼
┌─────────────────────────────────────────────────────────────┐
│  LAYER 3: Database Security                                  │
│  • Row Level Security (RLS) enabled                          │
│  • Prepared statements (SQL injection prevention)            │
│  • Encrypted at rest                                         │
│  • Regular backups                                           │
└─────────────────────────────────────────────────────────────┘
                        ▼
┌─────────────────────────────────────────────────────────────┐
│  LAYER 4: Container Isolation (Phase 2)                      │
│  • No network access during build                            │
│  • Resource limits (memory, CPU, time)                       │
│  • Read-only root filesystem                                 │
│  • Dropped capabilities                                      │
└─────────────────────────────────────────────────────────────┘
                        ▼
┌─────────────────────────────────────────────────────────────┐
│  LAYER 5: Code Scanning                                      │
│  • npm audit (dependencies)                                  │
│  • Secret detection (API keys)                               │
│  • Dangerous pattern detection                               │
│  • OWASP checks                                              │
└─────────────────────────────────────────────────────────────┘
                        ▼
┌─────────────────────────────────────────────────────────────┐
│  LAYER 6: Runtime Security                                   │
│  • Subdomain isolation                                       │
│  • CORS restrictions                                         │
│  • CSP headers                                               │
│  • HTTPS only                                                │
└─────────────────────────────────────────────────────────────┘
```

---

## 💰 Cost Scaling

```
PHASE 1 (Validation)
Users:        0-1,000
Builds:       0
Monthly Cost: $1
───────────────────────────────────────
Vercel        $0
Supabase      $0
Domain        $1
───────────────────────────────────────

PHASE 2 (4 builds/month)
Users:        1,000-10,000
Builds:       4/month
Monthly Cost: $116
───────────────────────────────────────
Vercel Pro    $20
Supabase Pro  $25
Claude API    $60  (4 × $15)
Railway       $10  (if needed)
Domain        $1
───────────────────────────────────────

PHASE 3 (8 builds/month)
Users:        10,000-50,000
Builds:       8/month
Monthly Cost: $220
───────────────────────────────────────
Vercel Pro    $20
Supabase Pro  $25
Redis         $10  (Upstash)
Claude API    $120 (8 × $15)
Railway       $20  (2 apps)
Fly.io        $20  (2 apps)
Monitoring    $5   (Sentry)
───────────────────────────────────────

Revenue Options:
• Sponsorships: $50/build
• GitHub Sponsors: $100-500/month
• Premium: Future phases
```

---

## 📈 Scaling Path

```
MONTH 1-2: Phase 1 MVP
├─ Deploy leaderboard
├─ Get 1,000 users
├─ Collect 50+ ideas
└─ Validate concept
    │
    ├─ Success? (50%+ return rate)
    │   └─> Continue to Phase 2
    │
    └─ Failure? (<30% return rate)
        └─> Pivot or stop

MONTH 3-5: Phase 2 Semi-Auto
├─ Add AI generation
├─ Build first 5 apps manually-assisted
├─ Learn and refine prompts
└─ Document process
    │
    ├─ Success? (80%+ build success)
    │   └─> Prepare Phase 3
    │
    └─ Challenges? (50% build success)
        └─> Iterate on Phase 2

MONTH 6-12: Phase 3 Full Auto
├─ Automate PRD approval
├─ Auto-deploy on success
├─ Scale to 2-4 builds/week
└─ Add monetization
    │
    └─> Sustainable, growing platform
```

---

## 🎯 Critical Success Factors

```
┌────────────────────────────────────────────┐
│  PHASE 1 SUCCESS = COMMUNITY ENGAGEMENT    │
│  • 60%+ return visitor rate                │
│  • 3+ ideas reach 50 votes                 │
│  • Active discussion in submissions        │
└────────────────────────────────────────────┘

┌────────────────────────────────────────────┐
│  PHASE 2 SUCCESS = BUILD RELIABILITY       │
│  • 80%+ successful builds                  │
│  • <$20 average cost per build             │
│  • <2 hours human time per build           │
└────────────────────────────────────────────┘

┌────────────────────────────────────────────┐
│  PHASE 3 SUCCESS = SUSTAINABILITY          │
│  • Revenue > Costs                         │
│  • 10,000+ active users                    │
│  • 90%+ build automation                   │
└────────────────────────────────────────────┘
```

---

## 🚀 Deployment Architecture

```
PRODUCTION ENVIRONMENT
═══════════════════════════════════════════════

┌─────────────────────────────────────────────┐
│  VERCEL EDGE NETWORK                        │
│  • Global CDN                               │
│  • Auto SSL                                 │
│  • Serverless Functions                     │
│  • Edge caching                             │
└─────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────┐
│  SUPABASE (Multi-region)                    │
│  • PostgreSQL 15                            │
│  • Realtime subscriptions                   │
│  • Auto backups (daily)                     │
│  • Point-in-time recovery                   │
└─────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────┐
│  DEPLOYED APPS (Multi-platform)             │
│  ┌────────────┐  ┌─────────┐  ┌──────────┐ │
│  │  Vercel    │  │ Railway │  │  Fly.io  │ │
│  │  (Static)  │  │ (Apps)  │  │  (APIs)  │ │
│  └────────────┘  └─────────┘  └──────────┘ │
└─────────────────────────────────────────────┘

MONITORING
═══════════════════════════════════════════════
• Vercel Analytics (performance)
• Sentry (error tracking)
• Custom dashboard (build metrics)
• Supabase logs
```

---

**This architecture supports:**
- ✅ Rapid MVP development (Phase 1)
- ✅ Safe, isolated builds (Phase 2)
- ✅ Multi-platform deployment
- ✅ Cost-effective scaling
- ✅ Security-first design
- ✅ Future automation ready

**Next:** Start with Phase 1 → Validate → Build Phase 2
