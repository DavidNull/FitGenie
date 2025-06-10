package handlers

import (
	"net/http"
	"strconv"

	"fitgenie/internal/database"
	"fitgenie/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserHandler handles user-related API endpoints
type UserHandler struct{}

// NewUserHandler creates a new user handler
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate UUID for new user
	user.ID = uuid.New()

	// Create user in database
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUser retrieves a user by ID
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("userId")

	var user models.User
	if err := database.DB.Preload("StyleProfile").Preload("ColorProfile").
		Preload("ClothingItems").Preload("Outfits").
		First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser updates an existing user
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("userId")

	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var updateData models.User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update allowed fields
	user.Name = updateData.Name
	user.Email = updateData.Email

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("userId")

	if err := database.DB.Delete(&models.User{}, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// CreateStyleProfile creates or updates a user's style profile
func (h *UserHandler) CreateStyleProfile(c *gin.Context) {
	userID := c.Param("userId")

	// Verify user exists
	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var profile models.StyleProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile.ID = uuid.New()
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	profile.UserID = userUUID

	// Check if profile already exists
	var existingProfile models.StyleProfile
	if err := database.DB.First(&existingProfile, "user_id = ?", userID).Error; err == nil {
		// Update existing profile
		profile.ID = existingProfile.ID
		if err := database.DB.Save(&profile).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update style profile"})
			return
		}
	} else {
		// Create new profile
		if err := database.DB.Create(&profile).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create style profile"})
			return
		}
	}

	c.JSON(http.StatusCreated, profile)
}

// GetStyleProfile retrieves a user's style profile
func (h *UserHandler) GetStyleProfile(c *gin.Context) {
	userID := c.Param("userId")

	var profile models.StyleProfile
	if err := database.DB.First(&profile, "user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Style profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// CreateColorProfile creates or updates a user's color profile
func (h *UserHandler) CreateColorProfile(c *gin.Context) {
	userID := c.Param("userId")

	// Verify user exists
	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var profile models.ColorProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile.ID = uuid.New()
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	profile.UserID = userUUID

	// Check if profile already exists
	var existingProfile models.ColorProfile
	if err := database.DB.First(&existingProfile, "user_id = ?", userID).Error; err == nil {
		// Update existing profile
		profile.ID = existingProfile.ID
		if err := database.DB.Save(&profile).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update color profile"})
			return
		}
	} else {
		// Create new profile
		if err := database.DB.Create(&profile).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create color profile"})
			return
		}
	}

	c.JSON(http.StatusCreated, profile)
}

// GetColorProfile retrieves a user's color profile
func (h *UserHandler) GetColorProfile(c *gin.Context) {
	userID := c.Param("userId")

	var profile models.ColorProfile
	if err := database.DB.First(&profile, "user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Color profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// ListUsers retrieves all users with pagination
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var users []models.User
	var total int64

	// Get total count
	database.DB.Model(&models.User{}).Count(&total)

	// Get paginated users
	if err := database.DB.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}
