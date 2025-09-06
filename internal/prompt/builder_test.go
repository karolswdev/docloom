package prompt

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewBuilder tests the creation of a new Builder instance
func TestNewBuilder(t *testing.T) {
	builder := NewBuilder()
	assert.NotNil(t, builder, "NewBuilder should return a non-nil Builder")
}

// TestBuildGenerationPrompt tests the generation prompt building functionality
func TestBuildGenerationPrompt(t *testing.T) {
	builder := NewBuilder()

	tests := []struct {
		validatePrompt func(t *testing.T, prompt string)
		schema         interface{}
		name           string
		sourceContent  string
		templatePrompt string
		expectError    bool
	}{
		{
			name:           "Basic prompt generation with map schema",
			sourceContent:  "This is the source content for the document.",
			templatePrompt: "Generate a technical summary",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title":   map[string]interface{}{"type": "string"},
					"summary": map[string]interface{}{"type": "string"},
				},
				"required": []string{"title", "summary"},
			},
			expectError: false,
			validatePrompt: func(t *testing.T, prompt string) {
				assert.Contains(t, prompt, "You are a technical documentation generator")
				assert.Contains(t, prompt, "Generate a technical summary")
				assert.Contains(t, prompt, "This is the source content for the document.")
				assert.Contains(t, prompt, "\"title\"")
				assert.Contains(t, prompt, "\"summary\"")
				assert.Contains(t, prompt, "## Template Instructions")
				assert.Contains(t, prompt, "## JSON Schema")
				assert.Contains(t, prompt, "## Source Documents")
				assert.Contains(t, prompt, "## Instructions")
			},
		},
		{
			name:           "Prompt generation with string schema",
			sourceContent:  "Source data",
			templatePrompt: "Template instructions",
			schema:         `{"type": "object", "properties": {"field": {"type": "string"}}}`,
			expectError:    false,
			validatePrompt: func(t *testing.T, prompt string) {
				assert.Contains(t, prompt, `"type": "object"`)
				assert.Contains(t, prompt, `"field"`)
				assert.Contains(t, prompt, "Source data")
				assert.Contains(t, prompt, "Template instructions")
			},
		},
		{
			name:           "Prompt generation with byte slice schema",
			sourceContent:  "Test content",
			templatePrompt: "Test template",
			schema:         []byte(`{"type": "object"}`),
			expectError:    false,
			validatePrompt: func(t *testing.T, prompt string) {
				assert.Contains(t, prompt, `"type": "object"`)
				assert.Contains(t, prompt, "Test content")
				assert.Contains(t, prompt, "Test template")
			},
		},
		{
			name:           "Empty source content",
			sourceContent:  "",
			templatePrompt: "Generate something",
			schema:         map[string]interface{}{"type": "object"},
			expectError:    false,
			validatePrompt: func(t *testing.T, prompt string) {
				assert.Contains(t, prompt, "Generate something")
				assert.Contains(t, prompt, "## Source Documents")
			},
		},
		{
			name:           "Complex nested schema",
			sourceContent:  "Complex document",
			templatePrompt: "Process complex data",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"metadata": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"author": map[string]interface{}{"type": "string"},
							"date":   map[string]interface{}{"type": "string"},
						},
					},
					"sections": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"title":   map[string]interface{}{"type": "string"},
								"content": map[string]interface{}{"type": "string"},
							},
						},
					},
				},
			},
			expectError: false,
			validatePrompt: func(t *testing.T, prompt string) {
				assert.Contains(t, prompt, "metadata")
				assert.Contains(t, prompt, "sections")
				assert.Contains(t, prompt, "author")
				assert.Contains(t, prompt, "Complex document")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt, err := builder.BuildGenerationPrompt(tt.sourceContent, tt.templatePrompt, tt.schema)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, prompt)

				// Validate the structure of the prompt
				assert.True(t, strings.HasPrefix(prompt, "You are a technical documentation generator"))

				// Run custom validation if provided
				if tt.validatePrompt != nil {
					tt.validatePrompt(t, prompt)
				}
			}
		})
	}
}

