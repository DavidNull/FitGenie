package services

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"fitgenie/internal/models"

	"github.com/lucasb-eyer/go-colorful"
)

// ColorTheoryService handles color analysis and recommendations
type ColorTheoryService struct {
	colorSeasons   map[string]ColorSeason
	colorHarmonies map[string]ColorHarmony
}

// ColorSeason represents seasonal color analysis
type ColorSeason struct {
	Name        string
	Temperature string // warm, cool
	Intensity   string // light, deep, soft, clear
	Colors      []string
	Description string
}

// ColorHarmony represents color harmony rules
type ColorHarmony struct {
	Name        string
	Description string
	Rule        func(baseColor colorful.Color) []colorful.Color
}

// NewColorTheoryService creates a new color theory service
func NewColorTheoryService() *ColorTheoryService {
	service := &ColorTheoryService{
		colorSeasons:   make(map[string]ColorSeason),
		colorHarmonies: make(map[string]ColorHarmony),
	}

	service.initializeColorSeasons()
	service.initializeColorHarmonies()

	return service
}

// initializeColorSeasons sets up the seasonal color analysis system
func (c *ColorTheoryService) initializeColorSeasons() {
	c.colorSeasons = map[string]ColorSeason{
		"spring": {
			Name:        "Spring",
			Temperature: "warm",
			Intensity:   "clear",
			Colors: []string{
				"#FFB6C1", "#FF69B4", "#FF1493", "#DC143C", "#B22222",
				"#FF4500", "#FF8C00", "#FFA500", "#FFD700", "#FFFF00",
				"#ADFF2F", "#32CD32", "#00FF00", "#00FA9A", "#00CED1",
				"#87CEEB", "#4169E1", "#0000FF", "#8A2BE2", "#9400D3",
			},
			Description: "Bright, clear, warm colors with yellow undertones",
		},
		"summer": {
			Name:        "Summer",
			Temperature: "cool",
			Intensity:   "soft",
			Colors: []string{
				"#F0F8FF", "#E6E6FA", "#DDA0DD", "#DA70D6", "#BA55D3",
				"#9370DB", "#8B008B", "#4B0082", "#483D8B", "#6495ED",
				"#4682B4", "#5F9EA0", "#008B8B", "#2F4F4F", "#708090",
				"#778899", "#B0C4DE", "#D3D3D3", "#DCDCDC", "#F5F5F5",
			},
			Description: "Soft, muted colors with blue undertones",
		},
		"autumn": {
			Name:        "Autumn",
			Temperature: "warm",
			Intensity:   "deep",
			Colors: []string{
				"#8B4513", "#A0522D", "#CD853F", "#D2691E", "#FF8C00",
				"#FF7F50", "#DC143C", "#B22222", "#8B0000", "#800000",
				"#556B2F", "#6B8E23", "#808000", "#BDB76B", "#F0E68C",
				"#DAA520", "#B8860B", "#CD853F", "#DEB887", "#F4A460",
			},
			Description: "Rich, warm, earthy colors with golden undertones",
		},
		"winter": {
			Name:        "Winter",
			Temperature: "cool",
			Intensity:   "clear",
			Colors: []string{
				"#000000", "#2F4F4F", "#191970", "#000080", "#00008B",
				"#0000CD", "#0000FF", "#4169E1", "#6495ED", "#87CEEB",
				"#B0E0E6", "#FFFFFF", "#F8F8FF", "#E6E6FA", "#DDA0DD",
				"#FF1493", "#DC143C", "#B22222", "#8B0000", "#800000",
			},
			Description: "Clear, intense colors with blue undertones",
		},
	}
}

