# Agent Definition Schema

## Overview

DocLoom agents are defined using YAML files that specify how they operate and what tools they provide. The agent system has evolved from a simple runner-based architecture to a powerful tool-based paradigm that enables AI-orchestrated analysis.

## Schema Version

Current version: `docloom.io/v1alpha1`

## Tool-Based Architecture

Agents now expose multiple **tools** that can be invoked independently by the AI during analysis. Each tool has:
- A unique name (used as identifier)
- A description (used by the AI to understand when to use the tool)
- A command to execute
- Optional arguments

## File Structure

Agent definitions must be placed in a file named `agent.agent.yaml` within the agent's directory.

```yaml
apiVersion: docloom.io/v1alpha1
kind: Agent
metadata:
  name: agent-name
  description: Human-readable description of the agent
spec:
  # Tools array - defines the capabilities of the agent
  tools:
    - name: tool_name
      description: Clear description for the AI to understand the tool's purpose
      command: /path/to/executable
      args: 
        - subcommand
        - "${PARAMETER_NAME}"
  
  # Parameters that can be passed to tools
  parameters:
    - name: parameter_name
      type: string|int|bool
      required: true|false
      default: default_value
      description: Description of the parameter
```

## Field Definitions

### Root Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `apiVersion` | string | Yes | Schema version (currently `docloom.io/v1alpha1`) |
| `kind` | string | Yes | Resource type (must be `Agent`) |
| `metadata` | object | Yes | Agent metadata |
| `spec` | object | Yes | Agent specification |

### Metadata Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Unique agent identifier |
| `description` | string | Yes | Human-readable description |

### Spec Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `tools` | array | No* | List of tools the agent provides |
| `runner` | object | No* | Legacy runner configuration (deprecated) |
| `parameters` | array | No | Input parameters for the agent |

*Note: Either `tools` or `runner` must be specified. New agents should use `tools`.

### Tool Definition

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Tool identifier (e.g., `list_projects`) |
| `description` | string | Yes | AI-readable description of what the tool does |
| `command` | string | Yes | Executable path or command |
| `args` | array | No | Command arguments (supports parameter substitution) |

### Parameter Definition

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Parameter name |
| `type` | string | Yes | Data type: `string`, `int`, or `bool` |
| `required` | bool | No | Whether the parameter is required (default: false) |
| `default` | any | No | Default value if not provided |
| `description` | string | Yes | Description of the parameter |

## Example: Multi-Tool C# Analyzer

```yaml
apiVersion: docloom.io/v1alpha1
kind: Agent
metadata:
  name: csharp-analyzer
  description: Analyzes C# repositories for architecture and code quality
spec:
  tools:
    - name: summarize_readme
      description: Finds and summarizes README files in the repository
      command: docloom-agent-csharp
      args:
        - summarize_readme
        - "${SOURCE_PATH}"
    
    - name: list_projects
      description: Lists all C# project files (.csproj) in the repository
      command: docloom-agent-csharp
      args:
        - list_projects
        - "${SOURCE_PATH}"
    
    - name: get_dependencies
      description: Extracts NuGet package dependencies from all projects
      command: docloom-agent-csharp
      args:
        - get_dependencies
        - "${SOURCE_PATH}"
    
    - name: get_api_surface
      description: Analyzes the public API surface of the codebase
      command: docloom-agent-csharp
      args:
        - get_api_surface
        - "${SOURCE_PATH}"
    
    - name: get_file_content
      description: Retrieves the content of a specific file
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

## Tool Description Best Practices

The `description` field for each tool is **critical** - it's the primary interface between your agent and the AI. Write clear, actionable descriptions that help the AI understand:

1. **What the tool does** - Be specific about the tool's function
2. **When to use it** - Include context about appropriate use cases
3. **What it returns** - Describe the output format and content
4. **Any prerequisites** - Mention if other tools should be called first

### Good Examples

```yaml
- name: list_projects
  description: "Lists all C# project files (.csproj) in the repository. Returns a JSON array of relative paths. Use this first to understand the repository structure."

- name: get_dependencies
  description: "Analyzes all .csproj files to extract NuGet package dependencies. Returns a map of project paths to their dependencies. Call list_projects first to see available projects."
```

### Poor Examples

```yaml
- name: list_projects
  description: "Lists projects"  # Too vague

- name: get_dependencies
  description: "Gets dependencies from projects"  # Doesn't explain format or prerequisites
```

## Parameter Substitution

Tools support parameter substitution in their arguments using the `${PARAMETER_NAME}` syntax. Parameters can come from:

1. Agent-level parameters (defined in `spec.parameters`)
2. Runtime parameters passed by the executor
3. Special variables like `${SOURCE_PATH}` and `${OUTPUT_PATH}`

## Backward Compatibility

Agents using the legacy `runner` field will continue to work but should be migrated to the tool-based architecture:

```yaml
# Legacy (deprecated)
spec:
  runner:
    command: /path/to/agent
    args: ["${SOURCE_PATH}", "${OUTPUT_PATH}"]

# New (recommended)
spec:
  tools:
    - name: analyze
      description: Performs comprehensive analysis of the codebase
      command: /path/to/agent
      args: ["analyze", "${SOURCE_PATH}", "${OUTPUT_PATH}"]
```

## Tool Execution

When a tool is invoked:

1. The executor looks up the tool by name in the agent definition
2. Parameters are substituted in the command arguments
3. Environment variables are set for all parameters (prefixed with `PARAM_`)
4. The command is executed and output is captured
5. The output (stdout) is returned to the caller

Tools should:
- Output structured data (preferably JSON) to stdout
- Use stderr for logging and diagnostics
- Return non-zero exit codes on failure
- Complete execution within a reasonable timeout