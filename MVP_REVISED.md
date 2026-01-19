# Idea Hamster - Revised MVP (Phased Automation) 🐹

**Version:** 2.0  
**Date:** January 2026  
**Status:** Ready for Development  

---

## 🎯 Core Concept (Simplified)

**Phase 1 MVP:** Beautiful 90's retro leaderboard for voting on app ideas  
**Phase 2:** Semi-automated build workflow with human oversight  
**Phase 3:** Full autonomous orchestration (future)  

---

## ✨ What Changed from Original PRD

### ❌ Removed from MVP
- Full autonomous AI building (too risky, too expensive)
- Real-time build logs streaming
- Automatic 48-hour build cycles
- Complex orchestration system

### ✅ Added to MVP
- **50 unique votes threshold** for build eligibility
- **Manual trigger** for build process
- **Phased workflow:** PRD → Tech Stack → Build → Deploy
- **Isolated dev containers** for safe building
- **Semi-automated approach** (AI-assisted, human-approved)

---

## 🚀 Phase 1: Community Validation MVP (Weeks 1-6)

### Goal
Prove that people want to submit ideas and vote on them, without the complexity of autonomous building.

### Features

#### 1. Idea Submission ✅
- Simple form with:
  - Title (80 chars)
  - Description (500 chars)
  - Category (Frontend/Backend/Full Stack)
  - Tags (max 5)
  - Submitter name (optional)
  
- Validation:
  - Profanity filter
  - Duplicate detection
  - Rate limiting (3 per user per week)

#### 2. Voting System ✅
- **Unique voter tracking:**
  - Email verification (lightweight, one-time)
  - One vote per email per idea
  - No full account needed
  
- **Vote display:**
  - Real-time count updates
  - Progress bars
  - Leaderboard ranking

#### 3. Leaderboard ✅
- Top 10 ideas displayed
- Sort by: Votes, Recent, Category
- Filter by category
- Expandable cards for details
- **Build Eligibility Badge:** Shows when idea hits 50 votes

#### 4. 90's Retro UI ✅
- Neon colors (pink, cyan, purple)
- Pixel fonts (Press Start 2P, VT323)
- CRT scanline effects
- Retro buttons and cards
- Animated hamster mascot
- Marquee text elements

#### 5. Build Queue ✅
- Separate page showing:
  - Ideas awaiting build (50+ votes)
  - Ideas in progress
  - Completed builds
  
- **Manual trigger button** (admin only)
- Status tracking

### Success Metrics (Phase 1)
- 500+ unique visitors in Month 1
- 50+ ideas submitted
- 1,000+ votes cast
- At least 3 ideas reach 50 votes
- 60%+ return visitor rate

---

## 🔧 Phase 2: Semi-Automated Building (Weeks 7-16)

### Goal
Build the top-voted ideas with AI assistance and human oversight, documenting the process for future automation.

### Workflow

#### Step 1: Manual Selection
**When an idea reaches 50 votes:**
1. Admin reviews idea for buildability
2. Clicks "Approve for Build" button
3. Idea moves to "Build Queue"

#### Step 2: AI-Generated PRD
**Using Claude API:**
```javascript
const prompt = `
You are a product manager. Create a detailed PRD for this app idea:

Title: ${idea.title}
Description: ${idea.description}
Category: ${idea.category}
Tags: ${idea.tags}

Include:
1. User stories
2. Core features (MVP scope)
3. Technical requirements
4. Success criteria
5. Estimated complexity (simple/medium/complex)
`;

const prd = await claudeAPI.generatePRD(prompt);
```

**Output:** Markdown PRD saved to GitHub repo

#### Step 3: AI-Generated Tech Stack
**Using Claude API:**
```javascript
const techStackPrompt = `
Based on this PRD, recommend:
1. Frontend framework
2. Backend framework (if needed)
3. Database choice
4. Hosting platform
5. Key dependencies
6. File structure

PRD: ${prd}
`;

const techStack = await claudeAPI.generateTechStack(techStackPrompt);
```

**Output:** `TECH_STACK.md` file

#### Step 4: Human Review & Approval
**Admin dashboard shows:**
- Generated PRD
- Proposed tech stack
- Estimated build time
- Estimated costs

**Admin can:**
- ✅ Approve as-is
- ✏️ Edit and refine
- ❌ Reject (back to queue)

#### Step 5: Containerized Build
**Using isolated dev containers:**

