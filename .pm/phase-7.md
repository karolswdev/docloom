## [ ] PHASE-7: Go-Native C# Analyzer Agent with Tree-sitter

---

### **1. Phase Context (What & Why)**

| ID | Title |
| :--- | :--- |
| PHASE-7 | Go-Native C# Analyzer Agent with Tree-sitter |

> **As a** .NET Software Architect, **I want** a self-contained, Go-native research agent that can automatically analyze my C# repositories, **so that** I can rapidly generate high-quality architecture documents without needing to install the .NET SDK or any external dependencies.

---

### **2. Phase Scope & Test Case Definitions (The Contract)**

*   **Requirement:** **PROD-017** - [Initial C# Analyzer Agent](./SRS.md#PROD-017)
    *   **Test Case ID:** `TC-23.1`
        *   **Test Method Signature:** `func TestCSharpParser_ExtractAPISurface(t *testing.T)`
        *   **Test Logic:** (Go Unit Test) (Arrange) Provide a Go string containing sample C# source code with classes, methods, and XML doc comments. (Act) Run the new Go-based C# parser on this string. (Assert) The parser must correctly extract a structured representation of the public API, including class names, method signatures, and the content of the doc comments.
        *   **Required Proof of Passing:** `go test` output showing the unit test passes, proving the core parsing logic is accurate.
    *   **Test Case ID:** `TC-23.2`
        *   **Test Method Signature:** `func TestCSharpAgent_E2E_Integration(t *testing.T)`
        *   **Test Logic:** (Full E2E Integration Test) (Arrange) Create a small, complete sample C# project. Place the new, Go-based `csharp-analyzer.agent.yaml` in a discoverable path. Mock the AI Client to expect summary data from the sample project. (Act) Run `docloom generate --agent csharp-analyzer --source /path/to/sample/csharp/project ...`. (Assert) The final rendered HTML must contain specific class and method names from the sample C# code, proving the entire Go-native agent workflow is successful.
        *   **Required Proof of Passing:** The final rendered HTML output, or a test that asserts its content is correct.

---

### **3. Implementation Plan (The Execution)**

**Note:** This phase involves creating a new agent as a Go sub-package or command within our existing `docloom` repository.

#### [x] STORY-7.1: Integrating the Tree-sitter C# Parser

1.  **Task:** Add Tree-sitter and the C# grammar to the project.
    *   **Instruction:** `Add the 'smacker/go-tree-sitter' Go module. Create a new package /internal/agents/csharp/parser. Set up the necessary boilerplate to load the official Tree-sitter grammar for C#.`
    *   **Fulfills:** **PROD-017** (foundation).
    *   **Verification via Test Cases:** N/A (Structural, verified in next task).
    *   **Documentation:**
        *   [ ] **Documentation Updated:** **Instruction:** `Create a new document at docs/agents/csharp-analyzer.md. In this document, begin the "Implementation Details" section, explaining that the agent is built in Go and uses the Tree-sitter parsing library to achieve language analysis without external dependencies.` **Evidence:** The content of the new `csharp-analyzer.md` file.
2.  **Task:** Implement the core C# code parser in Go.
    *   **Instruction:** `In the /internal/agents/csharp/parser package, create a Parser struct. Implement the core analysis logic that takes C# source code as a string, parses it using Tree-sitter, and traverses the resulting syntax tree. Use Tree-sitter queries to efficiently find and extract nodes for namespaces, classes, methods, and XML doc comments.`
    *   **Fulfills:** **PROD-017**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-23.1`:**
            *   [ ] **Test Method Created:** **Evidence:** Provide the Go unit test code for `TestCSharpParser_ExtractAPISurface`.
            *   [ ] **Test Method Passed:** **Evidence:** Console output from `go test`.
    *   **Documentation:**
        *   [ ] **Documentation Updated:** **Instruction:** `In docs/agents/authoring-guide.md, add a new section titled "Building a Go-Native Analyzer with Tree-sitter". Provide a high-level overview and code snippets showing how to load a grammar and query a syntax tree. This establishes a world-class pattern for future Go-based agents.` **Evidence:** Diff of `authoring-guide.md`.

---
> ### **Story Completion: STORY-7.1**
> 1.  **Run Full Regression Test:**
>     *   [ ] **All Prior Tests Passed:** **Instruction:** `Execute 'go test ./...'.` **Evidence:** Full summary output.
> 2.  **Create Git Commit:**
>     *   [ ] **Work Committed:** **Instruction:** `git commit -m "feat(agent-csharp): implement go-native c# parser with tree-sitter"`. **Evidence:** Commit hash.
> 3.  **Finalize Story:**
>     *   **Instruction:** Update this story's main checkbox to `[x]`.

---

#### [ ] STORY-7.2: Building the Executable Agent and Final Integration

1.  **Task:** Implement the agent's main logic and artifact generation.
    *   **Instruction:** `Create a new command in /cmd/docloom-agent-csharp/main.go. This command will be the agent's entry point. It will parse the source/output path arguments and environment variables for parameters. It will then use the parser from Story 7.1 to analyze all .cs files in the source directory and write the summary artifacts (e.g., ProjectSummary.md, ApiSurface.md) to the output path.`
    *   **Fulfills:** **PROD-017**.
    *   **Verification via Test Cases:** N/A (Verified via E2E test).
    *   **Documentation:**
        *   [ ] **Documentation Updated:** **Instruction:** `In docs/agents/csharp-analyzer.md, add a comprehensive "Output Artifacts" section describing each markdown file the agent produces and the information contained within.` **Evidence:** Diff of `csharp-analyzer.md`.
2.  **Task:** Create the agent definition file and update the build process.
    *   **Instruction:** `Create the csharp-analyzer.agent.yaml file in a new /agents directory. The 'command' in the runner spec will be ["./docloom-agent-csharp"]. Update the project's Makefile and GoReleaser configuration to build and package this new agent binary alongside the main docloom binary.`
    *   **Fulfills:** **PROD-017**.
    *   **Verification via Test Cases:** N/A (Structural, verified by next task).
    *   **Documentation:**
        *   [ ] **Documentation Updated:** **Instruction:** `In docs/agents/csharp-analyzer.md, add a "Configuration" section showing the complete csharp-analyzer.agent.yaml file and explaining its parameters.` **Evidence:** Diff of `csharp-analyzer.md`.
3.  **Task:** Write and pass the final end-to-end integration test.
    *   **Instruction:** `With the Go-native agent building and the YAML definition in place, implement the final E2E test. The test must set up a sample C# project, invoke docloom with the agent, and assert that the final rendered document contains the correctly analyzed information.`
    *   **Fulfills:** **PROD-017**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-23.2`:**
            *   [ ] **Test Method Created:** **Evidence:** Provide the E2E Go test code `TestCSharpAgent_E2E_Integration`.
            *   [ ] **Test Method Passed:** **Evidence:** Console output from `go test`.
    *   **Documentation:**
        *   [ ] **Documentation Updated:** **Instruction:** `Update the main README.md to proudly announce the C# Analyzer agent. Highlight that it is fully self-contained and requires no external dependencies, and link to its detailed documentation page.` **Evidence:** Diff of the main `README.md`.

---
> ### **Story Completion: STORY-7.2**
> 1.  **Run Full Regression Test:**
>     *   [ ] **All Prior Tests Passed:** **Instruction:** `Execute 'go test ./...'.` **Evidence:** Summary output from the Go test runner.
> 2.  **Create Git Commit:**
>     *   [ ] **Work Committed:** **Instruction:** `git commit -m "feat(agent-csharp): complete go-native agent and e2e integration"`. **Evidence:** Commit hash.
> 3.  **Finalize Story:**
>     *   **Instruction:** Update this story's main checkbox to `[x]`.

---

### **4. Definition of Done**

This Phase is officially complete **only when all `STORY-7.x` checkboxes in Section 3 are marked `[x]` AND the Final Acceptance Gate below is passed.**

#### Final Acceptance Gate

*   [ ] **Final Full Regression Test Passed:**
    *   **Instruction:** `Execute 'go test ./...'.`
    *   **Evidence:** Provide the full, final summary output from the test runner.

*   **Final Instruction:** Once the `Final Full Regression Test Passed` checkbox is marked `[x]`, modify the main title to `[x] PHASE-7`.