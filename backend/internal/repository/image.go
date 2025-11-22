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
	// ErrImageNotFound is returned when image is not found
	ErrImageNotFound = errors.New("image not found")
)

// ImageRepository defines the interface for image repository operations
type ImageRepository interface {
	Create(ctx context.Context, image *entity.Image) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Image, error)
	Update(ctx context.Context, image *entity.Image) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.ImageStatus) error
}

type imageRepository struct {
	db *gorm.DB
}

// NewImageRepository creates a new ImageRepository
func NewImageRepository(db *gorm.DB) ImageRepository {
	return &imageRepository{db: db}
}

// Create creates a new image record
func (r *imageRepository) Create(ctx context.Context, image *entity.Image) error {
	if err := r.db.WithContext(ctx).Create(image).Error; err != nil {
		return err
	}
	return nil
}

// GetByID retrieves an image by ID
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

// Update updates an existing image record
func (r *imageRepository) Update(ctx context.Context, image *entity.Image) error {
	image.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Save(image).Error; err != nil {
		return err
	}
	return nil
}

// UpdateStatus updates only the status of an image
func (r *imageRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.ImageStatus) error {
	result := r.db.WithContext(ctx).Model(&entity.Image{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
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
