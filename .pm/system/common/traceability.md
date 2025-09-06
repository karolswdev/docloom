# Traceability Matrix

## Overview
This document maintains the traceability links between requirements, test cases, and evidence for the DocLoom project.

## Format
Each entry follows the format:
- **Requirement ID** → **Test Case(s)** → **Evidence Path(s)**

## Phase 5 Traceability

### STORY-5.1: Agent Definition & Registry
- **PROD-015** (Agent Definition Files) → TC-17.1 → internal/agent/types_test.go:TestAgentDefinition_ParseYAML
- **TECH-P-008** (Agent Definition Files) → TC-17.1 → internal/agent/types_test.go:TestAgentDefinition_ParseYAML
- **ARCH-009** (Agent Registry & Discovery) → TC-18.1 → internal/agent/registry_test.go:TestAgentRegistry_DiscoverAgents

### STORY-5.2: CLI for Agent Management
- **USER-009** (Agent Management CLI) → TC-19.1 → internal/cli/agents_test.go:TestAgentsListCmd_E2E
- **USER-009** (Agent Management CLI) → TC-19.2 → internal/cli/agents_test.go:TestAgentsDescribeCmd_E2E

## Phase 6 Traceability

### STORY-6.1: Agent Execution Engine
- **ARCH-008** (Decoupled Executor) → TC-20.1 → internal/agent/executor_test.go:TestAgentExecutor_RunCommand
- **ARCH-010** (Caching) → TC-20.1 → internal/agent/cache.go + internal/agent/executor_test.go:TestAgentExecutor_RunCommand
- **USER-010** (Parameter Overrides) → TC-22.1 → internal/agent/executor_test.go:TestAgentExecutor_ParameterOverrides

### STORY-6.2: Integration with generate Command
- **PROD-014** (Integrated Agent Execution) → TC-21.1 → internal/cli/generate_agent_integration_test.go:TestGenerateCmd_WithAgent_E2E
- **PROD-016** (Agent Workflow) → TC-21.1 → internal/cli/generate_agent_integration_test.go:TestAgentWorkflowIntegration
- **USER-008** (Agent CLI Flags) → TC-21.1 → internal/cli/generate.go (--agent, --agent-param flags)
- **USER-010** (Parameter Overrides) → TC-22.1 → internal/agent/executor_test.go:TestAgentExecutor_ParameterOverrides

## Phase 7 Traceability

### STORY-7.1: Integrating the Tree-sitter C# Parser
- **PROD-017** (Initial C# Analyzer Agent) → TC-23.1 → internal/agents/csharp/parser/parser_test.go:TestCSharpParser_ExtractAPISurface

### STORY-7.2: Building the Executable Agent and Final Integration  
- **PROD-017** (Initial C# Analyzer Agent) → TC-23.2 → internal/cli/generate_agent_csharp_test.go:TestCSharpAgent_E2E_Integration

## Phase 8 Traceability

### STORY-8.1: Evolving the Agent into a Toolkit
- **Agent-as-Toolkit** (Implied by Vision) → TC-24.1 → internal/agent/executor_tool_test.go:TestAgentExecutor_RunTool
- **Documentation:** docs/agents/schema.md (tool-based architecture), docs/agents/csharp-analyzer.md (multi-tool docs), docs/agents/authoring-guide.md (tool paradigm)

### STORY-8.2: The AI Analysis Loop
- **LLM-Orchestrated Analysis Loop** (Implied by Vision) → TC-25.1 → internal/generate/orchestrator_analysis_test.go:TestOrchestrator_AnalysisLoop (pending full integration)
- **Documentation:** docs/architecture/analysis-loop.md (comprehensive workflow with Mermaid diagram)

### STORY-8.3: Template-Driven Intelligence
- **Template-Defined Intelligence** (Implied by Vision) → TC-26.1 → internal/cli/generate_template_analysis_test.go:TestGenerateCmd_UsesTemplateAnalysisPrompt
- **Documentation:** docs/guides/creating-intelligent-templates.md (world-class guide on analysis prompts)

## Phase 9 Traceability

### STORY-9.1: Defining the Claude Code Agent Contract
- **PROD-018** (External CLI Tool Agent - Contract) → TC-27.1 → docs/agents/artifact-spec-claude-code.md + .docloom/agents/csharp-cc-cli.agent.yaml

### STORY-9.2: Mock Implementation and E2E Testing
- **PROD-018** (External CLI Tool Agent - Mock) → TC-27.2 → internal/cli/generate_agent_claude_test.go:TestCSharpCCAgent_E2E_WithMock

## Phase 10 Traceability

### STORY-10.1: cc-cli Scaffolding and Core Logic
- **PROD-018** (External CLI Tool Agent - Real Implementation) → TC-28.1 → tools/claude-code-cli/cmd/*_test.go (Unit tests for scanner, prompt, claude, writer)
  - Evidence: evidence/PHASE-GO-10/STORY-10.1/test-output.txt

### STORY-10.2: Claude LLM Interaction and Artifact Generation
- **PROD-018** (External CLI Tool Agent - Real Implementation) → TC-28.2 → tools/claude-code-cli acceptance test
  - Evidence: evidence/PHASE-GO-10/STORY-10.2/test-tc-28.2-run.txt
- **PROD-018** (External CLI Tool Agent - Real Implementation) → TC-28.3 → Full E2E integration with main docloom
  - Evidence: evidence/PHASE-GO-10/final-regression.txt

## Previous Phases
[Previous phase traceability would be maintained here]