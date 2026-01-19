# 🛠️ Development Guide

## Quick Start

### 1. First Time Setup

```bash
# Install dependencies
make install

# Copy environment variables
cp .env.example .env

# Edit .env and add your Supabase credentials
# Get them from: https://supabase.com/dashboard/project/_/settings/api

# Generate Templ templates
make templ

# Build Tailwind CSS
make css
```

### 2. Set up Supabase

1. Go to [supabase.com](https://supabase.com) and create a new project
2. Wait for the project to be ready (~2 minutes)
3. Go to **Settings → API** and copy:
   - Project URL → `SUPABASE_URL`
   - `anon` `public` key → `SUPABASE_ANON_KEY`
4. Go to **Settings → Database** and copy:
   - Connection string → `DATABASE_URL`
   - Replace `[YOUR-PASSWORD]` with your database password
5. Go to **SQL Editor** and run `migrations/001_initial_schema.sql`

### 3. Run Development Server

**Option A: Basic (manual rebuild)**
```bash
# Terminal 1: Start server
make dev

# Terminal 2: Watch CSS changes
make css-watch
```

**Option B: Hot Reload (recommended)**
```bash
# Install Air
go install github.com/air-verse/air@latest

# Terminal 1: Start server with hot reload
air

# Terminal 2: Watch CSS changes
make css-watch
```

Visit **http://localhost:3000** 🎉

## Project Structure

```
idea-hamster/
├── cmd/server/          # Main application
│   └── main.go         # Entry point, router setup
├── internal/           # Private application code
│   ├── handlers/       # HTTP handlers (controllers)
│   │   └── home.go
│   ├── models/         # Data models
│   │   └── idea.go
│   ├── database/       # Database connection
│   │   └── database.go
│   └── middleware/     # Custom middleware (future)
├── web/
│   ├── templates/      # Templ templates (*.templ)
│   │   ├── layout.templ
│   │   └── home.templ
│   └── static/         # Static assets
│       ├── css/
│       │   ├── input.css    # Tailwind source
│       │   └── output.css   # Generated (gitignored)
│       ├── js/
│       └── images/
└── migrations/         # SQL migration files
```

## Development Workflow

### 1. Adding a New Page

```bash
# 1. Create template
touch web/templates/leaderboard.templ
```

```go
// web/templates/leaderboard.templ
package templates

templ Leaderboard() {
    @Layout("Leaderboard") {
        <div class="container mx-auto px-4 py-8">
            <h1 class="text-4xl font-pixel neon-glow-pink">Leaderboard</h1>
            <!-- Your content -->
        </div>
    }
}
```

```bash
# 2. Generate Templ code
make templ
```

```go
// 3. Create handler
// internal/handlers/leaderboard.go
package handlers

import (
    "net/http"
    "github.com/abhishekrajuchamarthi/idea-hamster/web/templates"
)

func HandleLeaderboard(w http.ResponseWriter, r *http.Request) {
    component := templates.Leaderboard()
    component.Render(r.Context(), w)
}
```

```go
// 4. Add route
// cmd/server/main.go
r.Get("/leaderboard", handlers.HandleLeaderboard)
```

```bash
# 5. Test
make dev
# Visit: http://localhost:3000/leaderboard
```

### 2. Styling with Tailwind

Use the custom classes defined in `web/static/css/input.css`:

```html
<!-- Retro button -->
<button class="retro-button">Click Me!</button>

<!-- Retro card -->
<div class="retro-card">
    <h2>Card Title</h2>
</div>

<!-- Neon glow text -->
<h1 class="neon-glow-pink">Glowing Title</h1>

<!-- Custom colors -->
<div class="bg-gray-900 text-neon-cyan border-neon-pink">
    Custom colors
</div>
```

### 3. Working with Database

```go
// Example: Fetch all ideas
package handlers

import (
    "context"
    "github.com/abhishekrajuchamarthi/idea-hamster/internal/database"
    "github.com/abhishekrajuchamarthi/idea-hamster/internal/models"
)

func GetIdeas(ctx context.Context) ([]models.Idea, error) {
    query := `
        SELECT id, title, description, category, tags,
               vote_count, status, created_at
        FROM ideas
        ORDER BY vote_count DESC
        LIMIT 10
    `

    rows, err := database.DB.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var ideas []models.Idea
    for rows.Next() {
        var idea models.Idea
        err := rows.Scan(
            &idea.ID, &idea.Title, &idea.Description,
            &idea.Category, &idea.Tags, &idea.VoteCount,
            &idea.Status, &idea.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        ideas = append(ideas, idea)
    }

    return ideas, nil
}
```

### 4. HTMX Integration

```html
<!-- Example: Vote button with HTMX -->
<button
    hx-post="/api/vote"
    hx-vals='{"idea_id": "123"}'
    hx-swap="outerHTML"
    class="retro-button"
>
    Vote ⬆️
</button>

<!-- Response will replace the button -->
```

```go
// Handler example
func HandleVote(w http.ResponseWriter, r *http.Request) {
    // Process vote...

    // Return updated button
    w.Write([]byte(`
        <span class="text-neon-green">
            ✓ Voted!
        </span>
    `))
}
```

## Common Commands

```bash
# Development
make dev              # Run server (rebuild required)
make css-watch        # Watch CSS changes
air                   # Run with hot reload (recommended)

# Building
make build            # Build production binary
make templ            # Generate Templ templates
make css              # Build CSS (minified)

# Testing
make test             # Run all tests
make fmt              # Format code
go test -v ./...      # Run tests with verbose output

# Database
# Run migrations in Supabase SQL Editor
# See: migrations/001_initial_schema.sql

# Cleaning
make clean            # Remove build artifacts
```

## Debugging

### Server won't start

```bash
# Check if port is in use
lsof -i :3000

# Kill process on port 3000
kill -9 $(lsof -t -i:3000)

# Check logs
tail -f build-errors.log  # If using Air
```

### Templ templates not updating

```bash
# Regenerate templates
make templ

# Check for errors
templ generate -v
```

### CSS not updating

```bash
# Check if Tailwind is watching
npm run dev:css

# Rebuild CSS
make css

# Clear browser cache (Cmd+Shift+R on Mac)
```

### Database connection fails

```bash
# Test connection
psql "postgresql://postgres:[password]@db.project.supabase.co:5432/postgres"

# Check .env file
cat .env | grep DATABASE_URL

# Verify Supabase project is running
# Visit: https://supabase.com/dashboard
```

## Git Workflow

```bash
# Create feature branch
git checkout -b feature/idea-submission-form

# Make changes, commit regularly
git add .
git commit -m "feat: add idea submission form"

# Push to remote
git push origin feature/idea-submission-form

# Create PR on GitHub
```

## Code Style

- **Go**: Follow [Effective Go](https://go.dev/doc/effective_go)
- **Templ**: Use semantic HTML, keep components small
- **CSS**: Use Tailwind utilities, create custom classes for repeated patterns
- **Commits**: Use [Conventional Commits](https://www.conventionalcommits.org/)
  - `feat:` - New feature
  - `fix:` - Bug fix
  - `docs:` - Documentation
  - `style:` - Formatting
  - `refactor:` - Code refactoring
  - `test:` - Tests
  - `chore:` - Maintenance

## Performance Tips

1. **Database**: Use indexes, limit queries, use pagination
2. **CSS**: Only include used Tailwind classes (automatic)
3. **Images**: Use WebP format, optimize with ImageOptim
4. **HTMX**: Use `hx-target` and `hx-swap` for partial updates
5. **Caching**: Add caching headers for static assets

## Resources

- [Templ Docs](https://templ.guide/)
- [HTMX Docs](https://htmx.org/)
- [Chi Router Docs](https://go-chi.io/)
- [Tailwind CSS Docs](https://tailwindcss.com/)
- [Supabase Docs](https://supabase.com/docs)

---

**Happy coding! 🐹✨**
