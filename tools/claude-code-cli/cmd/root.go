package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	repoPath   string
	outputPath string
	verbose    bool
	apiKey     string
	model      string
	maxTokens  int
)

var rootCmd = &cobra.Command{
	Use:   "cc-cli",
	Short: "Claude Code CLI - Analyze C# repositories using Claude LLM",
	Long: `The Claude Code CLI (cc-cli) is a powerful tool that uses the Claude LLM 
to perform deep analysis of C# repositories. It generates structured artifacts 
that provide rich context for automated document generation.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Setup logging
		if verbose {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}
		
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	},
	Run: analyze,
}

func init() {
	rootCmd.Flags().StringVar(&repoPath, "repo-path", "", "Path to the C# repository to analyze (required)")
	rootCmd.Flags().StringVar(&outputPath, "output-path", "", "Path where artifacts will be written (required)")
	rootCmd.Flags().BoolVar(&verbose, "verbose", false, "Enable verbose output")
	rootCmd.Flags().StringVar(&apiKey, "api-key", "", "Claude API key (can also be set via CLAUDE_API_KEY env var)")
	rootCmd.Flags().StringVar(&model, "model", "claude-3-opus-20240229", "Claude model to use")
	rootCmd.Flags().IntVar(&maxTokens, "max-tokens", 4096, "Maximum tokens for Claude response")
	
	rootCmd.MarkFlagRequired("repo-path")
	rootCmd.MarkFlagRequired("output-path")
}

func Execute() error {
	return rootCmd.Execute()
}

func analyze(cmd *cobra.Command, args []string) {
	// Get API key from environment if not provided
	if apiKey == "" {
		apiKey = os.Getenv("CLAUDE_API_KEY")
		if apiKey == "" {
			log.Fatal().Msg("Claude API key is required. Set via --api-key flag or CLAUDE_API_KEY environment variable")
		}
	}
	
	log.Info().
		Str("repo", repoPath).
		Str("output", outputPath).
		Str("model", model).
		Msg("Starting C# repository analysis")
	
	// Initialize scanner
	scanner := NewScanner(repoPath)
	
	// Scan repository
	scanResult, err := scanner.Scan()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to scan repository")
	}
	
	log.Info().
		Int("files_found", len(scanResult.Files)).
		Msg("Repository scan complete")
	
	// Generate prompt
	promptGen := NewPromptGenerator()
	prompt := promptGen.Generate(scanResult)
	
	log.Debug().Str("prompt_preview", fmt.Sprintf("%.200s...", prompt)).Msg("Generated analysis prompt")
	
	// Call Claude API
	client := NewClaudeClient(apiKey, model, maxTokens)
	response, err := client.Analyze(prompt)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get response from Claude")
	}
	
	log.Info().Msg("Received analysis from Claude")
	
	// Write artifacts
	writer := NewArtifactWriter(outputPath)
	if err := writer.Write(response); err != nil {
		log.Fatal().Err(err).Msg("Failed to write artifacts")
	}
	
	log.Info().
		Str("output", outputPath).
		Msg("Successfully wrote analysis artifacts")
}