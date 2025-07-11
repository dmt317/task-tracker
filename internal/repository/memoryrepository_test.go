package repository

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"task-tracker/internal/models"
)

type TestResultGet struct {
	resultTasks  []models.Task
	resultErrors []error
}

type TestResultGetAll struct {
	resultTasks []models.Task
	resultError error
}

func TestStorage_Add(t *testing.T) {
	tests := map[string]struct {
		inputTasks []*models.Task
		storage    *MemoryTaskRepository
		result     []error
	}{
		"add valid task": {
			inputTasks: []*models.Task{
				{
					ID:          "task1",
					Title:       "Title",
					Description: "Valid task",
					Status:      "Todo",
				},
			},
			storage: NewMemoryTaskRepository(),
			result:  []error{nil},
		},

		"add task with duplicate id": {
			inputTasks: []*models.Task{
				{
					ID:          "task1",
					Title:       "Title",
					Description: "First task",
					Status:      "Todo",
				},
				{
					ID:          "task1",
					Title:       "Title",
					Description: "Duplicate task",
					Status:      "Todo",
				},
			},
			storage: NewMemoryTaskRepository(),
			result:  []error{nil, models.ErrTaskExists},
		},

		"add task with empty description": {
			inputTasks: []*models.Task{
				{
					ID:          "task2",
					Title:       "Empty description",
					Description: "",
					Status:      "Todo",
				},
			},
			storage: NewMemoryTaskRepository(),
			result:  []error{nil},
		},

		"add multiple valid tasks": {
			inputTasks: []*models.Task{
				{
					ID:          "task3",
					Title:       "Task 3",
					Description: "Task 3",
					Status:      "Todo",
				},
				{
					ID:          "task4",
					Title:       "Task 4",
					Description: "Task 4",
					Status:      "Todo",
				},
				{
					ID:          "task5",
					Title:       "Task 5",
					Description: "Task 5",
					Status:      "Todo",
				},
			},
			storage: NewMemoryTaskRepository(),
			result:  []error{nil, nil, nil},
		},

		"add duplicate task when task already in storage": {
			inputTasks: []*models.Task{
				{
					ID:          "task1",
					Title:       "Duplicate task",
					Description: "Duplicate task",
					Status:      "Todo",
				},
			},
			storage: &MemoryTaskRepository{
				store: map[string]models.Task{
					"task1": {
						ID:          "task1",
						Title:       "First task",
						Description: "First task",
						Status:      "Todo",
						CreatedAt:   time.Now().Format(time.RFC3339Nano),
						UpdatedAt:   time.Now().Format(time.RFC3339Nano),
					}},
			},
			result: []error{models.ErrTaskExists},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			for i, task := range test.inputTasks {
				gotErr := test.storage.Add(task)
				expectedErr := test.result[i]

				if !errors.Is(gotErr, expectedErr) {
					t.Fatalf("test-case: (%q); returned %q; expected %q", name, gotErr, expectedErr)
				}

				if gotErr != nil {
					return
				}

				currTask, err := test.storage.Get(task.ID)
				if err != nil {
					t.Fatalf("test-case: (%q); unexpected error: %q", name, err)
				}

				fmt.Println(currTask)
			}
		})
	}
}

