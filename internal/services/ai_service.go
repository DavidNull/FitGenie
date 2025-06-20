package services

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"fitgenie/internal/models"

	"github.com/google/uuid"
	"github.com/lucasb-eyer/go-colorful"
)

// AIService combines color theory and style analysis for intelligent recommendations
type AIService struct {
	colorService *ColorTheoryService
	styleService *StyleService
}

// OutfitRecommendationRequest represents a request for outfit recommendations
type OutfitRecommendationRequest struct {
	UserID   string
	Occasion string
	Season   string
	Weather  string
	Style    string
	Colors   []string
	MaxItems int
}

// OutfitScore represents the scoring of an outfit combination
type OutfitScore struct {
	Outfit              *models.Outfit
	ColorHarmonyScore   float64
	StyleCoherenceScore float64
	SeasonScore         float64
	OccasionScore       float64
	OverallScore        float64
	Reasoning           string
}

// OutfitSuggestion represents a suggested outfit
type OutfitSuggestion struct {
	Items      []models.ClothingItem `json:"items"`
	Confidence float64               `json:"confidence"`
	Reasoning  string                `json:"reasoning"`
	StyleScore float64               `json:"style_score"`
	ColorScore float64               `json:"color_score"`
	Occasion   string                `json:"occasion"`
	Season     string                `json:"season"`
}

// NewAIService creates a new AI service
func NewAIService(colorService *ColorTheoryService, styleService *StyleService) *AIService {
	return &AIService{
		colorService: colorService,
		styleService: styleService,
	}
}

//  creates AI-powered outfit recommendations
func (ai *AIService) GenerateOutfitRecommendations(userID uuid.UUID, items []models.ClothingItem, preferences *models.StyleProfile) ([]models.OutfitRecommendation, error) {
	if len(items) < 2 {
		return nil, fmt.Errorf("need at least 2 items to create outfit recommendations")
	}

	recommendations := []models.OutfitRecommendation{}
	maxRecommendations := 5

	for i := 0; i < maxRecommendations && len(recommendations) < maxRecommendations; i++ {
		outfit, err := ai.generateSingleOutfit(items, preferences)
		if err != nil {
			continue
		}

		if ai.isOutfitUnique(outfit, recommendations) {
			recommendation := models.OutfitRecommendation{
				ID:         uuid.New(),
				OutfitID:   outfit.Items[0].ID,
				Confidence: outfit.Confidence,
				Reasoning:  outfit.Reasoning,
				CreatedAt:  time.Now(),
			}
			recommendations = append(recommendations, recommendation)
		}
	}

	if len(recommendations) == 0 {
		return nil, fmt.Errorf("could not generate any suitable outfit recommendations")
	}

	return recommendations, nil
}

// AnalyzeOutfitCompatibility analyzes how well clothing items work together
func (ai *AIService) AnalyzeOutfitCompatibility(items []models.ClothingItem) (float64, string, error) {
	if len(items) == 0 {
		return 0.0, "No items provided for analysis", nil
	}

	if len(items) == 1 {
		return 1.0, "Single item - perfect compatibility", nil
	}

	styleScore := ai.styleService.AnalyzeOutfitStyleCoherence(items)
	colorScore := ai.analyzeColorHarmony(items)

	overallScore := (styleScore*0.6 + colorScore*0.4)

	reasoning := ai.generateCompatibilityReasoning(styleScore, colorScore, items)

	return overallScore, reasoning, nil
}

// GeneratePersonalizedRecommendations creates recommendations based on user preferences
func (ai *AIService) GeneratePersonalizedRecommendations(
	userProfile *models.StyleProfile,
	colorProfile *models.ColorProfile,
	clothingItems []models.ClothingItem,
) ([]string, error) {

	recommendations := []string{}

	// Color-based recommendations
	colorRecs, err := ai.colorService.GenerateColorRecommendations(colorProfile)
	if err == nil {
		for i, rec := range colorRecs {
			if i < 3 { // Limit to top 3
				recommendations = append(recommendations, fmt.Sprintf(
					"Consider adding %s (%s) to your wardrobe - %s",
					rec.ColorName, rec.ColorHex, rec.Reason,
				))
			}
		}
	}

	// Style-based recommendations
	styleRecs, err := ai.styleService.GenerateStyleRecommendations(userProfile)
	if err == nil {
		recommendations = append(recommendations, styleRecs...)
	}

	// Wardrobe gap analysis
	gapAnalysis := ai.analyzeWardrobeGaps(clothingItems, userProfile)
	recommendations = append(recommendations, gapAnalysis...)

	// Versatility recommendations
	versatilityRecs := ai.generateVersatilityRecommendations(clothingItems)
	recommendations = append(recommendations, versatilityRecs...)

	return recommendations, nil
}

