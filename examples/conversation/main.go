package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/grokify/gollm"
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
	messages := []gollm.Message{
		{
			Role:    gollm.RoleSystem,
			Content: "You are a helpful assistant. Keep your responses concise and friendly.",
		},
	}

	currentProvider := gollm.ProviderNameOpenAI
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
		messages = append(messages, gollm.Message{
			Role:    gollm.RoleUser,
			Content: input,
		})

		// Get response
		response, err := client.CreateChatCompletion(context.Background(), &gollm.ChatCompletionRequest{
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
		messages = append(messages, gollm.Message{
			Role:    gollm.RoleAssistant,
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

func createClient(provider gollm.ProviderName) (*gollm.ChatClient, error) {
	config := gollm.ClientConfig{Provider: provider}
	
	switch provider {
	case gollm.ProviderNameOpenAI:
		config.APIKey = os.Getenv("OPENAI_API_KEY")
	case gollm.ProviderNameAnthropic:
		config.APIKey = os.Getenv("ANTHROPIC_API_KEY")
	case gollm.ProviderNameBedrock:
		config.Region = "us-east-1"
	}
	
	return gollm.NewClient(config)
}

func getNextProvider(current gollm.ProviderName) gollm.ProviderName {
	switch current {
	case gollm.ProviderNameOpenAI:
		return gollm.ProviderNameAnthropic
	case gollm.ProviderNameAnthropic:
		return gollm.ProviderNameBedrock
	case gollm.ProviderNameBedrock:
		return gollm.ProviderNameOpenAI
	default:
		return gollm.ProviderNameOpenAI
	}
}

func getModelForProvider(provider gollm.ProviderName) string {
	switch provider {
	case gollm.ProviderNameOpenAI:
		return gollm.ModelGPT4oMini
	case gollm.ProviderNameAnthropic:
		return gollm.ModelClaude3Haiku
	case gollm.ProviderNameBedrock:
		return gollm.ModelBedrockClaude3Sonnet
	default:
		return gollm.ModelGPT4oMini
	}
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}