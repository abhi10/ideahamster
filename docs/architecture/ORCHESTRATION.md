# Idea Hamster - Orchestration & Deployment Architecture 🏗️

**Version:** 1.0  
**Date:** January 2026  
**Purpose:** Technical deep-dive on build orchestration, container management, and deployment pipeline  

---

## 🎯 System Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                       IDEA HAMSTER PLATFORM                     │
│                     (ideahamster.dev - Vercel)                  │
└──────────────────────┬──────────────────────────────────────────┘
                       │
                       │ Idea reaches 50 votes
                       ▼
┌─────────────────────────────────────────────────────────────────┐
│                    BUILD ORCHESTRATION LAYER                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │  BullMQ      │  │  Redis       │  │  PostgreSQL  │         │
│  │  Job Queue   │  │  Cache       │  │  Job State   │         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
└──────────────────────┬──────────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────────┐
│                    AI GENERATION PHASE                           │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Anthropic Claude API (Sonnet 4)                         │  │
│  │  • Generate PRD                                          │  │
│  │  • Generate Tech Stack                                   │  │
│  │  • Generate Code Files (iterative)                       │  │
│  └──────────────────────────────────────────────────────────┘  │
└──────────────────────┬──────────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────────┐
│                    ISOLATED BUILD ENVIRONMENT                    │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Docker Container (Isolated Network)                     │  │
│  │  • Install dependencies                                  │  │
│  │  • Run build command                                     │  │
│  │  • Execute tests                                         │  │
│  │  • Security scan                                         │  │
│  └──────────────────────────────────────────────────────────┘  │
└──────────────────────┬──────────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────────┐
│                    DEPLOYMENT TARGETS                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │  Vercel      │  │  Railway     │  │  Fly.io      │         │
│  │  (Frontend)  │  │  (Full Stack)│  │  (Backend)   │         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
└─────────────────────────────────────────────────────────────────┘
                       │
                       ▼
                ┌──────────────────┐
                │  Built App Lives │
                │  {slug}.ideaham  │
                │  ster.dev        │
                └──────────────────┘
```

---

## 🔧 Orchestration Components

### 1. Job Queue System (BullMQ)

**Why BullMQ?**
- Built on Redis (fast, reliable)
- Retries and backoff strategies
- Priority queues
- Progress tracking
- Job events and hooks
- UI dashboard available (Bull Board)

**Setup:**

```typescript
// lib/queue/buildQueue.ts
import { Queue, Worker, QueueEvents } from 'bullmq';
import { Redis } from 'ioredis';

// Redis connection
const connection = new Redis({
  host: process.env.REDIS_HOST,
  port: parseInt(process.env.REDIS_PORT || '6379'),
  password: process.env.REDIS_PASSWORD,
  maxRetriesPerRequest: null,
  enableReadyCheck: false
});

// Build queue
export const buildQueue = new Queue('idea-builds', {
  connection,
  defaultJobOptions: {
    attempts: 3,
    backoff: {
      type: 'exponential',
      delay: 5000
    },
    removeOnComplete: 100, // Keep last 100 completed jobs
    removeOnFail: 200      // Keep last 200 failed jobs
  }
});

// Queue events for monitoring
const queueEvents = new QueueEvents('idea-builds', { connection });

queueEvents.on('completed', ({ jobId, returnvalue }) => {
  console.log(`Build ${jobId} completed successfully`);
  // Update database, send notifications
});

queueEvents.on('failed', ({ jobId, failedReason }) => {
  console.error(`Build ${jobId} failed: ${failedReason}`);
  // Alert admin, update status
});

queueEvents.on('progress', ({ jobId, data }) => {
  console.log(`Build ${jobId} progress: ${data}%`);
  // Update UI in real-time via websocket
});
```

**Job Data Structure:**

```typescript
interface BuildJobData {
  ideaId: string;
  title: string;
  description: string;
  category: 'frontend' | 'backend' | 'fullstack';
  tags: string[];
  submitter: string;
  voteCount: number;
  approvedBy: string; // Admin email
  approvedAt: Date;
}

