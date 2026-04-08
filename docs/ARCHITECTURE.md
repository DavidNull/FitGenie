# Arquitectura de FitGenie

Documentación técnica detallada de la arquitectura del sistema.

## Visión General

FitGenie sigue una arquitectura de **microservicios monolítico** (modular monolith) con clara separación de responsabilidades.

```
┌─────────────────────────────────────────────────────────────┐
│                        Clientes                             │
│  (App Móvil iOS/Android - NO IMPLEMENTADO)                 │
└──────────────────────┬──────────────────────────────────────┘
                       │ HTTPS
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                      API Gateway                            │
│  • Rate limiting                                            │
│  • Authentication (JWT - pendiente)                        │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                    FitGenie API (Go)                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐    │
│  │   Users     │  │  Clothing   │  │     Outfits     │    │
│  │   Module    │  │   Module    │  │     Module      │    │
│  └─────────────┘  └─────────────┘  └─────────────────┘    │
│                                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐    │
│  │    Color    │  │    Style    │  │       AI        │    │
│  │   Theory    │  │   Service   │  │   Recommender   │    │
│  └─────────────┘  └─────────────┘  └─────────────────┘    │
└──────────────────────┬──────────────────────────────────────┘
                       │
           ┌───────────┼───────────┐
           │           │           │
           ▼           ▼           ▼
┌──────────────┐ ┌──────────┐ ┌──────────────┐
│  PostgreSQL  │ │ LocalStack│ │  Prometheus  │
│  + pgvector  │ │    S3     │ │  (métricas)  │
└──────────────┘ └──────────┘ └──────────────┘
```

## Capas de la Arquitectura

### 1. Presentation Layer (Transport)

**Responsabilidad**: Manejo de requests HTTP, serialización JSON

**Implementación**:
- Framework: Gin v1.9
- Middleware: Logger, Recovery, Prometheus
- Validación: Go-Playground Validator

**Flujo de request**:
```
HTTP Request
    ↓
Gin Router (routes.go)
    ↓
Middleware Chain
    ↓
Handler (handlers/*.go)
    ↓
Service/Repository
```

### 2. Business Logic Layer (Domain)

**Responsabilidad**: Reglas de negocio, algoritmos de recomendación

**Servicios**:

#### ColorTheoryService
- Análisis de colorimetría personal
- Cálculo de armonías de color
- Determinación de estación de color

#### StyleService
- Categorización de estilos de moda
- Guías por tipo de cuerpo
- Recomendaciones por ocasión

#### AIService
- Algoritmos de recomendación de outfits
- Scoring de compatibilidad
- Análisis de coherencia de estilo

### 3. Data Access Layer (Repository)

**Responsabilidad**: Persistencia y recuperación de datos

**Patrón**: Repository Pattern con interfaces

```go
type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
    // ...
}
```

**Implementaciones**:
- GORM v2 con driver PostgreSQL
- Soporte para context.Context (timeouts)
- Preload de relaciones

### 4. Infrastructure Layer

**Responsabilidad**: Conexiones externas, logging, storage

**Componentes**:
- `pkg/database`: Pool de conexiones PostgreSQL
- `pkg/storage`: Cliente S3 (LocalStack/AWS)
- `pkg/logger`: Logger estructurado (slog)

## Modelo de Datos

```
┌─────────────┐       ┌─────────────────┐       ┌─────────────────┐
│    User     │───────│  StyleProfile   │       │  ColorProfile   │
├─────────────┤       ├─────────────────┤       ├─────────────────┤
│ id (uuid)   │       │ id (uuid)       │       │ id (uuid)       │
│ email       │       │ user_id (FK)    │       │ user_id (FK)    │
│ name        │       │ preferred_styles│       │ color_season    │
│ created_at  │       │ body_type       │       │ skin_tone       │
└─────────────┘       │ lifestyle       │       │ undertone       │
                      │ occasion        │       │ favorite_colors │
                      └─────────────────┘       └─────────────────┘
                              │
                              │
┌─────────────┐       ┌─────┴─────────────┐       ┌─────────────────┐
│ClothingItem │       │      Outfit       │◄──────│OutfitRecommenda-│
├─────────────┤       ├───────────────────┤       │      tion       │
│ id (uuid)   │       │ id (uuid)         │       ├─────────────────┤
│ user_id(FK) │──────►│ user_id (FK)      │       │ id (uuid)       │
│ name        │       │ name              │       │ outfit_id (FK)  │
│ category    │       │ description       │       │ confidence      │
│ primary_col-│       │ style             │       │ reasoning       │
│ or          │       │ color_harmony_    │       │ viewed          │
│ image_embed-│       │ score             │       │ accepted        │
│ ding(vector)│       │ overall_score     │       └─────────────────┘
└─────────────┘       └───────────────────┘
```

### Campos Vectoriales (pgvector)

