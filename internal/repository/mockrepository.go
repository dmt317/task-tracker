package repository

import (
	"context"
	"fmt"

	"task-tracker/internal/models"
)

type MockTaskRepository struct {
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

func (repo *MockTaskRepository) Add(ctx context.Context, task *models.Task) error {
	if repo.ForceRepositoryError {
		return ErrAddingTask
	}

	return nil
}

func (repo *MockTaskRepository) Delete(ctx context.Context, id string) error {
	if repo.ForceRepositoryError {
		return ErrDeletingTask
	}

	return nil
}

func (repo *MockTaskRepository) Exists(ctx context.Context, id string) (bool, error) {
	return repo.IsExist, nil
}

func (repo *MockTaskRepository) Get(ctx context.Context, id string) (models.Task, error) {
	if repo.ForceRepositoryError {
		return models.Task{}, ErrGettingTask
	}

	return models.Task{ID: id, Title: "Mock Task"}, nil
}

func (repo *MockTaskRepository) GetAll(ctx context.Context) ([]models.Task, error) {
	if repo.ForceRepositoryError {
		return []models.Task{}, ErrGettingAllTasks
	}

	return []models.Task{{ID: "task1", Title: "Mock Task"}}, nil
}

func (repo *MockTaskRepository) Update(ctx context.Context, updatedTask *models.Task) error {
	if repo.ForceRepositoryError {
		return ErrUpdatingTask
	}

	return nil
}
