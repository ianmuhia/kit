package httputil

import (
	"fmt"
	"net/http"
)

// HTTPError represents an HTTP error with status code
type HTTPError struct {
	StatusCode int
	Code       string
	Message    string
	Err        error
}

// Error implements the error interface
func (e *HTTPError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *HTTPError) Unwrap() error {
	return e.Err
}

// NewHTTPError creates a new HTTPError
func NewHTTPError(statusCode int, code, message string) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}

// WrapError wraps an error with HTTP context
func WrapError(statusCode int, code, message string, err error) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

// WriteError writes an HTTPError as a JSON response
func WriteError(w http.ResponseWriter, err *HTTPError) error {
	details := ""
	if err.Err != nil {
		details = err.Err.Error()
	}
	return ErrorWithDetails(w, err.StatusCode, err.Code, err.Message, details)
}

// Common error constructors

// ErrBadRequest creates a 400 error
func ErrBadRequest(message string) *HTTPError {
	return NewHTTPError(http.StatusBadRequest, "BAD_REQUEST", message)
}

// ErrUnauthorized creates a 401 error
func ErrUnauthorized(message string) *HTTPError {
	return NewHTTPError(http.StatusUnauthorized, "UNAUTHORIZED", message)
}

// ErrForbidden creates a 403 error
func ErrForbidden(message string) *HTTPError {
	return NewHTTPError(http.StatusForbidden, "FORBIDDEN", message)
}

// ErrNotFound creates a 404 error
func ErrNotFound(message string) *HTTPError {
	return NewHTTPError(http.StatusNotFound, "NOT_FOUND", message)
}

// ErrConflict creates a 409 error
func ErrConflict(message string) *HTTPError {
	return NewHTTPError(http.StatusConflict, "CONFLICT", message)
}

// ErrInternal creates a 500 error
func ErrInternal(message string) *HTTPError {
	return NewHTTPError(http.StatusInternalServerError, "INTERNAL_ERROR", message)
}
