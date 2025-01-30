package models

import (
	"errors"
)

var (
	ErrTaskExists   = errors.New("task already exists")
	ErrTaskNotFound = errors.New("task not found")
	ErrIdIsEmpty    = errors.New("id is empty")
)
