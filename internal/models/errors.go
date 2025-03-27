package models

import (
	"errors"
)

var (
	// Storage errors
	ErrTaskExists            = errors.New("task already exists")
	ErrTaskNotFound          = errors.New("task not found")
	ErrIdIsEmpty             = errors.New("id is empty")
	ErrStorageNotInitialized = errors.New("storage is not initialized")

	// HTTP server errors
	ErrMethodNotAllowed = errors.New("method not allowed")
	ErrBadRequest       = errors.New("invalid request body")
	ErrEncodeJSON       = errors.New("failed to encode task to JSON")
)
