# FitGenie
![Logo](./docs/img/Banner.png)

App para recomendar outfits de ropa con IA. Sube tus prendas, recibe combinaciones inteligentes basadas en color, estilo y ocasión.

## 🏗️ Arquitectura

```
┌─────────────────┐     WiFi/Red local     ┌─────────────────┐
│  Flutter App    │  ═══════════════════►  │  Backend Go     │
│  (Móvil/PC)     │                        │  + PostgreSQL   │
│                 │  ◄═════════════════════│  + S3 (Local)   │
└─────────────────┘     HTTP API :8080     └─────────────────┘
```

**Flujo:**
- **Backend** corre en tu PC (Docker/WSL) con toda la lógica y la BD
- **Flutter App** corre en móvil/emulador y se conecta vía IP al backend
- **Ventaja:** Puedes desarrollar y probar sin deployar nada a la nube

## 📱 Características

### Funcionalidades Implementadas ✅

| Feature | Descripción |
|---------|-------------|
| **Galería de prendas** | Ver todas las prendas con imágenes 3:4 |
| **Detalle de prenda** | Al tocar, ver info completa (color, estilo, ocasión) |
| **Añadir prendas** | Dos modos: cámara/galería (móvil) o assets locales (PC) |
| **Eliminar prendas** | Botón rojo con confirmación en cada prenda |
| **AI Recomendaciones** | Outfits generados según ocasión y temporada |
| **Filtros** | Filtrar por ocasión (casual, formal, etc.) y temporada |
| **Guardar outfits** | Añadir outfits recomendados a favoritos |
| **Navegación** | Bottom nav persistente en todas las pantallas |

### Modos de Imagen

**Para Móvil/Emulador (con S3):**
1. Toca icono cámara o galería
2. Selecciona imagen del dispositivo
3. Se sube automáticamente a S3
4. Se guarda URL en BD

**Para Desarrollo Local (sin S3):**
1. Coloca imágenes en `mobile/assets/clothing/`
2. En pantalla "Cámara", sección "modo desarrollo"
3. Toca una imagen local
4. Se guarda path local (`assets/clothing/xxx.jpg`)
5. No requiere subir nada

**Archivos locales de ejemplo:**
- `c1.jpg` - Camiseta Azul
- `c2.jpg` - Camisa Blanca  
- `p1.jpg` - Pantalón Negro
- `p2.jpg` - Zapatillas

## 🛠️ Stack Técnico

### Backend
- **Go 1.23** + Gin framework
- **PostgreSQL** + pgvector para vectores de color
- **LocalStack S3** para almacenamiento de imágenes (dev)
- **Docker + Compose** para orquestación

### Frontend (Flutter)
- **Flutter 3.x** con Material 3
- **Provider** para estado
- **Image Picker** para cámara/galería
- **HTTP** + **http_parser** para API calls

## 🚀 Cómo usar

### 1. Iniciar Backend (PC)

```bash
# Docker (recomendado - todo automático)
make docker-run

# Ver servicios corriendo
docker ps
# - API en :8080
# - PostgreSQL en :5432
# - S3 (LocalStack) en :4566
```

### 2. Configurar IP del Backend

Edita `mobile/lib/services/api_service.dart`:

```dart
// Para emulador Android:
static String apiHost = '10.0.2.2';  // localhost del host

// Para iOS Simulator:
static String apiHost = 'localhost';

// Para dispositivo físico:
static String apiHost = '192.168.1.xxx';  // IP de tu PC

// Para Linux (desktop):
static String apiHost = '172.21.56.127';  // IP WSL
```

### 3. Ejecutar Flutter App

```bash
cd mobile

# Emulador Android
flutter run -d emulator

# iOS Simulator
flutter run -d ios

# Linux (desktop)
flutter run -d linux

# Tu dispositivo físico
flutter run -d <device-id>
```

### 4. Verificar Conexión

La app llama automáticamente a `GET /api/v1/users/me` al iniciar:
- Si hay error de conexión → muestra error
- Si funciona → crea usuario con `X-Device-ID` y carga prendas

## 📁 Estructura del Proyecto

```
FitGenie/
├── cmd/server/              # Entrypoint Go
├── internal/
│   ├── api/handlers/        # HTTP handlers (clothing, outfits, upload)
│   ├── models/              # Modelos GORM
│   ├── repository/          # Acceso a datos PostgreSQL
│   └── services/            # Lógica de negocio + IA
├── pkg/                     # Librerías reutilizables
├── mobile/
│   ├── lib/
│   │   ├── screens/         # UI (Home, Gallery, Camera, Recommendations, Detail)
│   │   ├── providers/       # AppProvider (estado global)
│   │   ├── services/        # ApiService (HTTP client)
│   │   └── models/          # ClothingItem, Outfit, etc.
│   └── assets/
│       ├── clothing/        # Imágenes locales para testing
│       └── *.png            # Icons de navegación
├── docker-compose.yml       # PostgreSQL + API + S3
└── Makefile                 # Comandos útiles
```

## 🔌 Endpoints API

### Auth
Todas las peticiones requieren header `X-Device-ID`.

### Users
- `GET /api/v1/users/me` - Usuario actual (auto-crea si no existe)

### Clothing (Prendas)
- `GET /api/v1/clothing?user_id=xxx` - Listar prendas del usuario
- `POST /api/v1/clothing` - Crear prenda
- `DELETE /api/v1/clothing/:id` - Eliminar prenda

### Outfits
- `POST /api/v1/outfits` - Crear outfit manualmente
- `POST /api/v1/users/:id/outfits/recommendations` - Generar con IA

### Upload
- `POST /api/v1/upload` - Subir imagen (multipart/form-data, campo `image`)
- Soporta: JPG, JPEG, PNG, WEBP (max 5MB)

## 📝 Notas de Desarrollo

### Solución de Problemas

**Error "unable to load asset":**
Las imágenes en `assets/clothing/` son solo para desarrollo local. Si la BD tiene referencias a assets que no existen, la app muestra placeholder en lugar de crashear.

**Error de conexión "Connection refused":**
Verifica que:
1. Docker está corriendo (`docker ps`)
2. La IP en `api_service.dart` es correcta para tu plataforma
3. El firewall no bloquea el puerto 8080

**Hot Reload vs Hot Restart:**
- Hot Reload (`r`) - Actualiza UI pero no estado
- Hot Restart (`R`) - Reinicia app, recarga assets y código

### Archivos Clave

| Archivo | Descripción |
|---------|-------------|
| `mobile/lib/services/api_service.dart` | Configurar IP del backend |
| `mobile/lib/providers/app_provider.dart` | Estado global de la app |
| `mobile/lib/screens/camera_screen.dart` | Lógica de añadir prendas (dual: upload/local) |
| `mobile/lib/screens/gallery_screen.dart` | Grid de prendas con navegación a detalle |
| `mobile/lib/screens/recommendations_screen.dart` | Outfits generados por IA |

---

**FitGenie** — Por DavidNull
