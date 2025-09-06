package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestValidator_Validate_ValidJSON tests validation of valid JSON against a schema.
func TestValidator_Validate_ValidJSON(t *testing.T) {
	// Arrange: Provide a valid JSON string and a corresponding JSON schema
	validJSON := `{
		"title": "Test Document",
		"summary": "This is a test summary",
		"sections": [
			{
				"heading": "Introduction",
				"content": "This is the introduction"
			}
		]
	}`

	schema := `{
		"type": "object",
		"properties": {
			"title": {
				"type": "string"
			},
			"summary": {
				"type": "string"
			},
			"sections": {
				"type": "array",
				"items": {
					"type": "object",
					"properties": {
						"heading": {
							"type": "string"
						},
						"content": {
							"type": "string"
						}
					},
					"required": ["heading", "content"]
				}
			}
		},
		"required": ["title", "summary", "sections"]
	}`

	validator := NewValidator()

	// Act: Run the validation
	err := validator.Validate(validJSON, schema)

	// Assert: The function returns no error
	assert.NoError(t, err)
}

// TestValidator_Validate_InvalidJSON tests validation of invalid JSON against a schema.
func TestValidator_Validate_InvalidJSON(t *testing.T) {
	// Arrange: Provide a JSON string with a type error (number where string is expected)
	invalidJSON := `{
		"title": 12345,
		"summary": "This is a test summary",
		"sections": []
	}`

	schema := `{
		"type": "object",
		"properties": {
			"title": {
				"type": "string"
			},
			"summary": {
				"type": "string"
			},
			"sections": {
				"type": "array"
			}
		},
		"required": ["title", "summary", "sections"]
	}`

	validator := NewValidator()

	// Act: Run validation
	err := validator.Validate(invalidJSON, schema)

	// Assert: The function returns a specific validation error
	require.Error(t, err)
	assert.Contains(t, err.Error(), "title")
	// The error should indicate the type mismatch
}

// TestValidator_Validate_MissingRequiredField tests validation when required fields are missing.
func TestValidator_Validate_MissingRequiredField(t *testing.T) {
	// Arrange: JSON missing a required field
	incompleteJSON := `{
		"title": "Test Document"
	}`

	schema := `{
		"type": "object",
		"properties": {
			"title": {
				"type": "string"
			},
			"summary": {
				"type": "string"
			}
		},
		"required": ["title", "summary"]
	}`

	validator := NewValidator()

	// Act
	err := validator.Validate(incompleteJSON, schema)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "summary")
}

// TestValidator_Validate_InvalidJSONSyntax tests handling of malformed JSON.
func TestValidator_Validate_InvalidJSONSyntax(t *testing.T) {
	// Arrange: Malformed JSON
	malformedJSON := `{
		"title": "Test",
		"summary": 
	}`

	schema := `{
		"type": "object",
		"properties": {
			"title": {"type": "string"},
			"summary": {"type": "string"}
		}
	}`

	validator := NewValidator()

	// Act
	err := validator.Validate(malformedJSON, schema)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid JSON")
}

// TestValidator_Validate_ComplexSchema tests validation with a more complex schema.
func TestValidator_Validate_ComplexSchema(t *testing.T) {
	// Arrange
	complexJSON := `{
		"name": "Project Alpha",
		"version": "1.0.0",
		"metadata": {
			"author": "John Doe",
			"tags": ["documentation", "technical"],
			"priority": 5
		},
		"components": [
			{
				"id": "comp-1",
				"type": "service",
				"config": {
					"port": 8080,
					"enabled": true
				}
			}
		]
	}`

	complexSchema := `{
		"type": "object",
		"properties": {
			"name": {"type": "string"},
			"version": {"type": "string", "pattern": "^\\d+\\.\\d+\\.\\d+$"},
			"metadata": {
				"type": "object",
				"properties": {
					"author": {"type": "string"},
					"tags": {
						"type": "array",
						"items": {"type": "string"}
					},
					"priority": {
						"type": "integer",
						"minimum": 1,
						"maximum": 10
					}
				},
				"required": ["author", "tags"]
			},
			"components": {
				"type": "array",
				"items": {
					"type": "object",
					"properties": {
						"id": {"type": "string"},
						"type": {"type": "string", "enum": ["service", "library", "tool"]},
						"config": {
							"type": "object",
							"properties": {
								"port": {"type": "integer"},
								"enabled": {"type": "boolean"}
							}
						}
					},
					"required": ["id", "type"]
				}
			}
		},
		"required": ["name", "version", "metadata"]
	}`

	validator := NewValidator()

	// Act
	err := validator.Validate(complexJSON, complexSchema)

	// Assert
	assert.NoError(t, err)
}

// TestValidator_Validate_EnumViolation tests validation when an enum constraint is violated.
func TestValidator_Validate_EnumViolation(t *testing.T) {
	// Arrange
	jsonWithInvalidEnum := `{
		"status": "unknown"
	}`

	schemaWithEnum := `{
		"type": "object",
		"properties": {
			"status": {
				"type": "string",
				"enum": ["pending", "approved", "rejected"]
			}
		},
		"required": ["status"]
	}`

	validator := NewValidator()

	// Act
	err := validator.Validate(jsonWithInvalidEnum, schemaWithEnum)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "status")
}

// TestValidator_ValidateWithDetails tests the detailed validation result.
func TestValidator_ValidateWithDetails(t *testing.T) {
	// Arrange
	invalidJSON := `{
		"title": 123,
		"summary": "Test"
	}`

	schema := `{
		"type": "object",
		"properties": {
			"title": {"type": "string"},
			"summary": {"type": "string"}
		},
		"required": ["title", "summary"]
	}`

	validator := NewValidator()

	// Act
	result, err := validator.ValidateWithDetails(invalidJSON, schema)

	// Assert
	require.NoError(t, err) // No error in calling the function
	assert.False(t, result.Valid)
	assert.NotEmpty(t, result.Errors)
	assert.Contains(t, result.Errors[0].Message, "title")
}

// TestValidator_Validate_AdditionalProperties tests handling of additional properties.
func TestValidator_Validate_AdditionalProperties(t *testing.T) {
	// Test with additionalProperties: false
	jsonWithExtra := `{
		"title": "Test",
		"summary": "Summary",
		"extra": "This should not be here"
	}`

	strictSchema := `{
		"type": "object",
		"properties": {
			"title": {"type": "string"},
			"summary": {"type": "string"}
		},
		"required": ["title", "summary"],
		"additionalProperties": false
	}`

	validator := NewValidator()

	// Should fail with strict schema
	err := validator.Validate(jsonWithExtra, strictSchema)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "additional")

	// Test with additionalProperties: true (or not specified)
	lenientSchema := `{
		"type": "object",
		"properties": {
			"title": {"type": "string"},
			"summary": {"type": "string"}
		},
		"required": ["title", "summary"]
	}`

	// Should pass with lenient schema
	err = validator.Validate(jsonWithExtra, lenientSchema)
	assert.NoError(t, err)
}