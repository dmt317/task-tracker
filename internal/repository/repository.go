package repository

import (
	"context"

	"task-tracker/internal/models"
)

type TaskRepository interface {
	Add(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
	Get(ctx context.Context, id string) (models.Task, error)
	GetAll(ctx context.Context) ([]models.Task, error)
	Update(ctx context.Context, updatedTask *models.Task) error
}
