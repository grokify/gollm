package gollm

import (
	"context"
	"fmt"
	"time"

	"github.com/grokify/gollm/providers/anthropic"
	"github.com/grokify/gollm/providers/bedrock"
	"github.com/grokify/gollm/providers/openai"
)

// OpenAI Provider Adapter
type openAIProvider struct {
	client *openai.Client
}

func newOpenAIProvider(config ClientConfig) (Provider, error) {
	if config.APIKey == "" {
		return nil, ErrEmptyAPIKey
	}

	client := openai.New(config.APIKey, config.BaseURL)
	return &openAIProvider{client: client}, nil
}

func (p *openAIProvider) Name() string {
	return p.client.Name()
}

func (p *openAIProvider) CreateChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// Convert from unified format to OpenAI format
	openaiReq := &openai.Request{
		Model:       req.Model,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		Stop:        req.Stop,
	}

	// Convert messages
	for _, msg := range req.Messages {
		openaiReq.Messages = append(openaiReq.Messages, openai.Message{
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
	return &ChatCompletionResponse{
		ID:      resp.ID,
		Object:  resp.Object,
		Created: resp.Created,
		Model:   resp.Model,
		Choices: []ChatCompletionChoice{
			{
				Index: 0,
				Message: Message{
					Role:    Role(resp.Choices[0].Message.Role),
					Content: resp.Choices[0].Message.Content,
				},
				FinishReason: resp.Choices[0].FinishReason,
			},
		},
		Usage: Usage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
	}, nil
}

func (p *openAIProvider) CreateChatCompletionStream(ctx context.Context, req *ChatCompletionRequest) (ChatCompletionStream, error) {
	// Convert from unified format to OpenAI format
	openaiReq := &openai.Request{
		Model:       req.Model,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		Stop:        req.Stop,
	}

	// Convert messages
	for _, msg := range req.Messages {
		openaiReq.Messages = append(openaiReq.Messages, openai.Message{
			Role:    string(msg.Role),
			Content: msg.Content,
			Name:    msg.Name,
		})
	}

	stream, err := p.client.CreateCompletionStream(ctx, openaiReq)
	if err != nil {
		return nil, err
	}

	return &openAIStreamAdapter{stream: stream}, nil
}

func (p *openAIProvider) Close() error {
	return p.client.Close()
}

// OpenAI Stream Adapter
type openAIStreamAdapter struct {
	stream *openai.Stream
}

func (s *openAIStreamAdapter) Recv() (*ChatCompletionChunk, error) {
	chunk, err := s.stream.Recv()
	if err != nil {
		return nil, err
	}

	// Convert to unified format
	result := &ChatCompletionChunk{
		ID:      chunk.ID,
		Object:  chunk.Object,
		Created: chunk.Created,
		Model:   chunk.Model,
	}
	
	if chunk.Usage != nil {
		result.Usage = &Usage{
			PromptTokens:     chunk.Usage.PromptTokens,
			CompletionTokens: chunk.Usage.CompletionTokens,
			TotalTokens:      chunk.Usage.TotalTokens,
		}
	}

	for _, choice := range chunk.Choices {
		result.Choices = append(result.Choices, ChatCompletionChoice{
			Index:        choice.Index,
			FinishReason: choice.FinishReason,
		})
		if choice.Delta != nil {
			result.Choices[len(result.Choices)-1].Delta = &Message{
				Role:    Role(choice.Delta.Role),
				Content: choice.Delta.Content,
			}
		}
	}

	return result, nil
}

func (s *openAIStreamAdapter) Close() error {
	return s.stream.Close()
}

// Anthropic Provider Adapter
type anthropicProvider struct {
	client *anthropic.Client
}

func newAnthropicProvider(config ClientConfig) (Provider, error) {
	if config.APIKey == "" {
		return nil, ErrEmptyAPIKey
	}

	client := anthropic.New(config.APIKey, config.BaseURL)
	return &anthropicProvider{client: client}, nil
}

func (p *anthropicProvider) Name() string {
	return p.client.Name()
}

func (p *anthropicProvider) CreateChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// Convert from unified format to Anthropic format
	anthropicReq := &anthropic.Request{
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
		case RoleSystem:
			systemMessage = msg.Content
		case RoleUser, RoleAssistant:
			anthropicReq.Messages = append(anthropicReq.Messages, anthropic.Message{
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

	return &ChatCompletionResponse{
		ID:      resp.ID,
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   resp.Model,
		Choices: []ChatCompletionChoice{
			{
				Index: 0,
				Message: Message{
					Role:    RoleAssistant,
					Content: content,
				},
				FinishReason: &resp.StopReason,
			},
		},
		Usage: Usage{
			PromptTokens:     resp.Usage.InputTokens,
			CompletionTokens: resp.Usage.OutputTokens,
			TotalTokens:      resp.Usage.InputTokens + resp.Usage.OutputTokens,
		},
	}, nil
}

func (p *anthropicProvider) CreateChatCompletionStream(ctx context.Context, req *ChatCompletionRequest) (ChatCompletionStream, error) {
	return nil, fmt.Errorf("anthropic streaming not implemented in this demo")
}

func (p *anthropicProvider) Close() error {
	return p.client.Close()
}

// Bedrock Provider Adapter
type bedrockProvider struct {
	client *bedrock.Client
}

func newBedrockProvider(config ClientConfig) (Provider, error) {
	client, err := bedrock.New(config.Region)
	if err != nil {
		return nil, err
	}

	return &bedrockProvider{client: client}, nil
}

func (p *bedrockProvider) Name() string {
	return p.client.Name()
}

func (p *bedrockProvider) CreateChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	return nil, fmt.Errorf("bedrock implementation not fully implemented in this demo")
}

func (p *bedrockProvider) CreateChatCompletionStream(ctx context.Context, req *ChatCompletionRequest) (ChatCompletionStream, error) {
	return nil, fmt.Errorf("bedrock streaming not implemented in this demo")
}

func (p *bedrockProvider) Close() error {
	return p.client.Close()
}