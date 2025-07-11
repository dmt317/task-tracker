package service

import (
	"context"
	"errors"
	"testing"

	"task-tracker/internal/models"
	"task-tracker/internal/repository"
)

func TestAdd(t *testing.T) {
	tests := map[string]struct {
		service *DefaultTaskService
		result  error
	}{
		"successfully adds a valid task": {
			service: &DefaultTaskService{
				repo: &repository.MockTaskRepository{
					ForceRepositoryError: false,
					IsExist:              false,
				},
			},
			result: nil,
		},

		"add task fails when task already exists": {
			service: &DefaultTaskService{
				repo: &repository.MockTaskRepository{
					ForceRepositoryError: false,
					IsExist:              true,
				},
			},
			result: models.ErrTaskExists,
		},

		"add task fails due to repository error": {
			service: &DefaultTaskService{
				repo: &repository.MockTaskRepository{
					ForceRepositoryError: true,
					IsExist:              false,
				},
			},
			result: repository.ErrAddingTask,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := test.service.Add(context.Background(), &models.Task{})

			if !errors.Is(err, test.result) {
				t.Fatalf("test-case: (%q); returned %v; expected %v", name, err, test.result)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := map[string]struct {
		service *DefaultTaskService
		result  error
	}{
		"successfully delete a valid task": {
			service: &DefaultTaskService{
				repo: &repository.MockTaskRepository{
					ForceRepositoryError: false,
					IsExist:              true,
				},
			},
			result: nil,
		},

		"delete task fails when task doesn't exist": {
			service: &DefaultTaskService{
				repo: &repository.MockTaskRepository{
					ForceRepositoryError: false,
					IsExist:              false,
				},
			},
			result: models.ErrTaskNotFound,
		},

		"delete task fails due to repository error": {
			service: &DefaultTaskService{
				repo: &repository.MockTaskRepository{
					ForceRepositoryError: true,
					IsExist:              true,
				},
			},
			result: repository.ErrDeletingTask,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := test.service.Delete(context.Background(), "task1")

			if !errors.Is(err, test.result) {
				t.Fatalf("test-case: (%q); returned %v; expected %v", name, err, test.result)
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := map[string]struct {
		service *DefaultTaskService
		result  error
	}{
		"successfully get a valid task": {
			service: &DefaultTaskService{
				repo: &repository.MockTaskRepository{
					ForceRepositoryError: false,
					IsExist:              true,
				},
			},
			result: nil,
		},

		"get task fails when task doesn't exist": {
			service: &DefaultTaskService{
				repo: &repository.MockTaskRepository{
					ForceRepositoryError: false,
					IsExist:              false,
				},
			},
			result: models.ErrTaskNotFound,
		},

		"get task fails due to repository error": {
			service: &DefaultTaskService{
				repo: &repository.MockTaskRepository{
					ForceRepositoryError: true,
					IsExist:              true,
				},
			},
			result: repository.ErrGettingTask,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := test.service.Get(context.Background(), "task1")

			if !errors.Is(err, test.result) {
				t.Fatalf("test-case: (%q); returned %v; expected %v", name, err, test.result)
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	tests := map[string]struct {
		service *DefaultTaskService
		result  error
	}{
		"successfully get all tasks": {
			service: &DefaultTaskService{
				repo: &repository.MockTaskRepository{
					ForceRepositoryError: false,
				},
			},
			result: nil,
		},

		"get all tasks fails due to repository error": {
			service: &DefaultTaskService{
				repo: &repository.MockTaskRepository{
					ForceRepositoryError: true,
				},
			},
			result: repository.ErrGettingAllTasks,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := test.service.GetAll(context.Background())

			if !errors.Is(err, test.result) {
				t.Fatalf("test-case: (%q); returned %v; expected %v", name, err, test.result)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := map[string]struct {
		service *DefaultTaskService
		result  error
	}{
		"successfully update a valid task": {
			service: &DefaultTaskService{
				repo: &repository.MockTaskRepository{
					ForceRepositoryError: false,
					IsExist:              true,
				},
			},
			result: nil,
		},

		"update task fails when task doesn't exist": {
			service: &DefaultTaskService{
				repo: &repository.MockTaskRepository{
					ForceRepositoryError: false,
					IsExist:              false,
				},
			},
			result: models.ErrTaskNotFound,
		},

		"update task fails due to repository error": {
			service: &DefaultTaskService{
				repo: &repository.MockTaskRepository{
					ForceRepositoryError: true,
					IsExist:              true,
				},
			},
			result: repository.ErrUpdatingTask,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := test.service.Update(context.Background(), &models.Task{})

			if !errors.Is(err, test.result) {
				t.Fatalf("test-case: (%q); returned %v; expected %v", name, err, test.result)
			}
		})
	}
}
