# GoLLM - Unified Go SDK for Large Language Models

[![Build Status][build-status-svg]][build-status-url]
[![Lint Status][lint-status-svg]][lint-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![License][license-svg]][license-url]

GoLLM is a unified Go SDK that provides a consistent interface for interacting with multiple Large Language Model (LLM) providers including OpenAI, Anthropic (Claude), Google Gemini, AWS Bedrock, and Ollama. It implements the Chat Completions API pattern and offers both synchronous and streaming capabilities.

## ✨ Features

- **🔌 Multi-Provider Support**: OpenAI, Anthropic (Claude), Google Gemini, AWS Bedrock, and Ollama
- **🎯 Unified API**: Same interface across all providers
- **📡 Streaming Support**: Real-time response streaming
- **🧠 Conversation Memory**: Persistent conversation history using Key-Value Stores
- **🧪 Testable**: Clean interfaces that can be easily mocked
- **🔧 Extensible**: Easy to add new LLM providers
- **📦 Modular**: Provider-specific implementations in separate packages
- **🏗️ Reference Architecture**: Internal providers serve as reference implementations for external providers
- **🔌 3rd Party Friendly**: External providers can be injected without modifying core library
- **⚡ Type Safe**: Full Go type safety with comprehensive error handling

## 🏗️ Architecture

GoLLM uses a clean, modular architecture that separates concerns and enables easy extensibility:

```
gollm/
├── client.go            # Main ChatClient wrapper
├── providers.go         # Factory functions for built-in providers
├── types.go             # Type aliases for backward compatibility
├── memory.go            # Conversation memory management
├── errors.go            # Unified error handling
├── provider/            # 🎯 Public interface package for external providers
│   ├── interface.go     # Provider interface that all providers must implement
│   └── types.go         # Unified request/response types
└── providers/           # 📦 Individual provider packages (reference implementations)
    ├── openai/          # OpenAI implementation
    │   ├── openai.go    # HTTP client
    │   ├── types.go     # OpenAI-specific types
    │   └── adapter.go   # provider.Provider implementation
    ├── anthropic/       # Anthropic implementation
    │   ├── anthropic.go # HTTP client
    │   ├── types.go     # Anthropic-specific types
    │   └── adapter.go   # provider.Provider implementation
    ├── gemini/          # Google Gemini implementation
    │   ├── gemini.go    # HTTP client
    │   ├── types.go     # Gemini-specific types
    │   └── adapter.go   # provider.Provider implementation
    ├── bedrock/         # AWS Bedrock implementation
    │   ├── bedrock.go   # AWS client
    │   ├── types.go     # Bedrock-specific types
    │   └── adapter.go   # provider.Provider implementation
    └── ollama/          # Ollama implementation
        ├── ollama.go    # HTTP client
        ├── types.go     # Ollama-specific types
        └── adapter.go   # provider.Provider implementation
```

### Key Architecture Benefits

- **🎯 Public Interface**: The `provider` package exports the `Provider` interface that external packages can implement
- **🏗️ Reference Implementation**: Internal providers follow the exact same structure that external providers should use
- **🔌 Direct Injection**: External providers are injected via `ClientConfig.CustomProvider` without modifying core code
- **📦 Modular Design**: Each provider is self-contained with its own HTTP client, types, and adapter
- **🧪 Testable**: Clean interfaces that can be easily mocked and tested
- **🔧 Extensible**: New providers can be added without touching existing code

## 🚀 Quick Start

### Installation

```bash
go get github.com/grokify/gollm
```

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/grokify/gollm"
)

