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
	// ErrTaskNotFound is returned when task is not found
	ErrTaskNotFound = errors.New("task not found")
)

// TaskRepository defines the interface for processing task repository operations
type TaskRepository interface {
	Create(ctx context.Context, task *entity.ProcessingTask) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ProcessingTask, error)
	GetByImageID(ctx context.Context, imageID uuid.UUID) (*entity.ProcessingTask, error)
	Update(ctx context.Context, task *entity.ProcessingTask) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.TaskStatus, errorMsg *string) error
}

type taskRepository struct {
	db *gorm.DB
}

// NewTaskRepository creates a new TaskRepository
func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

// Create creates a new processing task record
func (r *taskRepository) Create(ctx context.Context, task *entity.ProcessingTask) error {
	if err := r.db.WithContext(ctx).Create(task).Error; err != nil {
		return err
	}
	return nil
}

// GetByID retrieves a processing task by ID
func (r *taskRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ProcessingTask, error) {
	var task entity.ProcessingTask
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return &task, nil
}

// GetByImageID retrieves a processing task by image ID
func (r *taskRepository) GetByImageID(ctx context.Context, imageID uuid.UUID) (*entity.ProcessingTask, error) {
	var task entity.ProcessingTask
	if err := r.db.WithContext(ctx).Where("image_id = ?", imageID).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return &task, nil
}

// Update updates an existing processing task record
func (r *taskRepository) Update(ctx context.Context, task *entity.ProcessingTask) error {
	task.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Save(task).Error; err != nil {
		return err
	}
	return nil
}

// UpdateStatus updates the status and optionally error message of a task
func (r *taskRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.TaskStatus, errorMsg *string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}
	if errorMsg != nil {
		updates["error_message"] = *errorMsg
	}

	result := r.db.WithContext(ctx).Model(&entity.ProcessingTask{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}
