# FitGenie - Complete Development Guide

This README contains all steps for development, production migration, and deployment options.

## Table of Contents

1. [Current State](#current-state)
2. [Development Setup (WSL)](#development-setup-wsl)
3. [Production Migration](#production-migration)
4. [Deployment Options](#deployment-options)
5. [Next Steps](#next-steps)

---

## Current State

**Professional Features Implemented:**
- ✅ **Flutter app** with gallery, camera, AI recommendations
- ✅ **Go backend** with clean architecture (handlers, services, repositories)
- ✅ **PostgreSQL** with **golang-migrate** SQL migrations
- ✅ **S3-compatible storage** (LocalStack) with presigned URLs
- ✅ **Auto-detect backend IP** - no more hardcoding WSL IP
- ✅ **Swagger/OpenAPI docs** at `/swagger/index.html`
- ✅ **CI/CD pipeline** with GitHub Actions (go fmt, flutter analyze, Docker build)
- ✅ **JWT Authentication** with access/refresh tokens
- ✅ **Docker Hub image** `davidnull/fitgenie:1.0`
- ✅ **Mobile-only branch** for quick testing

**Architecture Highlights:**
```
┌─────────────┐     HTTP/JSON     ┌──────────────────────────────┐
│  Flutter    │ ═══════════════► │ Go Backend (Gin + GORM)      │
│   Mobile    │ ◄═══════════════ │ ├─ REST API with Swagger     │
└─────────────┘                 │ ├─ JWT Authentication        │
                                │ ├─ PostgreSQL + Migrations   │
                                │ ├─ S3 Image Storage          │
                                │ └─ AI Outfit Recommendations │
                                └──────────────────────────────┘
```

---

## Development Setup (WSL)

### 1. Backend IP Configuration (Auto-Detect)

✅ **FIXED:** Flutter app now **auto-detects** the backend IP!

The app probes common addresses on startup:
- `localhost`
- `10.0.2.2` (Android Emulator)
- `172.17.0.1` (Docker bridge)
- `172.21.48.1` (WSL common)

**No manual configuration needed!**

If you need to override, edit `mobile/lib/services/api_service.dart`:
```dart
// Auto-detection (default)
static String get apiHost => _detectedHost ?? 'localhost';

// Manual override (if needed)
static String apiHost = '172.21.48.1';
```

### 2. Start Backend

```bash
# Option A: Docker Compose (with source code)
cd /home/david/FitGenie
docker compose up -d

# Option B: Docker Hub image only
docker pull davidnull/fitgenie:1.0
# Needs external postgres + localstack
```

### 3. Start Flutter

```bash
cd /home/david/FitGenie/mobile
flutter pub get
flutter run -d linux  # or chrome, android, ios
```

---

## Professional Features

### Database Migrations (golang-migrate)

**Industry standard** for database schema management:

```
migrations/
├── 000001_init_schema.up.sql      # Create tables
├── 000001_init_schema.down.sql    # Rollback
├── 000002_add_openai_config.up.sql
└── 000002_add_openai_config.down.sql
```

**Run migrations:**
```bash
# Auto-runs on backend startup
go run cmd/server/main.go

# Manual migration (using CLI tool)
docker run -v $(pwd)/migrations:/migrations \
  migrate/migrate -path=/migrations -database=postgres://... up
```

### API Documentation (Swagger/OpenAPI)

Interactive API docs available at: `http://localhost:8080/swagger/index.html`

**Features:**
- Browse all API endpoints
- Test endpoints directly from browser
- See request/response schemas
- Authentication with JWT

### CI/CD Pipeline (GitHub Actions)

**Checks on every commit:**
- ✅ `go fmt` - Go code formatting
- ✅ `go vet` - Static analysis
- ✅ `go test` - Unit tests
- ✅ `flutter analyze` - Dart linting
- ✅ `flutter build apk` - Android build
- ✅ Docker build test
- ✅ Security scan (Trivy)

**Auto-deploy to Docker Hub on main branch:**
```yaml
# .github/workflows/ci.yml
tags:
  - davidnull/fitgenie:latest
  - davidnull/fitgenie:${{ github.sha }}
```

### JWT Authentication

**Secure token-based auth** with access and refresh tokens:

```go
// pkg/auth/jwt.go
type TokenPair struct {
    AccessToken  string  // 24h expiry
    RefreshToken string  // 7d expiry
    ExpiresIn    int64
}
```

**Protected routes:**
```go
authMiddleware := middleware.NewAuthMiddleware(jwtService, log)

// Require valid JWT
cleanRouter := v1.Group("/outfits")
cleanRouter.Use(authMiddleware.RequireAuth())
cleanRouter.GET("", outfitHandler.ListOutfits)
```

**Flutter integration:**
```dart
// Store tokens
await storage.write(key: 'access_token', value: tokenPair.accessToken);

// Add to requests
headers['Authorization'] = 'Bearer $accessToken';
```

---

## Production Migration

### Phase 1: Firebase (Easiest)

**Goal:** Replace local backend with Firebase services

**Why Firebase:**
- No server maintenance
- Built-in auth, database, storage
- Global CDN for images
- Works offline

**Migration Steps:**

1. **Add Firebase to Flutter:**
```bash
flutter pub add firebase_core firebase_auth cloud_firestore firebase_storage
```

2. **Replace API calls with Firestore:**
```dart
// Before (REST API)
final items = await apiService.getClothingItems();

// After (Firestore)
final snapshot = await FirebaseFirestore.instance
  .collection('users')
  .doc(userId)
  .collection('clothing')
  .get();
```

3. **Replace S3 with Firebase Storage:**
```dart
// Before (S3 upload)
final url = await apiService.uploadImage(file);

// After (Firebase Storage)
final ref = FirebaseStorage.instance.ref('clothing/$fileName');
await ref.putFile(file);
final url = await ref.getDownloadURL();
```

4. **Add Firebase Auth:**
```dart
// Anonymous auth for quick start
final user = await FirebaseAuth.instance.signInAnonymously();
```

**Firebase Architecture:**
```
┌─────────────┐      Firebase      ┌─────────────┐
│  Flutter    │  ═══════════════►  │  Firestore  │
│   (app)     │  ═══════════════►  │   Storage   │
└─────────────┘                    │    Auth     │
                                   └─────────────┘
```

**Pros:**
- No backend code to maintain
- Scales automatically
- Works from anywhere (no local network needed)

**Cons:**
- Vendor lock-in
- Costs money at scale
- Limited backend logic

---

### Phase 2: Kubernetes (k3s)

**Goal:** Self-hosted production on k3s cluster

**Why k3s:**
- Lightweight Kubernetes
- Runs on Raspberry Pi or cheap VPS
- Full control over infrastructure

**Architecture:**
```
┌─────────────┐                    ┌─────────────────┐
│   Flutter   │  ════════════════► │   k3s Cluster   │
│    (app)    │                    │ ┌─────────────┐ │
└─────────────┘                    │ │  Ingress    │ │
                                   │ │    (TLS)    │ │
                                   │ └──────┬──────┘ │
                                   │        │        │
                                   │ ┌──────┴──────┐ │
                                   │ │  FitGenie   │ │
                                   │ │   (Go API)  │ │
                                   │ └──────┬──────┘ │
                                   │        │        │
                                   │ ┌──────┴──────┐ │
                                   │ │  Postgres   │ │
                                   │ │   (PVC)     │ │
                                   │ └─────────────┘ │
                                   └─────────────────┘
```

**Files to Create:**

1. `k8s/namespace.yaml`
2. `k8s/postgres-deployment.yaml`
3. `k8s/postgres-service.yaml`
4. `k8s/postgres-pvc.yaml`
5. `k8s/fitgenie-deployment.yaml`
6. `k8s/fitgenie-service.yaml`
7. `k8s/ingress.yaml`

**Example k3s Deployment:**

```yaml
# k8s/fitgenie-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: fitgenie-api
  namespace: fitgenie
spec:
  replicas: 2
  selector:
    matchLabels:
      app: fitgenie-api
  template:
    metadata:
      labels:
        app: fitgenie-api
    spec:
      containers:
      - name: api
        image: davidnull/fitgenie:1.0
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          value: "postgres://fitgenie:password@postgres:5432/fitgenie"
        - name: S3_ENDPOINT
          value: "https://s3.amazonaws.com"  # Real S3
        - name: S3_BUCKET
          value: "fitgenie-production"
```

**Deploy:**
```bash
# On k3s master node
git clone https://github.com/DavidNull/FitGenie.git
cd FitGenie/k8s
kubectl apply -f namespace.yaml
kubectl apply -f postgres-pvc.yaml
kubectl apply -f postgres-deployment.yaml
kubectl apply -f postgres-service.yaml
kubectl apply -f fitgenie-deployment.yaml
kubectl apply -f fitgenie-service.yaml
kubectl apply -f ingress.yaml
```

**Pros:**
- Full control
- Can run on cheap hardware
- Portable between cloud providers

**Cons:**
- More complex to manage
- Need to handle TLS, backups, monitoring

---

## Deployment Options Summary

| Option | Complexity | Cost | Best For |
|--------|-----------|------|----------|
| **Local (WSL)** | Low | Free | Development |
| **Firebase** | Low-Medium | Pay-as-you-go | MVP, quick launch |
| **k3s (VPS)** | Medium | $5-20/month | Production, learning |
| **AWS/GCP** | High | $50+/month | Enterprise |

---

## Next Steps

### Immediate (Before LinkedIn Post)

1. ✅ **Fix sample images** - Done
2. ✅ **Push Docker image** - Built (need to login to push)
3. ✅ **Create mobile branch** - Done
4. ⏳ **Test end-to-end flow** - Import samples → Gallery → Recommendations

### Short Term (Next 2 Weeks)

1. **Add GitHub Actions CI/CD**
```yaml
# .github/workflows/build.yml
name: Build and Push
on: [push]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Build Docker
        run: docker build -t davidnull/fitgenie:${{ github.sha }} .
      - name: Push to Docker Hub
        run: docker push davidnull/fitgenie:${{ github.sha }}
```

2. **Add Flutter tests**
```bash
cd mobile
flutter test
```

3. **Write blog post** about the architecture

### Long Term (Next 3 Months)

1. **Migrate to Firebase** for production
2. **Add user authentication** (Google, Apple)
3. **Implement AI recommendations** with OpenAI API
4. **Deploy to k3s** for portfolio demo
5. **Add monitoring** (Prometheus, Grafana)

---

## Useful Commands

```bash
# Docker
make docker-run        # Start all services
docker logs fitgenie-api -f  # Watch API logs

# Flutter
flutter run -d linux   # Desktop mode
flutter build apk      # Android release
flutter build ios      # iOS release

# k3s
kubectl get pods -n fitgenie
kubectl logs deployment/fitgenie-api -n fitgenie
```

---

## Resources

- [Flutter Documentation](https://docs.flutter.dev)
- [Firebase Flutter Setup](https://firebase.google.com/docs/flutter/setup)
- [k3s Quick Start](https://docs.k3s.io/quick-start)
- [Kubernetes Basics](https://kubernetes.io/docs/tutorials/kubernetes-basics/)

---

## Docker Hub Image

```bash
# Pull and run (once published)
docker pull davidnull/fitgenie:1.0
docker run -p 8080:8080 davidnull/fitgenie:1.0
```

**Branches:**
- `main` - Full project with backend
- `mobile-only` - Just Flutter app

---

## License

MIT License - See LICENSE file
