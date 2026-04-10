# scripts/

Scripts de inicialización.

## Archivos

- `init-db.sql` - Crea extensión pgvector en PostgreSQL
- `localstack-init.sh` - Crea bucket S3 en LocalStack
- `prometheus.yml` - Config de Prometheus

## Uso

```bash
psql $DATABASE_URL -f scripts/init-db.sql
./scripts/localstack-init.sh
```
