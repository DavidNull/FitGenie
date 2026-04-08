# pkg/

Librerías reutilizables de FitGenie - código genérico exportable.

## Principio

El directorio `pkg/` contiene código que podría ser usado por otros proyectos.
A diferencia de `internal/`, este código está diseñado para ser importado.

## Estructura

```
pkg/
├── database/          # Conexión y migraciones PostgreSQL
├── logger/            # Logger estructurado con slog
├── middleware/        # Middleware reutilizable para Gin
└── storage/           # Clientes de almacenamiento (S3)
```

## Paquetes

### `database/`

Abstracción de conexión a base de datos.

```go
db, err := database.NewConnection("postgres://...")
if err != nil {
    log.Fatal(err)
}
defer db.Close()

// Ejecutar migraciones
if err := db.Migrate(); err != nil {
    log.Fatal(err)
}
```

### `logger/`

Logger estructurado con soporte para JSON y texto plano.

```go
// Producción (JSON)
log := logger.NewLogger()
log.Info("servidor iniciado", "port", 8080)

// Desarrollo (texto plano)
log := logger.NewDevelopmentLogger()
log.Info("request recibido", "path", "/api/users")

// Con contexto
log.With("user_id", "123").Info("usuario autenticado")
```

### `middleware/`

Middleware reutilizable para framework Gin.

#### PrometheusMiddleware
Expone métricas HTTP:
- `http_requests_total` - contador con labels method/path/status
- `http_request_duration_seconds` - histograma de latencia
- `http_requests_in_flight` - gauge de requests activos

```go
router.Use(middleware.PrometheusMiddleware())
```

#### LoggerMiddleware
Logging de requests HTTP con formato estructurado:

```go
router.Use(middleware.LoggerMiddleware(log))
// Output: {"timestamp":"2024-...","method":"GET","path":"/api/users","status":200}
```

### `storage/`

Cliente S3 para almacenamiento de objetos.

```go
client, err := storage.NewS3Client(storage.S3Config{
    Endpoint:     "http://localhost:4566",  // LocalStack
    Region:       "us-east-1",
    Bucket:       "fitgenie-images",
    AccessKeyID:  "test",
    SecretAccessKey: "test",
    UsePathStyle: true,
})

// Subir archivo
err := client.Upload(ctx, "prendas/camisa.jpg", data, "image/jpeg")

// Generar URL presignada (válida 1 hora)
url, err := client.GetPresignedURL(ctx, "prendas/camisa.jpg", time.Hour)
```

## Uso en otros proyectos

```go
import (
    "github.com/davidnull/fitgenie/pkg/logger"
    "github.com/davidnull/fitgenie/pkg/database"
)
```
