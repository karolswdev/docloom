// Package agent provides types and utilities for defining and managing Research Agents.
package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ArtifactCache manages temporary directories for agent execution artifacts.
type ArtifactCache struct {
	baseDir string
}

// NewArtifactCache creates a new artifact cache manager.
func NewArtifactCache() (*ArtifactCache, error) {
	// Use system temp directory as base
	baseDir := filepath.Join(os.TempDir(), "docloom-agent-cache")

	// Ensure base directory exists
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache base directory: %w", err)
	}

	return &ArtifactCache{
		baseDir: baseDir,
	}, nil
}

// CreateRunDirectory creates a unique directory for an agent run.
func (c *ArtifactCache) CreateRunDirectory(agentName string) (string, error) {
	// Create timestamp-based unique directory
	timestamp := time.Now().Format("20060102-150405")
	runID := fmt.Sprintf("%s-%s-%d", agentName, timestamp, os.Getpid())
	runDir := filepath.Join(c.baseDir, runID)

	if err := os.MkdirAll(runDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create run directory: %w", err)
	}

	return runDir, nil
}

// Clean removes old cache directories (older than 24 hours).
func (c *ArtifactCache) Clean() error {
	entries, err := os.ReadDir(c.baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Nothing to clean
		}
		return fmt.Errorf("failed to read cache directory: %w", err)
	}

	cutoff := time.Now().Add(-24 * time.Hour)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue // Skip on error
		}

		if info.ModTime().Before(cutoff) {
			dirPath := filepath.Join(c.baseDir, entry.Name())
			_ = os.RemoveAll(dirPath) // Best effort cleanup
		}
	}

	return nil
}

// GetBaseDir returns the base directory for the cache.
func (c *ArtifactCache) GetBaseDir() string {
	return c.baseDir
}
