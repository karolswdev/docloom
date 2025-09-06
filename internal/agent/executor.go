package agent

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
)

// Executor handles the execution of agents as external processes.
type Executor struct {
	registry *Registry
	cache    *ArtifactCache
	logger   zerolog.Logger
}

// NewExecutor creates a new agent executor.
func NewExecutor(registry *Registry, cache *ArtifactCache, logger zerolog.Logger) *Executor {
	return &Executor{
		registry: registry,
		cache:    cache,
		logger:   logger,
	}
}

// RunOptions contains options for running an agent.
type RunOptions struct {
	AgentName  string            // Name of the agent to run
	SourcePath string            // Path to source files for analysis
	Parameters map[string]string // Parameter overrides
}

// RunResult contains the result of an agent execution.
type RunResult struct {
	OutputPath string // Path to the output directory containing artifacts
	ExitCode   int    // Exit code from the agent process
}

// Run executes an agent with the given options.
func (e *Executor) Run(opts RunOptions) (*RunResult, error) {
	// Look up agent in registry
	agent, exists := e.registry.Get(opts.AgentName)
	if !exists {
		return nil, fmt.Errorf("agent not found: %s", opts.AgentName)
	}

	e.logger.Info().
		Str("agent", opts.AgentName).
		Str("source", opts.SourcePath).
		Msg("Executing agent")

	// Create unique output directory for this run
	outputPath, err := e.cache.CreateRunDirectory(opts.AgentName)
	if err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	// Prepare command with arguments
	args := make([]string, 0, len(agent.Spec.Runner.Args)+2)
	
	// Add configured args from agent definition
	for _, arg := range agent.Spec.Runner.Args {
		// Replace placeholders
		arg = strings.ReplaceAll(arg, "${SOURCE_PATH}", opts.SourcePath)
		arg = strings.ReplaceAll(arg, "${OUTPUT_PATH}", outputPath)
		args = append(args, arg)
	}
	
	// If no args specified, use default pattern (source output)
	if len(args) == 0 {
		args = []string{opts.SourcePath, outputPath}
	}

	cmd := exec.Command(agent.Spec.Runner.Command, args...)

	// Set up environment variables for parameters
	env := os.Environ()
	
	// First, apply default values from agent definition
	for _, param := range agent.Spec.Parameters {
		if param.Default != nil {
			envKey := fmt.Sprintf("PARAM_%s", strings.ToUpper(param.Name))
			envVal := fmt.Sprintf("%v", param.Default)
			env = append(env, fmt.Sprintf("%s=%s", envKey, envVal))
		}
	}
	
	// Then apply overrides from options
	for key, value := range opts.Parameters {
		envKey := fmt.Sprintf("PARAM_%s", strings.ToUpper(key))
		env = append(env, fmt.Sprintf("%s=%s", envKey, value))
	}
	
	cmd.Env = env

	// Set up stdout and stderr pipes for logging
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start agent: %w", err)
	}

	// Stream output to logger
	go e.streamOutput(stdout, "stdout", opts.AgentName)
	go e.streamOutput(stderr, "stderr", opts.AgentName)

	// Wait for completion
	err = cmd.Wait()
	
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
			e.logger.Warn().
				Str("agent", opts.AgentName).
				Int("exit_code", exitCode).
				Msg("Agent exited with non-zero code")
		} else {
			return nil, fmt.Errorf("failed to wait for agent: %w", err)
		}
	}

	e.logger.Info().
		Str("agent", opts.AgentName).
		Str("output", outputPath).
		Int("exit_code", exitCode).
		Msg("Agent execution completed")

	return &RunResult{
		OutputPath: outputPath,
		ExitCode:   exitCode,
	}, nil
}

// streamOutput streams output from a reader to the logger.
func (e *Executor) streamOutput(reader io.Reader, stream string, agentName string) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		e.logger.Debug().
			Str("agent", agentName).
			Str("stream", stream).
			Str("line", scanner.Text()).
			Msg("Agent output")
	}
}

// ValidateOutput checks if the agent produced expected output files.
func (e *Executor) ValidateOutput(outputPath string) error {
	// Check if directory exists and has files
	entries, err := os.ReadDir(outputPath)
	if err != nil {
		return fmt.Errorf("failed to read output directory: %w", err)
	}

	if len(entries) == 0 {
		return fmt.Errorf("agent produced no output files")
	}

	// Check for at least one markdown file
	hasMarkdown := false
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".md" {
			hasMarkdown = true
			break
		}
	}

	if !hasMarkdown {
		return fmt.Errorf("agent did not produce any markdown files")
	}

	return nil
}