```yaml
# docker-compose.yml for each build
version: '3.8'
services:
  build-environment:
    image: node:20-alpine
    volumes:
      - ./workspace:/app
    working_dir: /app
    command: /bin/sh
    environment:
      - NODE_ENV=production
      - BUILD_ID=${idea.id}
    networks:
      - isolated-network
    
networks:
  isolated-network:
    driver: bridge
    internal: true  # No external access during build
```

**Build process:**
1. Spin up isolated container
2. AI generates code files (iterative approach)
3. Run tests inside container
4. Security scan (npm audit, OWASP checks)
5. Build succeeds → move to deployment
6. Build fails → human intervention

#### Step 6: Deployment
**Options based on app type:**

**Frontend Only:**
- Deploy to Vercel/Netlify
- Custom subdomain: `{idea-slug}.ideahamster.dev`

**Full Stack:**
- Frontend → Vercel
- Backend → Railway/Fly.io
- Database → Supabase (separate project per app)

**Backend API:**
- Deploy to Railway/Fly.io
- Provide API docs

#### Step 7: Post-Build
1. Update leaderboard with "Built!" badge
2. Add to archive with links
3. Send email to submitter
4. Tweet announcement
5. Add to showcase

### Technology Stack (Phase 2)

#### Build Orchestration
```javascript
// Recommended: BullMQ for job queue
import { Queue, Worker } from 'bullmq';

const buildQueue = new Queue('idea-builds', {
  connection: {
    host: 'redis-host',
    port: 6379
  }
});

// Add job when admin approves
await buildQueue.add('build-idea', {
  ideaId: idea.id,
  title: idea.title,
  description: idea.description
});

// Worker processes builds
const worker = new Worker('idea-builds', async (job) => {
  const { ideaId, title, description } = job.data;
  
  // Step 1: Generate PRD
  await job.updateProgress(10);
  const prd = await generatePRD(title, description);
  
  // Step 2: Generate tech stack
  await job.updateProgress(30);
  const techStack = await generateTechStack(prd);
  
  // Step 3: Wait for human approval
  await job.updateProgress(40);
  await waitForApproval(ideaId);
  
  // Step 4: Build in container
  await job.updateProgress(50);
  const buildResult = await buildInContainer(prd, techStack);
  
  // Step 5: Test
  await job.updateProgress(70);
  await runTests(buildResult);
  
  // Step 6: Security scan
  await job.updateProgress(80);
  await securityScan(buildResult);
  
  // Step 7: Deploy
  await job.updateProgress(90);
  const deployUrl = await deploy(buildResult);
  
  await job.updateProgress(100);
  return { success: true, url: deployUrl };
});
```

#### AI API Integration
```javascript
// api/ai/claude.ts
import Anthropic from '@anthropic-ai/sdk';

const anthropic = new Anthropic({
  apiKey: process.env.ANTHROPIC_API_KEY,
});

export async function generatePRD(title: string, description: string) {
  const message = await anthropic.messages.create({
    model: 'claude-sonnet-4-20250514',
    max_tokens: 8000,
    messages: [{
      role: 'user',
      content: PRD_PROMPT_TEMPLATE(title, description)
    }]
  });
  
  return message.content[0].text;
}

export async function generateCode(prd: string, fileName: string) {
  const message = await anthropic.messages.create({
    model: 'claude-sonnet-4-20250514',
    max_tokens: 4000,
    messages: [{
      role: 'user',
      content: CODE_GEN_PROMPT(prd, fileName)
    }]
  });
  
  return message.content[0].text;
}
```

#### Container Orchestration
```javascript
// utils/containerBuild.ts
import Docker from 'dockerode';
import { promisify } from 'util';

const docker = new Docker();

export async function buildInContainer(ideaId: string, files: FileMap) {
  // Create isolated network
  const network = await docker.createNetwork({
    Name: `build-${ideaId}`,
    Driver: 'bridge',
    Internal: true // No external access
  });
  
  // Create container
  const container = await docker.createContainer({
    Image: 'node:20-alpine',
    name: `build-${ideaId}`,
    WorkingDir: '/app',
    NetworkMode: network.id,
    HostConfig: {
      Memory: 512 * 1024 * 1024, // 512MB limit
      CpuQuota: 50000, // 50% CPU
      AutoRemove: true
    }
  });
  
  // Write files to container
  for (const [path, content] of Object.entries(files)) {
    await writeFileToContainer(container, path, content);
  }
  
  // Install dependencies
  await container.start();
  await execInContainer(container, 'npm install');
  
  // Run build
  const buildOutput = await execInContainer(container, 'npm run build');
  
  // Run tests
  const testOutput = await execInContainer(container, 'npm test');
  
  // Extract built files
  const builtFiles = await extractFromContainer(container, '/app/dist');
  
  // Cleanup
  await container.stop();
  await network.remove();
  
  return {
    success: true,
    files: builtFiles,
    logs: { build: buildOutput, test: testOutput }
  };
}
```

