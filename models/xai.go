package models

// X.AI Grok Model Documentation
const (
	// XAIModelsURL is the official X.AI models documentation page.
	// Use this to check for new models, deprecations, and model updates.
	XAIModelsURL = "https://docs.x.ai/docs/models"

	// XAIapiurl is the X.AI API reference page.
	XAIAPIURL = "https://docs.x.ai/docs"
)

// Grok 4.1 Family (Latest - November 2025)
const (
	// Grok4_1FastReasoning is the best tool-calling model with 2M context window.
	// Optimized for maximum intelligence and agentic tool calling.
	Grok4_1FastReasoning = "grok-4-1-fast-reasoning"

	// Grok4_1FastNonReasoning provides instant responses with 2M context window.
	// Optimized for speed without reasoning overhead.
	Grok4_1FastNonReasoning = "grok-4-1-fast-non-reasoning"
)

// Grok 4 Family (July 2025)
const (
	// Grok4_0709 is the flagship Grok 4 model with 256K context window.
	// Released July 9, 2025. Provides high-quality reasoning.
	Grok4_0709 = "grok-4-0709"

	// Grok4FastReasoning provides fast reasoning with 2M context window.
	// Performance on par with grok-4-0709 with larger context.
	Grok4FastReasoning = "grok-4-fast-reasoning"

	// Grok4FastNonReasoning provides fast non-reasoning with 2M context window.
	// Optimized for speed without reasoning overhead.
	Grok4FastNonReasoning = "grok-4-fast-non-reasoning"

	// GrokCodeFast1 is optimized for agentic coding with 256K context window.
	// Speedy and economical reasoning model for coding tasks.
	GrokCodeFast1 = "grok-code-fast-1"
)

// Grok 3 Family
const (
	Grok3     = "grok-3"      // Grok 3
	Grok3Mini = "grok-3-mini" // Grok 3 Mini (smaller, faster)
)

// Grok 2 Family
const (
	Grok2_1212   = "grok-2-1212"        // Grok 2 (December 2024)
	Grok2_Vision = "grok-2-vision-1212" // Grok 2 with vision capabilities
)

// Deprecated Grok Models
const (
	// GrokBeta is deprecated. Use Grok3 or Grok4_1FastReasoning instead.
	GrokBeta = "grok-beta" // Deprecated: use grok-3 or grok-4

	// GrokVision is deprecated. Use Grok2_Vision instead.
	GrokVision = "grok-vision-beta" // Deprecated
)
