package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/grokify/fluxllm"
)

// ProviderDemo holds configuration for demonstrating a specific provider
type ProviderDemo struct {
	Name     string
	Config   fluxllm.ClientConfig
	Model    string
	Messages []fluxllm.Message
}

func main() {
	// Define all provider demonstrations
	demos := []ProviderDemo{
		{
			Name: "OpenAI",
			Config: fluxllm.ClientConfig{
				Provider: fluxllm.ProviderNameOpenAI,
				APIKey:   os.Getenv("OPENAI_API_KEY"),
			},
			Model: fluxllm.ModelGPT4o,
			Messages: []fluxllm.Message{
				{
					Role:    fluxllm.RoleUser,
					Content: "Hello! Can you explain what a unified LLM SDK is?",
				},
			},
		},
		{
			Name: "Anthropic",
			Config: fluxllm.ClientConfig{
				Provider: fluxllm.ProviderNameAnthropic,
				APIKey:   os.Getenv("ANTHROPIC_API_KEY"),
			},
			Model: fluxllm.ModelClaude3Sonnet,
			Messages: []fluxllm.Message{
				{
					Role:    fluxllm.RoleSystem,
					Content: "You are a helpful assistant that explains technical concepts clearly.",
				},
				{
					Role:    fluxllm.RoleUser,
					Content: "What are the benefits of using a unified SDK for multiple LLM providers?",
				},
			},
		},
		{
			Name: "AWS Bedrock",
			Config: fluxllm.ClientConfig{
				Provider: fluxllm.ProviderNameBedrock,
				Region:   "us-east-1",
			},
			Model: fluxllm.ModelBedrockClaude3Sonnet,
			Messages: []fluxllm.Message{
				{
					Role:    fluxllm.RoleUser,
					Content: "Explain the advantages of using AWS Bedrock for LLM deployments.",
				},
			},
		},
		{
			Name: "Ollama (Local)",
			Config: fluxllm.ClientConfig{
				Provider: fluxllm.ProviderNameOllama,
				BaseURL:  "http://localhost:11434",
			},
			Model: "llama3", // Use the model name as it appears in "ollama list"
			Messages: []fluxllm.Message{
				{
					Role:    fluxllm.RoleUser,
					Content: "What are the benefits of running LLMs locally with Ollama?",
				},
			},
		},
	}

	// Run all demonstrations
	for _, demo := range demos {
		fmt.Printf("=== %s Example ===\n", demo.Name)
		if err := demonstrateProvider(demo); err != nil {
			log.Printf("%s error: %v", demo.Name, err)
		}
		fmt.Println()
	}
}

// demonstrateProvider is a generic function that works with any provider
func demonstrateProvider(demo ProviderDemo) error {
	client, err := fluxllm.NewClient(demo.Config)
	if err != nil {
		return err
	}
	defer client.Close()

	response, err := client.CreateChatCompletion(context.Background(), &fluxllm.ChatCompletionRequest{
		Model:       demo.Model,
		Messages:    demo.Messages,
		MaxTokens:   intPtr(150),
		Temperature: float64Ptr(0.7),
	})
	if err != nil {
		return err
	}

	fmt.Printf("Response: %s\n", response.Choices[0].Message.Content)
	fmt.Printf("Tokens used: %d\n", response.Usage.TotalTokens)

	return nil
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}