---

## 🔒 Security & Testing Strategy

### Security Measures

#### 1. Container Isolation
```yaml
# Strict resource limits
- Memory: 512MB max
- CPU: 50% quota
- Network: Isolated (no external access during build)
- Filesystem: Read-only except /app/workspace
- Time limit: 15 minutes max
```

#### 2. Code Scanning
**Before deployment, run:**
- `npm audit` - Check for known vulnerabilities
- `eslint` - Code quality and security rules
- `OWASP Dependency Check` - Security analysis
- Custom regex checks - No hardcoded secrets

#### 3. Runtime Security
**For deployed apps:**
- Rate limiting on all endpoints
- CORS restrictions
- Environment variable isolation (each app gets own)
- Subdomain isolation
- No direct database access (use Supabase RLS)

#### 4. Access Control
**Admin dashboard requires:**
- GitHub OAuth authentication
- Specific allow-list of admin emails
- 2FA enabled accounts only
- Audit log of all build approvals

### Testing Strategy

#### Automated Tests
```javascript
// tests/idea-build.test.ts
describe('Build Process', () => {
  it('should generate valid PRD', async () => {
    const prd = await generatePRD(mockIdea);
    expect(prd).toContain('## Features');
    expect(prd).toContain('## Technical Requirements');
  });
  
  it('should build in isolated container', async () => {
    const result = await buildInContainer(mockIdea.id, mockFiles);
    expect(result.success).toBe(true);
    expect(result.files).toHaveProperty('index.html');
  });
  
  it('should detect security issues', async () => {
    const scanResult = await securityScan(maliciousCode);
    expect(scanResult.issues).toHaveLength(1);
    expect(scanResult.severity).toBe('high');
  });
});
```

#### Manual Testing Checklist
Before marking build as "complete":
- [ ] App loads without errors
- [ ] Core functionality works
- [ ] Mobile responsive
- [ ] No console errors
- [ ] Security scan passed
- [ ] Performance acceptable (Lighthouse > 70)
- [ ] No exposed secrets

---

## 🌐 Hosting Architecture

### Main Platform (ideahamster.dev)
```
Platform: Vercel
├── Frontend: Next.js 15
├── API Routes: /api/*
├── Database: Supabase (PostgreSQL)
├── Cache: Upstash Redis
└── Storage: Supabase Storage
```

### Built Apps Hosting

#### Option 1: Vercel Projects (Recommended for MVP)
**Pros:**
- Easy deployment via API
- Free tier generous
- Custom domains (*.ideahamster.dev)
- Auto HTTPS
- Fast edge network

**Cons:**
- Costs scale with usage
- 100 projects limit on free tier

**Implementation:**
```javascript
import { VercelClient } from '@vercel/client';

const vercel = new VercelClient({ token: process.env.VERCEL_TOKEN });

async function deployToVercel(ideaSlug: string, files: FileMap) {
  const deployment = await vercel.createDeployment({
    name: ideaSlug,
    files: Object.entries(files).map(([path, content]) => ({
      file: path,
      data: content
    })),
    projectSettings: {
      framework: 'nextjs'
    }
  });
  
  // Set up custom domain
  await vercel.addDomain(deployment.id, `${ideaSlug}.ideahamster.dev`);
  
  return deployment.url;
}
```

#### Option 2: Railway/Fly.io (For Full Stack Apps)
**Pros:**
- Support databases, websockets
- Docker-based (consistency)
- Reasonable pricing ($5-10/app)

**Cons:**
- More complex setup
- Costs add up quickly

#### Option 3: Cloudflare Pages + Workers (Best for Scale)
**Pros:**
- Extremely cheap (potentially free)
- Fast global CDN
- Workers for backend logic
- D1 database, R2 storage

**Cons:**
- Different mental model
- Limitations on runtime

**Recommended for Phase 2:**
- **Frontend-only apps:** Vercel or Cloudflare Pages
- **Full-stack apps:** Railway (first 5), then evaluate based on usage
- **APIs only:** Fly.io

---

## 💰 Cost Analysis

### Phase 1 MVP Costs (Monthly)
| Service | Usage | Cost |
|---------|-------|------|
| Vercel (main platform) | Hobby tier | $0 |
| Supabase | Free tier | $0 |
| Domain (ideahamster.dev) | Annual/12 | $1 |
| **Total Phase 1** | | **~$1/mo** |

