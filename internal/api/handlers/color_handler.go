package handlers

import (
	"net/http"

	"fitgenie/internal/services"
	"fitgenie/pkg/logger"

	"github.com/gin-gonic/gin"
)

type ColorHandler struct {
	colorService *services.ColorTheoryService
	log          *logger.Logger
}

func NewColorHandler(colorService *services.ColorTheoryService, log *logger.Logger) *ColorHandler {
	return &ColorHandler{
		colorService: colorService,
		log:          log,
	}
}

func (h *ColorHandler) GetSeasons(c *gin.Context) {
	seasons := []string{"spring", "summer", "autumn", "winter"}
	c.JSON(http.StatusOK, gin.H{"seasons": seasons})
}

func (h *ColorHandler) GetHarmonies(c *gin.Context) {
	harmonies := []string{
		"monochromatic",
		"analogous",
		"complementary",
		"triadic",
		"split-complementary",
	}
	c.JSON(http.StatusOK, gin.H{"harmonies": harmonies})
}

func (h *ColorHandler) AnalyzeHarmony(c *gin.Context) {
	var req struct {
		Colors []string `json:"colors" binding:"required,min=2"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	score := h.colorService.AnalyzeOutfitColorHarmonyFromHex(req.Colors)

	c.JSON(http.StatusOK, gin.H{
		"harmony_score": score,
		"colors":        req.Colors,
	})
}

func (h *ColorHandler) GetRecommendations(c *gin.Context) {
	var req struct {
		BaseColor string `json:"base_color" binding:"required"`
		Harmony   string `json:"harmony" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	colors, err := h.colorService.GetHarmonyColors(req.BaseColor, req.Harmony)
	if err != nil {
		h.log.Error("failed to get harmony colors", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"base_color": req.BaseColor,
		"harmony":    req.Harmony,
		"colors":     colors,
	})
}
