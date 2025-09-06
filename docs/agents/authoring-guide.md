# Agent Authoring Guide

## Overview

This guide explains how to create powerful, AI-compatible agents for DocLoom using the tool-based architecture. Agents are external programs that provide specialized analysis capabilities, exposed as discrete tools that can be orchestrated by AI.

## Tool-Based Paradigm

Modern DocLoom agents expose multiple **tools** rather than a single monolithic function. This enables:

- **Granular Control**: AI can invoke specific capabilities as needed
- **Composability**: Tools can be combined in different ways
- **Efficiency**: Only necessary tools are executed
- **Discoverability**: AI understands each tool's purpose from its description

## Agent Structure

An agent consists of:

1. **Executable Binary**: A program that implements the tools
2. **Agent Definition**: A YAML file describing the tools
3. **Documentation**: User and developer guides

## Creating an Agent Binary

Your agent binary should:

1. **Use Subcommands**: Each tool is a subcommand
2. **Output JSON**: Tools should output structured JSON to stdout
3. **Log to stderr**: Use stderr for debugging and progress
4. **Exit Cleanly**: Return 0 on success, non-zero on error

### Example Structure (Go with Cobra)

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "my-agent",
    Short: "My specialized agent for DocLoom",
}

var listCmd = &cobra.Command{
    Use:   "list_items [path]",
    Short: "Lists items in the repository",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        path := args[0]
        
        // Your analysis logic here
        items := analyzeRepository(path)
        
        // Output JSON to stdout
        output := map[string]interface{}{
            "itemCount": len(items),
            "items":     items,
        }
        json.NewEncoder(os.Stdout).Encode(output)
    },
}

var analyzeCmd = &cobra.Command{
    Use:   "analyze_item [path]",
    Short: "Analyzes a specific item",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        path := args[0]
        
        // Detailed analysis
        analysis := performAnalysis(path)
        
        // Output JSON
        json.NewEncoder(os.Stdout).Encode(analysis)
    },
}

