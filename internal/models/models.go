package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email         string         `json:"email" gorm:"uniqueIndex;not null"`
	Name          string         `json:"name" gorm:"not null"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	StyleProfile  StyleProfile   `json:"style_profile" gorm:"foreignKey:UserID"`
	ColorProfile  ColorProfile   `json:"color_profile" gorm:"foreignKey:UserID"`
	ClothingItems []ClothingItem `json:"clothing_items" gorm:"foreignKey:UserID"`
	Outfits       []Outfit       `json:"outfits" gorm:"foreignKey:UserID"`
}

type StyleProfile struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID          uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	PreferredStyles []string  `json:"preferred_styles" gorm:"type:text[]"`
	BodyType        string    `json:"body_type"`
	Lifestyle       string    `json:"lifestyle"`
	Occasion        []string  `json:"occasion" gorm:"type:text[]"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ColorProfile struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID            uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	ColorSeason       string    `json:"color_season"`
	SkinTone          string    `json:"skin_tone"`
	Undertone         string    `json:"undertone"`
	FavoriteColors    []string  `json:"favorite_colors" gorm:"type:text[]"`
	AvoidColors       []string  `json:"avoid_colors" gorm:"type:text[]"`
	RecommendedColors []string  `json:"recommended_colors" gorm:"type:text[]"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type ClothingItem struct {
	ID             uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID         uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Name           string         `json:"name" gorm:"not null"`
	Category       string         `json:"category" gorm:"not null"`
	Brand          string         `json:"brand"`
	Size           string         `json:"size"`
	PrimaryColor   string         `json:"primary_color"`
	SecondaryColor string         `json:"secondary_color"`
	Material       string         `json:"material"`
	Style          string         `json:"style"`
	Season         pq.StringArray `json:"season" gorm:"type:text[]"`
	Occasion       pq.StringArray `json:"occasion" gorm:"type:text[]"`
	Notes          string         `json:"notes"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type Outfit struct {
	ID                  uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID              uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Name                string         `json:"name" gorm:"not null"`
	Description         string         `json:"description"`
	Style               string         `json:"style"`
	Occasion            pq.StringArray `json:"occasion" gorm:"type:text[]"`
	Season              pq.StringArray `json:"season" gorm:"type:text[]"`
	Weather             string         `json:"weather"`
	ColorHarmonyScore   *float64       `json:"color_harmony_score" gorm:"type:decimal"`
	StyleCoherenceScore *float64       `json:"style_coherence_score" gorm:"type:decimal"`
	OverallScore        *float64       `json:"overall_score" gorm:"type:decimal"`
	Rating              *int           `json:"rating"`
	Worn                bool           `json:"worn"`
	Favorite            bool           `json:"favorite"`
	Notes               string         `json:"notes"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	ClothingItems       []ClothingItem `json:"clothing_items" gorm:"many2many:outfit_clothing_items;"`
}

type OutfitRecommendation struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OutfitID          uuid.UUID `json:"outfit_id" gorm:"type:uuid;not null"`
	Outfit            Outfit    `json:"outfit" gorm:"foreignKey:OutfitID"`
	Confidence        float64   `json:"confidence"`
	Reasoning         string    `json:"reasoning"`
	RequestedOccasion string    `json:"requested_occasion"`
	RequestedSeason   string    `json:"requested_season"`
	RequestedWeather  string    `json:"requested_weather"`
	RequestedStyle    string    `json:"requested_style"`
	Viewed            bool      `json:"viewed"`
	Accepted          bool      `json:"accepted"`
	CreatedAt         time.Time `json:"created_at"`
}

type FavoriteOutfit struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	OutfitID  uuid.UUID `json:"outfit_id" gorm:"type:uuid;not null"`
	CreatedAt time.Time `json:"created_at"`
}

// HSL represents a color in HSL format
type HSL struct {
	H float64 `json:"h"` // Hue (0-360)
	S float64 `json:"s"` // Saturation (0-100)
	L float64 `json:"l"` // Lightness (0-100)
}

// RGB represents a color in RGB format
type RGB struct {
	R uint8 `json:"r"` // Red (0-255)
	G uint8 `json:"g"` // Green (0-255)
	B uint8 `json:"b"` // Blue (0-255)
}

// ColorInfo represents a color extracted from an image
type ColorInfo struct {
	Hex        string  `json:"hex"`
	Name       string  `json:"name"`
	Percentage float64 `json:"percentage"`
	R          uint8   `json:"r"`
	G          uint8   `json:"g"`
	B          uint8   `json:"b"`
	HSL        HSL     `json:"hsl"`
	RGB        RGB     `json:"rgb"`
}

// ColorAnalysis contains analysis of colors from an image
type ColorAnalysis struct {
	DominantColors   []ColorInfo `json:"dominant_colors"`
	ColorHarmony     string      `json:"color_harmony"`
	ColorTemperature string      `json:"color_temperature"`
	ColorSaturation  string      `json:"color_saturation"`
	ColorBrightness  string      `json:"color_brightness"`
}

// ColorRecommendation represents a color recommendation for a user
type ColorRecommendation struct {
	ColorHex   string  `json:"color_hex"`
	ColorName  string  `json:"color_name"`
	Category   string  `json:"category"` // Primary, Secondary, Neutral
	Reason     string  `json:"reason"`
	Confidence float64 `json:"confidence"`
}

// StyleAnalysis contains style analysis for a clothing item
type StyleAnalysis struct {
	ClothingItemID      uuid.UUID          `json:"clothing_item_id"`
	StyleCategory       string             `json:"style_category"`
	Formality           string             `json:"formality"`
	Versatility         float64            `json:"versatility"`
	SeasonSuitability   map[string]float64 `json:"season_suitability"`
	OccasionSuitability map[string]float64 `json:"occasion_suitability"`
	StyleMatchScore     float64            `json:"style_match_score"`
}

// BeforeCreate hooks for UUID generation
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (sp *StyleProfile) BeforeCreate(tx *gorm.DB) error {
	if sp.ID == uuid.Nil {
		sp.ID = uuid.New()
	}
	return nil
}

func (cp *ColorProfile) BeforeCreate(tx *gorm.DB) error {
	if cp.ID == uuid.Nil {
		cp.ID = uuid.New()
	}
	return nil
}

func (ci *ClothingItem) BeforeCreate(tx *gorm.DB) error {
	if ci.ID == uuid.Nil {
		ci.ID = uuid.New()
	}
	return nil
}

func (o *Outfit) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

func (or *OutfitRecommendation) BeforeCreate(tx *gorm.DB) error {
	if or.ID == uuid.Nil {
		or.ID = uuid.New()
	}
	return nil
}

func (fo *FavoriteOutfit) BeforeCreate(tx *gorm.DB) error {
	if fo.ID == uuid.Nil {
		fo.ID = uuid.New()
	}
	return nil
}
