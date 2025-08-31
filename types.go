package gollm

// Role represents the role of a message sender
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

// Message represents a chat message
type Message struct {
	Role       Role       `json:"role"`
	Content    string     `json:"content"`
	Name       *string    `json:"name,omitempty"`
	ToolCallID *string    `json:"tool_call_id,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
}

// ToolCall represents a tool function call
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

// ToolFunction represents the function being called
type ToolFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ChatCompletionRequest represents a request for chat completion
type ChatCompletionRequest struct {
	Model            string         `json:"model"`
	Messages         []Message      `json:"messages"`
	MaxTokens        *int           `json:"max_tokens,omitempty"`
	Temperature      *float64       `json:"temperature,omitempty"`
	TopP             *float64       `json:"top_p,omitempty"`
	Stream           *bool          `json:"stream,omitempty"`
	Stop             []string       `json:"stop,omitempty"`
	PresencePenalty  *float64       `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float64       `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int `json:"logit_bias,omitempty"`
	User             *string        `json:"user,omitempty"`
	Tools            []Tool         `json:"tools,omitempty"`
	ToolChoice       interface{}    `json:"tool_choice,omitempty"`
}

// Tool represents a tool that can be called
type Tool struct {
	Type     string   `json:"type"`
	Function ToolSpec `json:"function"`
}

// ToolSpec defines a tool specification
type ToolSpec struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  interface{} `json:"parameters"`
}

// ChatCompletionResponse represents a response from chat completion
type ChatCompletionResponse struct {
	ID                string                 `json:"id"`
	Object            string                 `json:"object"`
	Created           int64                  `json:"created"`
	Model             string                 `json:"model"`
	SystemFingerprint *string                `json:"system_fingerprint,omitempty"`
	Choices           []ChatCompletionChoice `json:"choices"`
	Usage             Usage                  `json:"usage"`
}

// ChatCompletionChoice represents a single choice in the response
type ChatCompletionChoice struct {
	Index        int         `json:"index"`
	Message      Message     `json:"message"`
	Delta        *Message    `json:"delta,omitempty"`
	FinishReason *string     `json:"finish_reason"`
	Logprobs     interface{} `json:"logprobs,omitempty"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatCompletionChunk represents a chunk in streaming response
type ChatCompletionChunk struct {
	ID                string                 `json:"id"`
	Object            string                 `json:"object"`
	Created           int64                  `json:"created"`
	Model             string                 `json:"model"`
	SystemFingerprint *string                `json:"system_fingerprint,omitempty"`
	Choices           []ChatCompletionChoice `json:"choices"`
	Usage             *Usage                 `json:"usage,omitempty"`
}

// ModelInfo represents information about a model
type ModelInfo struct {
	ID        string       `json:"id"`
	Provider  ProviderName `json:"provider"`
	Name      string       `json:"name"`
	MaxTokens int          `json:"max_tokens"`
}

// GetModelInfo returns model information
func GetModelInfo(modelID string) *ModelInfo {
	modelMap := map[string]ModelInfo{
		ModelGPT4o: {
			ID:        ModelGPT4o,
			Provider:  ProviderNameOpenAI,
			Name:      "GPT-4o",
			MaxTokens: 128000,
		},
		ModelClaude3Opus: {
			ID:        ModelClaude3Opus,
			Provider:  ProviderNameAnthropic,
			Name:      "Claude 3 Opus",
			MaxTokens: 200000,
		},
		ModelBedrockClaude3Sonnet: {
			ID:        ModelBedrockClaude3Sonnet,
			Provider:  ProviderNameBedrock,
			Name:      "Claude 3 Sonnet (Bedrock)",
			MaxTokens: 200000,
		},
		ModelOllamaLlama3_8B: {
			ID:        ModelOllamaLlama3_8B,
			Provider:  ProviderNameOllama,
			Name:      "Llama 3 8B",
			MaxTokens: 8192,
		},
		ModelOllamaMistral7B: {
			ID:        ModelOllamaMistral7B,
			Provider:  ProviderNameOllama,
			Name:      "Mistral 7B",
			MaxTokens: 32768,
		},
		ModelOllamaCodeLlama: {
			ID:        ModelOllamaCodeLlama,
			Provider:  ProviderNameOllama,
			Name:      "CodeLlama 13B",
			MaxTokens: 16384,
		},
	}

	if info, exists := modelMap[modelID]; exists {
		return &info
	}
	return nil
}
