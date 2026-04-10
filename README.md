# FitGenie - API de Recomendación de Outfits

FitGenie es una API REST moderna para recomendaciones inteligentes de outfits de ropa, utilizando análisis de colorimetría, teoría del estilo e IA. Diseñada con arquitectura de microservicios, observabilidad completa y despliegue en contenedores.

### Cómo funciona el sistema

1. **Device se conecta**: App envía `X-Device-ID` en cada request
2. **Auto-registro**: Backend crea usuario automáticamente si es nuevo
3. **Sube prendas**: Usuario fotografía ropa → se sube a S3 → se guarda en PostgreSQL
4. **Análisis IA**: El sistema analiza colorimetría y estilo de cada prenda
5. **Recomendaciones**: Algoritmo genera outfits basados en perfiles + compatibilidad

### Flujo de datos

```
┌─────────┐    ┌──────────┐    ┌─────────────┐    ┌──────────┐
│  App    │───▶│   API    │───▶│  PostgreSQL │◀──▶│   S3     │
│ Flutter │◀───│  (Go)    │    │  + pgvector │    │  Images  │
└─────────┘    └──────────┘    └─────────────┘    └──────────┘
                     │
                     ▼
               ┌──────────┐
               │ Prometheus│
               │ Métricas  │
               └──────────┘
```

### Componentes Principales

| Componente | Tecnología | Descripción |
|------------|-----------|-------------|
| **API** | Go 1.23 + Gin | Servidor HTTP con middleware de métricas |
| **Base de Datos** | PostgreSQL + pgvector | Persistencia + búsqueda vectorial |
| **Storage** | LocalStack S3 | Almacenamiento de fotos de prendas |
| **Auth** | Device ID | Autenticación sin password por dispositivo |
| **Observabilidad** | Prometheus | Métricas de requests (2xx/4xx/5xx) y latencia |
| **Contenedores** | Docker + Compose | Orquestación local |

## 📁 Estructura del Proyecto

```
fitgenie/
├── cmd/server/              # Punto de entrada de la aplicación
├── internal/               # Código privado del proyecto
│   ├── api/handlers/      # Handlers HTTP (User, Clothing, Outfit, Color)
│   ├── config/            # Configuración por variables de entorno
│   ├── models/            # Modelos de datos GORM
│   ├── repository/        # Capa de acceso a datos (Repository Pattern)
│   └── services/          # Lógica de negocio (Color Theory, Style, AI)
├── pkg/                    # Librerías reutilizables
│   ├── database/          # Conexión y migraciones
│   ├── logger/            # Logger estructurado (slog)
│   ├── middleware/        # Middleware de Gin (Prometheus, Logger)
│   └── storage/           # Cliente S3 para almacenamiento de imágenes
├── scripts/                # Scripts de inicialización
├── .github/workflows/      # CI/CD con GitHub Actions
├── Dockerfile             # Build multi-etapa (distroless, ~15MB)
├── docker-compose.yml     # Stack completo local
└── Makefile               # Automatización de tareas
```

## �️ Herramientas Utilizadas

### Lenguajes y Frameworks
- **Go 1.23** - Lenguaje principal del backend
- **Gin** - Web framework HTTP (middleware, routing, validación)
- **GORM** - ORM para PostgreSQL (migraciones, modelos, queries)
- **go-colorful** - Manipulación y análisis de colores

### Base de Datos
- **PostgreSQL 15** - Base de datos relacional
- **pgvector** - Extensión para vectores (búsqueda por similitud futura)

### Storage y Cloud
- **LocalStack** - Simulación de AWS S3 para desarrollo local
- **AWS SDK v2** - Cliente S3 para Go

### Observabilidad
- **Prometheus** - Métricas y monitoreo
- **client_golang/prometheus** - Instrumentación HTTP

### DevOps
- **Docker** - Contenerización
- **Docker Compose** - Orquestación multi-servicio
- **Distroless** - Imagen minimalista para producción (~15MB)
- **GitHub Actions** - CI/CD con linting y tests

### Testing
- **testing (stdlib)** - Tests unitarios
- **testify** - Assertions y mocks

## � Inicio Rápido

### Requisitos

