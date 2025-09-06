# C# Analyzer Agent

## Overview

The C# Analyzer is a multi-tool agent that provides comprehensive analysis capabilities for C# repositories. It uses Tree-sitter for robust parsing and exposes multiple tools that can be invoked independently by the AI during document generation.

## Tools

The agent provides the following tools:

### summarize_readme

**Purpose**: Finds and summarizes README files in the repository

**Usage**: `docloom-agent-csharp summarize_readme <path>`

**Output**: JSON object containing:
- `readmeCount`: Number of README files found
- `readmePaths`: Array of paths to README files
- `summary`: Brief summary text

### list_projects

**Purpose**: Lists all C# project files (.csproj) in the repository

**Usage**: `docloom-agent-csharp list_projects <path>`

**Output**: JSON object containing:
- `projectCount`: Number of projects found
- `projects`: Array of relative paths to .csproj files

### get_dependencies

**Purpose**: Analyzes project dependencies from .csproj files

**Usage**: `docloom-agent-csharp get_dependencies <path>`

**Output**: JSON object containing:
- `dependencies`: Map of project paths to their NuGet dependencies
- `summary`: Analysis summary

### get_api_surface

**Purpose**: Extracts the public API surface of the codebase

**Usage**: `docloom-agent-csharp get_api_surface <path>`

**Output**: JSON object containing:
- `summary`: Statistics about namespaces, classes, methods, etc.
- `apiSurface`: Detailed API structure with namespaces, classes, methods, and properties

### get_file_content

**Purpose**: Retrieves the content of a specific file

**Usage**: `docloom-agent-csharp get_file_content <file_path>`

**Output**: JSON object containing:
- `path`: File path
- `size`: File size in bytes
- `content`: File content as string

## Legacy Mode

For backward compatibility, the agent still supports the original analyze mode:

**Usage**: `docloom-agent-csharp <source_path> <output_path>`

This mode generates three markdown files:
- `ProjectSummary.md`: High-level project statistics
- `ApiSurface.md`: Detailed API documentation
- `ArchitecturalInsights.md`: Detected patterns and recommendations

## Environment Parameters

The agent accepts parameters through environment variables:

- `PARAM_INCLUDE_INTERNAL`: Include internal classes (default: false)
- `PARAM_MAX_DEPTH`: Maximum parsing depth (default: 10)
- `PARAM_EXTRACT_METRICS`: Extract code metrics (default: true)
- `PARAM_SOURCE_PATH`: Path to source repository
- `PARAM_FILE_PATH`: Specific file path (for get_file_content)

## Agent Definition

The agent should be configured in `agent.agent.yaml`:

```yaml
apiVersion: docloom.io/v1alpha1
kind: Agent
metadata:
  name: csharp-analyzer
  description: Analyzes C# repositories for architecture and code quality
spec:
  tools:
    - name: summarize_readme
      description: Finds and summarizes README files in the repository. Returns JSON with readme paths and count.
      command: docloom-agent-csharp
      args:
        - summarize_readme
        - "${SOURCE_PATH}"
    
    - name: list_projects
      description: Lists all C# project files (.csproj) in the repository. Returns JSON array of relative paths. Use this first to understand repository structure.
      command: docloom-agent-csharp
      args:
        - list_projects
        - "${SOURCE_PATH}"
    
    - name: get_dependencies
      description: Analyzes all .csproj files to extract NuGet package dependencies. Returns a map of project paths to their dependencies. Call list_projects first.
      command: docloom-agent-csharp
      args:
        - get_dependencies
        - "${SOURCE_PATH}"
    
    - name: get_api_surface
      description: Analyzes the public API surface of the codebase. Returns detailed namespace, class, method, and property information in JSON format.
      command: docloom-agent-csharp
      args:
        - get_api_surface
        - "${SOURCE_PATH}"
    
    - name: get_file_content
      description: Retrieves the content of a specific file. Useful for examining particular source files identified by other tools.
      command: docloom-agent-csharp
      args:
        - get_file_content
        - "${FILE_PATH}"
  
  parameters:
    - name: source_path
      type: string
      required: true
      description: Path to the source code repository
    
    - name: file_path
      type: string
      required: false
      description: Path to a specific file (used by get_file_content tool)
```

## Implementation Details

### Parser

The agent uses Tree-sitter with the C# grammar to parse source files. The parser extracts:

- Namespaces
- Classes (including interfaces and abstract classes)
- Methods with signatures
- Properties with types
- XML documentation comments
- Access modifiers (public, private, etc.)

### Pattern Detection

The agent detects common architectural patterns:

- Repository Pattern
- Factory Pattern
- Service Layer
- MVC/API Controllers

Detection is based on class naming conventions and structural analysis.

### Performance Considerations

- Skips common non-source directories (bin, obj, packages)
- Processes files in parallel where possible
- Caches parsing results within a single invocation
- Outputs structured JSON for efficient processing

## Error Handling

- Invalid file paths return error JSON
- Parse errors are logged to stderr and skipped
- Non-zero exit codes indicate tool failure
- Timeout protection for long-running operations

## Integration with DocLoom

The C# Analyzer integrates seamlessly with DocLoom's AI-orchestrated analysis loop:

1. The AI can invoke any tool based on the template's analysis needs
2. Tools can be called in sequence to build understanding
3. Output from one tool can inform the next tool call
4. The final document incorporates insights from all tool invocations

## Development

### Building

```bash
go build -o docloom-agent-csharp cmd/docloom-agent-csharp/main.go
```

### Testing Individual Tools

```bash
# List projects
./docloom-agent-csharp list_projects /path/to/csharp/repo

# Get API surface
./docloom-agent-csharp get_api_surface /path/to/csharp/repo

# Get specific file
./docloom-agent-csharp get_file_content /path/to/file.cs
```

### Adding New Tools

1. Add a new cobra command in `main.go`
2. Implement the tool logic
3. Output JSON to stdout
4. Update the agent.agent.yaml with the new tool
5. Document the tool in this file