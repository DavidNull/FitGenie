# FitGenie

Aplicación para gestionar tu armario y recibir recomendaciones de outfits con IA.

![FitGenie Banner](./docs/img/Banner.png)

## Descripción

FitGenie te permite:
- **Gestionar tu armario**: Añadir, editar y eliminar prendas de ropa
- **Ver detalles**: Color, estilo, temporada de cada prenda
- **Recibir recomendaciones**: La IA sugiere outfits basados en ocasión y temporada
- **Guardar favoritos**: Almacenar outfits que te gusten

## Arquitectura de Desarrollo Local

```
┌─────────────┐     WiFi/Red     ┌─────────────┐
│ Flutter App │ ◄══════════════► │ Backend Go  │
│  (Móvil)    │                  │ + Postgres  │
└─────────────┘                  │ + S3 Local  │
                                 └─────────────┘
                                       PC (WSL)
```

**Nota importante**: Esta configuración es para **desarrollo local**. Todo corre en tu PC y la app móvil se conecta vía IP.

## Requisitos

- Docker y Docker Compose
- Flutter SDK
- Dispositivo móvil o emulador en la misma red WiFi

## Inicio Rápido

### 1. Iniciar Backend

```bash
cd /home/david/FitGenie
make docker-run
```

Verifica que todo esté corriendo:
```bash
docker ps
```

Deberías ver:
- `fitgenie-api` (puerto 8080)
- `fitgenie-postgres` (puerto 5432)
- `fitgenie-localstack` (puerto 4566)

### 2. Configurar IP del Backend

Obtén tu IP de WSL:
```bash
ip route | grep default | awk '{print $3}'
```

Edita `mobile/lib/services/api_service.dart`:
```dart
static String apiHost = '172.21.48.1';  // Tu IP de WSL
```

**Para diferentes entornos:**
- **Emulador Android**: `10.0.2.2`
- **iOS Simulator**: `localhost`
- **Dispositivo físico**: IP de tu PC en la red WiFi
- **WSL/Linux**: IP que te da el comando anterior

### 3. Ejecutar Flutter

```bash
cd /home/david/FitGenie/mobile
flutter pub get
flutter run
```

## Uso de la App

### Primera vez

1. Ve a **Galería** (estará vacía)
2. Toca **"Usar imágenes de ejemplo"**
3. Se importarán 5 prendas de ejemplo con datos para recomendaciones
4. Ve a **Recomendaciones** para generar outfits con IA

### Añadir prendas propias

1. Ve a **Cámara** (icono + en el menú inferior)
2. Selecciona **Cámara** o **Galería** del dispositivo
3. Elige la categoría (Parte de arriba/Parte de abajo/Calzado)
4. La imagen se sube automáticamente y aparece en tu armario

### Generar Outfits

1. Ve a **Recomendaciones**
2. Selecciona **Ocasión** (Casual, Formal, Trabajo...)
3. Selecciona **Temporada** (Verano, Invierno...)
4. Toca **"Generar Outfits"**
5. La IA sugerirá combinaciones basadas en tu armario
6. Guarda los que te gusten en favoritos

## Stack Tecnológico

### Backend
- **Go 1.23** + Gin framework
- **PostgreSQL** para datos
- **LocalStack S3** para imágenes (modo desarrollo)
- **Docker Compose** para orquestación

### Frontend
- **Flutter 3.x** con Material 3
- **Provider** para gestión de estado
- **Image Picker** para acceso a cámara/galería
- **HTTP** para comunicación con API

## Mejoras Futuras (Firebase)

Para pasar de desarrollo local a producción, se recomienda:

### Firebase Integration
```
┌─────────────┐              ┌─────────────┐
│ Flutter App │ ◄──────────► │   Firebase  │
│             │              │  - Auth     │
│             │              │  - Firestore│
│             │              │  - Storage  │
└─────────────┘              └─────────────┘
```

**Ventajas:**
- **Auth**: Login con Google/Apple/email
- **Firestore**: Base de datos en la nube (sin servidor propio)
- **Storage**: Almacenamiento de imágenes (sin S3)
- **Hosting**: App web si se necesita
- **Functions**: Backend serverless si se necesita lógica extra

**Implementación:**
1. Crear proyecto en Firebase Console
2. Añadir `firebase_core`, `firebase_auth`, `cloud_firestore`, `firebase_storage`
3. Migrar modelos de PostgreSQL a Firestore
4. Reemplazar S3 por Firebase Storage
5. Deploy backend a Cloud Run o usar Firebase Functions

## Estructura del Proyecto

```
FitGenie/
├── cmd/api/              # Entry point backend
├── internal/
│   ├── api/              # HTTP handlers y rutas
│   ├── models/           # Modelos de datos
│   ├── repository/       # Acceso a BD
│   └── services/         # Lógica de negocio
├── pkg/
│   ├── logger/           # Logging
│   └── storage/          # Cliente S3
├── mobile/
│   ├── lib/
│   │   ├── screens/      # Pantallas Flutter
│   │   ├── providers/    # Estado con Provider
│   │   ├── services/     # API client
│   │   └── models/       # Modelos Dart
│   └── assets/           # Imágenes y recursos
├── docker-compose.yml    # Configuración Docker
└── README.md
```

## Comandos Útiles

```bash
# Backend
make docker-run           # Iniciar todo
make docker-stop          # Detener todo
docker logs fitgenie-api  # Ver logs del API

# Frontend
cd mobile
flutter pub get           # Instalar dependencias
flutter run               # Ejecutar app
flutter build apk         # Build Android
flutter build ios         # Build iOS

# BD
make migrate-up           # Aplicar migraciones
make migrate-down         # Revertir migraciones
```

## Solución de Problemas

### "Connection refused" o timeout
- Verifica que la IP en `api_service.dart` sea correcta
- Comprueba que los contenedores estén corriendo: `docker ps`
- Reinicia: `docker-compose restart api`

### Error al subir imágenes
- Verifica que LocalStack esté corriendo: `docker ps | grep localstack`
- Reinicia todos los servicios: `make docker-stop && make docker-run`

### La app no carga prendas
- Comprueba conexión: `curl http://TU_IP:8080/api/v1/users/me -H "X-Device-ID: test"`
- Verifica logs del backend: `docker logs fitgenie-api`

## Licencia

MIT License - Ver LICENSE para detalles.
