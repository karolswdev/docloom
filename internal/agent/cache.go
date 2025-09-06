package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ArtifactCache manages temporary directories for agent artifacts
type ArtifactCache struct {
	baseDir string
}

// NewArtifactCache creates a new artifact cache
func NewArtifactCache() (*ArtifactCache, error) {
	// Use system temp directory as base
	baseDir := filepath.Join(os.TempDir(), "docloom-agents")
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache base directory: %w", err)
	}
	
	return &ArtifactCache{
		baseDir: baseDir,
	}, nil
}

// CreateRunDirectory creates a unique directory for an agent run
func (c *ArtifactCache) CreateRunDirectory(agentName string) (string, error) {
	// Create unique directory name with timestamp
	timestamp := time.Now().Format("20060102-150405")
	runID := fmt.Sprintf("%s-%s-%d", agentName, timestamp, os.Getpid())
	
	runDir := filepath.Join(c.baseDir, runID)
	if err := os.MkdirAll(runDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create run directory: %w", err)
	}
	
	return runDir, nil
}

// GetBaseDir returns the base cache directory
func (c *ArtifactCache) GetBaseDir() string {
	return c.baseDir
}

// Clean removes old cache directories (optional, for maintenance)
func (c *ArtifactCache) Clean(maxAge time.Duration) error {
	entries, err := os.ReadDir(c.baseDir)
	if err != nil {
		return fmt.Errorf("failed to read cache directory: %w", err)
	}
	
	cutoff := time.Now().Add(-maxAge)
	for _, entry := range entries {
		if entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			if info.ModTime().Before(cutoff) {
				path := filepath.Join(c.baseDir, entry.Name())
				os.RemoveAll(path)
			}
		}
	}
	
	return nil
}