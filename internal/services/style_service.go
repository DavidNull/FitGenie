package services

import (
	"fitgenie/internal/models"
	"fmt"
	"strings"
)

type StyleService struct {
	styleCategories map[string]StyleCategory
	bodyTypeGuides  map[string]BodyTypeGuide
	occasionStyles  map[string][]string
}

type StyleCategory struct {
	Name        string
	Description string
	Keywords    []string
	Colors      []string
	Patterns    []string
	Materials   []string
	Silhouettes []string
	Accessories []string
	Occasions   []string
}

type BodyTypeGuide struct {
	Name            string
	Description     string
	RecommendedCuts []string
	AvoidCuts       []string
	BestSilhouettes []string
	StylingTips     []string
}

func NewStyleService() *StyleService {
	s := &StyleService{
		styleCategories: make(map[string]StyleCategory),
		bodyTypeGuides:  make(map[string]BodyTypeGuide),
		occasionStyles:  make(map[string][]string),
	}
	s.initializeStyleCategories()
	s.initializeBodyTypeGuides()
	s.initializeOccasionStyles()
	return s
}

func (s *StyleService) initializeStyleCategories() {
	s.styleCategories = map[string]StyleCategory{
		"casual": {
			Name:        "Casual",
			Description: "Relaxed, comfortable, everyday wear",
			Keywords:    []string{"comfortable", "relaxed", "everyday", "laid-back"},
			Colors:      []string{"#4A90E2", "#7ED321", "#F5A623", "#D0021B", "#9013FE"},
			Patterns:    []string{"solid", "stripes", "simple prints"},
			Materials:   []string{"cotton", "denim", "jersey", "fleece"},
			Silhouettes: []string{"relaxed fit", "straight cut", "loose"},
			Accessories: []string{"sneakers", "baseball cap", "backpack", "casual watch"},
			Occasions:   []string{"weekend", "shopping", "casual meetups", "home"},
		},
		"business": {
			Name:        "Business",
			Description: "Professional, polished, workplace appropriate",
			Keywords:    []string{"professional", "polished", "workplace", "formal"},
			Colors:      []string{"#000000", "#2C3E50", "#34495E", "#7F8C8D", "#FFFFFF"},
			Patterns:    []string{"solid", "pinstripes", "subtle checks"},
			Materials:   []string{"wool", "cotton blend", "silk", "polyester"},
			Silhouettes: []string{"tailored", "structured", "fitted"},
			Accessories: []string{"leather shoes", "briefcase", "watch", "belt"},
			Occasions:   []string{"office", "meetings", "presentations", "networking"},
		},
		"formal": {
			Name:        "Formal",
			Description: "Elegant, sophisticated, special occasions",
			Keywords:    []string{"elegant", "sophisticated", "dressy", "special occasion"},
			Colors:      []string{"#000000", "#FFFFFF", "#2C3E50", "#8E44AD", "#C0392B"},
			Patterns:    []string{"solid", "subtle textures", "minimal patterns"},
			Materials:   []string{"silk", "satin", "wool", "velvet", "chiffon"},
			Silhouettes: []string{"fitted", "A-line", "tailored", "flowing"},
			Accessories: []string{"dress shoes", "jewelry", "clutch", "formal watch"},
			Occasions:   []string{"weddings", "galas", "formal dinners", "theater"},
		},
		"bohemian": {
			Name:        "Bohemian",
			Description: "Free-spirited, artistic, eclectic mix",
			Keywords:    []string{"free-spirited", "artistic", "eclectic", "boho"},
			Colors:      []string{"#D35400", "#E67E22", "#F39C12", "#27AE60", "#8E44AD"},
			Patterns:    []string{"paisley", "floral", "ethnic prints", "tie-dye"},
			Materials:   []string{"cotton", "linen", "suede", "fringe", "crochet"},
			Silhouettes: []string{"flowing", "loose", "layered", "asymmetrical"},
			Accessories: []string{"sandals", "layered jewelry", "headbands", "fringe bags"},
			Occasions:   []string{"festivals", "art events", "casual outings", "travel"},
		},
		"minimalist": {
			Name:        "Minimalist",
			Description: "Clean, simple, understated elegance",
			Keywords:    []string{"clean", "simple", "understated", "minimal"},
			Colors:      []string{"#FFFFFF", "#000000", "#95A5A6", "#BDC3C7", "#ECF0F1"},
			Patterns:    []string{"solid", "geometric", "minimal"},
			Materials:   []string{"cotton", "linen", "wool", "cashmere"},
			Silhouettes: []string{"clean lines", "structured", "simple cuts"},
			Accessories: []string{"simple jewelry", "leather goods", "minimal bags"},
			Occasions:   []string{"work", "casual", "modern events", "gallery openings"},
		},
		"sporty": {
			Name:        "Sporty",
			Description: "Athletic, active, performance-oriented",
			Keywords:    []string{"athletic", "active", "performance", "sporty"},
			Colors:      []string{"#E74C3C", "#3498DB", "#2ECC71", "#F39C12", "#9B59B6"},
			Patterns:    []string{"solid", "stripes", "color blocking"},
			Materials:   []string{"polyester", "spandex", "moisture-wicking", "mesh"},
			Silhouettes: []string{"fitted", "flexible", "streamlined"},
			Accessories: []string{"sneakers", "sports watch", "gym bag", "cap"},
			Occasions:   []string{"gym", "sports", "outdoor activities", "casual wear"},
		},
		"romantic": {
			Name:        "Romantic",
			Description: "Feminine, soft, delicate details",
			Keywords:    []string{"feminine", "soft", "delicate", "romantic"},
			Colors:      []string{"#F8BBD9", "#E8F5E8", "#FFF2CC", "#E1BEE7", "#FFE5CC"},
			Patterns:    []string{"floral", "lace", "polka dots", "soft prints"},
			Materials:   []string{"chiffon", "lace", "silk", "satin", "tulle"},
			Silhouettes: []string{"flowing", "fitted waist", "A-line", "ruffles"},
			Accessories: []string{"delicate jewelry", "ballet flats", "small bags", "scarves"},
			Occasions:   []string{"dates", "brunches", "garden parties", "romantic dinners"},
		},
		"edgy": {
			Name:        "Edgy",
			Description: "Bold, unconventional, and statement-making pieces",
			Keywords:    []string{"leather", "studs", "asymmetric", "bold", "unconventional"},
			Colors:      []string{"#000000", "#8B0000", "#4B0082", "#2F4F4F"},
			Patterns:    []string{"geometric", "abstract", "animal print"},
			Materials:   []string{"leather", "denim", "metal", "vinyl"},
			Silhouettes: []string{"fitted", "asymmetrical", "structured"},
			Accessories: []string{"boots", "statement jewelry", "leather bags"},
			Occasions:   []string{"party", "casual", "date"},
		},
	}
}

