package storage

import (
	"fmt"
	"task-cli/errors"
	"task-cli/task"
	"testing"
	"time"
)

func TestStorage_Add(t *testing.T) {
	tests := map[string]struct {
		inputTasks []task.Task
		initMap    Storage
		result     []error
	}{
		"add valid task": {
			inputTasks: []task.Task{
				{Id: "task1", Description: "Valid task", CreatedAt: time.Now().Format(time.RFC3339Nano)},
			},
			initMap: Storage{store: make(map[string]task.Task)},
			result:  []error{nil},
		},

		"add task with empty id": {
			inputTasks: []task.Task{
				{Id: "", Description: "No ID", CreatedAt: time.Now().Format(time.RFC3339Nano)},
			},
			initMap: Storage{store: make(map[string]task.Task)},
			result:  []error{errors.ErrIdIsEmpty},
		},

		"add task with duplicate id": {
			inputTasks: []task.Task{
				{Id: "task1", Description: "First task", CreatedAt: time.Now().Format(time.RFC3339Nano)},
				{Id: "task1", Description: "Duplicate task", CreatedAt: time.Now().Format(time.RFC3339Nano)},
			},
			initMap: Storage{store: map[string]task.Task{}},
			result:  []error{nil, errors.ErrTaskExists},
		},

		"add task with empty description": {
			inputTasks: []task.Task{
				{Id: "task2", Description: "", CreatedAt: time.Now().Format(time.RFC3339Nano)},
			},
			initMap: Storage{store: make(map[string]task.Task)},
			result:  []error{nil},
		},

		"add multiple valid tasks": {
			inputTasks: []task.Task{
				{Id: "task3", Description: "Task 3", CreatedAt: time.Now().Format(time.RFC3339Nano)},
				{Id: "task4", Description: "Task 4", CreatedAt: time.Now().Format(time.RFC3339Nano)},
				{Id: "task5", Description: "Task 5", CreatedAt: time.Now().Format(time.RFC3339Nano)},
			},
			initMap: Storage{store: make(map[string]task.Task)},
			result:  []error{nil, nil, nil},
		},

		"add duplicate task when task already in storage": {
			inputTasks: []task.Task{
				{Id: "task1", Description: "Duplicate task", CreatedAt: time.Now().Format(time.RFC3339Nano)},
			},
			initMap: Storage{store: map[string]task.Task{
				"task1": {Id: "task1", Description: "First task", CreatedAt: time.Now().Format(time.RFC3339Nano)}}},
			result: []error{errors.ErrTaskExists},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s := tc.initMap
			for i, task := range tc.inputTasks {
				got := s.Add(task)
				expected := tc.result[i]
				if !errors.Is(got, expected) {
					t.Fatalf("test-case: (%q); returned %q; expected %q", name, got, expected)
				}
				if got == nil {
					currTask, err := s.Get(task.Id)
					if err == nil && currTask == task {
						fmt.Println(currTask)
					}
				}
			}
		})
	}
}

func TestStorage_Get(t *testing.T) {
	s := Storage{store: make(map[string]task.Task)}

	s.store["task1"] = task.Task{Id: "task1", Description: "Complete the boss's task", CreatedAt: time.Now().Format(time.RFC3339Nano)}

	tests := map[string]struct {
		input  string
		result struct {
			task task.Task
			err  error
		}
	}{
		"get task successfully": {
			input: "task1",
			result: struct {
				task task.Task
				err  error
			}{
				task: s.store["task1"], err: nil,
			},
		},
		"get task that does not exist": {
			input: "task2",
			result: struct {
				task task.Task
				err  error
			}{
				task: task.Task{}, err: errors.ErrTaskNotFound,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := s.Get(test.input)
			expected := test.result
			if got != expected.task && !errors.Is(err, expected.err) {
				t.Fatalf("test-case: (%q); returned task: %q, err: %q; expected task: %q, err: %q", name, got, err, expected.task, expected.err)
			}
		})
	}
}

func TestStorage_Update(t *testing.T) {
	s := Storage{store: make(map[string]task.Task)}

	s.store["task1"] = task.Task{Id: "task1", Description: "Complete the boss's task", CreatedAt: time.Now().Format(time.RFC3339Nano)}

	tests := map[string]struct {
		input  task.Task
		result error
	}{
		"update task successfully": {
			input:  task.Task{Id: "task1", Description: "Add new functionality", CreatedAt: time.Now().Format(time.RFC3339Nano), UpdatedAt: time.Now().Format(time.RFC3339Nano)},
			result: nil,
		},
		"update task that doesn't exists": {
			input:  task.Task{Id: "task2", Description: "Add new functionality", CreatedAt: time.Now().Format(time.RFC3339Nano), UpdatedAt: time.Now().Format(time.RFC3339Nano)},
			result: errors.ErrTaskNotFound,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got, expected := s.Update(test.input), test.result; !errors.Is(got, expected) {
				t.Fatalf("test-case: (%q); returned %q; expected %q", name, got, expected)
			}
		})
	}
}

func TestStorage_Delete(t *testing.T) {
	s := Storage{store: make(map[string]task.Task)}

	s.store["task1"] = task.Task{Id: "task1", Description: "Complete the boss's task", CreatedAt: time.Now().Format(time.RFC3339Nano)}

	tests := map[string]struct {
		input  string
		result error
	}{
		"delete task successfully": {
			input:  "task1",
			result: nil,
		},
		"delete task that does not exist": {
			input:  "task2",
			result: errors.ErrTaskNotFound,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got, expected := s.Delete(test.input), test.result; !errors.Is(got, expected) {
				t.Fatalf("test-case: (%q); returned %q; expected %q", name, got, expected)
			}

		})
	}
}
