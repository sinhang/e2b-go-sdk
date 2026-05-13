package e2b

import "strings"

// APIErrorPayload is a generic shape commonly returned by E2B-like APIs.
type APIErrorPayload struct {
	Code    string `json:"code,omitempty"`
	Type    string `json:"type,omitempty"`
	Message string `json:"message,omitempty"`
}

// FromAPIError maps API error metadata into a concrete SDK exception.
// exceptionName should match the Python SDK naming, e.g. "TimeoutException".
func FromAPIError(exceptionName string, payload APIErrorPayload, cause error) error {
	errType := ErrorType(payload.Type)
	name := strings.TrimSpace(exceptionName)
	msg := strings.TrimSpace(payload.Message)
	if msg == "" {
		msg = name
	}

	switch name {
	case "SandboxException":
		return NewSandboxException(msg, errType, cause)
	case "TimeoutException":
		return NewTimeoutException(msg, errType, cause)
	case "InvalidArgumentException":
		return NewInvalidArgumentException(msg, errType, cause)
	case "NotEnoughSpaceException":
		return NewNotEnoughSpaceException(msg, errType, cause)
	case "NotFoundException":
		return NewNotFoundException(msg, errType, cause)
	case "AuthenticationException":
		return NewAuthenticationException(msg, errType, cause)
	case "TemplateException":
		return NewTemplateException(msg, errType, cause)
	case "RateLimitException":
		return NewRateLimitException(msg, errType, cause)
	case "BuildException":
		return NewBuildException(msg, errType, cause)
	case "FileUploadException":
		return NewFileUploadException(msg, errType, cause)
	default:
		return &BaseError{Type: errType, Message: msg, Cause: cause}
	}
}