// initializeBodyTypeGuides sets up style recommendations for different body types
func (s *StyleService) initializeBodyTypeGuides() {
	s.bodyTypeGuides = map[string]BodyTypeGuide{
		"pear": {
			Name:        "Pear",
			Description: "Narrower shoulders, wider hips",
			RecommendedCuts: []string{
				"A-line tops", "boat neck", "off-shoulder", "wide leg pants",
				"bootcut jeans", "empire waist", "wrap tops",
			},
			AvoidCuts: []string{
				"skinny jeans", "tight bottoms", "hip-hugging styles",
				"low-rise pants", "pencil skirts",
			},
			BestSilhouettes: []string{
				"balance upper body", "emphasize shoulders", "create waist definition",
			},
			StylingTips: []string{
				"Add volume to upper body", "Use bright colors on top",
				"Choose darker colors for bottom", "Emphasize your waist",
			},
		},
		"apple": {
			Name:        "Apple",
			Description: "Fuller midsection, narrower hips",
			RecommendedCuts: []string{
				"empire waist", "A-line dresses", "V-neck tops", "wrap dresses",
				"straight leg pants", "high-waisted bottoms",
			},
			AvoidCuts: []string{
				"tight fitting tops", "horizontal stripes across middle",
				"low-rise pants", "cropped tops",
			},
			BestSilhouettes: []string{
				"elongate torso", "define waist above natural line",
				"draw attention to legs and arms",
			},
			StylingTips: []string{
				"Create vertical lines", "Show off your legs",
				"Use accessories to draw attention up", "Choose flowing fabrics",
			},
		},
		"hourglass": {
			Name:        "Hourglass",
			Description: "Balanced shoulders and hips, defined waist",
			RecommendedCuts: []string{
				"fitted tops", "wrap dresses", "belted styles", "pencil skirts",
				"high-waisted pants", "bodycon dresses",
			},
			AvoidCuts: []string{
				"boxy tops", "loose fitting clothes", "drop waist",
				"baggy styles that hide your shape",
			},
			BestSilhouettes: []string{
				"emphasize waist", "show your curves", "fitted styles",
			},
			StylingTips: []string{
				"Highlight your waist with belts", "Choose fitted styles",
				"Avoid hiding your natural shape", "Balance is key",
			},
		},
		"rectangle": {
			Name:        "Rectangle",
			Description: "Similar measurements for shoulders, waist, and hips",
			RecommendedCuts: []string{
				"peplum tops", "ruffled details", "layered looks", "cropped jackets",
				"wide leg pants", "A-line skirts",
			},
			AvoidCuts: []string{
				"straight, boxy cuts", "shapeless dresses",
				"baggy clothes without structure",
			},
			BestSilhouettes: []string{
				"create curves", "add volume strategically", "define waist",
			},
			StylingTips: []string{
				"Create the illusion of curves", "Use belts to define waist",
				"Add texture and layers", "Choose interesting necklines",
			},
		},
		"inverted_triangle": {
			Name:        "Inverted Triangle",
			Description: "Broader shoulders, narrower hips",
			RecommendedCuts: []string{
				"straight leg pants", "wide leg pants", "A-line skirts",
				"scoop neck tops", "V-neck styles", "hip details",
			},
			AvoidCuts: []string{
				"shoulder pads", "boat neck", "horizontal stripes on top",
				"skinny jeans", "narrow bottoms",
			},
			BestSilhouettes: []string{
				"balance proportions", "add volume to lower body",
				"minimize shoulder width",
			},
			StylingTips: []string{
				"Add volume to hips and legs", "Choose darker colors on top",
				"Use bright colors and patterns on bottom", "Avoid shoulder emphasis",
			},
		},
	}
}

