package gollm

import (
	"context"
)

// ChatClient is the main client interface that wraps a Provider
type ChatClient struct {
	provider Provider
}

// ChatCompletionStream represents a streaming chat completion response
type ChatCompletionStream interface {
	// Recv receives the next chunk from the stream
	Recv() (*ChatCompletionChunk, error)

	// Close closes the stream
	Close() error
}

// ClientConfig holds configuration for creating a client
type ClientConfig struct {
	Provider ProviderName
	APIKey   string
	BaseURL  string
	Region   string // For AWS Bedrock

	// Provider-specific configurations can be added here
	Extra map[string]interface{}
}

// NewClient creates a new ChatClient based on the provider
func NewClient(config ClientConfig) (*ChatClient, error) {
	var provider Provider
	var err error

	switch config.Provider {
	case ProviderNameOpenAI:
		provider, err = newOpenAIProvider(config)
	case ProviderNameAnthropic:
		provider, err = newAnthropicProvider(config)
	case ProviderNameBedrock:
		provider, err = newBedrockProvider(config)
	default:
		return nil, ErrUnsupportedProvider
	}

	if err != nil {
		return nil, err
	}

	return &ChatClient{provider: provider}, nil
}

// CreateChatCompletion creates a chat completion
func (c *ChatClient) CreateChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	return c.provider.CreateChatCompletion(ctx, req)
}

// CreateChatCompletionStream creates a streaming chat completion
func (c *ChatClient) CreateChatCompletionStream(ctx context.Context, req *ChatCompletionRequest) (ChatCompletionStream, error) {
	return c.provider.CreateChatCompletionStream(ctx, req)
}

// Close closes the client
func (c *ChatClient) Close() error {
	return c.provider.Close()
}

// Provider returns the underlying provider
func (c *ChatClient) Provider() Provider {
	return c.provider
}
