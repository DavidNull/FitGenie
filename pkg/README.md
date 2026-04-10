# pkg/

Librerías reutilizables.

## Estructura

```
pkg/
├── database/      # PostgreSQL connection
├── logger/        # Structured logging
├── middleware/    # Gin middleware
└── storage/       # S3 client
```

## Uso

```go
// Database
db, _ := database.NewConnection("postgres://...")
defer db.Close()
db.Migrate()

// Logger
log := logger.NewLogger()
log.Info("msg", "key", "value")

// Middleware
router.Use(middleware.PrometheusMiddleware())
router.Use(middleware.LoggerMiddleware(log))

// S3
client, _ := storage.NewS3Client(storage.S3Config{
    Endpoint: "http://localhost:4566",
    Region:   "us-east-1",
    Bucket:   "fitgenie-images",
})
```
