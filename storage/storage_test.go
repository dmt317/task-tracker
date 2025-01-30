package storage

import (
	"task-cli/errors"
	"task-cli/task"
	"testing"
	"time"
)

func TestStorage_Add(t *testing.T) {
	s := Storage{store: make(map[string]task.Task)}

	tests := map[string]struct {
		input  task.Task
		result error
	}{
		"add task successfully": {
			input:  task.Task{Id: "task1", Description: "Complete the boss's task", CreatedAt: time.Now().Format(time.RFC3339Nano)},
			result: nil,
		},
		"add task that already exists": {
			input:  task.Task{Id: "task1", Description: "Complete the boss's task", CreatedAt: time.Now().Format(time.RFC3339Nano)},
			result: errors.ErrTaskExists,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got, expected := s.Add(test.input), test.result; !errors.Is(got, expected) {
				t.Fatalf("test-case: (%q); returned %q; expected %q", name, got, expected)
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
