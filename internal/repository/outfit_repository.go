package repository

import (
	"context"
	"fmt"

	"fitgenie/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OutfitRepository interface {
	Create(ctx context.Context, outfit *models.Outfit) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Outfit, error)
	ListByUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]models.Outfit, int64, error)
	Update(ctx context.Context, outfit *models.Outfit) error
	Delete(ctx context.Context, id uuid.UUID) error
	AddClothingItem(ctx context.Context, outfitID, clothingItemID uuid.UUID) error
	RemoveClothingItem(ctx context.Context, outfitID, clothingItemID uuid.UUID) error
}

type outfitRepository struct {
	db *gorm.DB
}

func NewOutfitRepository(db *gorm.DB) OutfitRepository {
	return &outfitRepository{db: db}
}

func (r *outfitRepository) Create(ctx context.Context, outfit *models.Outfit) error {
	if err := r.db.WithContext(ctx).Create(outfit).Error; err != nil {
		return fmt.Errorf("failed to create outfit: %w", err)
	}
	return nil
}

func (r *outfitRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Outfit, error) {
	var outfit models.Outfit
	if err := r.db.WithContext(ctx).
		Preload("ClothingItems").
		First(&outfit, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("failed to get outfit: %w", err)
	}
	return &outfit, nil
}

func (r *outfitRepository) ListByUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]models.Outfit, int64, error) {
	var outfits []models.Outfit
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.Outfit{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count outfits: %w", err)
	}

	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("ClothingItems").
		Offset(offset).
		Limit(limit).
		Find(&outfits).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list outfits: %w", err)
	}

	return outfits, total, nil
}

func (r *outfitRepository) Update(ctx context.Context, outfit *models.Outfit) error {
	if err := r.db.WithContext(ctx).Save(outfit).Error; err != nil {
		return fmt.Errorf("failed to update outfit: %w", err)
	}
	return nil
}

func (r *outfitRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.Outfit{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete outfit: %w", err)
	}
	return nil
}

func (r *outfitRepository) AddClothingItem(ctx context.Context, outfitID, clothingItemID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Exec("INSERT INTO outfit_clothing_items (outfit_id, clothing_item_id) VALUES (?, ?)",
			outfitID, clothingItemID).Error
}

func (r *outfitRepository) RemoveClothingItem(ctx context.Context, outfitID, clothingItemID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Exec("DELETE FROM outfit_clothing_items WHERE outfit_id = ? AND clothing_item_id = ?",
			outfitID, clothingItemID).Error
}
