package gollm

const (
	EnvVarAnthropicAPIKey = "ANTHROPIC_API_KEY" // #nosec G101
	EnvVarOpenAIAPIKey    = "OPENAI_API_KEY"    // #nosec G101
)

// Common model constants for each provider
const (
	// OpenAI Models
	ModelGPT4o      = "gpt-4o"
	ModelGPT4oMini  = "gpt-4o-mini"
	ModelGPT4Turbo  = "gpt-4-turbo"
	ModelGPT35Turbo = "gpt-3.5-turbo"

	// Claude Models
	ModelClaude3Opus   = "claude-3-opus-20240229"
	ModelClaude3Sonnet = "claude-3-sonnet-20240229"
	ModelClaude3Haiku  = "claude-3-haiku-20240307"
	ModelClaudeSonnet4 = "claude-sonnet-4-20250514"
	ModelClaudeOpus4   = "claude-opus-4-20250514"

	// Bedrock Models (these would be the actual Bedrock model IDs)
	ModelBedrockClaude3Opus   = "anthropic.claude-3-opus-20240229-v1:0"
	ModelBedrockClaude3Sonnet = "anthropic.claude-3-sonnet-20240229-v1:0"
	ModelBedrockClaudeOpus4   = "anthropic.claude-opus-4-20250514-v1:0"
	ModelBedrockTitan         = "amazon.titan-text-express-v1"

	// Vertex AI Models
	ModelVertexClaudeOpus4 = "claude-opus-4@20250514"

	// Ollama Models (popular models that run well on Apple Silicon)
	ModelOllamaLlama3_8B   = "llama3:8b"
	ModelOllamaLlama3_70B  = "llama3:70b"
	ModelOllamaMistral7B   = "mistral:7b"
	ModelOllamaMixtral8x7B = "mixtral:8x7b"
	ModelOllamaCodeLlama   = "codellama:13b"
	ModelOllamaGemma2B     = "gemma:2b"
	ModelOllamaGemma7B     = "gemma:7b"
	ModelOllamaQwen2_5     = "qwen2.5:7b"
	ModelOllamaDeepSeek    = "deepseek-coder:6.7b"
)
