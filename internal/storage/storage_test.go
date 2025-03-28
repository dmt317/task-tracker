package storage

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"task-tracker/internal/models"
)

type TestResult_Get struct {
	resultTasks  []models.Task
	resultErrors []error
}

type TestResult_GetAll struct {
	resultTasks []models.Task
	resultError error
}

func TestStorage_Add(t *testing.T) {
	tests := map[string]struct {
		inputTasks []models.Task
		initMap    *Storage
		result     []error
	}{
		"add valid task": {
			inputTasks: []models.Task{
				{ID: "task1", Title: "Title", Description: "Valid task", Status: "Todo"},
			},
			initMap: &Storage{store: make(map[string]models.Task)},
			result:  []error{nil},
		},

		"add task with empty id": {
			inputTasks: []models.Task{
				{ID: "", Title: "No ID", Description: "No ID", Status: "Todo"},
			},
			initMap: &Storage{store: make(map[string]models.Task)},
			result:  []error{models.ErrIDIsEmpty},
		},

		"add task with duplicate id": {
			inputTasks: []models.Task{
				{ID: "task1", Title: "Title", Description: "First task", Status: "Todo"},
				{ID: "task1", Title: "Title", Description: "Duplicate task", Status: "Todo"},
			},
			initMap: &Storage{store: map[string]models.Task{}},
			result:  []error{nil, models.ErrTaskExists},
		},

		"add task with empty description": {
			inputTasks: []models.Task{
				{ID: "task2", Title: "Empty description", Description: "", Status: "Todo"},
			},
			initMap: &Storage{store: make(map[string]models.Task)},
			result:  []error{nil},
		},

		"add multiple valid tasks": {
			inputTasks: []models.Task{
				{ID: "task3", Title: "Task 3", Description: "Task 3", Status: "Todo"},
				{ID: "task4", Title: "Task 4", Description: "Task 4", Status: "Todo"},
				{ID: "task5", Title: "Task 5", Description: "Task 5", Status: "Todo"},
			},
			initMap: &Storage{store: make(map[string]models.Task)},
			result:  []error{nil, nil, nil},
		},

		"add duplicate task when task already in storage": {
			inputTasks: []models.Task{
				{ID: "task1", Title: "Duplicate task", Description: "Duplicate task", Status: "Todo"},
			},
			initMap: &Storage{store: map[string]models.Task{
				"task1": {ID: "task1", Title: "First task", Description: "First task", Status: "Todo", CreatedAt: time.Now().Format(time.RFC3339Nano), UpdatedAt: time.Now().Format(time.RFC3339Nano)}}},
			result: []error{models.ErrTaskExists},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			for i, task := range test.inputTasks {
				gotErr := test.initMap.Add(&task)
				expectedErr := test.result[i]
				if !errors.Is(gotErr, expectedErr) {
					t.Fatalf("test-case: (%q); returned %q; expected %q", name, gotErr, expectedErr)
				}

				if gotErr != nil {
					return
				}

				currTask, err := test.initMap.Get(task.ID)
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
		inputIds []string
		initMap  *Storage
		result   TestResult_Get
	}{
		"get existing task": {
			inputIds: []string{"task1"},
			initMap: &Storage{store: map[string]models.Task{
				"task1": {ID: "task1", Title: "First task", Description: "First task", Status: "Todo", CreatedAt: fixedTime, UpdatedAt: fixedTime},
			}},
			result: TestResult_Get{
				resultTasks:  []models.Task{{ID: "task1", Title: "First task", Description: "First task", Status: "Todo", CreatedAt: fixedTime, UpdatedAt: fixedTime}},
				resultErrors: []error{nil},
			},
		},

		"get non-existing task when there are multiple tasks in the map": {
			inputIds: []string{"task2"},
			initMap: &Storage{store: map[string]models.Task{
				"task1": {ID: "task1", Title: "First task", Description: "First task", Status: "Todo", CreatedAt: fixedTime, UpdatedAt: fixedTime},
				"task3": {ID: "task3", Title: "Third task", Description: "Third task", Status: "Todo", CreatedAt: fixedTime, UpdatedAt: fixedTime},
			}},
			result: TestResult_Get{
				resultTasks:  []models.Task{{}},
				resultErrors: []error{models.ErrTaskNotFound},
			},
		},

		"get non-existing task when there are no tasks in the map": {
			inputIds: []string{"task1"},
			initMap:  &Storage{store: map[string]models.Task{}},
			result: TestResult_Get{
				resultTasks:  []models.Task{{}},
				resultErrors: []error{models.ErrTaskNotFound},
			},
		},

		"get task with empty id": {
			inputIds: []string{""},
			initMap: &Storage{store: map[string]models.Task{
				"task1": {ID: "task1", Title: "First task", Description: "First task", Status: "Todo", CreatedAt: fixedTime, UpdatedAt: fixedTime},
				"task2": {ID: "task2", Title: "Second task", Description: "Second task", Status: "Todo", CreatedAt: fixedTime, UpdatedAt: fixedTime},
			}},
			result: TestResult_Get{
				resultTasks:  []models.Task{{}},
				resultErrors: []error{models.ErrIDIsEmpty},
			},
		},

		"get multiple tasks with duplicates": {
			inputIds: []string{"task1", "task2", "task3", "task1", "task2", "task3"},
			initMap: &Storage{store: map[string]models.Task{
				"task1": {ID: "task1", Title: "First task", Description: "First task", Status: "Todo", CreatedAt: fixedTime, UpdatedAt: fixedTime},
				"task2": {ID: "task2", Title: "Second task", Description: "Second task", Status: "Todo", CreatedAt: fixedTime, UpdatedAt: fixedTime},
				"task3": {ID: "task3", Title: "Third task", Description: "Third task", Status: "Todo", CreatedAt: fixedTime, UpdatedAt: fixedTime},
			}},
			result: TestResult_Get{
				resultTasks: []models.Task{
					{ID: "task1", Title: "First task", Description: "First task", Status: "Todo", CreatedAt: fixedTime, UpdatedAt: fixedTime},
					{ID: "task2", Title: "Second task", Description: "Second task", Status: "Todo", CreatedAt: fixedTime, UpdatedAt: fixedTime},
					{ID: "task3", Title: "Third task", Description: "Third task", Status: "Todo", CreatedAt: fixedTime, UpdatedAt: fixedTime},
					{ID: "task1", Title: "First task", Description: "First task", Status: "Todo", CreatedAt: fixedTime, UpdatedAt: fixedTime},
					{ID: "task2", Title: "Second task", Description: "Second task", Status: "Todo", CreatedAt: fixedTime, UpdatedAt: fixedTime},
					{ID: "task3", Title: "Third task", Description: "Third task", Status: "Todo", CreatedAt: fixedTime, UpdatedAt: fixedTime},
				},
				resultErrors: []error{nil, nil, nil, nil, nil, nil},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			for i, id := range test.inputIds {
				task, err := test.initMap.Get(id)
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
		inputTasks []models.Task
		initMap    *Storage
		result     []error
	}{
		"update existing task": {
			inputTasks: []models.Task{
				{ID: "task1", Title: "New title", Description: "New description", Status: "Done"},
			},
			initMap: &Storage{store: map[string]models.Task{
				"task1": {ID: "task1", Title: "Title", Description: "Description", Status: "Todo", CreatedAt: time.Now().Format(time.RFC3339Nano), UpdatedAt: time.Now().Format(time.RFC3339Nano)},
			}},
			result: []error{nil},
		},

		"update non-existing task": {
			inputTasks: []models.Task{
				{ID: "task2", Title: "New title", Description: "New description", Status: "Done"},
			},
			initMap: &Storage{store: map[string]models.Task{
				"task1": {ID: "task1", Title: "Title", Description: "Description", Status: "Todo", CreatedAt: time.Now().Format(time.RFC3339Nano), UpdatedAt: time.Now().Format(time.RFC3339Nano)},
			}},
			result: []error{models.ErrTaskNotFound},
		},

		"update task with empty id": {
			inputTasks: []models.Task{
				{ID: "", Title: "New title", Description: "New description", Status: "Done"},
			},
			initMap: &Storage{store: map[string]models.Task{
				"task1": {ID: "task1", Title: "Title", Description: "Description", Status: "Todo", CreatedAt: time.Now().Format(time.RFC3339Nano), UpdatedAt: time.Now().Format(time.RFC3339Nano)},
			}},
			result: []error{models.ErrIDIsEmpty},
		},

		"update multiple tasks": {
			inputTasks: []models.Task{
				{ID: "task1", Title: "New title", Description: "New description", Status: "Done"},
				{ID: "task2", Title: "New title", Description: "New description", Status: "Done"},
				{ID: "task3", Title: "New title", Description: "New description", Status: "Done"},
			},
			initMap: &Storage{store: map[string]models.Task{
				"task1": {ID: "task1", Title: "Title", Description: "Description", Status: "Todo", CreatedAt: time.Now().Format(time.RFC3339Nano), UpdatedAt: time.Now().Format(time.RFC3339Nano)},
				"task2": {ID: "task2", Title: "Title", Description: "Description", Status: "Todo", CreatedAt: time.Now().Format(time.RFC3339Nano), UpdatedAt: time.Now().Format(time.RFC3339Nano)},
				"task3": {ID: "task3", Title: "Title", Description: "Description", Status: "Todo", CreatedAt: time.Now().Format(time.RFC3339Nano), UpdatedAt: time.Now().Format(time.RFC3339Nano)},
			}},
			result: []error{nil, nil, nil},
		},

		"update task using the same data": {
			inputTasks: []models.Task{
				{ID: "task1", Title: "Title", Description: "Description", Status: "Todo"},
			},
			initMap: &Storage{store: map[string]models.Task{
				"task1": {ID: "task1", Title: "Title", Description: "Description", Status: "Todo", CreatedAt: time.Now().Format(time.RFC3339Nano), UpdatedAt: time.Now().Format(time.RFC3339Nano)},
			}},
			result: []error{nil},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			for i, task := range test.inputTasks {
				err := test.initMap.Update(task)
				result := test.result[i]
				if !errors.Is(err, result) {
					t.Fatalf("test-case: (%q); returned [%q %q]", name, err, result)
				}

				if err != nil {
					return
				}

				updatedTask, err := test.initMap.Get(task.ID)
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
		inputIds []string
		initMap  *Storage
		result   []error
	}{
		"delete existing task": {
			inputIds: []string{"task1"},
			initMap: &Storage{store: map[string]models.Task{
				"task1": {ID: "task1", Title: "Title", Description: "Description", Status: "Todo", CreatedAt: time.Now().Format(time.RFC3339Nano), UpdatedAt: time.Now().Format(time.RFC3339Nano)},
			}},
			result: []error{nil},
		},
		"delete non-existing task when there are no tasks in the map": {
			inputIds: []string{"task1"},
			initMap:  &Storage{store: map[string]models.Task{}},
			result:   []error{models.ErrTaskNotFound},
		},

		"delete non-existing task when there are some tasks in the map": {
			inputIds: []string{"task1"},
			initMap: &Storage{store: map[string]models.Task{
				"task2": {ID: "task2", Title: "Title", Description: "Description", Status: "Todo", CreatedAt: time.Now().Format(time.RFC3339Nano), UpdatedAt: time.Now().Format(time.RFC3339Nano)},
				"task3": {ID: "task3", Title: "Title", Description: "Description", Status: "Todo", CreatedAt: time.Now().Format(time.RFC3339Nano), UpdatedAt: time.Now().Format(time.RFC3339Nano)},
			}},
			result: []error{models.ErrTaskNotFound},
		},

		"delete task with empty id": {
			inputIds: []string{""},
			initMap: &Storage{store: map[string]models.Task{
				"task1": {ID: "task1", Title: "Title", Description: "Description", Status: "Todo", CreatedAt: time.Now().Format(time.RFC3339Nano), UpdatedAt: time.Now().Format(time.RFC3339Nano)},
			}},
			result: []error{models.ErrIDIsEmpty},
		},

		"delete multiple tasks": {
			inputIds: []string{"task1", "task2"},
			initMap: &Storage{store: map[string]models.Task{
				"task1": {ID: "task1", Title: "Title", Description: "Description", Status: "Todo", CreatedAt: time.Now().Format(time.RFC3339Nano), UpdatedAt: time.Now().Format(time.RFC3339Nano)},
				"task2": {ID: "task2", Title: "Title", Description: "Description", Status: "Todo", CreatedAt: time.Now().Format(time.RFC3339Nano), UpdatedAt: time.Now().Format(time.RFC3339Nano)},
			}},
			result: []error{nil, nil},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			for i, id := range test.inputIds {
				err := test.initMap.Delete(id)
				result := test.result[i]
				if !errors.Is(err, result) {
					t.Fatalf("test-case: (%q); returned [%q %q]", name, err, result)
				}

				if err != nil {
					return
				}

				_, err = test.initMap.Get(id)
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
		initMap *Storage
		result  TestResult_GetAll
	}{
		"get all tasks when storage has multiple tasks": {
			initMap: &Storage{store: map[string]models.Task{
				"task1": {ID: "task1", Title: "Title", Description: "Description", Status: "Todo", CreatedAt: fixedTime, UpdatedAt: fixedTime},
				"task2": {ID: "task2", Title: "Title", Description: "Description", Status: "Todo", CreatedAt: fixedTime, UpdatedAt: fixedTime},
				"task3": {ID: "task3", Title: "Title", Description: "Description", Status: "Todo", CreatedAt: fixedTime, UpdatedAt: fixedTime},
			}},
			result: TestResult_GetAll{
				resultTasks: []models.Task{
					{ID: "task1", Title: "Title", Description: "Description", Status: "Todo", CreatedAt: fixedTime, UpdatedAt: fixedTime},
					{ID: "task2", Title: "Title", Description: "Description", Status: "Todo", CreatedAt: fixedTime, UpdatedAt: fixedTime},
					{ID: "task3", Title: "Title", Description: "Description", Status: "Todo", CreatedAt: fixedTime, UpdatedAt: fixedTime},
				},
				resultError: nil,
			},
		},

		"get all tasks when storage is empty": {
			initMap: &Storage{store: map[string]models.Task{}},
			result: TestResult_GetAll{
				resultTasks: []models.Task{},
				resultError: nil,
			},
		},

		"get all tasks when storage has one task": {
			initMap: &Storage{store: map[string]models.Task{
				"task1": {ID: "task1", Title: "Title", Description: "Description", Status: "Todo", CreatedAt: fixedTime, UpdatedAt: fixedTime},
			}},
			result: TestResult_GetAll{
				resultTasks: []models.Task{
					{ID: "task1", Title: "Title", Description: "Description", Status: "Todo", CreatedAt: fixedTime, UpdatedAt: fixedTime},
				},
				resultError: nil,
			},
		},

		"get all tasks when store map is nil": {
			initMap: &Storage{store: nil},
			result: TestResult_GetAll{
				resultTasks: nil,
				resultError: models.ErrStorageNotInitialized,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tasks, err := test.initMap.GetAll()
			resultTasks := test.result.resultTasks
			resultError := test.result.resultError

			if !errors.Is(err, resultError) {
				t.Fatalf("test-case: %q; expected %q returned %q", name, err, resultError)
			}

			if len(tasks) != len(resultTasks) {
				t.Fatalf("test-case: %q; expected %d tasks, got %d", name, len(resultTasks), len(tasks))
			}

			for _, expectedTask := range resultTasks {
				actualTask, found := test.initMap.store[expectedTask.ID]
				if !found || expectedTask != actualTask {
					t.Fatalf("test-case: %q; expected task %v, got %v", name, expectedTask, actualTask)
				}
			}
		})
	}
}
