package cli

import (
	"fmt"
	"text/tabwriter"

	"github.com/karolswdev/docloom/internal/agent"
	"github.com/spf13/cobra"
)

// agentsCmd represents the agents command
var agentsCmd = &cobra.Command{
	Use:   "agents",
	Short: "Manage and inspect Research Agents",
	Long: `The agents command provides tools for managing Research Agents.

Research Agents are external programs that can analyze your codebase and
generate specialized documentation. Use the subcommands to list available
agents and inspect their details.`,
}

// agentsListCmd represents the agents list command
var agentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available Research Agents",
	Long: `List all Research Agents discovered in the search paths.

Agents are discovered from:
  - Project-local: .docloom/agents/
  - User-home: ~/.docloom/agents/
  
The list shows each agent's name and description.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create and populate registry
		registry := agent.NewRegistry()
		if err := registry.Discover(); err != nil {
			return fmt.Errorf("failed to discover agents: %w", err)
		}

		// Get all agents
		agents := registry.List()
		if len(agents) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No agents found. Place agent definition files in .docloom/agents/ or ~/.docloom/agents/")
			return nil
		}

		// Print in tabular format
		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tDESCRIPTION")
		fmt.Fprintln(w, "----\t-----------")
		
		for _, agent := range agents {
			description := agent.Metadata.Description
			if description == "" {
				description = "(no description)"
			}
			fmt.Fprintf(w, "%s\t%s\n", agent.Metadata.Name, description)
		}
		
		return w.Flush()
	},
}

// agentsDescribeCmd represents the agents describe command
var agentsDescribeCmd = &cobra.Command{
	Use:   "describe <agent-name>",
	Short: "Show detailed information about a specific agent",
	Long: `Display the full details of a Research Agent including its
name, description, runner command, and parameters.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		agentName := args[0]

		// Create and populate registry
		registry := agent.NewRegistry()
		if err := registry.Discover(); err != nil {
			return fmt.Errorf("failed to discover agents: %w", err)
		}

		// Find the requested agent
		agentDef, exists := registry.Get(agentName)
		if !exists {
			return fmt.Errorf("agent '%s' not found", agentName)
		}

		// Print agent details in human-readable format
		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "Agent: %s\n", agentDef.Metadata.Name)
		fmt.Fprintf(out, "API Version: %s\n", agentDef.APIVersion)
		fmt.Fprintf(out, "Kind: %s\n", agentDef.Kind)
		
		if agentDef.Metadata.Description != "" {
			fmt.Fprintf(out, "Description: %s\n", agentDef.Metadata.Description)
		}
		
		fmt.Fprintf(out, "\nRunner:\n")
		fmt.Fprintf(out, "  Command: %s\n", agentDef.Spec.Runner.Command)
		if len(agentDef.Spec.Runner.Args) > 0 {
			fmt.Fprintf(out, "  Args:\n")
			for _, arg := range agentDef.Spec.Runner.Args {
				fmt.Fprintf(out, "    - %s\n", arg)
			}
		}
		
		if len(agentDef.Spec.Parameters) > 0 {
			fmt.Fprintf(out, "\nParameters:\n")
			for _, param := range agentDef.Spec.Parameters {
				fmt.Fprintf(out, "  - Name: %s\n", param.Name)
				fmt.Fprintf(out, "    Type: %s\n", param.Type)
				if param.Description != "" {
					fmt.Fprintf(out, "    Description: %s\n", param.Description)
				}
				fmt.Fprintf(out, "    Required: %v\n", param.Required)
				if param.Default != nil {
					fmt.Fprintf(out, "    Default: %v\n", param.Default)
				}
				fmt.Fprintln(out)
			}
		} else {
			fmt.Fprintf(out, "\nNo parameters defined.\n")
		}
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(agentsCmd)
	agentsCmd.AddCommand(agentsListCmd)
	agentsCmd.AddCommand(agentsDescribeCmd)
}