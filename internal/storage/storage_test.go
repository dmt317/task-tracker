package storage

import (
	"errors"
	"fmt"
	"task-tracker/internal/models"
	"testing"
	"time"
)

type TestResult struct {
	resultTasks  []models.Task
	resultErrors []error
}

func TestStorage_Add(t *testing.T) {
	tests := map[string]struct {
		inputTasks []models.Task
		initMap    Storage
		result     []error
	}{
		"add valid task": {
			inputTasks: []models.Task{
				{Id: "task1", Description: "Valid task"},
			},
			initMap: Storage{store: make(map[string]models.Task)},
			result:  []error{nil},
		},

		"add task with empty id": {
			inputTasks: []models.Task{
				{Id: "", Description: "No ID"},
			},
			initMap: Storage{store: make(map[string]models.Task)},
			result:  []error{models.ErrIdIsEmpty},
		},

		"add task with duplicate id": {
			inputTasks: []models.Task{
				{Id: "task1", Description: "First task"},
				{Id: "task1", Description: "Duplicate task"},
			},
			initMap: Storage{store: map[string]models.Task{}},
			result:  []error{nil, models.ErrTaskExists},
		},

		"add task with empty description": {
			inputTasks: []models.Task{
				{Id: "task2", Description: ""},
			},
			initMap: Storage{store: make(map[string]models.Task)},
			result:  []error{nil},
		},

		"add multiple valid tasks": {
			inputTasks: []models.Task{
				{Id: "task3", Description: "Task 3"},
				{Id: "task4", Description: "Task 4"},
				{Id: "task5", Description: "Task 5"},
			},
			initMap: Storage{store: make(map[string]models.Task)},
			result:  []error{nil, nil, nil},
		},

		"add duplicate task when task already in storage": {
			inputTasks: []models.Task{
				{Id: "task1", Description: "Duplicate task"},
			},
			initMap: Storage{store: map[string]models.Task{
				"task1": {Id: "task1", Description: "First task", CreatedAt: time.Now().Format(time.RFC3339Nano)}}},
			result: []error{models.ErrTaskExists},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			for i, task := range test.inputTasks {
				gotErr := test.initMap.Add(task)
				expectedErr := test.result[i]
				if !errors.Is(gotErr, expectedErr) {
					t.Fatalf("test-case: (%q); returned %q; expected %q", name, gotErr, expectedErr)
				}

				if gotErr != nil {
					return
				}

				currTask, err := test.initMap.Get(task.Id)
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
		initMap  Storage
		result   TestResult
	}{
		"get existing task": {
			inputIds: []string{"task1"},
			initMap: Storage{store: map[string]models.Task{
				"task1": {Id: "task1", Description: "First task", CreatedAt: fixedTime},
			}},
			result: TestResult{
				resultTasks:  []models.Task{{Id: "task1", Description: "First task", CreatedAt: fixedTime}},
				resultErrors: []error{nil},
			},
		},

		"get non-existing task when there are multiple tasks in the map": {
			inputIds: []string{"task2"},
			initMap: Storage{store: map[string]models.Task{
				"task1": {Id: "task1", Description: "First task", CreatedAt: fixedTime},
				"task3": {Id: "task3", Description: "Third task", CreatedAt: fixedTime},
			}},
			result: TestResult{
				resultTasks:  []models.Task{{}},
				resultErrors: []error{models.ErrTaskNotFound},
			},
		},

		"get non-existing task when there are no tasks in the map": {
			inputIds: []string{"task1"},
			initMap:  Storage{store: map[string]models.Task{}},
			result: TestResult{
				resultTasks:  []models.Task{{}},
				resultErrors: []error{models.ErrTaskNotFound},
			},
		},

		"get task with empty id": {
			inputIds: []string{""},
			initMap: Storage{store: map[string]models.Task{
				"task1": {Id: "task1", Description: "First task", CreatedAt: fixedTime},
				"task2": {Id: "task2", Description: "Second task", CreatedAt: fixedTime},
			}},
			result: TestResult{
				resultTasks:  []models.Task{{}},
				resultErrors: []error{models.ErrIdIsEmpty},
			},
		},

		"get multiple tasks with duplicates": {
			inputIds: []string{"task1", "task2", "task3", "task1", "task2", "task3"},
			initMap: Storage{store: map[string]models.Task{
				"task1": {Id: "task1", Description: "First task", CreatedAt: fixedTime},
				"task2": {Id: "task2", Description: "Second task", CreatedAt: fixedTime},
				"task3": {Id: "task3", Description: "Third task", CreatedAt: fixedTime},
			}},
			result: TestResult{
				resultTasks: []models.Task{
					{Id: "task1", Description: "First task", CreatedAt: fixedTime},
					{Id: "task2", Description: "Second task", CreatedAt: fixedTime},
					{Id: "task3", Description: "Third task", CreatedAt: fixedTime},
					{Id: "task1", Description: "First task", CreatedAt: fixedTime},
					{Id: "task2", Description: "Second task", CreatedAt: fixedTime},
					{Id: "task3", Description: "Third task", CreatedAt: fixedTime},
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
		initMap    Storage
		result     []error
	}{
		"update existing task": {
			inputTasks: []models.Task{
				{Id: "task1", Description: "New description"},
			},
			initMap: Storage{store: map[string]models.Task{
				"task1": {Id: "task1", Description: "Description", CreatedAt: time.Now().Format(time.RFC3339Nano)},
			}},
			result: []error{nil},
		},

		"update non-existing task": {
			inputTasks: []models.Task{
				{Id: "task2", Description: "New description"},
			},
			initMap: Storage{store: map[string]models.Task{
				"task1": {Id: "task1", Description: "Description", CreatedAt: time.Now().Format(time.RFC3339Nano)},
			}},
			result: []error{models.ErrTaskNotFound},
		},

		"update task with empty id": {
			inputTasks: []models.Task{
				{Id: "", Description: "New description"},
			},
			initMap: Storage{store: map[string]models.Task{
				"task1": {Id: "task1", Description: "Description", CreatedAt: time.Now().Format(time.RFC3339Nano)},
			}},
			result: []error{models.ErrIdIsEmpty},
		},

		"update multiple tasks": {
			inputTasks: []models.Task{
				{Id: "task1", Description: "New description"},
				{Id: "task2", Description: "New description"},
				{Id: "task3", Description: "New description"},
			},
			initMap: Storage{store: map[string]models.Task{
				"task1": {Id: "task1", Description: "Description", CreatedAt: time.Now().Format(time.RFC3339Nano)},
				"task2": {Id: "task2", Description: "New description", CreatedAt: time.Now().Format(time.RFC3339Nano)},
				"task3": {Id: "task3", Description: "New description", CreatedAt: time.Now().Format(time.RFC3339Nano)},
			}},
			result: []error{nil, nil, nil},
		},

		"update task using the same data": {
			inputTasks: []models.Task{
				{Id: "task1", Description: "Description"},
			},
			initMap: Storage{store: map[string]models.Task{
				"task1": {Id: "task1", Description: "Description", CreatedAt: time.Now().Format(time.RFC3339Nano)},
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

				updatedTask, err := test.initMap.Get(task.Id)
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
		initMap  Storage
		result   []error
	}{
		"delete existing task": {
			inputIds: []string{"task1"},
			initMap: Storage{store: map[string]models.Task{
				"task1": {Id: "task1", Description: "Description", CreatedAt: time.Now().Format(time.RFC3339Nano)},
			}},
			result: []error{nil},
		},
		"delete non-existing task when there are no tasks in the map": {
			inputIds: []string{"task1"},
			initMap:  Storage{store: map[string]models.Task{}},
			result:   []error{models.ErrTaskNotFound},
		},

		"delete non-existing task when there are some tasks in the map": {
			inputIds: []string{"task1"},
			initMap: Storage{store: map[string]models.Task{
				"task2": {Id: "task2", Description: "Description", CreatedAt: time.Now().Format(time.RFC3339Nano)},
				"task3": {Id: "task3", Description: "Description", CreatedAt: time.Now().Format(time.RFC3339Nano)},
			}},
			result: []error{models.ErrTaskNotFound},
		},

		"delete task with empty id": {
			inputIds: []string{""},
			initMap: Storage{store: map[string]models.Task{
				"task1": {Id: "task1", Description: "Description", CreatedAt: time.Now().Format(time.RFC3339Nano)},
			}},
			result: []error{models.ErrIdIsEmpty},
		},

		"delete multiple tasks": {
			inputIds: []string{"task1", "task2"},
			initMap: Storage{store: map[string]models.Task{
				"task1": {Id: "task1", Description: "Description", CreatedAt: time.Now().Format(time.RFC3339Nano)},
				"task2": {Id: "task2", Description: "Description", CreatedAt: time.Now().Format(time.RFC3339Nano)},
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
