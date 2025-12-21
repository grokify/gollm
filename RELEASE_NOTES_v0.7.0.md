# Release Notes - FluxLLM v0.7.0

**Release Date:** 2025-12-21
**Base Version:** v0.6.1

## Overview

Version 0.7.0 is a major release that renames the module from `gollm` to `fluxllm` and introduces comprehensive observability features including hooks for tracing/logging/metrics, injectable structured logging via `slog`, and context-aware logging for request-scoped correlation.

**Summary:**
- **Module Rename**: `github.com/grokify/gollm` → `github.com/grokify/fluxllm`
- **Observability Hooks**: New `ObservabilityHook` interface for non-invasive tracing, logging, and metrics
- **Injectable Logging**: `*slog.Logger` support with null logger default
- **Context-Aware Logging**: Request-scoped logging via `slogutil.ContextWithLogger`
- **Call Correlation**: Unique `CallID` in `LLMCallInfo` for correlating hook calls
- **Bug Fix**: Memory-aware methods now properly invoke observability hooks

---

## Breaking Changes

### Module Rename

The module has been renamed from `gollm` to `fluxllm`:

**Before:**
```go
import "github.com/grokify/gollm"

client, err := gollm.NewClient(gollm.ClientConfig{...})
```

**After:**
```go
import "github.com/grokify/fluxllm"

client, err := fluxllm.NewClient(fluxllm.ClientConfig{...})
```

**Migration:**
1. Update import paths: `github.com/grokify/gollm` → `github.com/grokify/fluxllm`
2. Update type prefixes: `gollm.` → `fluxllm.`
3. Update go.mod: `go get github.com/grokify/fluxllm@v0.7.0`

---

## New Features

### 1. Observability Hooks

New `ObservabilityHook` interface allows external packages to observe LLM calls without modifying the core library. This enables integration with OpenTelemetry, Datadog, custom metrics systems, and more.

**New File:** `observability.go`

**Interface:**
```go
type LLMCallInfo struct {
    CallID       string    // Unique identifier for correlating BeforeRequest/AfterResponse
    ProviderName string    // e.g., "openai", "anthropic"
    StartTime    time.Time // When the call started
}

type ObservabilityHook interface {
    // BeforeRequest is called before each LLM call.
    // Returns a new context for trace/span propagation.
    BeforeRequest(ctx context.Context, info LLMCallInfo, req *provider.ChatCompletionRequest) context.Context

    // AfterResponse is called after each LLM call completes (success or failure).
    AfterResponse(ctx context.Context, info LLMCallInfo, req *provider.ChatCompletionRequest, resp *provider.ChatCompletionResponse, err error)

    // WrapStream wraps a stream for observability of streaming responses.
    WrapStream(ctx context.Context, info LLMCallInfo, req *provider.ChatCompletionRequest, stream provider.ChatCompletionStream) provider.ChatCompletionStream
}
```

**Usage:**
```go
type LoggingHook struct{}

func (h *LoggingHook) BeforeRequest(ctx context.Context, info fluxllm.LLMCallInfo, req *fluxllm.ChatCompletionRequest) context.Context {
    log.Printf("[%s] LLM call started: provider=%s model=%s", info.CallID, info.ProviderName, req.Model)
    return ctx
}

func (h *LoggingHook) AfterResponse(ctx context.Context, info fluxllm.LLMCallInfo, req *fluxllm.ChatCompletionRequest, resp *fluxllm.ChatCompletionResponse, err error) {
    duration := time.Since(info.StartTime)
    log.Printf("[%s] LLM call completed: duration=%v", info.CallID, duration)
}

func (h *LoggingHook) WrapStream(ctx context.Context, info fluxllm.LLMCallInfo, req *fluxllm.ChatCompletionRequest, stream fluxllm.ChatCompletionStream) fluxllm.ChatCompletionStream {
    return stream
}

// Configure hook
client, err := fluxllm.NewClient(fluxllm.ClientConfig{
    Provider:          fluxllm.ProviderNameOpenAI,
    APIKey:            "your-api-key",
    ObservabilityHook: &LoggingHook{},
})
```

**Key Features:**
- Non-invasive: Add observability without modifying core library code
- Provider agnostic: Works with all LLM providers
- Context propagation: Pass trace context through the entire call chain
- Streaming support: Wrap streams to observe streaming responses
- Unique CallID: Correlate BeforeRequest/AfterResponse in concurrent scenarios

### 2. Injectable Structured Logging

FluxLLM now supports injectable logging via Go's standard `log/slog` package.

**New Fields:**
- `ClientConfig.Logger` - Optional `*slog.Logger` for internal logging
- `ChatClient.Logger()` - Accessor method to retrieve the logger

**Behavior:**
- If no logger is provided, a null logger is used (zero overhead, no output)
- Logger is used for non-critical errors that shouldn't fail the main request

**Usage:**
```go
import (
    "log/slog"
    "os"
    "github.com/grokify/fluxllm"
)

logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))

client, err := fluxllm.NewClient(fluxllm.ClientConfig{
    Provider: fluxllm.ProviderNameOpenAI,
    APIKey:   "your-api-key",
    Logger:   logger,
})

// Access logger if needed
client.Logger().Info("client initialized")
```

