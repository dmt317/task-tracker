package storage

import (
	"sync"
	"task-tracker/models"
	"task-tracker/task"
	"time"
)

type Storage struct {
	store map[string]task.Task
	mu    sync.Mutex
}

func (s *Storage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if id == "" {
		return models.ErrIdIsEmpty
	}

	if _, found := s.store[id]; !found {
		return models.ErrTaskNotFound
	}

	delete(s.store, id)

	return nil
}

func (s *Storage) Add(task task.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task.Id == "" {
		return models.ErrIdIsEmpty
	}

	if _, found := s.store[task.Id]; found {
		return models.ErrTaskExists
	}

	task.CreatedAt = time.Now().Format(time.RFC3339Nano)
	task.UpdatedAt = task.CreatedAt

	s.store[task.Id] = task

	return nil
}

func (s *Storage) Update(updatedTask task.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if updatedTask.Id == "" {
		return models.ErrIdIsEmpty
	}

	task, found := s.store[updatedTask.Id]

	if !found {
		return models.ErrTaskNotFound
	}

	task.Description = updatedTask.Description
	task.UpdatedAt = time.Now().Format(time.RFC3339Nano)

	s.store[updatedTask.Id] = task

	return nil
}

func (s *Storage) Get(id string) (task.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if id == "" {
		return task.Task{}, models.ErrIdIsEmpty
	}

	if _, found := s.store[id]; !found {
		return task.Task{}, models.ErrTaskNotFound
	}

	return s.store[id], nil
}
