# FluxLLM Feature Roadmap

## High Value

### 1. Retry with Backoff ✅
Automatic retries for transient failures (rate limits, 5xx errors).

**Status:** Implemented in v0.7.0 via `ClientConfig.HTTPClient` with `retryhttp.RetryTransport`.

```go
rt := retryhttp.NewWithOptions(
    retryhttp.WithMaxRetries(5),
    retryhttp.WithInitialBackoff(500 * time.Millisecond),
)
client, err := fluxllm.NewClient(fluxllm.ClientConfig{
    Provider:   fluxllm.ProviderNameOpenAI,
    APIKey:     "...",
    HTTPClient: &http.Client{Transport: rt},
})
```

### 2. Request Timeouts
Per-request timeout configuration (currently relies on context).

```go
ClientConfig{
    Timeout: 30 * time.Second,
}
```

### 3. Fallback Providers
Automatic failover when primary provider fails.

```go
ClientConfig{
    Provider: fluxllm.ProviderNameOpenAI,
    FallbackProviders: []ProviderConfig{
        {Provider: fluxllm.ProviderNameAnthropic, APIKey: "..."},
    },
}
```

## Medium Value

### 4. Rate Limiting
Client-side rate limiter to respect provider limits.

### 5. Token Counting/Estimation
Estimate tokens before sending to avoid limit errors.

### 6. Response Caching
Cache identical requests to reduce costs (with TTL).

### 7. Circuit Breaker
Prevent cascading failures when provider is unhealthy.

## Nice to Have

### 8. Embeddings API
Unified interface for text embeddings.

### 9. Structured Output Validation
JSON schema validation for responses.

### 10. Batch Processing
Efficient batch request handling.

---

## Questions to Consider

- What's the primary use case? (chatbot, batch processing, real-time?)
- Is cost optimization important? (caching, token counting)
- How critical is uptime? (fallbacks, circuit breaker)