// initializeColorHarmonies sets up color harmony rules
func (c *ColorTheoryService) initializeColorHarmonies() {
	c.colorHarmonies = map[string]ColorHarmony{
		"monochromatic": {
			Name:        "Monochromatic",
			Description: "Different shades and tints of the same color",
			Rule: func(baseColor colorful.Color) []colorful.Color {
				h, s, l := baseColor.Hsl()
				colors := []colorful.Color{baseColor}

				// Add lighter and darker versions
				for i := 1; i <= 3; i++ {
					lighter := colorful.Hsl(h, s, math.Min(1.0, l+float64(i)*0.15))
					darker := colorful.Hsl(h, s, math.Max(0.0, l-float64(i)*0.15))
					colors = append(colors, lighter, darker)
				}

				return colors
			},
		},
		"complementary": {
			Name:        "Complementary",
			Description: "Colors opposite on the color wheel",
			Rule: func(baseColor colorful.Color) []colorful.Color {
				h, s, l := baseColor.Hsl()
				complementary := colorful.Hsl(math.Mod(h+180, 360), s, l)
				return []colorful.Color{baseColor, complementary}
			},
		},
		"triadic": {
			Name:        "Triadic",
			Description: "Three colors evenly spaced on the color wheel",
			Rule: func(baseColor colorful.Color) []colorful.Color {
				h, s, l := baseColor.Hsl()
				color2 := colorful.Hsl(math.Mod(h+120, 360), s, l)
				color3 := colorful.Hsl(math.Mod(h+240, 360), s, l)
				return []colorful.Color{baseColor, color2, color3}
			},
		},
		"analogous": {
			Name:        "Analogous",
			Description: "Colors adjacent on the color wheel",
			Rule: func(baseColor colorful.Color) []colorful.Color {
				h, s, l := baseColor.Hsl()
				color2 := colorful.Hsl(math.Mod(h+30, 360), s, l)
				color3 := colorful.Hsl(math.Mod(h-30, 360), s, l)
				return []colorful.Color{baseColor, color2, color3}
			},
		},
		"split_complementary": {
			Name:        "Split Complementary",
			Description: "Base color plus two colors adjacent to its complement",
			Rule: func(baseColor colorful.Color) []colorful.Color {
				h, s, l := baseColor.Hsl()
				color2 := colorful.Hsl(math.Mod(h+150, 360), s, l)
				color3 := colorful.Hsl(math.Mod(h+210, 360), s, l)
				return []colorful.Color{baseColor, color2, color3}
			},
		},
	}
}

// AnalyzeImageColors extracts and analyzes colors from clothing item image
func (c *ColorTheoryService) AnalyzeImageColors(imageData []byte) (*models.ColorAnalysis, error) {
	// This would integrate with image processing library
	// For now, we'll simulate the analysis

	dominantColors := []models.ColorInfo{
		{
			Hex:        "#FF6B6B",
			Name:       "Coral Red",
			Percentage: 45.2,
			HSL:        models.HSL{H: 0, S: 100, L: 70},
			RGB:        models.RGB{R: 255, G: 107, B: 107},
		},
		{
			Hex:        "#4ECDC4",
			Name:       "Turquoise",
			Percentage: 32.8,
			HSL:        models.HSL{H: 174, S: 65, L: 60},
			RGB:        models.RGB{R: 78, G: 205, B: 196},
		},
		{
			Hex:        "#45B7D1",
			Name:       "Sky Blue",
			Percentage: 22.0,
			HSL:        models.HSL{H: 200, S: 60, L: 55},
			RGB:        models.RGB{R: 69, G: 183, B: 209},
		},
	}

	analysis := &models.ColorAnalysis{
		DominantColors:   dominantColors,
		ColorHarmony:     c.determineColorHarmony(dominantColors),
		ColorTemperature: c.determineColorTemperature(dominantColors),
		ColorSaturation:  c.determineColorSaturation(dominantColors),
		ColorBrightness:  c.determineColorBrightness(dominantColors),
	}

	return analysis, nil
}

