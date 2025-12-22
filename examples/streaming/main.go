package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/grokify/metallm"
)

func main() {
	// Example: Streaming with OpenAI
	fmt.Println("=== OpenAI Streaming Example ===")
	if err := demonstrateOpenAIStreaming(); err != nil {
		log.Printf("OpenAI streaming error: %v", err)
	}

	// Example: Streaming with Anthropic
	fmt.Println("\n=== Anthropic Streaming Example ===")
	if err := demonstrateAnthropicStreaming(); err != nil {
		log.Printf("Anthropic streaming error: %v", err)
	}
}

func demonstrateOpenAIStreaming() error {
	client, err := metallm.NewClient(metallm.ClientConfig{
		Provider: metallm.ProviderNameOpenAI,
		APIKey:   os.Getenv("OPENAI_API_KEY"),
	})
	if err != nil {
		return err
	}
	defer client.Close()

	stream, err := client.CreateChatCompletionStream(context.Background(), &metallm.ChatCompletionRequest{
		Model: metallm.ModelGPT4o,
		Messages: []metallm.Message{
			{
				Role:    metallm.RoleUser,
				Content: "Write a short story about a robot learning to paint. Keep it under 100 words.",
			},
		},
		MaxTokens:   intPtr(150),
		Temperature: float64Ptr(0.8),
	})
	if err != nil {
		return err
	}
	defer stream.Close()

	fmt.Print("OpenAI Response: ")
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
	}
	fmt.Println()

	return nil
}

func demonstrateAnthropicStreaming() error {
	client, err := metallm.NewClient(metallm.ClientConfig{
		Provider: metallm.ProviderNameAnthropic,
		APIKey:   os.Getenv("ANTHROPIC_API_KEY"),
	})
	if err != nil {
		return err
	}
	defer client.Close()

	stream, err := client.CreateChatCompletionStream(context.Background(), &metallm.ChatCompletionRequest{
		Model: metallm.ModelClaude3Haiku,
		Messages: []metallm.Message{
			{
				Role:    metallm.RoleSystem,
				Content: "You are a creative writing assistant.",
			},
			{
				Role:    metallm.RoleUser,
				Content: "Write a haiku about programming. Make it thoughtful and concise.",
			},
		},
		MaxTokens:   intPtr(100),
		Temperature: float64Ptr(0.9),
	})
	if err != nil {
		return err
	}
	defer stream.Close()

	fmt.Print("Claude Response: ")
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
	}
	fmt.Println()

	return nil
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}
