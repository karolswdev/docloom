// Package agent provides types and utilities for defining and managing Research Agents.
package agent

// Definition represents a complete agent definition from a .agent.yaml file.
type Definition struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

// Metadata contains the agent's basic information.
type Metadata struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// Spec defines the agent's execution specification.
type Spec struct {
	Runner     Runner      `yaml:"runner,omitempty"` // Deprecated: Use Tools instead
	Tools      []Tool      `yaml:"tools,omitempty"`
	Parameters []Parameter `yaml:"parameters"`
}

// Runner specifies how to execute the agent.
type Runner struct {
	Command string   `yaml:"command"`
	Args    []string `yaml:"args,omitempty"`
}

// Tool represents a specific capability that an agent exposes.
// Each tool can be invoked independently by the LLM during analysis.
type Tool struct {
	Name        string   `yaml:"name"`           // Tool identifier (e.g., "list_projects")
	Description string   `yaml:"description"`    // LLM-facing description of what the tool does
	Command     string   `yaml:"command"`        // Command to execute (can include the agent binary path)
	Args        []string `yaml:"args,omitempty"` // Additional arguments to pass
}

// Parameter defines an input parameter for the agent.
type Parameter struct {
	Name        string      `yaml:"name"`
	Type        string      `yaml:"type"`
	Required    bool        `yaml:"required"`
	Default     interface{} `yaml:"default,omitempty"`
	Description string      `yaml:"description"`
}
