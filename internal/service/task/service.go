package task

import (
	"context"
	"time"

	"github.com/google/uuid"

	"task-tracker/internal/models"
	taskrepo "task-tracker/internal/repository/task"
)

type Service interface {
	Add(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (models.Task, error)
	GetAll(ctx context.Context) ([]models.Task, error)
	Update(ctx context.Context, updatedTask *models.Task) error
}

type DefaultService struct {
	repo taskrepo.Repository
}

func NewDefaultService(repo taskrepo.Repository) *DefaultService {
	return &DefaultService{
		repo: repo,
	}
}

func (s *DefaultService) Add(ctx context.Context, task *models.Task) error {
	if exists, _ := s.repo.Exists(ctx, task.ID); exists {
		return models.ErrTaskExists
	}

	task.ID = uuid.New().String()
	task.CreatedAt = time.Now().Format(time.RFC3339Nano)
	task.UpdatedAt = task.CreatedAt

	return s.repo.Add(ctx, task)
}

func (s *DefaultService) Delete(ctx context.Context, id string) error {
	if exists, _ := s.repo.Exists(ctx, id); !exists {
		return models.ErrTaskNotFound
	}

	return s.repo.Delete(ctx, id)
}

func (s *DefaultService) Get(ctx context.Context, id string) (models.Task, error) {
	if exists, _ := s.repo.Exists(ctx, id); !exists {
		return models.Task{}, models.ErrTaskNotFound
	}

	return s.repo.Get(ctx, id)
}

func (s *DefaultService) GetAll(ctx context.Context) ([]models.Task, error) {
	return s.repo.GetAll(ctx)
}

func (s *DefaultService) Update(ctx context.Context, updatedTask *models.Task) error {
	if exists, _ := s.repo.Exists(ctx, updatedTask.ID); !exists {
		return models.ErrTaskNotFound
	}

	return s.repo.Update(ctx, updatedTask)
}
