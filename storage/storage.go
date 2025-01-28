package storage

import (
	"fmt"
	"task-cli/task"
)

type Storage struct {
	store map[string]task.Task
}

func (s *Storage) Delete(id string) error {
	if _, found := s.store[id]; !found {
		return fmt.Errorf("task with id %s does not exist", id)
	}

	delete(s.store, id)

	return nil
}

func (s *Storage) Add(task task.Task) error {
	if _, found := s.store[task.Id]; found {
		return fmt.Errorf("task with id %s already exists", task.Id)
	}

	s.store[task.Id] = task

	return nil
}

func (s *Storage) Update(task task.Task) error {
	if _, found := s.store[task.Id]; !found {
		return fmt.Errorf("task with id %s does not exist", task.Id)
	}

	s.store[task.Id] = task

	return nil
}

func (s *Storage) Get(id string) (task.Task, error) {
	if _, found := s.store[id]; !found {
		return task.Task{}, fmt.Errorf("task with id %s does not exist", id)
	}

	return s.store[id], nil
}
