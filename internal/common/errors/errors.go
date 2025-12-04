package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application error
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	Details    any    `json:"details,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewAppError creates a new application error
func NewAppError(code, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details any) *AppError {
	e.Details = details
	return e
}

// Common error codes
const (
	ErrCodeBadRequest          = "BAD_REQUEST"
	ErrCodeUnauthorized        = "UNAUTHORIZED"
	ErrCodeForbidden           = "FORBIDDEN"
	ErrCodeNotFound            = "NOT_FOUND"
	ErrCodeConflict            = "CONFLICT"
	ErrCodeValidation          = "VALIDATION_ERROR"
	ErrCodeInternal            = "INTERNAL_ERROR"
	ErrCodeDatabaseError       = "DATABASE_ERROR"
	ErrCodeInvalidCredentials  = "INVALID_CREDENTIALS"
	ErrCodeTokenExpired        = "TOKEN_EXPIRED"
	ErrCodeTokenInvalid        = "TOKEN_INVALID"
	ErrCodeMFARequired         = "MFA_REQUIRED"
	ErrCodeMFAInvalid          = "MFA_INVALID"
	ErrCodeRateLimitExceeded   = "RATE_LIMIT_EXCEEDED"
	ErrCodeServiceUnavailable  = "SERVICE_UNAVAILABLE"
)

// Predefined errors
var (
	ErrBadRequest = NewAppError(
		ErrCodeBadRequest,
		"Bad request",
		http.StatusBadRequest,
	)

	ErrUnauthorized = NewAppError(
		ErrCodeUnauthorized,
		"Unauthorized access",
		http.StatusUnauthorized,
	)

	ErrForbidden = NewAppError(
		ErrCodeForbidden,
		"Access forbidden",
		http.StatusForbidden,
	)

	ErrNotFound = NewAppError(
		ErrCodeNotFound,
		"Resource not found",
		http.StatusNotFound,
	)

	ErrConflict = NewAppError(
		ErrCodeConflict,
		"Resource conflict",
		http.StatusConflict,
	)

	ErrValidation = NewAppError(
		ErrCodeValidation,
		"Validation error",
		http.StatusBadRequest,
	)

	ErrInternal = NewAppError(
		ErrCodeInternal,
		"Internal server error",
		http.StatusInternalServerError,
	)

	ErrDatabaseError = NewAppError(
		ErrCodeDatabaseError,
		"Database error",
		http.StatusInternalServerError,
	)

	ErrInvalidCredentials = NewAppError(
		ErrCodeInvalidCredentials,
		"Invalid credentials",
		http.StatusUnauthorized,
	)

	ErrTokenExpired = NewAppError(
		ErrCodeTokenExpired,
		"Token has expired",
		http.StatusUnauthorized,
	)

	ErrTokenInvalid = NewAppError(
		ErrCodeTokenInvalid,
		"Invalid token",
		http.StatusUnauthorized,
	)

	ErrMFARequired = NewAppError(
		ErrCodeMFARequired,
		"Multi-factor authentication required",
		http.StatusUnauthorized,
	)

	ErrMFAInvalid = NewAppError(
		ErrCodeMFAInvalid,
		"Invalid MFA code",
		http.StatusUnauthorized,
	)

	ErrRateLimitExceeded = NewAppError(
		ErrCodeRateLimitExceeded,
		"Rate limit exceeded",
		http.StatusTooManyRequests,
	)

	ErrServiceUnavailable = NewAppError(
		ErrCodeServiceUnavailable,
		"Service temporarily unavailable",
		http.StatusServiceUnavailable,
	)
)

// Domain-specific errors

// Patient errors
func ErrPatientNotFound(id string) *AppError {
	return NewAppError(
		"PATIENT_NOT_FOUND",
		fmt.Sprintf("Patient with ID %s not found", id),
		http.StatusNotFound,
	)
}

func ErrPatientAlreadyExists(mrn string) *AppError {
	return NewAppError(
		"PATIENT_ALREADY_EXISTS",
		fmt.Sprintf("Patient with MRN %s already exists", mrn),
		http.StatusConflict,
	)
}

// Encounter errors
func ErrEncounterNotFound(id string) *AppError {
	return NewAppError(
		"ENCOUNTER_NOT_FOUND",
		fmt.Sprintf("Encounter with ID %s not found", id),
		http.StatusNotFound,
	)
}

// Appointment errors
func ErrAppointmentNotFound(id string) *AppError {
	return NewAppError(
		"APPOINTMENT_NOT_FOUND",
		fmt.Sprintf("Appointment with ID %s not found", id),
		http.StatusNotFound,
	)
}

func ErrAppointmentConflict() *AppError {
	return NewAppError(
		"APPOINTMENT_CONFLICT",
		"Appointment slot is not available",
		http.StatusConflict,
	)
}

// User errors
func ErrUserNotFound(id string) *AppError {
	return NewAppError(
		"USER_NOT_FOUND",
		fmt.Sprintf("User with ID %s not found", id),
		http.StatusNotFound,
	)
}

func ErrUserAlreadyExists(email string) *AppError {
	return NewAppError(
		"USER_ALREADY_EXISTS",
		fmt.Sprintf("User with email %s already exists", email),
		http.StatusConflict,
	)
}

// Permission errors
func ErrInsufficientPermissions() *AppError {
	return NewAppError(
		"INSUFFICIENT_PERMISSIONS",
		"You do not have permission to perform this action",
		http.StatusForbidden,
	)
}
