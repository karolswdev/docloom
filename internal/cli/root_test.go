package cli

import (
	"bytes"
	"strings"
	"testing"
)

// TC-1.1: Test that --help flag works correctly
func TestRootCmd_HelpFlag(t *testing.T) {
	// Arrange
	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--help"})
	
	// Act
	err := cmd.Execute()
	
	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	output := buf.String()
	
	// Debug: Print actual output to see what we get
	t.Logf("Help output:\n%s", output)
	
	// Check that output contains usage information
	if !strings.Contains(output, "docloom") {
		t.Error("Help output should contain 'docloom'")
	}
	
	if !strings.Contains(output, "generate") {
		t.Error("Help output should contain 'generate' command")
	}
	
	if !strings.Contains(output, "technical folks to generate high-quality documents") {
		t.Error("Help output should contain the description")
	}
	
	if !strings.Contains(output, "Available Commands:") {
		t.Error("Help output should list available commands")
	}
	
	if !strings.Contains(output, "Flags:") {
		t.Error("Help output should list flags")
	}
	
	if !strings.Contains(output, "--verbose") || !strings.Contains(output, "-v") {
		t.Error("Help output should show verbose flag")
	}
}