// Add job to queue
export async function queueBuildJob(ideaId: string) {
  const idea = await db.idea.findUnique({ where: { id: ideaId } });
  
  const job = await buildQueue.add('build-idea', {
    ideaId: idea.id,
    title: idea.title,
    description: idea.description,
    category: idea.category,
    tags: idea.tags,
    submitter: idea.submittedBy,
    voteCount: idea.voteCount,
    approvedBy: 'admin@ideahamster.dev',
    approvedAt: new Date()
  }, {
    jobId: `build-${ideaId}`,
    priority: idea.voteCount // Higher votes = higher priority
  });
  
  return job.id;
}
```

---

### 2. Build Worker

**Worker Implementation:**

```typescript
// workers/buildWorker.ts
import { Worker, Job } from 'bullmq';
import { generatePRD, generateTechStack, generateCodeFiles } from '@/lib/ai/claude';
import { buildInContainer } from '@/lib/docker/containerBuild';
import { securityScan } from '@/lib/security/scanner';
import { deployApp } from '@/lib/deploy/deployer';
import { db } from '@/lib/db';

const buildWorker = new Worker('idea-builds', async (job: Job<BuildJobData>) => {
  const { ideaId, title, description, category, tags } = job.data;
  
  try {
    // Phase 1: Generate PRD (10% progress)
    await job.updateProgress(5);
    await updateBuildStatus(ideaId, 'GENERATING_PRD');
    
    const prd = await generatePRD({
      title,
      description,
      category,
      tags
    });
    
    await saveBuildArtifact(ideaId, 'PRD.md', prd);
    await job.updateProgress(10);
    
    // Phase 2: Generate Tech Stack (20% progress)
    await updateBuildStatus(ideaId, 'GENERATING_TECH_STACK');
    
    const techStack = await generateTechStack(prd);
    await saveBuildArtifact(ideaId, 'TECH_STACK.md', techStack);
    await job.updateProgress(20);
    
    // Phase 3: Wait for Human Approval (30% progress)
    await updateBuildStatus(ideaId, 'AWAITING_APPROVAL');
    await job.updateProgress(30);
    
    // Poll for approval (timeout after 72 hours)
    const approved = await waitForApproval(ideaId, 72 * 60 * 60 * 1000);
    if (!approved) {
      throw new Error('Build approval timeout');
    }
    
    // Phase 4: Generate Code Files (40-60% progress)
    await updateBuildStatus(ideaId, 'GENERATING_CODE');
    
    const files = await generateCodeFiles(prd, techStack, (progress) => {
      job.updateProgress(40 + (progress * 0.2)); // 40-60%
    });
    
    await job.updateProgress(60);
    
    // Phase 5: Build in Container (70% progress)
    await updateBuildStatus(ideaId, 'BUILDING');
    
    const buildResult = await buildInContainer(ideaId, files);
    if (!buildResult.success) {
      throw new Error(`Build failed: ${buildResult.error}`);
    }
    
    await job.updateProgress(70);
    
    // Phase 6: Security Scan (80% progress)
    await updateBuildStatus(ideaId, 'SECURITY_SCAN');
    
    const scanResult = await securityScan(buildResult.outputPath);
    if (scanResult.criticalIssues > 0) {
      throw new Error(`Security issues found: ${scanResult.issues.join(', ')}`);
    }
    
    await job.updateProgress(80);
    
    // Phase 7: Deploy (90% progress)
    await updateBuildStatus(ideaId, 'DEPLOYING');
    
    const deployUrl = await deployApp({
      ideaId,
      files: buildResult.outputPath,
      category,
      envVars: {} // Empty for now, can add per-app config later
    });
    
    await job.updateProgress(90);
    
    // Phase 8: Finalize (100% progress)
    await updateBuildStatus(ideaId, 'COMPLETED');
    await db.build.update({
      where: { ideaId },
      data: {
        deployUrl,
        completedAt: new Date(),
        status: 'COMPLETED'
      }
    });
    
    await job.updateProgress(100);
    
    // Post-build actions
    await sendSuccessNotifications(ideaId, deployUrl);
    await tweetAnnouncement(title, deployUrl);
    
    return {
      success: true,
      ideaId,
      deployUrl,
      duration: Date.now() - job.timestamp
    };
    
  } catch (error) {
    await updateBuildStatus(ideaId, 'FAILED');
    await db.build.update({
      where: { ideaId },
      data: {
        status: 'FAILED',
        errorMessage: error.message,
        failedAt: new Date()
      }
    });
    
    throw error; // BullMQ will handle retry logic
  }
}, {
  connection,
  concurrency: 2, // Max 2 builds at once
  limiter: {
    max: 10,      // Max 10 jobs
    duration: 1000 // Per second
  }
});

