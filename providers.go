package gollm

import (
	"github.com/grokify/gollm/provider"
	"github.com/grokify/gollm/providers/anthropic"
	"github.com/grokify/gollm/providers/bedrock"
	"github.com/grokify/gollm/providers/gemini"
	"github.com/grokify/gollm/providers/ollama"
	"github.com/grokify/gollm/providers/openai"
)

// newOpenAIProvider creates a new OpenAI provider adapter
func newOpenAIProvider(config ClientConfig) (provider.Provider, error) {
	if config.APIKey == "" {
		return nil, ErrEmptyAPIKey
	}
	return openai.NewProvider(config.APIKey, config.BaseURL), nil
}

// newAnthropicProvider creates a new Anthropic provider adapter
func newAnthropicProvider(config ClientConfig) (provider.Provider, error) {
	if config.APIKey == "" {
		return nil, ErrEmptyAPIKey
	}
	return anthropic.NewProvider(config.APIKey, config.BaseURL), nil
}

// newBedrockProvider creates a new Bedrock provider adapter
func newBedrockProvider(config ClientConfig) (provider.Provider, error) {
	return bedrock.NewProvider(config.Region)
}

// newOllamaProvider creates a new Ollama provider adapter
func newOllamaProvider(config ClientConfig) (provider.Provider, error) { //nolint:unparam // `error` added to fulfill interface requirements
	return ollama.NewProvider(config.BaseURL), nil
}

// newGeminiProvider creates a new Gemini provider adapter
func newGeminiProvider(config ClientConfig) (provider.Provider, error) {
	if config.APIKey == "" {
		return nil, ErrEmptyAPIKey
	}
	return gemini.NewProvider(config.APIKey), nil
}
