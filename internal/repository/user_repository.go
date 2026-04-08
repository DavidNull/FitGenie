package repository

import (
	"context"
	"fmt"

	"fitgenie/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]models.User, int64, error)
	CreateStyleProfile(ctx context.Context, profile *models.StyleProfile) error
	GetStyleProfile(ctx context.Context, userID uuid.UUID) (*models.StyleProfile, error)
	CreateColorProfile(ctx context.Context, profile *models.ColorProfile) error
	GetColorProfile(ctx context.Context, userID uuid.UUID) (*models.ColorProfile, error)
	AddFavoriteOutfit(ctx context.Context, favorite *models.FavoriteOutfit) error
	RemoveFavoriteOutfit(ctx context.Context, userID, outfitID uuid.UUID) error
	ListFavoriteOutfits(ctx context.Context, userID uuid.UUID) ([]models.FavoriteOutfit, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).
		Preload("StyleProfile").
		Preload("ColorProfile").
		Preload("ClothingItems").
		Preload("Outfits").
		First(&user, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (r *userRepository) List(ctx context.Context, offset, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}

func (r *userRepository) CreateStyleProfile(ctx context.Context, profile *models.StyleProfile) error {
	var existing models.StyleProfile
	err := r.db.WithContext(ctx).First(&existing, "user_id = ?", profile.UserID).Error
	if err == nil {
		profile.ID = existing.ID
		return r.db.WithContext(ctx).Save(profile).Error
	}

	if err := r.db.WithContext(ctx).Create(profile).Error; err != nil {
		return fmt.Errorf("failed to create style profile: %w", err)
	}
	return nil
}

func (r *userRepository) GetStyleProfile(ctx context.Context, userID uuid.UUID) (*models.StyleProfile, error) {
	var profile models.StyleProfile
	if err := r.db.WithContext(ctx).First(&profile, "user_id = ?", userID).Error; err != nil {
		return nil, fmt.Errorf("failed to get style profile: %w", err)
	}
	return &profile, nil
}

func (r *userRepository) CreateColorProfile(ctx context.Context, profile *models.ColorProfile) error {
	var existing models.ColorProfile
	err := r.db.WithContext(ctx).First(&existing, "user_id = ?", profile.UserID).Error
	if err == nil {
		profile.ID = existing.ID
		return r.db.WithContext(ctx).Save(profile).Error
	}

	if err := r.db.WithContext(ctx).Create(profile).Error; err != nil {
		return fmt.Errorf("failed to create color profile: %w", err)
	}
	return nil
}

func (r *userRepository) GetColorProfile(ctx context.Context, userID uuid.UUID) (*models.ColorProfile, error) {
	var profile models.ColorProfile
	if err := r.db.WithContext(ctx).First(&profile, "user_id = ?", userID).Error; err != nil {
		return nil, fmt.Errorf("failed to get color profile: %w", err)
	}
	return &profile, nil
}

func (r *userRepository) AddFavoriteOutfit(ctx context.Context, favorite *models.FavoriteOutfit) error {
	var existing models.FavoriteOutfit
	if err := r.db.WithContext(ctx).
		First(&existing, "user_id = ? AND outfit_id = ?", favorite.UserID, favorite.OutfitID).Error; err == nil {
		return fmt.Errorf("outfit already favorited")
	}

	if err := r.db.WithContext(ctx).Create(favorite).Error; err != nil {
		return fmt.Errorf("failed to add favorite: %w", err)
	}
	return nil
}

func (r *userRepository) RemoveFavoriteOutfit(ctx context.Context, userID, outfitID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Delete(&models.FavoriteOutfit{}, "user_id = ? AND outfit_id = ?", userID, outfitID).Error; err != nil {
		return fmt.Errorf("failed to remove favorite: %w", err)
	}
	return nil
}

func (r *userRepository) ListFavoriteOutfits(ctx context.Context, userID uuid.UUID) ([]models.FavoriteOutfit, error) {
	var favorites []models.FavoriteOutfit
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&favorites).Error; err != nil {
		return nil, fmt.Errorf("failed to list favorites: %w", err)
	}
	return favorites, nil
}