func TestStorage_Get(t *testing.T) {
	fixedTime := time.Now().Format(time.RFC3339Nano)

	tests := map[string]struct {
		inputIDs []string
		storage  *MemoryTaskRepository
		result   TestResultGet
	}{
		"get existing task": {
			inputIDs: []string{"task1"},
			storage: &MemoryTaskRepository{
				store: map[string]models.Task{
					"task1": {
						ID:          "task1",
						Title:       "First task",
						Description: "First task",
						Status:      "Todo",
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					},
				},
			},
			result: TestResultGet{
				resultTasks: []models.Task{{
					ID:          "task1",
					Title:       "First task",
					Description: "First task",
					Status:      "Todo",
					CreatedAt:   fixedTime,
					UpdatedAt:   fixedTime,
				}},
				resultErrors: []error{nil},
			},
		},

		"get non-existing task when there are multiple tasks in the map": {
			inputIDs: []string{"task2"},
			storage: &MemoryTaskRepository{
				store: map[string]models.Task{
					"task1": {
						ID:          "task1",
						Title:       "First task",
						Description: "First task",
						Status:      "Todo",
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					},
					"task3": {
						ID:          "task3",
						Title:       "Third task",
						Description: "Third task",
						Status:      "Todo",
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					},
				},
			},
			result: TestResultGet{
				resultTasks:  []models.Task{{}},
				resultErrors: []error{models.ErrTaskNotFound},
			},
		},

		"get non-existing task when there are no tasks in the map": {
			inputIDs: []string{"task1"},
			storage:  NewMemoryTaskRepository(),
			result: TestResultGet{
				resultTasks:  []models.Task{{}},
				resultErrors: []error{models.ErrTaskNotFound},
			},
		},

		"get multiple tasks with duplicates": {
			inputIDs: []string{"task1", "task2", "task3", "task1", "task2", "task3"},
			storage: &MemoryTaskRepository{
				store: map[string]models.Task{
					"task1": {
						ID:          "task1",
						Title:       "First task",
						Description: "First task",
						Status:      "Todo",
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					},
					"task2": {
						ID:          "task2",
						Title:       "Second task",
						Description: "Second task",
						Status:      "Todo",
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					},
					"task3": {
						ID:          "task3",
						Title:       "Third task",
						Description: "Third task",
						Status:      "Todo",
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					},
				}},
			result: TestResultGet{
				resultTasks: []models.Task{
					{
						ID:          "task1",
						Title:       "First task",
						Description: "First task",
						Status:      "Todo",
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					},
					{
						ID:          "task2",
						Title:       "Second task",
						Description: "Second task",
						Status:      "Todo",
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					},
					{
						ID:          "task3",
						Title:       "Third task",
						Description: "Third task",
						Status:      "Todo",
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					},
					{
						ID:          "task1",
						Title:       "First task",
						Description: "First task",
						Status:      "Todo",
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					},
					{
						ID:          "task2",
						Title:       "Second task",
						Description: "Second task",
						Status:      "Todo",
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					},
					{
						ID:          "task3",
						Title:       "Third task",
						Description: "Third task",
						Status:      "Todo",
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					},
				},
				resultErrors: []error{nil, nil, nil, nil, nil, nil},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			for i, id := range test.inputIDs {
				task, err := test.storage.Get(id)
				resultTask := test.result.resultTasks[i]
				resultError := test.result.resultErrors[i]

				if !errors.Is(err, resultError) || task != resultTask {
					t.Fatalf("test-case: (%q); returned [%q %q]; expected [%q %q]", name, task, err, resultTask, resultError)
				}
			}
		})
	}
}

