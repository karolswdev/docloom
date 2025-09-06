// Package ai provides a provider-agnostic AI client for interacting with OpenAI-compatible APIs.
package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	openai "github.com/sashabaranov/go-openai"
)

// Client defines the interface for AI model interactions.
type Client interface {
	// GenerateJSON sends a prompt to the AI model and returns the generated JSON response.
	GenerateJSON(ctx context.Context, prompt string) (string, error)
}

// Config holds the configuration for the AI client.
type Config struct {
	Seed        *int
	BaseURL     string
	APIKey      string
	Model       string
	MaxTokens   int
	MaxRetries  int
	RetryDelay  time.Duration
	Temperature float32
}

// OpenAIClient implements the Client interface using the go-openai library.
type OpenAIClient struct {
	client *openai.Client
	config Config
}

// NewOpenAIClient creates a new OpenAI-compatible client.
func NewOpenAIClient(config Config) (*OpenAIClient, error) {
	if config.APIKey == "" {
		return nil, errors.New("API key is required")
	}
	if config.Model == "" {
		return nil, errors.New("model is required")
	}
	if config.BaseURL == "" {
		config.BaseURL = "https://api.openai.com/v1"
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = time.Second
	}
	if config.MaxTokens == 0 {
		config.MaxTokens = 4096
	}

	clientConfig := openai.DefaultConfig(config.APIKey)
	clientConfig.BaseURL = config.BaseURL

	return &OpenAIClient{
		client: openai.NewClientWithConfig(clientConfig),
		config: config,
	}, nil
}

// GenerateJSON implements the Client interface.
func (c *OpenAIClient) GenerateJSON(ctx context.Context, prompt string) (string, error) {
	var lastErr error
	delay := c.config.RetryDelay

	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			log.Info().
				Int("attempt", attempt).
				Dur("delay", delay).
				Msg("Retrying AI request after delay")

			select {
			case <-time.After(delay):
				// Continue with retry
			case <-ctx.Done():
				return "", ctx.Err()
			}

			// Exponential backoff
			delay *= 2
		}

		response, err := c.makeRequest(ctx, prompt)
		if err == nil {
			return response, nil
		}

		lastErr = err

		// Check if error is retryable
		if !isRetryableError(err) {
			return "", err
		}

		log.Warn().
			Err(err).
			Int("attempt", attempt).
			Int("max_retries", c.config.MaxRetries).
			Msg("AI request failed, will retry")
	}

	return "", fmt.Errorf("failed after %d retries: %w", c.config.MaxRetries+1, lastErr)
}

func (c *OpenAIClient) makeRequest(ctx context.Context, prompt string) (string, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are a helpful assistant that generates structured JSON output based on the provided instructions. Always respond with valid JSON only, no additional text.",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
	}

	req := openai.ChatCompletionRequest{
		Model:       c.config.Model,
		Messages:    messages,
		Temperature: c.config.Temperature,
		MaxTokens:   c.config.MaxTokens,
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
	}

	if c.config.Seed != nil {
		req.Seed = c.config.Seed
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("AI request failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no response choices from AI model")
	}

	content := resp.Choices[0].Message.Content

	// Validate that the response is valid JSON
	var jsonCheck interface{}
	if err := json.Unmarshal([]byte(content), &jsonCheck); err != nil {
		return "", fmt.Errorf("AI response is not valid JSON: %w", err)
	}

	return content, nil
}

// isRetryableError determines if an error should trigger a retry.
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Check for specific OpenAI API errors that are retryable
	var apiErr *openai.APIError
	if errors.As(err, &apiErr) {
		// Retry on rate limit, server errors, and service unavailable
		switch apiErr.HTTPStatusCode {
		case 429, // Too Many Requests
			500, // Internal Server Error
			502, // Bad Gateway
			503, // Service Unavailable
			504: // Gateway Timeout
			return true
		}
	}

	// Check for context errors (don't retry on cancellation)
	if errors.Is(err, context.Canceled) {
		return false
	}

	// Retry on timeout
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	return false
}
