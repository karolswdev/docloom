package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Registry manages discovered agent definitions.
type Registry struct {
	agents      map[string]*Definition
	searchPaths []string
}

// NewRegistry creates a new agent registry with default search paths.
func NewRegistry() *Registry {
	homeDir, err := os.UserHomeDir()
	searchPaths := []string{
		".docloom/agents", // Project-local agents
	}
	if err == nil && homeDir != "" {
		searchPaths = append(searchPaths, filepath.Join(homeDir, ".docloom", "agents")) // User-home agents
	}

	return &Registry{
		agents:      make(map[string]*Definition),
		searchPaths: searchPaths,
	}
}

// AddSearchPath adds a custom search path for agent discovery.
func (r *Registry) AddSearchPath(path string) {
	r.searchPaths = append(r.searchPaths, path)
}

// Discover searches for and loads agent definition files.
func (r *Registry) Discover() error {
	for _, searchPath := range r.searchPaths {
		if err := r.discoverInPath(searchPath); err != nil {
			// Log but don't fail if a path doesn't exist
			if !os.IsNotExist(err) {
				return fmt.Errorf("error discovering agents in %s: %w", searchPath, err)
			}
		}
	}
	return nil
}

// discoverInPath searches for agent files in a specific directory.
func (r *Registry) discoverInPath(searchPath string) error {
	entries, err := os.ReadDir(searchPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if strings.HasSuffix(entry.Name(), ".agent.yaml") || strings.HasSuffix(entry.Name(), ".agent.yml") {
			fullPath := filepath.Join(searchPath, entry.Name())
			if err := r.loadAgent(fullPath); err != nil {
				return fmt.Errorf("error loading agent %s: %w", fullPath, err)
			}
		}
	}

	return nil
}

// loadAgent loads a single agent definition from a file.
func (r *Registry) loadAgent(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var def Definition
	if err := yaml.Unmarshal(data, &def); err != nil {
		return fmt.Errorf("error parsing YAML: %w", err)
	}

	// Validate the definition
	if def.APIVersion == "" {
		return fmt.Errorf("missing apiVersion")
	}
	if def.Kind != "ResearchAgent" {
		return fmt.Errorf("invalid kind: %s (expected ResearchAgent)", def.Kind)
	}
	if def.Metadata.Name == "" {
		return fmt.Errorf("missing metadata.name")
	}

	r.agents[def.Metadata.Name] = &def
	return nil
}

// Get retrieves an agent definition by name.
func (r *Registry) Get(name string) (*Definition, bool) {
	agent, exists := r.agents[name]
	return agent, exists
}

// List returns all discovered agents.
func (r *Registry) List() []*Definition {
	agents := make([]*Definition, 0, len(r.agents))
	for _, agent := range r.agents {
		agents = append(agents, agent)
	}
	return agents
}
