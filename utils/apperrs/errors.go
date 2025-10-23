package apperrs

import "errors"

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrNotFound     = errors.New("not found")
	ErrBadRequest   = errors.New("bad request")
	ErrInternal     = errors.New("internal server error")
)

// AppError represents an application error
type AppError struct {
	Code    string
	Message string
	Cause   error
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func (e *AppError) New() *AppError {
	return &AppError{
		Code:    e.Code,
		Message: e.Message,
	}
}

func (e *AppError) WithCause(cause error) *AppError {
	e.Cause = cause
	return e
}
