package handlers

import (
	"net/http"
	"strconv"

	"fitgenie/internal/models"
	"fitgenie/internal/repository"
	"fitgenie/internal/services"
	"fitgenie/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// OutfitHandler handles outfit-related API endpoints
type OutfitHandler struct {
	outfitRepo   repository.OutfitRepository
	clothingRepo repository.ClothingRepository
	userRepo     repository.UserRepository
	aiService    *services.AIService
	log          *logger.Logger
}

// NewOutfitHandler creates a new outfit handler with dependencies
func NewOutfitHandler(
	outfitRepo repository.OutfitRepository,
	clothingRepo repository.ClothingRepository,
	userRepo repository.UserRepository,
	aiService *services.AIService,
	log *logger.Logger,
) *OutfitHandler {
	return &OutfitHandler{
		outfitRepo:   outfitRepo,
		clothingRepo: clothingRepo,
		userRepo:     userRepo,
		aiService:    aiService,
		log:          log,
	}
}


// CreateOutfit creates a new outfit
func (h *OutfitHandler) CreateOutfit(c *gin.Context) {
	var outfit models.Outfit
	if err := c.ShouldBindJSON(&outfit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	outfit.ID = uuid.New()

	if err := h.outfitRepo.Create(c.Request.Context(), &outfit); err != nil {
		h.log.Error("failed to create outfit", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create outfit"})
		return
	}

	c.JSON(http.StatusCreated, outfit)
}

// GetOutfit retrieves an outfit by ID
func (h *OutfitHandler) GetOutfit(c *gin.Context) {
	outfitID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid outfit ID"})
		return
	}

	outfit, err := h.outfitRepo.GetByID(c.Request.Context(), outfitID)
	if err != nil {
		h.log.Error("outfit not found", "error", err, "outfit_id", outfitID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Outfit not found"})
		return
	}

	c.JSON(http.StatusOK, outfit)
}

// ListOutfits retrieves all outfits for a user
func (h *OutfitHandler) ListOutfits(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	outfits, total, err := h.outfitRepo.ListByUser(c.Request.Context(), userID, offset, limit)
	if err != nil {
		h.log.Error("failed to list outfits", "error", err)
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
	outfitID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid outfit ID"})
		return
	}

	if err := h.outfitRepo.Delete(c.Request.Context(), outfitID); err != nil {
		h.log.Error("failed to delete outfit", "error", err, "outfit_id", outfitID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete outfit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Outfit deleted successfully"})
}

// GetOutfitRecommendations generates AI-powered outfit recommendations
func (h *OutfitHandler) GetOutfitRecommendations(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var request services.OutfitRecommendationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userUUID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userClothingItems, _, err := h.clothingRepo.ListByUser(c.Request.Context(), userUUID, 0, 1000)
	if err != nil {
		h.log.Error("failed to retrieve user clothing items", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user clothing items"})
		return
	}

	if len(userClothingItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No clothing items found for user"})
		return
	}

	styleProfile, err := h.userRepo.GetStyleProfile(c.Request.Context(), userID)
	if err != nil {
		styleProfile = &models.StyleProfile{
			ID:     uuid.New(),
			UserID: userID,
		}
	}

	colorProfile, err := h.userRepo.GetColorProfile(c.Request.Context(), userID)
	if err != nil {
		colorProfile = &models.ColorProfile{
			ID:     uuid.New(),
			UserID: userID,
		}
	}

	recommendations, err := h.aiService.GenerateOutfitRecommendations(
		userID,
		userClothingItems,
		styleProfile,
	)
	if err != nil {
		h.log.Error("failed to generate recommendations", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate recommendations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recommendations": recommendations,
		"total":           len(recommendations),
	})
}


