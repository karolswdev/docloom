> ### **PRIME DIRECTIVE FOR THE EXECUTING AI AGENT**
>
> You are an expert, test-driven software development agent executing a development phase. You **MUST** adhere to the following methodology without deviation:
>
> 1.  **Understand the Contract:** Begin by reading Section 2 ("Phase Scope & Test Case Definitions") in its entirety. This is your reference library for **what** to test and **how** to prove success.
> 2.  **Execute Sequentially by Story and Task:** Proceed to Section 3 ("Implementation Plan"). Address each **Story** in order. Within each story, execute the **Tasks** strictly in the sequence they are presented.
> 3.  **Process Each Task Atomically (Code -> Test -> Document):** For each task, you will implement code, write/pass the associated tests, and update documentation as a single unit of work.
> 4.  **Escalate Testing (Story & Phase Regression):**
>     a.  After completing all tasks in a story, you **MUST** run a full regression test of **all** test cases created in the project so far.
>     b.  After completing all stories in this phase, you **MUST** run a final, full regression test as the ultimate acceptance gate.
> 5.  **Commit Work:** You **MUST** create a Git commit at the completion of each story. This is a non-negotiable step.
> 6.  **Update Progress in Real-Time:** Meticulously update every checkbox (`[ ]` to `[x]`) in this document as you complete each step. Your progress tracking must be flawless.

## [ ] PHASE-5: Agent System Foundation & Discovery

---

### **1. Phase Context (What & Why)**

| ID | Title |
| :--- | :--- |
| PHASE-5 | Agent System Foundation & Discovery |

> **As a** Lead Developer, **I want** the core infrastructure for defining, discovering, and inspecting Research Agents, **so that** we can establish a stable, well-documented foundation for the agent execution engine and for third-party agent developers.

---

### **2. Phase Scope & Test Case Definitions (The Contract)**