// initializeOccasionStyles maps occasions to appropriate styles
func (s *StyleService) initializeOccasionStyles() {
	s.occasionStyles = map[string][]string{
		"work":        {"business", "minimalist", "casual"},
		"date":        {"romantic", "casual", "minimalist"},
		"party":       {"formal", "edgy", "romantic"},
		"wedding":     {"formal", "romantic", "minimalist"},
		"casual":      {"casual", "bohemian", "sporty"},
		"gym":         {"sporty"},
		"beach":       {"casual", "bohemian"},
		"travel":      {"casual", "minimalist", "bohemian"},
		"interview":   {"business", "minimalist"},
		"dinner":      {"formal", "business", "romantic"},
		"shopping":    {"casual", "minimalist"},
		"concert":     {"edgy", "casual", "bohemian"},
		"art_gallery": {"minimalist", "bohemian", "edgy"},
		"brunch":      {"casual", "romantic", "minimalist"},
		"networking":  {"business", "minimalist"},
	}
}

// AnalyzeClothingStyle analyzes the style of a clothing item
func (s *StyleService) AnalyzeClothingStyle(item *models.ClothingItem) (*models.StyleAnalysis, error) {
	analysis := &models.StyleAnalysis{
		ClothingItemID:      item.ID,
		SeasonSuitability:   make(map[string]float64),
		OccasionSuitability: make(map[string]float64),
	}

	// Determine style category
	styleScore := make(map[string]float64)

	for styleName, style := range s.styleCategories {
		score := 0.0

		// Check keywords in item description/name
		itemText := strings.ToLower(item.Name + " " + item.Category)
		for _, keyword := range style.Keywords {
			if strings.Contains(itemText, keyword) {
				score += 2.0
			}
		}

		// Check colors
		for _, styleColor := range style.Colors {
			if item.PrimaryColor == styleColor || item.SecondaryColor == styleColor {
				score += 1.5
			}
		}

		// Check materials
		for _, material := range style.Materials {
			if strings.Contains(strings.ToLower(item.Material), material) {
				score += 1.0
			}
		}

		styleScore[styleName] = score
	}

	// Find the best matching style
	maxScore := 0.0
	bestStyle := "casual"
	for style, score := range styleScore {
		if score > maxScore {
			maxScore = score
			bestStyle = style
		}
	}

	analysis.StyleCategory = bestStyle

	// Determine formality level
	analysis.Formality = s.determineFormalityLevel(item, bestStyle)

	// Calculate versatility score
	analysis.Versatility = s.calculateVersatility(item, bestStyle)

	// Analyze season suitability
	analysis.SeasonSuitability = s.analyzeSeasonSuitability(item)

	// Analyze occasion suitability
	analysis.OccasionSuitability = s.analyzeOccasionSuitability(item, bestStyle)

	return analysis, nil
}

