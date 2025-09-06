# Agent Definition Schema

## Overview

Research Agents in DocLoom are defined using YAML files with a `.agent.yaml` or `.agent.yml` extension. This document describes the official schema for agent definition files, which serves as the contract for agent developers.

## Schema Structure

### Top-Level Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `apiVersion` | string | Yes | The API version for the agent schema. Currently must be `v1`. |
| `kind` | string | Yes | The resource type. Must be `ResearchAgent`. |
| `metadata` | object | Yes | Contains basic information about the agent. |
| `spec` | object | Yes | Defines the agent's execution specification. |

### Metadata Section

The `metadata` section contains identifying information:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Unique identifier for the agent (lowercase, hyphens allowed). |
| `description` | string | No | Human-readable description of what the agent does. |

### Spec Section

The `spec` section defines how the agent executes:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `runner` | object | Yes | Specifies the command to execute the agent. |
| `parameters` | array | No | List of input parameters the agent accepts. |

#### Runner Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `command` | string | Yes | The executable command (e.g., `python3`, `node`, `./script.sh`). |
| `args` | array | No | Additional arguments to pass to the command. |

#### Parameter Object

Each parameter in the `parameters` array has:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Parameter name (used as environment variable). |
| `type` | string | Yes | Data type (`string`, `integer`, `boolean`, `float`). |
| `description` | string | No | Description of the parameter's purpose. |
| `required` | boolean | No | Whether the parameter must be provided (default: false). |
| `default` | any | No | Default value if not provided by user. |

## Complete Example

```yaml
apiVersion: v1
kind: ResearchAgent
metadata:
  name: academic-research-agent
  description: Performs academic literature research on specified topics
spec:
  runner:
    command: python3
    args:
      - /agents/academic_research.py
  parameters:
    - name: topic
      type: string
      description: The research topic to investigate
      required: true
    - name: depth
      type: integer
      description: How many levels deep to search (1-5)
      required: false
      default: 3
    - name: include_citations
      type: boolean
      description: Whether to include full citations
      required: false
      default: true
    - name: max_results
      type: integer
      description: Maximum number of results to return
      required: false
      default: 100
```

## Parameter Passing

When an agent is executed, parameters are passed as environment variables with the prefix `PARAM_`. For example:
- `topic` becomes `PARAM_TOPIC`
- `max_results` becomes `PARAM_MAX_RESULTS`

The agent also receives two command-line arguments:
1. Source path - The path to the codebase being analyzed
2. Output path - Where the agent should write its artifacts

## Validation Rules

1. **API Version**: Must be a supported version (currently only `v1`).
2. **Kind**: Must be exactly `ResearchAgent`.
3. **Name**: Must be unique within the registry, lowercase, may contain hyphens.
4. **Command**: Must be a valid executable available in the system PATH or an absolute path.
5. **Parameter Types**: Must be one of `string`, `integer`, `boolean`, or `float`.
6. **Default Values**: Must match the declared type.

## File Naming Convention

Agent definition files must follow this naming pattern:
- `<agent-name>.agent.yaml`
- `<agent-name>.agent.yml`

The `<agent-name>` portion should match the `metadata.name` field for consistency.

## Discovery Locations

DocLoom searches for agent definitions in the following locations (in order):
1. `.docloom/agents/` - Project-specific agents
2. `~/.docloom/agents/` - User-wide agents
3. Custom paths added via configuration

## Best Practices

1. **Descriptive Names**: Use clear, descriptive names that indicate the agent's purpose.
2. **Comprehensive Descriptions**: Provide detailed descriptions for both the agent and its parameters.
3. **Sensible Defaults**: Include reasonable default values for optional parameters.
4. **Type Safety**: Always specify the correct type for parameters to enable proper validation.
5. **Documentation**: Include a README in your agent directory explaining setup and usage.
6. **Version Compatibility**: Always specify `apiVersion: v1` for compatibility.

## Error Handling

If an agent definition file is invalid, it will be skipped during discovery with an error message indicating:
- Missing required fields
- Invalid field values
- Malformed YAML syntax
- Unsupported API versions or kinds

## Future Enhancements

The schema may be extended in future versions to support:
- Multiple runner types (Docker, WebAssembly, etc.)
- Resource limits (CPU, memory, timeout)
- Output format specifications
- Dependency declarations
- Authentication mechanisms