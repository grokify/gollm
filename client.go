package gollm

import (
	"context"
	"fmt"
	"strings"

	"github.com/grokify/sogo/database/kvs"
)

// ChatClient is the main client interface that wraps a Provider
type ChatClient struct {
	provider Provider
	memory   *MemoryManager
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

	// Memory configuration (optional)
	Memory       kvs.Client
	MemoryConfig *MemoryConfig

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
	case ProviderNameOllama:
		provider, err = newOllamaProvider(config)
	default:
		return nil, ErrUnsupportedProvider
	}

	if err != nil {
		return nil, err
	}

	client := &ChatClient{provider: provider}

	// Initialize memory if provided
	if config.Memory != nil {
		memoryConfig := DefaultMemoryConfig()
		if config.MemoryConfig != nil {
			memoryConfig = *config.MemoryConfig
		}
		client.memory = NewMemoryManager(config.Memory, memoryConfig)
	}

	return client, nil
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

// Memory returns the memory manager (nil if not configured)
func (c *ChatClient) Memory() *MemoryManager {
	return c.memory
}

// HasMemory returns true if memory is configured
func (c *ChatClient) HasMemory() bool {
	return c.memory != nil
}

// CreateChatCompletionWithMemory creates a chat completion using conversation memory
func (c *ChatClient) CreateChatCompletionWithMemory(ctx context.Context, sessionID string, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	if !c.HasMemory() {
		return c.CreateChatCompletion(ctx, req)
	}

	// Load existing conversation
	conversation, err := c.memory.LoadConversation(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Merge stored messages with request messages
	allMessages := append(conversation.Messages, req.Messages...)

	// Create new request with combined messages
	memoryReq := *req
	memoryReq.Messages = allMessages

	// Get response
	response, err := c.provider.CreateChatCompletion(ctx, &memoryReq)
	if err != nil {
		return nil, err
	}

	// Save the conversation with new messages and response
	if len(response.Choices) > 0 {
		// Save request messages and response
		messagesToSave := append(req.Messages, response.Choices[0].Message)
		err = c.memory.AppendMessages(ctx, sessionID, messagesToSave)
		if err != nil {
			// Log error but don't fail the request
			// In a production environment, you might want to use a proper logger here
		}
	}

	return response, nil
}

// CreateChatCompletionStreamWithMemory creates a streaming chat completion using conversation memory
func (c *ChatClient) CreateChatCompletionStreamWithMemory(ctx context.Context, sessionID string, req *ChatCompletionRequest) (ChatCompletionStream, error) {
	if !c.HasMemory() {
		return c.CreateChatCompletionStream(ctx, req)
	}

	// Load existing conversation
	conversation, err := c.memory.LoadConversation(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Merge stored messages with request messages
	allMessages := append(conversation.Messages, req.Messages...)

	// Create new request with combined messages
	memoryReq := *req
	memoryReq.Messages = allMessages

	// Get stream response
	stream, err := c.provider.CreateChatCompletionStream(ctx, &memoryReq)
	if err != nil {
		return nil, err
	}

	// Wrap the stream to capture the response for memory storage
	return &memoryAwareStream{
		stream:      stream,
		memory:      c.memory,
		sessionID:   sessionID,
		reqMessages: req.Messages,
		ctx:         ctx,
	}, nil
}

// LoadConversation loads a conversation from memory
func (c *ChatClient) LoadConversation(ctx context.Context, sessionID string) (*ConversationMemory, error) {
	if !c.HasMemory() {
		return nil, fmt.Errorf("memory not configured")
	}
	return c.memory.LoadConversation(ctx, sessionID)
}

// SaveConversation saves a conversation to memory
func (c *ChatClient) SaveConversation(ctx context.Context, conversation *ConversationMemory) error {
	if !c.HasMemory() {
		return fmt.Errorf("memory not configured")
	}
	return c.memory.SaveConversation(ctx, conversation)
}

// AppendMessage appends a message to a conversation in memory
func (c *ChatClient) AppendMessage(ctx context.Context, sessionID string, message Message) error {
	if !c.HasMemory() {
		return fmt.Errorf("memory not configured")
	}
	return c.memory.AppendMessage(ctx, sessionID, message)
}

// GetConversationMessages retrieves messages from a conversation
func (c *ChatClient) GetConversationMessages(ctx context.Context, sessionID string) ([]Message, error) {
	if !c.HasMemory() {
		return nil, fmt.Errorf("memory not configured")
	}
	return c.memory.GetMessages(ctx, sessionID)
}

// CreateConversationWithSystemMessage creates a new conversation with a system message
func (c *ChatClient) CreateConversationWithSystemMessage(ctx context.Context, sessionID, systemMessage string) error {
	if !c.HasMemory() {
		return fmt.Errorf("memory not configured")
	}
	return c.memory.CreateConversationWithSystemMessage(ctx, sessionID, systemMessage)
}

// DeleteConversation removes a conversation from memory
func (c *ChatClient) DeleteConversation(ctx context.Context, sessionID string) error {
	if !c.HasMemory() {
		return fmt.Errorf("memory not configured")
	}
	return c.memory.DeleteConversation(ctx, sessionID)
}

// memoryAwareStream wraps a ChatCompletionStream to capture responses for memory storage
type memoryAwareStream struct {
	stream      ChatCompletionStream
	memory      *MemoryManager
	sessionID   string
	reqMessages []Message
	ctx         context.Context

	// Buffer to collect the complete response
	responseBuffer strings.Builder
	streamClosed   bool
}

// Recv receives the next chunk from the stream and buffers the response
func (s *memoryAwareStream) Recv() (*ChatCompletionChunk, error) {
	chunk, err := s.stream.Recv()
	if err != nil {
		// If we hit EOF and haven't saved the response yet, save it now
		if err.Error() == "EOF" && !s.streamClosed {
			s.saveBufferedResponse()
			s.streamClosed = true
		}
		return chunk, err
	}

	// Buffer the response content
	if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
		s.responseBuffer.WriteString(chunk.Choices[0].Delta.Content)
	}

	return chunk, nil
}

// Close closes the stream and saves the complete response to memory
func (s *memoryAwareStream) Close() error {
	if !s.streamClosed {
		s.saveBufferedResponse()
		s.streamClosed = true
	}
	return s.stream.Close()
}

// saveBufferedResponse saves the complete buffered response to memory
func (s *memoryAwareStream) saveBufferedResponse() {
	if s.responseBuffer.Len() > 0 {
		// Create assistant message from buffered response
		assistantMessage := Message{
			Role:    RoleAssistant,
			Content: s.responseBuffer.String(),
		}

		// Save request messages and response
		messagesToSave := append(s.reqMessages, assistantMessage)
		err := s.memory.AppendMessages(s.ctx, s.sessionID, messagesToSave)
		if err != nil {
			// Log error but don't fail the stream
			// In a production environment, you might want to use a proper logger here
		}
	}
}
