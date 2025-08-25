package main

import (
	"fmt"
	"log"

	"github.com/grokify/gollm"
)

func main() {
	fmt.Println("=== GoLLM Provider Architecture Demo ===")
	fmt.Println()
	
	// Show the current architecture
	fmt.Println("Current Architecture:")
	fmt.Println("📁 gollm/ (main package)")
	fmt.Println("  ├── client.go        - ChatClient wrapper")
	fmt.Println("  ├── provider.go      - Provider interface")
	fmt.Println("  ├── providers.go     - Provider adapters")
	fmt.Println("  ├── types.go         - Unified types")
	fmt.Println("  ├── errors.go        - Error handling")
	fmt.Println("  └── providers/")
	fmt.Println("      ├── openai/      - OpenAI implementation")
	fmt.Println("      ├── anthropic/   - Claude implementation")
	fmt.Println("      └── bedrock/     - AWS Bedrock implementation")
	fmt.Println()
	
	// Demonstrate creating clients for different providers
	fmt.Println("Creating clients for different providers...")
	
	// OpenAI client (won't work without real API key)
	openaiClient, err := gollm.NewClient(gollm.ClientConfig{
		Provider: gollm.ProviderNameOpenAI,
		APIKey:   "demo-openai-key",
	})
	if err != nil {
		log.Printf("Failed to create OpenAI client: %v", err)
	} else {
		fmt.Printf("✅ OpenAI client created: %s\n", openaiClient.Provider().Name())
		openaiClient.Close()
	}
	
	// Anthropic client (won't work without real API key)
	anthropicClient, err := gollm.NewClient(gollm.ClientConfig{
		Provider: gollm.ProviderNameAnthropic,
		APIKey:   "demo-anthropic-key",
	})
	if err != nil {
		log.Printf("Failed to create Anthropic client: %v", err)
	} else {
		fmt.Printf("✅ Anthropic client created: %s\n", anthropicClient.Provider().Name())
		anthropicClient.Close()
	}
	
	// Bedrock client (requires AWS credentials)
	bedrockClient, err := gollm.NewClient(gollm.ClientConfig{
		Provider: gollm.ProviderNameBedrock,
		Region:   "us-east-1",
	})
	if err != nil {
		log.Printf("⚠️  Bedrock client creation failed (expected without AWS creds): %v", err)
	} else {
		fmt.Printf("✅ Bedrock client created: %s\n", bedrockClient.Provider().Name())
		bedrockClient.Close()
	}
	
	fmt.Println()
	fmt.Println("Benefits of this architecture:")
	fmt.Println("1. 🔌 Pluggable: Easy to add new LLM providers")
	fmt.Println("2. 🎯 Unified: Same API for all providers") 
	fmt.Println("3. 🧪 Testable: Provider interface can be mocked")
	fmt.Println("4. 📦 Modular: Each provider is self-contained")
	fmt.Println("5. 🔧 Maintainable: Clear separation of concerns")
	
	// Show example request structure
	fmt.Println()
	fmt.Println("Example unified request structure:")
	req := &gollm.ChatCompletionRequest{
		Model: gollm.ModelGPT4o,
		Messages: []gollm.Message{
			{
				Role:    gollm.RoleSystem,
				Content: "You are a helpful assistant.",
			},
			{
				Role:    gollm.RoleUser,
				Content: "Hello, world!",
			},
		},
		MaxTokens:   &[]int{100}[0],
		Temperature: &[]float64{0.7}[0],
	}
	
	fmt.Printf("Request: %+v\n", req)
	fmt.Println()
	fmt.Println("This same request structure works with all providers!")
}