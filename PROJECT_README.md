# 🐹 Idea Hamster - Phase 1 Implementation

> **Your Ideas. Our Hamster. Pure 90's Magic!**

A community-driven leaderboard for voting on app ideas, built with **Go + Templ + HTMX + Tailwind CSS**.

## 🎯 What is this?

Idea Hamster is a retro-themed platform where users:
1. **Submit** app ideas
2. **Vote** on their favorites
3. **Watch** top-voted ideas (50+ votes) get queued for building in Phase 2

This is **Phase 1** - focused on community validation with a beautiful 90's aesthetic.

## 🛠️ Tech Stack

| Component | Technology |
|-----------|------------|
| **Language** | Go 1.21+ |
| **Templates** | Templ (type-safe HTML) |
| **Frontend** | HTMX + Tailwind CSS |
| **Database** | Supabase (PostgreSQL) |
| **Router** | Chi (lightweight, idiomatic) |
| **Hosting** | Railway / Fly.io / VPS |

## 📁 Project Structure

```
idea-hamster/
├── cmd/
│   └── server/          # Application entry point
│       └── main.go
├── internal/
│   ├── handlers/        # HTTP request handlers
│   ├── models/          # Data models
│   ├── database/        # Database connection
│   └── middleware/      # Custom middleware
├── web/
│   ├── templates/       # Templ templates (*.templ)
│   └── static/          # Static assets
│       ├── css/         # Tailwind CSS
│       ├── js/          # JavaScript (minimal)
│       └── images/      # Images and icons
├── migrations/          # SQL migration files
├── docs/                # Documentation (PRD, etc.)
├── Makefile             # Development commands
└── go.mod               # Go dependencies
```

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+** - [Install Go](https://go.dev/doc/install)
- **Node.js 18+** - [Install Node](https://nodejs.org/) (for Tailwind CSS)
- **Templ CLI** - Install via: `go install github.com/a-h/templ/cmd/templ@latest`
- **Supabase Account** - [Sign up](https://supabase.com/)

### 1. Clone & Setup

```bash
# Clone the repository
git clone https://github.com/abhishekrajuchamarthi/idea-hamster.git
cd idea-hamster

# Install dependencies
make install

# Create .env file
cp .env.example .env
```

### 2. Configure Supabase

1. Create a new Supabase project at [supabase.com](https://supabase.com)
2. Go to **Settings → Database** and copy the connection string
3. Go to **Settings → API** and copy:
   - `SUPABASE_URL`
   - `SUPABASE_ANON_KEY`
4. Update your `.env` file with these values

### 3. Run Database Migrations

```bash
# Go to Supabase SQL Editor
# Copy and paste: migrations/001_initial_schema.sql
# Run the SQL
```

### 4. Start Development Server

```bash
# Terminal 1: Start the server
make dev

# Terminal 2: Watch CSS changes
make css-watch
```

Visit **http://localhost:3000** 🎉

## 📝 Available Commands

```bash
make help         # Show all available commands
make install      # Install all dependencies
make dev          # Run development server
make build        # Build production binary
make test         # Run tests
make clean        # Clean build artifacts
make fmt          # Format code
```

## 🗄️ Database Schema

### `ideas` table
| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| title | VARCHAR(80) | Idea title |
| description | VARCHAR(500) | Idea description |
| category | VARCHAR(20) | Frontend/Backend/Full Stack |
| tags | TEXT[] | Array of tags |
| submitter | VARCHAR(100) | Optional submitter name |
| vote_count | INTEGER | Number of votes |
| status | VARCHAR(20) | pending/eligible/building/built |
| created_at | TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | Last update timestamp |

### `votes` table
| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| idea_id | UUID | Foreign key to ideas |
| email | VARCHAR(255) | Voter email |
| verified | BOOLEAN | Email verification status |
| created_at | TIMESTAMP | Vote timestamp |

**Unique constraint:** One vote per email per idea

## 🎨 Design System

### Colors (Neon 90's Palette)
- **Pink:** `#FF10F0` - Primary accent
- **Cyan:** `#00FFFF` - Secondary accent
- **Purple:** `#9D00FF` - Tertiary accent
- **Green:** `#39FF14` - Success states
- **Yellow:** `#FFFF00` - Warnings

### Fonts
- **Headings:** Press Start 2P (pixel font)
- **Body:** VT323 (monospace retro)

### Components
- Retro buttons with glow effects
- CRT scanline overlay
- Pixelated image rendering
- Animated hamster mascot (coming soon)

## 🔧 Development Workflow

### 1. Create a Feature
```bash
# Create a new branch
git checkout -b feature/your-feature-name

# Make changes
# Edit .templ files for templates
# Edit .go files for logic

# Generate Templ code
make templ

# Test your changes
make dev
```

### 2. Commit Your Changes
```bash
git add .
git commit -m "feat: add your feature description"
git push origin feature/your-feature-name
```

### 3. Testing
```bash
# Run tests
make test

# Format code
make fmt
```

## 📦 Production Build

```bash
# Build binary
make build

# Run production server
./bin/idea-hamster
```

### Deploy to Railway
```bash
# Install Railway CLI
npm i -g @railway/cli

# Login
railway login

# Deploy
railway up
```

### Deploy to Fly.io
```bash
# Install Fly CLI
curl -L https://fly.io/install.sh | sh

# Login
fly auth login

# Deploy
fly deploy
```

## 🌟 Phase 1 Features

- [x] Project structure and setup
- [ ] Idea submission form
- [ ] Voting system with email verification
- [ ] Leaderboard with sorting/filtering
- [ ] Build queue page (shows 50+ vote ideas)
- [ ] 90's retro UI polish
- [ ] Mobile responsive design
- [ ] Rate limiting and security
- [ ] Analytics tracking

## 🚀 Phase 2 (Future)

Once Phase 1 validates community interest:
- AI-assisted app building (Claude API)
- Docker container isolation
- Semi-automated deployment
- Build progress tracking
- Deployed app showcase

## 📚 Documentation

- [PRD.md](./PRD.md) - Original product requirements
- [MVP_REVISED.md](./MVP_REVISED.md) - Phased approach strategy
- [DESIGN_GUIDE.md](./DESIGN_GUIDE.md) - 90's aesthetic guidelines
- [ARCHITECTURE.md](./ARCHITECTURE.md) - Technical architecture
- [QUICK_START.md](./QUICK_START.md) - Quick setup guide

## 🤝 Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## 📄 License

MIT License - see [LICENSE](./LICENSE) for details

## 🙏 Acknowledgments

- Inspired by [Gastown](https://github.com/steveyegge/gastown)
- Retro aesthetic inspired by 90's web culture
- Built with ❤️ and nostalgia

## 📧 Contact

Questions? Ideas? Reach out:
- **GitHub:** [@abhishekrajuchamarthi](https://github.com/abhishekrajuchamarthi)
- **Website:** [araju.dev](https://www.araju.dev)

---

**Let's build something retro! 🐹✨**

*Status: Phase 1 in Development*
*Started: January 2026*
