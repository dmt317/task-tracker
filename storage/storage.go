package storage

import (
	"sync"
	"task-tracker/models"
	"task-tracker/task"
)

type Storage struct {
	store map[string]task.Task
	mu    sync.Mutex
}

func (s *Storage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

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

	s.store[task.Id] = task

	return nil
}

func (s *Storage) Update(task task.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, found := s.store[task.Id]; !found {
		return models.ErrTaskNotFound
	}

	s.store[task.Id] = task

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
