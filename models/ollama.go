package models

// Ollama Model Documentation
const (
	// OllamaModelsURL is the official Ollama models library page.
	// Use this to check for new models and model updates.
	OllamaModelsURL = "https://ollama.com/library"

	// OllamaAPIURL is the Ollama API reference page.
	OllamaAPIURL = "https://github.com/ollama/ollama/blob/main/docs/api.md"
)

// Ollama Llama Models
const (
	OllamaLlama3_8B  = "llama3:8b"  // Llama 3 8B
	OllamaLlama3_70B = "llama3:70b" // Llama 3 70B
)

// Ollama Mistral Models
const (
	OllamaMistral7B   = "mistral:7b"   // Mistral 7B
	OllamaMixtral8x7B = "mixtral:8x7b" // Mixtral 8x7B
)

// Ollama Code Models
const (
	OllamaCodeLlama = "codellama:13b"       // CodeLlama 13B
	OllamaDeepSeek  = "deepseek-coder:6.7b" // DeepSeek Coder 6.7B
)

// Ollama Gemma Models
const (
	OllamaGemma2B = "gemma:2b" // Gemma 2B
	OllamaGemma7B = "gemma:7b" // Gemma 7B
)

// Ollama Other Models
const (
	OllamaQwen2_5 = "qwen2.5:7b" // Qwen 2.5 7B
)