// AnalyzeOutfitStyleCoherence analyzes how well the styles in an outfit work together
func (s *StyleService) AnalyzeOutfitStyleCoherence(items []models.ClothingItem) float64 {
	if len(items) == 0 {
		return 0.0
	}

	if len(items) == 1 {
		return 1.0
	}

	styles := make(map[string]int)
	for _, item := range items {
		if item.Style != "" {
			styles[strings.ToLower(item.Style)]++
		}
	}

	if len(styles) == 0 {
		return 0.5
	}

	if len(styles) == 1 {
		return 1.0
	}

	coherenceScore := s.calculateStyleCoherence(styles, len(items))
	return coherenceScore
}

// GenerateStyleRecommendations generates style recommendations based on user profile
func (s *StyleService) GenerateStyleRecommendations(profile *models.StyleProfile) ([]string, error) {
	recommendations := []string{}

	if len(profile.PreferredStyles) == 0 {
		recommendations = append(recommendations, "Consider defining your preferred style categories to get better recommendations")
		return recommendations, nil
	}

	for _, style := range profile.PreferredStyles {
		styleRecs := s.getStyleSpecificRecommendations(style, profile)
		recommendations = append(recommendations, styleRecs...)
	}

	bodyTypeRecs := s.getBodyTypeRecommendations(profile.BodyType)
	recommendations = append(recommendations, bodyTypeRecs...)

	lifestyleRecs := s.getLifestyleRecommendations(profile.Lifestyle)
	recommendations = append(recommendations, lifestyleRecs...)

	return recommendations, nil
}

// Helper methods

func (s *StyleService) determineFormalityLevel(item *models.ClothingItem, style string) string {
	formalityMap := map[string]string{
		"formal":     "very formal",
		"business":   "formal",
		"minimalist": "smart casual",
		"romantic":   "smart casual",
		"casual":     "casual",
		"bohemian":   "casual",
		"sporty":     "very casual",
		"edgy":       "casual",
	}

	if formality, exists := formalityMap[style]; exists {
		return formality
	}

	// Fallback based on category
	switch strings.ToLower(item.Category) {
	case "suit", "blazer", "dress shirt":
		return "formal"
	case "dress", "blouse":
		return "smart casual"
	case "t-shirt", "jeans", "sneakers":
		return "casual"
	case "activewear", "shorts":
		return "very casual"
	default:
		return "casual"
	}
}

func (s *StyleService) calculateVersatility(item *models.ClothingItem, style string) float64 {
	versatilityScore := 0.5 // Base score

	// Neutral colors are more versatile
	neutralColors := []string{"#FFFFFF", "#000000", "#808080", "#2F4F4F", "#F5F5F5"}
	for _, neutral := range neutralColors {
		if item.PrimaryColor == neutral {
			versatilityScore += 0.2
			break
		}
	}

	// Basic categories are more versatile
	basicCategories := []string{"shirt", "pants", "jeans", "blazer", "cardigan"}
	for _, basic := range basicCategories {
		if strings.Contains(strings.ToLower(item.Category), basic) {
			versatilityScore += 0.1
			break
		}
	}

	// Multiple seasons increase versatility
	if len(item.Season) > 2 {
		versatilityScore += 0.1
	}

	return versatilityScore
}

