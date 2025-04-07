package models

import (
	"errors"
)

var (
	// Storage errors.
	ErrTaskExists   = errors.New("task already exists")
	ErrTaskNotFound = errors.New("task not found")
	ErrIDIsEmpty    = errors.New("id is empty")

	// HTTP server errors.
	ErrMethodNotAllowed = errors.New("method not allowed")
	ErrBadRequest       = errors.New("invalid request body")
	ErrEncodeJSON       = errors.New("failed to encode task to JSON")
)
