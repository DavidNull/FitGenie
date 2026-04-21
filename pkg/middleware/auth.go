package middleware

import (
	"net/http"
	"strings"

	"fitgenie/pkg/auth"
	"fitgenie/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthMiddleware creates JWT authentication middleware
type AuthMiddleware struct {
	jwtService *auth.Service
	log        *logger.Logger
}

// NewAuthMiddleware creates a new auth middleware instance
func NewAuthMiddleware(jwtService *auth.Service, log *logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
		log:        log,
	}
}

// RequireAuth middleware validates JWT token from Authorization header
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		// Extract Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := m.jwtService.ValidateToken(tokenString)
		if err != nil {
			m.log.Warn("invalid token", "error", err, "ip", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// Check token type (must be access token)
		if claims.TokenType != "access" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token type"})
			c.Abort()
			return
		}

		// Store user info in context
		userID, _ := uuid.Parse(claims.UserID)
		c.Set("userID", userID)
		c.Set("userEmail", claims.Email)
		c.Set("deviceID", claims.DeviceID)

		c.Next()
	}
}

// OptionalAuth allows requests without auth but extracts token if present
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Next()
			return
		}

		claims, err := m.jwtService.ValidateToken(parts[1])
		if err == nil && claims.TokenType == "access" {
			userID, _ := uuid.Parse(claims.UserID)
			c.Set("userID", userID)
			c.Set("userEmail", claims.Email)
			c.Set("deviceID", claims.DeviceID)
			c.Set("authenticated", true)
		}

		c.Next()
	}
}

// GetUserID extracts user ID from context (must be used after RequireAuth)
func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, false
	}
	id, ok := userID.(uuid.UUID)
	return id, ok
}

// GetUserEmail extracts user email from context
func GetUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get("userEmail")
	if !exists {
		return "", false
	}
	str, ok := email.(string)
	return str, ok
}

// GetDeviceID extracts device ID from context
func GetDeviceID(c *gin.Context) (string, bool) {
	deviceID, exists := c.Get("deviceID")
	if !exists {
		return "", false
	}
	str, ok := deviceID.(string)
	return str, ok
}
