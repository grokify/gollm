package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/grokify/gollm"
)

func main() {
	// Example 1: OpenAI
	fmt.Println("=== OpenAI Example ===")
	if err := demonstrateOpenAI(); err != nil {
		log.Printf("OpenAI error: %v", err)
	}

	// Example 2: Anthropic (Claude)
	fmt.Println("\n=== Anthropic Example ===")
	if err := demonstrateAnthropic(); err != nil {
		log.Printf("Anthropic error: %v", err)
	}

	// Example 3: AWS Bedrock
	fmt.Println("\n=== AWS Bedrock Example ===")
	if err := demonstrateBedrock(); err != nil {
		log.Printf("Bedrock error: %v", err)
	}
}

func demonstrateOpenAI() error {
	client, err := gollm.NewClient(gollm.ClientConfig{
		Provider: gollm.ProviderNameOpenAI,
		APIKey:   os.Getenv("OPENAI_API_KEY"),
	})
	if err != nil {
		return err
	}
	defer client.Close()

	response, err := client.CreateChatCompletion(context.Background(), &gollm.ChatCompletionRequest{
		Model: gollm.ModelGPT4o,
		Messages: []gollm.Message{
			{
				Role:    gollm.RoleUser,
				Content: "Hello! Can you explain what a unified LLM SDK is?",
			},
		},
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

func demonstrateAnthropic() error {
	client, err := gollm.NewClient(gollm.ClientConfig{
		Provider: gollm.ProviderNameAnthropic,
		APIKey:   os.Getenv("ANTHROPIC_API_KEY"),
	})
	if err != nil {
		return err
	}
	defer client.Close()

	response, err := client.CreateChatCompletion(context.Background(), &gollm.ChatCompletionRequest{
		Model: gollm.ModelClaude3Sonnet,
		Messages: []gollm.Message{
			{
				Role:    gollm.RoleSystem,
				Content: "You are a helpful assistant that explains technical concepts clearly.",
			},
			{
				Role:    gollm.RoleUser,
				Content: "What are the benefits of using a unified SDK for multiple LLM providers?",
			},
		},
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

func demonstrateBedrock() error {
	client, err := gollm.NewClient(gollm.ClientConfig{
		Provider: gollm.ProviderNameBedrock,
		Region:   "us-east-1", // AWS region
	})
	if err != nil {
		return err
	}
	defer client.Close()

	response, err := client.CreateChatCompletion(context.Background(), &gollm.ChatCompletionRequest{
		Model: gollm.ModelBedrockClaude3Sonnet,
		Messages: []gollm.Message{
			{
				Role:    gollm.RoleUser,
				Content: "Explain the advantages of using AWS Bedrock for LLM deployments.",
			},
		},
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
