package main

import (
	"fmt"
	"log"

	"github.com/grokify/gollm"
)

func main() {
	fmt.Println("=== GoLLM Architecture Demo ===")

	// Create an OpenAI client
	client, err := gollm.NewClient(gollm.ClientConfig{
		Provider: gollm.ProviderNameOpenAI,
		APIKey:   "demo-key", // This won't work without a real key, but shows the structure
	})
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return
	}
	defer client.Close()

	// Show provider name
	fmt.Printf("Created client with provider: %s\n", client.Provider().Name())

	// Show how the interface works (this will fail without real credentials)
	req := &gollm.ChatCompletionRequest{
		Model: gollm.ModelGPT4o,
		Messages: []gollm.Message{
			{
				Role:    gollm.RoleUser,
				Content: "Hello, world!",
			},
		},
		MaxTokens: &[]int{50}[0],
	}

	fmt.Printf("Request structure: %+v\n", req)
	fmt.Println("\nThe architecture now has:")
	fmt.Println("1. A unified Provider interface")
	fmt.Println("2. Provider-specific implementations (OpenAI, Anthropic, Bedrock)")
	fmt.Println("3. A ChatClient wrapper that provides a consistent API")
	fmt.Println("4. All core types in the main gollm package")

	// Demonstrate model info functionality
	if info := gollm.GetModelInfo(gollm.ModelGPT4o); info != nil {
		fmt.Printf("\nModel info for %s:\n", info.ID)
		fmt.Printf("  Provider: %s\n", info.Provider)
		fmt.Printf("  Name: %s\n", info.Name)
		fmt.Printf("  Max Tokens: %d\n", info.MaxTokens)
	}
}
