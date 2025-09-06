// Package generate orchestrates the document generation workflow.
package generate

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/karolswdev/docloom/internal/agent"
	"github.com/karolswdev/docloom/internal/ai"
	"github.com/karolswdev/docloom/internal/ingest"
	"github.com/karolswdev/docloom/internal/prompt"
	"github.com/karolswdev/docloom/internal/render"
	"github.com/karolswdev/docloom/internal/templates"
	"github.com/karolswdev/docloom/internal/validate"
)

// Options contains configuration for the generation process.
type Options struct {
	Seed         *int
	TemplateType string
	OutputFile   string
	Model        string
	BaseURL      string
	APIKey       string
	Sources      []string
	MaxRetries   int
	MaxRepairs   int
	Temperature  float32
	DryRun       bool
	Force        bool
}

// Orchestrator coordinates the document generation workflow.
type Orchestrator struct {
	aiClient      ai.Client
	ingester      *ingest.Ingester
	builder       *prompt.Builder
	validator     *validate.Validator
	registry      *templates.Registry
	renderer      *render.Renderer
	outputDir     string
	agentRegistry *agent.Registry
	agentExecutor *agent.Executor
}

// NewOrchestrator creates a new generation orchestrator.
func NewOrchestrator(aiClient ai.Client) *Orchestrator {
	registry := templates.NewRegistry()
	if err := registry.LoadDefaults(); err != nil {
		// Log warning but continue - templates can be loaded later
		log.Warn().Err(err).Msg("Failed to load default templates")
	}

	// Initialize agent support
	agentRegistry := agent.NewRegistry()
	agentCache, _ := agent.NewArtifactCache()
	agentExecutor := agent.NewExecutor(agentRegistry, agentCache, log.Logger)

	return &Orchestrator{
		aiClient:      aiClient,
		ingester:      ingest.NewIngester(),
		builder:       prompt.NewBuilder(),
		validator:     validate.NewValidator(),
		registry:      registry,
		renderer:      render.NewRenderer("output"), // Default output directory
		outputDir:     "output",
		agentRegistry: agentRegistry,
		agentExecutor: agentExecutor,
	}
}

// generateWithRetries attempts to generate valid JSON with retries
func (o *Orchestrator) generateWithRetries(ctx context.Context, generationPrompt string, tmpl *templates.Template, opts Options) (string, error) {
	var generatedJSON string
	var lastError error
	maxAttempts := opts.MaxRepairs + 1 // Initial attempt + repairs

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		var currentPrompt string

		if attempt == 1 {
			currentPrompt = generationPrompt
			log.Info().Msg("Calling AI model for initial generation")
			log.Debug().Str("model", opts.Model).Float32("temperature", opts.Temperature).Msg("Model parameters")
		} else {
			// Build repair prompt
			log.Info().Int("attempt", attempt).Int("max_attempts", maxAttempts).Msg("Attempting repair")
			repairPrompt, err := o.builder.BuildRepairPrompt(generationPrompt, generatedJSON, lastError.Error(), tmpl.Schema)
			if err != nil {
				return "", fmt.Errorf("failed to build repair prompt: %w", err)
			}
			currentPrompt = repairPrompt
		}

		// Call AI model
		startTime := time.Now()
		generatedJSON, err := o.aiClient.GenerateJSON(ctx, currentPrompt)
		if err != nil {
			return "", fmt.Errorf("AI generation failed: %w", err)
		}
		log.Info().Dur("duration", time.Since(startTime)).Int("response_bytes", len(generatedJSON)).Msg("Received AI response")

		// Validate the generated JSON
		schemaStr, schemaErr := json.Marshal(tmpl.Schema)
		if schemaErr != nil {
			return "", fmt.Errorf("failed to marshal schema: %w", schemaErr)
		}

		validationErr := o.validator.Validate(generatedJSON, string(schemaStr))
		if validationErr == nil {
			log.Info().Msg("JSON validation successful")
			return generatedJSON, nil
		}

		lastError = validationErr
		log.Warn().Err(validationErr).Int("attempt", attempt).Msg("JSON validation failed")
	}

	return "", fmt.Errorf("failed to generate valid JSON after %d attempts: %w", maxAttempts, lastError)
}

// handleDryRun prints dry-run information and returns
func (o *Orchestrator) handleDryRun(opts Options, tmpl *templates.Template, generationPrompt string) error {
	fmt.Println("\n=== DRY RUN MODE ===")
	fmt.Printf("Template: %s\n", opts.TemplateType)
	fmt.Printf("Sources: %v\n", opts.Sources)
	fmt.Printf("Output: %s\n", opts.OutputFile)
	fmt.Printf("Model: %s\n", opts.Model)
	fmt.Printf("Estimated tokens: %d\n", o.builder.EstimateTokens(generationPrompt))
	fmt.Println("\n=== PROMPT PREVIEW (first 1000 chars) ===")
	if len(generationPrompt) > 1000 {
		fmt.Println(generationPrompt[:1000] + "...")
	} else {
		fmt.Println(generationPrompt)
	}
	fmt.Println("\n=== SCHEMA ===")
	schemaBytes, schemaErr := json.MarshalIndent(tmpl.Schema, "", "  ")
	if schemaErr != nil {
		log.Warn().Err(schemaErr).Msg("Failed to marshal schema for display")
		schemaBytes = []byte("{}")
	}
	fmt.Println(string(schemaBytes))
	return nil
}

