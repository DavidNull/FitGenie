# internal/

Código privado del proyecto.

## Estructura

```
internal/
├── api/handlers/    # HTTP handlers (User, Clothing, Outfit, Color)
├── config/         # Env vars
├── models/         # GORM models
├── repository/     # DB access (User, Clothing, Outfit repos)
└── services/       # Business logic (Color, Style, AI)
```

## Flujo

```
Request → Handler → Service → Repository → Database
```

## Convenciones

- Handlers usan dependency injection
- Repositories usan context.Context
- Services no conocen HTTP
- Models usan hooks GORM para UUIDs
