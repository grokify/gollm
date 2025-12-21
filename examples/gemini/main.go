package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/grokify/fluxllm"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is required")
	}

	// Create a Gemini client
	client, err := fluxllm.NewClient(fluxllm.ClientConfig{
		Provider: fluxllm.ProviderNameGemini,
		APIKey:   apiKey,
	})
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}
	defer client.Close()

	// Create a chat completion request
	response, err := client.CreateChatCompletion(context.Background(), &fluxllm.ChatCompletionRequest{
		Model: fluxllm.ModelGemini1_5Flash,
		Messages: []fluxllm.Message{
			{
				Role:    fluxllm.RoleUser,
				Content: "Hello! Can you tell me a short joke?",
			},
		},
		MaxTokens:   &[]int{150}[0],
		Temperature: &[]float64{0.7}[0],
	})
	if err != nil {
		log.Fatal("Failed to create completion:", err)
	}

	// Print the response
	fmt.Printf("Response: %s\n", response.Choices[0].Message.Content)
	fmt.Printf("Tokens used: %d\n", response.Usage.TotalTokens)
}
