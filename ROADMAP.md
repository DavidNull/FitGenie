# Roadmap FitGenie

Lista detallada de features pendientes y planificación de desarrollo.

## Leyenda

- 🔴 **Crítico**: Bloqueante para MVP
- 🟡 **Alto**: Importante para v1.0
- 🟢 **Medio**: Mejora experiencia
- ⚪ **Bajo**: Nice to have

---

## Backend API

### Autenticación & Seguridad 🔴

| Feature | Estado | Descripción | Estimación |
|---------|--------|-------------|------------|
| JWT Authentication | ❌ No iniciado | Login/registro con tokens JWT | 2 días |
| Password Hashing (bcrypt) | ❌ No iniciado | Almacenar contraseñas seguras | 0.5 días |
| Refresh Tokens | ❌ No iniciado | Rotación de tokens de sesión | 1 día |
| Auth Middleware | ❌ No iniciado | Proteger rutas privadas | 1 día |
| Rate Limiting | ❌ No iniciado | Prevenir abuso de API | 1 día |
| CORS Configuration | ❌ No iniciado | Configurar orígenes permitidos | 0.5 días |

**Total estimado**: 6 días

### Gestión de Imágenes 🟡

| Feature | Estado | Descripción | Estimación |
|---------|--------|-------------|------------|
| Upload Endpoint | ❌ No iniciado | POST /api/v1/clothing/:id/upload | 1 día |
| Image Validation | ❌ No iniciado | Validar tipo, tamaño, dimensiones | 0.5 días |
| S3 Integration | 🟡 Parcial | Cliente S3 existe, falta integrar en handlers | 1 día |
| Image Resizing | ❌ No iniciado | Generar thumbnails (Go/ImageMagick) | 1 día |
| Color Extraction | ❌ No iniciado | Extraer colores dominantes de imagen | 2 días |
| WebP Conversion | ❌ No iniciado | Optimizar formato de imagen | 1 día |

**Total estimado**: 6.5 días

### Búsqueda Vectorial (pgvector) 🟡

| Feature | Estado | Descripción | Estimación |
|---------|--------|-------------|------------|
| Image Embeddings | ❌ No iniciado | Generar vectores de imágenes (CLIP/ResNet) | 3 días |
| Vector Storage | 🟡 Parcial | pgvector configurado, falta almacenar embeddings | 0.5 días |
| Similarity Search | ❌ No iniciado | Buscar prendas similares por imagen | 1 día |
| Visual Recommendations | ❌ No iniciado | "Parecido a esto que te gusta" | 2 días |

**Total estimado**: 6.5 días

### Features de Negocio 🟢

| Feature | Estado | Descripción | Estimación |
|---------|--------|-------------|------------|
| Weather Integration | ❌ No iniciado | Conectar API meteorológica para recomendaciones según clima | 1 día |
| Calendar Events | ❌ No iniciado | Sincronizar con calendario para outfit según evento | 2 días |
| Outfit History | ❌ No iniciado | Registrar outfits usados con fecha | 1 día |
| Packing Lists | ❌ No iniciado | Generar lista de ropa para viajes | 2 días |
| Laundry Tracker | ❌ No iniciado | Seguimiento de prendas en lavandería | 1 día |
| Wishlist | ❌ No iniciado | Lista de deseos de prendas (links externos) | 1 día |

**Total estimado**: 8 días

### Performance & Cache ⚪

| Feature | Estado | Descripción | Estimación |
|---------|--------|-------------|------------|
| Redis Cache | ❌ No iniciado | Cache de recomendaciones y perfiles | 2 días |
| DB Connection Pool | 🟢 Done | Ya implementado en database.go | ✅ |
| Request Caching | ❌ No iniciado | Cache de respuestas HTTP (ETag) | 1 día |
| Background Jobs | ❌ No iniciado | Worker pool para tareas pesadas | 2 días |
| CDN Integration | ❌ No iniciado | CloudFront/CloudFlare para imágenes | 1 día |

**Total estimado**: 6 días

---

## Frontend Móvil 🔴

**FitGenie actualmente es SOLO backend. No existe aplicación móvil.**

### Estructura del Proyecto Móvil

Recomendación: Crear carpeta `mobile/` en la raíz:

