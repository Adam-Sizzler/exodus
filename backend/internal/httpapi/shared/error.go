package shared

import (
	"fmt"
	"net/http"

	"exodus/internal/config"
)

// APIError represents a strongly typed HTTP error implementing the error interface.
type APIError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	ErrorText  string `json:"error"`
	Code       string `json:"code,omitempty"`
	Details    any    `json:"details,omitempty"`
	Cause      error  `json:"-"`
}

func (e *APIError) Error() string {
	if e == nil {
		return ""
	}
	if e.Cause != nil {
		return fmt.Sprintf("[%d] %s: %v", e.StatusCode, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%d] %s", e.StatusCode, e.Message)
}

func (e *APIError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// WithCause returns a copy of e with Cause set to cause, leaving e itself
// (a shared, package-level registry entry) untouched. Registry entries are
// declared once as package-level *APIError values and read concurrently by
// every request — mutating e.Cause directly would race across requests, so
// every call site must go through WithCause instead of assigning to a
// registry entry's Cause field.
func (e *APIError) WithCause(cause error) *APIError {
	if e == nil {
		return nil
	}
	cp := *e
	cp.Cause = cause
	return &cp
}

func (e *APIError) ToResponse() ErrorResponse {
	if e == nil {
		return ErrorResponse{}
	}
	errText := e.ErrorText
	if errText == "" {
		errText = e.Message
	}
	return ErrorResponse{
		StatusCode: e.StatusCode,
		Message:    e.Message,
		Error:      errText,
		Code:       e.Code,
		Details:    e.Details,
	}
}

// Factory functions for common HTTP errors

func NewBadRequestError(message string, err error) *APIError {
	return &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
		ErrorText:  message,
		Cause:      err,
	}
}

func NewValidationError(message string, details any) *APIError {
	return &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
		ErrorText:  message,
		Details:    details,
	}
}

func NewNotFoundError(resource string, err error) *APIError {
	msg := fmt.Sprintf("%s not found", resource)
	return &APIError{
		StatusCode: http.StatusNotFound,
		Message:    msg,
		ErrorText:  msg,
		Cause:      err,
	}
}

func NewUnauthorizedError(message string, err error) *APIError {
	if message == "" {
		message = "Unauthorized"
	}
	return &APIError{
		StatusCode: http.StatusUnauthorized,
		Message:    message,
		ErrorText:  message,
		Cause:      err,
	}
}

func NewForbiddenError(message string, err error) *APIError {
	if message == "" {
		message = "Forbidden"
	}
	return &APIError{
		StatusCode: http.StatusForbidden,
		Message:    message,
		ErrorText:  message,
		Cause:      err,
	}
}

func NewInternalError(message string, err error) *APIError {
	if message == "" {
		message = "Internal Server Error"
	}
	return &APIError{
		StatusCode: http.StatusInternalServerError,
		Message:    message,
		ErrorText:  "Internal Server Error",
		Cause:      err,
	}
}

func NewAPIError(statusCode int, code, message string, err error) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		ErrorText:  message,
		Code:       code,
		Cause:      err,
	}
}

// SendAPIError sends a typed APIError as JSON response and logs it appropriately.
func SendAPIError(w http.ResponseWriter, apiErr *APIError, cfg *config.BackendConfig) {
	if apiErr == nil {
		return
	}
	if apiErr.Cause != nil && cfg != nil && cfg.Logger != nil {
		if apiErr.StatusCode >= 500 {
			cfg.Logger.Error(apiErr.Message, "error", apiErr.Cause, "code", apiErr.Code)
		} else {
			cfg.Logger.Debug(apiErr.Message, "error", apiErr.Cause, "code", apiErr.Code)
		}
	}
	WriteJSON(w, apiErr.StatusCode, apiErr.ToResponse())
}

// SendError is the general helper for sending error responses with status code and message.
func SendError(w http.ResponseWriter, code int, msg string, err error, cfg *config.BackendConfig) {
	SendAPIError(w, &APIError{
		StatusCode: code,
		Message:    msg,
		ErrorText:  msg,
		Cause:      err,
	}, cfg)
}
