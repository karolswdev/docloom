# Agent Authoring Guide

This guide explains how to create custom Research Agents for DocLoom.

## Overview

A Research Agent is an external program that analyzes source code and produces documentation artifacts. Agents receive source paths as command-line arguments and parameters as environment variables.

## Agent Contract

### Input

Agents receive two primary inputs:

1. **Source Path** (arg 1): The directory or file to analyze
2. **Output Path** (arg 2): The directory where the agent should write its output

### Parameters

Additional parameters are passed as environment variables prefixed with `PARAM_`:

- `PARAM_DEBUG`: Debug mode flag
- `PARAM_MAX_DEPTH`: Maximum analysis depth
- `PARAM_<NAME>`: Any custom parameter defined in your agent definition

### Output

Agents must produce at least one markdown file (`.md`) in the output directory. Common outputs include:

- `analysis.md`: Main analysis report
- `metrics.md`: Code metrics and statistics
- `recommendations.md`: Improvement suggestions

## Agent Definition File

Create a `.agent.yaml` file to define your agent:

```yaml
apiVersion: v1
kind: Agent
metadata:
  name: my-agent
  description: Analyzes code and produces insights
spec:
  runner:
    command: /usr/bin/python3
    args: 
      - /path/to/my-agent.py
      - ${SOURCE_PATH}
      - ${OUTPUT_PATH}
  parameters:
    - name: depth
      type: integer
      required: false
      default: 3
      description: Analysis depth
    - name: verbose
      type: boolean
      required: false
      default: false
      description: Enable verbose output
```

## Example: Shell Script Agent

Here's a simple shell script agent:

```bash
#!/bin/bash
# simple-analyzer.sh - A basic code analysis agent

SOURCE_PATH="$1"
OUTPUT_PATH="$2"

# Read parameters from environment
DEPTH="${PARAM_DEPTH:-3}"
VERBOSE="${PARAM_VERBOSE:-false}"

echo "Analyzing: $SOURCE_PATH" >&2
echo "Output to: $OUTPUT_PATH" >&2

# Create analysis report
cat > "$OUTPUT_PATH/analysis.md" << EOF
# Code Analysis Report

## Source Information
- **Path**: $SOURCE_PATH
- **Analysis Depth**: $DEPTH

## File Statistics
$(find "$SOURCE_PATH" -type f -name "*.go" | wc -l) Go files found

## Code Metrics
- Lines of Code: $(find "$SOURCE_PATH" -name "*.go" -exec wc -l {} + | tail -1 | awk '{print $1}')
- Test Files: $(find "$SOURCE_PATH" -name "*_test.go" | wc -l)

## Generated at
$(date)
EOF

echo "Analysis complete" >&2
exit 0
```

Save this as `simple-analyzer.sh` and create the agent definition:

```yaml
apiVersion: v1
kind: Agent
metadata:
  name: simple-analyzer
  description: Simple code analyzer
spec:
  runner:
    command: /path/to/simple-analyzer.sh
  parameters:
    - name: depth
      type: integer
      required: false
      default: 3
      description: Analysis depth
```

## Example: Python Agent

For more complex analysis, use Python:

```python
#!/usr/bin/env python3
# code-insights.py - Advanced code analysis agent

import os
import sys
import json
from pathlib import Path

def main():
    if len(sys.argv) < 3:
        print("Usage: code-insights.py <source_path> <output_path>", file=sys.stderr)
        sys.exit(1)
    
    source_path = Path(sys.argv[1])
    output_path = Path(sys.argv[2])
    
    # Read parameters
    max_files = int(os.environ.get('PARAM_MAX_FILES', '100'))
    include_tests = os.environ.get('PARAM_INCLUDE_TESTS', 'false').lower() == 'true'
    
    # Perform analysis
    go_files = list(source_path.rglob('*.go'))
    if not include_tests:
        go_files = [f for f in go_files if not f.name.endswith('_test.go')]
    
    # Generate report
    report = f"""# Code Insights Report

## Repository Overview
- **Total Go Files**: {len(go_files)}
- **Analysis Parameters**:
  - Max Files: {max_files}
  - Include Tests: {include_tests}

## File List
"""
    
    for i, file in enumerate(go_files[:max_files]):
        rel_path = file.relative_to(source_path)
        report += f"- `{rel_path}`\n"
    
    # Write output
    output_file = output_path / 'insights.md'
    output_file.write_text(report)
    
    print(f"Analysis complete. Report written to {output_file}", file=sys.stderr)
    return 0

if __name__ == '__main__':
    sys.exit(main())
```

## Best Practices

1. **Error Handling**: Always validate inputs and handle errors gracefully
2. **Logging**: Write progress messages to stderr, not stdout
3. **Output Validation**: Ensure at least one `.md` file is created
4. **Performance**: Be mindful of large repositories; implement reasonable limits
5. **Documentation**: Include clear descriptions in your agent definition
6. **Testing**: Test your agent with various repository sizes and structures

## Installation

Place your agent definition files in one of these locations:

- `$HOME/.docloom/agents/` - User-specific agents
- `<project>/.docloom/agents/` - Project-specific agents

DocLoom will automatically discover and register agents from these directories.

## Debugging

Use the verbose flag to see agent execution details:

```bash
docloom generate --agent my-agent --source ./src --verbose
```

Check the agent's stderr output for debugging information.

## Advanced Features

### Placeholder Variables

In your agent's `args` configuration, you can use:

- `${SOURCE_PATH}`: Replaced with the actual source path
- `${OUTPUT_PATH}`: Replaced with the cache directory path

### Parameter Types

Supported parameter types:

- `string`: Text values
- `integer`: Numeric values
- `boolean`: true/false values
- `float`: Decimal numbers

### Exit Codes

- `0`: Success
- Non-zero: Failure (DocLoom will report the error)

## Security Considerations

- Agents run with the same permissions as DocLoom
- Validate and sanitize all inputs
- Be cautious with file system operations
- Never expose sensitive information in outputs