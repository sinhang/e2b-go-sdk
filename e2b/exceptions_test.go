package e2b

import (
	"errors"
	"testing"
)

func TestFromAPIError_AllExceptions(t *testing.T) {
	cause := errors.New("root cause")
	payload := APIErrorPayload{Type: "deadline_exceeded", Message: "timed out"}

	tests := []struct {
		name          string
		exceptionName string
		assert        func(t *testing.T, err error)
	}{
		{
			name:          "SandboxException",
			exceptionName: "SandboxException",
			assert: func(t *testing.T, err error) {
				var target *SandboxException
				if !errors.As(err, &target) {
					t.Fatalf("expected SandboxException, got %T", err)
				}
			},
		},
		{
			name:          "TimeoutException",
			exceptionName: "TimeoutException",
			assert: func(t *testing.T, err error) {
				var target *TimeoutException
				if !errors.As(err, &target) {
					t.Fatalf("expected TimeoutException, got %T", err)
				}
			},
		},
		{
			name:          "InvalidArgumentException",
			exceptionName: "InvalidArgumentException",
			assert: func(t *testing.T, err error) {
				var target *InvalidArgumentException
				if !errors.As(err, &target) {
					t.Fatalf("expected InvalidArgumentException, got %T", err)
				}
			},
		},
		{
			name:          "NotEnoughSpaceException",
			exceptionName: "NotEnoughSpaceException",
			assert: func(t *testing.T, err error) {
				var target *NotEnoughSpaceException
				if !errors.As(err, &target) {
					t.Fatalf("expected NotEnoughSpaceException, got %T", err)
				}
			},
		},
		{
			name:          "NotFoundException",
			exceptionName: "NotFoundException",
			assert: func(t *testing.T, err error) {
				var target *NotFoundException
				if !errors.As(err, &target) {
					t.Fatalf("expected NotFoundException, got %T", err)
				}
			},
		},
		{
			name:          "AuthenticationException",
			exceptionName: "AuthenticationException",
			assert: func(t *testing.T, err error) {
				var target *AuthenticationException
				if !errors.As(err, &target) {
					t.Fatalf("expected AuthenticationException, got %T", err)
				}
			},
		},
		{
			name:          "TemplateException",
			exceptionName: "TemplateException",
			assert: func(t *testing.T, err error) {
				var target *TemplateException
				if !errors.As(err, &target) {
					t.Fatalf("expected TemplateException, got %T", err)
				}
			},
		},
		{
			name:          "RateLimitException",
			exceptionName: "RateLimitException",
			assert: func(t *testing.T, err error) {
				var target *RateLimitException
				if !errors.As(err, &target) {
					t.Fatalf("expected RateLimitException, got %T", err)
				}
			},
		},
		{
			name:          "BuildException",
			exceptionName: "BuildException",
			assert: func(t *testing.T, err error) {
				var target *BuildException
				if !errors.As(err, &target) {
					t.Fatalf("expected BuildException, got %T", err)
				}
			},
		},
		{
			name:          "FileUploadException",
			exceptionName: "FileUploadException",
			assert: func(t *testing.T, err error) {
				var uploadTarget *FileUploadException
				if !errors.As(err, &uploadTarget) {
					t.Fatalf("expected FileUploadException, got %T", err)
				}
				var buildTarget *BuildException
				if !errors.As(err, &buildTarget) {
					t.Fatalf("expected FileUploadException to be a BuildException")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := FromAPIError(tt.exceptionName, payload, cause)
			tt.assert(t, err)
			if !errors.Is(err, cause) {
				t.Fatalf("expected wrapped cause")
			}
		})
	}
}

func TestFromAPIError_UnknownFallback(t *testing.T) {
	err := FromAPIError("RandomException", APIErrorPayload{Type: "unknown", Message: "m"}, nil)
	var base *BaseError
	if !errors.As(err, &base) {
		t.Fatalf("expected BaseError fallback, got %T", err)
	}
}
