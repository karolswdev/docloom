package cli

import (
	"fmt"

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
		// For now, just print a message indicating the command would run
		if dryRun {
			fmt.Println("Dry-run mode: would generate document with the following settings:")
			fmt.Printf("  Template: %s\n", templateType)
			fmt.Printf("  Sources: %v\n", sources)
			fmt.Printf("  Output: %s\n", outputFile)
			fmt.Printf("  Model: %s\n", model)
			if baseURL != "" {
				fmt.Printf("  Base URL: %s\n", baseURL)
			}
			fmt.Printf("  Temperature: %.2f\n", temperature)
			if seed > 0 {
				fmt.Printf("  Seed: %d\n", seed)
			}
			fmt.Printf("  Max retries: %d\n", maxRetries)
			return nil
		}
		
		// Placeholder for actual generation logic
		return fmt.Errorf("generate command not yet implemented")
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
	generateCmd.MarkFlagRequired("type")
	generateCmd.MarkFlagRequired("out")
}