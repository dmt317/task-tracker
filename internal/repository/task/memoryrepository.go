package task

import (
	"context"
	"sync"
	"time"

	"task-tracker/internal/models"
)

type MemoryRepository struct {
	store map[string]models.Task
	mu    sync.Mutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		store: make(map[string]models.Task),
	}
}

func (repo *MemoryRepository) Add(_ context.Context, task *models.Task) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.store[task.ID] = *task

	return nil
}

func (repo *MemoryRepository) Delete(_ context.Context, id string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	delete(repo.store, id)

	return nil
}

func (repo *MemoryRepository) Exists(_ context.Context, id string) (bool, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	_, found := repo.store[id]

	return found, nil
}

func (repo *MemoryRepository) Get(_ context.Context, id string) (models.Task, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	task := repo.store[id]

	return task, nil
}

func (repo *MemoryRepository) GetAll(_ context.Context) ([]models.Task, error) {
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

func (repo *MemoryRepository) Update(_ context.Context, updatedTask *models.Task) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	task := repo.store[updatedTask.ID]

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
