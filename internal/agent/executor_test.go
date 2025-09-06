package agent

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestAgentExecutor_RunCommand tests the basic agent execution flow
func TestAgentExecutor_RunCommand(t *testing.T) {
	// Skip in CI environments where bash might not be available
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping test in CI environment")
	}

	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Create mock agent script
	mockAgentPath := filepath.Join(tmpDir, "mock-agent.sh")
	mockAgentContent := `#!/bin/bash
# Mock agent that copies a file and logs parameters
SOURCE_PATH=$1
OUTPUT_PATH=$2

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_PATH"

# Create a test file in the output
echo "Test artifact content" > "$OUTPUT_PATH/artifact.md"

# Log the parameters to verify they were passed correctly
echo "Source: $SOURCE_PATH" > "$OUTPUT_PATH/agent.log"
echo "Output: $OUTPUT_PATH" >> "$OUTPUT_PATH/agent.log"
echo "Environment variables:" >> "$OUTPUT_PATH/agent.log"
env | grep PARAM_ >> "$OUTPUT_PATH/agent.log" || true

exit 0
`
	if err := os.WriteFile(mockAgentPath, []byte(mockAgentContent), 0755); err != nil {
		t.Fatalf("Failed to create mock agent: %v", err)
	}

	// Create test source directory
	sourceDir := filepath.Join(tmpDir, "source")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	// Create a logger for testing
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Create executor
	executor, err := NewExecutor(logger)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	// Define test agent
	def := &SimpleDefinition{
		Name:    "mock-agent",
		Command: mockAgentPath,
		Parameters: map[string]string{
			"test_param": "test_value",
		},
	}

	// Execute agent
	outputPath, err := executor.Execute(def, ExecuteOptions{
		SourcePath: sourceDir,
		Parameters: map[string]string{
			"custom_param": "custom_value",
		},
	})
	if err != nil {
		t.Fatalf("Failed to execute agent: %v", err)
	}

	// Verify output artifact was created
	artifactPath := filepath.Join(outputPath, "artifact.md")
	if _, statErr := os.Stat(artifactPath); os.IsNotExist(statErr) {
		t.Errorf("Expected artifact file was not created at %s", artifactPath)
	}

	// Verify log file contains correct information
	logPath := filepath.Join(outputPath, "agent.log")
	logContent, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read agent log: %v", err)
	}

	logStr := string(logContent)
	if !strings.Contains(logStr, sourceDir) {
		t.Errorf("Log does not contain source path: %s", logStr)
	}
	if !strings.Contains(logStr, outputPath) {
		t.Errorf("Log does not contain output path: %s", logStr)
	}

	// Verify artifact content
	artifactContent, err := os.ReadFile(artifactPath)
	if err != nil {
		t.Fatalf("Failed to read artifact: %v", err)
	}
	if !strings.Contains(string(artifactContent), "Test artifact content") {
		t.Errorf("Artifact has unexpected content: %s", artifactContent)
	}

	t.Logf("Test passed - Agent executed successfully")
	t.Logf("Output path: %s", outputPath)
	t.Logf("Log content:\n%s", logStr)
}

// TestAgentExecutor_ParameterOverrides tests parameter override functionality
func TestAgentExecutor_ParameterOverrides(t *testing.T) {
	// Skip in CI environments where bash might not be available
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping test in CI environment")
	}

	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Create mock agent that echoes environment variables
	mockAgentPath := filepath.Join(tmpDir, "mock-agent.sh")
	mockAgentContent := `#!/bin/bash
SOURCE_PATH=$1
OUTPUT_PATH=$2

mkdir -p "$OUTPUT_PATH"

# Echo all PARAM_ environment variables to a file
echo "Parameters received:" > "$OUTPUT_PATH/params.log"
env | grep PARAM_ | sort >> "$OUTPUT_PATH/params.log"

exit 0
`
	if err := os.WriteFile(mockAgentPath, []byte(mockAgentContent), 0755); err != nil {
		t.Fatalf("Failed to create mock agent: %v", err)
	}

	// Create test source directory
	sourceDir := filepath.Join(tmpDir, "source")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	// Create a logger for testing
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Create executor
	executor, err := NewExecutor(logger)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	// Define test agent with default parameter
	def := &SimpleDefinition{
		Name:    "mock-agent",
		Command: mockAgentPath,
		Parameters: map[string]string{
			"my_param":      "default_value",
			"another_param": "another_default",
		},
	}

	// Execute agent with parameter override
	outputPath, err := executor.Execute(def, ExecuteOptions{
		SourcePath: sourceDir,
		Parameters: map[string]string{
			"my_param": "overridden",
		},
	})
	if err != nil {
		t.Fatalf("Failed to execute agent: %v", err)
	}

	// Read parameter log
	paramsPath := filepath.Join(outputPath, "params.log")
	paramsContent, err := os.ReadFile(paramsPath)
	if err != nil {
		t.Fatalf("Failed to read params log: %v", err)
	}

	paramsStr := string(paramsContent)
	t.Logf("Parameters log:\n%s", paramsStr)

	// Verify the overridden parameter
	if !strings.Contains(paramsStr, "PARAM_MY_PARAM=overridden") {
		t.Errorf("Parameter was not overridden correctly. Expected 'overridden', log:\n%s", paramsStr)
	}

	// Verify the non-overridden parameter keeps default
	if !strings.Contains(paramsStr, "PARAM_ANOTHER_PARAM=another_default") {
		t.Errorf("Default parameter was not preserved. Expected 'another_default', log:\n%s", paramsStr)
	}

	t.Logf("Test passed - Parameter overrides work correctly")
}
