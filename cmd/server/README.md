# cmd/server

Punto de entrada principal de la aplicación FitGenie.

## ¿Qué hace?

Este paquete inicializa y arranca el servidor HTTP de FitGenie:

1. **Carga configuración**: Lee variables de entorno vía `internal/config`
2. **Conecta a la base de datos**: Establece conexión PostgreSQL con `pkg/database`
3. **Ejecuta migraciones**: Crea/actualiza tablas automáticamente
4. **Configura el router**: Inicializa Gin con middleware de logging y métricas
5. **Expone métricas**: Endpoint `/metrics` para Prometheus
6. **Graceful shutdown**: Manejo elegante de señales SIGINT/SIGTERM

## Estructura

```
cmd/server/
└── main.go          # Único archivo - punto de entrada
```

## Cómo ejecutar

```bash
# Desde la raíz del proyecto
go run ./cmd/server/main.go

# O usando Make
make run

# O compilando primero
make build
./build/server
```

## Variables de entorno requeridas

| Variable | Requerida | Descripción |
|----------|-----------|-------------|
| `DATABASE_URL` | Sí | URL de PostgreSQL |
| `PORT` | No | Puerto HTTP (default: 8080) |
| `GIN_MODE` | No | debug/release (default: debug) |

## Ciclo de vida del servidor

```
main()
  ├── Cargar config
  ├── Conectar DB
  ├── Migrar modelos
  ├── Crear router Gin
  ├── Iniciar HTTP server
  └── Esperar señal de cierre (graceful shutdown)
```
