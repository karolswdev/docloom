# Agent Registry & Discovery

## Overview

The Agent Registry is DocLoom's central mechanism for discovering, loading, and managing Research Agents. It automatically scans predefined directories for agent definition files and makes them available to the DocLoom CLI.

## Discovery Mechanism

### Search Paths

The registry searches for agents in the following locations, in order of precedence:

1. **Project-Local Agents** (`.docloom/agents/`)
   - Agents specific to the current project
   - Highest precedence - overrides user and system agents
   - Ideal for project-specific research tools

2. **User-Home Agents** (`~/.docloom/agents/`)
   - Personal agents available across all projects
   - Shared between all DocLoom projects for the current user
   - Good for commonly used research agents

3. **Custom Paths** (via configuration)
   - Additional directories can be added programmatically
   - Useful for shared team agents or organizational repositories

### Discovery Process

1. **Scanning**: The registry recursively scans each search path
2. **Filtering**: Only files matching `*.agent.yaml` or `*.agent.yml` are considered
3. **Validation**: Each file is parsed and validated against the schema
4. **Registration**: Valid agents are added to the registry's internal map
5. **Conflict Resolution**: If duplicate names exist, the first discovered agent wins (based on search path precedence)

## File Detection

The registry identifies agent files by their extension:
- `.agent.yaml` (preferred)
- `.agent.yml` (alternative)

All other files are ignored during discovery.

## Validation

During discovery, each agent file is validated for:

1. **Syntax**: Valid YAML structure
2. **Required Fields**:
   - `apiVersion` must be present and supported
   - `kind` must be `ResearchAgent`
   - `metadata.name` must be present and unique
3. **Type Checking**: Field values must match expected types
4. **Command Validation**: Runner commands are checked for basic validity

Invalid agents are skipped with error logging, but don't prevent other agents from loading.

## Precedence Rules

When multiple agents with the same name exist in different locations:

```
Project (.docloom/agents/) > User (~/.docloom/agents/) > Custom paths
```

This allows projects to override standard agents with project-specific versions.

## Adding Custom Agents

### For Individual Projects

1. Create the directory structure:
   ```bash
   mkdir -p .docloom/agents
   ```

2. Add your agent definition file:
   ```bash
   cp my-agent.agent.yaml .docloom/agents/
   ```

3. Verify discovery:
   ```bash
   docloom agents list
   ```

### For All Projects (User-Wide)

1. Create the user agents directory:
   ```bash
   mkdir -p ~/.docloom/agents
   ```

2. Add agent definitions:
   ```bash
   cp my-agent.agent.yaml ~/.docloom/agents/
   ```

3. Available in all projects:
   ```bash
   cd /any/project
   docloom agents list
   ```

## Registry API

The registry provides a simple API for agent management:

### Core Methods

- `Discover()` - Scan all search paths and load agents
- `Get(name)` - Retrieve a specific agent by name
- `List()` - Get all discovered agents
- `AddSearchPath(path)` - Add a custom search directory

### Usage Example

```go
// Create a new registry
registry := agent.NewRegistry()

// Add a custom search path
registry.AddSearchPath("/opt/shared-agents")

// Discover all agents
if err := registry.Discover(); err != nil {
    log.Fatal(err)
}

// Get a specific agent
if agent, exists := registry.Get("research-agent"); exists {
    // Use the agent
}

// List all agents
for _, agent := range registry.List() {
    fmt.Printf("%s: %s\n", agent.Metadata.Name, agent.Metadata.Description)
}
```

## Error Handling

The registry handles various error conditions gracefully:

1. **Missing Directories**: Non-existent search paths are silently skipped
2. **Invalid YAML**: Files with syntax errors are skipped with error logging
3. **Schema Violations**: Invalid agents are rejected with descriptive errors
4. **Duplicate Names**: Later discoveries of the same name are ignored
5. **Permission Errors**: Unreadable files are skipped with warnings

## Performance Considerations

- **Lazy Loading**: Agents are discovered on-demand, not at startup
- **Caching**: Discovered agents are cached in memory
- **Minimal I/O**: Only `.agent.yaml/yml` files are read
- **Fast Parsing**: YAML parsing is optimized for agent schemas

## Debugging

To troubleshoot agent discovery issues:

1. **Check File Location**: Ensure agents are in the correct directory
2. **Validate YAML**: Use a YAML validator to check syntax
3. **Verify Extension**: Files must end in `.agent.yaml` or `.agent.yml`
4. **Review Logs**: Check debug logs for discovery errors
5. **List Agents**: Run `docloom agents list` to see what was discovered

## Best Practices

1. **Organize by Function**: Group related agents in subdirectories
2. **Use Descriptive Names**: Agent names should indicate their purpose
3. **Document Dependencies**: Include setup instructions with agents
4. **Version Control**: Track project agents in version control
5. **Test Locally**: Validate agents work before sharing

## Intermediate Artifact Cache

The registry works in conjunction with an artifact cache system to manage temporary files produced by agents.

### Cache Directory Structure

```
/tmp/docloom-agent-cache/
├── agent-name-20240106-150405-12345/
│   ├── analysis.md
│   ├── metrics.md
│   └── recommendations.md
└── another-agent-20240106-151230-12346/
    └── report.md
```

### Cache Lifecycle

1. **Creation**: A unique directory is created for each agent run
2. **Naming**: Directories follow the pattern: `{agent-name}-{timestamp}-{pid}`
3. **Location**: Cache resides in the system temp directory (`/tmp` or `%TEMP%`)
4. **Usage**: Agents write their output files to this directory
5. **Consumption**: DocLoom reads artifacts from the cache for document generation
6. **Cleanup**: Directories older than 24 hours are automatically cleaned

### Cache Management

The cache system provides:

- **Isolation**: Each agent run gets its own directory
- **Uniqueness**: Timestamp and PID prevent collisions
- **Automatic Cleanup**: Old artifacts are purged daily
- **No Manual Intervention**: Users don't need to manage cache files

### Benefits

- **Clean Source Directories**: No intermediate files in user's project
- **Parallel Execution**: Multiple agents can run simultaneously
- **Debugging**: Artifacts persist for inspection if needed
- **Performance**: Fast local disk I/O for artifacts

## Future Enhancements

Planned improvements to the registry include:

- **Remote Repositories**: Fetch agents from Git repositories
- **Agent Versioning**: Support multiple versions of the same agent
- **Dependency Resolution**: Automatic installation of agent requirements
- **Hot Reload**: Detect and load new agents without restart
- **Registry Plugins**: Extensible discovery mechanisms
- **Agent Marketplace**: Central repository for community agents