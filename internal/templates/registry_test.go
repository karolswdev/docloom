package templates

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TC-3.1: Test that the template registry can load templates
// validateTemplate validates a template has all required fields
func validateTemplate(t *testing.T, tmpl *Template, expectedName string) {
	t.Helper()

	if tmpl == nil {
		t.Fatal("Template should not be nil")
	}

	if tmpl.Name != expectedName {
		t.Errorf("Expected template name '%s', got '%s'", expectedName, tmpl.Name)
	}

	if tmpl.Description == "" {
		t.Errorf("Template '%s' missing description", expectedName)
	}

	if tmpl.HTMLContent == "" {
		t.Errorf("Template '%s' missing HTML content", expectedName)
	}

	if len(tmpl.Schema) == 0 {
		t.Errorf("Template '%s' missing schema", expectedName)
	}

	// Validate the schema is valid JSON
	var schemaObj interface{}
	if err := json.Unmarshal(tmpl.Schema, &schemaObj); err != nil {
		t.Errorf("Template '%s' schema should be valid JSON: %v", expectedName, err)
	}

	if tmpl.Prompt == "" {
		t.Errorf("Template '%s' missing prompt", expectedName)
	}
}

func TestTemplateRegistry_LoadDefaults(t *testing.T) {
	// Arrange
	registry := NewRegistry()

	// Act
	err := registry.LoadDefaults()

	// Assert
	if err != nil {
		t.Fatalf("Failed to load default templates: %v", err)
	}

	// Check that architecture-vision template is loaded
	tmpl, err := registry.Get("architecture-vision")
	if err != nil {
		t.Errorf("Expected 'architecture-vision' template to be loaded, got error: %v", err)
	}

	validateTemplate(t, tmpl, "architecture-vision")
}

func TestTemplateRegistry_AllRequiredTemplatesPresent(t *testing.T) {
	// Arrange
	registry := NewRegistry()
	registry.LoadDefaults()

	requiredTemplates := []string{
		"architecture-vision",
		"technical-debt-summary",
		"reference-architecture",
	}

	// Act & Assert
	for _, name := range requiredTemplates {
		tmpl, err := registry.Get(name)
		if err != nil {
			t.Errorf("Required template '%s' not found: %v", name, err)
			continue
		}
		validateTemplate(t, tmpl, name)
	}
}

func TestTemplateRegistry_LoadFromDirectory(t *testing.T) {
	// Arrange - create a temporary test directory with template files
	tmpDir, err := os.MkdirTemp("", "docloom-test-templates")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test template directory
	testTemplateDir := filepath.Join(tmpDir, "test-template")
	if mkdirErr := os.MkdirAll(testTemplateDir, 0755); mkdirErr != nil {
		t.Fatalf("Failed to create test template dir: %v", mkdirErr)
	}

	// Create template.json to mark it as a template
	templateJSON := filepath.Join(testTemplateDir, "template.json")
	if writeErr := os.WriteFile(templateJSON, []byte(`{"name": "test-template"}`), 0644); writeErr != nil {
		t.Fatalf("Failed to create template.json: %v", writeErr)
	}

	// Create other template files
	htmlFile := filepath.Join(testTemplateDir, "test-template.html")
	if writeErr := os.WriteFile(htmlFile, []byte("<html></html>"), 0644); writeErr != nil {
		t.Fatalf("Failed to create HTML file: %v", writeErr)
	}

	schemaFile := filepath.Join(testTemplateDir, "schema.json")
	if schemaErr := os.WriteFile(schemaFile, []byte(`{"type": "object"}`), 0644); schemaErr != nil {
		t.Fatalf("Failed to create schema file: %v", schemaErr)
	}

	promptFile := filepath.Join(testTemplateDir, "prompt.txt")
	if promptErr := os.WriteFile(promptFile, []byte("Test prompt"), 0644); promptErr != nil {
		t.Fatalf("Failed to create prompt file: %v", promptErr)
	}

	registry := NewRegistry()

	// Act
	err = registry.LoadFromDirectory(tmpDir)

	// Assert
	if err != nil {
		t.Fatalf("Failed to load templates from directory: %v", err)
	}

	// Check that the test template was loaded
	tmpl, err := registry.Get("test-template")
	if err != nil {
		t.Errorf("Expected 'test-template' to be loaded, got error: %v", err)
	}

	if tmpl == nil {
		t.Fatal("Template should not be nil")
	}

	if tmpl.Name != "test-template" {
		t.Errorf("Expected template name 'test-template', got '%s'", tmpl.Name)
	}
}

func TestTemplateRegistry_List(t *testing.T) {
	// Arrange
	registry := NewRegistry()
	registry.LoadDefaults()

	// Act
	templates := registry.List()

	// Assert
	if len(templates) < 3 {
		t.Errorf("Expected at least 3 templates, got %d", len(templates))
	}

	// Check that expected templates are in the list
	expectedTemplates := map[string]bool{
		"architecture-vision":    false,
		"technical-debt-summary": false,
		"reference-architecture": false,
	}

	for _, name := range templates {
		if _, ok := expectedTemplates[name]; ok {
			expectedTemplates[name] = true
		}
	}

	for name, found := range expectedTemplates {
		if !found {
			t.Errorf("Expected template '%s' not found in list", name)
		}
	}
}

func TestTemplateRegistry_GetNonExistent(t *testing.T) {
	// Arrange
	registry := NewRegistry()
	registry.LoadDefaults()

	// Act
	tmpl, err := registry.Get("non-existent-template")

	// Assert
	if err == nil {
		t.Error("Expected error when getting non-existent template")
	}

	if tmpl != nil {
		t.Error("Template should be nil for non-existent template")
	}
}
