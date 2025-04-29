package service

import (
	"time"

	"github.com/google/uuid"

	"task-tracker/internal/models"
	"task-tracker/internal/repository"
)

type TaskService interface {
	Add(task *models.Task) error
	Delete(id string) error
	Get(id string) (models.Task, error)
	GetAll() ([]models.Task, error)
	Update(updatedTask *models.Task) error
}

type DefaultTaskService struct {
	repo repository.TaskRepository
}

func NewDefaultTaskService(repo repository.TaskRepository) *DefaultTaskService {
	return &DefaultTaskService{
		repo: repo,
	}
}

func (s *DefaultTaskService) Add(task *models.Task) error {
	task.ID = uuid.New().String()
	task.CreatedAt = time.Now().Format(time.RFC3339Nano)
	task.UpdatedAt = task.CreatedAt

	return s.repo.Add(task)
}

func (s *DefaultTaskService) Delete(id string) error {
	if id == "" {
		return models.ErrIDIsEmpty
	}

	return s.repo.Delete(id)
}

func (s *DefaultTaskService) Get(id string) (models.Task, error) {
	if id == "" {
		return models.Task{}, models.ErrIDIsEmpty
	}

	return s.repo.Get(id)
}

func (s *DefaultTaskService) GetAll() ([]models.Task, error) {
	return s.repo.GetAll()
}

func (s *DefaultTaskService) Update(updatedTask *models.Task) error {
	if updatedTask.ID == "" {
		return models.ErrIDIsEmpty
	}

	return s.repo.Update(updatedTask)
}
