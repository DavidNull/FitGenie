package handlers

import (
	"net/http"
	"strconv"

	"fitgenie/internal/models"
	"fitgenie/internal/repository"
	"fitgenie/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserHandler handles user-related API endpoints
type UserHandler struct {
	repo repository.UserRepository
	log  *logger.Logger
}

// NewUserHandler creates a new user handler with dependencies
func NewUserHandler(repo repository.UserRepository, log *logger.Logger) *UserHandler {
	return &UserHandler{
		repo: repo,
		log:  log,
	}
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = uuid.New()

	if err := h.repo.Create(c.Request.Context(), &user); err != nil {
		h.log.Error("failed to create user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUser retrieves a user by ID
func (h *UserHandler) GetUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.repo.GetByID(c.Request.Context(), userID)
	if err != nil {
		h.log.Error("failed to get user", "error", err, "user_id", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser updates an existing user
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.repo.GetByID(c.Request.Context(), userID)
	if err != nil {
		h.log.Error("failed to get user for update", "error", err, "user_id", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var updateData models.User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Name = updateData.Name
	user.Email = updateData.Email

	if err := h.repo.Update(c.Request.Context(), user); err != nil {
		h.log.Error("failed to update user", "error", err, "user_id", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), userID); err != nil {
		h.log.Error("failed to delete user", "error", err, "user_id", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// CreateStyleProfile creates or updates a user's style profile
func (h *UserHandler) CreateStyleProfile(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Verify user exists
	_, err = h.repo.GetByID(c.Request.Context(), userID)
	if err != nil {
		h.log.Error("user not found for style profile", "error", err, "user_id", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var profile models.StyleProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile.ID = uuid.New()
	profile.UserID = userID

	if err := h.repo.CreateStyleProfile(c.Request.Context(), &profile); err != nil {
		h.log.Error("failed to create style profile", "error", err, "user_id", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create style profile"})
		return
	}

	c.JSON(http.StatusCreated, profile)
}

// GetStyleProfile retrieves a user's style profile
func (h *UserHandler) GetStyleProfile(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	profile, err := h.repo.GetStyleProfile(c.Request.Context(), userID)
	if err != nil {
		h.log.Error("style profile not found", "error", err, "user_id", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Style profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// CreateColorProfile creates or updates a user's color profile
func (h *UserHandler) CreateColorProfile(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Verify user exists
	_, err = h.repo.GetByID(c.Request.Context(), userID)
	if err != nil {
		h.log.Error("user not found for color profile", "error", err, "user_id", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var profile models.ColorProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile.ID = uuid.New()
	profile.UserID = userID

	if err := h.repo.CreateColorProfile(c.Request.Context(), &profile); err != nil {
		h.log.Error("failed to create color profile", "error", err, "user_id", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create color profile"})
		return
	}

	c.JSON(http.StatusCreated, profile)
}

// GetColorProfile retrieves a user's color profile
func (h *UserHandler) GetColorProfile(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	profile, err := h.repo.GetColorProfile(c.Request.Context(), userID)
	if err != nil {
		h.log.Error("color profile not found", "error", err, "user_id", userID)
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

	users, total, err := h.repo.List(c.Request.Context(), offset, limit)
	if err != nil {
		h.log.Error("failed to list users", "error", err)
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

// AddFavoriteOutfit adds an outfit to the user's favorites
func (h *UserHandler) AddFavoriteOutfit(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	outfitID, err := uuid.Parse(c.Param("outfitId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid outfit ID"})
		return
	}

	favorite := &models.FavoriteOutfit{
		ID:       uuid.New(),
		UserID:   userID,
		OutfitID: outfitID,
	}

	if err := h.repo.AddFavoriteOutfit(c.Request.Context(), favorite); err != nil {
		h.log.Error("failed to add favorite", "error", err, "user_id", userID, "outfit_id", outfitID)
		c.JSON(http.StatusConflict, gin.H{"error": "Outfit already favorited"})
		return
	}

	c.JSON(http.StatusCreated, favorite)
}

// RemoveFavoriteOutfit removes an outfit from the user's favorites
func (h *UserHandler) RemoveFavoriteOutfit(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	outfitID, err := uuid.Parse(c.Param("outfitId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid outfit ID"})
		return
	}

	if err := h.repo.RemoveFavoriteOutfit(c.Request.Context(), userID, outfitID); err != nil {
		h.log.Error("failed to remove favorite", "error", err, "user_id", userID, "outfit_id", outfitID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove favorite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Favorite removed successfully"})
}

// ListFavoriteOutfits lists all favorite outfits for a user
func (h *UserHandler) ListFavoriteOutfits(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	favorites, err := h.repo.ListFavoriteOutfits(c.Request.Context(), userID)
	if err != nil {
		h.log.Error("failed to list favorites", "error", err, "user_id", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve favorites"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"favorites": favorites})
}
