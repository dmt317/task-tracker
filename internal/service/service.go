package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"task-tracker/internal/models"
	"task-tracker/internal/repository"
)

type TaskService interface {
	Add(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (models.Task, error)
	GetAll(ctx context.Context) ([]models.Task, error)
	Update(ctx context.Context, updatedTask *models.Task) error
}

type DefaultTaskService struct {
	repo repository.TaskRepository
}

func NewDefaultTaskService(repo repository.TaskRepository) *DefaultTaskService {
	return &DefaultTaskService{
		repo: repo,
	}
}

func (s *DefaultTaskService) Add(ctx context.Context, task *models.Task) error {
	if exists, _ := s.repo.Exists(ctx, task.ID); exists {
		return models.ErrTaskExists
	}

	task.ID = uuid.New().String()
	task.CreatedAt = time.Now().Format(time.RFC3339Nano)
	task.UpdatedAt = task.CreatedAt

	return s.repo.Add(ctx, task)
}

func (s *DefaultTaskService) Delete(ctx context.Context, id string) error {
	if exists, _ := s.repo.Exists(ctx, id); !exists {
		return models.ErrTaskNotFound
	}

	return s.repo.Delete(ctx, id)
}

func (s *DefaultTaskService) Get(ctx context.Context, id string) (models.Task, error) {
	if exists, _ := s.repo.Exists(ctx, id); !exists {
		return models.Task{}, models.ErrTaskNotFound
	}

	return s.repo.Get(ctx, id)
}

func (s *DefaultTaskService) GetAll(ctx context.Context) ([]models.Task, error) {
	return s.repo.GetAll(ctx)
}

func (s *DefaultTaskService) Update(ctx context.Context, updatedTask *models.Task) error {
	if exists, _ := s.repo.Exists(ctx, updatedTask.ID); !exists {
		return models.ErrTaskNotFound
	}

	return s.repo.Update(ctx, updatedTask)
}
