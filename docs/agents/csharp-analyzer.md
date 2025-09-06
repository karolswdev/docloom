# C# Analyzer Agent

## Overview

The C# Analyzer is a Go-native agent that provides automatic analysis of C# repositories, extracting the public API surface and generating comprehensive documentation. This agent is fully self-contained and requires no external dependencies like the .NET SDK.

## Implementation Details

The agent is built entirely in Go and uses the Tree-sitter parsing library to achieve language analysis without external dependencies. Tree-sitter provides a robust, error-tolerant parser that can handle incomplete or syntactically incorrect code while still extracting meaningful information.

### Key Technologies

- **Language**: Go 1.22+
- **Parser**: Tree-sitter with C# grammar
- **Output Format**: Markdown and JSON

### Architecture

The agent consists of two main components:

1. **Parser Package** (`internal/agents/csharp/parser`): Core parsing logic that uses Tree-sitter to analyze C# source code and extract:
   - Namespaces and their hierarchy
   - Classes, interfaces, and structs
   - Public methods with signatures
   - Properties and their types
   - XML documentation comments

2. **Agent Command** (`cmd/docloom-agent-csharp`): The executable entry point that:
   - Accepts source and output path arguments
   - Recursively scans for `.cs` files
   - Generates structured documentation artifacts

## Output Artifacts

The agent produces the following markdown files in the output directory:

### ProjectSummary.md
A high-level overview of the analyzed C# project including:
- Total number of namespaces
- Count of classes, interfaces, and structs
- Public API statistics
- Key architectural patterns detected

### ApiSurface.md
Detailed documentation of the public API surface:
- Hierarchical namespace structure
- Complete class definitions with:
  - XML documentation comments
  - Public methods and their signatures
  - Properties with types
  - Inheritance relationships
- Interface contracts
- Static utility classes

### ArchitecturalInsights.md
Analysis of architectural patterns and design decisions:
- Detected design patterns (Repository, Factory, etc.)
- Dependency injection usage
- Async/await patterns
- SOLID principle adherence indicators

## Configuration

The agent is configured via the `csharp-analyzer.agent.yaml` file:

```yaml
apiVersion: v1
kind: Agent
metadata:
  name: csharp-analyzer
  description: Go-native C# code analyzer using tree-sitter
spec:
  runner:
    command: ["./docloom-agent-csharp"]
    args: []
  parameters:
    - name: include-internal
      description: Include internal classes in analysis
      type: boolean
      default: false
    - name: max-depth
      description: Maximum namespace depth to analyze
      type: integer
      default: 10
    - name: extract-metrics
      description: Calculate code complexity metrics
      type: boolean
      default: true
```

## Usage

The agent can be invoked through the main docloom CLI:

```bash
# Analyze a C# project
docloom generate --agent csharp-analyzer --source /path/to/csharp/project --out architecture-doc.html

# With custom parameters
docloom generate --agent csharp-analyzer \
  --source /path/to/project \
  --agent-param include-internal=true \
  --agent-param extract-metrics=false \
  --out doc.html
```

## Benefits

1. **No External Dependencies**: Completely self-contained, no .NET SDK required
2. **Fast Analysis**: Tree-sitter provides blazing-fast parsing performance
3. **Error Tolerant**: Can analyze incomplete or work-in-progress code
4. **Cross-Platform**: Works on any platform where Go runs
5. **Integrated Workflow**: Seamlessly integrates with docloom's document generation pipeline