### 3. Context-Aware Logging

Support for request-scoped logging via context, enabling trace IDs, user IDs, and other request-specific attributes to flow through all log output.

**New Functions (in `github.com/grokify/mogo/log/slogutil` v0.72.5):**
- `slogutil.ContextWithLogger(ctx, logger)` - Attach a logger to context
- `slogutil.LoggerFromContext(ctx, fallback)` - Retrieve logger from context

**Usage:**
```go
import (
    "log/slog"
    "github.com/grokify/fluxllm"
    "github.com/grokify/mogo/log/slogutil"
)

// Create request-scoped logger with trace context
reqLogger := slog.Default().With(
    slog.String("trace_id", traceID),
    slog.String("user_id", userID),
    slog.String("request_id", requestID),
)

// Attach to context
ctx = slogutil.ContextWithLogger(ctx, reqLogger)

// All internal logging will now include trace_id, user_id, request_id
response, err := client.CreateChatCompletionWithMemory(ctx, sessionID, req)
```

---

## Bug Fixes

### Memory-Aware Methods Now Invoke Observability Hooks

**Issue:** `CreateChatCompletionWithMemory` and `CreateChatCompletionStreamWithMemory` were calling the provider directly, bypassing the observability hook.

**Before (broken):**
```go
response, err := c.provider.CreateChatCompletion(ctx, &memoryReq)  // Hook not called!
```

**After (fixed):**
```go
response, err := c.CreateChatCompletion(ctx, &memoryReq)  // Hook is called
```

This ensures all LLM calls are properly observed regardless of whether memory is used.

---

## Dependencies

### Updated
- `github.com/grokify/mogo` v0.72.4 → v0.72.5
  - Added `slogutil.ContextWithLogger()` and `slogutil.LoggerFromContext()`

---

## Files Changed

| File | Changes |
|------|---------|
| `client.go` | Added `ObservabilityHook` field, `Logger` field, hook invocations, context-aware logging |
| `observability.go` | **New file** - `LLMCallInfo`, `ObservabilityHook` interface, `newCallID()` |
| `go.mod` | Module rename, mogo upgrade to v0.72.5 |
| `README.md` | Renamed to FluxLLM, added Observability Hooks section, Logging Configuration, Context-Aware Logging |

---

## Migration Guide

### Upgrading from v0.6.x

**1. Update import paths:**
```bash
# Using sed (macOS/Linux)
find . -name "*.go" -exec sed -i '' 's|github.com/grokify/gollm|github.com/grokify/fluxllm|g' {} +
find . -name "*.go" -exec sed -i '' 's|gollm\.|fluxllm.|g' {} +
```

**2. Update go.mod:**
```bash
go get github.com/grokify/fluxllm@v0.7.0
go mod tidy
```

**3. (Optional) Add observability:**
```go
client, err := fluxllm.NewClient(fluxllm.ClientConfig{
    Provider:          fluxllm.ProviderNameOpenAI,
    APIKey:            apiKey,
    ObservabilityHook: &YourHook{},  // New
    Logger:            slog.Default(), // New
})
```

---

## API Additions

### New Types
- `LLMCallInfo` - Metadata about LLM calls (CallID, ProviderName, StartTime)
- `ObservabilityHook` - Interface for observing LLM calls

### New Functions
- `newCallID()` - Generates unique call IDs (internal)

### New Fields
- `ClientConfig.ObservabilityHook` - Hook for tracing/logging/metrics
- `ClientConfig.Logger` - Injectable `*slog.Logger`
- `ChatClient.hook` - Stores the observability hook
- `ChatClient.logger` - Stores the logger
- `LLMCallInfo.CallID` - Unique identifier for call correlation

### New Methods
- `ChatClient.Logger()` - Returns the client's logger

---

## Full Changelog

**v0.6.1...v0.7.0**

- feat: rename module from `gollm` to `fluxllm`
- feat: add `ObservabilityHook` interface for tracing/logging/metrics
- feat: add `LLMCallInfo` with `CallID` for call correlation
- feat: add injectable `*slog.Logger` with null logger default
- feat: add context-aware logging via `slogutil.ContextWithLogger`
- fix: memory-aware methods now properly invoke observability hooks
- docs: update README.md for FluxLLM rename
- docs: add Observability Hooks documentation
- docs: add Logging Configuration documentation
- docs: add Context-Aware Logging documentation
- chore: upgrade `github.com/grokify/mogo` to v0.72.5

---

## Acknowledgments

This release focuses on production-readiness with comprehensive observability support, enabling teams to integrate FluxLLM with their existing monitoring and tracing infrastructure.

---

## Resources

- **Documentation:** [README.md](README.md)
- **Examples:** See `/examples` directory
- **Issues:** [GitHub Issues](https://github.com/grokify/fluxllm/issues)

---

**Full Changelog:** [v0.6.1...v0.7.0](https://github.com/grokify/fluxllm/compare/v0.6.1...v0.7.0)