La extensión `pgvector` permite almacenar embeddings de imágenes para búsqueda por similitud:

```sql
ALTER TABLE clothing_items 
ADD COLUMN image_embedding vector(512);

CREATE INDEX idx_clothing_image_embedding 
ON clothing_items USING ivfflat (image_embedding vector_cosine_ops);
```

**Casos de uso futuros**:
- Búsqueda: "Encontrar prendas similares a esta imagen"
- Recomendación basada en embeddings visuales
- Clustering automático de prendas por similitud visual

## Flujo de Recomendación de Outfits

```
1. Usuario solicita recomendación
   POST /api/v1/users/{id}/outfits/recommendations
   Body: { occasion: "casual", season: "summer" }

2. Handler recibe request
   outfit_handler.go:GetOutfitRecommendations()

3. Recuperar datos del usuario
   ├── userRepo.GetByID() → Datos básicos
   ├── userRepo.GetStyleProfile() → Preferencias de estilo
   └── userRepo.GetColorProfile() → Paleta de colores personal

4. Recuperar prendas disponibles
   clothingRepo.ListByUser() → []ClothingItem

5. AI Service genera recomendaciones
   aiService.GenerateOutfitRecommendations()
   ├── Filtrar por ocasión/estación
   ├── Calcular compatibilidad de colores
   ├── Calcular coherencia de estilo
   └── Score y ordenar opciones

6. Persistir recomendaciones
   outfitRepo.Create() para cada recomendación

7. Responder al cliente
   JSON: { recommendations: [...], total: N }
```

## Decisiones de Diseño

### ¿Por qué Repository Pattern?

**Problema**: Los handlers accediendo directamente a GORM son difíciles de testear.

**Solución**: Interfaces de repository permiten mocks:

```go
// En tests
mockRepo := &mockUserRepository{
    users: []models.User{...}
}
handler := NewUserHandler(mockRepo, log)
```

### ¿Por qué Dependency Injection?

**Beneficios**:
1. **Testabilidad**: Inyectar mocks en tests
2. **Flexibilidad**: Cambiar implementaciones sin modificar handlers
3. **Claridad**: Dependencias explícitas en constructores

### ¿Por qué Gin sobre Gorilla Mux?

| Característica | Gin | Gorilla Mux |
|----------------|-----|-------------|
| Performance | ⭐⭐⭐ | ⭐⭐ |
| Middleware | Built-in | Manual |
| Validación | Binding automático | Manual |
| Context | *gin.Context | http.ResponseWriter |
| JSON | Serialización automática | Manual |

### ¿Por qué Distroless?

| Aspecto | Alpine (~20MB) | Distroless (~15MB) |
|---------|---------------|-------------------|
| Tamaño | 20MB | 15MB |
| Shell | sh | ❌ Ninguno |
| Herramientas | apk, curl, etc | ❌ Ninguna |
| Superficie ataque | Media | Mínima |
| Debug | Fácil | Difícil |

Para producción, distroless es más seguro.

## Escalabilidad

### Horizontal Scaling

La API es **stateless**, permitiendo escalar horizontalmente:

```yaml
# docker-compose.yml con 3 réplicas
services:
  api:
    deploy:
      replicas: 3
    environment:
      DATABASE_URL: postgres://...  # Misma DB
```

### Base de Datos

Para alta disponibilidad:
- PostgreSQL Primary-Replica (streaming replication)
- Connection pooling (PgBouncer)
- Read replicas para queries pesados

### Cache (Pendiente)

Redis para cachear:
- Recomendaciones populares
- Perfiles de usuario frecuentes
- Resultados de análisis de color

## Seguridad (Pendiente)

### Implementar

1. **Autenticación JWT**
   - Login con email/password
   - Refresh tokens
   - Middleware de protección de rutas

2. **Autorización RBAC**
   - Roles: user, admin
   - Usuarios solo acceden a sus datos

3. **Validación de Inputs**
   - Sanitización de strings
   - Limitar tamaño de payloads
   - Rate limiting por IP/usuario

4. **Storage Seguro**
   - URLs presignadas con expiración
   - Validación de tipos MIME
   - Escaneo antivirus (ClamAV)

## Observabilidad

### Métricas Prometheus

```promql
# Request rate por endpoint
rate(http_requests_total[5m])

# Latencia p95
histogram_quantile(0.95, 
  rate(http_request_duration_seconds_bucket[5m]))

# Error rate
rate(http_requests_total{status=~"5.."}[5m])
```

### Logs Estructurados

```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "level": "INFO",
  "msg": "request completed",
  "method": "POST",
  "path": "/api/v1/users",
  "status": 201,
  "duration_ms": 45,
  "user_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

## Referencias

- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Gin Web Framework](https://gin-gonic.com/)
- [GORM Documentation](https://gorm.io/)
- [pgvector](https://github.com/pgvector/pgvector)
