package task

import (
	"context"
	"fmt"

	"task-tracker/internal/models"
)

type MockRepository struct {
	ForceRepositoryError bool
	IsExist              bool
}

var (
	ErrAddingTask      = fmt.Errorf("error adding task")
	ErrDeletingTask    = fmt.Errorf("error deleting task")
	ErrGettingTask     = fmt.Errorf("error getting task")
	ErrGettingAllTasks = fmt.Errorf("error getting all tasks")
	ErrUpdatingTask    = fmt.Errorf("error updating task")
)

func (repo *MockRepository) Add(_ context.Context, _ *models.Task) error {
	if repo.ForceRepositoryError {
		return ErrAddingTask
	}

	return nil
}

func (repo *MockRepository) Delete(_ context.Context, _ string) error {
	if repo.ForceRepositoryError {
		return ErrDeletingTask
	}

	return nil
}

func (repo *MockRepository) Exists(_ context.Context, _ string) (bool, error) {
	return repo.IsExist, nil
}

func (repo *MockRepository) Get(_ context.Context, id string) (models.Task, error) {
	if repo.ForceRepositoryError {
		return models.Task{}, ErrGettingTask
	}

	return models.Task{ID: id, Title: "Mock Task"}, nil
}

func (repo *MockRepository) GetAll(_ context.Context) ([]models.Task, error) {
	if repo.ForceRepositoryError {
		return []models.Task{}, ErrGettingAllTasks
	}

	return []models.Task{{ID: "task1", Title: "Mock Task"}}, nil
}

func (repo *MockRepository) Update(_ context.Context, _ *models.Task) error {
	if repo.ForceRepositoryError {
		return ErrUpdatingTask
	}

	return nil
}
