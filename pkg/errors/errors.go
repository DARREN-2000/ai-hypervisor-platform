package errors

import "fmt"

// ErrorCode represents a specific error classification
type ErrorCode string

const (
	// VM errors
	ErrVMNotFound         ErrorCode = "vm_not_found"
	ErrVMAlreadyExists    ErrorCode = "vm_already_exists"
	ErrVMInvalidState     ErrorCode = "vm_invalid_state"
	ErrVMOperationFailed  ErrorCode = "vm_operation_failed"

	// GPU errors
	ErrGPUNotFound        ErrorCode = "gpu_not_found"
	ErrGPUUnavailable     ErrorCode = "gpu_unavailable"
	ErrGPUAllocationFailed ErrorCode = "gpu_allocation_failed"
	ErrInsufficientGPUs   ErrorCode = "insufficient_gpus"

	// Scheduler errors
	ErrNoCapacity         ErrorCode = "no_capacity"
	ErrSchedulingFailed   ErrorCode = "scheduling_failed"
	ErrNodeNotFound       ErrorCode = "node_not_found"

	// Task errors
	ErrTaskNotFound       ErrorCode = "task_not_found"
	ErrTaskFailed         ErrorCode = "task_failed"
	ErrTaskMaxRetries     ErrorCode = "task_max_retries"

	// Configuration errors
	ErrConfigNotFound     ErrorCode = "config_not_found"
	ErrInvalidConfig      ErrorCode = "invalid_config"

	// Generic errors
	ErrInternal           ErrorCode = "internal_error"
	ErrValidation         ErrorCode = "validation_error"
	ErrUnauthorized       ErrorCode = "unauthorized"
	ErrForbidden          ErrorCode = "forbidden"
	ErrTimeout            ErrorCode = "timeout"
	ErrConflict           ErrorCode = "conflict"
)

// APIError represents an API-level error with context
type APIError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	StatusCode int       `json:"status_code"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Cause      error     `json:"-"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewAPIError creates a new API error
func NewAPIError(code ErrorCode, message string, statusCode int) *APIError {
	return &APIError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Details:    make(map[string]interface{}),
	}
}

// WithCause adds the underlying cause to the error
func (e *APIError) WithCause(err error) *APIError {
	e.Cause = err
	return e
}

// WithDetail adds detail information
func (e *APIError) WithDetail(key string, value interface{}) *APIError {
	e.Details[key] = value
	return e
}

// Common error constructors

// NotFound creates a not found error
func NotFound(resource string) *APIError {
	return NewAPIError(ErrVMNotFound, fmt.Sprintf("%s not found", resource), 404)
}

// ValidationError creates a validation error
func ValidationError(message string) *APIError {
	return NewAPIError(ErrValidation, message, 400)
}

// ConflictError creates a conflict error
func ConflictError(message string) *APIError {
	return NewAPIError(ErrConflict, message, 409)
}

// InternalError creates an internal error
func InternalError(message string) *APIError {
	return NewAPIError(ErrInternal, message, 500)
}

// UnauthorizedError creates an unauthorized error
func UnauthorizedError(message string) *APIError {
	return NewAPIError(ErrUnauthorized, message, 401)
}

// ForbiddenError creates a forbidden error
func ForbiddenError(message string) *APIError {
	return NewAPIError(ErrForbidden, message, 403)
}

// TimeoutError creates a timeout error
func TimeoutError(operation string) *APIError {
	return NewAPIError(ErrTimeout, fmt.Sprintf("operation '%s' timed out", operation), 504)
}

// InsufficientResourcesError creates a resource exhaustion error
func InsufficientResourcesError(resource string) *APIError {
	return NewAPIError(ErrNoCapacity, fmt.Sprintf("insufficient %s available", resource), 507)
}

// InvalidStateError creates an invalid state error
func InvalidStateError(resource, currentState, targetState string) *APIError {
	return NewAPIError(
		ErrVMInvalidState,
		fmt.Sprintf("cannot transition %s from %s to %s", resource, currentState, targetState),
		400,
	)
}