// DetermineColorSeason analyzes user's color profile to determine their color season
func (c *ColorTheoryService) DetermineColorSeason(skinTone, undertone string, favoriteColors []string) string {
	scores := make(map[string]float64)

	// Analyze based on undertone
	for seasonName, season := range c.colorSeasons {
		score := 0.0

		// Undertone matching
		if (undertone == "warm" || undertone == "yellow") && season.Temperature == "warm" {
			score += 3.0
		} else if (undertone == "cool" || undertone == "pink") && season.Temperature == "cool" {
			score += 3.0
		} else if undertone == "neutral" {
			score += 1.5
		}

		// Skin tone matching
		if skinTone == "light" && (seasonName == "spring" || seasonName == "summer") {
			score += 2.0
		} else if skinTone == "medium" {
			score += 1.5
		} else if skinTone == "dark" && (seasonName == "autumn" || seasonName == "winter") {
			score += 2.0
		}

		// Favorite colors analysis
		for _, favColor := range favoriteColors {
			for _, seasonColor := range season.Colors {
				if c.colorsAreSimilar(favColor, seasonColor) {
					score += 1.0
				}
			}
		}

		scores[seasonName] = score
	}

	// Find the season with highest score
	maxScore := 0.0
	bestSeason := "spring"
	for season, score := range scores {
		if score > maxScore {
			maxScore = score
			bestSeason = season
		}
	}

	return bestSeason
}

// GenerateColorRecommendations creates personalized color recommendations
func (c *ColorTheoryService) GenerateColorRecommendations(colorProfile *models.ColorProfile) ([]models.ColorRecommendation, error) {
	season := c.colorSeasons[strings.ToLower(colorProfile.ColorSeason)]
	recommendations := []models.ColorRecommendation{}

	// Primary colors (best matches)
	primaryColors := season.Colors[:5]
	for i, colorHex := range primaryColors {
		rec := models.ColorRecommendation{
			ColorHex:   colorHex,
			ColorName:  c.getColorName(colorHex),
			Category:   "Primary",
			Confidence: 0.9 - float64(i)*0.1,
			Reason:     fmt.Sprintf("Perfect match for %s season palette", season.Name),
		}
		recommendations = append(recommendations, rec)
	}

	// Secondary colors (good matches)
	secondaryColors := season.Colors[5:10]
	for i, colorHex := range secondaryColors {
		rec := models.ColorRecommendation{
			ColorHex:   colorHex,
			ColorName:  c.getColorName(colorHex),
			Category:   "Secondary",
			Confidence: 0.7 - float64(i)*0.05,
			Reason:     fmt.Sprintf("Good complement for %s season", season.Name),
		}
		recommendations = append(recommendations, rec)
	}

	// Neutral colors
	neutrals := []string{"#FFFFFF", "#F5F5F5", "#D3D3D3", "#808080", "#2F2F2F", "#000000"}
	for i, colorHex := range neutrals {
		rec := models.ColorRecommendation{
			ColorHex:   colorHex,
			ColorName:  c.getColorName(colorHex),
			Category:   "Neutral",
			Confidence: 0.8,
			Reason:     "Universal neutral that works with any season",
		}
		recommendations = append(recommendations, rec)
		if i >= 2 { // Limit neutrals
			break
		}
	}

	return recommendations, nil
}

// AnalyzeOutfitColorHarmony evaluates color harmony in an outfit
func (c *ColorTheoryService) AnalyzeOutfitColorHarmony(clothingItems []models.ClothingItem) float64 {
	if len(clothingItems) < 2 {
		return 1.0 // Single item is always harmonious
	}

	colors := []colorful.Color{}
	for _, item := range clothingItems {
		if item.PrimaryColor != "" {
			if color, err := colorful.Hex(item.PrimaryColor); err == nil {
				colors = append(colors, color)
			}
		}
	}

	if len(colors) < 2 {
		return 0.5
	}

	// Check for known harmonious combinations
	harmonyScore := 0.0
	totalComparisons := 0

	for i := 0; i < len(colors); i++ {
		for j := i + 1; j < len(colors); j++ {
			score := c.calculateColorHarmonyScore(colors[i], colors[j])
			harmonyScore += score
			totalComparisons++
		}
	}

	if totalComparisons == 0 {
		return 0.5
	}

	return harmonyScore / float64(totalComparisons)
}

