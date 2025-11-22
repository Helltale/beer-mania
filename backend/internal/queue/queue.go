package queue

import (
	"context"

	"github.com/google/uuid"
)

type Queue interface {
	PublishTask(ctx context.Context, taskID uuid.UUID, imageID uuid.UUID) error
	ConsumeTasks(ctx context.Context, handler func(taskID uuid.UUID, imageID uuid.UUID) error) error
	Close() error
}
