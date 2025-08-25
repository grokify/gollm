package gollm

import "context"

// Provider defines the interface that all LLM providers must implement
type Provider interface {
	// CreateChatCompletion creates a new chat completion
	CreateChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error)
	
	// CreateChatCompletionStream creates a streaming chat completion
	CreateChatCompletionStream(ctx context.Context, req *ChatCompletionRequest) (ChatCompletionStream, error)
	
	// Close closes the provider and cleans up resources
	Close() error
	
	// Name returns the provider name
	Name() string
}

// ProviderName represents the different LLM provider names
type ProviderName string

const (
	ProviderNameOpenAI    ProviderName = "openai"
	ProviderNameAnthropic ProviderName = "anthropic" 
	ProviderNameBedrock   ProviderName = "bedrock"
)