func main() {
    // Create a client for OpenAI
    client, err := gollm.NewClient(gollm.ClientConfig{
        Provider: gollm.ProviderNameOpenAI,
        APIKey:   "your-openai-api-key",
    })
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Create a chat completion request
    response, err := client.CreateChatCompletion(context.Background(), &gollm.ChatCompletionRequest{
        Model: gollm.ModelGPT4o,
        Messages: []gollm.Message{
            {
                Role:    gollm.RoleUser,
                Content: "Hello! How can you help me today?",
            },
        },
        MaxTokens:   &[]int{150}[0],
        Temperature: &[]float64{0.7}[0],
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Response: %s\n", response.Choices[0].Message.Content)
    fmt.Printf("Tokens used: %d\n", response.Usage.TotalTokens)
}
```

## 🔧 Supported Providers

### OpenAI

- **Models**: GPT-4o, GPT-4o-mini, GPT-4-turbo, GPT-3.5-turbo
- **Features**: Chat completions, streaming, function calling

```go
client, err := gollm.NewClient(gollm.ClientConfig{
    Provider: gollm.ProviderNameOpenAI,
    APIKey:   "your-openai-api-key",
    BaseURL:  "https://api.openai.com/v1", // optional
})
```

### Anthropic (Claude)

- **Models**: Claude-3-Opus, Claude-3-Sonnet, Claude-3-Haiku, Claude-Sonnet-4
- **Features**: Chat completions, streaming, system message support

```go
client, err := gollm.NewClient(gollm.ClientConfig{
    Provider: gollm.ProviderNameAnthropic,
    APIKey:   "your-anthropic-api-key",
    BaseURL:  "https://api.anthropic.com", // optional
})
```

### Google Gemini

- **Models**: Gemini-2.5-Pro, Gemini-2.5-Flash, Gemini-1.5-Pro, Gemini-1.5-Flash, Gemini-Pro
- **Features**: Chat completions, streaming

```go
client, err := gollm.NewClient(gollm.ClientConfig{
    Provider: gollm.ProviderNameGemini,
    APIKey:   "your-gemini-api-key",
})
```

### AWS Bedrock

- **Models**: Anthropic Claude models, Amazon Titan
- **Features**: AWS IAM-based authentication, multiple model families

```go
client, err := gollm.NewClient(gollm.ClientConfig{
    Provider: gollm.ProviderNameBedrock,
    Region:   "us-east-1", // AWS region
})
```

### Ollama (Local Models)

- **Models**: Llama 3, Mistral, CodeLlama, Gemma, Qwen2.5, DeepSeek-Coder
- **Features**: Local inference, no API keys required, optimized for Apple Silicon

```go
client, err := gollm.NewClient(gollm.ClientConfig{
    Provider: gollm.ProviderNameOllama,
    BaseURL:  "http://localhost:11434", // default Ollama endpoint
})
```

## 📡 Streaming Example

```go
stream, err := client.CreateChatCompletionStream(context.Background(), &gollm.ChatCompletionRequest{
    Model: gollm.ModelGPT4o,
    Messages: []gollm.Message{
        {
            Role:    gollm.RoleUser,
            Content: "Tell me a short story about AI.",
        },
    },
    MaxTokens:   &[]int{200}[0],
    Temperature: &[]float64{0.8}[0],
})
if err != nil {
    log.Fatal(err)
}
defer stream.Close()

fmt.Print("AI Response: ")
for {
    chunk, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatal(err)
    }
    
    if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
        fmt.Print(chunk.Choices[0].Delta.Content)
    }
}
fmt.Println()
```

## 🧠 Conversation Memory

GoLLM supports persistent conversation memory using any Key-Value Store that implements the [Sogo KVS interface](https://github.com/grokify/sogo/blob/master/database/kvs/definitions.go). This enables multi-turn conversations that persist across application restarts.

### Memory Configuration

```go
// Configure memory settings
memoryConfig := gollm.MemoryConfig{
    MaxMessages: 50,                    // Keep last 50 messages per session
    TTL:         24 * time.Hour,       // Messages expire after 24 hours
    KeyPrefix:   "myapp:conversations", // Custom key prefix
}

// Create client with memory (using Redis, DynamoDB, etc.)
client, err := gollm.NewClient(gollm.ClientConfig{
    Provider:     gollm.ProviderNameOpenAI,
    APIKey:       "your-api-key",
    Memory:       kvsClient,          // Your KVS implementation
    MemoryConfig: &memoryConfig,
})
```

### Memory-Aware Completions

```go
// Create a session with system message
err = client.CreateConversationWithSystemMessage(ctx, "user-123", 
    "You are a helpful assistant that remembers our conversation history.")

// Use memory-aware completion - automatically loads conversation history
response, err := client.CreateChatCompletionWithMemory(ctx, "user-123", &gollm.ChatCompletionRequest{
    Model: gollm.ModelGPT4o,
    Messages: []gollm.Message{
        {Role: gollm.RoleUser, Content: "What did we discuss last time?"},
    },
    MaxTokens: &[]int{200}[0],
})

// The response will include context from previous conversations in this session
```

### Memory Management

```go
// Load conversation history
conversation, err := client.LoadConversation(ctx, "user-123")

// Get just the messages
messages, err := client.GetConversationMessages(ctx, "user-123")

// Manually append messages
err = client.AppendMessage(ctx, "user-123", gollm.Message{
    Role:    gollm.RoleUser,
    Content: "Remember this important fact: I prefer JSON responses.",
})

