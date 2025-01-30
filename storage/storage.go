package storage

import (
	"sync"
	"task-cli/errors"
	"task-cli/task"
)

type Storage struct {
	store map[string]task.Task
	mu    sync.Mutex
}

func (s *Storage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, found := s.store[id]; !found {
		return errors.ErrTaskNotFound
	}

	delete(s.store, id)

	return nil
}

func (s *Storage) Add(task task.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, found := s.store[task.Id]; found {
		return errors.ErrTaskExists
	}

	s.store[task.Id] = task

	return nil
}

func (s *Storage) Update(task task.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, found := s.store[task.Id]; !found {
		return errors.ErrTaskNotFound
	}

	s.store[task.Id] = task

	return nil
}

func (s *Storage) Get(id string) (task.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, found := s.store[id]; !found {
		return task.Task{}, errors.ErrTaskNotFound
	}

	return s.store[id], nil
}
