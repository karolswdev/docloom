package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// TC-17.1: Test agent definition YAML parsing
func TestAgentDefinition_ParseYAML(t *testing.T) {
	t.Run("valid agent definition", func(t *testing.T) {
		// Arrange: Create a valid test.agent.yaml string
		validYAML := `
apiVersion: v1
kind: ResearchAgent
metadata:
  name: test-agent
  description: A test research agent
spec:
  runner:
    command: python3
    args:
      - research.py
  parameters:
    - name: topic
      type: string
      description: The research topic
      required: true
    - name: depth
      type: integer
      description: Research depth level
      required: false
      default: 3
`

		// Act: Unmarshal the YAML into the new Go structs
		var def Definition
		err := yaml.Unmarshal([]byte(validYAML), &def)

		// Assert: The structs must be populated with the correct metadata and spec
		require.NoError(t, err)
		assert.Equal(t, "v1", def.APIVersion)
		assert.Equal(t, "ResearchAgent", def.Kind)
		assert.Equal(t, "test-agent", def.Metadata.Name)
		assert.Equal(t, "A test research agent", def.Metadata.Description)
		assert.Equal(t, "python3", def.Spec.Runner.Command)
		assert.Equal(t, []string{"research.py"}, def.Spec.Runner.Args)
		assert.Len(t, def.Spec.Parameters, 2)

		// Validate first parameter
		assert.Equal(t, "topic", def.Spec.Parameters[0].Name)
		assert.Equal(t, "string", def.Spec.Parameters[0].Type)
		assert.Equal(t, "The research topic", def.Spec.Parameters[0].Description)
		assert.True(t, def.Spec.Parameters[0].Required)

		// Validate second parameter
		assert.Equal(t, "depth", def.Spec.Parameters[1].Name)
		assert.Equal(t, "integer", def.Spec.Parameters[1].Type)
		assert.Equal(t, "Research depth level", def.Spec.Parameters[1].Description)
		assert.False(t, def.Spec.Parameters[1].Required)
		assert.Equal(t, 3, def.Spec.Parameters[1].Default)
	})

	t.Run("malformed YAML", func(t *testing.T) {
		// Arrange: Create malformed YAML
		malformedYAML := `
apiVersion: v1
kind: ResearchAgent
metadata:
  name: [this is invalid
  description: broken yaml
`

		// Act: Attempt to unmarshal
		var def Definition
		err := yaml.Unmarshal([]byte(malformedYAML), &def)

		// Assert: Test failure on malformed YAML
		assert.Error(t, err)
	})

	t.Run("missing required fields", func(t *testing.T) {
		// Arrange: YAML missing required fields
		incompleteYAML := `
kind: ResearchAgent
metadata:
  description: Missing name and apiVersion
`

		// Act: Unmarshal the YAML
		var def Definition
		err := yaml.Unmarshal([]byte(incompleteYAML), &def)

		// Assert: Parsing succeeds but fields are empty
		require.NoError(t, err)
		assert.Empty(t, def.APIVersion)
		assert.Empty(t, def.Metadata.Name)
	})
}