// Helper functions
async function waitForApproval(ideaId: string, timeout: number): Promise<boolean> {
  const startTime = Date.now();
  
  while (Date.now() - startTime < timeout) {
    const build = await db.build.findUnique({ where: { ideaId } });
    if (build?.adminApproved) {
      return true;
    }
    await new Promise(resolve => setTimeout(resolve, 10000)); // Check every 10s
  }
  
  return false;
}

async function updateBuildStatus(ideaId: string, status: string) {
  await db.build.update({
    where: { ideaId },
    data: { phase: status, updatedAt: new Date() }
  });
}

async function saveBuildArtifact(ideaId: string, fileName: string, content: string) {
  await db.buildArtifact.create({
    data: {
      buildId: ideaId,
      fileName,
      content,
      createdAt: new Date()
    }
  });
}
```

---

## 🐳 Docker Container Orchestration

### Container Security Configuration

```typescript
// lib/docker/containerBuild.ts
import Docker from 'dockerode';
import path from 'path';
import fs from 'fs-extra';
import { createHash } from 'crypto';

const docker = new Docker();

export interface BuildFiles {
  [path: string]: string;
}

export interface BuildResult {
  success: boolean;
  outputPath?: string;
  logs?: {
    install: string;
    build: string;
    test: string;
  };
  error?: string;
}

export async function buildInContainer(
  ideaId: string,
  files: BuildFiles
): Promise<BuildResult> {
  const buildId = createHash('md5').update(ideaId).digest('hex');
  const workspacePath = path.join('/tmp', `build-${buildId}`);
  
  // Prepare workspace
  await fs.ensureDir(workspacePath);
  
  // Write all files to workspace
  for (const [filePath, content] of Object.entries(files)) {
    const fullPath = path.join(workspacePath, filePath);
    await fs.ensureDir(path.dirname(fullPath));
    await fs.writeFile(fullPath, content, 'utf-8');
  }
  
  try {
    // Create isolated network
    const network = await docker.createNetwork({
      Name: `build-network-${buildId}`,
      Driver: 'bridge',
      Internal: true, // No external internet access
      EnableIPv6: false
    });
    
    // Create container with strict limits
    const container = await docker.createContainer({
      Image: 'node:20-alpine',
      name: `build-${buildId}`,
      WorkingDir: '/app',
      Cmd: ['/bin/sh'],
      Tty: false,
      AttachStdout: true,
      AttachStderr: true,
      HostConfig: {
        // Resource limits
        Memory: 512 * 1024 * 1024,      // 512MB RAM
        MemorySwap: 512 * 1024 * 1024,  // No swap
        CpuQuota: 50000,                 // 50% CPU
        CpuPeriod: 100000,
        PidsLimit: 100,                  // Max 100 processes
        
        // Security
        ReadonlyRootfs: false, // Need write access for npm install
        SecurityOpt: ['no-new-privileges'],
        CapDrop: ['ALL'],      // Drop all capabilities
        CapAdd: ['CHOWN', 'SETUID', 'SETGID'], // Only needed for npm
        
        // Network isolation
        NetworkMode: `build-network-${buildId}`,
        
        // Filesystem
        Binds: [
          `${workspacePath}:/app:rw`
        ],
        
        // Time limit (handled by timeout below)
        AutoRemove: true
      },
      Env: [
        'NODE_ENV=production',
        'CI=true',
        'NO_UPDATE_NOTIFIER=true'
      ]
    });
    
    await container.start();
    
    // Set overall timeout
    const timeoutPromise = new Promise<never>((_, reject) => {
      setTimeout(() => {
        container.kill().catch(() => {});
        reject(new Error('Build timeout (15 minutes exceeded)'));
      }, 15 * 60 * 1000); // 15 minutes
    });
    
    // Run build steps
    const buildPromise = (async () => {
      // Step 1: Install dependencies
      const installLogs = await execInContainer(container, 'npm install --production');
      
      // Step 2: Run build (if package.json has build script)
      let buildLogs = '';
      try {
        buildLogs = await execInContainer(container, 'npm run build');
      } catch (e) {
        // Build script might not exist, that's OK
        buildLogs = 'No build script found';
      }
      
      // Step 3: Run tests
      let testLogs = '';
      try {
        testLogs = await execInContainer(container, 'npm test');
      } catch (e) {
        // Tests might fail, log but don't block
        testLogs = `Tests failed: ${e.message}`;
      }
      
      return {
        install: installLogs,
        build: buildLogs,
        test: testLogs
      };
    })();
    
    const logs = await Promise.race([buildPromise, timeoutPromise]);
    
    // Extract built files
    const outputPath = path.join('/tmp', `output-${buildId}`);
    await fs.ensureDir(outputPath);
    
    // Copy build output or source files
    const stream = await container.getArchive({ path: '/app' });
    const tarStream = fs.createWriteStream(path.join(outputPath, 'app.tar'));
    stream.pipe(tarStream);
    
    await new Promise((resolve, reject) => {
      tarStream.on('finish', resolve);
      tarStream.on('error', reject);
    });
    
    // Extract tar
    const tar = require('tar');
    await tar.extract({
      file: path.join(outputPath, 'app.tar'),
      cwd: outputPath
    });
    
    // Cleanup
    await container.stop();
    await network.remove();
    await fs.remove(workspacePath);
    
    return {
      success: true,
      outputPath,
      logs
    };
    
  } catch (error) {
    // Cleanup on error
    try {
      const containers = await docker.listContainers({ 
        all: true,
        filters: { name: [`build-${buildId}`] }
      });
      for (const containerInfo of containers) {
        const container = docker.getContainer(containerInfo.Id);
        await container.remove({ force: true });
      }
      
      const networks = await docker.listNetworks({
        filters: { name: [`build-network-${buildId}`] }
      });
      for (const networkInfo of networks) {
        const network = docker.getNetwork(networkInfo.Id);
        await network.remove();
      }
    } catch (cleanupError) {
      console.error('Cleanup error:', cleanupError);
    }
    
    await fs.remove(workspacePath);
    
    return {
      success: false,
      error: error.message
    };
  }
}

