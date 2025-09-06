package render

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

// Renderer handles rendering HTML templates with field data
type Renderer struct {
	outputDir string
}

// NewRenderer creates a new renderer instance
func NewRenderer(outputDir string) *Renderer {
	return &Renderer{
		outputDir: outputDir,
	}
}

// RenderHTML takes an HTML template and field data, replacing placeholders with actual values
// This function is pure - it has no side effects other than returning the rendered string
func RenderHTML(htmlTemplate string, fields map[string]interface{}) (string, error) {
	// Create a flat map of field paths to values
	flatFields := flattenMap(fields, "")
	
	// Regular expression to match data-field comments
	// Matches patterns like: <!-- data-field="document.title" -->
	fieldPattern := regexp.MustCompile(`<!--\s*data-field="([^"]+)"\s*-->`)
	
	// Replace each placeholder with its corresponding value
	rendered := fieldPattern.ReplaceAllStringFunc(htmlTemplate, func(match string) string {
		// Extract the field path from the match
		matches := fieldPattern.FindStringSubmatch(match)
		if len(matches) < 2 {
			return match // Return unchanged if no field path found
		}
		
		fieldPath := matches[1]
		
		// Look up the value in our flattened fields
		if value, exists := flatFields[fieldPath]; exists {
			// Convert the value to string
			switch v := value.(type) {
			case string:
				return v
			case []byte:
				return string(v)
			default:
				// For other types, use JSON encoding for proper representation
				jsonBytes, err := json.Marshal(v)
				if err != nil {
					log.Warn().Err(err).Str("field", fieldPath).Msg("Failed to marshal field value")
					return match
				}
				// If it's a string in JSON, remove the quotes
				str := string(jsonBytes)
				if strings.HasPrefix(str, `"`) && strings.HasSuffix(str, `"`) {
					str = str[1 : len(str)-1]
				}
				return str
			}
		}
		
		log.Debug().Str("field", fieldPath).Msg("Field not found in data, leaving placeholder")
		return match // Return unchanged if field not found
	})
	
	return rendered, nil
}

// Render renders an HTML template with the given fields and saves both HTML and JSON outputs
func (r *Renderer) Render(templateHTML string, fields map[string]interface{}, outputPath string) error {
	// Render the HTML
	renderedHTML, err := RenderHTML(templateHTML, fields)
	if err != nil {
		return fmt.Errorf("failed to render HTML: %w", err)
	}
	
	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if mkdirErr := os.MkdirAll(outputDir, 0755); mkdirErr != nil {
		return fmt.Errorf("failed to create output directory: %w", mkdirErr)
	}
	
	// Write the rendered HTML
	if err := os.WriteFile(outputPath, []byte(renderedHTML), 0600); err != nil {
		return fmt.Errorf("failed to write HTML output: %w", err)
	}
	
	// Generate JSON sidecar path (same name, .json extension)
	jsonPath := strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".json"
	
	// Marshal fields to JSON
	jsonData, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal fields to JSON: %w", err)
	}
	
	// Write the JSON sidecar
	if err := os.WriteFile(jsonPath, jsonData, 0600); err != nil {
		return fmt.Errorf("failed to write JSON sidecar: %w", err)
	}
	
	log.Info().
		Str("html", outputPath).
		Str("json", jsonPath).
		Msg("Successfully rendered template")
	
	return nil
}

// RenderFromFiles renders using file paths instead of content
func (r *Renderer) RenderFromFiles(templatePath string, fieldsPath string, outputPath string) error {
	// Read the HTML template
	templateBytes, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}
	
	// Read the fields JSON
	fieldsBytes, err := os.ReadFile(fieldsPath)
	if err != nil {
		return fmt.Errorf("failed to read fields file: %w", err)
	}
	
	// Parse the fields JSON
	var fields map[string]interface{}
	if err := json.Unmarshal(fieldsBytes, &fields); err != nil {
		return fmt.Errorf("failed to parse fields JSON: %w", err)
	}
	
	// Render using the main Render function
	return r.Render(string(templateBytes), fields, outputPath)
}

// RenderToWriter renders the template and writes to the provided writers
func RenderToWriter(htmlTemplate string, fields map[string]interface{}, htmlWriter, jsonWriter io.Writer) error {
	// Render the HTML
	renderedHTML, err := RenderHTML(htmlTemplate, fields)
	if err != nil {
		return fmt.Errorf("failed to render HTML: %w", err)
	}
	
	// Write the rendered HTML
	if _, err := htmlWriter.Write([]byte(renderedHTML)); err != nil {
		return fmt.Errorf("failed to write HTML: %w", err)
	}
	
	// Marshal and write the JSON
	encoder := json.NewEncoder(jsonWriter)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(fields); err != nil {
		return fmt.Errorf("failed to write JSON: %w", err)
	}
	
	return nil
}

// flattenMap flattens a nested map into a single-level map with dot-separated keys
func flattenMap(m map[string]interface{}, prefix string) map[string]interface{} {
	result := make(map[string]interface{})
	
	for key, value := range m {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}
		
		switch v := value.(type) {
		case map[string]interface{}:
			// Recursively flatten nested maps
			for k, val := range flattenMap(v, fullKey) {
				result[k] = val
			}
		default:
			// Add non-map values directly
			result[fullKey] = value
		}
	}
	
	return result
}