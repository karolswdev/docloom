// Package generate provides the AI analysis loop for agent-based document generation.
package generate

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/karolswdev/docloom/internal/agent"
	"github.com/karolswdev/docloom/internal/ai"
	"github.com/karolswdev/docloom/internal/templates"
)

// AnalysisOptions contains configuration for the analysis loop.
type AnalysisOptions struct {
	AgentName   string
	Template    *templates.Template
	SourcePath  string
	MaxTurns    int
	AgentParams map[string]string
}

// RunAnalysisLoop executes the multi-turn conversation between AI and agent tools.
func (o *Orchestrator) RunAnalysisLoop(ctx context.Context, opts AnalysisOptions) (string, error) {
	// Get the agent definition
	agentDef, exists := o.agentRegistry.Get(opts.AgentName)
	if !exists {
		return "", fmt.Errorf("agent not found: %s", opts.AgentName)
	}

	// Convert agent tools to AI tools
	aiTools := convertAgentTools(agentDef)

	// Initialize conversation
	messages := o.initializeConversation(opts)

	log.Info().
		Str("agent", opts.AgentName).
		Int("tools", len(aiTools)).
		Msg("Starting AI analysis loop")

	// Analysis loop
	for turn := 0; turn < opts.MaxTurns; turn++ {
		result, shouldContinue, err := o.executeAnalysisTurn(ctx, turn, messages, aiTools, opts)
		if err != nil {
			return "", err
		}
		if result != "" {
			return result, nil
		}
		if !shouldContinue {
			break
		}
	}

	return "", fmt.Errorf("analysis loop reached maximum turns (%d) without completion", opts.MaxTurns)
}

// initializeConversation sets up the initial message context.
func (o *Orchestrator) initializeConversation(opts AnalysisOptions) []ai.ChatMessage {
	messages := []ai.ChatMessage{
		{
			Role:    "system",
			Content: opts.Template.Analysis.SystemPrompt,
		},
		{
			Role:    "user",
			Content: opts.Template.Analysis.InitialUserPrompt,
		},
	}

	// Add source context if available
	if opts.SourcePath != "" {
		messages[1].Content += fmt.Sprintf("\n\nRepository path: %s", opts.SourcePath)
	}

	return messages
}

// executeAnalysisTurn performs a single turn of the analysis loop.
func (o *Orchestrator) executeAnalysisTurn(ctx context.Context, turn int, messages []ai.ChatMessage, aiTools []ai.Tool, opts AnalysisOptions) (string, bool, error) {
	log.Debug().
		Int("turn", turn+1).
		Int("messages", len(messages)).
		Msg("Sending request to AI")

	// Get AI response
	openaiClient, ok := o.aiClient.(*ai.OpenAIClient)
	if !ok {
		return "", false, fmt.Errorf("AI client does not support tool calling")
	}

	response, err := openaiClient.ChatWithTools(ctx, messages, aiTools)
	if err != nil {
		return "", false, fmt.Errorf("AI request failed: %w", err)
	}

	// Check if AI wants to call tools
	if len(response.ToolCalls) > 0 {
		err := o.handleToolCalls(response.ToolCalls, &messages, opts)
		if err != nil {
			return "", false, err
		}
		return "", true, nil
	}

	// AI provided a response
	return o.handleAIResponse(response, &messages)
}

// handleToolCalls processes and executes requested tool calls.
func (o *Orchestrator) handleToolCalls(toolCalls []ai.ToolCall, messages *[]ai.ChatMessage, opts AnalysisOptions) error {
	log.Debug().
		Int("tool_calls", len(toolCalls)).
		Msg("AI requested tool calls")

	// Add assistant message with tool calls to history
	assistantMsg := ai.ChatMessage{
		Role:      "assistant",
		ToolCalls: toolCalls,
	}
	*messages = append(*messages, assistantMsg)

	// Execute each tool call
	for _, toolCall := range toolCalls {
		err := o.executeSingleTool(toolCall, messages, opts)
		if err != nil {
			return err
		}
	}

	return nil
}

