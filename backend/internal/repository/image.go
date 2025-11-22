package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Helltale/beer-mania/backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrImageNotFound = errors.New("image not found")
)

type ImageRepository interface {
	Create(ctx context.Context, image *entity.Image) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Image, error)
	Update(ctx context.Context, image *entity.Image) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.ImageStatus) error
}

type imageRepository struct {
	db *gorm.DB
}

func NewImageRepository(db *gorm.DB) ImageRepository {
	return &imageRepository{db: db}
}

func (r *imageRepository) Create(ctx context.Context, image *entity.Image) error {
	if err := r.db.WithContext(ctx).Create(image).Error; err != nil {
		return err
	}
	return nil
}

func (r *imageRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Image, error) {
	var image entity.Image
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&image).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrImageNotFound
		}
		return nil, err
	}
	return &image, nil
}

func (r *imageRepository) Update(ctx context.Context, image *entity.Image) error {
	image.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Save(image).Error; err != nil {
		return err
	}
	return nil
}

func (r *imageRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.ImageStatus) error {
	result := r.db.WithContext(ctx).Model(&entity.Image{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":     status,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrImageNotFound
	}

	return nil
}
