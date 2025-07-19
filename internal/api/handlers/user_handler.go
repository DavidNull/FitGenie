package handlers

import (
	"net/http"
	"strconv"

	"fitgenie/internal/database"
	"fitgenie/internal/models"

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

//  adds an outfit to the user's favorites
func (h *UserHandler) AddFavoriteOutfit(c *gin.Context) {
	userID := c.Param("userId")
	outfitID := c.Param("outfitId")

	// Validate UUIDs
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	outfitUUID, err := uuid.Parse(outfitID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid outfit ID"})
		return
	}

	// Check if already favorited
	var existing models.FavoriteOutfit
	if err := database.DB.First(&existing, "user_id = ? AND outfit_id = ?", userUUID, outfitUUID).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Outfit already favorited"})
		return
	}

	favorite := models.FavoriteOutfit{
		UserID:   userUUID,
		OutfitID: outfitUUID,
	}
	if err := database.DB.Create(&favorite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add favorite"})
		return
	}

	c.JSON(http.StatusCreated, favorite)
}

//  removes an outfit from the user's favorites
func (h *UserHandler) RemoveFavoriteOutfit(c *gin.Context) {
	userID := c.Param("userId")
	outfitID := c.Param("outfitId")

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	outfitUUID, err := uuid.Parse(outfitID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid outfit ID"})
		return
	}

	if err := database.DB.Delete(&models.FavoriteOutfit{}, "user_id = ? AND outfit_id = ?", userUUID, outfitUUID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove favorite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Favorite removed successfully"})
}

// ListFavoriteOutfits lists all favorite outfits for a user
func (h *UserHandler) ListFavoriteOutfits(c *gin.Context) {
	userID := c.Param("userId")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var favorites []models.FavoriteOutfit
	if err := database.DB.Where("user_id = ?", userUUID).Find(&favorites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve favorites"})
		return
	}

	// Collect outfit IDs
	outfitIDs := make([]uuid.UUID, 0, len(favorites))
	for _, fav := range favorites {
		outfitIDs = append(outfitIDs, fav.OutfitID)
	}

	var outfits []models.Outfit
	if len(outfitIDs) > 0 {
		if err := database.DB.Where("id IN ?", outfitIDs).Find(&outfits).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve outfits"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"favorites": outfits})
}
