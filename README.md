# FitGenie Mobile

App Flutter para gestionar tu armario y recibir sugerencias de outfits.

## Requisitos

- Flutter SDK 3.x
- Backend FitGenie corriendo (ver abajo)

## Probarlo

### 1. Levantar el backend

```bash
# Opción A: Docker Compose completo (ver repo principal)
docker compose up -d

# Opción B: Solo imagen Docker
docker pull davidnull/fitgenie:latest
docker run -p 8080:8080 davidnull/fitgenie:latest
```

El backend debe estar en `http://localhost:8080`.

### 2. Ejecutar la app

```bash
flutter pub get
flutter run -d linux      # Linux
flutter run -d chrome     # Web
flutter run               # Android
```

**La app detecta automáticamente la IP del backend.** No necesitas configurar nada.

### 3. Usar

1. "Usar imágenes de ejemplo" → carga datos de prueba
2. "Armario" → ver tu ropa
3. "Recomendaciones" → pedir sugerencias

## Conexión backend

La app prueba automáticamente estas IPs al arrancar:
- `localhost`
- `10.0.2.2` (Android emulator)
- `172.17.0.1`, `172.21.48.1` (WSL/Docker)

Si falla, se puede forzar en `lib/services/api_service.dart`:

```dart
static String apiHost = 'TU_IP';  // manual override
```

## Estructura

```
lib/
├── main.dart              # Entrypoint
├── models/               # ClothingItem, Outfit...
├── providers/            # AppProvider (estado)
├── services/             # ApiService (HTTP client)
└── screens/              # UI: Home, Gallery, Camera...
```

## Conectar a tu backend

Por defecto la app busca el backend en IPs locales. Para producción:

```dart
// lib/services/api_service.dart
static String apiHost = 'api.tudominio.com';
```

## Build release

```bash
# Android APK
flutter build apk --release

# Web
flutter build web --release
```

## Repo principal

Código completo (backend + mobile): https://github.com/DavidNull/FitGenie
