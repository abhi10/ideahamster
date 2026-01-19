# Idea Hamster - Quick Start Guide 🚀

**Get Phase 1 MVP running in under 30 minutes**

---

## 🎯 What You're Building (Phase 1)

A beautiful 90's retro-themed leaderboard where users can:
- Submit app ideas
- Vote on their favorites
- See which ideas hit the 50-vote threshold for building

**No AI automation yet** - that's Phase 2. First, we validate the concept!

---

## 📋 Prerequisites

- Node.js 20+ installed
- Git installed
- Supabase account (free tier)
- Vercel account (free tier)
- Code editor (VS Code recommended)

---

## 🚀 Setup (30 minutes)

### Step 1: Create Supabase Project (5 min)

1. Go to [supabase.com](https://supabase.com)
2. Click "New Project"
3. Name: `idea-hamster`
4. Database Password: (save this somewhere safe)
5. Region: Choose closest to you
6. Click "Create new project"

**While it's setting up, continue to Step 2...**

### Step 2: Initialize Next.js Project (5 min)

```bash
# Create new Next.js project
npx create-next-app@latest idea-hamster

# Options to select:
# ✓ TypeScript? Yes
# ✓ ESLint? Yes
# ✓ Tailwind CSS? Yes
# ✓ `src/` directory? Yes
# ✓ App Router? Yes
# ✓ Customize default import alias? No

cd idea-hamster

# Install additional dependencies
npm install @supabase/supabase-js @supabase/ssr zustand framer-motion
npm install -D @types/node
```

### Step 3: Configure Supabase (10 min)

1. Go back to Supabase dashboard
2. Click on your project (should be ready now)
3. Go to **Settings** → **API**
4. Copy these values:
   - Project URL
   - `anon` public key

5. Create `.env.local` in your project root:

```bash
NEXT_PUBLIC_SUPABASE_URL=your-project-url
NEXT_PUBLIC_SUPABASE_ANON_KEY=your-anon-key
```

6. Create database schema:
   - Go to **SQL Editor** in Supabase
   - Click "New query"
   - Paste this SQL:

```sql
-- Ideas table
CREATE TABLE ideas (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title TEXT NOT NULL CHECK (char_length(title) <= 80),
  description TEXT NOT NULL CHECK (char_length(description) <= 500),
  category TEXT NOT NULL CHECK (category IN ('frontend', 'backend', 'fullstack')),
  tags TEXT[] DEFAULT '{}',
  submitted_by TEXT DEFAULT 'Anonymous',
  vote_count INTEGER DEFAULT 0,
  status TEXT DEFAULT 'active' CHECK (status IN ('active', 'queued', 'building', 'built')),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Votes table (track unique votes)
CREATE TABLE votes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  idea_id UUID REFERENCES ideas(id) ON DELETE CASCADE,
  voter_email TEXT NOT NULL,
  voted_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(idea_id, voter_email)
);

-- Enable Row Level Security
ALTER TABLE ideas ENABLE ROW LEVEL SECURITY;
ALTER TABLE votes ENABLE ROW LEVEL SECURITY;

-- Policies (allow all for MVP, tighten later)
CREATE POLICY "Ideas are viewable by everyone" ON ideas
  FOR SELECT USING (true);

CREATE POLICY "Anyone can insert ideas" ON ideas
  FOR INSERT WITH CHECK (true);

CREATE POLICY "Votes are viewable by everyone" ON votes
  FOR SELECT USING (true);

CREATE POLICY "Anyone can insert votes" ON votes
  FOR INSERT WITH CHECK (true);

-- Indexes for performance
CREATE INDEX idx_ideas_vote_count ON ideas(vote_count DESC);
CREATE INDEX idx_ideas_created_at ON ideas(created_at DESC);
CREATE INDEX idx_votes_idea_id ON votes(idea_id);
CREATE INDEX idx_votes_email ON votes(voter_email);

-- Function to update vote count
CREATE OR REPLACE FUNCTION update_idea_vote_count()
RETURNS TRIGGER AS $$
BEGIN
  UPDATE ideas
  SET vote_count = (
    SELECT COUNT(*) FROM votes WHERE idea_id = NEW.idea_id
  ),
  updated_at = NOW()
  WHERE id = NEW.idea_id;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to auto-update vote count
CREATE TRIGGER vote_count_trigger
AFTER INSERT ON votes
FOR EACH ROW
EXECUTE FUNCTION update_idea_vote_count();
```

7. Click "Run" - you should see "Success. No rows returned"

### Step 4: Create Supabase Client (5 min)

Create `src/lib/supabase.ts`:

```typescript
import { createClient } from '@supabase/supabase-js';

const supabaseUrl = process.env.NEXT_PUBLIC_SUPABASE_URL!;
const supabaseKey = process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!;

export const supabase = createClient(supabaseUrl, supabaseKey);

// Types
export interface Idea {
  id: string;
  title: string;
  description: string;
  category: 'frontend' | 'backend' | 'fullstack';
  tags: string[];
  submitted_by: string;
  vote_count: number;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface Vote {
  id: string;
  idea_id: string;
  voter_email: string;
  voted_at: string;
}
```

### Step 5: Basic Retro Styles (5 min)

Update `src/app/globals.css`:

```css
@import url('https://fonts.googleapis.com/css2?family=Press+Start+2P&family=VT323&display=swap');

@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  :root {
    --neon-pink: #FF10F0;
    --electric-blue: #00FFFF;
    --lime-green: #39FF14;
    --cyber-purple: #9D00FF;
    --deep-black: #0A0A0A;
    --space-gray: #1A1A2E;
  }
  
  body {
    background: var(--deep-black);
    color: var(--electric-blue);
    font-family: 'VT323', monospace;
  }
}

@layer components {
  .pixel-font {
    font-family: 'Press Start 2P', monospace;
  }
  
  .neon-border {
    border: 2px solid var(--electric-blue);
    box-shadow: 0 0 10px var(--electric-blue);
  }
  
  .retro-button {
    @apply pixel-font text-xs px-4 py-2 uppercase cursor-pointer transition-all;
    background: linear-gradient(180deg, #9D00FF 0%, #6B00B3 100%);
    border: 2px solid var(--neon-pink);
    color: white;
    box-shadow: 0 0 10px var(--neon-pink);
  }
  
  .retro-button:hover {
    box-shadow: 0 0 20px var(--neon-pink);
    transform: translateY(-2px);
  }
}
```

### Step 6: Deploy to Vercel (5 min)

```bash
# Install Vercel CLI
npm i -g vercel

# Login
vercel login

# Deploy
vercel

# Follow prompts:
# Set up and deploy? Yes
# Which scope? Your account
# Link to existing project? No
# Project name? idea-hamster
# Directory? ./
# Override settings? No

# Deploy to production
vercel --prod
```

**Done!** Your basic setup is complete.

---

## 🎨 Next Steps

### Create Your First Component

Create `src/components/Leaderboard.tsx`:

```typescript
'use client';

import { useEffect, useState } from 'react';
import { supabase, Idea } from '@/lib/supabase';

export default function Leaderboard() {
  const [ideas, setIdeas] = useState<Idea[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchIdeas();
  }, []);

  async function fetchIdeas() {
    const { data, error } = await supabase
      .from('ideas')
      .select('*')
      .order('vote_count', { ascending: false })
      .limit(10);

    if (error) {
      console.error('Error fetching ideas:', error);
    } else {
      setIdeas(data || []);
    }
    setLoading(false);
  }

  if (loading) {
    return (
      <div className="text-center pixel-font text-sm">
        LOADING...
      </div>
    );
  }

  return (
    <div className="neon-border p-6">
      <h2 className="pixel-font text-center text-xl mb-6">
        ▓▓▓ LEADERBOARD ▓▓▓
      </h2>
      
      {ideas.map((idea, index) => (
        <div 
          key={idea.id} 
          className="border-b border-cyan-800 p-4 hover:bg-cyan-900/10 transition"
        >
          <div className="flex items-center gap-4">
            <div className="pixel-font text-2xl w-12 text-center">
              {index === 0 && '🥇'}
              {index === 1 && '🥈'}
              {index === 2 && '🥉'}
              {index > 2 && `#${index + 1}`}
            </div>
            
            <div className="flex-1">
              <h3 className="text-xl text-pink-400">{idea.title}</h3>
              <p className="text-sm text-gray-400 mt-1">{idea.description}</p>
              <div className="flex gap-2 mt-2">
                <span className="text-xs bg-purple-900/50 px-2 py-1 rounded">
                  {idea.category}
                </span>
                {idea.tags.map(tag => (
                  <span key={tag} className="text-xs bg-cyan-900/50 px-2 py-1 rounded">
                    {tag}
                  </span>
                ))}
              </div>
            </div>
            
            <div className="text-right">
              <div className="pixel-font text-2xl text-lime-400">
                ▲ {idea.vote_count}
              </div>
              {idea.vote_count >= 50 && (
                <div className="text-xs text-yellow-400 mt-1">
                  ELIGIBLE!
                </div>
              )}
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
```

### Update Homepage

Update `src/app/page.tsx`:

```typescript
import Leaderboard from '@/components/Leaderboard';

export default function Home() {
  return (
    <main className="min-h-screen p-8">
      <div className="max-w-5xl mx-auto">
        <header className="text-center mb-12">
          <h1 className="pixel-font text-4xl mb-4" style={{ color: 'var(--neon-pink)' }}>
            🐹 IDEA HAMSTER
          </h1>
          <p className="text-xl">
            Your Ideas. Our Hamster. Pure 90's Magic!
          </p>
        </header>
        
        <Leaderboard />
      </div>
    </main>
  );
}
```

### Test It!

```bash
npm run dev
```

Visit `http://localhost:3000` - you should see your retro leaderboard!

---

## 📝 Add Your First Idea (Test Data)

Go to Supabase SQL Editor and run:

```sql
INSERT INTO ideas (title, description, category, tags, submitted_by, vote_count)
VALUES 
  ('AI-Powered Plant Watering System', 'Smart IoT device that waters plants based on soil moisture and weather forecasts', 'fullstack', ARRAY['AI', 'IoT', 'Hardware'], 'GreenThumb99', 234),
  ('Retro Pixel Art Generator', 'Web app that converts photos into pixel art with customizable color palettes', 'frontend', ARRAY['Art', 'Canvas', 'Creative'], 'PixelMaster', 187),
  ('API Rate Limiter Dashboard', 'Real-time dashboard for monitoring and managing API rate limits across services', 'backend', ARRAY['API', 'DevTools', 'Monitoring'], 'DevOpsGuru', 142);
```

Refresh your app - the ideas should appear!

---

## 🎯 Week 1 Checklist

- [x] Set up Supabase
- [x] Initialize Next.js project
- [x] Create basic leaderboard
- [x] Deploy to Vercel
- [ ] Add submission form
- [ ] Add voting functionality
- [ ] Add email verification
- [ ] Add filtering/sorting
- [ ] Polish retro aesthetic
- [ ] Add animations

---

## 📚 Helpful Resources

- **Supabase Docs:** https://supabase.com/docs
- **Next.js Docs:** https://nextjs.org/docs
- **Tailwind CSS:** https://tailwindcss.com/docs
- **Framer Motion:** https://www.framer.com/motion/

---

## 🆘 Common Issues

### "Module not found" errors
```bash
rm -rf node_modules package-lock.json
npm install
```

### Supabase connection issues
- Check `.env.local` has correct values
- Ensure Supabase project is active (not paused)

### Styles not applying
- Restart dev server: `npm run dev`
- Clear browser cache

---

## 🚀 Next: Build Submission Form

See `PRD.md` for full feature specifications.

**You're ready to build!** 🎉

---

**Questions?** Open an issue or check the main PRD for details.