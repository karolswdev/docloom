# Claude Code CLI (cc-cli)

## Overview

The Claude Code CLI (`cc-cli`) is a powerful, standalone tool that uses the Claude LLM to perform deep analysis of C# repositories. It generates structured artifacts that provide rich, accurate, and synthesized context for fully automated document generation.

## Purpose

This tool is designed to work as an agent for the DocLoom documentation generation system. It analyzes C# codebases and produces comprehensive insights about:

- Project architecture and structure
- API endpoints and data models  
- Dependencies and frameworks
- Technical debt and recommendations
- Security considerations
- Testing coverage and deployment configuration

## Installation

### From Source

```bash
cd tools/claude-code-cli
go build -o cc-cli
```

### From Release

Download the pre-built binary from the releases page (when available).

## Usage

### Basic Command

```bash
cc-cli --repo-path /path/to/csharp/repo --output-path /path/to/output
```

### Full Options

```bash
cc-cli \
  --repo-path /path/to/csharp/repo \
  --output-path /path/to/output \
  --api-key YOUR_CLAUDE_API_KEY \
  --model claude-3-opus-20240229 \
  --max-tokens 4096 \
  --verbose
```

### Environment Variables

- `CLAUDE_API_KEY`: Your Claude API key (alternative to --api-key flag)

## Output Structure

The tool generates artifacts conforming to the DocLoom Agent Artifact Specification:

```
output-path/
├── metadata.json                 # Agent metadata and execution info
├── analysis/
│   ├── summary.json             # Complete structured analysis
│   └── project-overview.md      # Human-readable project summary
├── repository-context/
│   ├── architecture.md          # Architecture documentation
│   └── api-endpoints.md         # API documentation (if applicable)
└── technical-insights/
    ├── technical-debt.md        # Technical debt analysis
    └── recommendations.md       # Improvement recommendations
```

### Output Format Details

#### metadata.json
Contains execution metadata including agent version, timestamp, and model used.

#### analysis/summary.json
Complete structured JSON response with all analysis data:
- Project name and description
- Architecture patterns and layers
- Dependencies and frameworks
- Features and APIs
- Testing and deployment information
- Security considerations
- Technical debt items
- Recommendations

#### Project Documentation Files
Human-readable markdown files generated from the analysis for easy review and inclusion in documentation.

## Integration with DocLoom

This tool is designed to be used as an agent within the DocLoom system. The corresponding agent definition file (`csharp-cc-cli.agent.yaml`) should be configured to point to this tool:

```yaml
apiVersion: v1
kind: Agent
metadata:
  name: csharp-cc-cli
  description: Claude Code CLI for deep C# repository analysis
spec:
  runner:
    command: cc-cli
    args:
      - "--repo-path"
      - "{{.SourcePath}}"
      - "--output-path"  
      - "{{.OutputPath}}"
  parameters:
    - name: model
      description: Claude model to use
      default: claude-3-opus-20240229
    - name: max-tokens
      description: Maximum tokens for response
      default: 4096
```

## How It Works

1. **Repository Scanning**: The tool scans the target C# repository for key files:
   - Solution files (*.sln)
   - Project files (*.csproj)
   - README files
   - Configuration files
   - Key source files (Program.cs, Startup.cs, Controllers, Services, etc.)

2. **Prompt Generation**: Constructs a sophisticated prompt that includes:
   - Repository structure overview
   - File contents (with token limits)
   - Specific analysis instructions

3. **Claude Analysis**: Sends the prompt to Claude LLM for deep analysis
   - Extracts architectural patterns
   - Identifies features and APIs
   - Analyzes technical debt
   - Provides recommendations

4. **Artifact Generation**: Writes structured artifacts to disk
   - JSON for machine processing
   - Markdown for human readability
   - Follows DocLoom specification exactly

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o cc-cli
```

## Requirements

- Go 1.22 or higher
- Claude API key with access to Claude 3 models
- Target repository must be a valid C# project

## Limitations

- File size limited to 1MB per file to prevent memory issues
- Token limits apply based on Claude model constraints
- Currently focused on C# repositories only

## License

Part of the DocLoom project. See main project LICENSE for details.