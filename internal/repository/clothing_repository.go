package repository

import (
	"context"
	"fmt"
	"log"

	"fitgenie/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClothingRepository interface {
	Create(ctx context.Context, item *models.ClothingItem) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.ClothingItem, error)
	ListByUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]models.ClothingItem, int64, error)
	Update(ctx context.Context, item *models.ClothingItem) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type clothingRepository struct {
	db *gorm.DB
}

func NewClothingRepository(db *gorm.DB) ClothingRepository {
	return &clothingRepository{db: db}
}

func (r *clothingRepository) Create(ctx context.Context, item *models.ClothingItem) error {
	log.Printf("[DEBUG] Creating clothing item: user_id=%s, name=%s, category=%s", item.UserID, item.Name, item.Category)
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		log.Printf("[ERROR] Failed to create clothing item: %v", err)
		return fmt.Errorf("failed to create clothing item: %w", err)
	}
	log.Printf("[DEBUG] Successfully created clothing item: id=%s", item.ID)
	return nil
}

func (r *clothingRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.ClothingItem, error) {
	var item models.ClothingItem
	if err := r.db.WithContext(ctx).First(&item, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("failed to get clothing item: %w", err)
	}
	return &item, nil
}

func (r *clothingRepository) ListByUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]models.ClothingItem, int64, error) {
	var items []models.ClothingItem
	var total int64

	userIDStr := userID.String()

	// Use Raw SQL for debugging
	countQuery := "SELECT COUNT(*) FROM clothing_items WHERE user_id = $1"
	if err := r.db.WithContext(ctx).Raw(countQuery, userIDStr).Scan(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count clothing items: %w", err)
	}

	query := "SELECT * FROM clothing_items WHERE user_id = $1 LIMIT $2 OFFSET $3"
	if err := r.db.WithContext(ctx).Raw(query, userIDStr, limit, offset).Scan(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list clothing items: %w", err)
	}

	return items, total, nil
}

func (r *clothingRepository) Update(ctx context.Context, item *models.ClothingItem) error {
	if err := r.db.WithContext(ctx).Save(item).Error; err != nil {
		return fmt.Errorf("failed to update clothing item: %w", err)
	}
	return nil
}

func (r *clothingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.ClothingItem{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete clothing item: %w", err)
	}
	return nil
}
