## [x] PHASE-6: Agent Execution & Workflow Integration

---

### **1. Phase Context (What & Why)**

| ID | Title |
| :--- | :--- |
| PHASE-6 | Agent Execution & Workflow Integration |

> **As a** Document Author, **I want** to seamlessly run a research agent on my source code repository as part of the document generation command, **so that** I can get high-quality analysis artifacts without a complex multi-step manual process.

---

### **2. Phase Scope & Test Case Definitions (The Contract)**

*   **Requirement:** **ARCH-008, ARCH-010** - [Decoupled Executor & Caching](./SRS.md#ARCH-008)
    *   **Test Case ID:** `TC-20.1`
        *   **Test Method Signature:** `func TestAgentExecutor_RunCommand(t *testing.T)`
        *   **Test Logic:** (Arrange) Create a simple shell script (`mock-agent.sh`) that accepts a source and output path, copies a file, and prints its parameters to a log. Define a test agent YAML that points to this script. (Act) Run the Agent Executor with the test agent. (Assert) The executor must correctly spawn the process, pass the paths as arguments, pass parameters as environment variables, and the script must successfully create its output artifact.
        *   **Required Proof of Passing:** Test runner output, and verification that the mock agent's output file was created correctly in a temporary cache directory.
*   **Requirement:** **PROD-014, PROD-016, USER-008** - [Integrated Agent Execution](./SRS.md#PROD-014)
    *   **Test Case ID:** `TC-21.1`
        *   **Test Method Signature:** `func TestGenerateCmd_WithAgent_E2E(t *testing.T)`
        *   **Test Logic:** (E2E Test) (Arrange) Use the `mock-agent.sh` from `TC-20.1` that produces a known markdown file. Mock the AI client to expect content from this markdown file. (Act) Run `docloom generate --agent mock-agent --source /path/to/mock/repo ...`. (Assert) The system must first run the agent, then ingest the agent's *output artifact*, and finally pass this artifact content to the AI client in the prompt.
        *   **Required Proof of Passing:** Test runner output showing the full workflow succeeded. Debug logs should show the agent running first, followed by ingestion of the cached artifact.
*   **Requirement:** **USER-010** - [Agent Parameter Overrides](./SRS.md#USER-010)
    *   **Test Case ID:** `TC-22.1`
        *   **Test Method Signature:** `func TestAgentExecutor_ParameterOverrides(t *testing.T)`
        *   **Test Logic:** (Arrange) Use the `mock-agent.sh` which echoes its environment variables to a file. Define an agent with a default parameter. (Act) Run the Agent Executor using `--agent-param "my_param=overridden"`. (Assert) The agent's output log must show that the environment variable for the parameter was set to `overridden`, not the default value.
        *   **Required Proof of Passing:** The content of the mock agent's output log file.

---

### **3. Implementation Plan (The Execution)**

#### [x] STORY-6.1: Agent Execution Engine

1.  **Task:** Implement the agent artifact cache.
    *   **Instruction:** `In /internal/agent, create an ArtifactCache component that manages a temporary directory (e.g., in the system's temp folder). It should be responsible for creating a unique subdirectory for each agent run and providing the path for the agent's output.`
    *   **Fulfills:** **ARCH-010**.
    *   **Verification via Test Cases:** N/A (Verified via executor tests).
    *   **Documentation:**
        *   [x] **Documentation Updated:** **Instruction:** `In docs/agents/registry.md, add a subsection explaining how intermediate artifacts are handled, where they are stored temporarily, and the lifecycle of the cache.` **Evidence:** Diff of `registry.md`.
2.  **Task:** Implement the decoupled Agent Executor.
    *   **Instruction:** `Create an Executor in /internal/agent. It will take an agent definition, source/output paths, and parameter overrides. It must: 1) Prepare environment variables for parameters. 2) Use os/exec to run the agent's command. 3) Stream the agent's stdout/stderr to docloom's logger. 4) Wait for completion and handle exit codes.`
    *   **Fulfills:** **ARCH-008, USER-010**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-20.1`:**
            *   [x] **Test Method Created:** **Evidence:** Provide the test code for `TestAgentExecutor_RunCommand` and the content of `mock-agent.sh`.
            *   [x] **Test Method Passed:** **Evidence:** Console output from `go test`.
        *   **Test Case `TC-22.1`:**
            *   [x] **Test Method Created:** **Evidence:** Provide the test code for `TestAgentExecutor_ParameterOverrides`.
            *   [x] **Test Method Passed:** **Evidence:** Console output from `go test` and the captured agent log.
    *   **Documentation:**
        *   [x] **Documentation Updated:** **Instruction:** `Create docs/agents/authoring-guide.md. In this guide, explain the contract for agent developers: how they receive source/output paths as arguments and parameters as environment variables (prefixed with PARAM_). Provide a simple shell script example.` **Evidence:** Content of `authoring-guide.md`.

---
> ### **Story Completion: STORY-6.1**
> 1.  **Run Full Regression Test:**
>     *   [x] **All Prior Tests Passed:** **Instruction:** `Execute 'go test ./...'.` **Evidence:** Full summary output.
> 2.  **Create Git Commit:**
>     *   [x] **Work Committed:** **Instruction:** `git commit -m "feat(agent): implement decoupled agent executor and artifact cache"`. **Evidence:** Commit hash 809c64a.
> 3.  **Finalize Story:**
>     *   **Instruction:** Update this story's main checkbox to `[x]`.

---

#### [x] STORY-6.2: Integration with `generate` Command

1.  **Task:** Add `--agent` and `--agent-param` flags to the `generate` command.
    *   **Instruction:** `Using Cobra, add the --agent <string> and --agent-param <key=value> flags to the generate command. Ensure --agent-param can be specified multiple times.`
    *   **Fulfills:** **USER-008, USER-010**.
    *   **Verification via Test Cases:** N/A (Verified via E2E test).
    *   **Documentation:**
        *   [x] **Documentation Updated:** **Instruction:** `Update the README.md "Generating Documents" section with a new example showing a full agent invocation: docloom generate --agent <name> --source ... --agent-param "key=value".` **Evidence:** Diff of `README.md`.
2.  **Task:** Integrate agent execution into the `generate` workflow.
    *   **Instruction:** `Modify the generate command's RunE function. If the --agent flag is used: 1) Look up the agent in the registry. 2) Invoke the Agent Executor. 3) On success, replace the user's --source paths with the path to the agent's output artifacts directory. 4) Proceed with the rest of the generation process (ingestion, AI call, rendering).`
    *   **Fulfills:** **PROD-014, PROD-016**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-21.1`:**
            *   [x] **Test Method Created:** **Evidence:** Provide the E2E test code for `TestGenerateCmd_WithAgent_E2E`.
            *   [x] **Test Method Passed:** **Evidence:** Console output from the test runner.
    *   **Documentation:**
        *   [x] **Documentation Updated:** **Instruction:** `Update docs/agents/registry.md with a "Workflow" diagram that visually explains the new generate process: CLI -> Agent Executor -> Artifact Cache -> Ingester -> AI Core -> Renderer.` **Evidence:** A mermaid diagram or description added to the documentation.

---
> ### **Story Completion: STORY-6.2**
> 1.  **Run Full Regression Test:**
>     *   [x] **All Prior Tests Passed:** **Instruction:** `Execute 'go test ./...'.` **Evidence:** Full summary output.
> 2.  **Create Git Commit:**
>     *   [x] **Work Committed:** **Instruction:** `git commit -m "feat(generate): integrate agent execution into generate command"`. **Evidence:** Commit hash 7f5f79a.
> 3.  **Finalize Story:**
>     *   **Instruction:** Update this story's main checkbox to `[x]`.

---

### **4. Definition of Done**

This Phase is officially complete **only when all `STORY-6.x` checkboxes in Section 3 are marked `[x]` AND the Final Acceptance Gate below is passed.**

#### Final Acceptance Gate

*   [x] **Final Full Regression Test Passed:**
    *   **Instruction:** `Execute 'go test ./...'.`
    *   **Evidence:** All packages pass: 13 packages tested, 0 failures.

*   **Final Instruction:** Once the `Final Full Regression Test Passed` checkbox is marked `[x]`, modify the main title to `[x] PHASE-6`.