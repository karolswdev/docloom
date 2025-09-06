// Package ai provides tool calling support for AI-orchestrated agent interaction.
package ai

import (
	"context"
	"encoding/json"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

// Tool represents a tool that can be called by the AI.
type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters,omitempty"`
}

// ToolCall represents a request from the AI to call a specific tool.
type ToolCall struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

// ToolResult represents the result of executing a tool.
type ToolResult struct {
	ToolCallID string `json:"tool_call_id"`
	Content    string `json:"content"`
}

// ChatMessage represents a message in the conversation.
type ChatMessage struct {
	Role       string     `json:"role"` // "system", "user", "assistant", "tool"
	Content    string     `json:"content,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"` // For tool response messages
}

// ChatResponse represents the AI's response which can be either a message or tool calls.
type ChatResponse struct {
	Message      string     `json:"message,omitempty"`
	ToolCalls    []ToolCall `json:"tool_calls,omitempty"`
	FinishReason string     `json:"finish_reason,omitempty"`
}

// ChatWithTools sends a conversation with available tools to the AI and returns its response.
func (c *OpenAIClient) ChatWithTools(ctx context.Context, messages []ChatMessage, tools []Tool) (*ChatResponse, error) {
	// Convert our messages to OpenAI format
	openaiMessages := make([]openai.ChatCompletionMessage, 0, len(messages))
	for _, msg := range messages {
		openaiMsg := openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}

		// Handle tool messages
		if msg.Role == "tool" {
			openaiMsg.Role = openai.ChatMessageRoleTool
			openaiMsg.ToolCallID = msg.ToolCallID
		}

		// Handle assistant messages with tool calls
		if len(msg.ToolCalls) > 0 {
			openaiMsg.ToolCalls = make([]openai.ToolCall, len(msg.ToolCalls))
			for i, tc := range msg.ToolCalls {
				openaiMsg.ToolCalls[i] = openai.ToolCall{
					ID:   tc.ID,
					Type: "function",
					Function: openai.FunctionCall{
						Name:      tc.Name,
						Arguments: string(tc.Arguments),
					},
				}
			}
		}

		openaiMessages = append(openaiMessages, openaiMsg)
	}

	// Convert tools to OpenAI format
	openaiTools := make([]openai.Tool, 0, len(tools))
	for _, tool := range tools {
		// Create a simple parameter schema if not provided
		params := tool.Parameters
		if params == nil {
			params = map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			}
		}

		openaiTools = append(openaiTools, openai.Tool{
			Type: "function",
			Function: &openai.FunctionDefinition{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  params,
			},
		})
	}

	// Create the request
	req := openai.ChatCompletionRequest{
		Model:       c.config.Model,
		Messages:    openaiMessages,
		MaxTokens:   c.config.MaxTokens,
		Temperature: c.config.Temperature,
	}

	// Add tools if provided
	if len(openaiTools) > 0 {
		req.Tools = openaiTools
	}

	// Add seed if configured
	if c.config.Seed != nil {
		req.Seed = c.config.Seed
	}

	// Make the API call
	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("chat completion failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned from AI")
	}

	choice := resp.Choices[0]
	result := &ChatResponse{
		FinishReason: string(choice.FinishReason),
	}

	// Check if the response contains tool calls
	if len(choice.Message.ToolCalls) > 0 {
		result.ToolCalls = make([]ToolCall, len(choice.Message.ToolCalls))
		for i, tc := range choice.Message.ToolCalls {
			result.ToolCalls[i] = ToolCall{
				ID:        tc.ID,
				Name:      tc.Function.Name,
				Arguments: json.RawMessage(tc.Function.Arguments),
			}
		}
	} else {
		// Regular message response
		result.Message = choice.Message.Content
	}

	return result, nil
}

// ConvertAgentToolsToAITools converts agent tool definitions to AI-compatible tool definitions.
func ConvertAgentToolsToAITools(agentTools []interface{}) []Tool {
	aiTools := make([]Tool, 0)

	// This would be populated from the agent's tool definitions
	// For now, returning empty slice as a placeholder
	// The actual implementation would parse the agent.Spec.Tools

	return aiTools
}
