package main

import (
	"context"
	"fmt"
	"log"

	"github.com/grokify/gollm"
)

// Example of a 3rd party provider implementation
// This could be in an external package like github.com/someone/gollm-custom

// customProvider implements the gollm.Provider interface
type customProvider struct {
	name   string
	apiKey string
}

// NewCustomProvider creates a new custom provider (this would be in external package)
func NewCustomProvider(name, apiKey string) gollm.Provider {
	return &customProvider{
		name:   name,
		apiKey: apiKey,
	}
}

func (p *customProvider) Name() string {
	return p.name
}

func (p *customProvider) CreateChatCompletion(ctx context.Context, req *gollm.ChatCompletionRequest) (*gollm.ChatCompletionResponse, error) {
	// Mock implementation for demonstration
	return &gollm.ChatCompletionResponse{
		ID:      "custom-123",
		Object:  "chat.completion",
		Created: 1234567890,
		Model:   req.Model,
		Choices: []gollm.ChatCompletionChoice{
			{
				Index: 0,
				Message: gollm.Message{
					Role:    gollm.RoleAssistant,
					Content: fmt.Sprintf("Hello from %s! You asked: %s", p.name, req.Messages[len(req.Messages)-1].Content),
				},
				FinishReason: &[]string{"stop"}[0],
			},
		},
		Usage: gollm.Usage{
			PromptTokens:     10,
			CompletionTokens: 20,
			TotalTokens:      30,
		},
	}, nil
}

func (p *customProvider) CreateChatCompletionStream(ctx context.Context, req *gollm.ChatCompletionRequest) (gollm.ChatCompletionStream, error) {
	return nil, fmt.Errorf("streaming not implemented in custom provider demo")
}

func (p *customProvider) Close() error {
	return nil
}

func main() {
	fmt.Println("=== 3rd Party Custom Provider Example ===")

	// Create a custom provider (this could be from an external package)
	customProv := NewCustomProvider("MyCustomLLM", "custom-api-key")

	// Inject the custom provider directly into gollm
	client, err := gollm.NewClient(gollm.ClientConfig{
		CustomProvider: customProv, // Direct provider injection!
		// Note: Provider field is ignored when CustomProvider is set
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	fmt.Printf("Using custom provider: %s\n", client.Provider().Name())

	// Use the same gollm API with the custom provider
	response, err := client.CreateChatCompletion(context.Background(), &gollm.ChatCompletionRequest{
		Model: "custom-model-v1",
		Messages: []gollm.Message{
			{
				Role:    gollm.RoleUser,
				Content: "Hello from a 3rd party provider!",
			},
		},
		MaxTokens:   &[]int{50}[0],
		Temperature: &[]float64{0.7}[0],
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response: %s\n", response.Choices[0].Message.Content)
	fmt.Printf("Tokens used: %d\n", response.Usage.TotalTokens)

	fmt.Println()
	fmt.Println("🎉 3rd party providers can now extend gollm without modifying core!")
	fmt.Println("   - No registry needed")
	fmt.Println("   - No global state")
	fmt.Println("   - Compile-time type safety")
	fmt.Println("   - Clean dependency injection")
}