// Generate performs the complete document generation workflow.
func (o *Orchestrator) Generate(ctx context.Context, opts Options) error {
	// Validate options
	if err := o.validateOptions(opts); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}

	// Check if output file exists and handle force flag
	if !opts.Force {
		if _, err := os.Stat(opts.OutputFile); err == nil {
			return fmt.Errorf("output file %s already exists (use --force to overwrite)", opts.OutputFile)
		}
	}

	// Get template from registry
	tmpl, err := o.registry.Get(opts.TemplateType)
	if err != nil {
		return fmt.Errorf("failed to get template %s: %w", opts.TemplateType, err)
	}

	// Step 1: Ingest source documents
	log.Info().Strs("sources", opts.Sources).Msg("Ingesting source documents")
	log.Debug().Str("template", opts.TemplateType).Msg("Using template for generation")
	log.Debug().Str("model", opts.Model).Msg("Selected AI model")
	log.Debug().Int("max_repairs", opts.MaxRepairs).Msg("Maximum repair attempts configured")
	sourceContent, err := o.ingester.IngestSources(opts.Sources)
	if err != nil {
		return fmt.Errorf("failed to ingest sources: %w", err)
	}
	log.Info().Int("bytes", len(sourceContent)).Msg("Source ingestion complete")
	log.Debug().Int("source_files", len(opts.Sources)).Msg("Total source files processed")

	// Step 2: Build generation prompt
	log.Info().Msg("Building generation prompt")
	log.Debug().Str("template_prompt", tmpl.Prompt[:min(100, len(tmpl.Prompt))]).Msg("Template prompt preview")
	generationPrompt, err := o.builder.BuildGenerationPrompt(sourceContent, tmpl.Prompt, tmpl.Schema)
	if err != nil {
		return fmt.Errorf("failed to build prompt: %w", err)
	}
	log.Debug().Int("prompt_length", len(generationPrompt)).Msg("Generation prompt built")

	if opts.DryRun {
		return o.handleDryRun(opts, tmpl, generationPrompt)
	}

	// Step 3: Generate with validation and repair loop
	generatedJSON, err := o.generateWithRetries(ctx, generationPrompt, tmpl, opts)
	if err != nil {
		return err
	}

	// Step 4: Save JSON sidecar file
	jsonFile := strings.TrimSuffix(opts.OutputFile, ".html") + ".json"
	if err := os.WriteFile(jsonFile, []byte(generatedJSON), 0600); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}
	log.Info().Str("file", jsonFile).Msg("Saved JSON sidecar file")

	// Step 5: Render HTML output
	log.Info().Msg("Rendering HTML output")
	log.Debug().Msg("Parsing generated JSON for rendering")

	// Parse the JSON into a map for rendering
	var fields map[string]interface{}
	if err := json.Unmarshal([]byte(generatedJSON), &fields); err != nil {
		return fmt.Errorf("failed to parse generated JSON: %w", err)
	}
	log.Debug().Int("field_count", len(fields)).Msg("Parsed JSON fields")

	// Use the renderer to render and save both HTML and JSON
	log.Debug().Str("output_file", opts.OutputFile).Msg("Writing rendered HTML")
	if err := o.renderer.Render(tmpl.HTMLContent, fields, opts.OutputFile); err != nil {
		return fmt.Errorf("failed to render output: %w", err)
	}

	log.Info().
		Str("html_file", opts.OutputFile).
		Str("json_file", jsonFile).
		Msg("Document generation complete")
	log.Debug().Msg("Generation workflow completed successfully")

	return nil
}

// writeOutput writes the generated HTML and JSON to files.
func (o *Orchestrator) writeOutput(outputFile string, htmlContent string, fields map[string]interface{}, force bool) error {
	// Check if output file exists and force flag
	if !force {
		if _, err := os.Stat(outputFile); err == nil {
			return fmt.Errorf("output file '%s' already exists (use --force to overwrite)", outputFile)
		}
	}

	// Write HTML file
	if err := os.WriteFile(outputFile, []byte(htmlContent), 0644); err != nil {
		return fmt.Errorf("failed to write HTML file: %w", err)
	}

	// Write JSON sidecar file
	jsonFile := strings.TrimSuffix(outputFile, filepath.Ext(outputFile)) + ".json"
	jsonData, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	if err := os.WriteFile(jsonFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	return nil
}

// validateOptions checks that all required options are provided.
func (o *Orchestrator) validateOptions(opts Options) error {
	if opts.TemplateType == "" {
		return fmt.Errorf("template type is required")
	}
	if len(opts.Sources) == 0 {
		return fmt.Errorf("at least one source is required")
	}
	if opts.OutputFile == "" {
		return fmt.Errorf("output file is required")
	}
	if !opts.DryRun && opts.APIKey == "" {
		// Check environment variable
		if os.Getenv("OPENAI_API_KEY") == "" {
			return fmt.Errorf("API key is required (use --api-key or OPENAI_API_KEY env var)")
		}
	}
	if opts.MaxRepairs < 0 {
		return fmt.Errorf("max repairs must be non-negative")
	}
	return nil
}
