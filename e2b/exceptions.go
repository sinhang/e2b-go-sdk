package e2b

import "fmt"

// ErrorType is the machine-readable API error type.
type ErrorType string

const (
	ErrorTypeUnavailable      ErrorType = "unavailable"
	ErrorTypeCanceled         ErrorType = "canceled"
	ErrorTypeDeadlineExceeded ErrorType = "deadline_exceeded"
	ErrorTypeUnknown          ErrorType = "unknown"
)

// BaseError holds common fields for all SDK errors.
type BaseError struct {
	Type    ErrorType
	Message string
	Cause   error
}

func (e *BaseError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Type != "" && e.Message != "" {
		return fmt.Sprintf("%s: %s", e.Type, e.Message)
	}
	if e.Message != "" {
		return e.Message
	}
	if e.Type != "" {
		return string(e.Type)
	}
	return "e2b error"
}

func (e *BaseError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// SandboxException is the base class for all sandbox errors.
type SandboxException struct {
	*BaseError
}

func NewSandboxException(msg string, errType ErrorType, cause error) *SandboxException {
	return &SandboxException{BaseError: &BaseError{Type: errType, Message: msg, Cause: cause}}
}

func (e *SandboxException) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.BaseError
}

// TimeoutException is raised when a timeout occurs.
type TimeoutException struct {
	*SandboxException
}

func NewTimeoutException(msg string, errType ErrorType, cause error) *TimeoutException {
	return &TimeoutException{SandboxException: NewSandboxException(msg, errType, cause)}
}

func (e *TimeoutException) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.SandboxException
}

// InvalidArgumentException is raised when an invalid argument is provided.
type InvalidArgumentException struct {
	*SandboxException
}

func NewInvalidArgumentException(msg string, errType ErrorType, cause error) *InvalidArgumentException {
	return &InvalidArgumentException{SandboxException: NewSandboxException(msg, errType, cause)}
}

func (e *InvalidArgumentException) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.SandboxException
}

// NotEnoughSpaceException is raised when there is not enough disk space.
type NotEnoughSpaceException struct {
	*SandboxException
}

func NewNotEnoughSpaceException(msg string, errType ErrorType, cause error) *NotEnoughSpaceException {
	return &NotEnoughSpaceException{SandboxException: NewSandboxException(msg, errType, cause)}
}

func (e *NotEnoughSpaceException) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.SandboxException
}

// NotFoundException is raised when a resource is not found.
type NotFoundException struct {
	*SandboxException
}

func NewNotFoundException(msg string, errType ErrorType, cause error) *NotFoundException {
	return &NotFoundException{SandboxException: NewSandboxException(msg, errType, cause)}
}

func (e *NotFoundException) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.SandboxException
}

// AuthenticationException is raised when authentication fails.
type AuthenticationException struct {
	*BaseError
}

func NewAuthenticationException(msg string, errType ErrorType, cause error) *AuthenticationException {
	return &AuthenticationException{BaseError: &BaseError{Type: errType, Message: msg, Cause: cause}}
}

// TemplateException is raised when template is incompatible with the SDK.
type TemplateException struct {
	*SandboxException
}

func NewTemplateException(msg string, errType ErrorType, cause error) *TemplateException {
	return &TemplateException{SandboxException: NewSandboxException(msg, errType, cause)}
}

func (e *TemplateException) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.SandboxException
}

// RateLimitException is raised when the API rate limit is exceeded.
type RateLimitException struct {
	*SandboxException
}

func NewRateLimitException(msg string, errType ErrorType, cause error) *RateLimitException {
	return &RateLimitException{SandboxException: NewSandboxException(msg, errType, cause)}
}

func (e *RateLimitException) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.SandboxException
}

// BuildException is raised when a build fails.
type BuildException struct {
	*BaseError
}

func NewBuildException(msg string, errType ErrorType, cause error) *BuildException {
	return &BuildException{BaseError: &BaseError{Type: errType, Message: msg, Cause: cause}}
}

// FileUploadException is raised when file upload fails.
type FileUploadException struct {
	*BuildException
}

func NewFileUploadException(msg string, errType ErrorType, cause error) *FileUploadException {
	return &FileUploadException{BuildException: NewBuildException(msg, errType, cause)}
}

func (e *FileUploadException) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.BuildException
}