// Delete conversation
err = client.DeleteConversation(ctx, "user-123")
```

### KVS Backend Support

Memory works with any KVS implementation:
- **Redis**: For high-performance, distributed memory
- **DynamoDB**: For AWS-native storage
- **In-Memory**: For testing and development
- **Custom**: Any implementation of the Sogo KVS interface

```go
// Example with Redis (using a hypothetical Redis KVS implementation)
redisKVS := redis.NewKVSClient("localhost:6379")
client, err := gollm.NewClient(gollm.ClientConfig{
    Provider: gollm.ProviderNameOpenAI,
    APIKey:   "your-key",
    Memory:   redisKVS,
})
```

## 🔄 Provider Switching

The unified interface makes it easy to switch between providers:

```go
// Same request works with any provider
request := &gollm.ChatCompletionRequest{
    Model: gollm.ModelGPT4o, // or gollm.ModelClaude3Sonnet, etc.
    Messages: []gollm.Message{
        {Role: gollm.RoleUser, Content: "Hello, world!"},
    },
    MaxTokens: &[]int{100}[0],
}

// OpenAI
openaiClient, _ := gollm.NewClient(gollm.ClientConfig{
    Provider: gollm.ProviderNameOpenAI,
    APIKey:   "openai-key",
})

// Anthropic  
anthropicClient, _ := gollm.NewClient(gollm.ClientConfig{
    Provider: gollm.ProviderNameAnthropic,
    APIKey:   "anthropic-key",
})

// Gemini
geminiClient, _ := gollm.NewClient(gollm.ClientConfig{
    Provider: gollm.ProviderNameGemini,
    APIKey:   "gemini-key",
})

// Same API call for all providers
response1, _ := openaiClient.CreateChatCompletion(ctx, request)
response2, _ := anthropicClient.CreateChatCompletion(ctx, request)
response3, _ := geminiClient.CreateChatCompletion(ctx, request)
```

## 🧪 Testing

The clean interface design makes testing straightforward:

```go
// Mock the Provider interface for testing
type mockProvider struct{}

func (m *mockProvider) CreateChatCompletion(ctx context.Context, req *gollm.ChatCompletionRequest) (*gollm.ChatCompletionResponse, error) {
    return &gollm.ChatCompletionResponse{
        Choices: []gollm.ChatCompletionChoice{
            {
                Message: gollm.Message{
                    Role:    gollm.RoleAssistant,
                    Content: "Mock response",
                },
            },
        },
    }, nil
}

func (m *mockProvider) CreateChatCompletionStream(ctx context.Context, req *gollm.ChatCompletionRequest) (gollm.ChatCompletionStream, error) {
    return nil, nil
}

func (m *mockProvider) Close() error { return nil }
func (m *mockProvider) Name() string { return "mock" }
```

## 📚 Examples

The repository includes comprehensive examples:

- **Basic Usage**: Simple chat completions with each provider
- **Streaming**: Real-time response handling
- **Conversation**: Multi-turn conversations with context
- **Memory Demo**: Persistent conversation memory with KVS backend
- **Architecture Demo**: Overview of the provider architecture
- **Custom Provider**: How to create and use 3rd party providers

Run examples:
```bash
go run examples/basic/main.go
go run examples/streaming/main.go
go run examples/anthropic_streaming/main.go
go run examples/conversation/main.go
go run examples/memory_demo/main.go
go run examples/providers_demo/main.go
go run examples/ollama/main.go
go run examples/ollama_streaming/main.go
go run examples/gemini/main.go
go run examples/custom_provider/main.go
```

## 🔧 Configuration

### Environment Variables
- `OPENAI_API_KEY`: Your OpenAI API key
- `ANTHROPIC_API_KEY`: Your Anthropic API key
- `GEMINI_API_KEY`: Your Google Gemini API key
- AWS credentials for Bedrock (via AWS CLI/SDK configuration)

### Advanced Configuration

```go
config := gollm.ClientConfig{
    Provider: gollm.ProviderNameOpenAI,
    APIKey:   "your-api-key",
    BaseURL:  "https://custom-endpoint.com/v1",
    Extra: map[string]interface{}{
        "timeout": 60, // Custom provider-specific settings
    },
}
```

## 🏗️ Adding New Providers

### 🎯 3rd Party Providers (Recommended)

External packages can create providers without modifying the core library. This is the recommended approach for most use cases:

#### Step 1: Create Your Provider Package

```go
// In your external package (e.g., github.com/yourname/gollm-gemini)
package gemini

import (
    "context"
    "github.com/grokify/gollm/provider"
)

// Step 1: HTTP Client (like providers/openai/openai.go)
type Client struct {
    apiKey string
    // your HTTP client implementation
}

func New(apiKey string) *Client {
    return &Client{apiKey: apiKey}
}

// Step 2: Provider Adapter (like providers/openai/adapter.go)
type Provider struct {
    client *Client
}

func NewProvider(apiKey string) provider.Provider {
    return &Provider{client: New(apiKey)}
}

