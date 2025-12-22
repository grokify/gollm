package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/grokify/metallm"
)

// ProviderDemo holds configuration for demonstrating a specific provider
type ProviderDemo struct {
	Name     string
	Config   metallm.ClientConfig
	Model    string
	Messages []metallm.Message
}

func main() {
	// Define all provider demonstrations
	demos := []ProviderDemo{
		{
			Name: "OpenAI",
			Config: metallm.ClientConfig{
				Provider: metallm.ProviderNameOpenAI,
				APIKey:   os.Getenv("OPENAI_API_KEY"),
			},
			Model: metallm.ModelGPT4o,
			Messages: []metallm.Message{
				{
					Role:    metallm.RoleUser,
					Content: "Hello! Can you explain what a unified LLM SDK is?",
				},
			},
		},
		{
			Name: "Anthropic",
			Config: metallm.ClientConfig{
				Provider: metallm.ProviderNameAnthropic,
				APIKey:   os.Getenv("ANTHROPIC_API_KEY"),
			},
			Model: metallm.ModelClaude3Sonnet,
			Messages: []metallm.Message{
				{
					Role:    metallm.RoleSystem,
					Content: "You are a helpful assistant that explains technical concepts clearly.",
				},
				{
					Role:    metallm.RoleUser,
					Content: "What are the benefits of using a unified SDK for multiple LLM providers?",
				},
			},
		},
		{
			Name: "AWS Bedrock",
			Config: metallm.ClientConfig{
				Provider: metallm.ProviderNameBedrock,
				Region:   "us-east-1",
			},
			Model: metallm.ModelBedrockClaude3Sonnet,
			Messages: []metallm.Message{
				{
					Role:    metallm.RoleUser,
					Content: "Explain the advantages of using AWS Bedrock for LLM deployments.",
				},
			},
		},
		{
			Name: "Ollama (Local)",
			Config: metallm.ClientConfig{
				Provider: metallm.ProviderNameOllama,
				BaseURL:  "http://localhost:11434",
			},
			Model: "llama3", // Use the model name as it appears in "ollama list"
			Messages: []metallm.Message{
				{
					Role:    metallm.RoleUser,
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
	client, err := metallm.NewClient(demo.Config)
	if err != nil {
		return err
	}
	defer client.Close()

	response, err := client.CreateChatCompletion(context.Background(), &metallm.ChatCompletionRequest{
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