// executeSingleTool executes a single tool call and adds the result to messages.
func (o *Orchestrator) executeSingleTool(toolCall ai.ToolCall, messages *[]ai.ChatMessage, opts AnalysisOptions) error {
	log.Info().
		Str("tool", toolCall.Name).
		Str("id", toolCall.ID).
		Msg("Executing tool")

	// Parse and prepare arguments
	args := o.prepareToolArguments(toolCall, opts)

	// Execute the tool
	toolOutput, err := o.agentExecutor.RunTool(opts.AgentName, toolCall.Name, args)
	if err != nil {
		// Add error as tool response
		toolOutput = fmt.Sprintf("Error executing tool: %v", err)
		log.Error().
			Err(err).
			Str("tool", toolCall.Name).
			Msg("Tool execution failed")
	}

	// Add tool response to conversation
	toolMsg := ai.ChatMessage{
		Role:       "tool",
		Content:    toolOutput,
		ToolCallID: toolCall.ID,
	}
	*messages = append(*messages, toolMsg)

	log.Debug().
		Str("tool", toolCall.Name).
		Int("output_len", len(toolOutput)).
		Msg("Tool executed successfully")

	return nil
}

// prepareToolArguments prepares arguments for tool execution.
func (o *Orchestrator) prepareToolArguments(toolCall ai.ToolCall, opts AnalysisOptions) map[string]string {
	// Parse tool arguments
	var args map[string]string
	if err := json.Unmarshal(toolCall.Arguments, &args); err != nil {
		// If arguments are not a map, try as a simple string
		args = map[string]string{"input": string(toolCall.Arguments)}
	}

	// Merge with agent parameters
	for k, v := range opts.AgentParams {
		if _, exists := args[k]; !exists {
			args[k] = v
		}
	}

	// Add source path if not present
	if _, exists := args["SOURCE_PATH"]; !exists && opts.SourcePath != "" {
		args["SOURCE_PATH"] = opts.SourcePath
	}

	return args
}

// handleAIResponse processes the AI's final response.
func (o *Orchestrator) handleAIResponse(response *ai.ChatResponse, messages *[]ai.ChatMessage) (string, bool, error) {
	log.Info().Msg("AI provided final response")

	// Add final assistant message to history
	*messages = append(*messages, ai.ChatMessage{
		Role:    "assistant",
		Content: response.Message,
	})

	// Check if response is valid JSON
	var jsonCheck interface{}
	if err := json.Unmarshal([]byte(response.Message), &jsonCheck); err != nil {
		// If not JSON, ask AI to format as JSON
		*messages = append(*messages, ai.ChatMessage{
			Role:    "user",
			Content: "Please format your response as valid JSON matching the template schema.",
		})
		return "", true, nil
	}

	// Check if we should stop (e.g., finish_reason is "stop")
	if response.FinishReason == "stop" && response.Message != "" {
		return response.Message, false, nil
	}

	return response.Message, false, nil
}

// convertAgentTools converts agent tool definitions to AI tool format.
func convertAgentTools(agentDef *agent.Definition) []ai.Tool {
	tools := make([]ai.Tool, 0, len(agentDef.Spec.Tools))

	for _, agentTool := range agentDef.Spec.Tools {
		// Create parameter schema based on the tool's expected inputs
		params := map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		}

		// Parse arguments to determine expected parameters
		for _, arg := range agentTool.Args {
			if strings.Contains(arg, "${") {
				// Extract parameter name
				start := strings.Index(arg, "${")
				end := strings.Index(arg[start:], "}")
				if end > 0 {
					paramName := arg[start+2 : start+end]
					// Convert to lowercase for consistency
					paramKey := strings.ToLower(paramName)
					if props, ok := params["properties"].(map[string]interface{}); ok {
						props[paramKey] = map[string]interface{}{
							"type":        "string",
							"description": fmt.Sprintf("Value for %s", paramName),
						}
					}
				}
			}
		}

		aiTool := ai.Tool{
			Name:        agentTool.Name,
			Description: agentTool.Description,
			Parameters:  params,
		}

		tools = append(tools, aiTool)
	}

	return tools
}
