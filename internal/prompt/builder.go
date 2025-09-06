// Package prompt provides functionality for assembling prompts for AI model interactions.
package prompt

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Builder is responsible for constructing prompts for the AI model.
type Builder struct{}

// NewBuilder creates a new prompt builder.
func NewBuilder() *Builder {
	return &Builder{}
}

// BuildGenerationPrompt assembles a prompt for generating JSON content based on source documents and a template.
func (b *Builder) BuildGenerationPrompt(sourceContent string, templatePrompt string, schema interface{}) (string, error) {
	// Convert schema to JSON string if it's not already
	var schemaJSON string
	switch v := schema.(type) {
	case string:
		schemaJSON = v
	case []byte:
		schemaJSON = string(v)
	default:
		schemaBytes, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal schema: %w", err)
		}
		schemaJSON = string(schemaBytes)
	}

	// Build the prompt with clear sections
	var promptBuilder strings.Builder
	
	promptBuilder.WriteString("You are a technical documentation generator. ")
	promptBuilder.WriteString("Your task is to generate structured JSON content based on the provided source documents and template requirements.\n\n")
	
	promptBuilder.WriteString("## Template Instructions\n")
	promptBuilder.WriteString(templatePrompt)
	promptBuilder.WriteString("\n\n")
	
	promptBuilder.WriteString("## JSON Schema\n")
	promptBuilder.WriteString("Your response MUST conform to the following JSON schema:\n")
	promptBuilder.WriteString("```json\n")
	promptBuilder.WriteString(schemaJSON)
	promptBuilder.WriteString("\n```\n\n")
	
	promptBuilder.WriteString("## Source Documents\n")
	promptBuilder.WriteString("Use the following source content to generate the JSON fields:\n")
	promptBuilder.WriteString("```\n")
	promptBuilder.WriteString(sourceContent)
	promptBuilder.WriteString("\n```\n\n")
	
	promptBuilder.WriteString("## Instructions\n")
	promptBuilder.WriteString("1. Analyze the source documents carefully\n")
	promptBuilder.WriteString("2. Generate JSON that matches the schema exactly\n")
	promptBuilder.WriteString("3. Use information from the source documents to populate the fields\n")
	promptBuilder.WriteString("4. Ensure all required fields are present\n")
	promptBuilder.WriteString("5. Return ONLY valid JSON, no additional text or markdown formatting\n")

	return promptBuilder.String(), nil
}

// BuildRepairPrompt creates a prompt for repairing invalid JSON based on validation errors.
func (b *Builder) BuildRepairPrompt(originalPrompt string, invalidJSON string, validationError string, schema interface{}) (string, error) {
	// Convert schema to JSON string if needed
	var schemaJSON string
	switch v := schema.(type) {
	case string:
		schemaJSON = v
	case []byte:
		schemaJSON = string(v)
	default:
		schemaBytes, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal schema: %w", err)
		}
		schemaJSON = string(schemaBytes)
	}

	var promptBuilder strings.Builder
	
	promptBuilder.WriteString("The previously generated JSON failed validation. Please fix the issues and generate valid JSON.\n\n")
	
	promptBuilder.WriteString("## Validation Error\n")
	promptBuilder.WriteString("The following validation error occurred:\n")
	promptBuilder.WriteString("```\n")
	promptBuilder.WriteString(validationError)
	promptBuilder.WriteString("\n```\n\n")
	
	promptBuilder.WriteString("## Invalid JSON\n")
	promptBuilder.WriteString("This was the invalid JSON that was generated:\n")
	promptBuilder.WriteString("```json\n")
	promptBuilder.WriteString(invalidJSON)
	promptBuilder.WriteString("\n```\n\n")
	
	promptBuilder.WriteString("## Required Schema\n")
	promptBuilder.WriteString("The JSON MUST conform to this schema:\n")
	promptBuilder.WriteString("```json\n")
	promptBuilder.WriteString(schemaJSON)
	promptBuilder.WriteString("\n```\n\n")
	
	promptBuilder.WriteString("## Original Context\n")
	promptBuilder.WriteString(originalPrompt)
	promptBuilder.WriteString("\n\n")
	
	promptBuilder.WriteString("## Repair Instructions\n")
	promptBuilder.WriteString("1. Identify the validation error in the JSON\n")
	promptBuilder.WriteString("2. Fix the specific issue mentioned in the error\n")
	promptBuilder.WriteString("3. Ensure the repaired JSON matches the schema exactly\n")
	promptBuilder.WriteString("4. Preserve all valid content from the original JSON\n")
	promptBuilder.WriteString("5. Return ONLY the repaired JSON, no additional text\n")

	return promptBuilder.String(), nil
}

// EstimateTokens provides a rough estimate of the number of tokens in a prompt.
// This is a simple heuristic and not exact.
func (b *Builder) EstimateTokens(prompt string) int {
	// Rough estimate: ~1 token per 4 characters or ~0.75 tokens per word
	// Using character-based estimation for simplicity
	return len(prompt) / 4
}