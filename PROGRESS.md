# FitGenie - Progreso y Pendientes

**Fecha:** 2026-04-14  
**Sprint:** Conexión App ↔ IA Backend

---

## ✅ Estado Actual (2026-04-14 10:30) - BUGS FIXEADOS

### Bugs Resueltos
| Bug | Estado | Fix |
|-----|--------|-----|
| `ListClothing` vacío | ✅ **FIXED** | Usar `pq.StringArray` en modelos para PostgreSQL arrays |
| `CreateClothing` error 500 | ✅ **FIXED** | Mismo fix - arrays PostgreSQL |
| Rutas Gin orden | ✅ **FIXED** | Mover `GET ""` antes que `GET "/:id"` en `routes.go` |
| Device auth duplicate key | ✅ **FIXED** | Reintentar `GetByID` si `Create` falla por duplicate key |

### Root Cause
El problema principal era que PostgreSQL arrays (`text[]`) no se escaneaban correctamente en Go `[]string`. El driver estándar devolvía strings como `{"summer","spring"}` pero GORM no los convertía automáticamente.

**Solución:** Usar `pq.StringArray` de `github.com/lib/pq` en los modelos:
```go
Season   pq.StringArray `json:"season" gorm:"type:text[]"`
Occasion pq.StringArray `json:"occasion" gorm:"type:text[]"`
```

### Estado Endpoints
| Endpoint | Estado |
|----------|--------|
| `GET /api/v1/clothing` | ✅ **Funcionando** |
| `POST /api/v1/clothing` | ✅ **Funcionando** |
| `POST /api/v1/users/:id/outfits/recommendations` | ✅ **Funcionando** |

### Backend 100% Funcional 🎉
Todos los bugs del backend Go han sido resueltos. El sistema ahora puede:
- Crear prendas de ropa
- Listar prendas por usuario
- Generar recomendaciones de outfits con IA

---

## ✅ Completado Hoy

### UI/UX Flutter
- [x] Logo más grande en Home y Camera
- [x] Bottom navigation con iconos circulares y sombra
- [x] Gallery con gradiente animado (blanco → azul al hacer scroll)
- [x] Cámara con fondo claro (#F5F8FA)

### Arquitectura App Móvil
- [x] Modelos creados:
  - `ClothingItem` - Representa prendas del armario
  - `Outfit` - Conjuntos de ropa
  - `OutfitRecommendation` - Recomendaciones de la IA
- [x] `ApiService` - Cliente HTTP para backend
- [x] `AppProvider` - Gestión de estado con Provider
- [x] Integración Provider en `main.dart`

### Backend (Go)
- [x] `clothing_handler.go` - Añadido `userID` desde contexto
- [x] `outfit_handler.go` - Añadido logging de debug
- [x] Fix tabla `clothing_items` creada correctamente
- [x] Usuario de prueba creado en DB con 4 prendas

---

## 🔴 Bugs Críticos (Opción A - Fixear)

### Bug 1: ListClothing devuelve vacío
**Síntoma:** `GET /api/v1/clothing?user_id=XXX` retorna `{"items":[],"total":0}` aunque la DB tiene 4 prendas.

**Investigación:**
```sql
SELECT COUNT(*) FROM clothing_items WHERE user_id = '172f2ee4-ddea-4351-8b03-fa05fd28d05d';
-- Result: 4 (correcto)
```

**Causa probable:** 
- El `user_id` en la query GORM no está haciendo match con el formato UUID
- Posible problema de conversión entre `uuid.UUID` y string en PostgreSQL

**Fix propuesto:**
```go
// En clothing_repository.go - método ListByUser
// Cambiar: Where("user_id = ?", userID)
// A:      Where("user_id::text = ?", userID.String())
```

**Archivos:**
- `internal/repository/clothing_repository.go`

---

### Bug 2: CreateClothing error 500
**Síntoma:** `POST /api/v1/clothing` retorna `{"error":"Failed to create clothing item"}`

**Investigación:**
- El handler ya obtiene `userID` del contexto correctamente
- El error ocurre en `repo.Create()`
- Posible constraint violation o tipo de dato

**Causa probable:**
- Las columnas `season` y `occasion` son `text[]` en PostgreSQL
- GORM puede no estar manejando correctamente los arrays de strings
- O falta el `user_id` en el INSERT

**Fix propuesto:**
1. Verificar que `item.UserID` no es nil antes de crear
2. Loggear el error exacto de PostgreSQL
3. Si es problema de arrays, usar `pq.StringArray` o similar

**Archivos:**
- `internal/repository/clothing_repository.go`
- `internal/api/handlers/clothing_handler.go`

---

## 📋 Pendientes Mañana

### Prioridad 1: Fix Backend (Opción A)
- [ ] Debug y fix `ListByUser` query en `clothing_repository.go`
- [ ] Debug y fix `Create` en `clothing_repository.go`
- [ ] Verificar que `user_id` se pasa correctamente a todas las queries
- [ ] Test endpoint: `GET /api/v1/clothing`
- [ ] Test endpoint: `POST /api/v1/clothing`
- [ ] Test endpoint: `POST /api/v1/users/:userId/outfits/recommendations`

### Prioridad 2: Integrar Imágenes
**Nuevos assets recibidos:**
- `c*.jpg` - Prendas parte superior (camisetas, chaquetas...)
- `p*.jpg` - Prendas parte inferior (pantalones, shorts...)

**Tareas:**
- [ ] Copiar imágenes a `mobile/assets/clothing/`
- [ ] Añadir a `pubspec.yaml`
- [ ] Actualizar `GalleryScreen` para mostrar imágenes reales desde DB
- [ ] Actualizar `CameraScreen` para capturar y subir fotos

### Prioridad 3: UI Recomendaciones
- [ ] Crear pantalla `RecommendationScreen`
- [ ] Mostrar outfits recomendados por IA
- [ ] Botón "Aceptar outfit" (guarda en favoritos)
- [ ] Botón "Pedir otra recomendación"

### Prioridad 4: Testing End-to-End
- [ ] Flujo: Añadir prenda → Ver en galería → Pedir recomendación → Ver outfit IA

---

## 🔧 Comandos Útiles

```bash
# Ver logs del API
docker logs fitgenie-api --tail 20

# Ver prendas en DB
docker exec fitgenie-postgres psql -U fitgenie -d fitgenie -c "SELECT * FROM clothing_items;"

# Test API manual
curl -s http://localhost:8080/api/v1/users/me -H "X-Device-ID: flutter-test-device"

# Rebuild y reiniciar API
cd /home/david/FitGenie && docker compose build api && docker compose up -d api

# Correr Flutter
cd /home/david/FitGenie/mobile && flutter run -d linux
```

---

## 🔗 URLs Importantes

- API Health: http://localhost:8080/health
- API Docs: http://localhost:8080/api/v1 (endpoints)
- Flutter App: `flutter run -d linux`

---

## 🎯 User ID de Prueba

```
Device ID: flutter-test-device
User ID: 172f2ee4-ddea-4351-8b03-fa05fd28d05d
Prendas en DB: 4 (Camiseta Azul, Pantalón Negro, Zapatillas Blancas, Chaqueta Vaquera)
```

---

## 📝 Notas

- El middleware de auth (`DeviceAuthMiddleware`) funciona correctamente
- La DB tiene las tablas creadas por GORM
- El problema parece ser específicamente en las queries de `clothing_repository.go`
- Considerar añadir `log.Printf` en cada query para debug
