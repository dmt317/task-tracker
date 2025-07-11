package repository

import (
	"context"
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

func (repo *MemoryTaskRepository) Add(ctx context.Context, task *models.Task) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.store[task.ID] = *task

	return nil
}

func (repo *MemoryTaskRepository) Delete(ctx context.Context, id string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	delete(repo.store, id)

	return nil
}

func (repo *MemoryTaskRepository) Exists(ctx context.Context, id string) (bool, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	_, found := repo.store[id]

	return found, nil
}

func (repo *MemoryTaskRepository) Get(ctx context.Context, id string) (models.Task, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	task := repo.store[id]

	return task, nil
}

func (repo *MemoryTaskRepository) GetAll(ctx context.Context) ([]models.Task, error) {
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

func (repo *MemoryTaskRepository) Update(ctx context.Context, updatedTask *models.Task) error {
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
