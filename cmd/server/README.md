# cmd/server

Punto de entrada de la aplicación.

## Uso

```bash
go run ./cmd/server/main.go
make run
```

## Env vars

- `DATABASE_URL` - PostgreSQL connection string
- `PORT` - default 8080
- `GIN_MODE` - debug/release