func (s *StyleService) analyzeSeasonSuitability(item *models.ClothingItem) map[string]float64 {
	suitability := map[string]float64{
		"spring": 0.5,
		"summer": 0.5,
		"autumn": 0.5,
		"winter": 0.5,
	}

	// Boost scores for explicitly mentioned seasons
	for _, season := range item.Season {
		if score, exists := suitability[strings.ToLower(season)]; exists {
			suitability[strings.ToLower(season)] = score + 0.4
		}
	}

	// Material-based adjustments
	material := strings.ToLower(item.Material)
	if strings.Contains(material, "wool") || strings.Contains(material, "cashmere") {
		suitability["winter"] += 0.3
		suitability["autumn"] += 0.2
	}
	if strings.Contains(material, "cotton") || strings.Contains(material, "linen") {
		suitability["summer"] += 0.3
		suitability["spring"] += 0.2
	}

	// Normalize scores to max 1.0
	for season := range suitability {
		if suitability[season] > 1.0 {
			suitability[season] = 1.0
		}
	}

	return suitability
}

func (s *StyleService) analyzeOccasionSuitability(item *models.ClothingItem, style string) map[string]float64 {
	suitability := make(map[string]float64)

	// Initialize all occasions with base score
	for occasion := range s.occasionStyles {
		suitability[occasion] = 0.1
	}

	// Boost scores based on style category
	if styleCategory, exists := s.styleCategories[style]; exists {
		for _, occasion := range styleCategory.Occasions {
			if _, exists := suitability[occasion]; exists {
				suitability[occasion] += 0.6
			}
		}
	}

	// Boost scores for explicitly mentioned occasions
	for _, occasion := range item.Occasion {
		occasionKey := strings.ToLower(strings.ReplaceAll(occasion, " ", "_"))
		if _, exists := suitability[occasionKey]; exists {
			suitability[occasionKey] += 0.3
		}
	}

	// Normalize scores
	for occasion := range suitability {
		if suitability[occasion] > 1.0 {
			suitability[occasion] = 1.0
		}
	}

	return suitability
}

func (s *StyleService) calculateStyleCoherence(styles map[string]int, totalItems int) float64 {
	if len(styles) == 1 {
		return 1.0
	}

	maxCount := 0
	for _, count := range styles {
		if count > maxCount {
			maxCount = count
		}
	}

	dominantRatio := float64(maxCount) / float64(totalItems)

	compatibilityScore := 0.0
	comparisons := 0

	styleList := make([]string, 0, len(styles))
	for style := range styles {
		styleList = append(styleList, style)
	}

	for i := 0; i < len(styleList); i++ {
		for j := i + 1; j < len(styleList); j++ {
			compatibility := s.AnalyzeStyleCompatibility(styleList[i], styleList[j])
			compatibilityScore += compatibility
			comparisons++
		}
	}

	if comparisons > 0 {
		compatibilityScore /= float64(comparisons)
	} else {
		compatibilityScore = 1.0
	}

	coherenceScore := (dominantRatio*0.6 + compatibilityScore*0.4)
	return coherenceScore
}

func (s *StyleService) AnalyzeStyleCompatibility(style1, style2 string) float64 {
	style1 = strings.ToLower(strings.TrimSpace(style1))
	style2 = strings.ToLower(strings.TrimSpace(style2))

	if style1 == style2 {
		return 1.0
	}

	compatibilityMatrix := map[string]map[string]float64{
		"casual": {
			"sporty":     0.8,
			"bohemian":   0.7,
			"minimalist": 0.6,
			"vintage":    0.5,
			"business":   0.3,
			"formal":     0.2,
		},
		"business": {
			"formal":     0.9,
			"minimalist": 0.8,
			"classic":    0.7,
			"casual":     0.3,
			"sporty":     0.2,
			"bohemian":   0.2,
		},
		"formal": {
			"business": 0.9,
			"classic":  0.8,
			"elegant":  0.7,
			"casual":   0.1,
			"sporty":   0.1,
		},
	}

	if matrix, exists := compatibilityMatrix[style1]; exists {
		if score, exists := matrix[style2]; exists {
			return score
		}
	}

	if matrix, exists := compatibilityMatrix[style2]; exists {
		if score, exists := matrix[style1]; exists {
			return score
		}
	}

	return 0.5
}

