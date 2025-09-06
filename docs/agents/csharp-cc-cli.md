# Claude Code CLI Agent

## Overview

The `csharp-cc-cli` agent orchestrates the powerful external **Claude Code CLI (`cc-cli`)** tool to perform comprehensive code analysis and generate structured insights about your codebase. This agent serves as a bridge between DocLoom's documentation generation pipeline and the advanced analysis capabilities of Claude Code.

## Purpose

While DocLoom excels at generating documentation from source materials, the Claude Code CLI provides deep code analysis that goes beyond simple parsing. It detects patterns, anti-patterns, calculates complexity metrics, maps dependencies, and provides actionable recommendations. By integrating this analysis into DocLoom's workflow, we can generate documentation that is:

- **Accurate**: Based on actual code analysis, not assumptions
- **Insightful**: Includes technical debt identification and improvement recommendations
- **Comprehensive**: Covers architecture, dependencies, complexity, and API surface
- **Actionable**: Provides specific guidance for code improvements

## How It Works

1. **DocLoom invokes the agent** with your source code path and desired parameters
2. **The agent executes Claude Code CLI** which analyzes the codebase
3. **Claude Code generates artifacts** in a structured directory format (see [Artifact Specification](./artifact-spec-claude-code.md))
4. **DocLoom ingests these artifacts** and includes them in the AI prompt
5. **The AI model uses the analysis** to generate informed, accurate documentation

## Configuration

The agent is defined in `.docloom/agents/csharp-cc-cli.agent.yaml` and accepts the following parameters:

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `language` | string | "csharp" | Programming language hint (e.g., csharp, go, python) |
| `depth` | integer | 2 | Analysis depth level (1-3, where 3 is most detailed) |
| `include-tests` | boolean | false | Include test files in the analysis |
| `timeout` | integer | 300 | Maximum execution time in seconds |

## Usage

### Basic Usage

```bash
docloom generate \
  --type architecture-vision \
  --source ./src \
  --agent csharp-cc-cli \
  --out architecture.html
```

### With Custom Parameters

```bash
docloom generate \
  --type technical-debt-summary \
  --source ./src \
  --agent csharp-cc-cli \
  --agent-param depth=3 \
  --agent-param include-tests=true \
  --out tech-debt.html
```

### Dry Run Mode

To preview what the agent will analyze without actually running it:

```bash
docloom generate \
  --type architecture-vision \
  --source ./src \
  --agent csharp-cc-cli \
  --dry-run
```

## Output Structure

The Claude Code CLI generates a comprehensive set of artifacts. See the complete [Artifact Specification](./artifact-spec-claude-code.md) for details. Key outputs include:

- **overview.md**: High-level project summary
- **structure.json**: Hierarchical code structure
- **dependencies.json**: Internal and external dependency mapping
- **complexity.json**: Code complexity metrics and hotspots
- **api-surface.json**: Public API documentation
- **insights/**: Directory containing pattern detection and recommendations

## Integration with Templates

The `csharp-cc-cli` agent works particularly well with these DocLoom templates:

### Architecture Vision
Combines code analysis with architectural documentation to create comprehensive technical vision documents.

### Technical Debt Summary
Uses complexity metrics and anti-pattern detection to generate actionable technical debt reports.

### API Documentation
Leverages the extracted API surface to create accurate API documentation.

### Code Review Report
Synthesizes all analysis insights into comprehensive code review documentation.

## Mock Implementation

For testing and development, a mock implementation (`mock-cc-cli.sh`) simulates the Claude Code CLI behavior. This allows:

- Testing the integration without the actual Claude Code CLI
- Validating the artifact ingestion pipeline
- Developing new features against a predictable baseline

To use the mock implementation, the agent definition points to `./mock-cc-cli.sh` instead of the production CLI.

## Best Practices

1. **Choose the right depth level**
   - Level 1: Quick overview for high-level documentation
   - Level 2: Balanced analysis for most use cases
   - Level 3: Deep analysis for comprehensive technical documentation

2. **Include tests when relevant**
   - Enable `include-tests` for test coverage documentation
   - Disable for production code-only analysis

3. **Cache results**
   - The agent caches results for 1 hour by default
   - Results are keyed by source hash and parameters

4. **Combine with other agents**
   - Use multiple agents for different perspectives
   - The csharp-analyzer agent provides complementary tree-sitter based analysis

## Troubleshooting

### Agent not found
Ensure the agent definition file exists in `.docloom/agents/csharp-cc-cli.agent.yaml`

### Timeout errors
Increase the timeout parameter for large codebases:
```bash
--agent-param timeout=600
```

### Missing artifacts
Check that the Claude Code CLI is properly installed and accessible in your PATH

### Mock vs Production
To switch between mock and production:
1. Edit `.docloom/agents/csharp-cc-cli.agent.yaml`
2. Change `command` from `./mock-cc-cli.sh` to `cc-cli`

## Future Enhancements

The Claude Code CLI and this agent are under active development. Planned enhancements include:

- Security vulnerability analysis
- License compliance checking
- Performance profiling integration
- Test coverage mapping
- Documentation quality metrics
- Multi-language support expansion

## Related Documentation

- [Artifact Specification](./artifact-spec-claude-code.md) - Detailed specification of Claude Code CLI outputs
- [Agent Authoring Guide](./authoring-guide.md) - How to create and test agents
- [Agent Schema](./schema.md) - Agent definition file format
- [C# Analyzer Agent](./csharp-analyzer.md) - Alternative tree-sitter based analyzer