// TestBuildRepairPrompt tests the repair prompt building functionality
func TestBuildRepairPrompt(t *testing.T) {
	builder := NewBuilder()

	tests := []struct {
		validatePrompt  func(t *testing.T, prompt string)
		schema          interface{}
		name            string
		originalPrompt  string
		invalidJSON     string
		validationError string
		expectError     bool
	}{
		{
			name:            "Basic repair prompt",
			originalPrompt:  "Original generation prompt",
			invalidJSON:     `{"title": 123}`,
			validationError: "title should be string, not number",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{"type": "string"},
				},
			},
			expectError: false,
			validatePrompt: func(t *testing.T, prompt string) {
				assert.Contains(t, prompt, "The previously generated JSON failed validation")
				assert.Contains(t, prompt, "title should be string, not number")
				assert.Contains(t, prompt, `{"title": 123}`)
				assert.Contains(t, prompt, "Original generation prompt")
				assert.Contains(t, prompt, "## Validation Error")
				assert.Contains(t, prompt, "## Invalid JSON")
				assert.Contains(t, prompt, "## Required Schema")
				assert.Contains(t, prompt, "## Original Context")
				assert.Contains(t, prompt, "## Repair Instructions")
			},
		},
		{
			name:            "Repair with string schema",
			originalPrompt:  "Test prompt",
			invalidJSON:     `{"field": null}`,
			validationError: "field cannot be null",
			schema:          `{"type": "object", "properties": {"field": {"type": "string"}}}`,
			expectError:     false,
			validatePrompt: func(t *testing.T, prompt string) {
				assert.Contains(t, prompt, "field cannot be null")
				assert.Contains(t, prompt, `{"field": null}`)
				assert.Contains(t, prompt, `"type": "object"`)
			},
		},
		{
			name:            "Repair with byte slice schema",
			originalPrompt:  "Original",
			invalidJSON:     `{}`,
			validationError: "missing required field: title",
			schema:          []byte(`{"type": "object", "required": ["title"]}`),
			expectError:     false,
			validatePrompt: func(t *testing.T, prompt string) {
				assert.Contains(t, prompt, "missing required field: title")
				assert.Contains(t, prompt, `"required": ["title"]`)
			},
		},
		{
			name:            "Complex validation error",
			originalPrompt:  "Complex prompt",
			invalidJSON:     `{"sections": "should be array"}`,
			validationError: "sections: expected array, got string",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"sections": map[string]interface{}{
						"type": "array",
					},
				},
			},
			expectError: false,
			validatePrompt: func(t *testing.T, prompt string) {
				assert.Contains(t, prompt, "sections: expected array, got string")
				assert.Contains(t, prompt, `"sections": "should be array"`)
				assert.Contains(t, prompt, "Identify the validation error")
				assert.Contains(t, prompt, "Fix the specific issue")
				assert.Contains(t, prompt, "Preserve all valid content")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt, err := builder.BuildRepairPrompt(
				tt.originalPrompt,
				tt.invalidJSON,
				tt.validationError,
				tt.schema,
			)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, prompt)

				// Run custom validation if provided
				if tt.validatePrompt != nil {
					tt.validatePrompt(t, prompt)
				}
			}
		})
	}
}

// TestEstimateTokens tests the token estimation functionality
func TestEstimateTokens(t *testing.T) {
	builder := NewBuilder()

	tests := []struct {
		name        string
		prompt      string
		expectedMin int
		expectedMax int
	}{
		{
			name:        "Empty prompt",
			prompt:      "",
			expectedMin: 0,
			expectedMax: 0,
		},
		{
			name:        "Short prompt",
			prompt:      "Hello world",
			expectedMin: 2, // 11 chars / 4 = 2.75, rounded down
			expectedMax: 3,
		},
		{
			name:        "Medium prompt",
			prompt:      strings.Repeat("This is a test. ", 10),
			expectedMin: 35, // 160 chars / 4 = 40
			expectedMax: 45,
		},
		{
			name:        "Long prompt",
			prompt:      strings.Repeat("Generate a comprehensive technical document. ", 50),
			expectedMin: 500, // 2250 chars / 4 = 562
			expectedMax: 600,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := builder.EstimateTokens(tt.prompt)

			assert.GreaterOrEqual(t, tokens, tt.expectedMin,
				"Token count should be at least %d", tt.expectedMin)
			assert.LessOrEqual(t, tokens, tt.expectedMax,
				"Token count should be at most %d", tt.expectedMax)
		})
	}
}

