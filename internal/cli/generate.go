package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/karolswdev/docloom/internal/agent"
	"github.com/karolswdev/docloom/internal/ai"
	"github.com/karolswdev/docloom/internal/generate"
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
	agentName    string
	agentParams  []string
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a document from sources and template",
	Long: `Generate a complete document by combining provided source materials 
with a selected template type and AI-generated content mapped to that template's field schema.

Example:
  docloom generate --type architecture-vision --source ./docs --out output.html
  docloom generate --agent research-agent --source ./repo --type report --out analysis.html`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
		
		// If agent is specified, run it first
		actualSources := sources
		if agentName != "" {
			// Parse agent parameters
			params := make(map[string]string)
			for _, param := range agentParams {
				parts := strings.SplitN(param, "=", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid agent parameter format: %s (expected key=value)", param)
				}
				params[parts[0]] = parts[1]
			}

			// Create agent registry and discover agents
			registry := agent.NewRegistry()
			if err := registry.Discover(); err != nil {
				return fmt.Errorf("failed to discover agents: %w", err)
			}

			// Create artifact cache
			cache, err := agent.NewArtifactCache()
			if err != nil {
				return fmt.Errorf("failed to create artifact cache: %w", err)
			}

			// Create executor
			executor := agent.NewExecutor(registry, cache, logger)

			// Prepare source path (use first source or current directory)
			sourcePath := "."
			if len(sources) > 0 {
				sourcePath = sources[0]
			}

			// Run the agent
			fmt.Printf("Running agent '%s' on source: %s\n", agentName, sourcePath)
			result, err := executor.Run(agent.RunOptions{
				AgentName:  agentName,
				SourcePath: sourcePath,
				Parameters: params,
			})
			if err != nil {
				return fmt.Errorf("agent execution failed: %w", err)
			}

			// Validate agent output
			if err := executor.ValidateOutput(result.OutputPath); err != nil {
				return fmt.Errorf("agent output validation failed: %w", err)
			}

			// Replace sources with agent output directory
			actualSources = []string{result.OutputPath}
			fmt.Printf("Agent completed. Using artifacts from: %s\n", result.OutputPath)
		}
		
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
			Sources:      actualSources,
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

	// Agent flags
	generateCmd.Flags().StringVar(&agentName, "agent", "", "Research agent to run before generation")
	generateCmd.Flags().StringSliceVar(&agentParams, "agent-param", []string{}, "Agent parameters (format: key=value, can be specified multiple times)")

	// Mark required flags
	if err := generateCmd.MarkFlagRequired("type"); err != nil {
		panic(fmt.Sprintf("failed to mark type flag as required: %v", err))
	}
	if err := generateCmd.MarkFlagRequired("out"); err != nil {
		panic(fmt.Sprintf("failed to mark out flag as required: %v", err))
	}
}
