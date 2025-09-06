package agent

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// SimpleDefinition represents a simple agent configuration (legacy)
type SimpleDefinition struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Command     string            `yaml:"command"`
	Parameters  map[string]string `yaml:"parameters,omitempty"`
}

// Executor runs agent commands
type Executor struct {
	cache  *ArtifactCache
	logger *slog.Logger
}

// NewExecutor creates a new agent executor
func NewExecutor(logger *slog.Logger) (*Executor, error) {
	cache, err := NewArtifactCache()
	if err != nil {
		return nil, fmt.Errorf("failed to create artifact cache: %w", err)
	}
	
	return &Executor{
		cache:  cache,
		logger: logger,
	}, nil
}

// ExecuteOptions contains options for executing an agent
type ExecuteOptions struct {
	SourcePath string            // Path to source code repository
	OutputPath string            // Path where agent should write output
	Parameters map[string]string // Parameter overrides
}

// Execute runs an agent with the given definition and options
func (e *Executor) Execute(def *SimpleDefinition, opts ExecuteOptions) (string, error) {
	e.logger.Info("executing agent", "name", def.Name, "command", def.Command)
	
	// Create run directory if output path not specified
	outputPath := opts.OutputPath
	if outputPath == "" {
		var err error
		outputPath, err = e.cache.CreateRunDirectory(def.Name)
		if err != nil {
			return "", fmt.Errorf("failed to create run directory: %w", err)
		}
		e.logger.Debug("created artifact directory", "path", outputPath)
	}
	
	// Prepare command with arguments
	// Agent contract: command receives source and output paths as arguments
	cmdParts := strings.Fields(def.Command)
	if len(cmdParts) == 0 {
		return "", fmt.Errorf("empty command")
	}
	
	cmd := exec.Command(cmdParts[0], append(cmdParts[1:], opts.SourcePath, outputPath)...)
	
	// Set up environment variables for parameters
	env := os.Environ()
	
	// First add default parameters from definition
	for key, value := range def.Parameters {
		env = append(env, fmt.Sprintf("PARAM_%s=%s", strings.ToUpper(key), value))
	}
	
	// Then override with user-provided parameters
	for key, value := range opts.Parameters {
		// Find and replace existing or append new
		envKey := fmt.Sprintf("PARAM_%s=", strings.ToUpper(key))
		found := false
		for i, envVar := range env {
			if strings.HasPrefix(envVar, envKey) {
				env[i] = envKey + value
				found = true
				break
			}
		}
		if !found {
			env = append(env, envKey+value)
		}
	}
	
	cmd.Env = env
	
	// Set up stdout and stderr pipes for streaming
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stderr pipe: %w", err)
	}
	
	// Start the command
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start agent command: %w", err)
	}
	
	// Stream output to logger
	go e.streamOutput(stdout, "stdout", def.Name)
	go e.streamOutput(stderr, "stderr", def.Name)
	
	// Wait for completion
	if err := cmd.Wait(); err != nil {
		return "", fmt.Errorf("agent command failed: %w", err)
	}
	
	e.logger.Info("agent execution completed", "name", def.Name, "output", outputPath)
	
	return outputPath, nil
}

// streamOutput streams command output to the logger
func (e *Executor) streamOutput(reader io.Reader, stream string, agentName string) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		e.logger.Debug("agent output",
			"agent", agentName,
			"stream", stream,
			"line", scanner.Text())
	}
}

// LoadDefinition loads an agent definition from a YAML file
func LoadDefinition(path string) (*SimpleDefinition, error) {
	// For now, return a simple mock for testing
	// Real implementation would parse YAML
	return &SimpleDefinition{
		Name:    filepath.Base(path),
		Command: path,
		Parameters: map[string]string{
			"default_param": "default_value",
		},
	}, nil
}