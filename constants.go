package gollm

const (
	EnvVarAnthropicAPIKey = "ANTHROPIC_API_KEY" // #nosec G101
	EnvVarOpenAIAPIKey    = "OPENAI_API_KEY"    // #nosec G101
	EnvVarGeminiAPIKey    = "GEMINI_API_KEY"    // #nosec G101
	EnvVarXAIAPIKey       = "XAI_API_KEY"       // #nosec G101
)

// ProviderName represents the different LLM provider names
type ProviderName string

const (
	ProviderNameOpenAI    ProviderName = "openai"
	ProviderNameAnthropic ProviderName = "anthropic"
	ProviderNameBedrock   ProviderName = "bedrock"
	ProviderNameOllama    ProviderName = "ollama"
	ProviderNameGemini    ProviderName = "gemini"
	ProviderNameXAI       ProviderName = "xai"
)

// Common model constants for each provider
const (
	// Bedrock Models (these would be the actual Bedrock model IDs)
	ModelBedrockClaude3Opus   = "anthropic.claude-3-opus-20240229-v1:0"
	ModelBedrockClaude3Sonnet = "anthropic.claude-3-sonnet-20240229-v1:0"
	ModelBedrockClaudeOpus4   = "anthropic.claude-opus-4-20250514-v1:0"
	ModelBedrockTitan         = "amazon.titan-text-express-v1"

	// Claude Models
	ModelClaudeOpus4_1   = "claude-opus-4-1-20250805"
	ModelClaudeOpus4     = "claude-opus-4-20250514"
	ModelClaudeSonnet4   = "claude-sonnet-4-20250514"
	ModelClaude3_7Sonnet = "claude-3-7-sonnet-20250219"
	ModelClaude3_5Haiku  = "claude-3-5-haiku-20241022"
	ModelClaude3Opus     = "claude-3-opus-20240229"
	ModelClaude3Sonnet   = "claude-3-sonnet-20240229"
	ModelClaude3Haiku    = "claude-3-haiku-20240307"

	// Gemini Models
	ModelGemini2_5Pro       = "gemini-2.5-pro"        // Stable, advanced reasoning
	ModelGemini2_5Flash     = "gemini-2.5-flash"      // Stable, balanced performance
	ModelGeminiLive2_5Flash = "gemini-live-2.5-flash" // Stable Live API (private GA)
	ModelGemini1_5Pro       = "gemini-1.5-pro"
	ModelGemini1_5Flash     = "gemini-1.5-flash"
	ModelGeminiPro          = "gemini-pro"

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

	// OpenAI Models
	ModelGPT5           = "gpt-5"
	ModelGPT5Mini       = "gpt-5-mini"
	ModelGPT5Nano       = "gpt-5-nano"
	ModelGPT5ChatLatest = "gpt-5-chat-latest"
	ModelGPT4_1         = "gpt-4.1"
	ModelGPT4_1Mini     = "gpt-4.1-mini"
	ModelGPT4_1Nano     = "gpt-4.1-nano"
	ModelGPT4o          = "gpt-4o"
	ModelGPT4oMini      = "gpt-4o-mini"
	ModelGPT4Turbo      = "gpt-4-turbo"
	ModelGPT35Turbo     = "gpt-3.5-turbo"

	// Vertex AI Models
	ModelVertexClaudeOpus4 = "claude-opus-4@20250514"

	// X.AI Grok Models
	ModelGrok3          = "grok-3"              // Latest Grok model
	ModelGrok3Mini      = "grok-3-mini"         // Smaller, faster Grok model
	ModelGrok2_1212     = "grok-2-1212"         // Grok 2 (December 2024)
	ModelGrok2_Vision   = "grok-2-vision-1212"  // Grok 2 with vision
	ModelGrokBeta       = "grok-beta"           // Deprecated: use grok-3
	ModelGrokVision     = "grok-vision-beta"    // Deprecated
)
