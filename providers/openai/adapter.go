// Package openai provides OpenAI provider adapter for the gollm unified interface
package openai

import (
	"context"

	"github.com/grokify/gollm/provider"
)

// Provider represents the OpenAI provider adapter
type Provider struct {
	client *Client
}

// NewProvider creates a new OpenAI provider adapter
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
	// Convert from unified format to OpenAI format
	openaiReq := &Request{
		Model:       req.Model,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		Stop:        req.Stop,
	}

	// Convert messages
	for _, msg := range req.Messages {
		openaiReq.Messages = append(openaiReq.Messages, Message{
			Role:    string(msg.Role),
			Content: msg.Content,
			Name:    msg.Name,
		})
	}

	resp, err := p.client.CreateCompletion(ctx, openaiReq)
	if err != nil {
		return nil, err
	}

	// Convert back to unified format
	return &provider.ChatCompletionResponse{
		ID:      resp.ID,
		Object:  resp.Object,
		Created: resp.Created,
		Model:   resp.Model,
		Choices: []provider.ChatCompletionChoice{
			{
				Index: 0,
				Message: provider.Message{
					Role:    provider.Role(resp.Choices[0].Message.Role),
					Content: resp.Choices[0].Message.Content,
				},
				FinishReason: resp.Choices[0].FinishReason,
			},
		},
		Usage: provider.Usage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
	}, nil
}

// CreateChatCompletionStream creates a streaming chat completion
func (p *Provider) CreateChatCompletionStream(ctx context.Context, req *provider.ChatCompletionRequest) (provider.ChatCompletionStream, error) {
	// Convert from unified format to OpenAI format
	openaiReq := &Request{
		Model:       req.Model,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		Stop:        req.Stop,
	}

	// Convert messages
	for _, msg := range req.Messages {
		openaiReq.Messages = append(openaiReq.Messages, Message{
			Role:    string(msg.Role),
			Content: msg.Content,
			Name:    msg.Name,
		})
	}

	stream, err := p.client.CreateCompletionStream(ctx, openaiReq)
	if err != nil {
		return nil, err
	}

	return &StreamAdapter{stream: stream}, nil
}

// Close closes the provider
func (p *Provider) Close() error {
	return p.client.Close()
}

// StreamAdapter adapts OpenAI stream to unified interface
type StreamAdapter struct {
	stream *Stream
}

// Recv receives the next chunk from the stream
func (s *StreamAdapter) Recv() (*provider.ChatCompletionChunk, error) {
	chunk, err := s.stream.Recv()
	if err != nil {
		return nil, err
	}

	// Convert to unified format
	result := &provider.ChatCompletionChunk{
		ID:      chunk.ID,
		Object:  chunk.Object,
		Created: chunk.Created,
		Model:   chunk.Model,
	}

	if chunk.Usage != nil {
		result.Usage = &provider.Usage{
			PromptTokens:     chunk.Usage.PromptTokens,
			CompletionTokens: chunk.Usage.CompletionTokens,
			TotalTokens:      chunk.Usage.TotalTokens,
		}
	}

	for _, choice := range chunk.Choices {
		result.Choices = append(result.Choices, provider.ChatCompletionChoice{
			Index:        choice.Index,
			FinishReason: choice.FinishReason,
		})
		if choice.Delta != nil {
			result.Choices[len(result.Choices)-1].Delta = &provider.Message{
				Role:    provider.Role(choice.Delta.Role),
				Content: choice.Delta.Content,
			}
		}
	}

	return result, nil
}

// Close closes the stream
func (s *StreamAdapter) Close() error {
	return s.stream.Close()
}