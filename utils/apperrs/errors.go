package apperrs

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
	Cause   error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) New() *AppError {
	return &AppError{
		Code:    e.Code,
		Message: e.Message,
		Status:  e.Status,
	}
}

func (e *AppError) WithCause(cause error) *AppError {
	e.Cause = cause
	return e
}

func ErrInternalServer() *AppError {
	return &AppError{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "Internal server error",
		Status:  http.StatusInternalServerError,
	}
}

func ErrBadRequest() *AppError {
	return &AppError{
		Code:    "BAD_REQUEST",
		Message: "Bad request",
		Status:  http.StatusBadRequest,
	}
}

func ErrNotFound() *AppError {
	return &AppError{
		Code:    "NOT_FOUND",
		Message: "Resource not found",
		Status:  http.StatusNotFound,
	}
}

func ErrMissingHeader() *AppError {
	return &AppError{
		Code:    "MISSING_HEADER",
		Message: "Required header is missing",
		Status:  http.StatusBadRequest,
	}
}

