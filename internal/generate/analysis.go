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
	AgentName    string
	Template     *templates.Template
	SourcePath   string
	MaxTurns     int
	AgentParams  map[string]string
}

// runAnalysisLoop executes the multi-turn conversation between AI and agent tools.
func (o *Orchestrator) runAnalysisLoop(ctx context.Context, opts AnalysisOptions) (string, error) {
	// Get the agent definition
	agentDef, exists := o.agentRegistry.Get(opts.AgentName)
	if !exists {
		return "", fmt.Errorf("agent not found: %s", opts.AgentName)
	}

	// Convert agent tools to AI tools
	aiTools := convertAgentTools(agentDef)
	
	// Initialize conversation with system prompt
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

	log.Info().
		Str("agent", opts.AgentName).
		Int("tools", len(aiTools)).
		Msg("Starting AI analysis loop")

	// Analysis loop
	for turn := 0; turn < opts.MaxTurns; turn++ {
		log.Debug().
			Int("turn", turn+1).
			Int("messages", len(messages)).
			Msg("Sending request to AI")

		// Get AI response
		openaiClient, ok := o.aiClient.(*ai.OpenAIClient)
		if !ok {
			return "", fmt.Errorf("AI client does not support tool calling")
		}

		response, err := openaiClient.ChatWithTools(ctx, messages, aiTools)
		if err != nil {
			return "", fmt.Errorf("AI request failed: %w", err)
		}

		// Check if AI wants to call tools
		if len(response.ToolCalls) > 0 {
			log.Debug().
				Int("tool_calls", len(response.ToolCalls)).
				Msg("AI requested tool calls")

			// Add assistant message with tool calls to history
			assistantMsg := ai.ChatMessage{
				Role:      "assistant",
				ToolCalls: response.ToolCalls,
			}
			messages = append(messages, assistantMsg)

			// Execute each tool call
			for _, toolCall := range response.ToolCalls {
				log.Info().
					Str("tool", toolCall.Name).
					Str("id", toolCall.ID).
					Msg("Executing tool")

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
				messages = append(messages, toolMsg)

				log.Debug().
					Str("tool", toolCall.Name).
					Int("output_len", len(toolOutput)).
					Msg("Tool executed successfully")
			}
		} else {
			// AI provided final answer
			log.Info().Msg("AI provided final response")
			
			// Add final assistant message to history
			messages = append(messages, ai.ChatMessage{
				Role:    "assistant",
				Content: response.Message,
			})

			// Check if response is valid JSON
			var jsonCheck interface{}
			if err := json.Unmarshal([]byte(response.Message), &jsonCheck); err != nil {
				// If not JSON, ask AI to format as JSON
				messages = append(messages, ai.ChatMessage{
					Role:    "user",
					Content: "Please format your response as valid JSON matching the template schema.",
				})
				continue
			}

			return response.Message, nil
		}

		// Check if we should stop (e.g., finish_reason is "stop")
		if response.FinishReason == "stop" && response.Message != "" {
			return response.Message, nil
		}
	}

	return "", fmt.Errorf("analysis loop reached maximum turns (%d) without completion", opts.MaxTurns)
}

// convertAgentTools converts agent tool definitions to AI tool format.
func convertAgentTools(agentDef *agent.Definition) []ai.Tool {
	tools := make([]ai.Tool, 0, len(agentDef.Spec.Tools))
	
	for _, agentTool := range agentDef.Spec.Tools {
		// Create parameter schema based on the tool's expected inputs
		params := map[string]interface{}{
			"type": "object",
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
					params["properties"].(map[string]interface{})[paramKey] = map[string]interface{}{
						"type":        "string",
						"description": fmt.Sprintf("Value for %s", paramName),
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

// GenerateWithAgent generates a document using an agent for analysis.
func (o *Orchestrator) GenerateWithAgent(ctx context.Context, opts Options, agentName string, agentParams map[string]string) error {
	// Load template
	template, err := o.registry.Load(opts.TemplateType)
	if err != nil {
		return fmt.Errorf("failed to load template: %w", err)
	}

	// Check if template has analysis prompts
	if template.Analysis == nil || template.Analysis.SystemPrompt == "" {
		return fmt.Errorf("template '%s' does not support agent-based analysis", opts.TemplateType)
	}

	// Determine source path
	sourcePath := ""
	if len(opts.Sources) > 0 {
		sourcePath = opts.Sources[0] // Use first source as primary path
	}

	// Run the analysis loop
	analysisOpts := AnalysisOptions{
		AgentName:   agentName,
		Template:    template,
		SourcePath:  sourcePath,
		MaxTurns:    10, // Default max turns
		AgentParams: agentParams,
	}

	jsonResponse, err := o.runAnalysisLoop(ctx, analysisOpts)
	if err != nil {
		return fmt.Errorf("analysis loop failed: %w", err)
	}

	// Validate the response against template schema
	if err := o.validator.ValidateJSON(jsonResponse, template.FieldSchema); err != nil {
		log.Warn().Err(err).Msg("AI response validation failed, attempting repair")
		
		// Attempt to repair the JSON
		repairedJSON, repairErr := o.validator.RepairJSON(jsonResponse, template.FieldSchema, opts.MaxRepairs)
		if repairErr != nil {
			return fmt.Errorf("failed to repair JSON: %w", repairErr)
		}
		jsonResponse = repairedJSON
	}

	// Parse the JSON response
	var fields map[string]interface{}
	if err := json.Unmarshal([]byte(jsonResponse), &fields); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Render the document
	htmlContent, err := o.renderer.RenderHTML(template.HTMLTemplate, fields)
	if err != nil {
		return fmt.Errorf("failed to render HTML: %w", err)
	}

	// Write output files
	if err := o.writeOutput(opts.OutputFile, htmlContent, fields, opts.Force); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	log.Info().
		Str("output", opts.OutputFile).
		Msg("Document generated successfully with agent analysis")

	return nil
}

// Additional fields needed in Orchestrator
type enhancedOrchestrator struct {
	*Orchestrator
	agentRegistry *agent.Registry
	agentExecutor *agent.Executor
}