- [Docker](https://docs.docker.com/get-docker/) 20.10+
- [Docker Compose](https://docs.docker.com/compose/install/) 2.0+
- Go 1.23+ (solo para desarrollo local)

### Opción 1: Docker Compose (Recomendado)

```bash
# Clonar repositorio
git clone https://github.com/davidnull/fitgenie.git
cd fitgenie

# Iniciar stack completo (PostgreSQL + LocalStack S3 + API)
make docker-run

# Ver logs
make docker-logs

# Detener servicios
make docker-down
```

Servicios disponibles:
- **API**: http://localhost:8080
- **PostgreSQL**: localhost:5432
- **LocalStack S3**: http://localhost:4566
- **Prometheus** (opcional): http://localhost:9090

### Opción 2: Desarrollo Local

```bash
# Instalar dependencias
go mod download

# Configurar variables de entorno (ejemplos en .env.example)
cp .env.example .env
# Editar .env con tus credenciales reales

# Ejecutar
make run
```

## 🔧 Variables de Entorno

| Variable | Descripción | Default |
|----------|-------------|---------|
| `DATABASE_URL` | URL de conexión PostgreSQL | `postgres://fitgenie:fitgenie123@localhost:5432/fitgenie?sslmode=disable` |
| `PORT` | Puerto del servidor HTTP | `8080` |
| `GIN_MODE` | Modo de Gin (debug/release) | `debug` |
| `S3_ENDPOINT` | Endpoint de S3 (LocalStack/AWS) | `` |
| `S3_BUCKET` | Nombre del bucket S3 | `fitgenie-images` |
| `AWS_ACCESS_KEY_ID` | AWS Access Key | `` |
| `AWS_SECRET_ACCESS_KEY` | AWS Secret Key | `` |

> ⚠️ **Seguridad**: Los valores por defecto son para **desarrollo local únicamente**.  
> En producción, configura variables de entorno con secrets fuertes.  
> Ver `.env.example` para referencia.

## 🔐 Autenticación

FitGenie utiliza **Device ID** para autenticación sin password. Cada dispositivo envía un header `X-Device-ID` que identifica únicamente al usuario.

### Cómo funciona

1. **Primera vez**: El dispositivo envía cualquier string como `X-Device-ID`
2. **Backend**: Genera un UUID determinístico a partir del Device ID
3. **Usuario**: Se crea automáticamente un usuario asociado a ese Device ID
4. **Sesiones posteriores**: Mismo Device ID = mismo usuario

### Ejemplo

```bash
# Cualquier string funciona como device ID
curl -H "X-Device-ID: mi-movil-android" http://localhost:8080/api/v1/users/me

# Sin header, el backend genera uno nuevo
curl http://localhost:8080/api/v1/users/me
```

## 📡 API Endpoints

### Autenticación / Usuarios
```
GET    /api/v1/users/me                 # Obtener usuario actual (por Device ID)
POST   /api/v1/users                    # Crear usuario manual
GET    /api/v1/users                    # Listar usuarios
GET    /api/v1/users/:userId            # Obtener usuario por ID
PUT    /api/v1/users/:userId            # Actualizar usuario
DELETE /api/v1/users/:userId            # Eliminar usuario
```

### Perfiles
```
POST   /api/v1/users/:userId/style-profile   # Crear/actualizar perfil de estilo
GET    /api/v1/users/:userId/style-profile  # Obtener perfil de estilo
POST   /api/v1/users/:userId/color-profile   # Crear/actualizar perfil de color
GET    /api/v1/users/:userId/color-profile   # Obtener perfil de color
```

### Ropa
```
POST   /api/v1/clothing                    # Añadir prenda
GET    /api/v1/clothing/:id                 # Obtener prenda
PUT    /api/v1/clothing/:id                 # Actualizar prenda
DELETE /api/v1/clothing/:id                 # Eliminar prenda
GET    /api/v1/clothing?user_id=xxx        # Listar prendas de usuario
```

### Outfits
```
POST   /api/v1/outfits                      # Crear outfit
GET    /api/v1/outfits/:id                 # Obtener outfit
DELETE /api/v1/outfits/:id                 # Eliminar outfit
GET    /api/v1/users/:userId/outfits      # Listar outfits de usuario
POST   /api/v1/users/:userId/outfits/recommendations  # Generar recomendaciones IA
```

### Favoritos
```
POST   /api/v1/users/:userId/favorites/:outfitId   # Añadir a favoritos
DELETE /api/v1/users/:userId/favorites/:outfitId   # Quitar de favoritos
GET    /api/v1/users/:userId/favorites              # Listar favoritos
```

### Upload de Imágenes
```
POST   /api/v1/upload                      # Subir imagen (multipart/form-data, campo 'image')
GET    /api/v1/images/:path                # Obtener URL presignada de imagen
```

**Headers requeridos**: `X-Device-ID: tu-device-id`

**Ejemplo**:
```bash
# Subir foto de prenda
curl -X POST http://localhost:8080/api/v1/upload \
  -H "X-Device-ID: test-device-123" \
  -F "image=@/ruta/a/prenda.jpg"

# Respuesta
{
  "success": true,
  "file_path": "users/8669a828-.../53455e6b-....jpg",
  "url": "http://localstack:4566/fitgenie-images/...",
  "size": 17793
}
```

### Teoría del Color
```
GET    /api/v1/color-theory/seasons              # Listar estaciones de color
GET    /api/v1/color-theory/harmonies             # Listar armonías de color
POST   /api/v1/color-theory/analyze-harmony       # Analizar armonía de colores
POST   /api/v1/color-theory/recommendations       # Obtener colores recomendados
```

### Health & Métricas
```
GET    /health                          # Health check
GET    /metrics                         # Métricas Prometheus
```

## 🔨 Comandos Make

```bash
make build          # Compilar binario a ./build/server
make test           # Ejecutar tests con cobertura
make test-short     # Tests sin cobertura (más rápido)
make lint           # Ejecutar golangci-lint
make lint-fix       # Lint con auto-fix
make run            # Ejecutar localmente
make docker-run     # Iniciar con Docker Compose
make docker-down    # Detener Docker Compose
make docker-logs    # Ver logs de contenedores
make docker-clean   # Limpiar recursos Docker
make deps           # Descargar dependencias
make tidy           # go mod tidy + vendor
make ci             # Ejecutar todos los checks CI
```

## 🧪 Testing

```bash
# Ejecutar todos los tests
go test ./...

# Tests con cobertura
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## 📊 Observabilidad

### Métricas Prometheus Exponidas

| Métrica | Descripción |
|---------|-------------|
| `http_requests_total` | Total de requests HTTP (con labels: method, path, status) |
| `http_request_duration_seconds` | Latencia de requests en segundos |
| `http_requests_in_flight` | Requests en proceso actualmente |

### Configuración Prometheus Local

```bash
# Iniciar con stack de observabilidad
docker-compose --profile observability up -d

# Acceder a Prometheus
open http://localhost:9090
```

## 🏗️ Patrones de Arquitectura

### Repository Pattern
La capa `internal/repository/` abstrae el acceso a datos, permitiendo:
- **Testabilidad**: Mocks fáciles para tests unitarios
- **Desacoplamiento**: Handlers no conocen detalles de la DB
- **Flexibilidad**: Cambio de DB sin modificar lógica de negocio

### Dependency Injection
Todos los handlers reciben sus dependencias vía constructores:
```go
func NewUserHandler(repo repository.UserRepository, log *logger.Logger) *UserHandler
```

### Middleware Chain
```go
router.Use(gin.Recovery())
router.Use(middleware.LoggerMiddleware(log))
router.Use(middleware.PrometheusMiddleware())
```

## 🐳 Docker

### Build Multi-Etapa

El `Dockerfile` utiliza 2 etapas:
1. **Builder**: `golang:1.23-alpine` - Compila el binario estático
2. **Runtime**: `gcr.io/distroless/static:nonroot` - Imagen final minimalista (~15MB)

### Imagen Distroless
- Sin shell, sin utilidades del sistema
- Usuario no-root (`nonroot:nonroot`)
- Superficie de ataque mínima

```bash
# Build manual
docker build -t fitgenie:latest .

# Verificar tamaño
docker images fitgenie:latest
```

## 🗺️ Roadmap - Lo que Falta Implementar

### Frontend Móvil ❌ NO IMPLEMENTADO
FitGenie actualmente es **solo backend API**. Falta:

#### Aplicación Móvil (iOS/Android)
- **Tecnología recomendada**: Flutter o React Native
- **Funcionalidades necesarias**:
  - 📷 Cámara para fotografiar prendas
  - 🎨 Selector de colores para análisis de paleta personal
  - 👤 Gestión de perfiles de usuario
  - 👗 Visualizador de outfits recomendados
  - ❤️ Sistema de favoritos
  - 📊 Historial de outfits usados

#### Aplicación Web (Opcional)
- Panel de administración
- Dashboard de métricas de uso
- Gestión de catálogo de prendas

### Features Backend - Estado

| Feature | Estado | Descripción |
|---------|--------|-------------|
| **Device Auth** | ✅ Listo | Autenticación por Device ID, sin password |
| **Upload de imágenes** | ✅ Listo | Endpoint para subir fotos a S3 |
| **Análisis de imágenes** | 🔄 Pendiente | Extraer colores dominantes de fotos |
| **JWT Auth** | ❌ No requerido | Login/registro tradicional (innecesario con Device Auth) |
| **Embeddings vectoriales** | Media | Almacenar vectores de imágenes en pgvector |
| **Búsqueda por similitud** | Media | Encontrar prendas similares por imagen |
| **Notificaciones push** | Baja | Alertas de recomendaciones diarias |
| **Cache Redis** | Baja | Cache de recomendaciones frecuentes |

### Infraestructura

| Componente | Estado | Notas |
|------------|--------|-------|
| **Kubernetes** | ❌ No | Manifests K8s para despliegue en cloud |
| **Terraform** | ❌ No | IaC para AWS/GCP/Azure |
| **GitOps (ArgoCD)** | ❌ No | Despliegue automático |
| **Tracing (Jaeger)** | ❌ No | Distributed tracing |
| **Logs centralizados** | ❌ No | ELK/Loki stack |

## 📚 Documentación Adicional

- [`cmd/server/README.md`](cmd/server/README.md) - Punto de entrada
- [`internal/README.md`](internal/README.md) - Arquitectura interna
- [`pkg/README.md`](pkg/README.md) - Librerías compartidas
- [`scripts/README.md`](scripts/README.md) - Scripts de inicialización

## 🤝 Contribución

1. Fork el repositorio
2. Crea una rama (`git checkout -b feature/nueva-feature`)
3. Commit cambios (`git commit -am 'Añadir nueva feature'`)
4. Push a la rama (`git push origin feature/nueva-feature`)
5. Abre un Pull Request

## 📄 Licencia

MIT License - ver [LICENSE](LICENSE) para detalles

---

**FitGenie** — DavidNull 🐇 | Arquitectura de microservicios para moda inteligente
