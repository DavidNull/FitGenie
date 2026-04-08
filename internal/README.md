# internal/

Código privado de FitGenie - lógica de negocio interna no exportable.

## Estructura

```
internal/
├── api/
│   └── handlers/          # Handlers HTTP (MVC Controller)
├── config/                # Configuración por variables de entorno
├── models/                # Modelos de datos (GORM structs)
├── repository/            # Capa de acceso a datos (Repository Pattern)
└── services/              # Lógica de negocio (Color, Style, AI)
```

## Diagrama de flujo de datos

```
Request → Handler → Service → Repository → Database
              ↓
         Response
```

## Descripción de paquetes

### `api/handlers/`
Handlers HTTP que procesan requests REST:
- **user_handler.go**: CRUD usuarios, perfiles de estilo/color, favoritos
- **clothing_handler.go**: Gestión de prendas de ropa
- **outfit_handler.go**: Creación de outfits y recomendaciones IA
- **color_handler.go**: Endpoints de teoría del color

### `config/`
Configuración centralizada:
- Lectura de variables de entorno
- Validación de configuración
- Valores por defecto

### `models/`
Structs de GORM que definen el esquema de la base de datos:
- User, StyleProfile, ColorProfile
- ClothingItem, Outfit, OutfitRecommendation
- FavoriteOutfit

### `repository/`
Abstracción de persistencia (Repository Pattern):
- **user_repository.go**: Operaciones CRUD de usuarios
- **clothing_repository.go**: Operaciones de prendas
- **outfit_repository.go**: Operaciones de outfits

Beneficios:
- Tests unitarios con mocks fáciles
- Cambio de DB sin modificar handlers
- Código desacoplado

### `services/`
Lógica de negocio compleja:
- **color_theory.go**: Análisis de colorimetría, armonías, estaciones
- **style_service.go**: Análisis de estilo personal, categorías de moda
- **ai_service.go**: Algoritmos de recomendación de outfits

## Convenciones

1. **Handlers** reciben dependencias por constructor (DI)
2. **Repositories** usan context.Context para timeouts/cancelación
3. **Services** no conocen HTTP, operan solo con structs de models
4. **Models** definen hooks GORM (BeforeCreate para UUIDs)

## Ejemplo de flujo

```go
// 1. Handler recibe request
func (h *UserHandler) CreateUser(c *gin.Context) {
    // 2. Parsea input
    var user models.User
    c.ShouldBindJSON(&user)
    
    // 3. Llama a repository
    err := h.repo.Create(c.Request.Context(), &user)
    
    // 4. Responde
    c.JSON(201, user)
}
```
