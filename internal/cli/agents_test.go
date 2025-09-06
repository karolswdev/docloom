package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TC-19.1: Test agents list command E2E
func TestAgentsListCmd_E2E(t *testing.T) {
	// Arrange: Point the agent registry to a test directory with two known agent files
	testDir := t.TempDir()
	agentsDir := filepath.Join(testDir, ".docloom", "agents")
	require.NoError(t, os.MkdirAll(agentsDir, 0755))

	// Create two test agent files
	agent1 := `
apiVersion: v1
kind: ResearchAgent
metadata:
  name: test-agent-one
  description: First test agent for listing
spec:
  runner:
    command: echo
    args: ["test1"]
  parameters:
    - name: param1
      type: string
      description: Test parameter
      required: true
`

	agent2 := `
apiVersion: v1
kind: ResearchAgent
metadata:
  name: test-agent-two
  description: Second test agent for listing
spec:
  runner:
    command: echo
    args: ["test2"]
`

	require.NoError(t, os.WriteFile(filepath.Join(agentsDir, "agent1.agent.yaml"), []byte(agent1), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(agentsDir, "agent2.agent.yaml"), []byte(agent2), 0644))

	// Change to test directory
	originalWd, _ := os.Getwd()
	require.NoError(t, os.Chdir(testDir))
	defer os.Chdir(originalWd)

	// Act: Execute the `docloom agents list` command
	rootCmd.SetArgs([]string{"agents", "list"})

	// Capture output
	var stdout bytes.Buffer
	rootCmd.SetOut(&stdout)
	rootCmd.SetErr(&stdout)

	err := rootCmd.Execute()

	// Assert: The command's standard output must contain the names and descriptions of the two test agents in a clean, tabular format
	require.NoError(t, err)

	output := stdout.String()
	t.Logf("agents list output:\n%s", output)

	// Check for header
	assert.Contains(t, output, "NAME")
	assert.Contains(t, output, "DESCRIPTION")

	// Check for both agents
	assert.Contains(t, output, "test-agent-one")
	assert.Contains(t, output, "First test agent for listing")
	assert.Contains(t, output, "test-agent-two")
	assert.Contains(t, output, "Second test agent for listing")

	// Verify tabular format (check for consistent spacing/alignment)
	lines := strings.Split(strings.TrimSpace(output), "\n")
	assert.GreaterOrEqual(t, len(lines), 4) // Header, separator, and at least 2 agents
}

// TC-19.2: Test agents describe command E2E
func TestAgentsDescribeCmd_E2E(t *testing.T) {
	// Arrange: Point the agent registry to a test directory with a known agent file
	testDir := t.TempDir()
	agentsDir := filepath.Join(testDir, ".docloom", "agents")
	require.NoError(t, os.MkdirAll(agentsDir, 0755))

	// Create a test agent file with full details
	agentDef := `
apiVersion: v1
kind: ResearchAgent
metadata:
  name: detailed-test-agent
  description: A comprehensive test agent for the describe command
spec:
  runner:
    command: python3
    args:
      - /path/to/script.py
      - --verbose
  parameters:
    - name: input_file
      type: string
      description: Path to the input file
      required: true
    - name: depth_level
      type: integer
      description: How deep to analyze
      required: false
      default: 5
    - name: enable_cache
      type: boolean
      description: Whether to use caching
      required: false
      default: true
`

	require.NoError(t, os.WriteFile(filepath.Join(agentsDir, "detailed.agent.yaml"), []byte(agentDef), 0644))

	// Change to test directory
	originalWd, _ := os.Getwd()
	require.NoError(t, os.Chdir(testDir))
	defer os.Chdir(originalWd)

	// Act: Execute `docloom agents describe <agent-name>`
	rootCmd.SetArgs([]string{"agents", "describe", "detailed-test-agent"})

	// Capture output
	var stdout bytes.Buffer
	rootCmd.SetOut(&stdout)
	rootCmd.SetErr(&stdout)

	err := rootCmd.Execute()

	// Assert: The standard output must contain the agent's full details
	require.NoError(t, err)

	output := stdout.String()
	t.Logf("agents describe output:\n%s", output)

	// Check agent metadata
	assert.Contains(t, output, "Agent: detailed-test-agent")
	assert.Contains(t, output, "API Version: v1")
	assert.Contains(t, output, "Kind: ResearchAgent")
	assert.Contains(t, output, "Description: A comprehensive test agent for the describe command")

	// Check runner details
	assert.Contains(t, output, "Runner:")
	assert.Contains(t, output, "Command: python3")
	assert.Contains(t, output, "/path/to/script.py")
	assert.Contains(t, output, "--verbose")

	// Check parameters
	assert.Contains(t, output, "Parameters:")
	assert.Contains(t, output, "Name: input_file")
	assert.Contains(t, output, "Type: string")
	assert.Contains(t, output, "Description: Path to the input file")
	assert.Contains(t, output, "Required: true")

	assert.Contains(t, output, "Name: depth_level")
	assert.Contains(t, output, "Type: integer")
	assert.Contains(t, output, "Description: How deep to analyze")
	assert.Contains(t, output, "Required: false")
	assert.Contains(t, output, "Default: 5")

	assert.Contains(t, output, "Name: enable_cache")
	assert.Contains(t, output, "Type: boolean")
	assert.Contains(t, output, "Description: Whether to use caching")
	assert.Contains(t, output, "Required: false")
	assert.Contains(t, output, "Default: true")
}

func TestAgentsDescribeCmd_NotFound(t *testing.T) {
	// Test error handling when agent doesn't exist
	testDir := t.TempDir()
	agentsDir := filepath.Join(testDir, ".docloom", "agents")
	require.NoError(t, os.MkdirAll(agentsDir, 0755))

	// Change to test directory
	originalWd, _ := os.Getwd()
	require.NoError(t, os.Chdir(testDir))
	defer os.Chdir(originalWd)

	// Try to describe non-existent agent
	rootCmd.SetArgs([]string{"agents", "describe", "non-existent-agent"})

	var stderr bytes.Buffer
	rootCmd.SetOut(&stderr)
	rootCmd.SetErr(&stderr)

	err := rootCmd.Execute()

	// Should return an error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "agent 'non-existent-agent' not found")
}

func TestAgentsListCmd_Empty(t *testing.T) {
	// Test list command when no agents are present
	testDir := t.TempDir()
	agentsDir := filepath.Join(testDir, ".docloom", "agents")
	require.NoError(t, os.MkdirAll(agentsDir, 0755))

	// Change to test directory
	originalWd, _ := os.Getwd()
	require.NoError(t, os.Chdir(testDir))
	defer os.Chdir(originalWd)

	// Execute list command
	rootCmd.SetArgs([]string{"agents", "list"})

	var stdout bytes.Buffer
	rootCmd.SetOut(&stdout)
	rootCmd.SetErr(&stdout)

	err := rootCmd.Execute()

	// Should succeed but show message about no agents
	require.NoError(t, err)
	output := stdout.String()
	assert.Contains(t, output, "No agents found")
	assert.Contains(t, output, ".docloom/agents/")
}
