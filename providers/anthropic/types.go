package anthropic

// Request represents an Anthropic API request
type Request struct {
	Model       string    `json:"model"`
	MaxTokens   int       `json:"max_tokens"`
	Messages    []Message `json:"messages"`
	System      string    `json:"system,omitempty"`
	Temperature *float64  `json:"temperature,omitempty"`
	TopP        *float64  `json:"top_p,omitempty"`
}

// Message represents a message in Anthropic format
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Response represents an Anthropic API response
type Response struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	Role       string    `json:"role"`
	Content    []Content `json:"content"`
	Model      string    `json:"model"`
	StopReason string    `json:"stop_reason"`
	Usage      Usage     `json:"usage"`
}

// Content represents content in Anthropic response
type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Usage represents token usage in Anthropic response
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}