// TestPromptStructure tests that prompts have the expected structure
func TestPromptStructure(t *testing.T) {
	builder := NewBuilder()

	t.Run("Generation prompt structure", func(t *testing.T) {
		prompt, err := builder.BuildGenerationPrompt(
			"Source content",
			"Template prompt",
			map[string]interface{}{"type": "object"},
		)
		require.NoError(t, err)

		// Check that all expected sections are present and in order
		sections := []string{
			"## Template Instructions",
			"## JSON Schema",
			"## Source Documents",
			"## Instructions",
		}

		lastIndex := -1
		for _, section := range sections {
			index := strings.Index(prompt, section)
			assert.Greater(t, index, lastIndex,
				"Section '%s' should appear after previous sections", section)
			lastIndex = index
		}

		// Check that JSON is properly formatted
		assert.Contains(t, prompt, "```json")
		assert.Contains(t, prompt, "```")
	})

	t.Run("Repair prompt structure", func(t *testing.T) {
		prompt, err := builder.BuildRepairPrompt(
			"Original prompt",
			`{"invalid": true}`,
			"Validation failed",
			map[string]interface{}{"type": "object"},
		)
		require.NoError(t, err)

		// Check that all expected sections are present and in order
		sections := []string{
			"## Validation Error",
			"## Invalid JSON",
			"## Required Schema",
			"## Original Context",
			"## Repair Instructions",
		}

		lastIndex := -1
		for _, section := range sections {
			index := strings.Index(prompt, section)
			assert.Greater(t, index, lastIndex,
				"Section '%s' should appear after previous sections", section)
			lastIndex = index
		}
	})
}

// TestSchemaMarshaling tests that different schema types are handled correctly
func TestSchemaMarshaling(t *testing.T) {
	builder := NewBuilder()

	t.Run("Invalid schema object", func(t *testing.T) {
		// Create a schema that cannot be marshaled to JSON
		type invalidType struct {
			Channel chan int `json:"channel"`
		}

		schema := invalidType{Channel: make(chan int)}

		_, err := builder.BuildGenerationPrompt(
			"content",
			"prompt",
			schema,
		)

		// Should get an error when trying to marshal the channel
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to marshal schema")
	})

	t.Run("Valid complex schema", func(t *testing.T) {
		// Complex but valid schema
		schema := map[string]interface{}{
			"$schema": "http://json-schema.org/draft-07/schema#",
			"type":    "object",
			"properties": map[string]interface{}{
				"id":     map[string]interface{}{"type": "integer"},
				"name":   map[string]interface{}{"type": "string", "minLength": 1},
				"email":  map[string]interface{}{"type": "string", "format": "email"},
				"active": map[string]interface{}{"type": "boolean"},
				"tags":   map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
			},
			"required": []string{"id", "name"},
		}

		prompt, err := builder.BuildGenerationPrompt(
			"User data",
			"Generate user profile",
			schema,
		)

		require.NoError(t, err)
		assert.Contains(t, prompt, `"$schema"`)
		assert.Contains(t, prompt, `"format": "email"`)
		assert.Contains(t, prompt, `"minLength": 1`)
	})

	t.Run("JSON RawMessage schema", func(t *testing.T) {
		// Test with json.RawMessage (common in real usage)
		rawSchema := json.RawMessage(`{"type": "object", "properties": {"test": {"type": "string"}}}`)

		prompt, err := builder.BuildGenerationPrompt(
			"content",
			"prompt",
			rawSchema,
		)

		require.NoError(t, err)
		assert.Contains(t, prompt, `"test"`)
		assert.Contains(t, prompt, `"type": "string"`)
	})
}

// TestPromptConsistency tests that prompts are consistent across calls
func TestPromptConsistency(t *testing.T) {
	builder := NewBuilder()

	sourceContent := "Test content"
	templatePrompt := "Test template"
	schema := map[string]interface{}{"type": "object"}

	// Generate the same prompt multiple times
	prompt1, err1 := builder.BuildGenerationPrompt(sourceContent, templatePrompt, schema)
	prompt2, err2 := builder.BuildGenerationPrompt(sourceContent, templatePrompt, schema)

	require.NoError(t, err1)
	require.NoError(t, err2)

	// Prompts should be identical for the same inputs
	assert.Equal(t, prompt1, prompt2, "Prompts should be consistent for the same inputs")
}

// BenchmarkBuildGenerationPrompt benchmarks prompt generation
func BenchmarkBuildGenerationPrompt(b *testing.B) {
	builder := NewBuilder()
	sourceContent := strings.Repeat("This is test content. ", 100)
	templatePrompt := "Generate a technical document"
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"title":   map[string]interface{}{"type": "string"},
			"summary": map[string]interface{}{"type": "string"},
			"sections": map[string]interface{}{
				"type":  "array",
				"items": map[string]interface{}{"type": "string"},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = builder.BuildGenerationPrompt(sourceContent, templatePrompt, schema)
	}
}

// BenchmarkEstimateTokens benchmarks token estimation
func BenchmarkEstimateTokens(b *testing.B) {
	builder := NewBuilder()
	prompt := strings.Repeat("This is a test prompt for token estimation. ", 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = builder.EstimateTokens(prompt)
	}
}
