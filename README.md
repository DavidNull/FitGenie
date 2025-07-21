# FitGenie

FitGenie es una API para recomendaciones de outfits usando IA, color y estilo.

## Arquitectura

![Arquitectura de FitGenie](docs/img/FitGenie%20Arquitectura.png)

## Requisitos

- Go 1.23+
- Docker (para la base de datos PostgreSQL)

## Uso rápido

1. Clona el repo y entra a la carpeta:
   ```bash
   git clone <url>
   cd FitGenie
   ```

2. Instala dependencias:
   ```bash
   go mod tidy
   ```

3. Inicia la app (esto levanta la base de datos y la API):
   ```bash
   ./scripts/run-local.sh
   ```

4. API disponible en: [http://localhost:8080](http://localhost:8080)

## Endpoints principales

- `GET /api/v1/health` — Estado del servicio
- `POST /api/v1/users` — Crear usuario
- `POST /api/v1/users/:userId/clothing` — Añadir prenda
- `POST /api/v1/users/:userId/recommendations/outfits` — Recomendar outfits

## Configuración

Copia `.env.example` a `.env` y edítalo si lo necesitas

---

**FitGenie** — DavidNull 🐇
