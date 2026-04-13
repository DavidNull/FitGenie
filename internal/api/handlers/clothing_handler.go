package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"fitgenie/internal/models"
	"fitgenie/internal/repository"
	"fitgenie/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ClothingHandler struct {
	repo repository.ClothingRepository
	log  *logger.Logger
}

func NewClothingHandler(repo repository.ClothingRepository, log *logger.Logger) *ClothingHandler {
	return &ClothingHandler{
		repo: repo,
		log:  log,
	}
}

func (h *ClothingHandler) CreateClothing(c *gin.Context) {
	fmt.Printf("[DEBUG] CreateClothing called\n")
	var item models.ClothingItem
	if err := c.ShouldBindJSON(&item); err != nil {
		fmt.Printf("[ERROR] Failed to bind JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item.ID = uuid.New()

	// Get userID from context (set by DeviceAuthMiddleware)
	if userID, exists := c.Get("userID"); exists {
		item.UserID = userID.(uuid.UUID)
	}

	if err := h.repo.Create(c.Request.Context(), &item); err != nil {
		h.log.Error("failed to create clothing item", "error", err, "user_id", item.UserID, "item", item)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *ClothingHandler) GetClothing(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	item, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		h.log.Error("clothing item not found", "error", err, "id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Clothing item not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *ClothingHandler) ListClothing(c *gin.Context) {
	userIDStr := c.Query("user_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id", "received": userIDStr})
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

	items, total, err := h.repo.ListByUser(c.Request.Context(), userID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "user_id": userID.String()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     items,
		"items_len": len(items),
		"total":     total,
		"user_id":   userID.String(),
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func (h *ClothingHandler) UpdateClothing(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	item, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		h.log.Error("clothing item not found for update", "error", err, "id", id)
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

	if err := h.repo.Update(c.Request.Context(), item); err != nil {
		h.log.Error("failed to update clothing item", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update clothing item"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *ClothingHandler) DeleteClothing(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		h.log.Error("failed to delete clothing item", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete clothing item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Clothing item deleted successfully"})
}