### Phase 2 Costs (Monthly, building 4 apps/month)
| Service | Usage | Cost |
|---------|-------|------|
| Vercel (main platform) | Pro tier (needed for team) | $20 |
| Supabase | Pro tier (more DB usage) | $25 |
| Anthropic Claude API | 4 builds × $15 avg | $60 |
| Vercel (deployed apps) | 4 projects | $0 |
| Railway (if needed) | 2 full-stack apps | $10 |
| Domain + DNS | | $1 |
| **Total Phase 2** | | **~$116/mo** |

### Cost Reduction Strategies
1. **Weekly builds instead of continuous:** 4/month vs 14/month = 71% savings
2. **Use Codex/GPT-4 for simpler prompts:** Cheaper than Claude for some tasks
3. **Community sponsors:** "This build sponsored by [Company]" ($50/build)
4. **Limit build complexity:** Cap at 2,000 lines of code per build
5. **Pause/archive old apps:** After 6 months, take down low-traffic apps

---

## 🎨 Updated UI/UX Flows

### Flow 1: Submit & Vote (Phase 1)
```
1. User visits ideahamster.dev
   └─> Sees leaderboard with current top ideas
   
2. Clicks "Submit Idea" (glowing neon button)
   └─> Fills form (title, description, category, tags)
   └─> Sees preview in retro card style
   └─> Submits (hamster animation plays)
   
3. Redirected to leaderboard
   └─> New idea appears at bottom
   └─> "Share your idea!" prompt with social buttons
   
4. User votes on other ideas
   └─> Clicks vote button (glitch effect)
   └─> Email verification modal (first time only)
   └─> Verifies email (one-time code)
   └─> Vote counted, leaderboard updates
   
5. Idea reaches 50 votes
   └─> "ELIGIBLE FOR BUILD!" badge appears
   └─> Confetti animation
   └─> Moves to "Build Queue" section
```

### Flow 2: Build Process (Phase 2)
```
Admin View:
1. Opens admin dashboard (auth required)
   └─> Sees "Build Queue" with eligible ideas
   └─> Reviews top idea details
   
2. Clicks "Generate PRD"
   └─> AI generates PRD (30 seconds)
   └─> Reviews PRD in markdown editor
   └─> Can edit/refine
   └─> Clicks "Approve PRD"
   
3. Clicks "Generate Tech Stack"
   └─> AI recommends stack
   └─> Shows estimated complexity
   └─> Approves or modifies
   
4. Clicks "Start Build"
   └─> Container spins up (visible logs)
   └─> AI generates code files (progress: 0-100%)
   └─> Tests run automatically
   └─> Security scan completes
   
5. Reviews build output
   └─> Previews app in iframe
   └─> Checks logs for errors
   └─> Runs manual tests
   └─> Approves deployment
   
6. App goes live
   └─> URL: {slug}.ideahamster.dev
   └─> Appears in "Built Apps" archive
   └─> Submitter notified via email
   └─> Announcement tweeted

Public View:
1. User checks "Build Queue" page
   └─> Sees their idea status: "PRD Generated ✅"
   └─> Watches progress bar: "Building... 45%"
   └─> Gets email: "Your idea is live! 🎉"
   └─> Visits deployed app
   └─> Shares on social media
```

---

## 📋 Phase 1 Development Checklist

### Week 1: Foundation
- [ ] Initialize Next.js 15 project
- [ ] Set up Supabase project
- [ ] Create database schema (ideas, votes, users)
- [ ] Deploy to Vercel
- [ ] Set up custom domain

### Week 2: Core Features
- [ ] Build idea submission form
- [ ] Implement vote system (with email verification)
- [ ] Create leaderboard component
- [ ] Add filtering and sorting
- [ ] Build "Build Queue" page (static for now)

### Week 3: Retro UI
- [ ] Design system (colors, fonts, components)
- [ ] Create retro component library
- [ ] Add CRT effects and scanlines
- [ ] Implement animated hamster mascot
- [ ] Polish animations and transitions

### Week 4: Testing & Polish
- [ ] Write unit tests
- [ ] E2E tests (Playwright)
- [ ] Mobile responsiveness
- [ ] Performance optimization
- [ ] SEO meta tags

### Week 5: Soft Launch
- [ ] Beta test with 20 users
- [ ] Fix critical bugs
- [ ] Pre-seed 10 quality ideas
- [ ] Create social media content
- [ ] Write launch blog post

