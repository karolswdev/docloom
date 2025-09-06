package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// templatesCmd represents the templates command
var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "Manage document templates",
	Long:  `Commands for managing and listing available document templates.`,
}

// listCmd represents the templates list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available templates",
	Long:  `List all available document templates that can be used with the generate command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Placeholder list of templates
		templates := []struct {
			Name        string
			Description string
		}{
			{
				Name:        "architecture-vision",
				Description: "Architecture Vision document template",
			},
			{
				Name:        "technical-debt-summary",
				Description: "Technical Debt Summary template",
			},
			{
				Name:        "reference-architecture",
				Description: "Reference Architecture template",
			},
		}

		fmt.Println("Available templates:")
		fmt.Println()
		for _, tmpl := range templates {
			fmt.Printf("  %s\n    %s\n\n", tmpl.Name, tmpl.Description)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(templatesCmd)
	templatesCmd.AddCommand(listCmd)
}
