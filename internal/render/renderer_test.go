package render

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TC-4.1: Golden file test for HTML rendering
func TestRenderer_RenderHTML_Golden(t *testing.T) {
	// Arrange
	// Read the template HTML
	templatePath := filepath.Join("testdata", "architecture-vision.html")
	templateBytes, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}
	
	// Read the fields JSON
	fieldsPath := filepath.Join("testdata", "fields.json")
	fieldsBytes, err := os.ReadFile(fieldsPath)
	if err != nil {
		t.Fatalf("Failed to read fields file: %v", err)
	}
	
	// Parse the fields
	var fields map[string]interface{}
	if unmarshalErr := json.Unmarshal(fieldsBytes, &fields); unmarshalErr != nil {
		t.Fatalf("Failed to parse fields JSON: %v", unmarshalErr)
	}
	
	// Read the expected output (golden file)
	expectedPath := filepath.Join("testdata", "expected.html")
	expectedBytes, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("Failed to read expected output file: %v", err)
	}
	expected := string(expectedBytes)
	
	// Act
	rendered, err := RenderHTML(string(templateBytes), fields)
	
	// Assert
	if err != nil {
		t.Fatalf("RenderHTML failed: %v", err)
	}
	
	if rendered != expected {
		// For debugging, write the actual output to a file
		actualPath := filepath.Join("testdata", "actual.html")
		os.WriteFile(actualPath, []byte(rendered), 0644)
		
		t.Errorf("Rendered HTML does not match expected golden file\n"+
			"Expected output saved to: %s\n"+
			"Actual output saved to: %s\n"+
			"Diff:\nExpected length: %d\nActual length: %d",
			expectedPath, actualPath, len(expected), len(rendered))
		
		// Show first difference
		for i := 0; i < len(expected) && i < len(rendered); i++ {
			if expected[i] != rendered[i] {
				t.Errorf("First difference at position %d: expected %q, got %q",
					i, expected[max(0, i-20):min(len(expected), i+20)],
					rendered[max(0, i-20):min(len(rendered), i+20)])
				break
			}
		}
	}
	
	// Test that the renderer also saves JSON sidecar
	t.Run("JSON sidecar saved", func(t *testing.T) {
		// Create a temporary directory for output
		tmpDir := t.TempDir()
		outputPath := filepath.Join(tmpDir, "output.html")
		
		// Create renderer and render
		renderer := NewRenderer(tmpDir)
		err := renderer.Render(string(templateBytes), fields, outputPath)
		
		if err != nil {
			t.Fatalf("Render failed: %v", err)
		}
		
		// Check that HTML file exists
		if _, statErr := os.Stat(outputPath); os.IsNotExist(statErr) {
			t.Error("HTML output file was not created")
		}
		
		// Check that JSON sidecar exists
		jsonPath := filepath.Join(tmpDir, "output.json")
		if _, statErr := os.Stat(jsonPath); os.IsNotExist(statErr) {
			t.Error("JSON sidecar file was not created")
		}
		
		// Verify JSON content
		jsonBytes, err := os.ReadFile(jsonPath)
		if err != nil {
			t.Fatalf("Failed to read JSON sidecar: %v", err)
		}
		
		var savedFields map[string]interface{}
		if err := json.Unmarshal(jsonBytes, &savedFields); err != nil {
			t.Fatalf("Failed to parse saved JSON: %v", err)
		}
		
		// Check that essential fields are present
		if doc, ok := savedFields["document"].(map[string]interface{}); ok {
			if title, ok := doc["title"].(string); !ok || title != "NextGen Platform Architecture Vision" {
				t.Error("JSON sidecar missing or incorrect document.title")
			}
		} else {
			t.Error("JSON sidecar missing document field")
		}
	})
}

// Test pure rendering function
func TestRenderHTML_Pure(t *testing.T) {
	tests := []struct {
		name     string
		template string
		fields   map[string]interface{}
		expected string
	}{
		{
			name:     "Simple replacement",
			template: `<h1><!-- data-field="title" --></h1>`,
			fields:   map[string]interface{}{"title": "Hello World"},
			expected: `<h1>Hello World</h1>`,
		},
		{
			name:     "Nested field",
			template: `<p><!-- data-field="user.name" --></p>`,
			fields: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "John Doe",
				},
			},
			expected: `<p>John Doe</p>`,
		},
		{
			name:     "Multiple replacements",
			template: `<div><!-- data-field="first" --> and <!-- data-field="second" --></div>`,
			fields: map[string]interface{}{
				"first":  "One",
				"second": "Two",
			},
			expected: `<div>One and Two</div>`,
		},
		{
			name:     "Missing field leaves placeholder",
			template: `<span><!-- data-field="missing" --></span>`,
			fields:   map[string]interface{}{},
			expected: `<span><!-- data-field="missing" --></span>`,
		},
		{
			name:     "Number field",
			template: `<span>Year: <!-- data-field="year" --></span>`,
			fields:   map[string]interface{}{"year": 2025},
			expected: `<span>Year: 2025</span>`,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RenderHTML(tt.template, tt.fields)
			if err != nil {
				t.Fatalf("RenderHTML failed: %v", err)
			}
			
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}

// Test idempotency - rendering the same inputs always produces the same output
func TestRenderHTML_Idempotent(t *testing.T) {
	template := `<h1><!-- data-field="title" --></h1><p><!-- data-field="content" --></p>`
	fields := map[string]interface{}{
		"title":   "Test Document",
		"content": "This is test content",
	}
	
	// Render multiple times
	results := make([]string, 5)
	for i := 0; i < 5; i++ {
		result, err := RenderHTML(template, fields)
		if err != nil {
			t.Fatalf("RenderHTML failed on iteration %d: %v", i, err)
		}
		results[i] = result
	}
	
	// All results should be identical
	for i := 1; i < len(results); i++ {
		if results[i] != results[0] {
			t.Errorf("Rendering is not idempotent. Result %d differs from result 0", i)
		}
	}
}

// Test RenderToWriter
func TestRenderToWriter(t *testing.T) {
	template := `<h1><!-- data-field="title" --></h1>`
	fields := map[string]interface{}{"title": "Test Title"}
	
	htmlBuf := &bytes.Buffer{}
	jsonBuf := &bytes.Buffer{}
	
	err := RenderToWriter(template, fields, htmlBuf, jsonBuf)
	if err != nil {
		t.Fatalf("RenderToWriter failed: %v", err)
	}
	
	// Check HTML output
	expectedHTML := `<h1>Test Title</h1>`
	if htmlBuf.String() != expectedHTML {
		t.Errorf("HTML output mismatch. Expected: %s, Got: %s", expectedHTML, htmlBuf.String())
	}
	
	// Check JSON output
	var savedFields map[string]interface{}
	if err := json.Unmarshal(jsonBuf.Bytes(), &savedFields); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}
	
	if title, ok := savedFields["title"].(string); !ok || title != "Test Title" {
		t.Error("JSON output missing or incorrect title field")
	}
}

// Helper functions
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}