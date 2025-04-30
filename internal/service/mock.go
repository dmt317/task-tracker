package service

import (
	"errors"
	"net/http"

	"task-tracker/internal/models"
)

var ErrInternalMock = models.Error{
	Err:        errors.New("internal server error"),
	StatusCode: http.StatusInternalServerError,
}

const NotFound = "not_found"

type TaskServiceMock struct {
	ForceInternalError bool
}

func (m *TaskServiceMock) Add(task *models.Task) error {
	if task.Title == "existing" {
		return models.ErrTaskExists
	}

	if m.ForceInternalError {
		return ErrInternalMock
	}

	return nil
}

func (m *TaskServiceMock) Delete(id string) error {
	if id == NotFound {
		return models.ErrTaskNotFound
	}

	if m.ForceInternalError {
		return ErrInternalMock
	}

	return nil
}

func (m *TaskServiceMock) Get(id string) (models.Task, error) {
	if id == NotFound {
		return models.Task{}, models.ErrTaskNotFound
	}

	if m.ForceInternalError {
		return models.Task{}, ErrInternalMock
	}

	return models.Task{ID: id, Title: "Mock Task"}, nil
}

func (m *TaskServiceMock) GetAll() ([]models.Task, error) {
	if m.ForceInternalError {
		return []models.Task{}, ErrInternalMock
	}

	return []models.Task{{ID: "task1", Title: "Mock Task"}}, nil
}

func (m *TaskServiceMock) Update(updatedTask *models.Task) error {
	if updatedTask.ID == NotFound {
		return models.ErrTaskNotFound
	}

	if m.ForceInternalError {
		return ErrInternalMock
	}

	return nil
}
