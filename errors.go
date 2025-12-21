package fluxllm

import (
	"errors"
	"fmt"
)

var (
	// Common errors
	ErrUnsupportedProvider  = errors.New("unsupported provider")
	ErrInvalidConfiguration = errors.New("invalid configuration")
	ErrEmptyAPIKey          = errors.New("API key cannot be empty")
	ErrEmptyModel           = errors.New("model cannot be empty")
	ErrEmptyMessages        = errors.New("messages cannot be empty")
	ErrStreamClosed         = errors.New("stream is closed")
	ErrInvalidResponse      = errors.New("invalid response format")
	ErrRateLimitExceeded    = errors.New("rate limit exceeded")
	ErrQuotaExceeded        = errors.New("quota exceeded")
	ErrInvalidRequest       = errors.New("invalid request")
	ErrModelNotFound        = errors.New("model not found")
	ErrServerError          = errors.New("server error")
	ErrNetworkError         = errors.New("network error")
)

// APIError represents an error response from the API
type APIError struct {
	StatusCode int          `json:"status_code"`
	Message    string       `json:"message"`
	Type       string       `json:"type"`
	Code       string       `json:"code"`
	Provider   ProviderName `json:"provider"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("[%s] %s (status: %d, type: %s, code: %s)",
		e.Provider, e.Message, e.StatusCode, e.Type, e.Code)
}

// NewAPIError creates a new API error
func NewAPIError(provider ProviderName, statusCode int, message, errorType, code string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		Type:       errorType,
		Code:       code,
		Provider:   provider,
	}
}