// Helper methods

func (ai *AIService) filterClothingItems(items []models.ClothingItem, request OutfitRecommendationRequest) []models.ClothingItem {
	filtered := []models.ClothingItem{}

	for _, item := range items {
		include := true

		// Filter by season
		if request.Season != "" {
			seasonMatch := false
			for _, season := range item.Season {
				if season == request.Season {
					seasonMatch = true
					break
				}
			}
			if !seasonMatch && len(item.Season) > 0 {
				include = false
			}
		}

		// Filter by occasion
		if request.Occasion != "" {
			occasionMatch := false
			for _, occasion := range item.Occasion {
				if occasion == request.Occasion {
					occasionMatch = true
					break
				}
			}
			if !occasionMatch && len(item.Occasion) > 0 {
				include = false
			}
		}

		// Filter by style
		if request.Style != "" && item.Style != "" {
			if item.Style != request.Style {
				include = false
			}
		}

		// Filter by colors (if specific colors requested)
		if len(request.Colors) > 0 {
			colorMatch := false
			for _, requestedColor := range request.Colors {
				if item.PrimaryColor == requestedColor || item.SecondaryColor == requestedColor {
					colorMatch = true
					break
				}
			}
			if !colorMatch {
				include = false
			}
		}

		if include {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

func (ai *AIService) generateOutfitCombinations(items []models.ClothingItem, maxItems int) [][]models.ClothingItem {
	if maxItems == 0 {
		maxItems = 4 // Default max items per outfit
	}

	// Group items by category
	categories := make(map[string][]models.ClothingItem)
	for _, item := range items {
		categories[item.Category] = append(categories[item.Category], item)
	}

	combinations := [][]models.ClothingItem{}

	// Generate basic outfit templates
	templates := [][]string{
		{"shirt", "pants", "shoes"},
		{"dress", "shoes"},
		{"shirt", "skirt", "shoes"},
		{"sweater", "jeans", "shoes"},
		{"blazer", "shirt", "pants", "shoes"},
		{"t-shirt", "jeans", "jacket", "shoes"},
	}

	for _, template := range templates {
		combination := ai.buildCombinationFromTemplate(template, categories)
		if len(combination) >= 2 && len(combination) <= maxItems {
			combinations = append(combinations, combination)
		}
	}

	// Generate additional random combinations if needed
	if len(combinations) < 20 {
		additionalCombos := ai.generateRandomCombinations(items, maxItems, 20-len(combinations))
		combinations = append(combinations, additionalCombos...)
	}

	return combinations
}

func (ai *AIService) buildCombinationFromTemplate(template []string, categories map[string][]models.ClothingItem) []models.ClothingItem {
	combination := []models.ClothingItem{}

	for _, category := range template {
		if items, exists := categories[category]; exists && len(items) > 0 {
			// Pick the first item from this category (could be randomized)
			combination = append(combination, items[0])
		}
	}

	return combination
}

func (ai *AIService) generateRandomCombinations(items []models.ClothingItem, maxItems, count int) [][]models.ClothingItem {
	combinations := [][]models.ClothingItem{}

	for i := 0; i < count && i < 50; i++ { // Limit to prevent infinite loops
		combination := []models.ClothingItem{}
		usedCategories := make(map[string]bool)

		// Randomly select items ensuring no duplicate categories
		for len(combination) < maxItems && len(combination) < len(items) {
			for _, item := range items {
				if !usedCategories[item.Category] && len(combination) < maxItems {
					combination = append(combination, item)
					usedCategories[item.Category] = true
				}
			}
			break // Prevent infinite loop
		}

		if len(combination) >= 2 {
			combinations = append(combinations, combination)
		}
	}

	return combinations
}

func (ai *AIService) scoreOutfitCombination(
	items []models.ClothingItem,
	request OutfitRecommendationRequest,
	userProfile *models.StyleProfile,
	colorProfile *models.ColorProfile,
) OutfitScore {

	outfit := &models.Outfit{
		ClothingItems: items,
	}

	// Score color harmony
	colorScore := ai.colorService.AnalyzeOutfitColorHarmony(items)

	// Score style coherence
	styleScore := ai.styleService.AnalyzeOutfitStyleCoherence(items)

	// Score season appropriateness
	seasonScore := ai.scoreSeasonAppropriateness(items, request.Season)

	// Score occasion appropriateness
	occasionScore := ai.scoreOccasionAppropriateness(items, request.Occasion)

	// Calculate weighted overall score
	overallScore := (colorScore*0.25 + styleScore*0.25 + seasonScore*0.25 + occasionScore*0.25)

	// Generate reasoning
	reasoning := ai.generateOutfitReasoning(colorScore, styleScore, seasonScore, occasionScore)

	return OutfitScore{
		Outfit:              outfit,
		ColorHarmonyScore:   colorScore,
		StyleCoherenceScore: styleScore,
		SeasonScore:         seasonScore,
		OccasionScore:       occasionScore,
		OverallScore:        overallScore,
		Reasoning:           reasoning,
	}
}

func (ai *AIService) scoreSeasonAppropriateness(items []models.ClothingItem, season string) float64 {
	if season == "" {
		return 1.0
	}

	totalScore := 0.0
	for _, item := range items {
		itemScore := 0.5 // Base score

		for _, itemSeason := range item.Season {
			if itemSeason == season {
				itemScore = 1.0
				break
			}
		}

		totalScore += itemScore
	}

	return totalScore / float64(len(items))
}

func (ai *AIService) scoreOccasionAppropriateness(items []models.ClothingItem, occasion string) float64 {
	if occasion == "" {
		return 1.0
	}

	totalScore := 0.0
	for _, item := range items {
		itemScore := 0.5 // Base score

		for _, itemOccasion := range item.Occasion {
			if itemOccasion == occasion {
				itemScore = 1.0
				break
			}
		}

		totalScore += itemScore
	}

	return totalScore / float64(len(items))
}

func (ai *AIService) generateOutfitReasoning(colorScore, styleScore, seasonScore, occasionScore float64) string {
	reasons := []string{}

	if colorScore > 0.8 {
		reasons = append(reasons, "excellent color harmony")
	} else if colorScore > 0.6 {
		reasons = append(reasons, "good color coordination")
	} else {
		reasons = append(reasons, "acceptable color combination")
	}

	if styleScore > 0.8 {
		reasons = append(reasons, "highly coherent style")
	} else if styleScore > 0.6 {
		reasons = append(reasons, "well-matched pieces")
	} else {
		reasons = append(reasons, "mixed style elements")
	}

	if seasonScore > 0.8 {
		reasons = append(reasons, "perfect for the season")
	} else if seasonScore > 0.6 {
		reasons = append(reasons, "season-appropriate")
	}

	if occasionScore > 0.8 {
		reasons = append(reasons, "ideal for the occasion")
	} else if occasionScore > 0.6 {
		reasons = append(reasons, "suitable for the event")
	}

	if len(reasons) == 0 {
		return "A basic outfit combination"
	}

	return fmt.Sprintf("This outfit features %s", joinReasons(reasons))
}

func joinReasons(reasons []string) string {
	if len(reasons) == 1 {
		return reasons[0]
	}
	if len(reasons) == 2 {
		return reasons[0] + " and " + reasons[1]
	}

	result := ""
	for i, reason := range reasons {
		if i == len(reasons)-1 {
			result += "and " + reason
		} else if i == 0 {
			result += reason
		} else {
			result += ", " + reason
		}
	}
	return result
}

func (ai *AIService) calculateOverallOutfitScore(outfit *models.Outfit) float64 {
	return (outfit.ColorHarmonyScore + outfit.StyleCoherenceScore) / 2.0
}

func (ai *AIService) generateOutfitName(items []models.ClothingItem) string {
	if len(items) == 0 {
		return "Empty Outfit"
	}

	// Simple naming based on main pieces
	mainPieces := []string{}
	for _, item := range items {
		if item.Category == "dress" || item.Category == "suit" {
			return fmt.Sprintf("%s %s Outfit", item.Brand, item.Category)
		}
		if item.Category == "shirt" || item.Category == "pants" || item.Category == "skirt" {
			mainPieces = append(mainPieces, item.Category)
		}
	}

	if len(mainPieces) > 0 {
		return fmt.Sprintf("%s Combination", mainPieces[0])
	}

	return "Casual Outfit"
}

func (ai *AIService) generateOutfitDescription(items []models.ClothingItem) string {
	if len(items) == 0 {
		return "No items in this outfit"
	}

	descriptions := []string{}
	for _, item := range items {
		desc := item.Name
		if item.Brand != "" {
			desc = item.Brand + " " + desc
		}
		descriptions = append(descriptions, desc)
	}

	return fmt.Sprintf("Outfit featuring: %s", joinReasons(descriptions))
}

func (ai *AIService) determineOutfitStyle(items []models.ClothingItem) string {
	styleCounts := make(map[string]int)

	for _, item := range items {
		if item.Style != "" {
			styleCounts[item.Style]++
		}
	}

	maxCount := 0
	dominantStyle := "casual"
	for style, count := range styleCounts {
		if count > maxCount {
			maxCount = count
			dominantStyle = style
		}
	}

	return dominantStyle
}

func (ai *AIService) determineOutfitSeason(items []models.ClothingItem) []string {
	seasonCounts := make(map[string]int)

	for _, item := range items {
		for _, season := range item.Season {
			seasonCounts[season]++
		}
	}

	// Return seasons that appear in at least half the items
	threshold := len(items) / 2
	seasons := []string{}
	for season, count := range seasonCounts {
		if count >= threshold {
			seasons = append(seasons, season)
		}
	}

	if len(seasons) == 0 {
		seasons = []string{"all-season"}
	}

	return seasons
}

func (ai *AIService) determineOutfitOccasion(items []models.ClothingItem) []string {
	occasionCounts := make(map[string]int)

	for _, item := range items {
		for _, occasion := range item.Occasion {
			occasionCounts[occasion]++
		}
	}

	// Return occasions that appear in at least half the items
	threshold := len(items) / 2
	occasions := []string{}
	for occasion, count := range occasionCounts {
		if count >= threshold {
			occasions = append(occasions, occasion)
		}
	}

	if len(occasions) == 0 {
		occasions = []string{"casual"}
	}

	return occasions
}

func (ai *AIService) analyzeWardrobeGaps(items []models.ClothingItem, profile *models.StyleProfile) []string {
	recommendations := []string{}

	// Count items by category
	categoryCounts := make(map[string]int)
	for _, item := range items {
		categoryCounts[item.Category]++
	}

	// Essential categories for a complete wardrobe
	essentials := map[string]int{
		"shirt":   3,
		"pants":   3,
		"dress":   2,
		"shoes":   3,
		"jacket":  2,
		"sweater": 2,
	}

	for category, recommended := range essentials {
		current := categoryCounts[category]
		if current < recommended {
			gap := recommended - current
			recommendations = append(recommendations, fmt.Sprintf(
				"Consider adding %d more %s(s) to complete your wardrobe basics",
				gap, category,
			))
		}
	}

	return recommendations
}

func (ai *AIService) generateVersatilityRecommendations(items []models.ClothingItem) []string {
	recommendations := []string{}

	// Analyze versatility of current items
	lowVersatilityCount := 0
	for _, item := range items {
		// Simple versatility check based on color and pattern
		if item.PrimaryColor != "#FFFFFF" && item.PrimaryColor != "#000000" &&
			item.PrimaryColor != "#808080" && item.Pattern != "solid" {
			lowVersatilityCount++
		}
	}

	if lowVersatilityCount > len(items)/2 {
		recommendations = append(recommendations,
			"Consider adding more neutral-colored, solid pieces for better mix-and-match versatility")
	}

	// Check for basic neutral pieces
	hasBasicWhite := false
	hasBasicBlack := false
	for _, item := range items {
		if item.PrimaryColor == "#FFFFFF" && item.Pattern == "solid" {
			hasBasicWhite = true
		}
		if item.PrimaryColor == "#000000" && item.Pattern == "solid" {
			hasBasicBlack = true
		}
	}

	if !hasBasicWhite {
		recommendations = append(recommendations, "A basic white shirt or top would greatly increase your outfit options")
	}
	if !hasBasicBlack {
		recommendations = append(recommendations, "A basic black piece would serve as a versatile foundation for many outfits")
	}

	return recommendations
}

func (ai *AIService) generateSingleOutfit(items []models.ClothingItem, preferences *models.StyleProfile) (*OutfitSuggestion, error) {
	rand.Seed(time.Now().UnixNano())

	requiredCategories := []string{"tops", "bottoms"}
	optionalCategories := []string{"outerwear", "shoes", "accessories"}

	outfitItems := []models.ClothingItem{}

	for _, category := range requiredCategories {
		categoryItems := ai.filterItemsByCategory(items, category)
		if len(categoryItems) == 0 {
			return nil, fmt.Errorf("no items found for required category: %s", category)
		}

		selectedItem := categoryItems[rand.Intn(len(categoryItems))]
		outfitItems = append(outfitItems, selectedItem)
	}

	for _, category := range optionalCategories {
		if rand.Float64() < 0.7 {
			categoryItems := ai.filterItemsByCategory(items, category)
			if len(categoryItems) > 0 {
				selectedItem := categoryItems[rand.Intn(len(categoryItems))]
				outfitItems = append(outfitItems, selectedItem)
			}
		}
	}

	compatibility, reasoning, err := ai.AnalyzeOutfitCompatibility(outfitItems)
	if err != nil {
		return nil, err
	}

	styleScore := ai.styleService.AnalyzeOutfitStyleCoherence(outfitItems)
	colorScore := ai.analyzeColorHarmony(outfitItems)

	outfit := &OutfitSuggestion{
		Items:      outfitItems,
		Confidence: compatibility,
		Reasoning:  reasoning,
		StyleScore: styleScore,
		ColorScore: colorScore,
		Occasion:   ai.determineOccasion(outfitItems),
		Season:     ai.determineSeason(outfitItems),
	}

	return outfit, nil
}

func (ai *AIService) filterItemsByCategory(items []models.ClothingItem, category string) []models.ClothingItem {
	filtered := []models.ClothingItem{}
	for _, item := range items {
		if strings.ToLower(item.Category) == strings.ToLower(category) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func (ai *AIService) analyzeColorHarmony(items []models.ClothingItem) float64 {
	if len(items) <= 1 {
		return 1.0
	}

	colors := []string{}
	for _, item := range items {
		if item.PrimaryColor != "" {
			colors = append(colors, item.PrimaryColor)
		}
	}

	if len(colors) <= 1 {
		return 1.0
	}

	harmonyScore := 0.0
	comparisons := 0

	for i := 0; i < len(colors); i++ {
		for j := i + 1; j < len(colors); j++ {
			color1, err1 := colorful.Hex(colors[i])
			color2, err2 := colorful.Hex(colors[j])
			if err1 == nil && err2 == nil {
				compatibility := ai.colorService.calculateColorHarmonyScore(color1, color2)
				harmonyScore += compatibility
				comparisons++
			}
		}
	}

	if comparisons > 0 {
		harmonyScore /= float64(comparisons)
	}

	return harmonyScore
}

func (ai *AIService) generateCompatibilityReasoning(styleScore, colorScore float64, items []models.ClothingItem) string {
	reasons := []string{}

	if styleScore >= 0.8 {
		reasons = append(reasons, "Excellent style coherence")
	} else if styleScore >= 0.6 {
		reasons = append(reasons, "Good style compatibility")
	} else {
		reasons = append(reasons, "Mixed style elements")
	}

	if colorScore >= 0.8 {
		reasons = append(reasons, "Harmonious color palette")
	} else if colorScore >= 0.6 {
		reasons = append(reasons, "Complementary colors")
	} else {
		reasons = append(reasons, "Contrasting color scheme")
	}

	styles := make(map[string]int)
	for _, item := range items {
		if item.Style != "" {
			styles[item.Style]++
		}
	}

	if len(styles) == 1 {
		for style := range styles {
			reasons = append(reasons, fmt.Sprintf("Consistent %s style throughout", style))
		}
	} else if len(styles) > 1 {
		styleNames := make([]string, 0, len(styles))
		for style := range styles {
			styleNames = append(styleNames, style)
		}
		reasons = append(reasons, fmt.Sprintf("Combines %s styles", strings.Join(styleNames, " and ")))
	}

	return strings.Join(reasons, ". ")
}

func (ai *AIService) determineOccasion(items []models.ClothingItem) string {
	formalityScore := 0
	casualScore := 0

	for _, item := range items {
		switch strings.ToLower(item.Category) {
		case "suits", "dress_shirts", "dress_shoes":
			formalityScore += 2
		case "blazers", "dress_pants":
			formalityScore += 1
		case "t-shirts", "jeans", "sneakers":
			casualScore += 2
		case "casual_shirts", "chinos":
			casualScore += 1
		}

		if strings.Contains(strings.ToLower(item.Style), "formal") {
			formalityScore += 1
		} else if strings.Contains(strings.ToLower(item.Style), "casual") {
			casualScore += 1
		}
	}

	if formalityScore > int(float64(casualScore)*1.5) {
		return "formal"
	} else if casualScore > int(float64(formalityScore)*1.5) {
		return "casual"
	} else {
		return "smart-casual"
	}
}

func (ai *AIService) determineSeason(items []models.ClothingItem) string {
	seasonScores := map[string]int{
		"spring": 0,
		"summer": 0,
		"fall":   0,
		"winter": 0,
	}

	for _, item := range items {
		switch strings.ToLower(item.Category) {
		case "coats", "sweaters", "boots":
			seasonScores["winter"] += 2
			seasonScores["fall"] += 1
		case "t-shirts", "shorts", "sandals":
			seasonScores["summer"] += 2
			seasonScores["spring"] += 1
		case "light_jackets", "cardigans":
			seasonScores["spring"] += 2
			seasonScores["fall"] += 2
		}
	}

	maxScore := 0
	bestSeason := "all-season"
	for season, score := range seasonScores {
		if score > maxScore {
			maxScore = score
			bestSeason = season
		}
	}

	return bestSeason
}

func (ai *AIService) isOutfitUnique(newOutfit *OutfitSuggestion, existingRecommendations []models.OutfitRecommendation) bool {
	newItemIDs := make(map[uuid.UUID]bool)
	for _, item := range newOutfit.Items {
		newItemIDs[item.ID] = true
	}

	for _, existing := range existingRecommendations {
		if existing.OutfitID == newOutfit.Items[0].ID {
			return false
		}
	}

	return true
}

func (ai *AIService) GetStyleInsights(items []models.ClothingItem) (map[string]interface{}, error) {
	insights := make(map[string]interface{})

	if len(items) == 0 {
		return insights, nil
	}

	styleDistribution := make(map[string]int)
	colorDistribution := make(map[string]int)
	categoryDistribution := make(map[string]int)

	for _, item := range items {
		if item.Style != "" {
			styleDistribution[item.Style]++
		}
		if item.PrimaryColor != "" {
			colorDistribution[item.PrimaryColor]++
		}
		if item.Category != "" {
			categoryDistribution[item.Category]++
		}
	}

	insights["style_distribution"] = styleDistribution
	insights["color_distribution"] = colorDistribution
	insights["category_distribution"] = categoryDistribution

	dominantStyle := ""
	maxStyleCount := 0
	for style, count := range styleDistribution {
		if count > maxStyleCount {
			maxStyleCount = count
			dominantStyle = style
		}
	}
	insights["dominant_style"] = dominantStyle

	dominantColor := ""
	maxColorCount := 0
	for color, count := range colorDistribution {
		if count > maxColorCount {
			maxColorCount = count
			dominantColor = color
		}
	}
	insights["dominant_color"] = dominantColor

	recommendations, err := ai.styleService.GenerateStyleRecommendations(&models.StyleProfile{
		PreferredStyles: []string{dominantStyle},
	})
	if err == nil {
		insights["recommendations"] = recommendations
	}

	return insights, nil
}

func (ai *AIService) AnalyzeWardrobeGaps(items []models.ClothingItem, preferences *models.StyleProfile) ([]string, error) {
	gaps := []string{}

	essentialCategories := map[string]int{
		"tops":        3,
		"bottoms":     3,
		"outerwear":   2,
		"shoes":       2,
		"accessories": 1,
	}

	categoryCount := make(map[string]int)
	for _, item := range items {
		categoryCount[strings.ToLower(item.Category)]++
	}

	for category, minCount := range essentialCategories {
		if categoryCount[category] < minCount {
			needed := minCount - categoryCount[category]
			gaps = append(gaps, fmt.Sprintf("Need %d more %s items", needed, category))
		}
	}

	if preferences != nil && len(preferences.PreferredStyles) > 0 {
		styleCount := make(map[string]int)
		for _, item := range items {
			if item.Style != "" {
				styleCount[strings.ToLower(item.Style)]++
			}
		}

		for _, preferredStyle := range preferences.PreferredStyles {
			if styleCount[strings.ToLower(preferredStyle)] == 0 {
				gaps = append(gaps, fmt.Sprintf("No items in preferred %s style", preferredStyle))
			}
		}
	}

	return gaps, nil
}
