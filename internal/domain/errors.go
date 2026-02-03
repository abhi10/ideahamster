// Package domain contains core business logic and domain types.
// This package should have no external dependencies beyond the standard library.
package domain

import (
	"errors"
	"fmt"
)

// ===========================================
// Error Types
// ===========================================

// ErrorType distinguishes between error categories
type ErrorType int

const (
	// ErrorTypeInternal represents internal/system errors (500)
	ErrorTypeInternal ErrorType = iota
	// ErrorTypeValidation represents input validation errors (400)
	ErrorTypeValidation
	// ErrorTypeNotFound represents resource not found errors (404)
	ErrorTypeNotFound
	// ErrorTypeUnauthorized represents auth errors (401)
	ErrorTypeUnauthorized
	// ErrorTypeForbidden represents permission errors (403)
	ErrorTypeForbidden
	// ErrorTypeConflict represents conflict errors (409)
	ErrorTypeConflict
	// ErrorTypeRateLimit represents rate limiting errors (429)
	ErrorTypeRateLimit
)

// AppError is a custom error type that distinguishes between
// internal errors (for logging) and user-facing errors (for responses).
type AppError struct {
	// Type categorizes the error for HTTP status code mapping
	Type ErrorType
	// Message is the user-safe message to return in responses
	Message string
	// Internal is the detailed error for logging (may contain sensitive info)
	Internal error
	// Code is an optional application-specific error code
	Code string
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Internal)
	}
	return e.Message
}

// Unwrap returns the underlying error for errors.Is/As support
func (e *AppError) Unwrap() error {
	return e.Internal
}

// UserMessage returns the safe message for API responses
func (e *AppError) UserMessage() string {
	return e.Message
}

// HTTPStatus returns the appropriate HTTP status code
func (e *AppError) HTTPStatus() int {
	switch e.Type {
	case ErrorTypeValidation:
		return 400
	case ErrorTypeUnauthorized:
		return 401
	case ErrorTypeForbidden:
		return 403
	case ErrorTypeNotFound:
		return 404
	case ErrorTypeConflict:
		return 409
	case ErrorTypeRateLimit:
		return 429
	default:
		return 500
	}
}

// ===========================================
// Error Constructors
// ===========================================

// NewInternalError creates an internal server error
func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Type:     ErrorTypeInternal,
		Message:  message,
		Internal: err,
	}
}

// NewValidationError creates a validation error
func NewValidationError(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeValidation,
		Message: message,
	}
}

// NewNotFoundError creates a not found error
func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Type:    ErrorTypeNotFound,
		Message: fmt.Sprintf("%s not found", resource),
	}
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeUnauthorized,
		Message: message,
	}
}

// NewForbiddenError creates a forbidden error
func NewForbiddenError(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeForbidden,
		Message: message,
	}
}

// NewConflictError creates a conflict error
func NewConflictError(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeConflict,
		Message: message,
	}
}

// NewRateLimitError creates a rate limit error
func NewRateLimitError() *AppError {
	return &AppError{
		Type:    ErrorTypeRateLimit,
		Message: "Too many requests, please try again later",
	}
}

// ===========================================
// Error Helpers
// ===========================================

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// GetAppError extracts AppError from an error chain
func GetAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	// Wrap unknown errors as internal errors
	return NewInternalError("An unexpected error occurred", err)
}

// Wrap wraps an error with additional context
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	if appErr, ok := err.(*AppError); ok {
		return &AppError{
			Type:     appErr.Type,
			Message:  message,
			Internal: appErr,
			Code:     appErr.Code,
		}
	}
	return &AppError{
		Type:     ErrorTypeInternal,
		Message:  message,
		Internal: err,
	}
}
