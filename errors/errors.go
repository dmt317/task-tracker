package errors

import (
	"errors"
)

var (
	ErrTaskExists   = errors.New("task already exists")
	ErrTaskNotFound = errors.New("task not found")
	ErrIdIsEmpty    = errors.New("id is empty")
)

func Is(err1 error, err2 error) bool {
	return errors.Is(err1, err2)
}
