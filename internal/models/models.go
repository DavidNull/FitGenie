package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Name      string    `json:"name" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Style preferences
	StyleProfile StyleProfile `json:"style_profile" gorm:"foreignKey:UserID"`
	ColorProfile ColorProfile `json:"color_profile" gorm:"foreignKey:UserID"`

	// User's wardrobe
	ClothingItems []ClothingItem `json:"clothing_items" gorm:"foreignKey:UserID"`
	Outfits       []Outfit       `json:"outfits" gorm:"foreignKey:UserID"`
}

// StyleProfile represents user's style preferences
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

// ColorProfile represents user's color preferences and analysis
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

// ClothingItem represents a piece of clothing
type ClothingItem struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID         uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Name           string    `json:"name" gorm:"not null"`
	Category       string    `json:"category" gorm:"not null"`
	Brand          string    `json:"brand"`
	Size           string    `json:"size"`
	PrimaryColor   string    `json:"primary_color"`
	SecondaryColor string    `json:"secondary_color"`
	Material       string    `json:"material"`
	Style          string    `json:"style"`
	Season         []string  `json:"season" gorm:"type:text[]"`
	Occasion       []string  `json:"occasion" gorm:"type:text[]"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Outfit represents a complete outfit combination
type Outfit struct {
	ID                  uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID              uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Name                string         `json:"name" gorm:"not null"`
	Description         string         `json:"description"`
	Style               string         `json:"style"`
	Occasion            []string       `json:"occasion" gorm:"type:text[]"`
	Season              []string       `json:"season" gorm:"type:text[]"`
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

// OutfitRecommendation represents an AI-generated outfit recommendation
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

// FavoriteOutfit represents a user's favorite outfit (many-to-many relationship)
type FavoriteOutfit struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	OutfitID  uuid.UUID `json:"outfit_id" gorm:"type:uuid;not null"`
	CreatedAt time.Time `json:"created_at"`
}

// BeforeCreate hook for User model
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// BeforeCreate hooks for other models
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
