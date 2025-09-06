// Package generate orchestrates the document generation workflow.
package generate

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/karolswdev/docloom/internal/ai"
	"github.com/karolswdev/docloom/internal/ingest"
	"github.com/karolswdev/docloom/internal/prompt"
	"github.com/karolswdev/docloom/internal/render"
	"github.com/karolswdev/docloom/internal/templates"
	"github.com/karolswdev/docloom/internal/validate"
	"github.com/rs/zerolog/log"
)

// Options contains configuration for the generation process.
type Options struct {
	TemplateType string
	Sources      []string
	OutputFile   string
	
	// AI Configuration
	Model       string
	BaseURL     string
	APIKey      string
	Temperature float32
	Seed        *int
	MaxRetries  int
	
	// Behavior
	DryRun      bool
	Force       bool
	MaxRepairs  int
}

// Orchestrator coordinates the document generation workflow.
type Orchestrator struct {
	aiClient  ai.Client
	ingester  *ingest.Ingester
	builder   *prompt.Builder
	validator *validate.Validator
	registry  *templates.Registry
	renderer  *render.Renderer
	outputDir string
}

// NewOrchestrator creates a new generation orchestrator.
func NewOrchestrator(aiClient ai.Client) *Orchestrator {
	registry := templates.NewRegistry()
	registry.LoadDefaults() // Load default templates
	
	return &Orchestrator{
		aiClient:  aiClient,
		ingester:  ingest.NewIngester(),
		builder:   prompt.NewBuilder(),
		validator: validate.NewValidator(),
		registry:  registry,
		renderer:  render.NewRenderer("output"), // Default output directory
		outputDir: "output",
	}
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
	sourceContent, err := o.ingester.IngestSources(opts.Sources)
	if err != nil {
		return fmt.Errorf("failed to ingest sources: %w", err)
	}
	log.Info().Int("bytes", len(sourceContent)).Msg("Source ingestion complete")

	// Step 2: Build generation prompt
	log.Info().Msg("Building generation prompt")
	generationPrompt, err := o.builder.BuildGenerationPrompt(sourceContent, tmpl.Prompt, tmpl.Schema)
	if err != nil {
		return fmt.Errorf("failed to build prompt: %w", err)
	}
	
	if opts.DryRun {
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
		schemaBytes, _ := json.MarshalIndent(tmpl.Schema, "", "  ")
		fmt.Println(string(schemaBytes))
		return nil
	}

	// Step 3: Generate with validation and repair loop
	var generatedJSON string
	var lastError error
	maxAttempts := opts.MaxRepairs + 1 // Initial attempt + repairs
	
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		var currentPrompt string
		
		if attempt == 1 {
			// First attempt - use the original prompt
			currentPrompt = generationPrompt
			log.Info().Msg("Calling AI model for initial generation")
		} else {
			// Repair attempt - build repair prompt
			log.Info().
				Int("attempt", attempt).
				Int("max_attempts", maxAttempts).
				Msg("Validation failed, attempting repair")
			
			repairPrompt, err := o.builder.BuildRepairPrompt(
				generationPrompt, 
				generatedJSON, 
				lastError.Error(),
				tmpl.Schema,
			)
			if err != nil {
				return fmt.Errorf("failed to build repair prompt: %w", err)
			}
			currentPrompt = repairPrompt
		}

		// Call AI model
		startTime := time.Now()
		generatedJSON, err = o.aiClient.GenerateJSON(ctx, currentPrompt)
		if err != nil {
			return fmt.Errorf("AI generation failed: %w", err)
		}
		duration := time.Since(startTime)
		
		log.Info().
			Dur("duration", duration).
			Int("response_bytes", len(generatedJSON)).
			Msg("Received AI response")

		// Validate the generated JSON
		schemaStr, err := json.Marshal(tmpl.Schema)
		if err != nil {
			return fmt.Errorf("failed to marshal schema: %w", err)
		}
		
		validationErr := o.validator.Validate(generatedJSON, string(schemaStr))
		if validationErr == nil {
			// Validation passed!
			log.Info().Msg("JSON validation successful")
			break
		}
		
		// Validation failed
		lastError = validationErr
		log.Warn().
			Err(validationErr).
			Int("attempt", attempt).
			Msg("JSON validation failed")
		
		if attempt == maxAttempts {
			return fmt.Errorf("failed to generate valid JSON after %d attempts: %w", maxAttempts, lastError)
		}
	}

	// Step 4: Save JSON sidecar file
	jsonFile := strings.TrimSuffix(opts.OutputFile, ".html") + ".json"
	if err := os.WriteFile(jsonFile, []byte(generatedJSON), 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}
	log.Info().Str("file", jsonFile).Msg("Saved JSON sidecar file")

	// Step 5: Render HTML output
	log.Info().Msg("Rendering HTML output")
	
	// Parse the JSON into a map for rendering
	var fields map[string]interface{}
	if err := json.Unmarshal([]byte(generatedJSON), &fields); err != nil {
		return fmt.Errorf("failed to parse generated JSON: %w", err)
	}
	
	// Use the renderer to render and save both HTML and JSON
	if err := o.renderer.Render(tmpl.HTMLContent, fields, opts.OutputFile); err != nil {
		return fmt.Errorf("failed to render output: %w", err)
	}
	
	log.Info().
		Str("html_file", opts.OutputFile).
		Str("json_file", jsonFile).
		Msg("Document generation complete")
	
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