package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/karolswdev/docloom/internal/version"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Long:  `Display detailed version information about DocLoom including build metadata.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.Info())
	},
}

var versionFlag bool

func init() {
	rootCmd.AddCommand(versionCmd)

	// Also add --version flag to root command
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "", false, "Display version information")
}
