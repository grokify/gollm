package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/grokify/fluxllm"
)

func main() {
	// Interactive conversation example
	fmt.Println("=== Interactive Conversation Example ===")
	fmt.Println("This example demonstrates maintaining conversation context across multiple providers.")
	fmt.Println("Type 'quit' to exit, 'switch' to change provider")
	fmt.Println()

	if err := runConversation(); err != nil {
		log.Fatal(err)
	}
}

func runConversation() error {
	scanner := bufio.NewScanner(os.Stdin)
	messages := []fluxllm.Message{
		{
			Role:    fluxllm.RoleSystem,
			Content: "You are a helpful assistant. Keep your responses concise and friendly.",
		},
	}

	currentProvider := fluxllm.ProviderNameOpenAI
	client, err := createClient(currentProvider)
	if err != nil {
		return err
	}
	defer client.Close()

	fmt.Printf("Current provider: %s\n", currentProvider)
	fmt.Print("You: ")

	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			fmt.Print("You: ")
			continue
		}

		if input == "quit" {
			break
		}

		if input == "switch" {
			// Switch to next provider
			client.Close()
			currentProvider = getNextProvider(currentProvider)
			client, err = createClient(currentProvider)
			if err != nil {
				log.Printf("Failed to switch provider: %v", err)
				fmt.Print("You: ")
				continue
			}
			fmt.Printf("\nSwitched to provider: %s\n", currentProvider)
			fmt.Print("You: ")
			continue
		}

		// Add user message
		messages = append(messages, fluxllm.Message{
			Role:    fluxllm.RoleUser,
			Content: input,
		})

		// Get response
		response, err := client.CreateChatCompletion(context.Background(), &fluxllm.ChatCompletionRequest{
			Model:       getModelForProvider(currentProvider),
			Messages:    messages,
			MaxTokens:   intPtr(150),
			Temperature: float64Ptr(0.7),
		})
		if err != nil {
			log.Printf("Error: %v", err)
			fmt.Print("You: ")
			continue
		}

		assistantMessage := response.Choices[0].Message.Content
		fmt.Printf("Assistant (%s): %s\n", currentProvider, assistantMessage)

		// Add assistant response to conversation
		messages = append(messages, fluxllm.Message{
			Role:    fluxllm.RoleAssistant,
			Content: assistantMessage,
		})

		// Keep conversation history manageable (last 10 messages + system message)
		if len(messages) > 11 {
			messages = append(messages[:1], messages[len(messages)-10:]...)
		}

		fmt.Print("You: ")
	}

	return nil
}

func createClient(provider fluxllm.ProviderName) (*fluxllm.ChatClient, error) {
	config := fluxllm.ClientConfig{Provider: provider}

	switch provider {
	case fluxllm.ProviderNameOpenAI:
		config.APIKey = os.Getenv("OPENAI_API_KEY")
	case fluxllm.ProviderNameAnthropic:
		config.APIKey = os.Getenv("ANTHROPIC_API_KEY")
	case fluxllm.ProviderNameBedrock:
		config.Region = "us-east-1"
	}

	return fluxllm.NewClient(config)
}

func getNextProvider(current fluxllm.ProviderName) fluxllm.ProviderName {
	switch current {
	case fluxllm.ProviderNameOpenAI:
		return fluxllm.ProviderNameAnthropic
	case fluxllm.ProviderNameAnthropic:
		return fluxllm.ProviderNameBedrock
	case fluxllm.ProviderNameBedrock:
		return fluxllm.ProviderNameOpenAI
	default:
		return fluxllm.ProviderNameOpenAI
	}
}

func getModelForProvider(provider fluxllm.ProviderName) string {
	switch provider {
	case fluxllm.ProviderNameOpenAI:
		return fluxllm.ModelGPT4oMini
	case fluxllm.ProviderNameAnthropic:
		return fluxllm.ModelClaude3Haiku
	case fluxllm.ProviderNameBedrock:
		return fluxllm.ModelBedrockClaude3Sonnet
	default:
		return fluxllm.ModelGPT4oMini
	}
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}
