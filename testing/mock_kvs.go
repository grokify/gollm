// Package testing provides mock implementations for testing
package testing

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// MockKVS is a simple in-memory key-value store for testing
type MockKVS struct {
	mu    sync.RWMutex
	store map[string]string
}

// NewMockKVS creates a new mock KVS client
func NewMockKVS() *MockKVS {
	return &MockKVS{
		store: make(map[string]string),
	}
}

// SetString stores a string value
func (m *MockKVS) SetString(ctx context.Context, key, val string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.store[key] = val
	return nil
}

// GetString retrieves a string value
func (m *MockKVS) GetString(ctx context.Context, key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, exists := m.store[key]
	if !exists {
		return "", fmt.Errorf("key not found: %s", key)
	}
	return val, nil
}

// GetOrDefaultString retrieves a string value or returns default
func (m *MockKVS) GetOrDefaultString(ctx context.Context, key, def string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, exists := m.store[key]
	if !exists {
		return def
	}
	return val
}

// SetAny stores any value as JSON
func (m *MockKVS) SetAny(ctx context.Context, key string, val any) error {
	data, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	return m.SetString(ctx, key, string(data))
}

// GetAny retrieves a value and unmarshals it
func (m *MockKVS) GetAny(ctx context.Context, key string, val any) error {
	str, err := m.GetString(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(str), val)
}

// Delete removes a key (helper for testing)
func (m *MockKVS) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.store, key)
}

// Clear removes all keys (helper for testing)
func (m *MockKVS) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.store = make(map[string]string)
}

// Keys returns all keys (helper for testing)
func (m *MockKVS) Keys() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	keys := make([]string, 0, len(m.store))
	for k := range m.store {
		keys = append(keys, k)
	}
	return keys
}

// Size returns the number of keys (helper for testing)
func (m *MockKVS) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.store)
}
