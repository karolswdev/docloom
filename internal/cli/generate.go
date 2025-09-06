package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/karolswdev/docloom/internal/ai"
	"github.com/karolswdev/docloom/internal/generate"
	"github.com/spf13/cobra"
)

var (
	templateType string
	sources      []string
	outputFile   string
	model        string
	baseURL      string
	apiKey       string
	temperature  float64
	seed         int
	maxRetries   int
	dryRun       bool
	force        bool
	configFile   string
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a document from sources and template",
	Long: `Generate a complete document by combining provided source materials 
with a selected template type and AI-generated content mapped to that template's field schema.

Example:
  docloom generate --type architecture-vision --source ./docs --out output.html`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get API key from flag or environment
		if apiKey == "" {
			apiKey = os.Getenv("OPENAI_API_KEY")
			if apiKey == "" {
				apiKey = os.Getenv("DOCLOOM_API_KEY")
			}
		}

		// For dry-run, we don't need to create a real AI client
		var aiClient ai.Client
		if !dryRun {
			// Create AI client configuration
			aiConfig := ai.Config{
				BaseURL:     baseURL,
				APIKey:      apiKey,
				Model:       model,
				Temperature: float32(temperature),
				MaxTokens:   4096,
				MaxRetries:  maxRetries,
			}

			if seed > 0 {
				aiConfig.Seed = &seed
			}

			// Create AI client
			var err error
			aiClient, err = ai.NewOpenAIClient(aiConfig)
			if err != nil {
				return fmt.Errorf("failed to create AI client: %w", err)
			}
		}

		// Create orchestrator
		orchestrator := generate.NewOrchestrator(aiClient)

		// Prepare options
		opts := generate.Options{
			TemplateType: templateType,
			Sources:      sources,
			OutputFile:   outputFile,
			Model:        model,
			BaseURL:      baseURL,
			APIKey:       apiKey,
			Temperature:  float32(temperature),
			MaxRetries:   maxRetries,
			DryRun:       dryRun,
			Force:        force,
			MaxRepairs:   3, // Default to 3 repair attempts
		}

		if seed > 0 {
			opts.Seed = &seed
		}

		// Run generation
		ctx := context.Background()
		if err := orchestrator.Generate(ctx, opts); err != nil {
			return err
		}

		if !dryRun {
			fmt.Printf("Successfully generated document: %s\n", outputFile)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Required flags
	generateCmd.Flags().StringVarP(&templateType, "type", "t", "", "Template type to use (required)")
	generateCmd.Flags().StringSliceVarP(&sources, "source", "s", []string{}, "Source paths (files or directories)")
	generateCmd.Flags().StringVarP(&outputFile, "out", "o", "", "Output file path (required)")

	// Model configuration flags
	generateCmd.Flags().StringVar(&model, "model", "gpt-4", "Model to use for generation")
	generateCmd.Flags().StringVar(&baseURL, "base-url", "", "Base URL for OpenAI-compatible API")
	generateCmd.Flags().StringVar(&apiKey, "api-key", "", "API key (can also use OPENAI_API_KEY env var)")
	generateCmd.Flags().Float64Var(&temperature, "temperature", 0.7, "Temperature for model generation")
	generateCmd.Flags().IntVar(&seed, "seed", 0, "Seed for reproducible generation")
	generateCmd.Flags().IntVar(&maxRetries, "retries", 3, "Maximum number of retries for model calls")

	// Operational flags
	generateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview without making API calls")
	generateCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing output files")
	generateCmd.Flags().StringVar(&configFile, "config", "", "Config file path")

	// Mark required flags
	if err := generateCmd.MarkFlagRequired("type"); err != nil {
		panic(fmt.Sprintf("failed to mark type flag as required: %v", err))
	}
	if err := generateCmd.MarkFlagRequired("out"); err != nil {
		panic(fmt.Sprintf("failed to mark out flag as required: %v", err))
	}
}
