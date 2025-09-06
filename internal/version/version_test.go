package version

import (
	"strings"
	"testing"
)

func TestInfo(t *testing.T) {
	// Save original values
	origVersion := Version
	origCommit := GitCommit
	origDate := BuildDate

	// Test with custom values
	Version = "1.0.0"
	GitCommit = "abc123"
	BuildDate = "2024-01-01T00:00:00Z"

	defer func() {
		// Restore original values
		Version = origVersion
		GitCommit = origCommit
		BuildDate = origDate
	}()

	info := Info()

	// Check that all expected fields are present
	if !strings.Contains(info, "DocLoom version 1.0.0") {
		t.Errorf("Info() missing version, got: %s", info)
	}

	if !strings.Contains(info, "Build Date: 2024-01-01T00:00:00Z") {
		t.Errorf("Info() missing build date, got: %s", info)
	}

	if !strings.Contains(info, "Git Commit: abc123") {
		t.Errorf("Info() missing git commit, got: %s", info)
	}

	if !strings.Contains(info, "Go Version:") {
		t.Errorf("Info() missing Go version, got: %s", info)
	}

	if !strings.Contains(info, "Platform:") {
		t.Errorf("Info() missing platform, got: %s", info)
	}
}

func TestShort(t *testing.T) {
	// Save original value
	origVersion := Version

	// Test with custom value
	Version = "2.0.0"

	defer func() {
		// Restore original value
		Version = origVersion
	}()

	if got := Short(); got != "2.0.0" {
		t.Errorf("Short() = %v, want %v", got, "2.0.0")
	}
}

func TestDefaultValues(t *testing.T) {
	// Check that default values are set
	if Version == "" {
		t.Error("Version should not be empty")
	}

	if GitCommit == "" {
		t.Error("GitCommit should not be empty")
	}

	if BuildDate == "" {
		t.Error("BuildDate should not be empty")
	}

	if GoVersion == "" {
		t.Error("GoVersion should not be empty")
	}

	if Platform == "" {
		t.Error("Platform should not be empty")
	}
}