### Week 6: Public Launch
- [ ] Launch on Product Hunt
- [ ] Post to Hacker News
- [ ] Share on Twitter/LinkedIn
- [ ] Monitor analytics
- [ ] Engage with community

---

## 📋 Phase 2 Development Checklist

### Week 7-8: AI Integration
- [ ] Set up Anthropic API
- [ ] Create PRD generation prompts
- [ ] Test with 5 sample ideas
- [ ] Build admin dashboard (auth)
- [ ] Create manual trigger UI

### Week 9-10: Container Setup
- [ ] Set up Docker environment
- [ ] Create isolated build containers
- [ ] Implement file writing/reading
- [ ] Add resource limits
- [ ] Test container builds

### Week 11-12: Security & Testing
- [ ] Implement security scanning
- [ ] Add automated tests to builds
- [ ] Create deployment scripts
- [ ] Set up monitoring (Sentry)
- [ ] Error handling and retries

### Week 13-14: Deployment Pipeline
- [ ] Integrate Vercel API
- [ ] Set up custom subdomains
- [ ] Create build status tracking
- [ ] Email notifications
- [ ] Social announcements automation

### Week 15-16: First Builds
- [ ] Build first idea end-to-end
- [ ] Document process
- [ ] Refine prompts
- [ ] Optimize costs
- [ ] Launch Phase 2 publicly

---

## 🎯 Success Criteria

### Phase 1 (Validation)
- ✅ 1,000+ unique visitors in first month
- ✅ 50+ ideas submitted
- ✅ 2,000+ votes cast
- ✅ 5+ ideas reach 50 votes
- ✅ 70%+ positive feedback

### Phase 2 (Building)
- ✅ 5 successful builds completed
- ✅ 80%+ build success rate
- ✅ < 2 hours per build (human time)
- ✅ < $20 per build (avg cost)
- ✅ No security incidents

---

## 🚀 Future: Phase 3 (Full Automation)

Once Phase 2 is proven, automate:
- PRD approval (AI evaluates buildability)
- Code review (AI checks quality)
- Deployment decision (auto-deploy if tests pass)
- Monitoring and alerts
- Scaling deployed apps based on usage

**Timeline:** 6-12 months after Phase 2 launch

---

## 🎉 Why This Approach Works

### ✅ Validates Core Hypothesis
You'll learn if people actually want to vote on ideas **before** spending time/money on automation.

### ✅ Manages Risk
- Phase 1 is simple and cheap
- Phase 2 adds complexity incrementally
- Can stop/pivot if Phase 1 fails

### ✅ Builds Community First
- Early adopters invest in the platform
- By the time you're building, there's an audience waiting

### ✅ Iterates on Prompts
- Manual process lets you perfect AI prompts
- Learn what works before automating

### ✅ Sustainable Costs
- Phase 1: ~$1/month
- Phase 2: ~$116/month (manageable)
- Can add sponsorships to offset

---

## 📊 Recommended Next Steps

1. **This Week:**
   - [ ] Review and approve this revised PRD
   - [ ] Create Figma mockups for leaderboard
   - [ ] Set up Supabase project
   - [ ] Register ideahamster.dev domain

2. **Next Week:**
   - [ ] Initialize Next.js project
   - [ ] Build basic leaderboard (no voting yet)
   - [ ] Deploy to Vercel
   - [ ] Share progress on Twitter

3. **Week 3-6:**
   - [ ] Complete Phase 1 features
   - [ ] Polish 90's aesthetic
   - [ ] Beta test with friends
   - [ ] Prepare launch content

4. **Week 7+:**
   - [ ] Start Phase 2 if Phase 1 successful
   - [ ] Build first idea manually to document process
   - [ ] Implement AI assistance incrementally

---

## 💡 Final Thoughts

This phased approach is **significantly better** than trying to build everything at once:

- **Lower risk:** Prove the concept before heavy investment
- **Faster to market:** Phase 1 can launch in 6 weeks
- **Lower costs:** Start at $1/month instead of $100+/month
- **Better UX:** Perfect the core loop before adding complexity
- **Community-driven:** Let users guide what gets built

**The 90's aesthetic makes Phase 1 special even without the AI building.** It's a fun, unique voting platform that stands alone.

Then Phase 2 adds the "magic" of semi-automated building once you have an engaged community.

---

**Let's build this! 🐹✨**

**Status:** Ready to Start Phase 1  
**Timeline:** Launch Phase 1 in 6 weeks  
**Budget:** ~$1/month to start  
