package omnillm

import (
	"github.com/agentplexus/omnillm/provider"
	"github.com/agentplexus/omnillm/providers/anthropic"
	"github.com/agentplexus/omnillm/providers/gemini"
	"github.com/agentplexus/omnillm/providers/ollama"
	"github.com/agentplexus/omnillm/providers/openai"
	"github.com/agentplexus/omnillm/providers/xai"
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
