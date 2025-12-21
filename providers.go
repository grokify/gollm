package fluxllm

import (
	"github.com/grokify/fluxllm/provider"
	"github.com/grokify/fluxllm/providers/anthropic"
	"github.com/grokify/fluxllm/providers/bedrock"
	"github.com/grokify/fluxllm/providers/gemini"
	"github.com/grokify/fluxllm/providers/ollama"
	"github.com/grokify/fluxllm/providers/openai"
	"github.com/grokify/fluxllm/providers/xai"
)

// newOpenAIProvider creates a new OpenAI provider adapter
func newOpenAIProvider(config ClientConfig) (provider.Provider, error) {
	if config.APIKey == "" {
		return nil, ErrEmptyAPIKey
	}
	return openai.NewProvider(config.APIKey, config.BaseURL, config.HTTPClient), nil
}

// newAnthropicProvider creates a new Anthropic provider adapter
func newAnthropicProvider(config ClientConfig) (provider.Provider, error) {
	if config.APIKey == "" {
		return nil, ErrEmptyAPIKey
	}
	return anthropic.NewProvider(config.APIKey, config.BaseURL, config.HTTPClient), nil
}

// newBedrockProvider creates a new Bedrock provider adapter
func newBedrockProvider(config ClientConfig) (provider.Provider, error) {
	return bedrock.NewProvider(config.Region)
}

// newOllamaProvider creates a new Ollama provider adapter
func newOllamaProvider(config ClientConfig) (provider.Provider, error) { //nolint:unparam // `error` added to fulfill interface requirements
	return ollama.NewProvider(config.BaseURL, config.HTTPClient), nil
}

// newGeminiProvider creates a new Gemini provider adapter
func newGeminiProvider(config ClientConfig) (provider.Provider, error) {
	if config.APIKey == "" {
		return nil, ErrEmptyAPIKey
	}
	return gemini.NewProvider(config.APIKey), nil
}

// newXAIProvider creates a new X.AI provider adapter
func newXAIProvider(config ClientConfig) (provider.Provider, error) {
	if config.APIKey == "" {
		return nil, ErrEmptyAPIKey
	}
	return xai.NewProvider(config.APIKey, config.BaseURL, config.HTTPClient), nil
}
