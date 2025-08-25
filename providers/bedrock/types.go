package bedrock

// Request represents a Bedrock API request
type Request struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   *int      `json:"max_tokens,omitempty"`
	Temperature *float64  `json:"temperature,omitempty"`
}

// Message represents a message in Bedrock format
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Response represents a Bedrock API response
type Response struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Model   string `json:"model"`
	Usage   Usage  `json:"usage"`
}

// Usage represents token usage in Bedrock response
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}