func (s *StyleService) getStyleSpecificRecommendations(style string, profile *models.StyleProfile) []string {
	recommendations := []string{}

	styleRecs := map[string][]string{
		"casual": {
			"Add comfortable jeans in dark wash for versatility",
			"Invest in quality basic t-shirts in neutral colors",
			"Consider comfortable sneakers for everyday wear",
		},
		"business": {
			"Build a collection of well-fitted blazers",
			"Invest in quality dress shirts in white and light blue",
			"Add tailored trousers in navy and charcoal",
		},
		"formal": {
			"Invest in a well-tailored suit in navy or charcoal",
			"Add formal dress shoes in black and brown",
			"Consider elegant accessories like watches and cufflinks",
		},
		"bohemian": {
			"Add flowing fabrics and earth tones to your wardrobe",
			"Consider layered jewelry and accessories",
			"Look for pieces with interesting textures and patterns",
		},
		"minimalist": {
			"Focus on high-quality basics in neutral colors",
			"Choose pieces with clean lines and simple silhouettes",
			"Invest in versatile items that can be mixed and matched",
		},
	}

	if recs, exists := styleRecs[strings.ToLower(style)]; exists {
		recommendations = append(recommendations, recs...)
	} else {
		recommendations = append(recommendations, fmt.Sprintf("Explore %s style pieces that reflect your personality", style))
	}

	return recommendations
}

func (s *StyleService) getBodyTypeRecommendations(bodyType string) []string {
	recommendations := []string{}

	bodyTypeRecs := map[string][]string{
		"pear": {
			"Emphasize your upper body with statement tops",
			"Choose A-line skirts and wide-leg pants",
			"Add structured blazers to balance your silhouette",
		},
		"apple": {
			"Choose empire waist dresses and tops",
			"Look for V-necks to elongate your torso",
			"Add straight-leg pants and A-line skirts",
		},
		"hourglass": {
			"Emphasize your waist with fitted clothing",
			"Choose wrap dresses and belted styles",
			"Look for pieces that follow your natural curves",
		},
		"rectangle": {
			"Create curves with peplum tops and fit-and-flare dresses",
			"Add belts to define your waist",
			"Choose pieces with interesting details and textures",
		},
		"inverted-triangle": {
			"Balance your shoulders with wider bottom pieces",
			"Choose A-line skirts and wide-leg pants",
			"Look for tops with minimal shoulder details",
		},
	}

	if recs, exists := bodyTypeRecs[strings.ToLower(bodyType)]; exists {
		recommendations = append(recommendations, recs...)
	}

	return recommendations
}

func (s *StyleService) getLifestyleRecommendations(lifestyle string) []string {
	recommendations := []string{}

	lifestyleRecs := map[string][]string{
		"active": {
			"Invest in quality activewear that transitions to casual wear",
			"Choose comfortable shoes suitable for walking",
			"Look for wrinkle-resistant fabrics for busy days",
		},
		"professional": {
			"Build a capsule wardrobe of mix-and-match pieces",
			"Invest in quality basics that can be dressed up or down",
			"Choose pieces that are appropriate for your work environment",
		},
		"social": {
			"Add statement pieces for special occasions",
			"Invest in versatile dresses suitable for various events",
			"Choose accessories that can transform basic outfits",
		},
		"casual": {
			"Focus on comfort without sacrificing style",
			"Choose easy-care fabrics and relaxed fits",
			"Invest in quality basics that work for everyday activities",
		},
	}

	if recs, exists := lifestyleRecs[strings.ToLower(lifestyle)]; exists {
		recommendations = append(recommendations, recs...)
	}

	return recommendations
}
