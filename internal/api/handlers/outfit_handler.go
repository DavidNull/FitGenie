package handlers

import (
	"fitgenie/internal/database"
	"fitgenie/internal/models"
	"fitgenie/internal/services"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
)

// OutfitHandler handles outfit-related API endpoints
type OutfitHandler struct {
	aiService *services.AIService
	userService *services.UserService
}

// NewOutfitHandler creates a new outfit handler
func NewOutfitHandler(aiService *services.AIService, userService *services.UserService) *OutfitHandler {
	return &OutfitHandler{
		aiService: aiService,
		userService: userService,
	}
}

// CreateClothingItem creates a new clothing item
func (h *OutfitHandler) CreateClothingItem(c *gin.Context) {
	userID := c.Param("userId")
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var item models.ClothingItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item.ID = uuid.New()
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	item.UserID = userUUID

	if err := database.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create clothing item"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// GetClothingItem retrieves a clothing item by ID
func (h *OutfitHandler) GetClothingItem(c *gin.Context) {
	itemID := c.Param("id")

	var item models.ClothingItem
	if err := database.DB.First(&item, "id = ?", itemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Clothing item not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// UpdateClothingItem updates an existing clothing item
func (h *OutfitHandler) UpdateClothingItem(c *gin.Context) {
	itemID := c.Param("id")

	var item models.ClothingItem
	if err := database.DB.First(&item, "id = ?", itemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Clothing item not found"})
		return
	}

	var updateData models.ClothingItem
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item.Name = updateData.Name
	item.Category = updateData.Category
	item.Brand = updateData.Brand
	item.Size = updateData.Size
	item.PrimaryColor = updateData.PrimaryColor
	item.SecondaryColor = updateData.SecondaryColor
	item.Material = updateData.Material
	item.Style = updateData.Style
	item.Season = updateData.Season
	item.Occasion = updateData.Occasion
	item.Notes = updateData.Notes

	if err := database.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update clothing item"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DeleteClothingItem deletes a clothing item
func (h *OutfitHandler) DeleteClothingItem(c *gin.Context) {
	itemID := c.Param("id")

	if err := database.DB.Delete(&models.ClothingItem{}, "id = ?", itemID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete clothing item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Clothing item deleted successfully"})
}

// ListClothingItems retrieves all clothing items for a user
func (h *OutfitHandler) ListClothingItems(c *gin.Context) {
	userID := c.Param("userId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	var items []models.ClothingItem
	var total int64

	query := database.DB.Where("user_id = ?", userID)

	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}
	if style := c.Query("style"); style != "" {
		query = query.Where("style = ?", style)
	}
	if color := c.Query("color"); color != "" {
		query = query.Where("primary_color = ? OR secondary_color = ?", color, color)
	}

	query.Model(&models.ClothingItem{}).Count(&total)

	if err := query.Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve clothing items"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// CreateOutfit creates a new outfit
func (h *OutfitHandler) CreateOutfit(c *gin.Context) {
	userID := c.Param("userId")

	var outfit models.Outfit
	if err := c.ShouldBindJSON(&outfit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	outfit.ID = uuid.New()
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	outfit.UserID = userUUID

	if err := database.DB.Create(&outfit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create outfit"})
		return
	}

	c.JSON(http.StatusCreated, outfit)
}

// GetOutfit retrieves an outfit by ID
func (h *OutfitHandler) GetOutfit(c *gin.Context) {
	outfitID := c.Param("id")

	var outfit models.Outfit
	if err := database.DB.Preload("ClothingItems").First(&outfit, "id = ?", outfitID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Outfit not found"})
		return
	}

	c.JSON(http.StatusOK, outfit)
}

// ListOutfits retrieves all outfits for a user
func (h *OutfitHandler) ListOutfits(c *gin.Context) {
	userID := c.Param("userId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var outfits []models.Outfit
	var total int64

	query := database.DB.Where("user_id = ?", userID)

	if style := c.Query("style"); style != "" {
		query = query.Where("style = ?", style)
	}

	query.Model(&models.Outfit{}).Count(&total)

	if err := query.Preload("ClothingItems").Offset(offset).Limit(limit).Find(&outfits).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve outfits"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"outfits": outfits,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// DeleteOutfit deletes an outfit
func (h *OutfitHandler) DeleteOutfit(c *gin.Context) {
	outfitID := c.Param("id")

	if err := database.DB.Delete(&models.Outfit{}, "id = ?", outfitID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete outfit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Outfit deleted successfully"})
}

// GetOutfitRecommendations generates AI-powered outfit recommendations
func (h *OutfitHandler) GetOutfitRecommendations(c *gin.Context) {
	userID := c.Param("userId")

	var request services.OutfitRecommendationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request.UserID = userID

	var userClothingItems []models.ClothingItem
	if err := database.DB.Where("user_id = ?", userID).Find(&userClothingItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user clothing items"})
		return
	}

	if len(userClothingItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No clothing items found for user"})
		return
	}

	var styleProfile models.StyleProfile
	if err := database.DB.First(&styleProfile, "user_id = ?", userID).Error; err != nil {
		styleProfile = models.StyleProfile{
			ID:     uuid.New(),
			UserID: uuid.MustParse(userID),
		}
		database.DB.Create(&styleProfile)
	}

	var colorProfile models.ColorProfile
	if err := database.DB.First(&colorProfile, "user_id = ?", userID).Error; err != nil {
		colorProfile = models.ColorProfile{
			ID:     uuid.New(),
			UserID: uuid.MustParse(userID),
		}
		database.DB.Create(&colorProfile)
	}

	recommendations, err := h.aiService.GenerateOutfitRecommendations(
		request,
		userClothingItems,
		&styleProfile,
		&colorProfile,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate recommendations"})
		return
	}

	for i := range recommendations {
		if err := database.DB.Create(&recommendations[i]).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save recommendations"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"recommendations": recommendations,
		"total":           len(recommendations),
	})
}

// GetPersonalizedRecommendations gets personalized wardrobe recommendations
func (h *OutfitHandler) GetPersonalizedRecommendations(c *gin.Context) {
	userID := c.Param("userId")

	var styleProfile models.StyleProfile
	if err := database.DB.First(&styleProfile, "user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Style profile not found"})
		return
	}

	var colorProfile models.ColorProfile
	if err := database.DB.First(&colorProfile, "user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Color profile not found"})
		return
	}

	var clothingItems []models.ClothingItem
	if err := database.DB.Where("user_id = ?", userID).Find(&clothingItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve clothing items"})
		return
	}

	recommendations, err := h.aiService.GeneratePersonalizedRecommendations(
		&styleProfile,
		&colorProfile,
		clothingItems,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate personalized recommendations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recommendations": recommendations,
		"total":           len(recommendations),
	})
}

// MarkRecommendationViewed marks a recommendation as viewed
func (h *OutfitHandler) MarkRecommendationViewed(c *gin.Context) {
	recommendationID := c.Param("id")

	var recommendation models.OutfitRecommendation
	if err := database.DB.First(&recommendation, "id = ?", recommendationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recommendation not found"})
		return
	}

	recommendation.Viewed = true
	if err := database.DB.Save(&recommendation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update recommendation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recommendation marked as viewed"})
}

// AcceptRecommendation marks a recommendation as accepted
func (h *OutfitHandler) AcceptRecommendation(c *gin.Context) {
	recommendationID := c.Param("id")

	var recommendation models.OutfitRecommendation
	if err := database.DB.First(&recommendation, "id = ?", recommendationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recommendation not found"})
		return
	}

	recommendation.Accepted = true
	recommendation.Viewed = true
	if err := database.DB.Save(&recommendation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update recommendation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recommendation accepted"})
}
