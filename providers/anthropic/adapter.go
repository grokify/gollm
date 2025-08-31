// Package anthropic provides Anthropic provider adapter for the gollm unified interface
package anthropic

import (
	"context"
	"fmt"
	"time"

	"github.com/grokify/gollm/provider"
)

// Provider represents the Anthropic provider adapter
type Provider struct {
	client *Client
}

// NewProvider creates a new Anthropic provider adapter
func NewProvider(apiKey, baseURL string) provider.Provider {
	client := New(apiKey, baseURL)
	return &Provider{client: client}
}

// Name returns the provider name
func (p *Provider) Name() string {
	return p.client.Name()
}

// CreateChatCompletion creates a chat completion
func (p *Provider) CreateChatCompletion(ctx context.Context, req *provider.ChatCompletionRequest) (*provider.ChatCompletionResponse, error) {
	// Convert from unified format to Anthropic format
	anthropicReq := &Request{
		Model:       req.Model,
		MaxTokens:   4096, // Default
		Temperature: req.Temperature,
		TopP:        req.TopP,
	}

	if req.MaxTokens != nil {
		anthropicReq.MaxTokens = *req.MaxTokens
	}

	// Convert messages (Anthropic separates system messages)
	var systemMessage string
	for _, msg := range req.Messages {
		switch msg.Role {
		case provider.RoleSystem:
			systemMessage = msg.Content
		case provider.RoleUser, provider.RoleAssistant:
			anthropicReq.Messages = append(anthropicReq.Messages, Message{
				Role:    string(msg.Role),
				Content: msg.Content,
			})
		}
	}

	if systemMessage != "" {
		anthropicReq.System = systemMessage
	}

	resp, err := p.client.CreateCompletion(ctx, anthropicReq)
	if err != nil {
		return nil, err
	}

	// Convert back to unified format
	var content string
	if len(resp.Content) > 0 && resp.Content[0].Type == "text" {
		content = resp.Content[0].Text
	}

	return &provider.ChatCompletionResponse{
		ID:      resp.ID,
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   resp.Model,
		Choices: []provider.ChatCompletionChoice{
			{
				Index: 0,
				Message: provider.Message{
					Role:    provider.RoleAssistant,
					Content: content,
				},
				FinishReason: &resp.StopReason,
			},
		},
		Usage: provider.Usage{
			PromptTokens:     resp.Usage.InputTokens,
			CompletionTokens: resp.Usage.OutputTokens,
			TotalTokens:      resp.Usage.InputTokens + resp.Usage.OutputTokens,
		},
	}, nil
}

// CreateChatCompletionStream creates a streaming chat completion
func (p *Provider) CreateChatCompletionStream(ctx context.Context, req *provider.ChatCompletionRequest) (provider.ChatCompletionStream, error) {
	return nil, fmt.Errorf("anthropic streaming not implemented in this demo")
}

// Close closes the provider
func (p *Provider) Close() error {
	return p.client.Close()
}