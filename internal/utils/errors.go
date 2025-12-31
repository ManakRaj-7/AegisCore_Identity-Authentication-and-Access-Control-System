package utils

import (
	"errors"
	"net/http"
)

// AppError represents an application error with HTTP status code
type AppError struct {
	Message    string
	StatusCode int
}

func (e *AppError) Error() string {
	return e.Message
}

// Predefined application errors
var (
	ErrInvalidRequest   = &AppError{Message: "invalid request", StatusCode: http.StatusBadRequest}
	ErrUnauthorized     = &AppError{Message: "unauthorized", StatusCode: http.StatusUnauthorized}
	ErrForbidden        = &AppError{Message: "forbidden", StatusCode: http.StatusForbidden}
	ErrConflict         = &AppError{Message: "email already exists", StatusCode: http.StatusConflict}
	ErrInvalidCredentials = &AppError{Message: "invalid credentials", StatusCode: http.StatusUnauthorized}
	ErrInvalidToken     = &AppError{Message: "invalid or expired token", StatusCode: http.StatusUnauthorized}
	ErrInternalError    = &AppError{Message: "internal server error", StatusCode: http.StatusInternalServerError}
)

// ToAppError converts a standard error to AppError
func ToAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	// Map common error messages to AppError
	errMsg := err.Error()
	switch errMsg {
	case "email already exists":
		return ErrConflict
	case "invalid credentials":
		return ErrInvalidCredentials
	case "invalid or expired token":
		return ErrInvalidToken
	case "invalid email format", "password must be at least 8 characters long":
		return ErrInvalidRequest
	default:
		return ErrInternalError
	}
}

