## [ ] PHASE-9: `csharp-cc-cli` Agent for Claude Code & Mock Integration

---

### **1. Phase Context (What & Why)**

| ID | Title |
| :--- | :--- |
| PHASE-9 | `csharp-cc-cli` Agent for Claude Code & Mock Integration |

> **As a** Lead Architect, **I want** to define the data contract for a high-level code analysis agent and integrate a mock version of the **Claude Code CLI (`cc-cli`)**, **so that** we can prove the entire `docloom`-to-agent workflow is sound before investing in the development of the complex external tool.

---

### **2. Phase Scope & Test Case Definitions (The Contract)**

*   **Requirement:** **PROD-018 (Partial)** - External CLI Tool Agent (Contract & Mock)
    *   **Test Case ID:** `TC-27.1`
        *   **Test Method Signature:** N/A (Documentation & Schema Test)
        *   **Test Logic:** (Arrange) Create a formal "Artifact Specification" document that defines the directory structure and file formats the **Claude Code CLI (`cc-cli`)** is expected to produce. (Act) Create a `csharp-cc-cli.agent.yaml` file that references a mock script. (Assert) The YAML file must be parsable, and the artifact spec document must be clear.
        *   **Required Proof of Passing:** The final, checked-in content of `csharp-cc-cli.agent.yaml` and `docs/agents/artifact-spec-claude-code.md`.
    *   **Test Case ID:** `TC-27.2`
        *   **Test Method Signature:** `func TestCSharpCCAgent_E2E_WithMock(t *testing.T)`
        *   **Test Logic:** (E2E Test) (Arrange) Create a `mock-cc-cli.sh` script that simulates the **Claude Code CLI** by creating the exact directory and file structure defined in the Artifact Specification. (Act) Run `docloom generate --agent csharp-cc-cli ...`. (Assert) `docloom` must successfully invoke the mock script and ingest all generated mock artifacts into the AI prompt.
        *   **Required Proof of Passing:** The captured prompt sent to the (mocked) AI client must contain the placeholder text from the mock artifacts.

---

### **3. Implementation Plan (The Execution)**

#### [x] STORY-9.1: Defining the Claude Code Agent Contract

1.  **Task:** Design and document the **Claude Code CLI (`cc-cli`)** Artifact Specification.
    *   **Instruction:** `Create a new document at docs/agents/artifact-spec-claude-code.md. This is the official contract. Define the expected output structure of the Claude Code CLI, including filenames (e.g., overview.md, dependencies.json) and their content schemas.`
    *   **Fulfills:** `PROD-018` (Contract Definition).
    *   **Verification via Test Cases:** N/A (Documentation-as-code).
    *   **Documentation:**
        *   [ ] **Documentation Updated:** **Instruction:** `This task IS the documentation. The artifact specification is the single source of truth that will govern the development of the real Claude Code CLI.` **Evidence:** The complete content of `artifact-spec-claude-code.md`.
2.  **Task:** Create the `csharp-cc-cli.agent.yaml` definition file.
    *   **Instruction:** `In the /agents directory, create csharp-cc-cli.agent.yaml. Define the agent's metadata. For the 'command', point it to a mock script: ["./mock-cc-cli.sh"]. Add parameters that will eventually be passed to the real Claude Code CLI.`
    *   **Fulfills:** `PROD-018`.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-27.1`:**
            *   [ ] **Test Method Created:** This test is fulfilled by creating the agent file itself.
            *   [ ] **Test Method Passed:** **Evidence:** The content of `csharp-cc-cli.agent.yaml` and the output of `docloom agents describe csharp-cc-cli` running successfully.
    *   **Documentation:**
        *   [ ] **Documentation Updated:** **Instruction:** `Create docs/agents/csharp-cc-cli.md. Document this new agent, explaining that it orchestrates the powerful external **Claude Code CLI (`cc-cli`)**. Link to the Artifact Specification document.` **Evidence:** The new `csharp-cc-cli.md` file.

---
> ### **Story Completion: STORY-9.1**
> 1.  **Run Full Regression Test:**
>     *   [ ] **All Prior Tests Passed:** **Instruction:** `Execute 'go test ./...'.` **Evidence:** Full summary output.
> 2.  **Create Git Commit:**
>     *   [ ] **Work Committed:** **Instruction:** `git commit -m "feat(agent-claude): define artifact spec and agent for Claude Code CLI"`. **Evidence:** Commit hash.
> 3.  **Finalize Story:**
>     *   **Instruction:** Update this story's main checkbox to `[x]`.

---

#### [ ] STORY-9.2: Mock Implementation and E2E Testing

1.  **Task:** Create the mock `cc-cli` script.
    *   **Instruction:** `Create a shell script named mock-cc-cli.sh to simulate the behavior of the future Claude Code CLI. The script MUST accept source and output path arguments and create the directory and file structure exactly as defined in the Artifact Specification, filling them with placeholder text.`
    *   **Fulfills:** `PROD-018` (Mocking).
    *   **Verification via Test Cases:** N/A (Script is part of the test setup).
    *   **Documentation:** [ ] No documentation updates required for this task.
2.  **Task:** Implement the end-to-end integration test with the mock script.
    *   **Instruction:** `In Go, create the E2E test for the csharp-cc-cli agent. The test will run docloom generate with the agent, verify that the mock-cc-cli.sh was called correctly, and assert that the content passed to the AI prompt includes the placeholder text generated by the mock script.`
    *   **Fulfills:** `PROD-018`.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-27.2`:**
            *   [ ] **Test Method Created:** **Evidence:** Provide the Go test code for `TestCSharpCCAgent_E2E_WithMock`.
            *   [ ] **Test Method Passed:** **Evidence:** Console output from `go test`.
    *   **Documentation:**
        *   [ ] **Documentation Updated:** **Instruction:** `In docs/agents/authoring-guide.md, add a section on "Testing Agents with Mock Implementations". Explain this best practice and use the mock-cc-cli.sh script as the primary example.` **Evidence:** Diff of `authoring-guide.md`.

---
> ### **Story Completion: STORY-9.2**
> 1.  **Run Full Regression Test:**
>     *   [ ] **All Prior Tests Passed:** **Instruction:** `Execute 'go test ./...'.` **Evidence:** Full summary output.
> 2.  **Create Git Commit:**
>     *   [ ] **Work Committed:** **Instruction:** `git commit -m "test(agent-claude): add mock Claude Code CLI and e2e integration test"`. **Evidence:** Commit hash.
> 3.  **Finalize Story:**
>     *   **Instruction:** Update this story's main checkbox to `[x]`.

---

### **4. Definition of Done**

This Phase is officially complete **only when all `STORY-9.x` checkboxes in Section 3 are marked `[x]` AND the Final Acceptance Gate below is passed.**

#### Final Acceptance Gate

*   [ ] **Final Full Regression Test Passed:**
    *   **Instruction:** `Execute 'go test ./...'.`
    *   **Evidence:** Provide the full, final summary output from the test runner.

*   **Final Instruction:** Once the `Final Full Regression Test Passed` checkbox is marked `[x]`, modify the main title to `[x] PHASE-9`.