func TestStorage_Update(t *testing.T) {
	tests := map[string]struct {
		inputTasks []*models.Task
		storage    *MemoryTaskRepository
		result     []error
	}{
		"update existing task": {
			inputTasks: []*models.Task{
				{
					ID:          "task1",
					Title:       "New title",
					Description: "New description",
					Status:      "Done",
				},
			},
			storage: &MemoryTaskRepository{
				store: map[string]models.Task{
					"task1": {
						ID:          "task1",
						Title:       "Title",
						Description: "Description",
						Status:      "Todo",
						CreatedAt:   time.Now().Format(time.RFC3339Nano),
						UpdatedAt:   time.Now().Format(time.RFC3339Nano),
					},
				},
			},
			result: []error{nil},
		},

		"update non-existing task": {
			inputTasks: []*models.Task{
				{
					ID:          "task2",
					Title:       "New title",
					Description: "New description",
					Status:      "Done",
				},
			},
			storage: &MemoryTaskRepository{
				store: map[string]models.Task{
					"task1": {
						ID:          "task1",
						Title:       "Title",
						Description: "Description",
						Status:      "Todo",
						CreatedAt:   time.Now().Format(time.RFC3339Nano),
						UpdatedAt:   time.Now().Format(time.RFC3339Nano),
					},
				},
			},
			result: []error{models.ErrTaskNotFound},
		},

		"update multiple tasks": {
			inputTasks: []*models.Task{
				{
					ID:          "task1",
					Title:       "New title",
					Description: "New description",
					Status:      "Done",
				},
				{
					ID:          "task2",
					Title:       "New title",
					Description: "New description",
					Status:      "Done",
				},
				{
					ID:          "task3",
					Title:       "New title",
					Description: "New description",
					Status:      "Done",
				},
			},
			storage: &MemoryTaskRepository{
				store: map[string]models.Task{
					"task1": {
						ID:          "task1",
						Title:       "Title",
						Description: "Description",
						Status:      "Todo",
						CreatedAt:   time.Now().Format(time.RFC3339Nano),
						UpdatedAt:   time.Now().Format(time.RFC3339Nano),
					},
					"task2": {
						ID:          "task2",
						Title:       "Title",
						Description: "Description",
						Status:      "Todo",
						CreatedAt:   time.Now().Format(time.RFC3339Nano),
						UpdatedAt:   time.Now().Format(time.RFC3339Nano),
					},
					"task3": {
						ID:          "task3",
						Title:       "Title",
						Description: "Description",
						Status:      "Todo",
						CreatedAt:   time.Now().Format(time.RFC3339Nano),
						UpdatedAt:   time.Now().Format(time.RFC3339Nano),
					},
				},
			},
			result: []error{nil, nil, nil},
		},

		"update task using the same data": {
			inputTasks: []*models.Task{
				{
					ID:          "task1",
					Title:       "Title",
					Description: "Description",
					Status:      "Todo",
				},
			},
			storage: &MemoryTaskRepository{
				store: map[string]models.Task{
					"task1": {
						ID:          "task1",
						Title:       "Title",
						Description: "Description",
						Status:      "Todo",
						CreatedAt:   time.Now().Format(time.RFC3339Nano),
						UpdatedAt:   time.Now().Format(time.RFC3339Nano),
					},
				},
			},
			result: []error{nil},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			for i, task := range test.inputTasks {
				err := test.storage.Update(task)
				result := test.result[i]

				if !errors.Is(err, result) {
					t.Fatalf("test-case: (%q); returned [%q %q]", name, err, result)
				}

				if err != nil {
					return
				}

				updatedTask, err := test.storage.Get(task.ID)
				if err != nil {
					t.Fatalf("test-case: (%q); unexpected error: %q", name, err)
				}

				if updatedTask.Description != task.Description {
					t.Fatalf("test-case: (%q); task hasn't been updated; expected [%q]; got: [%q] ", name, task, updatedTask)
				}

				fmt.Println(updatedTask)
			}
		})
	}
}

