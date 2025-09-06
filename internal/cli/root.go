package cli

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	verbose bool
	logger  zerolog.Logger
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "docloom",
	Short: "Beautiful, template-driven technical documentation — fast",
	Long: `Docloom is a system for technical folks to generate high-quality documents 
by combining structured templates with source materials and model-assisted content. 
The aim is consistent, branded, and reviewable outputs that you can print, share, 
and iterate on quickly.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Configure logging based on verbose flag
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		
		if verbose {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}
		
		// Human-readable console output
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).
			With().
			Timestamp().
			Logger()
		
		log.Logger = logger
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Persistent flags available to all commands
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
}

// GetRootCmd returns the root command for testing
func GetRootCmd() *cobra.Command {
	return rootCmd
}

// GetLogger returns the configured logger
func GetLogger() zerolog.Logger {
	return logger
}