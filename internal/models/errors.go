package models

import (
	"errors"
	"net/http"
)

type Error struct {
	Err        error
	StatusCode int
}

func (e Error) Error() string {
	return e.Err.Error()
}

func NewError(errorMessage string, statusCode int) Error {
	return Error{
		Err:        errors.New(errorMessage),
		StatusCode: statusCode,
	}
}

var (
	// Storage errors.
	ErrTaskExists   = NewError("task already exists", http.StatusConflict)
	ErrTaskNotFound = NewError("task not found", http.StatusNotFound)
	ErrIDIsEmpty    = NewError("id is empty", http.StatusBadRequest)

	ErrTitleIsEmpty       = NewError("title field is empty", http.StatusBadRequest)
	ErrDescriptionIsEmpty = NewError("description field is empty", http.StatusBadRequest)
	ErrStatusIsEmpty      = NewError("status field is empty", http.StatusBadRequest)

	ErrMethodNotAllowed  = NewError("method not allowed", http.StatusBadRequest)
	ErrBadRequest        = NewError("invalid request body", http.StatusBadRequest)
	ErrSwaggerUINotFound = NewError("swagger UI not found", http.StatusNotFound)
)