func TestStorage_Delete(t *testing.T) {
	tests := map[string]struct {
		inputIDs []string
		storage  *MemoryTaskRepository
		result   []error
	}{
		"delete existing task": {
			inputIDs: []string{"task1"},
			storage: &MemoryTaskRepository{
				store: map[string]models.Task{
					"task1": {
						ID:          "task1",
						Title:       "Title",
						Description: "Description",
						Status:      "Todo",
						CreatedAt:   time.Now().Format(time.RFC3339Nano),
						UpdatedAt:   time.Now().Format(time.RFC3339Nano),
					},
				},
			},
			result: []error{nil},
		},
		"delete non-existing task when there are no tasks in the map": {
			inputIDs: []string{"task1"},
			storage:  NewMemoryTaskRepository(),
			result:   []error{models.ErrTaskNotFound},
		},

		"delete non-existing task when there are some tasks in the map": {
			inputIDs: []string{"task1"},
			storage: &MemoryTaskRepository{
				store: map[string]models.Task{
					"task2": {
						ID:          "task2",
						Title:       "Title",
						Description: "Description",
						Status:      "Todo",
						CreatedAt:   time.Now().Format(time.RFC3339Nano),
						UpdatedAt:   time.Now().Format(time.RFC3339Nano),
					},
					"task3": {
						ID:          "task3",
						Title:       "Title",
						Description: "Description",
						Status:      "Todo",
						CreatedAt:   time.Now().Format(time.RFC3339Nano),
						UpdatedAt:   time.Now().Format(time.RFC3339Nano),
					},
				},
			},
			result: []error{models.ErrTaskNotFound},
		},

		"delete multiple tasks": {
			inputIDs: []string{"task1", "task2"},
			storage: &MemoryTaskRepository{
				store: map[string]models.Task{
					"task1": {
						ID:          "task1",
						Title:       "Title",
						Description: "Description",
						Status:      "Todo",
						CreatedAt:   time.Now().Format(time.RFC3339Nano),
						UpdatedAt:   time.Now().Format(time.RFC3339Nano),
					},
					"task2": {
						ID:          "task2",
						Title:       "Title",
						Description: "Description",
						Status:      "Todo",
						CreatedAt:   time.Now().Format(time.RFC3339Nano),
						UpdatedAt:   time.Now().Format(time.RFC3339Nano),
					},
				},
			},
			result: []error{nil, nil},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			for i, id := range test.inputIDs {
				err := test.storage.Delete(id)
				result := test.result[i]

				if !errors.Is(err, result) {
					t.Fatalf("test-case: (%q); returned [%q %q]", name, err, result)
				}

				if err != nil {
					return
				}

				_, err = test.storage.Get(id)
				if !errors.Is(err, models.ErrTaskNotFound) {
					t.Fatalf("test-case: (%q); task hasn't been deleted, get method return: %q", name, err)
				}

				fmt.Println("Task was deleted successfully")
			}
		})
	}
}

func TestStorage_GetAll(t *testing.T) {
	fixedTime := time.Now().Format(time.RFC3339Nano)

	tests := map[string]struct {
		storage *MemoryTaskRepository
		result  TestResultGetAll
	}{
		"get all tasks when storage has multiple tasks": {
			storage: &MemoryTaskRepository{
				store: map[string]models.Task{
					"task1": {
						ID:          "task1",
						Title:       "Title",
						Description: "Description",
						Status:      "Todo",
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					},
					"task2": {
						ID:          "task2",
						Title:       "Title",
						Description: "Description",
						Status:      "Todo",
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					},
					"task3": {
						ID:          "task3",
						Title:       "Title",
						Description: "Description",
						Status:      "Todo",
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					},
				},
			},
			result: TestResultGetAll{
				resultTasks: []models.Task{
					{
						ID:          "task1",
						Title:       "Title",
						Description: "Description",
						Status:      "Todo",
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					},
					{
						ID:          "task2",
						Title:       "Title",
						Description: "Description",
						Status:      "Todo",
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					},
					{
						ID:          "task3",
						Title:       "Title",
						Description: "Description",
						Status:      "Todo",
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					},
				},
				resultError: nil,
			},
		},

		"get all tasks when storage is empty": {
			storage: NewMemoryTaskRepository(),
			result: TestResultGetAll{
				resultTasks: []models.Task{},
				resultError: nil,
			},
		},

		"get all tasks when storage has one task": {
			storage: &MemoryTaskRepository{
				store: map[string]models.Task{
					"task1": {
						ID:          "task1",
						Title:       "Title",
						Description: "Description",
						Status:      "Todo",
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					},
				},
			},
			result: TestResultGetAll{
				resultTasks: []models.Task{
					{
						ID:          "task1",
						Title:       "Title",
						Description: "Description",
						Status:      "Todo",
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					},
				},
				resultError: nil,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tasks, err := test.storage.GetAll()
			resultTasks := test.result.resultTasks
			resultError := test.result.resultError

			if !errors.Is(err, resultError) {
				t.Fatalf("test-case: %q; expected %q returned %q", name, err, resultError)
			}

			if len(tasks) != len(resultTasks) {
				t.Fatalf("test-case: %q; expected %d tasks, got %d", name, len(resultTasks), len(tasks))
			}

			for _, expectedTask := range resultTasks {
				actualTask, err := test.storage.Get(expectedTask.ID)
				if err != nil || expectedTask != actualTask {
					t.Fatalf("test-case: %q; expected task %v, got %v", name, expectedTask, actualTask)
				}
			}
		})
	}
}