async function execInContainer(
  container: Docker.Container,
  command: string
): Promise<string> {
  const exec = await container.exec({
    Cmd: ['/bin/sh', '-c', command],
    AttachStdout: true,
    AttachStderr: true
  });
  
  const stream = await exec.start({ Detach: false });
  
  let output = '';
  stream.on('data', (chunk) => {
    output += chunk.toString();
  });
  
  await new Promise((resolve, reject) => {
    stream.on('end', resolve);
    stream.on('error', reject);
  });
  
  const inspectResult = await exec.inspect();
  if (inspectResult.ExitCode !== 0) {
    throw new Error(`Command failed: ${command}\n${output}`);
  }
  
  return output;
}
```

---

## 🔒 Security & Testing

### Security Scanner

```typescript
// lib/security/scanner.ts
import { exec } from 'child_process';
import { promisify } from 'util';
import path from 'path';
import fs from 'fs-extra';

const execAsync = promisify(exec);

export interface SecurityScanResult {
  passed: boolean;
  criticalIssues: number;
  highIssues: number;
  mediumIssues: number;
  lowIssues: number;
  issues: Array<{
    severity: 'critical' | 'high' | 'medium' | 'low';
    type: string;
    description: string;
    file?: string;
    line?: number;
  }>;
}

export async function securityScan(buildPath: string): Promise<SecurityScanResult> {
  const results: SecurityScanResult = {
    passed: true,
    criticalIssues: 0,
    highIssues: 0,
    mediumIssues: 0,
    lowIssues: 0,
    issues: []
  };
  
  // 1. npm audit
  try {
    const { stdout } = await execAsync('npm audit --json', { 
      cwd: buildPath,
      timeout: 60000 
    });
    
    const auditResult = JSON.parse(stdout);
    if (auditResult.metadata) {
      const { vulnerabilities } = auditResult.metadata;
      results.criticalIssues += vulnerabilities.critical || 0;
      results.highIssues += vulnerabilities.high || 0;
      results.mediumIssues += vulnerabilities.moderate || 0;
      results.lowIssues += vulnerabilities.low || 0;
      
      // Add specific issues
      for (const [name, issue] of Object.entries(auditResult.vulnerabilities || {})) {
        results.issues.push({
          severity: issue.severity,
          type: 'dependency',
          description: `${name}: ${issue.via[0]?.title || 'Vulnerability'}`,
          file: 'package.json'
        });
      }
    }
  } catch (e) {
    // npm audit might exit with non-zero if vulnerabilities found
    console.warn('npm audit error:', e.message);
  }
  
  // 2. Check for exposed secrets
  const secretPatterns = [
    { pattern: /(sk|pk)_live_[a-zA-Z0-9]{24,}/, name: 'Stripe API Key' },
    { pattern: /AKIA[0-9A-Z]{16}/, name: 'AWS Access Key' },
    { pattern: /AIza[0-9A-Za-z\\-_]{35}/, name: 'Google API Key' },
    { pattern: /sk-[a-zA-Z0-9]{48}/, name: 'OpenAI API Key' },
    { pattern: /ghp_[a-zA-Z0-9]{36}/, name: 'GitHub Token' },
  ];
  
  const files = await getAllFiles(buildPath);
  for (const file of files) {
    const content = await fs.readFile(file, 'utf-8');
    const lines = content.split('\n');
    
    for (let i = 0; i < lines.length; i++) {
      for (const { pattern, name } of secretPatterns) {
        if (pattern.test(lines[i])) {
          results.issues.push({
            severity: 'critical',
            type: 'exposed-secret',
            description: `Potential ${name} found in code`,
            file: path.relative(buildPath, file),
            line: i + 1
          });
          results.criticalIssues++;
        }
      }
    }
  }
  
  // 3. Check for dangerous patterns
  const dangerousPatterns = [
    { pattern: /eval\(/, name: 'eval() usage', severity: 'high' },
    { pattern: /dangerouslySetInnerHTML/, name: 'dangerouslySetInnerHTML', severity: 'medium' },
    { pattern: /localStorage\.setItem\([^,]+,\s*password/i, name: 'Password in localStorage', severity: 'critical' },
  ];
  
  for (const file of files) {
    if (!file.endsWith('.js') && !file.endsWith('.ts') && !file.endsWith('.jsx') && !file.endsWith('.tsx')) {
      continue;
    }
    
    const content = await fs.readFile(file, 'utf-8');
    const lines = content.split('\n');
    
    for (let i = 0; i < lines.length; i++) {
      for (const { pattern, name, severity } of dangerousPatterns) {
        if (pattern.test(lines[i])) {
          results.issues.push({
            severity: severity as any,
            type: 'code-quality',
            description: `Dangerous pattern: ${name}`,
            file: path.relative(buildPath, file),
            line: i + 1
          });
          
          if (severity === 'critical') results.criticalIssues++;
          else if (severity === 'high') results.highIssues++;
          else if (severity === 'medium') results.mediumIssues++;
        }
      }
    }
  }
  
  // Determine if passed (no critical issues)
  results.passed = results.criticalIssues === 0;
  
  return results;
}

async function getAllFiles(dir: string): Promise<string[]> {
  const entries = await fs.readdir(dir, { withFileTypes: true });
  const files = await Promise.all(
    entries.map(async (entry) => {
      const fullPath = path.join(dir, entry.name);
      
      // Skip node_modules, .git, etc.
      if (entry.name === 'node_modules' || entry.name === '.git' || entry.name === 'dist') {
        return [];
      }
      
      if (entry.isDirectory()) {
        return getAllFiles(fullPath);
      } else {
        return [fullPath];
      }
    })
  );
  
  return files.flat();
}
```

---

## 🚀 Deployment Layer

### Multi-Platform Deployer

```typescript
// lib/deploy/deployer.ts
import { deployToVercel } from './vercel';
import { deployToRailway } from './railway';
import { deployToFly } from './fly';

export interface DeployConfig {
  ideaId: string;
  files: string; // Path to built files
  category: 'frontend' | 'backend' | 'fullstack';
  envVars?: Record<string, string>;
}

export async function deployApp(config: DeployConfig): Promise<string> {
  const { category, ideaId } = config;
  const slug = await generateSlug(ideaId);
  
  // Route to appropriate platform based on category
  switch (category) {
    case 'frontend':
      return await deployToVercel({ ...config, slug });
      
    case 'fullstack':
      return await deployToRailway({ ...config, slug });
      
    case 'backend':
      return await deployToFly({ ...config, slug });
      
    default:
      throw new Error(`Unknown category: ${category}`);
  }
}

async function generateSlug(ideaId: string): Promise<string> {
  const idea = await db.idea.findUnique({ where: { id: ideaId } });
  if (!idea) throw new Error('Idea not found');
  
  // Convert title to slug
  const slug = idea.title
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-|-$/g, '')
    .substring(0, 32);
  
  return slug;
}
```

### Vercel Deployment

```typescript
// lib/deploy/vercel.ts
import { createDeployment } from '@vercel/client';
import fs from 'fs-extra';
import path from 'path';

export async function deployToVercel(config: DeployConfig & { slug: string }): Promise<string> {
  const { files: filesPath, slug } = config;
  
  // Read all files
  const files = await readFilesRecursively(filesPath);
  
  // Create deployment
  const deployment = await createDeployment({
    name: slug,
    files: files.map(f => ({
      file: f.path,
      data: f.content
    })),
    projectSettings: {
      framework: detectFramework(filesPath),
      buildCommand: 'npm run build',
      outputDirectory: detectOutputDir(filesPath),
      installCommand: 'npm install'
    },
    target: 'production'
  }, {
    token: process.env.VERCEL_TOKEN!,
    teamId: process.env.VERCEL_TEAM_ID
  });
  
  // Wait for deployment to be ready
  await waitForDeployment(deployment.id);
  
  // Set custom domain
  const customDomain = `${slug}.ideahamster.dev`;
  await addCustomDomain(deployment.id, customDomain);
  
  return `https://${customDomain}`;
}

function detectFramework(filesPath: string): string {
  const packageJsonPath = path.join(filesPath, 'package.json');
  if (fs.existsSync(packageJsonPath)) {
    const packageJson = fs.readJsonSync(packageJsonPath);
    
    if (packageJson.dependencies?.next) return 'nextjs';
    if (packageJson.dependencies?.react) return 'create-react-app';
    if (packageJson.dependencies?.vue) return 'vue';
    if (packageJson.dependencies?.svelte) return 'svelte';
  }
  
  return 'static'; // Default to static site
}

function detectOutputDir(filesPath: string): string {
  // Check for common output directories
  const commonDirs = ['dist', 'build', 'out', 'public'];
  for (const dir of commonDirs) {
    if (fs.existsSync(path.join(filesPath, dir))) {
      return dir;
    }
  }
  return '.'; // Current directory
}
```

---

## 📊 Monitoring & Observability

### Build Metrics

```typescript
// lib/monitoring/metrics.ts
import { db } from '@/lib/db';

export async function trackBuildMetrics(buildId: string, metrics: {
  duration: number;
  tokensUsed: number;
  apiCost: number;
  success: boolean;
}) {
  await db.buildMetrics.create({
    data: {
      buildId,
      duration: metrics.duration,
      tokensUsed: metrics.tokensUsed,
      apiCost: metrics.apiCost,
      success: metrics.success,
      timestamp: new Date()
    }
  });
  
  // Calculate rolling averages
  const recentBuilds = await db.buildMetrics.findMany({
    where: {
      timestamp: {
        gte: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000) // Last 30 days
      }
    }
  });
  
  const avgDuration = recentBuilds.reduce((sum, b) => sum + b.duration, 0) / recentBuilds.length;
  const avgCost = recentBuilds.reduce((sum, b) => sum + b.apiCost, 0) / recentBuilds.length;
  const successRate = recentBuilds.filter(b => b.success).length / recentBuilds.length;
  
  console.log('Build Metrics (30 days):');
  console.log(`  Avg Duration: ${(avgDuration / 1000 / 60).toFixed(2)} minutes`);
  console.log(`  Avg Cost: $${avgCost.toFixed(2)}`);
  console.log(`  Success Rate: ${(successRate * 100).toFixed(1)}%`);
}