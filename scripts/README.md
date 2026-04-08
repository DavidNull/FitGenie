# scripts/

Scripts de utilidad para inicialización, desarrollo y despliegue.

## Archivos

### `init-db.sql`
Script de inicialización de PostgreSQL.

Se ejecuta automáticamente al iniciar el contenedor de PostgreSQL (vía volumen mount a `/docker-entrypoint-initdb.d/`).

**Hace:**
- Activa extensión `pgvector` para búsquedas vectoriales
- Añade columna `image_embedding` a tabla `clothing_items`
- Crea índice IVFFlat para búsqueda por similitud de coseno

### `localstack-init.sh`
Script de inicialización de LocalStack S3.

Se ejecuta cuando LocalStack está listo (vía volumen mount a `/etc/localstack/init/ready.d/`).

**Hace:**
- Crea bucket `fitgenie-images`
- Configura CORS para permitir uploads desde el frontend

### `prometheus.yml`
Configuración de Prometheus para scraping de métricas.

**Endpoints monitoreados:**
- Prometheus self-monitoring (`localhost:9090`)
- FitGenie API (`api:8080/metrics`)

## Uso manual

```bash
# Ejecutar init SQL manualmente
psql $DATABASE_URL -f scripts/init-db.sql

# Inicializar S3 en LocalStack corriendo
./scripts/localstack-init.sh
```

## Para añadir nuevos scripts

1. Crear archivo con shebang apropiado (`#!/bin/bash` o `#!/bin/sh`)
2. Añadir permisos de ejecución: `chmod +x scripts/nuevo-script.sh`
3. Documentar en este README
