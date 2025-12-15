package models

// Google Gemini Model Documentation
const (
	// GeminiModelsURL is the official Google Gemini models documentation page.
	// Use this to check for new models, deprecations, and model updates.
	GeminiModelsURL = "https://ai.google.dev/gemini-api/docs/models/gemini"

	// GeminiAPIURL is the Google Gemini API reference page.
	GeminiAPIURL = "https://ai.google.dev/gemini-api/docs"
)

// Gemini 2.5 Family (Latest)
const (
	// Gemini2_5Pro is stable with advanced reasoning capabilities.
	Gemini2_5Pro = "gemini-2.5-pro"

	// Gemini2_5Flash is stable with balanced performance.
	Gemini2_5Flash = "gemini-2.5-flash"

	// GeminiLive2_5Flash is the stable Live API model (private GA).
	GeminiLive2_5Flash = "gemini-live-2.5-flash"
)

// Gemini 1.5 Family
const (
	Gemini1_5Pro   = "gemini-1.5-pro"   // Gemini 1.5 Pro
	Gemini1_5Flash = "gemini-1.5-flash" // Gemini 1.5 Flash
)

// Legacy Gemini Models
const (
	GeminiPro = "gemini-pro" // Legacy Gemini Pro (use Gemini 2.5 instead)
)