func (p *Provider) CreateChatCompletion(ctx context.Context, req *provider.ChatCompletionRequest) (*provider.ChatCompletionResponse, error) {
    // Convert provider.ChatCompletionRequest to your API format
    // Make HTTP call via p.client
    // Convert response back to provider.ChatCompletionResponse
}

func (p *Provider) CreateChatCompletionStream(ctx context.Context, req *provider.ChatCompletionRequest) (provider.ChatCompletionStream, error) {
    // Your streaming implementation
}

func (p *Provider) Close() error { return p.client.Close() }
func (p *Provider) Name() string { return "gemini" }
```

#### Step 2: Use Your Provider

```go
import (
    "github.com/grokify/gollm"
    "github.com/yourname/gollm-gemini"
)

func main() {
    // Create your custom provider
    customProvider := gemini.NewProvider("your-api-key")
    
    // Inject it directly into gollm - no core modifications needed!
    client, err := gollm.NewClient(gollm.ClientConfig{
        CustomProvider: customProvider,
    })
    
    // Use the same gollm API
    response, err := client.CreateChatCompletion(ctx, &gollm.ChatCompletionRequest{
        Model: "gemini-pro",
        Messages: []gollm.Message{{Role: gollm.RoleUser, Content: "Hello!"}},
    })
}
```

### 🔧 Built-in Providers (For Core Contributors)

To add a built-in provider to the core library, follow the same structure as existing providers:

1. **Create Provider Package**: `providers/newprovider/`
   - `newprovider.go` - HTTP client implementation
   - `types.go` - Provider-specific request/response types
   - `adapter.go` - `provider.Provider` interface implementation

2. **Update Core Files**:
   - Add factory function in `providers.go`
   - Add provider constant in `constants.go`
   - Add model constants if needed

3. **Reference Implementation**: Look at any existing provider (e.g., `providers/openai/`) as they all follow the exact same pattern that external providers should use

### 🎯 Why This Architecture?

- **🔌 No Core Changes**: External providers don't require modifying the core library
- **🏗️ Reference Pattern**: Internal providers demonstrate the exact structure external providers should follow
- **🧪 Easy Testing**: Both internal and external providers use the same `provider.Provider` interface
- **📦 Self-Contained**: Each provider manages its own HTTP client, types, and adapter logic
- **🔧 Direct Injection**: Clean dependency injection via `ClientConfig.CustomProvider`

## 📊 Model Support

| Provider | Models | Features |
|----------|--------|----------|
| OpenAI | GPT-5, GPT-5-mini, GPT-5-nano, GPT-4.1, GPT-4.1-mini, GPT-4.1-nano | Chat, Streaming, Functions |
| Anthropic | Claude-Opus-4, Claude-Sonnet-4, Claude-3-Opus, Claude-3-Sonnet, Claude-3-Haiku | Chat, Streaming, System messages |
| Gemini | Gemini-2.5-Pro, Gemini-2.5-Flash, Gemini-Pro | Chat, Streaming |
| Bedrock | Claude models, Titan models | Chat, Multiple model families |
| Ollama | Llama 3, Mistral, CodeLlama, Gemma, Qwen2.5 | Chat, Streaming, Local inference |

## 🚨 Error Handling

GoLLM provides comprehensive error handling with provider-specific context:

```go
response, err := client.CreateChatCompletion(ctx, request)
if err != nil {
    if apiErr, ok := err.(*gollm.APIError); ok {
        fmt.Printf("Provider: %s, Status: %d, Message: %s\n", 
            apiErr.Provider, apiErr.StatusCode, apiErr.Message)
    }
}
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Related Projects

- [Anthropic Go SDK](https://github.com/anthropics/anthropic-sdk-go) - Official Anthropic SDK
- [OpenAI Go SDK](https://github.com/openai/openai-go) - Official OpenAI SDK
- [AWS SDK for Go](https://github.com/aws/aws-sdk-go-v2) - Official AWS SDK

---

**Made with ❤️ for the Go and AI community**

 [build-status-svg]: https://github.com/grokify/gollm/actions/workflows/ci.yaml/badge.svg?branch=main
 [build-status-url]: https://github.com/grokify/gollm/actions/workflows/ci.yaml
 [lint-status-svg]: https://github.com/grokify/gollm/actions/workflows/lint.yaml/badge.svg?branch=main
 [lint-status-url]: https://github.com/grokify/gollm/actions/workflows/lint.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/gollm
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/gollm
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/gollm
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/gollm
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/gollm/blob/master/LICENSE
 [used-by-svg]: https://sourcegraph.com/github.com/grokify/gollm/-/badge.svg
 [used-by-url]: https://sourcegraph.com/github.com/grokify/gollm?badge
