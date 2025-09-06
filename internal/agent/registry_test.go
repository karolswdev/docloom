package agent

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TC-18.1: Test agent registry discovery mechanism
func TestAgentRegistry_DiscoverAgents(t *testing.T) {
	// Arrange: Create a temporary directory structure with valid .agent.yaml files and some non-agent files
	tempDir := t.TempDir()
	agentsDir := filepath.Join(tempDir, "agents")
	require.NoError(t, os.MkdirAll(agentsDir, 0755))

	// Create valid agent files
	validAgent1 := `
apiVersion: v1
kind: ResearchAgent
metadata:
  name: research-agent
  description: Conducts research on topics
spec:
  runner:
    command: python3
    args: [research.py]
  parameters:
    - name: topic
      type: string
      description: Research topic
      required: true
`

	validAgent2 := `
apiVersion: v1
kind: ResearchAgent
metadata:
  name: analysis-agent
  description: Analyzes data and patterns
spec:
  runner:
    command: node
    args: [analyze.js]
  parameters:
    - name: dataset
      type: string
      description: Dataset to analyze
      required: true
`

	// Create an invalid file (not an agent)
	nonAgentFile := `
# This is just a regular YAML file
some_data:
  key: value
`

	// Write files
	require.NoError(t, os.WriteFile(filepath.Join(agentsDir, "research.agent.yaml"), []byte(validAgent1), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(agentsDir, "analysis.agent.yml"), []byte(validAgent2), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(agentsDir, "config.yaml"), []byte(nonAgentFile), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(agentsDir, "README.md"), []byte("# Agents Directory"), 0644))

	// Act: Run the registry's discovery mechanism on the temp directory
	registry := NewRegistry()
	registry.searchPaths = []string{agentsDir} // Override default paths
	err := registry.Discover()

	// Assert: The registry must contain only the valid agents and must ignore the other files
	require.NoError(t, err)
	assert.Len(t, registry.agents, 2)

	// Verify research-agent
	researchAgent, exists := registry.Get("research-agent")
	assert.True(t, exists)
	assert.Equal(t, "research-agent", researchAgent.Metadata.Name)
	assert.Equal(t, "Conducts research on topics", researchAgent.Metadata.Description)
	assert.Equal(t, "python3", researchAgent.Spec.Runner.Command)

	// Verify analysis-agent
	analysisAgent, exists := registry.Get("analysis-agent")
	assert.True(t, exists)
	assert.Equal(t, "analysis-agent", analysisAgent.Metadata.Name)
	assert.Equal(t, "Analyzes data and patterns", analysisAgent.Metadata.Description)
	assert.Equal(t, "node", analysisAgent.Spec.Runner.Command)

	// Verify non-agent files were ignored
	_, exists = registry.Get("config")
	assert.False(t, exists)
}

func TestAgentRegistry_InvalidAgents(t *testing.T) {
	t.Run("missing apiVersion", func(t *testing.T) {
		tempDir := t.TempDir()
		invalidAgent := `
kind: ResearchAgent
metadata:
  name: invalid-agent
  description: Missing apiVersion
spec:
  runner:
    command: echo
`
		require.NoError(t, os.WriteFile(filepath.Join(tempDir, "invalid.agent.yaml"), []byte(invalidAgent), 0644))

		registry := NewRegistry()
		registry.searchPaths = []string{tempDir}
		err := registry.Discover()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing apiVersion")
	})

	t.Run("wrong kind", func(t *testing.T) {
		tempDir := t.TempDir()
		wrongKind := `
apiVersion: v1
kind: SomethingElse
metadata:
  name: wrong-kind
  description: Wrong kind field
spec:
  runner:
    command: echo
`
		require.NoError(t, os.WriteFile(filepath.Join(tempDir, "wrong.agent.yaml"), []byte(wrongKind), 0644))

		registry := NewRegistry()
		registry.searchPaths = []string{tempDir}
		err := registry.Discover()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid kind")
	})

	t.Run("missing name", func(t *testing.T) {
		tempDir := t.TempDir()
		noName := `
apiVersion: v1
kind: ResearchAgent
metadata:
  description: Missing name
spec:
  runner:
    command: echo
`
		require.NoError(t, os.WriteFile(filepath.Join(tempDir, "noname.agent.yaml"), []byte(noName), 0644))

		registry := NewRegistry()
		registry.searchPaths = []string{tempDir}
		err := registry.Discover()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing metadata.name")
	})
}

func TestAgentRegistry_List(t *testing.T) {
	tempDir := t.TempDir()

	agent1 := `
apiVersion: v1
kind: ResearchAgent
metadata:
  name: agent-one
  description: First agent
spec:
  runner:
    command: echo
`

	agent2 := `
apiVersion: v1
kind: ResearchAgent
metadata:
  name: agent-two
  description: Second agent
spec:
  runner:
    command: echo
`

	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "one.agent.yaml"), []byte(agent1), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "two.agent.yaml"), []byte(agent2), 0644))

	registry := NewRegistry()
	registry.searchPaths = []string{tempDir}
	require.NoError(t, registry.Discover())

	agents := registry.List()
	assert.Len(t, agents, 2)

	// Check that both agents are in the list
	names := []string{agents[0].Metadata.Name, agents[1].Metadata.Name}
	assert.Contains(t, names, "agent-one")
	assert.Contains(t, names, "agent-two")
}
