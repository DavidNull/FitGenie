# FitGenie
![Logo](./mobile/assets/LOGO.png)
API para recomendar outfits de ropa. Analiza colores y estilo para sugerir combinaciones.

## Qué hace

- Usuario se identifica con `X-Device-ID` (sin password)
- Sube fotos de ropa → se guardan en S3
- Crea perfiles de estilo y color
- Recibe recomendaciones de outfits

## Stack

- Go 1.23 + Gin
- PostgreSQL + pgvector
- LocalStack S3 (para dev)
- Docker + Compose

## Estructura

```
cmd/server/       # Entrypoint
internal/         # Código privado
  api/handlers/   # HTTP handlers
  config/         # Config
  models/         # Modelos GORM
  repository/     # Acceso a datos
  services/       # Lógica de negocio
pkg/              # Librerías reutilizables
  database/
  logger/
  middleware/
  storage/
```

## Uso

```bash
# Docker (recomendado)
make docker-run
make docker-down

# Local
cp .env.example .env
make run
```

Servicios:
- API: http://localhost:8080
- DB: localhost:5432
- S3: http://localhost:4566

## Variables de entorno

Ver `.env.example`. Los defaults son para desarrollo local.

## Auth

Header `X-Device-ID` identifica al usuario. Cualquier string válido.

```bash
curl -H "X-Device-ID: mi-movil" http://localhost:8080/api/v1/users/me
```

## Endpoints

**Users**
- `GET /api/v1/users/me` - Usuario actual
- `POST /api/v1/users` - Crear usuario
- `GET/PUT/DELETE /api/v1/users/:id`

**Perfiles**
- `POST/GET /api/v1/users/:id/style-profile`
- `POST/GET /api/v1/users/:id/color-profile`

**Ropa**
- `POST /api/v1/clothing` - Añadir prenda
- `GET/PUT/DELETE /api/v1/clothing/:id`
- `GET /api/v1/clothing?user_id=xxx`

**Outfits**
- `POST /api/v1/outfits` - Crear outfit
- `GET /api/v1/users/:id/outfits`
- `POST /api/v1/users/:id/outfits/recommendations` - Recomendaciones

**Upload**
- `POST /api/v1/upload` - Subir imagen (multipart, campo 'image')
- `GET /api/v1/images/:path` - URL presignada

**Color**
- `GET /api/v1/color-theory/seasons`
- `GET /api/v1/color-theory/harmonies`
- `POST /api/v1/color-theory/analyze-harmony`

**Health**
- `GET /health`
- `GET /metrics` (Prometheus)

## Makefile

```bash
make build
make test
make run
make docker-run
make docker-down
```

## Falta hacer

- [ ] App móvil (Flutter/React Native)
- [ ] Extraer colores de fotos automáticamente
- [ ] Búsqueda por similitud con vectores

---

FitGenie — DavidNull

<p align="center">
  <img src="https://media0.giphy.com/media/v1.Y2lkPTc5MGI3NjExb284Znc1d2N2NWtuY2NxNTg0eWEwb3kwb2t3am50bHNpeWNqamptciZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/OccV3kjjhZnYnhJDCf/giphy.gif" width="120" /> 
</p>
