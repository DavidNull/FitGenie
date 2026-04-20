# FitGenie Mobile

App Flutter para FitGenie. Esta rama contiene solo el código móvil.

## Requisitos

- Flutter SDK 3.x
- Conexión al backend (ver instrucciones abajo)

## Configuración

### 1. Configurar IP del Backend

Edita `lib/services/api_service.dart`:

```dart
// Para emulador Android:
static String apiHost = '10.0.2.2';

// Para iOS Simulator:
static String apiHost = 'localhost';

// Para dispositivo físico:
static String apiHost = '192.168.1.xxx';  // IP de tu PC

// Para Linux (desktop):
static String apiHost = '172.21.48.1';  // IP de WSL
```

**Para WSL (Windows):**
```bash
ip route | grep default | awk '{print $3}'
# Usa esa IP en apiHost
```

### 2. Iniciar Backend

```bash
# Usando Docker Hub
docker pull davidnull/fitgenie:1.0
docker run -p 8080:8080 -e S3_PUBLIC_ENDPOINT=http://TU_IP:4566 davidnull/fitgenie:1.0
```

O usa docker-compose del repositorio principal.

### 3. Ejecutar App

```bash
cd mobile
flutter pub get
flutter run
```

## Arquitectura

```
┌─────────────┐      HTTP       ┌─────────────┐
│  Flutter    │  ═══════════════► │  Backend    │
│   (tú)      │  ◄═══════════════│  (Docker)   │
└─────────────┘                 └─────────────┘
```

## Estructura

- `lib/screens/` - Pantallas de la app
- `lib/providers/` - Estado con Provider
- `lib/services/` - Cliente API
- `lib/models/` - Modelos de datos
- `assets/` - Imágenes locales

## Licencia

MIT License