```
mobile/
├── flutter/              # Opción A: Flutter (recomendado)
│   ├── android/
│   ├── ios/
│   ├── lib/
│   └── pubspec.yaml
│
└── react-native/         # Opción B: React Native
    ├── android/
    ├── ios/
    ├── src/
    └── package.json
```

### Features App Móvil (Flutter)

| Feature | Prioridad | Descripción | Estimación |
|---------|-----------|-------------|------------|
| **Setup proyecto** | 🔴 | Inicializar Flutter con navegación | 1 día |
| **Onboarding** | 🔴 | Tutorial inicial, permisos cámara | 1 día |
| **Auth UI** | 🔴 | Login/Registro screens | 2 días |
| **Cámara** | 🔴 | Fotografiar prenda + preview | 2 días |
| **Galería** | 🔴 | Seleccionar foto existente | 1 día |
| **Color Picker** | 🔴 | Seleccionar colores de prenda manual | 2 días |
| **Catálogo prendas** | 🔴 | Grid/listado con filtros | 2 días |
| **Detalle prenda** | 🔴 | Ver/editar/eliminar prenda | 2 días |
| **Perfil estilo** | 🟡 | Wizard de preguntas de estilo | 3 días |
| **Perfil color** | 🟡 | Análisis de colorimetría personal | 2 días |
| **Recomendaciones** | 🔴 | Vista de outfits sugeridos | 2 días |
| **Favoritos** | 🟡 | Guardar outfits favoritos | 1 día |
| **Calendario** | 🟢 | Planificar outfits por fecha | 3 días |
| **Notificaciones** | ⚪ | Push diario de recomendaciones | 2 días |
| **Offline mode** | ⚪ | Cache local para usar sin red | 3 días |
| **Widget home** | ⚪ | Widget iOS/Android rápido | 2 días |

**Total Flutter estimado**: 30-35 días (1 desarrollador senior)

### Stack Tecnológico Recomendado (Flutter)

```yaml
Framework: Flutter 3.x
Dart: 3.x

State Management: Riverpod / BLoC
Navigation: GoRouter
HTTP Client: Dio
Local DB: Hive / SQLite
Image Cache: CachedNetworkImage
Camera: camera plugin
Permissions: permission_handler
Push: firebase_messaging
Auth: flutter_secure_storage
```

### Backend Requirements para Móvil

Para soportar la app móvil, el backend necesita:

1. **Endpoints faltantes**:
   ```
   POST /api/v1/auth/login
   POST /api/v1/auth/register
   POST /api/v1/auth/refresh
   POST /api/v1/clothing/:id/upload
   GET  /api/v1/users/me (perfil actual)
   ```

2. **Mobile-optimized responses**:
   - Paginación eficiente (cursor-based)
   - Imágenes en múltiples resoluciones (thumbnail, medium, full)
   - Caché headers apropiados

---

## Infraestructura & DevOps

### Cloud Deployment 🟡

| Feature | Estado | Descripción | Estimación |
|---------|--------|-------------|------------|
| Terraform AWS | ❌ No iniciado | IaC para RDS, EKS, S3 | 3 días |
| Kubernetes Manifests | ❌ No iniciado | Deployments, Services, Ingress | 2 días |
| Helm Charts | ❌ No iniciado | Templating de K8s | 2 días |
| GitOps (ArgoCD) | ❌ No iniciado | Despliegue automático | 2 días |
| Cert Manager | ❌ No iniciado | TLS automático con Let's Encrypt | 1 día |
| External DNS | ❌ No iniciado | Gestión DNS automática | 0.5 días |

**Total estimado**: 10.5 días

### Observabilidad Avanzada 🟢

| Feature | Estado | Descripción | Estimación |
|---------|--------|-------------|------------|
| Grafana Dashboards | ❌ No iniciado | Dashboards de métricas | 1 día |
| Alert Manager | ❌ No iniciado | Alertas Slack/PagerDuty | 1 día |
| Distributed Tracing (Jaeger) | ❌ No iniciado | Trace de requests entre servicios | 2 días |
| Log Aggregation (Loki) | ❌ No iniciado | Centralización de logs | 1 día |
| APM (NewRelic/Datadog) | ❌ No iniciado | Performance monitoring | 2 días |

