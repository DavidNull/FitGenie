package middleware

import (
	"crypto/md5"
	"net/http"
	"strings"

	"fitgenie/internal/models"
	"fitgenie/internal/repository"
	"fitgenie/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DeviceAuthMiddleware crea o recupera usuario basado en Device ID
// Acepta cualquier string como device ID (no requiere UUID format)
func DeviceAuthMiddleware(userRepo repository.UserRepository, log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.GetHeader("X-Device-ID")

		// Si no hay device ID, generar uno nuevo (UUID)
		if deviceID == "" {
			deviceID = uuid.New().String()
			c.Header("X-Device-ID", deviceID)
		}

		// Generar UUID determinístico a partir del deviceID string
		// Esto permite cualquier formato de device ID
		deviceUUID := generateUUIDFromString(deviceID)

		// Buscar usuario por device ID
		ctx := c.Request.Context()
		user, err := userRepo.GetByID(ctx, deviceUUID)

		if err != nil {
			// Usuario no existe, crear uno nuevo
			user = &models.User{
				ID:    deviceUUID,
				Email: sanitizeDeviceID(deviceID) + "@device.fitgenie",
				Name:  "User " + deviceID[:min(8, len(deviceID))],
			}
			if err := userRepo.Create(ctx, user); err != nil {
				// Si el error es por duplicate key, el usuario ya existe (race condition)
				// En ese caso, intentamos obtener el usuario de nuevo
				user, err = userRepo.GetByID(ctx, deviceUUID)
				if err != nil {
					log.Error("failed to create or get device user", "error", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
					c.Abort()
					return
				}
			}
		}

		// Guardar user ID en contexto
		c.Set("userID", user.ID)
		c.Set("deviceID", deviceUUID)
		c.Next()
	}
}

// generateUUIDFromString genera un UUID v4 determinístico a partir de un string
func generateUUIDFromString(s string) uuid.UUID {
	hash := md5.Sum([]byte(s))
	// Set version (4) and variant bits for valid UUID v4
	hash[6] = (hash[6] & 0x0f) | 0x40
	hash[8] = (hash[8] & 0x3f) | 0x80
	return uuid.UUID(hash[:])
}

// sanitizeDeviceID limpia el device ID para usarlo en email
func sanitizeDeviceID(s string) string {
	// Reemplazar caracteres no válidos para email
	s = strings.ReplaceAll(s, "@", "_")
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, ":", "_")
	s = strings.ReplaceAll(s, "\\", "_")
	return s
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
