# Railway Deployment Guide

This document captures learnings from deploying Idea Hamster to Railway.

## Overview

- **Platform**: Railway (railway.app)
- **Plan**: Basic ($5/month)
- **Build**: Dockerfile-based
- **Domain**: ideahamster.dev (custom domain via Porkbun)

## Configuration Files

### railway.toml
```toml
[build]
builder = "dockerfile"
dockerfilePath = "Dockerfile"

[deploy]
healthcheckPath = "/health"
healthcheckTimeout = 100
restartPolicyType = "on_failure"
restartPolicyMaxRetries = 3
```

### Dockerfile (Key Points)
```dockerfile
# Multi-stage build for small image size
FROM golang:1.24-alpine AS builder
# ... build steps ...

FROM alpine:3.19
# Railway injects PORT at runtime - no need to set defaults
CMD ["./server"]
```

## Critical Learnings

### 1. PORT Environment Variable is Required

**Issue**: App deployed but returned 502 errors.

**Cause**: Railway needs to know which port your app listens on.

**Solution**: Add `PORT=3001` in Railway Variables tab.

```go
// Your app must read PORT from environment
port := os.Getenv("PORT")
if port == "" {
    port = "3001"
}
server.ListenAndServe(":" + port)
```

**Key Point**: Don't hardcode ENV PORT in Dockerfile - let Railway inject it.

### 2. Static Files Must Be in Git

**Issue**: CSS returned 404 in production, site had no styling.

**Cause**: `output.css` (Tailwind compiled CSS) was in `.gitignore`.

**Solution**: Remove from `.gitignore` and commit the file.

```gitignore
# Before (broken)
web/static/css/output.css

# After (working)
# web/static/css/output.css  # Commented out - file is committed
```

**Key Point**: Docker builds don't have Node.js to regenerate Tailwind CSS.

### 3. SSL Certificate Provisioning Takes Time

**Issue**: "Certificate Authority is validating challenges" stuck for 10+ minutes.

**Cause**: Let's Encrypt ACME challenge was pending.

**Solution**: Delete and re-add the custom domain in Railway to restart the process.

**Timeline**: SSL provisioning typically takes 2-10 minutes after DNS propagates.

### 4. ALIAS Records for Root Domain

**Issue**: Can't use CNAME for root domain (ideahamster.dev).

**Solution**: Use ALIAS record (Porkbun supports this).

```
Type: ALIAS
Host: (blank for root)
Value: xxxx.up.railway.app
```

### 5. Different DNS Targets for Each Domain

**Observation**: Railway gives different DNS targets when you add custom domains.

```
ideahamster.dev     → 0hn5ka4p.up.railway.app
www.ideahamster.dev → tv2kelln.up.railway.app
```

This is normal - each domain gets its own routing target.

## DNS Configuration (Porkbun)

| Type  | Host | Value                     | TTL |
|-------|------|---------------------------|-----|
| ALIAS | @    | xxxx.up.railway.app       | 600 |
| CNAME | www  | yyyy.up.railway.app       | 600 |

**Note**: Values change if you remove and re-add domains in Railway.

## Health Check Endpoint

Railway uses the health check to determine deployment success:

```go
// internal/handlers/health.go
func HandleHealth(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "status":    "healthy",
        "timestamp": time.Now().UTC().Format(time.RFC3339),
        "version":   "1.0.0",
    })
}
```

Configured in `railway.toml`:
```toml
healthcheckPath = "/health"
healthcheckTimeout = 100
```

## Troubleshooting

### 502 Bad Gateway
1. Check PORT variable is set in Railway
2. Verify app logs for startup errors
3. Ensure health check endpoint returns 200

### SSL Certificate Errors
1. Wait 10 minutes for initial provisioning
2. If stuck, delete and re-add the custom domain
3. Verify DNS records match Railway's requirements

### Static Files 404
1. Check files are committed to git (not in .gitignore)
2. Verify Dockerfile copies static folder
3. Check static file route in main.go

### CSS Not Loading
1. Ensure `output.css` is in git
2. Or add Tailwind build step to Dockerfile (requires Node.js)

## Costs

| Plan    | Cost     | Custom Domains |
|---------|----------|----------------|
| Hobby   | Free     | 1              |
| Basic   | $5/month | Unlimited      |

## URLs

- **Production**: https://ideahamster.dev
- **Backup**: https://ideahamster-production.up.railway.app
- **Health**: https://ideahamster.dev/health

## Future Improvements

1. Add PostgreSQL addon in Railway
2. Set up staging environment
3. Configure CI/CD for automatic deployments
4. Add monitoring/alerting
