// Package validate provides JSON schema validation functionality.
package validate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

// Validator handles JSON schema validation.
type Validator struct{}

// NewValidator creates a new JSON schema validator.
func NewValidator() *Validator {
	return &Validator{}
}

// Validate checks if a JSON string conforms to the provided schema.
func (v *Validator) Validate(jsonStr string, schemaStr string) error {
	// Parse the JSON to validate
	var jsonData interface{}
	if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	// Create a new compiler for this validation
	compiler := jsonschema.NewCompiler()
	compiler.Draft = jsonschema.Draft7

	// Add the schema to the compiler using a reader
	schemaReader := bytes.NewReader([]byte(schemaStr))
	if err := compiler.AddResource("schema.json", schemaReader); err != nil {
		return fmt.Errorf("failed to add schema resource: %w", err)
	}

	// Compile the schema
	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return fmt.Errorf("failed to compile schema: %w", err)
	}

	// Validate the JSON against the schema
	if err := schema.Validate(jsonData); err != nil {
		return newValidationError(err)
	}

	return nil
}

// ValidationError provides detailed information about validation failures.
type ValidationError struct {
	Message  string
	Field    string
	Expected string
	Actual   string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation failed at field '%s': %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation failed: %s", e.Message)
}

// newValidationError converts a jsonschema validation error into our ValidationError.
func newValidationError(err error) error {
	if err == nil {
		return nil
	}

	// The jsonschema library returns detailed validation errors
	// We'll extract the relevant information
	errStr := err.Error()

	// Try to extract field path and message
	validationErr := &ValidationError{
		Message: errStr,
	}

	// The error format is typically like:
	// "jsonschema: '/field/path' does not validate with schema: error details"
	// or "doesn't validate with" followed by the error

	if strings.Contains(errStr, "doesn't validate with") {
		parts := strings.SplitN(errStr, "doesn't validate with", 2)
		if len(parts) >= 1 {
			// Extract the field path
			fieldPart := strings.TrimSpace(parts[0])
			if strings.HasPrefix(fieldPart, "jsonschema: '") {
				fieldPart = strings.TrimPrefix(fieldPart, "jsonschema: '")
				if idx := strings.Index(fieldPart, "'"); idx > 0 {
					validationErr.Field = fieldPart[:idx]
				}
			}
		}
		if len(parts) >= 2 {
			// Extract the error message
			validationErr.Message = strings.TrimSpace(parts[1])
		}
	}

	return validationErr
}

// ValidateWithDetails performs validation and returns detailed error information.
func (v *Validator) ValidateWithDetails(jsonStr string, schemaStr string) (*ValidationResult, error) {
	result := &ValidationResult{
		Valid: true,
	}

	// Parse the JSON
	var jsonData interface{}
	if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationIssue{
			Type:    "parse_error",
			Message: fmt.Sprintf("Invalid JSON: %v", err),
		})
		return result, nil
	}

	// Create a new compiler for this validation
	compiler := jsonschema.NewCompiler()
	compiler.Draft = jsonschema.Draft7

	// Add the schema to the compiler using a reader
	schemaReader := bytes.NewReader([]byte(schemaStr))
	if err := compiler.AddResource("schema.json", schemaReader); err != nil {
		return nil, fmt.Errorf("failed to add schema resource: %w", err)
	}

	// Compile the schema
	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return nil, fmt.Errorf("failed to compile schema: %w", err)
	}

	// Validate
	if err := schema.Validate(jsonData); err != nil {
		result.Valid = false
		// Extract validation issues
		issue := ValidationIssue{
			Type:    "validation_error",
			Message: err.Error(),
		}

		// Try to extract field information
		if validErr, ok := err.(*ValidationError); ok {
			issue.Field = validErr.Field
		}

		result.Errors = append(result.Errors, issue)
	}

	return result, nil
}

// ValidationResult contains the outcome of a validation check.
type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationIssue `json:"errors,omitempty"`
}

// ValidationIssue represents a single validation problem.
type ValidationIssue struct {
	Type    string `json:"type"`
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}
