# Key capabilities

AgentCore Memory provides:

1. **Core Infrastructure:** Serverless setup with built-in encryption and observability
1. **Event Storage:** Raw event storage (conversation history/checkpointing) with branching support
1. **Strategy Management:** Configurable extraction strategies (SEMANTIC, SUMMARY, USER_PREFERENCES, CUSTOM)
1. **Memory Records Extraction:** Automatic extraction of facts, preferences, and summaries based on configured strategies
1. **Semantic Search:** Vector-based retrieval of relevant memories using natural language queries

## Memory Strategy Types

AgentCore Memory supports four strategy types:

1. **Semantic Memory:** Stores factual information using vector embeddings for similarity search
1. **Summary Memory:** Creates and maintains conversation summaries for context preservation
1. **User Preference Memory:** Tracks user-specific preferences and settings
1. **Custom Memory:** Allows customization of extraction and consolidation logic

## What is Short-Term Memory?

Short-term memory focuses on:

1. **Session Continuity:** Maintaining context within a single conversation session
1. **Immediate Context:** Preserving recent conversation history for coherent responses
1. **Temporary State:** Managing transient information that's relevant for the current interaction
1. **Conversation Flow:** Ensuring smooth transitions between topics within a session

## Best Practices

1. **Context Window Management:** Monitor context usage to prevent overflow
1. **Session Boundaries:** Clearly define when sessions begin and end
1. **Memory Cleanup:** Implement appropriate cleanup policies for expired sessions
1. **Error Handling:** Handle memory retrieval failures gracefully
1. **Performance Optimization:** Use efficient querying patterns (e.g. via Summary Strategy in long term) for large conversation histories

## Understanding Namespaces

Namespaces organize memory records within strategies using a path-like structure. They can include variables that are dynamically replaced:

* support/facts/{sessionId}: Organizes facts by session
* user/{actorId}/preferences: Stores user preferences by actor ID
* meetings/{memoryId}/summaries/{sessionId}: Groups summaries by memory

The {actorId}, {sessionId}, and {memoryId} variables are automatically replaced with actual values when storing and retrieving memories.
