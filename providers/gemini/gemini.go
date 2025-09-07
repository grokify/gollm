// Package gemini provides Google Gemini API client implementation
package gemini

import (
	"context"
	"fmt"
	"io"
	"time"

	"google.golang.org/genai"
)

// Client implements Google Gemini API client
type Client struct {
	client  *genai.Client
	ctx     context.Context
	initErr error
}

// New creates a new Gemini client
func New(apiKey string) *Client {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})

	// For simplicity, we'll store the error and handle it during first use
	// In a production implementation, you might want to return the error here
	return &Client{
		client:  client,
		ctx:     ctx,
		initErr: err,
	}
}

// NewWithContext creates a new Gemini client with context
func NewWithContext(ctx context.Context, apiKey string) (*Client, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &Client{
		client: client,
		ctx:    ctx,
	}, nil
}

// Name returns the provider name
func (c *Client) Name() string {
	return "gemini"
}

// CreateCompletion creates a chat completion
func (c *Client) CreateCompletion(ctx context.Context, req *Request) (*Response, error) {
	if c.initErr != nil {
		return nil, fmt.Errorf("client initialization failed: %w", c.initErr)
	}
	if req.Model == "" {
		return nil, fmt.Errorf("model cannot be empty")
	}
	if len(req.Messages) == 0 {
		return nil, fmt.Errorf("messages cannot be empty")
	}

	// Create a chat session
	chat, err := c.client.Chats.Create(ctx, req.Model, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat: %w", err)
	}

	// Convert messages to Gemini format
	parts := make([]*genai.Part, 0, len(req.Messages))
	for _, msg := range req.Messages {
		if msg.Content != "" {
			parts = append(parts, genai.NewPartFromText(msg.Content))
		}
	}

	// Send the message and get response
	response, err := chat.Send(ctx, parts...)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	// Convert response to our format
	result := &Response{
		ID:      generateID(),
		Object:  "chat.completion",
		Created: currentTimestamp(),
		Model:   req.Model,
	}

	if response.Candidates != nil && len(response.Candidates) > 0 {
		candidate := response.Candidates[0]
		content := ""

		// Extract text content from the candidate
		if candidate.Content != nil && candidate.Content.Parts != nil {
			for _, part := range candidate.Content.Parts {
				if part.Text != "" {
					content += part.Text
				}
			}
		}

		choice := Choice{
			Index: 0,
			Message: Message{
				Role:    "assistant",
				Content: content,
			},
		}

		if candidate.FinishReason != "" {
			reason := string(candidate.FinishReason)
			choice.FinishReason = &reason
		}

		result.Choices = []Choice{choice}
	}

	// Set usage information (Gemini doesn't provide detailed token counts)
	result.Usage = Usage{
		PromptTokens:     estimateTokens(req.Messages),
		CompletionTokens: estimateTokens(result.Choices),
		TotalTokens:      0, // Will be calculated
	}
	result.Usage.TotalTokens = result.Usage.PromptTokens + result.Usage.CompletionTokens

	return result, nil
}

// CreateCompletionStream creates a streaming chat completion
func (c *Client) CreateCompletionStream(ctx context.Context, req *Request) (*Stream, error) {
	if c.initErr != nil {
		return nil, fmt.Errorf("client initialization failed: %w", c.initErr)
	}
	if req.Model == "" {
		return nil, fmt.Errorf("model cannot be empty")
	}
	if len(req.Messages) == 0 {
		return nil, fmt.Errorf("messages cannot be empty")
	}

	// Create a chat session
	chat, err := c.client.Chats.Create(ctx, req.Model, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat: %w", err)
	}

	// Convert messages to Gemini format
	parts := make([]*genai.Part, 0, len(req.Messages))
	for _, msg := range req.Messages {
		if msg.Content != "" {
			parts = append(parts, genai.NewPartFromText(msg.Content))
		}
	}

	// Send the message with streaming
	stream := chat.SendStream(ctx, parts...)

	// Collect all responses from the stream
	var responses []*genai.GenerateContentResponse
	var errors []error

	for response, err := range stream {
		if err != nil {
			errors = append(errors, err)
		} else {
			responses = append(responses, response)
		}
	}

	return &Stream{
		responses: responses,
		errors:    errors,
		model:     req.Model,
		index:     0,
	}, nil
}

// Close closes the client
func (c *Client) Close() error {
	// The genai.Client doesn't have a Close method, so we just return nil
	return nil
}

// Stream represents a streaming response
type Stream struct {
	responses []*genai.GenerateContentResponse
	errors    []error
	model     string
	index     int
}

// Recv receives the next chunk from the stream
func (s *Stream) Recv() (*Chunk, error) {
	if s.index >= len(s.responses) {
		return nil, io.EOF
	}

	if s.index < len(s.errors) && s.errors[s.index] != nil {
		err := s.errors[s.index]
		s.index++
		return nil, fmt.Errorf("failed to receive stream chunk: %w", err)
	}

	response := s.responses[s.index]
	s.index++

	chunk := &Chunk{
		ID:      generateID(),
		Object:  "chat.completion.chunk",
		Created: currentTimestamp(),
		Model:   s.model,
	}

	if response.Candidates != nil && len(response.Candidates) > 0 {
		candidate := response.Candidates[0]
		content := ""

		// Extract text content from the candidate
		if candidate.Content != nil && candidate.Content.Parts != nil {
			for _, part := range candidate.Content.Parts {
				if part.Text != "" {
					content += part.Text
				}
			}
		}

		choice := Choice{
			Index: 0,
			Delta: &Message{
				Role:    "assistant",
				Content: content,
			},
		}

		if candidate.FinishReason != "" {
			reason := string(candidate.FinishReason)
			choice.FinishReason = &reason
		}

		chunk.Choices = []Choice{choice}
	}

	return chunk, nil
}

// Close closes the stream
func (s *Stream) Close() error {
	// Gemini stream iterator doesn't have explicit close
	return nil
}

// Helper functions

func generateID() string {
	return fmt.Sprintf("chatcmpl-%d", currentTimestamp())
}

func currentTimestamp() int64 {
	return time.Now().Unix()
}

func estimateTokens(data interface{}) int {
	// Simple token estimation - in a real implementation you might use
	// a more sophisticated method or actual token counting
	switch v := data.(type) {
	case []Message:
		total := 0
		for _, msg := range v {
			total += len(msg.Content) / 4 // Rough estimation: 4 chars per token
		}
		return total
	case []Choice:
		total := 0
		for _, choice := range v {
			total += len(choice.Message.Content) / 4
		}
		return total
	default:
		return 0
	}
}
