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

**Working Features:**
- ✅ Flutter app with gallery, camera, recommendations
- ✅ Go backend with PostgreSQL + S3 (LocalStack)
- ✅ Sample images load as local assets
- ✅ Docker image built (`davidnull/fitgenie:1.0`)
- ✅ Mobile-only branch created

**Known Limitations:**
- Requires manual IP configuration for WSL
- LocalStack S3 needs public endpoint for Flutter access
- No CI/CD pipeline
- No production cloud deployment

---

## Development Setup (WSL)

### 1. Backend IP Configuration

**Problem:** WSL changes IP on restart

**Solution:** Hardcode IP in `mobile/lib/services/api_service.dart`:

```dart
class ApiService {
  // Get your WSL IP:
  // $ ip route | grep default | awk '{print $3}'
  // Example: 172.21.48.1
  
  static String apiHost = '172.21.48.1';  // <-- Update this
  static String get baseUrl => 'http://$apiHost:8080/api/v1';
}
```

**Environments:**
- **Android Emulator:** `10.0.2.2` (host localhost)
- **iOS Simulator:** `localhost`
- **Physical Device:** Your PC's WiFi IP
- **Linux Desktop:** WSL IP (`172.21.48.1`)

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