// Helper methods

func (c *ColorTheoryService) determineColorHarmony(colors []models.ColorInfo) string {
	if len(colors) < 2 {
		return "monochromatic"
	}

	// Convert to colorful.Color for analysis
	colorfulColors := []colorful.Color{}
	for _, color := range colors {
		if col, err := colorful.Hex(color.Hex); err == nil {
			colorfulColors = append(colorfulColors, col)
		}
	}

	// Check for complementary colors
	for i := 0; i < len(colorfulColors); i++ {
		for j := i + 1; j < len(colorfulColors); j++ {
			h1, _, _ := colorfulColors[i].Hsl()
			h2, _, _ := colorfulColors[j].Hsl()

			hueDiff := math.Abs(h1 - h2)
			if hueDiff > 180 {
				hueDiff = 360 - hueDiff
			}

			if hueDiff > 150 && hueDiff < 210 {
				return "complementary"
			}
		}
	}

	// Check for analogous colors
	hues := []float64{}
	for _, color := range colorfulColors {
		h, _, _ := color.Hsl()
		hues = append(hues, h)
	}

	sort.Float64s(hues)
	maxDiff := 0.0
	for i := 1; i < len(hues); i++ {
		diff := hues[i] - hues[i-1]
		if diff > maxDiff {
			maxDiff = diff
		}
	}

	if maxDiff < 60 {
		return "analogous"
	}

	return "mixed"
}

func (c *ColorTheoryService) determineColorTemperature(colors []models.ColorInfo) string {
	warmCount := 0
	coolCount := 0

	for _, color := range colors {
		if col, err := colorful.Hex(color.Hex); err == nil {
			h, _, _ := col.Hsl()

			// Warm colors: red, orange, yellow (0-60, 300-360)
			// Cool colors: green, blue, purple (60-300)
			if (h >= 0 && h <= 60) || (h >= 300 && h <= 360) {
				warmCount++
			} else {
				coolCount++
			}
		}
	}

	if warmCount > coolCount {
		return "warm"
	} else if coolCount > warmCount {
		return "cool"
	}
	return "neutral"
}

func (c *ColorTheoryService) determineColorSaturation(colors []models.ColorInfo) string {
	totalSaturation := 0.0
	count := 0

	for _, color := range colors {
		if col, err := colorful.Hex(color.Hex); err == nil {
			_, s, _ := col.Hsl()
			totalSaturation += s
			count++
		}
	}

	if count == 0 {
		return "medium"
	}

	avgSaturation := totalSaturation / float64(count)

	if avgSaturation > 0.7 {
		return "high"
	} else if avgSaturation < 0.3 {
		return "low"
	}
	return "medium"
}

func (c *ColorTheoryService) determineColorBrightness(colors []models.ColorInfo) string {
	totalLightness := 0.0
	count := 0

	for _, color := range colors {
		if col, err := colorful.Hex(color.Hex); err == nil {
			_, _, l := col.Hsl()
			totalLightness += l
			count++
		}
	}

	if count == 0 {
		return "medium"
	}

	avgLightness := totalLightness / float64(count)

	if avgLightness > 0.7 {
		return "bright"
	} else if avgLightness < 0.3 {
		return "dark"
	}
	return "medium"
}

func (c *ColorTheoryService) colorsAreSimilar(color1, color2 string) bool {
	col1, err1 := colorful.Hex(color1)
	col2, err2 := colorful.Hex(color2)

	if err1 != nil || err2 != nil {
		return false
	}

	// Calculate color difference using Delta E
	return col1.DistanceCIE76(col2) < 30 // Threshold for similarity
}

