# STORY-6.2 Evidence: Integration with Generate Command

## Implementation Summary
Successfully integrated agent execution into the generate command workflow.

## Test Results

### TC-21.1: Generate Command with Agent E2E
- **Test File**: `internal/cli/generate_agent_integration_test.go:TestGenerateCmd_WithAgent_E2E`
- **Result**: PASS
- **Evidence**: Full workflow tested - agent discovery, execution, artifact generation, and source replacement

## Integration Points

1. **CLI Flags Added**
   - `--agent <name>`: Specify agent to run
   - `--agent-param <key=value>`: Pass parameters (multi-use)

2. **Workflow Implementation** (`internal/cli/generate.go`)
   - Parse agent flags and parameters
   - Discover agent via registry
   - Execute agent with source path
   - Replace sources with artifact cache directory
   - Continue with normal generation flow

## Documentation Updates

1. **README.md**
   - Added "Using Research Agents" section with examples
   - Shows agent invocation with parameters

2. **docs/agents/registry.md**
   - Added workflow diagram (Mermaid)
   - Documented complete execution flow
   - Provided example scenarios

## Commit
- Hash: 7f5f79a
- Message: "feat(generate): integrate agent execution into generate command"