*   **Requirement:** **PROD-015, TECH-P-008** - [Agent Definition Files](./SRS.md#PROD-015)
    *   **Test Case ID:** `TC-17.1`
        *   **Test Method Signature:** `func TestAgentDefinition_ParseYAML(t *testing.T)`
        *   **Test Logic:** (Arrange) Create a valid `test.agent.yaml` string. (Act) Unmarshal the YAML into the new Go structs. (Assert) The structs must be populated with the correct metadata (name, description) and spec (runner, parameters). Test failure on malformed YAML.
        *   **Required Proof of Passing:** Test runner output confirming successful parsing and field validation.
*   **Requirement:** **ARCH-009** - [Agent Registry & Discovery](./SRS.md#ARCH-009)
    *   **Test Case ID:** `TC-18.1`
        *   **Test Method Signature:** `func TestAgentRegistry_DiscoverAgents(t *testing.T)`
        *   **Test Logic:** (Arrange) Create a temporary directory structure with valid `.agent.yaml` files and some non-agent files. (Act) Run the registry's discovery mechanism on the temp directory. (Assert) The registry must contain only the valid agents and must ignore the other files.
        *   **Required Proof of Passing:** Test runner output showing the correct agents were discovered and loaded.
*   **Requirement:** **USER-009** - [Agent Management CLI](./SRS.md#USER-009)
    *   **Test Case ID:** `TC-19.1`
        *   **Test Method Signature:** `func TestAgentsListCmd_E2E(t *testing.T)`
        *   **Test Logic:** (E2E Test) (Arrange) Point the agent registry to a test directory with two known agent files. (Act) Execute the `docloom agents list` command. (Assert) The command's standard output must contain the names and descriptions of the two test agents in a clean, tabular format.
        *   **Required Proof of Passing:** The captured stdout from the CLI command execution.
    *   **Test Case ID:** `TC-19.2`
        *   **Test Method Signature:** `func TestAgentsDescribeCmd_E2E(t *testing.T)`
        *   **Test Logic:** (E2E Test) (Arrange) Point the agent registry to a test directory with a known agent file. (Act) Execute `docloom agents describe <agent-name>`. (Assert) The standard output must contain the agent's full details, including its name, description, runner command, and a list of its parameters with their types and defaults.
        *   **Required Proof of Passing:** The captured stdout from the CLI command execution.

---

### **3. Implementation Plan (The Execution)**

#### [ ] STORY-5.1: Agent Definition & Registry

1.  **Task:** Define agent YAML schema and Go structs.
    *   **Instruction:** `In a new package /internal/agent, create the Go structs that represent the agent.agent.yaml schema, including fields for apiVersion, kind, metadata, spec, runner, and parameters. Add yaml tags for parsing.`
    *   **Fulfills:** **PROD-015, TECH-P-008**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-17.1`:**
            *   [ ] **Test Method Created:** **Evidence:** Provide the complete code for `TestAgentDefinition_ParseYAML`.
            *   [ ] **Test Method Passed:** **Evidence:** Console output from `go test`.
    *   **Documentation:**
        *   [ ] **Documentation Updated:** **Instruction:** `Create a new file docs/agents/schema.md. Document the agent.yaml file format in detail, explaining each field's purpose and providing a complete example. This is the official contract for agent developers.` **Evidence:** The full content of `schema.md`.
2.  **Task:** Implement the Agent Registry.
    *   **Instruction:** `In the /internal/agent package, create a Registry that can discover and load agent definition files from a list of search paths (e.g., project-local .docloom/agents, user-home .docloom/agents). Store the parsed agent definitions in a map.`
    *   **Fulfills:** **ARCH-009**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-18.1`:**
            *   [ ] **Test Method Created:** **Evidence:** Provide the complete code for `TestAgentRegistry_DiscoverAgents`.
            *   [ ] **Test Method Passed:** **Evidence:** Console output from `go test`.
    *   **Documentation:**
        *   [ ] **Documentation Updated:** **Instruction:** `Create a new file docs/agents/registry.md explaining the agent discovery mechanism, search paths, and precedence. Explain how users can add their own agents.` **Evidence:** The full content of `registry.md`.

---
> ### **Story Completion: STORY-5.1**
> 1.  **Run Full Regression Test:**
>     *   [ ] **All Prior Tests Passed:** **Instruction:** `Execute 'go test ./...'.` **Evidence:** Full summary output.
> 2.  **Create Git Commit:**
>     *   [ ] **Work Committed:** **Instruction:** `git commit -m "feat(agent): implement agent definition parsing and discovery registry"`. **Evidence:** Commit hash.
> 3.  **Finalize Story:**
>     *   **Instruction:** Update this story's main checkbox to `[x]`.

---

#### [ ] STORY-5.2: CLI for Agent Management

1.  **Task:** Implement the `docloom agents` subcommand.
    *   **Instruction:** `Using Cobra, add a new top-level command docloom agents. This command will serve as the parent for list and describe.`
    *   **Fulfills:** **USER-009**.
    *   **Verification via Test Cases:** N/A (Structural, verified by subcommands).
    *   **Documentation:**
        *   [ ] **Documentation Updated:** **Instruction:** `Update README.md by adding a new "Research Agents" section. Briefly introduce the concept and add the new docloom agents command to the main usage guide.` **Evidence:** Diff of `README.md`.
2.  **Task:** Implement `agents list` and `agents describe`.
    *   **Instruction:** `Create the list and describe subcommands. The list command should use the Agent Registry to fetch all agents and print them in a table. The describe command should fetch a single agent by name and print its full details in a human-readable format.`
    *   **Fulfills:** **USER-009**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-19.1`:**
            *   [ ] **Test Method Created:** **Evidence:** Provide the complete E2E test code for the `list` command.
            *   [ ] **Test Method Passed:** **Evidence:** Console output from the test runner.
        *   **Test Case `TC-19.2`:**
            *   [ ] **Test Method Created:** **Evidence:** Provide the complete E2E test code for the `describe` command.
            *   [ ] **Test Method Passed:** **Evidence:** Console output from the test runner.
    *   **Documentation:**
        *   [ ] **Documentation Updated:** **Instruction:** `Update the new "Research Agents" section in README.md with detailed examples for both docloom agents list and docloom agents describe.` **Evidence:** Diff of `README.md`.

---
> ### **Story Completion: STORY-5.2**
> 1.  **Run Full Regression Test:**
>     *   [ ] **All Prior Tests Passed:** **Instruction:** `Execute 'go test ./...'.` **Evidence:** Full summary output.
> 2.  **Create Git Commit:**
>     *   [ ] **Work Committed:** **Instruction:** `git commit -m "feat(cli): add 'agents list' and 'agents describe' commands"`. **Evidence:** Commit hash.
> 3.  **Finalize Story:**
>     *   **Instruction:** Update this story's main checkbox to `[x]`.

---

### **4. Definition of Done**

This Phase is officially complete **only when all `STORY-5.x` checkboxes in Section 3 are marked `[x]` AND the Final Acceptance Gate below is passed.**

#### Final Acceptance Gate

*   [ ] **Final Full Regression Test Passed:**
    *   **Instruction:** `Execute 'go test ./...'.`
    *   **Evidence:** Provide the full, final summary output from the test runner.

*   **Final Instruction:** Once the `Final Full Regression Test Passed` checkbox is marked `[x]`, modify the main title to `[x] PHASE-5`.
