package models

// AWS Bedrock Model Documentation
const (
	// BedrockModelsURL is the official AWS Bedrock models documentation page.
	// Use this to check for new models, deprecations, and model updates.
	BedrockModelsURL = "https://docs.aws.amazon.com/bedrock/latest/userguide/models-supported.html"

	// BedrockAPIURL is the AWS Bedrock API reference page.
	BedrockAPIURL = "https://docs.aws.amazon.com/bedrock/latest/APIReference/welcome.html"
)

// Bedrock Claude Models
const (
	// BedrockClaudeOpus4 is Claude Opus 4 on AWS Bedrock.
	BedrockClaudeOpus4 = "anthropic.claude-opus-4-20250514-v1:0"

	// BedrockClaude3Opus is Claude 3 Opus on AWS Bedrock.
	BedrockClaude3Opus = "anthropic.claude-3-opus-20240229-v1:0"

	// BedrockClaude3Sonnet is Claude 3 Sonnet on AWS Bedrock.
	BedrockClaude3Sonnet = "anthropic.claude-3-sonnet-20240229-v1:0"
)

// Bedrock Amazon Titan Models
const (
	// BedrockTitan is Amazon Titan Text Express on AWS Bedrock.
	BedrockTitan = "amazon.titan-text-express-v1"
)
