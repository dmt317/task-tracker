package storage

import (
	"errors"
	"fmt"
	"task-tracker/models"
	"task-tracker/task"
	"testing"
	"time"
)

type TestResult struct {
	resultTasks  []task.Task
	resultErrors []error
}

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
			result:  []error{models.ErrIdIsEmpty},
		},

		"add task with duplicate id": {
			inputTasks: []task.Task{
				{Id: "task1", Description: "First task", CreatedAt: time.Now().Format(time.RFC3339Nano)},
				{Id: "task1", Description: "Duplicate task", CreatedAt: time.Now().Format(time.RFC3339Nano)},
			},
			initMap: Storage{store: map[string]task.Task{}},
			result:  []error{nil, models.ErrTaskExists},
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
				if gotErr == nil {
					currTask, err := test.initMap.Get(task.Id)
					if err == nil && currTask == task {
						fmt.Println(currTask)
					}
				}
			}
		})
	}
}

func TestStorage_Get(t *testing.T) {
	fixedTime := time.Now().Add(-time.Hour).Format(time.RFC3339Nano)

	tests := map[string]struct {
		inputIds []string
		initMap  Storage
		result   TestResult
	}{
		"get existing task": {
			inputIds: []string{"task1"},
			initMap: Storage{store: map[string]task.Task{
				"task1": {Id: "task1", Description: "First task", CreatedAt: fixedTime},
			}},
			result: TestResult{
				resultTasks:  []task.Task{{Id: "task1", Description: "First task", CreatedAt: fixedTime}},
				resultErrors: []error{nil},
			},
		},

		"get non-existing task when there are multiple tasks in the map": {
			inputIds: []string{"task2"},
			initMap: Storage{store: map[string]task.Task{
				"task1": {Id: "task1", Description: "First task", CreatedAt: fixedTime},
				"task3": {Id: "task3", Description: "Third task", CreatedAt: fixedTime},
			}},
			result: TestResult{
				resultTasks:  []task.Task{{}},
				resultErrors: []error{models.ErrTaskNotFound},
			},
		},

		"get non-existing task when there are no tasks in the map": {
			inputIds: []string{"task1"},
			initMap:  Storage{store: map[string]task.Task{}},
			result: TestResult{
				resultTasks:  []task.Task{{}},
				resultErrors: []error{models.ErrTaskNotFound},
			},
		},

		"get task with empty id": {
			inputIds: []string{""},
			initMap: Storage{store: map[string]task.Task{
				"task1": {Id: "task1", Description: "First task", CreatedAt: fixedTime},
				"task2": {Id: "task2", Description: "Second task", CreatedAt: fixedTime},
			}},
			result: TestResult{
				resultTasks:  []task.Task{{}},
				resultErrors: []error{models.ErrIdIsEmpty},
			},
		},

		"ensuring thread-safe get operation with duplicate task requests": {
			inputIds: []string{"task1", "task2", "task3", "task1", "task2", "task3"},
			initMap: Storage{store: map[string]task.Task{
				"task1": {Id: "task1", Description: "First task", CreatedAt: fixedTime},
				"task2": {Id: "task2", Description: "Second task", CreatedAt: fixedTime},
				"task3": {Id: "task3", Description: "Third task", CreatedAt: fixedTime},
			}},
			result: TestResult{
				resultTasks: []task.Task{
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
				if err == nil && task == test.initMap.store[id] {
					fmt.Println(task)
				}
			}
		})
	}
}

func TestStorage_Update(t *testing.T) {

}

func TestStorage_Delete(t *testing.T) {

}
