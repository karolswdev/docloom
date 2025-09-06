package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestVersionCmd tests the version command (TC-14.1)
func TestVersionCmd(t *testing.T) {
	// Build the binary with version information
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "docloom")

	// Build command with ldflags (disable VCS stamping in container environments)
	buildCmd := exec.Command("go", "build",
		"-buildvcs=false",
		"-ldflags",
		"-X github.com/karolswdev/docloom/internal/version.Version=1.0.0-test "+
			"-X github.com/karolswdev/docloom/internal/version.GitCommit=test123 "+
			"-X github.com/karolswdev/docloom/internal/version.BuildDate=2024-01-01T00:00:00Z",
		"-o", binaryPath,
		"../../cmd/docloom",
	)

	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build binary: %v\nOutput: %s", err, buildOutput)
	}

	// Test 1: Run with 'version' subcommand
	versionCmd := exec.Command(binaryPath, "version")
	output, err := versionCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run version command: %v\nOutput: %s", err, output)
	}

	outputStr := string(output)

	// Verify all required information is present
	if !strings.Contains(outputStr, "DocLoom version") {
		t.Errorf("Version output missing 'DocLoom version', got: %s", outputStr)
	}

	if !strings.Contains(outputStr, "1.0.0-test") {
		t.Errorf("Version output missing version '1.0.0-test', got: %s", outputStr)
	}

	if !strings.Contains(outputStr, "test123") {
		t.Errorf("Version output missing commit 'test123', got: %s", outputStr)
	}

	if !strings.Contains(outputStr, "2024-01-01T00:00:00Z") {
		t.Errorf("Version output missing build date '2024-01-01T00:00:00Z', got: %s", outputStr)
	}

	if !strings.Contains(outputStr, "Go Version:") {
		t.Errorf("Version output missing 'Go Version:', got: %s", outputStr)
	}

	if !strings.Contains(outputStr, "Platform:") {
		t.Errorf("Version output missing 'Platform:', got: %s", outputStr)
	}

	// Test 2: Run with '--version' flag
	versionFlagCmd := exec.Command(binaryPath, "--version")
	flagOutput, err := versionFlagCmd.CombinedOutput()
	if err == nil || err.Error() == "exit status 1" {
		// The command returns an error to stop execution, but that's expected
		flagOutputStr := string(flagOutput)

		if !strings.Contains(flagOutputStr, "DocLoom version") {
			t.Errorf("--version flag output missing 'DocLoom version', got: %s", flagOutputStr)
		}

		if !strings.Contains(flagOutputStr, "1.0.0-test") {
			t.Errorf("--version flag output missing version '1.0.0-test', got: %s", flagOutputStr)
		}
	} else {
		t.Fatalf("Unexpected error running --version flag: %v\nOutput: %s", err, flagOutput)
	}

	// Save evidence
	evidenceDir := "../../evidence/PHASE-4/story-4.1/task-1"
	if err := os.MkdirAll(evidenceDir, 0755); err != nil {
		t.Logf("Warning: Could not create evidence directory: %v", err)
	} else {
		evidenceFile := filepath.Join(evidenceDir, "TC-14.1-version-output.txt")
		if err := os.WriteFile(evidenceFile, output, 0644); err != nil {
			t.Logf("Warning: Could not save evidence: %v", err)
		} else {
			t.Logf("Evidence saved to: %s", evidenceFile)
		}
	}
}
