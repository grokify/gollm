package main

import (
	"context"
	"fmt"
	"log"

	"github.com/grokify/gollm"
)

func main() {
	// Create a client for Ollama
	// Default BaseURL is http://localhost:11434, but you can customize it
	client, err := gollm.NewClient(gollm.ClientConfig{
		Provider: gollm.ProviderNameOllama,
		BaseURL:  "http://localhost:11434", // Optional - this is the default
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	fmt.Println("Testing Ollama provider with GoLLM...")
	fmt.Println("Make sure you have Ollama running locally with a model installed.")
	fmt.Println("Example: ollama run llama3:8b")
	fmt.Println()

	// Create a chat completion request
	response, err := client.CreateChatCompletion(context.Background(), &gollm.ChatCompletionRequest{
		Model: gollm.ModelOllamaLlama3_8B, // You can use any model you have installed
		Messages: []gollm.Message{
			{
				Role:    gollm.RoleUser,
				Content: "Hello! Can you tell me a short fact about Apple Silicon MacBooks?",
			},
		},
		MaxTokens:   &[]int{100}[0],
		Temperature: &[]float64{0.7}[0],
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response: %s\n", response.Choices[0].Message.Content)
	fmt.Printf("Model used: %s\n", response.Model)
	fmt.Printf("Tokens used: %d (prompt: %d, completion: %d)\n",
		response.Usage.TotalTokens,
		response.Usage.PromptTokens,
		response.Usage.CompletionTokens)
}
