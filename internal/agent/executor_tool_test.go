package agent

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// TestAgentExecutor_RunTool verifies that the executor can invoke specific tools from an agent.
// Test Case ID: TC-24.1
// Requirement: Agent-as-Toolkit
func TestAgentExecutor_RunTool(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	agentDir := tempDir // Place agent files directly in tempDir for discovery
	require.NoError(t, os.MkdirAll(agentDir, 0755))

	// Create a mock agent binary that responds to tool subcommands
	mockScript := filepath.Join(agentDir, "mock-agent.sh")
	scriptContent := `#!/bin/bash
if [ "$1" = "list_projects" ]; then
    echo "Project1"
    echo "Project2"
    echo "Project3"
    exit 0
elif [ "$1" = "get_file_content" ]; then
    echo "File content for: $2"
    exit 0
else
    echo "Unknown tool: $1" >&2
    exit 1
fi
`
	require.NoError(t, os.WriteFile(mockScript, []byte(scriptContent), 0755))

	// Create agent definition with multiple tools
	agentDef := Definition{
		APIVersion: "docloom.io/v1alpha1",
		Kind:       "Agent",
		Metadata: Metadata{
			Name:        "test-toolkit",
			Description: "Test agent with multiple tools",
		},
		Spec: Spec{
			Tools: []Tool{
				{
					Name:        "list_projects",
					Description: "Lists all projects in the repository",
					Command:     mockScript,
					Args:        []string{"list_projects"},
				},
				{
					Name:        "get_file_content",
					Description: "Gets the content of a specific file",
					Command:     mockScript,
					Args:        []string{"get_file_content", "${FILE_PATH}"},
				},
			},
		},
	}

	// Write agent.agent.yaml
	agentYAML, err := yaml.Marshal(agentDef)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(agentDir, "agent.agent.yaml"), agentYAML, 0644))

	// Create registry
	registry := NewRegistry()
	// Add the temp directory to search paths
	registry.AddSearchPath(tempDir)
	err = registry.Discover()
	require.NoError(t, err)

	// Create cache and executor
	cache, err := NewArtifactCache()
	require.NoError(t, err)
	executor := NewExecutor(registry, cache, zerolog.Nop())

	t.Run("list_projects tool", func(t *testing.T) {
		// Act
		output, err := executor.RunTool("test-toolkit", "list_projects", nil)

		// Assert
		require.NoError(t, err)
		assert.Contains(t, output, "Project1")
		assert.Contains(t, output, "Project2")
		assert.Contains(t, output, "Project3")
	})

	t.Run("get_file_content tool with parameter", func(t *testing.T) {
		// Act
		params := map[string]string{
			"FILE_PATH": "/test/file.go",
		}
		output, err := executor.RunTool("test-toolkit", "get_file_content", params)

		// Assert
		require.NoError(t, err)
		assert.Contains(t, output, "File content for: /test/file.go")
	})

	t.Run("non-existent tool", func(t *testing.T) {
		// Act
		_, err := executor.RunTool("test-toolkit", "unknown_tool", nil)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tool 'unknown_tool' not found")
	})

	t.Run("non-existent agent", func(t *testing.T) {
		// Act
		_, err := executor.RunTool("non-existent", "list_projects", nil)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "agent not found")
	})
}