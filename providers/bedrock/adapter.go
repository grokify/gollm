// Package bedrock provides AWS Bedrock provider adapter for the gollm unified interface
package bedrock

import (
	"context"
	"fmt"

	"github.com/grokify/gollm/provider"
)

// Provider represents the Bedrock provider adapter
type Provider struct {
	client *Client
}

// NewProvider creates a new Bedrock provider adapter
func NewProvider(region string) (provider.Provider, error) {
	client, err := New(region)
	if err != nil {
		return nil, err
	}
	return &Provider{client: client}, nil
}

// Name returns the provider name
func (p *Provider) Name() string {
	return p.client.Name()
}

// CreateChatCompletion creates a chat completion
func (p *Provider) CreateChatCompletion(ctx context.Context, req *provider.ChatCompletionRequest) (*provider.ChatCompletionResponse, error) {
	return nil, fmt.Errorf("bedrock implementation not fully implemented in this demo")
}

// CreateChatCompletionStream creates a streaming chat completion
func (p *Provider) CreateChatCompletionStream(ctx context.Context, req *provider.ChatCompletionRequest) (provider.ChatCompletionStream, error) {
	return nil, fmt.Errorf("bedrock streaming not implemented in this demo")
}

// Close closes the provider
func (p *Provider) Close() error {
	return p.client.Close()
}
