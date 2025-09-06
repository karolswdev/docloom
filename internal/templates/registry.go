package templates

import (
	// "embed" // Currently using programmatic templates
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

// Analysis contains prompts for AI-driven analysis
type Analysis struct {
	SystemPrompt      string `json:"system_prompt"`
	InitialUserPrompt string `json:"initial_user_prompt"`
}

// Template represents a document template with its assets
type Template struct {
	Assets       map[string][]byte `json:"-"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	HTMLContent  string            `json:"-"`
	HTMLTemplate string            `json:"-"` // For compatibility
	Prompt       string            `json:"prompt"`
	Schema       json.RawMessage   `json:"schema"`
	FieldSchema  json.RawMessage   `json:"-"` // Alias for Schema
	Analysis     *Analysis         `json:"analysis,omitempty"`
}

// Registry manages available templates
type Registry struct {
	templates map[string]*Template
}

// Embed default templates into the binary
// Note: Currently using programmatic templates
// //go:embed defaults/*
// var defaultTemplatesFS embed.FS

// NewRegistry creates a new template registry
func NewRegistry() *Registry {
	return &Registry{
		templates: make(map[string]*Template),
	}
}

// LoadDefaults loads the default embedded templates
func (r *Registry) LoadDefaults() error {
	log.Debug().Msg("Loading default embedded templates")

	// For now, we'll create sample templates programmatically
	// In production, these would be loaded from the embedded FS
	r.registerDefaultTemplates()

	return nil
}

// LoadFromDirectory loads templates from a directory
func (r *Registry) LoadFromDirectory(dir string) error {
	log.Debug().Str("dir", dir).Msg("Loading templates from directory")

	// Walk the directory looking for template definitions
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-template files
		if info.IsDir() {
			return nil
		}

		// Look for template directories with required files
		if strings.HasSuffix(path, "template.json") {
			templateDir := filepath.Dir(path)
			r.loadTemplate(templateDir)
			// TODO: When file reading is implemented, handle errors here
		}

		return nil
	})

	return err
}

// loadTemplate loads a single template from a directory
func (r *Registry) loadTemplate(dir string) {
	templateName := filepath.Base(dir)

	// Initialize template
	tmpl := &Template{
		Name:   templateName,
		Assets: make(map[string][]byte),
	}

	// Check for required files
	htmlPath := filepath.Join(dir, templateName+".html")
	schemaPath := filepath.Join(dir, "schema.json")
	promptPath := filepath.Join(dir, "prompt.txt")

	// For testing purposes, we'll accept templates even if files don't exist yet
	// In production, we would read these files
	log.Debug().
		Str("name", templateName).
		Str("html", htmlPath).
		Str("schema", schemaPath).
		Str("prompt", promptPath).
		Msg("Loading template")

	// Register the template
	r.templates[templateName] = tmpl

	// TODO: Return errors when file reading is implemented
}

// registerDefaultTemplates registers the built-in templates
func (r *Registry) registerDefaultTemplates() {
	// Architecture Vision template
	r.templates["architecture-vision"] = &Template{
		Name:        "architecture-vision",
		Description: "Architecture Vision document template",
		HTMLContent: architectureVisionHTML,
		Schema:      json.RawMessage(architectureVisionSchema),
		Prompt:      architectureVisionPrompt,
		Assets:      make(map[string][]byte),
	}

	// Technical Debt Summary template
	r.templates["technical-debt-summary"] = &Template{
		Name:        "technical-debt-summary",
		Description: "Technical Debt Summary template",
		HTMLContent: technicalDebtHTML,
		Schema:      json.RawMessage(technicalDebtSchema),
		Prompt:      technicalDebtPrompt,
		Assets:      make(map[string][]byte),
	}

	// Reference Architecture template
	r.templates["reference-architecture"] = &Template{
		Name:        "reference-architecture",
		Description: "Reference Architecture template",
		HTMLContent: referenceArchHTML,
		Schema:      json.RawMessage(referenceArchSchema),
		Prompt:      referenceArchPrompt,
		Assets:      make(map[string][]byte),
	}
}

// Get retrieves a template by name
func (r *Registry) Get(name string) (*Template, error) {
	tmpl, exists := r.templates[name]
	if !exists {
		return nil, fmt.Errorf("template '%s' not found", name)
	}
	// Set compatibility fields
	if tmpl.HTMLTemplate == "" {
		tmpl.HTMLTemplate = tmpl.HTMLContent
	}
	if tmpl.FieldSchema == nil {
		tmpl.FieldSchema = tmpl.Schema
	}
	return tmpl, nil
}

// Load is an alias for Get for compatibility
func (r *Registry) Load(name string) (*Template, error) {
	return r.Get(name)
}

// Register adds a new template to the registry
func (r *Registry) Register(name string, tmpl *Template) error {
	if _, exists := r.templates[name]; exists {
		return fmt.Errorf("template '%s' already exists", name)
	}
	r.templates[name] = tmpl
	return nil
}

// List returns all available template names
func (r *Registry) List() []string {
	names := make([]string, 0, len(r.templates))
	for name := range r.templates {
		names = append(names, name)
	}
	return names
}

// ListWithDescriptions returns all templates with their descriptions
func (r *Registry) ListWithDescriptions() []struct{ Name, Description string } {
	result := make([]struct{ Name, Description string }, 0, len(r.templates))
	for _, tmpl := range r.templates {
		result = append(result, struct{ Name, Description string }{
			Name:        tmpl.Name,
			Description: tmpl.Description,
		})
	}
	return result
}

// Sample template content (simplified for testing)
const (
	architectureVisionHTML = `<!DOCTYPE html>
<html>
<head><title>Architecture Vision</title></head>
<body>
<!-- data-field="document.title" -->
<!-- data-field="document.content" -->
</body>
</html>`

	architectureVisionSchema = `{
  "type": "object",
  "properties": {
    "document": {
      "type": "object",
      "properties": {
        "title": {"type": "string"},
        "content": {"type": "string"}
      }
    }
  }
}`

	architectureVisionPrompt = `Generate an architecture vision document based on the provided sources.`

	technicalDebtHTML = `<!DOCTYPE html>
<html>
<head><title>Technical Debt Summary</title></head>
<body>
<!-- data-field="summary.title" -->
<!-- data-field="summary.items" -->
</body>
</html>`

	technicalDebtSchema = `{
  "type": "object",
  "properties": {
    "summary": {
      "type": "object",
      "properties": {
        "title": {"type": "string"},
        "items": {"type": "array"}
      }
    }
  }
}`

	technicalDebtPrompt = `Analyze technical debt from the provided sources.`

	referenceArchHTML = `<!DOCTYPE html>
<html>
<head><title>Reference Architecture</title></head>
<body>
<!-- data-field="architecture.name" -->
<!-- data-field="architecture.components" -->
</body>
</html>`

	referenceArchSchema = `{
  "type": "object",
  "properties": {
    "architecture": {
      "type": "object",
      "properties": {
        "name": {"type": "string"},
        "components": {"type": "array"}
      }
    }
  }
}`

	referenceArchPrompt = `Create a reference architecture based on the provided sources.`
)