**Total estimado**: 7 días

### Seguridad Infraestructura 🟡

| Feature | Estado | Descripción | Estimación |
|---------|--------|-------------|------------|
| WAF (AWS/CloudFlare) | ❌ No iniciado | Protección contra ataques web | 1 día |
| Secrets Management | ❌ No iniciado | AWS Secrets Manager / Vault | 1 día |
| Network Policies | ❌ No iniciado | Isolamiento de servicios K8s | 1 día |
| Security Scanning | ❌ No iniciado | Trivy para imágenes Docker | 0.5 días |
| DDoS Protection | ❌ No iniciado | CloudFlare/AWS Shield | 0.5 días |

**Total estimado**: 4 días

---

## Integraciones Externas ⚪

### E-commerce

| Feature | Prioridad | Descripción | Estimación |
|---------|-----------|-------------|------------|
| Shopify Integration | ⚪ | Importar catálogo de tienda | 3 días |
| Amazon Affiliate | ⚪ | Links de compra recomendados | 2 días |
| ASOS API | ⚪ | Búsqueda de prendas similares | 3 días |

### Redes Sociales

| Feature | Prioridad | Descripción | Estimación |
|---------|-----------|-------------|------------|
| Instagram Share | ⚪ | Compartir outfit en stories | 2 días |
| Pinterest Save | ⚪ | Guardar outfit en Pinterest | 2 días |

---

## Milestones Propuestos

### MVP (Mínimo Producto Viable)
**Objetivo**: App funcional con usuario de prueba

**Fecha estimada**: 6-8 semanas (1 dev backend + 1 dev mobile)

**Features incluidos**:
- ✅ Backend API (actual)
- 🔴 JWT Auth
- 🔴 Upload imágenes
- 🔴 App móvil básica (Flutter)
- 🔴 CRUD prendas desde móvil
- 🔴 Recomendaciones simples

### v1.0 (Producto Completo)
**Objetivo**: Producto listo para launch público

**Fecha estimada**: 12-16 semanas (+ 2 devs, + QA)

**Features adicionales**:
- 🟡 Perfiles color/estilo
- 🟡 Favoritos
- 🟡 Search vectorial
- 🟡 Deploy producción AWS
- 🟡 Tests E2E automatizados

### v2.0 (Scale)
**Objetivo**: Escalar a miles de usuarios

**Fecha estimada**: 24+ semanas

**Features adicionales**:
- ⚪ Cache Redis
- ⚪ Background jobs
- ⚪ ML avanzado
- ⚪ Marketplace integración

---

## Recursos Necesarios

### Equipo Ideal

| Rol | Cantidad | Tiempo | Tarea Principal |
|-----|----------|--------|-----------------|
| Backend Senior (Go) | 1 | 100% | API, ML, Infra |
| Mobile Senior (Flutter) | 1 | 100% | App iOS/Android |
| DevOps | 0.5 | 50% | K8s, CI/CD, AWS |
| Diseñador UX/UI | 0.5 | 50% | Mockups, App design |
| QA Engineer | 0.5 | 50% | Tests, Automation |
| Product Manager | 0.25 | 25% | Roadmap, Priorización |

### Infraestructura AWS (estimado mensual)

| Servicio | Costo Estimado |
|----------|---------------|
| EKS (3 nodos t3.medium) | $150/mes |
| RDS PostgreSQL (db.t3.micro) | $15/mes |
| S3 (almacenamiento imágenes) | $10-50/mes |
| CloudFront (CDN) | $20/mes |
| ALB (Load Balancer) | $20/mes |
| **Total** | **~$200-250/mes** |

---

## Notas de Implementación

### Próximos Pasos Inmediatos

1. **Esta semana**:
   - Implementar JWT Auth (backend)
   - Crear endpoint upload imágenes
   - Setup proyecto Flutter

2. **Próximas 2 semanas**:
   - Screens de login/onboarding Flutter
   - Integración cámara/galería
   - Upload de primera prenda end-to-end

3. **Mes 1**:
   - CRUD completo prendas
   - Vista recomendaciones básica
   - Perfiles estilo/color

## Contacto

Para dudas sobre roadmap: [davidnull@example.com]

---

**Última actualización**: Abril 2026
**Versión documento**: 1.0
