package errors

import (
	"errors"
)

var (
	ErrTaskExists   = errors.New("task already exists")
	ErrTaskNotFound = errors.New("task not found")
)
