// Package httputil provides HTTP utilities and helpers
package httputil

import (
	"encoding/json"
	"net/http"
)

// Response represents a standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Error      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Error represents an error in the API response
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Meta represents pagination or additional metadata
type Meta struct {
	Page       int `json:"page,omitempty"`
	PageSize   int `json:"page_size,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
	TotalItems int `json:"total_items,omitempty"`
}

// JSON writes a JSON response to the writer
func JSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

// Success writes a successful JSON response
func Success(w http.ResponseWriter, data interface{}) error {
	return JSON(w, http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// SuccessWithMeta writes a successful JSON response with metadata
func SuccessWithMeta(w http.ResponseWriter, data interface{}, meta *Meta) error {
	return JSON(w, http.StatusOK, Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// ErrorResponse writes an error JSON response
func ErrorResponse(w http.ResponseWriter, statusCode int, code, message string) error {
	return JSON(w, statusCode, Response{
		Success: false,
		Error: &Error{
			Code:    code,
			Message: message,
		},
	})
}

// ErrorWithDetails writes an error JSON response with details
func ErrorWithDetails(w http.ResponseWriter, statusCode int, code, message, details string) error {
	return JSON(w, statusCode, Response{
		Success: false,
		Error: &Error{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// BadRequest writes a 400 Bad Request response
func BadRequest(w http.ResponseWriter, message string) error {
	return ErrorResponse(w, http.StatusBadRequest, "BAD_REQUEST", message)
}

// Unauthorized writes a 401 Unauthorized response
func Unauthorized(w http.ResponseWriter, message string) error {
	return ErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", message)
}

// Forbidden writes a 403 Forbidden response
func Forbidden(w http.ResponseWriter, message string) error {
	return ErrorResponse(w, http.StatusForbidden, "FORBIDDEN", message)
}

// NotFound writes a 404 Not Found response
func NotFound(w http.ResponseWriter, message string) error {
	return ErrorResponse(w, http.StatusNotFound, "NOT_FOUND", message)
}

// InternalServerError writes a 500 Internal Server Error response
func InternalServerError(w http.ResponseWriter, message string) error {
	return ErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", message)
}
