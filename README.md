# FitGenie

App para gestionar tu armario y recibir sugerencias de outfits. Flutter + Go backend.

## Qué hace

- Sube fotos de tu ropa (categorías: camisetas, pantalones, calzado...)
- Guarda en la app con color y estilo
- Pide sugerencias para ocasiones (trabajo, cena, ocio)
- La app recomienda combinaciones según lo que tienes

## Stack técnico

| Capa | Tecnología |
|------|-----------|
| Mobile | Flutter (Provider) |
| Backend | Go (Gin + GORM) |
| BD | PostgreSQL |
| Imágenes | S3 (LocalStack para local) |
| Auth | JWT (listo para usar) |
| Docs | Swagger en `/swagger/index.html` |

## Probarlo en local

### 1. Clonar y levantar backend

```bash
git clone https://github.com/DavidNull/FitGenie.git
cd FitGenie

# Con Docker (PostgreSQL + LocalStack S3 + Go API)
docker compose up -d

# Verificar que está funcionando
curl http://localhost:8080/health
```

### 2. Levantar Flutter

```bash
cd mobile
flutter pub get

# Elige tu plataforma
flutter run -d linux        # Linux desktop
flutter run -d chrome       # Web
flutter run                 # Android conectado
```

**Nota:** La app detecta automáticamente la IP del backend. No necesitas configurar nada.

### 3. Usar la app

1. Toca "Usar imágenes de ejemplo" para cargar datos de prueba
2. Ve a "Armario" para ver tu ropa
3. Ve a "Recomendaciones", elige ocasión y temporada
4. Recibe sugerencias de outfits

## API

Documentación interactiva: `http://localhost:8080/swagger/index.html`

Endpoints principales:
- `POST /api/v1/users` - Crear usuario
- `GET /api/v1/users/me` - Mi perfil
- `POST /api/v1/upload` - Subir imagen
- `GET /api/v1/clothing` - Listar mi ropa
- `POST /api/v1/recommendations` - Pedir sugerencias

## Escalar a producción

### Opción rápida: Firebase

Para no mantener backend propio, migrar a Firebase:

```dart
// Reemplazar llamadas REST por Firebase
FirebaseFirestore.instance
  .collection('users').doc(uid)
  .collection('clothing').add(item)
```

Pros: Sin servidores, escalado automático, funciona offline  
Contras: Vendor lock-in, costes a escala

### Opción propia: Kubernetes (k3s)

Para auto-alojamiento en VPS barato (~5€/mes):

```bash
# Desplegar en k3s cluster
kubectl apply -f k8s/
```

Incluye: Ingress con TLS, PostgreSQL con volumen, 2 réplicas del API.

## Docker

Imagen publicada:

```bash
docker pull davidnull/fitgenie:latest
```

## Estructura del proyecto

```
FitGenie/
├── cmd/server/          # Entrypoint Go
├── internal/
│   ├── api/handlers/    # HTTP handlers
│   ├── services/        # Lógica de negocio
│   ├── repository/      # Acceso a BD
│   └── models/          # Structs
├── pkg/
│   ├── auth/            # JWT
│   ├── database/        # Conexión PostgreSQL
│   ├── middleware/      # Auth, logging
│   └── storage/         # S3 client
├── migrations/          # SQL migrations
├── mobile/              # Flutter app
│   ├── lib/
│   │   ├── screens/     # UI screens
│   │   ├── providers/   # Estado (Provider)
│   │   └── services/    # API client
│   └── test/
├── k8s/                 # Kubernetes manifests
└── .github/workflows/   # CI/CD
```

## CI/CD

GitHub Actions ejecuta en cada push:
- `go fmt`, `go vet`, tests
- `flutter analyze`, `dart format`
- Build Docker y push a Docker Hub

## TODO / Siguientes pasos

- [ ] Login real en Flutter (ahora usa device-id)
- [ ] Análisis automático de imágenes (detectar color/tipo)
- [ ] Tests unitarios backend
- [ ] Tests widget Flutter

## Licencia

MIT - Libre para usar y modificar.
