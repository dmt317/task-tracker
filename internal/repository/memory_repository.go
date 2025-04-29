package repository

import (
	"sync"
	"time"

	"task-tracker/internal/models"
)

type MemoryTaskRepository struct {
	store map[string]models.Task
	mu    sync.Mutex
}

func NewMemoryTaskRepository() *MemoryTaskRepository {
	return &MemoryTaskRepository{
		store: make(map[string]models.Task),
	}
}

func (repo *MemoryTaskRepository) Add(task *models.Task) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, found := repo.store[task.ID]; found {
		return models.ErrTaskExists
	}

	repo.store[task.ID] = *task

	return nil
}

func (repo *MemoryTaskRepository) Delete(id string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, found := repo.store[id]; !found {
		return models.ErrTaskNotFound
	}

	delete(repo.store, id)

	return nil
}

func (repo *MemoryTaskRepository) Get(id string) (models.Task, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	task, found := repo.store[id]

	if !found {
		return models.Task{}, models.ErrTaskNotFound
	}

	return task, nil
}

func (repo *MemoryTaskRepository) GetAll() ([]models.Task, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	tasks := make([]models.Task, len(repo.store))
	i := 0

	for _, task := range repo.store {
		tasks[i] = task
		i++
	}

	return tasks, nil
}

func (repo *MemoryTaskRepository) Update(updatedTask *models.Task) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	task, found := repo.store[updatedTask.ID]

	if !found {
		return models.ErrTaskNotFound
	}

	updated := false

	if updatedTask.Title != "" && updatedTask.Title != task.Title {
		task.Title = updatedTask.Title
		updated = true
	}

	if updatedTask.Description != "" && updatedTask.Description != task.Description {
		task.Description = updatedTask.Description
		updated = true
	}

	if updatedTask.Status != "" && updatedTask.Status != task.Status {
		task.Status = updatedTask.Status
		updated = true
	}

	if updated {
		task.UpdatedAt = time.Now().Format(time.RFC3339Nano)
		repo.store[updatedTask.ID] = task
	}

	return nil
}
