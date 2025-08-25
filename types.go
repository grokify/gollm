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
	Role       Role        `json:"role"`
	Content    string      `json:"content"`
	Name       *string     `json:"name,omitempty"`
	ToolCallID *string     `json:"tool_call_id,omitempty"`
	ToolCalls  []ToolCall  `json:"tool_calls,omitempty"`
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
	Model            string    `json:"model"`
	Messages         []Message `json:"messages"`
	MaxTokens        *int      `json:"max_tokens,omitempty"`
	Temperature      *float64  `json:"temperature,omitempty"`
	TopP             *float64  `json:"top_p,omitempty"`
	Stream           *bool     `json:"stream,omitempty"`
	Stop             []string  `json:"stop,omitempty"`
	PresencePenalty  *float64  `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float64  `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int `json:"logit_bias,omitempty"`
	User             *string   `json:"user,omitempty"`
	Tools            []Tool    `json:"tools,omitempty"`
	ToolChoice       interface{} `json:"tool_choice,omitempty"`
}

// Tool represents a tool that can be called
type Tool struct {
	Type     string       `json:"type"`
	Function ToolSpec     `json:"function"`
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
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	Delta        *Message `json:"delta,omitempty"`
	FinishReason *string `json:"finish_reason"`
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
	ID       string `json:"id"`
	Provider ProviderName `json:"provider"`
	Name     string `json:"name"`
	MaxTokens int   `json:"max_tokens"`
}

// Common model constants for each provider
const (
	// OpenAI Models
	ModelGPT4o           = "gpt-4o"
	ModelGPT4oMini       = "gpt-4o-mini"
	ModelGPT4Turbo       = "gpt-4-turbo"
	ModelGPT35Turbo      = "gpt-3.5-turbo"
	
	// Claude Models
	ModelClaude3Opus     = "claude-3-opus-20240229"
	ModelClaude3Sonnet   = "claude-3-sonnet-20240229"
	ModelClaude3Haiku    = "claude-3-haiku-20240307"
	ModelClaudeSonnet4   = "claude-sonnet-4-20250514"
	
	// Bedrock Models (these would be the actual Bedrock model IDs)
	ModelBedrockClaude3Opus   = "anthropic.claude-3-opus-20240229-v1:0"
	ModelBedrockClaude3Sonnet = "anthropic.claude-3-sonnet-20240229-v1:0"
	ModelBedrockTitan         = "amazon.titan-text-express-v1"
)

// GetModelInfo returns model information
func GetModelInfo(modelID string) *ModelInfo {
	modelMap := map[string]ModelInfo{
		ModelGPT4o: {
			ID: ModelGPT4o,
			Provider: ProviderNameOpenAI,
			Name: "GPT-4o",
			MaxTokens: 128000,
		},
		ModelClaude3Opus: {
			ID: ModelClaude3Opus,
			Provider: ProviderNameAnthropic,
			Name: "Claude 3 Opus",
			MaxTokens: 200000,
		},
		ModelBedrockClaude3Sonnet: {
			ID: ModelBedrockClaude3Sonnet,
			Provider: ProviderNameBedrock,
			Name: "Claude 3 Sonnet (Bedrock)",
			MaxTokens: 200000,
		},
	}
	
	if info, exists := modelMap[modelID]; exists {
		return &info
	}
	return nil
}