func (c *ColorTheoryService) calculateColorHarmonyScore(color1, color2 colorful.Color) float64 {
	h1, s1, l1 := color1.Hsl()
	h2, s2, l2 := color2.Hsl()

	// Calculate hue difference
	hueDiff := math.Abs(h1 - h2)
	if hueDiff > 180 {
		hueDiff = 360 - hueDiff
	}

	// Score based on harmonic relationships
	score := 0.0

	// Complementary (opposite colors)
	if hueDiff > 150 && hueDiff < 210 {
		score = 0.9
	}
	// Triadic (120 degrees apart)
	if (hueDiff > 110 && hueDiff < 130) || (hueDiff > 230 && hueDiff < 250) {
		score = 0.8
	}
	// Analogous (adjacent colors)
	if hueDiff < 30 {
		score = 0.7
	}
	// Split complementary
	if (hueDiff > 140 && hueDiff < 160) || (hueDiff > 200 && hueDiff < 220) {
		score = 0.75
	}

	// Adjust score based on saturation and lightness similarity
	satDiff := math.Abs(s1 - s2)
	lightDiff := math.Abs(l1 - l2)

	if satDiff < 0.2 && lightDiff < 0.2 {
		score += 0.1 // Bonus for similar saturation and lightness
	}

	return math.Min(1.0, score)
}

func (c *ColorTheoryService) getColorName(hex string) string {
	// Simple color name mapping - in production, use a comprehensive color name database
	colorNames := map[string]string{
		"#FF0000": "Red",
		"#00FF00": "Green",
		"#0000FF": "Blue",
		"#FFFF00": "Yellow",
		"#FF00FF": "Magenta",
		"#00FFFF": "Cyan",
		"#FFFFFF": "White",
		"#000000": "Black",
		"#808080": "Gray",
		"#FFA500": "Orange",
		"#800080": "Purple",
		"#FFC0CB": "Pink",
		"#A52A2A": "Brown",
		"#F5F5DC": "Beige",
		"#000080": "Navy",
	}

	if name, exists := colorNames[strings.ToUpper(hex)]; exists {
		return name
	}

	return "Unknown Color"
}

// AnalyzeOutfitColorHarmonyFromHex analyzes color harmony from hex color strings
func (c *ColorTheoryService) AnalyzeOutfitColorHarmonyFromHex(colors []string) float64 {
	if len(colors) < 2 {
		return 1.0
	}

	colorfulColors := []colorful.Color{}
	for _, hex := range colors {
		if col, err := colorful.Hex(hex); err == nil {
			colorfulColors = append(colorfulColors, col)
		}
	}

	if len(colorfulColors) < 2 {
		return 0.5
	}

	harmonyScore := 0.0
	totalComparisons := 0

	for i := 0; i < len(colorfulColors); i++ {
		for j := i + 1; j < len(colorfulColors); j++ {
			score := c.calculateColorHarmonyScore(colorfulColors[i], colorfulColors[j])
			harmonyScore += score
			totalComparisons++
		}
	}

	if totalComparisons == 0 {
		return 0.5
	}

	return harmonyScore / float64(totalComparisons)
}

// GetHarmonyColors returns colors that harmonize with the base color
func (c *ColorTheoryService) GetHarmonyColors(baseHex, harmonyType string) ([]string, error) {
	baseColor, err := colorful.Hex(baseHex)
	if err != nil {
		return nil, err
	}

	harmony, exists := c.colorHarmonies[harmonyType]
	if !exists {
		return nil, fmt.Errorf("unknown harmony type: %s", harmonyType)
	}

	colors := harmony.Rule(baseColor)
	hexColors := make([]string, len(colors))
	for i, col := range colors {
		hexColors[i] = col.Hex()
	}

	return hexColors, nil
}
