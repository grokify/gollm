package metallm

import (
	"context"
	"fmt"
	"time"

	"github.com/grokify/sogo/database/kvs"
)

// MemoryConfig holds configuration for conversation memory
type MemoryConfig struct {
	// MaxMessages limits the number of messages to keep in memory per session
	MaxMessages int
	// TTL sets the time-to-live for stored conversations (0 for no expiration)
	TTL time.Duration
	// KeyPrefix allows customizing the key prefix for stored conversations
	KeyPrefix string
}

// DefaultMemoryConfig returns sensible defaults for memory configuration
func DefaultMemoryConfig() MemoryConfig {
	return MemoryConfig{
		MaxMessages: 50,
		TTL:         24 * time.Hour,
		KeyPrefix:   "metallm:session",
	}
}

// ConversationMemory represents stored conversation data
type ConversationMemory struct {
	SessionID string         `json:"session_id"`
	Messages  []Message      `json:"messages"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

// MemoryManager handles conversation persistence using KVS
type MemoryManager struct {
	kvs    kvs.Client
	config MemoryConfig
}

// NewMemoryManager creates a new memory manager with the given KVS client and config
func NewMemoryManager(kvsClient kvs.Client, config MemoryConfig) *MemoryManager {
	return &MemoryManager{
		kvs:    kvsClient,
		config: config,
	}
}

// LoadConversation retrieves a conversation from memory
func (m *MemoryManager) LoadConversation(ctx context.Context, sessionID string) (*ConversationMemory, error) {
	if m.kvs == nil {
		return nil, fmt.Errorf("memory not configured")
	}

	key := m.buildKey(sessionID)

	var conversation ConversationMemory
	err := m.kvs.GetAny(ctx, key, &conversation)
	if err != nil {
		// Return empty conversation if not found
		return &ConversationMemory{
			SessionID: sessionID,
			Messages:  []Message{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Metadata:  make(map[string]any),
		}, nil
	}

	return &conversation, nil
}

// SaveConversation stores a conversation in memory
func (m *MemoryManager) SaveConversation(ctx context.Context, conversation *ConversationMemory) error {
	if m.kvs == nil {
		return fmt.Errorf("memory not configured")
	}

	// Apply message limit
	if m.config.MaxMessages > 0 && len(conversation.Messages) > m.config.MaxMessages {
		// Keep system messages and limit the rest
		systemMessages := []Message{}
		otherMessages := []Message{}

		for _, msg := range conversation.Messages {
			if msg.Role == RoleSystem {
				systemMessages = append(systemMessages, msg)
			} else {
				otherMessages = append(otherMessages, msg)
			}
		}

		// Keep the most recent messages within the limit
		maxOthers := m.config.MaxMessages - len(systemMessages)
		if maxOthers > 0 && len(otherMessages) > maxOthers {
			otherMessages = otherMessages[len(otherMessages)-maxOthers:]
		}

		conversation.Messages = append(systemMessages, otherMessages...)
	}

	conversation.UpdatedAt = time.Now()
	key := m.buildKey(conversation.SessionID)

	return m.kvs.SetAny(ctx, key, conversation)
}

// AppendMessage adds a message to the conversation and saves it
func (m *MemoryManager) AppendMessage(ctx context.Context, sessionID string, message Message) error {
	conversation, err := m.LoadConversation(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to load conversation: %w", err)
	}

	conversation.Messages = append(conversation.Messages, message)

	return m.SaveConversation(ctx, conversation)
}

// AppendMessages adds multiple messages to the conversation and saves it
func (m *MemoryManager) AppendMessages(ctx context.Context, sessionID string, messages []Message) error {
	conversation, err := m.LoadConversation(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to load conversation: %w", err)
	}

	conversation.Messages = append(conversation.Messages, messages...)

	return m.SaveConversation(ctx, conversation)
}

// DeleteConversation removes a conversation from memory
func (m *MemoryManager) DeleteConversation(ctx context.Context, sessionID string) error {
	if m.kvs == nil {
		return fmt.Errorf("memory not configured")
	}

	key := m.buildKey(sessionID)

	// Since the KVS interface doesn't have a Delete method, we'll set an empty value
	// This is a limitation of the current KVS interface
	return m.kvs.SetString(ctx, key, "")
}

// GetMessages returns just the messages from a conversation
func (m *MemoryManager) GetMessages(ctx context.Context, sessionID string) ([]Message, error) {
	conversation, err := m.LoadConversation(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return conversation.Messages, nil
}

// SetMetadata sets metadata for a conversation
func (m *MemoryManager) SetMetadata(ctx context.Context, sessionID string, metadata map[string]any) error {
	conversation, err := m.LoadConversation(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to load conversation: %w", err)
	}

	if conversation.Metadata == nil {
		conversation.Metadata = make(map[string]any)
	}

	for k, v := range metadata {
		conversation.Metadata[k] = v
	}

	return m.SaveConversation(ctx, conversation)
}

// buildKey constructs the storage key for a session
func (m *MemoryManager) buildKey(sessionID string) string {
	return fmt.Sprintf("%s:%s", m.config.KeyPrefix, sessionID)
}

// CreateConversationWithSystemMessage creates a new conversation with a system message
func (m *MemoryManager) CreateConversationWithSystemMessage(ctx context.Context, sessionID, systemMessage string) error {
	conversation := &ConversationMemory{
		SessionID: sessionID,
		Messages: []Message{
			{
				Role:    RoleSystem,
				Content: systemMessage,
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata:  make(map[string]any),
	}

	return m.SaveConversation(ctx, conversation)
}
