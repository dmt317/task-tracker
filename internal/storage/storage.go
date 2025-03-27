package storage

import (
	"sync"
	"time"

	"task-tracker/internal/models"
)

type Storage struct {
	store map[string]models.Task
	mu    sync.Mutex
}

func NewStorage() *Storage {
	return &Storage{
		store: make(map[string]models.Task),
		mu:    sync.Mutex{},
	}
}

func (s *Storage) Delete(ID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if ID == "" {
		return models.ErrIdIsEmpty
	}

	if _, found := s.store[ID]; !found {
		return models.ErrTaskNotFound
	}

	delete(s.store, ID)

	return nil
}

func (s *Storage) Add(task *models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task.ID == "" {
		return models.ErrIdIsEmpty
	}

	if _, found := s.store[task.ID]; found {
		return models.ErrTaskExists
	}

	task.CreatedAt = time.Now().Format(time.RFC3339Nano)
	task.UpdatedAt = task.CreatedAt

	s.store[task.ID] = *task

	return nil
}

func (s *Storage) Update(updatedTask models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if updatedTask.ID == "" {
		return models.ErrIdIsEmpty
	}

	task, found := s.store[updatedTask.ID]

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
		s.store[updatedTask.ID] = task
	}

	return nil
}

func (s *Storage) Get(ID string) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if ID == "" {
		return models.Task{}, models.ErrIdIsEmpty
	}

	if _, found := s.store[ID]; !found {
		return models.Task{}, models.ErrTaskNotFound
	}

	return s.store[ID], nil
}

func (s *Storage) GetAll() ([]models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.store == nil {
		return nil, models.ErrStorageNotInitialized
	}

	tasks := make([]models.Task, len(s.store))
	i := 0
	for _, task := range s.store {
		tasks[i] = task
		i++
	}

	return tasks, nil
}