func init() {
    rootCmd.AddCommand(listCmd)
    rootCmd.AddCommand(analyzeCmd)
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

### Example Structure (Python with argparse)

```python
#!/usr/bin/env python3
import json
import sys
import argparse

def list_items(path):
    """Lists items in the repository."""
    # Your analysis logic
    items = analyze_repository(path)
    
    # Output JSON to stdout
    output = {
        "itemCount": len(items),
        "items": items
    }
    print(json.dumps(output))

def analyze_item(path):
    """Analyzes a specific item."""
    # Detailed analysis
    analysis = perform_analysis(path)
    
    # Output JSON
    print(json.dumps(analysis))

def main():
    parser = argparse.ArgumentParser(description='My specialized agent')
    subparsers = parser.add_subparsers(dest='command', help='Commands')
    
    # List command
    list_parser = subparsers.add_parser('list_items', help='List items')
    list_parser.add_argument('path', help='Repository path')
    
    # Analyze command
    analyze_parser = subparsers.add_parser('analyze_item', help='Analyze item')
    analyze_parser.add_argument('path', help='Item path')
    
    args = parser.parse_args()
    
    if args.command == 'list_items':
        list_items(args.path)
    elif args.command == 'analyze_item':
        analyze_item(args.path)
    else:
        parser.print_help()
        sys.exit(1)

if __name__ == '__main__':
    main()
```

## Writing the Agent Definition

Create `agent.agent.yaml` in your agent's directory:

```yaml
apiVersion: docloom.io/v1alpha1
kind: Agent
metadata:
  name: my-analyzer
  description: Analyzes repositories for specific patterns
spec:
  tools:
    - name: list_items
      description: |
        Lists all items in the repository that match our criteria.
        Returns a JSON array of item paths with basic metadata.
        Use this first to understand what's available for analysis.
      command: ./my-agent
      args:
        - list_items
        - "${SOURCE_PATH}"
    
    - name: analyze_item
      description: |
        Performs deep analysis on a specific item.
        Returns detailed metrics, patterns, and recommendations.
        Call list_items first to identify items to analyze.
      command: ./my-agent
      args:
        - analyze_item
        - "${ITEM_PATH}"
    
    - name: compare_items
      description: |
        Compares two items and identifies differences.
        Returns a structured diff with semantic analysis.
        Useful for understanding changes or variations.
      command: ./my-agent
      args:
        - compare_items
        - "${ITEM_A}"
        - "${ITEM_B}"
  
  parameters:
    - name: source_path
      type: string
      required: true
      description: Path to the repository to analyze
    
    - name: item_path
      type: string
      required: false
      description: Path to a specific item
    
    - name: depth
      type: int
      default: 3
      description: Analysis depth level
```

## Tool Description Best Practices

Write descriptions that help the AI understand:

### 1. Purpose and Function

```yaml
# Good
description: |
  Extracts all public API endpoints from the codebase.
  Identifies REST routes, GraphQL schemas, and gRPC services.

# Poor
description: Gets APIs
```

### 2. Output Format

```yaml
# Good
description: |
  Returns a JSON object with:
  - endpoints: Array of {method, path, handler}
  - total: Number of endpoints found
  - categories: Grouped by service type

# Poor
description: Returns endpoint data
```

### 3. Prerequisites and Dependencies

```yaml
# Good
description: |
  Analyzes test coverage for a specific module.
  Requires: call list_modules first to get module paths.
  Note: This runs the test suite, which may take time.

# Poor
description: Gets test coverage
```

### 4. Use Cases and Context

```yaml
# Good
description: |
  Identifies security vulnerabilities in dependencies.
  Use when: Generating security audit documents.
  Checks: CVE database, outdated packages, known issues.

# Poor
description: Checks security
```

## Parameter Handling

### Environment Variables

Parameters are passed as environment variables with `PARAM_` prefix:

```go
func getConfig() Config {
    return Config{
        Depth:     getEnvInt("PARAM_DEPTH", 3),
        Verbose:   getEnvBool("PARAM_VERBOSE", false),
        OutputDir: getEnvString("PARAM_OUTPUT_DIR", "./output"),
    }
}
```

### Parameter Substitution in Arguments

Use `${PARAMETER_NAME}` in tool arguments:

```yaml
args:
  - analyze
  - "${SOURCE_PATH}"
  - "--depth=${DEPTH}"
  - "--output=${OUTPUT_DIR}"
```

## Output Guidelines

### JSON Structure

Keep output structured and consistent:

```json
{
  "summary": {
    "totalItems": 42,
    "categories": ["typeA", "typeB"],
    "timestamp": "2024-01-15T10:30:00Z"
  },
  "items": [
    {
      "id": "item-1",
      "type": "typeA",
      "metrics": {
        "complexity": 5,
        "size": 1024
      }
    }
  ],
  "metadata": {
    "version": "1.0.0",
    "analysisTime": 1.234
  }
}
```

### Error Handling

Return errors as JSON with proper exit codes:

```go
if err != nil {
    errorOutput := map[string]interface{}{
        "error": err.Error(),
        "code": "INVALID_PATH",
        "details": "The specified path does not exist",
    }
    json.NewEncoder(os.Stdout).Encode(errorOutput)
    os.Exit(1)
}
```

## Testing Your Agent

### Unit Testing Tools

Test each tool independently:

```bash
# Test list tool
./my-agent list_items /path/to/test/repo | jq .

# Test analyze tool
./my-agent analyze_item /path/to/specific/item | jq .

# Test with parameters
PARAM_DEPTH=5 ./my-agent analyze_item /path/to/item | jq .
```

### Integration Testing

Test with the DocLoom executor:

```go
func TestAgentIntegration(t *testing.T) {
    executor := agent.NewExecutor(registry, cache, logger)
    
    // Test tool invocation
    output, err := executor.RunTool("my-analyzer", "list_items", 
        map[string]string{
            "SOURCE_PATH": "/test/repo",
        })
    
    require.NoError(t, err)
    
    var result map[string]interface{}
    err = json.Unmarshal([]byte(output), &result)
    require.NoError(t, err)
    
    assert.Contains(t, result, "items")
}
```

## Performance Considerations

### Caching

Implement caching for expensive operations:

```go
type Cache struct {
    results map[string]*AnalysisResult
    mu      sync.RWMutex
}

func (c *Cache) Get(key string) (*AnalysisResult, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    result, ok := c.results[key]
    return result, ok
}
```

### Parallel Processing

Use concurrency for independent operations:

```go
func analyzeFiles(files []string) []Result {
    results := make([]Result, len(files))
    var wg sync.WaitGroup
    
    for i, file := range files {
        wg.Add(1)
        go func(idx int, path string) {
            defer wg.Done()
            results[idx] = analyzeFile(path)
        }(i, file)
    }
    
    wg.Wait()
    return results
}
```

### Progress Reporting

Report progress to stderr for long operations:

```go
for i, item := range items {
    fmt.Fprintf(os.Stderr, "Processing %d/%d: %s\n", 
        i+1, len(items), item.Name)
    processItem(item)
}
```

## Advanced Features

### Tool Chaining

Design tools that work well together:

```yaml
tools:
  - name: discover_services
    description: Finds all microservices in the repository
    
  - name: analyze_service
    description: Analyzes a specific service (call discover_services first)
    
  - name: trace_dependencies
    description: Traces dependencies between services
```

### Incremental Analysis

Support incremental updates:

```go
func analyzeIncremental(path string, since time.Time) {
    // Only analyze files modified after 'since'
    files := findModifiedFiles(path, since)
    results := analyzeFiles(files)
    
    // Merge with cached results
    mergeResults(results)
}
```

### Multi-Format Support

Handle different input formats:

```go
func detectFormat(path string) Format {
    switch filepath.Ext(path) {
    case ".json":
        return JSONFormat
    case ".yaml", ".yml":
        return YAMLFormat
    case ".xml":
        return XMLFormat
    default:
        return AutoDetect
    }
}
```

## Deployment

### Binary Distribution

1. Build for target platforms:
```bash
GOOS=linux GOARCH=amd64 go build -o my-agent-linux
GOOS=darwin GOARCH=amd64 go build -o my-agent-darwin
GOOS=windows GOARCH=amd64 go build -o my-agent.exe
```

2. Include in agent directory:
```
agents/my-analyzer/
├── agent.agent.yaml
├── my-agent-linux
├── my-agent-darwin
├── my-agent.exe
└── README.md
```

### Container Packaging

Create a Dockerfile:

```dockerfile
FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o my-agent cmd/my-agent/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/my-agent /usr/local/bin/
COPY agent.agent.yaml /etc/my-agent/
ENTRYPOINT ["my-agent"]
```

## Debugging

### Verbose Logging

Add debug output to stderr:

```go
if verbose {
    fmt.Fprintf(os.Stderr, "[DEBUG] Processing file: %s\n", file)
    fmt.Fprintf(os.Stderr, "[DEBUG] Found %d items\n", len(items))
}
```

### Tool Testing Script

Create a test script:

```bash
#!/bin/bash
# test-agent.sh

echo "Testing list_items..."
./my-agent list_items ./test-data | jq .

echo "Testing analyze_item..."
./my-agent analyze_item ./test-data/sample.txt | jq .

echo "Testing with parameters..."
PARAM_DEPTH=5 PARAM_VERBOSE=true \
  ./my-agent analyze_item ./test-data/complex.txt | jq .
```

## Common Patterns

### Repository Analysis Pattern

```go
type RepoAnalyzer struct {
    // Shared state
}

func (r *RepoAnalyzer) ListProjects() []Project
func (r *RepoAnalyzer) AnalyzeProject(path string) ProjectAnalysis
func (r *RepoAnalyzer) CompareProjects(a, b string) Comparison
```

### File Processing Pattern

```go
func processFiles(root string, pattern string) {
    filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
            processFile(path)
        }
        return nil
    })
}
```

### Configuration Pattern

```go
type Config struct {
    // Loaded from environment
}

func LoadConfig() Config {
    return Config{
        SourcePath: os.Getenv("PARAM_SOURCE_PATH"),
        MaxDepth:   getEnvInt("PARAM_MAX_DEPTH", 10),
        // ...
    }
}
```

## Summary

Creating effective agents for DocLoom requires:

1. **Clear Tool Design**: Each tool should have a single, well-defined purpose
2. **Descriptive Documentation**: Help the AI understand when and how to use each tool
3. **Structured Output**: Consistent JSON output for reliable parsing
4. **Error Handling**: Graceful failures with informative messages
5. **Performance**: Efficient processing with caching and parallelism
6. **Testing**: Comprehensive tests for each tool and integration scenarios

By following these guidelines, you'll create agents that integrate seamlessly with DocLoom's AI-orchestrated document generation system.