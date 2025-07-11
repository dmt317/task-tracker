package repository

import (
	"task-tracker/internal/models"
)

type TaskRepository interface {
	Add(task *models.Task) error
	Delete(id string) error
	Exists(id string) (bool, error)
	Get(id string) (models.Task, error)
	GetAll() ([]models.Task, error)
	Update(updatedTask *models.Task) error
}
