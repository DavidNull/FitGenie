# .github/workflows/

Configuración de CI/CD con GitHub Actions.

## Workflows

### `ci.yml`
Pipeline de Integración Continua que se ejecuta en cada push/PR a `main` o `develop`.

#### Jobs

| Job | Descripción | Requisitos |
|-----|-------------|------------|
| **lint** | Ejecuta golangci-lint | Go 1.23 |
| **test** | Ejecuta tests con cobertura | Go 1.23 + PostgreSQL |
| **build** | Compila binario y build Docker | Necesita lint + test ✅ |

#### Flujo

```
Push/PR
   ├── lint (paralelo)
   ├── test (paralelo) ← PostgreSQL service
   └── build (secuencial, requiere lint+test)
       ├── Build binario
       └── Build Docker image
```

#### Tests con PostgreSQL

El job de tests levanta un servicio PostgreSQL con pgvector para tests de integración:

```yaml
services:
  postgres:
    image: pgvector/pgvector:pg15
    env:
      POSTGRES_DB: fitgenie_test
```

## Variables de entorno del workflow

| Variable | Valor | Uso |
|----------|-------|-----|
| `GO_VERSION` | `1.23` | Versión de Go |
| `GOLANGCI_LINT_VERSION` | `v1.59` | Versión del linter |

## Secrets requeridos

No se requieren secrets para el CI básico. Para despliegue automático a producción, añadir:

- `DOCKER_USERNAME` / `DOCKER_PASSWORD` - Para push a registry
- `KUBE_CONFIG` - Para despliegue a Kubernetes

## Añadir nuevos workflows

Crear archivo `.yml` en esta carpeta. Naming:
- `ci.yml` - Integración continua
- `cd.yml` - Despliegue continuo
- `release.yml` - Releases automáticos
