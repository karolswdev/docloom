## [ ] PHASE-10: Building the Claude Code CLI (`cc-cli`) Analysis Tool

---

### **1. Phase Context (What & Why)**

| ID | Title |
| :--- | :--- |
| PHASE-10 | Building the Claude Code CLI (`cc-cli`) Analysis Tool |

> **As a** `docloom` User, **I want** a powerful, standalone **Claude Code CLI (`cc-cli`)** that uses the Claude LLM to perform deep analysis of my C# repositories, **so that** I can provide my templates with rich, accurate, and synthesized context for fully automated document generation.

---

### **2. Phase Scope & Test Case Definitions (The Contract)**

*   **Requirement:** **PROD-018** - External CLI Tool Agent (Real Implementation)
    *   **Test Case ID:** `TC-28.1`
        *   **Test Method Signature:** N/A (Tool Unit Tests)
        *   **Test Logic:** The **Claude Code CLI (`cc-cli`)** project will have its own suite of unit tests for its internal logic (file discovery, prompt generation, Claude API response parsing).
        *   **Required Proof of Passing:** `go test` output from the `cc-cli`'s own module showing its tests pass.
    *   **Test Case ID:** `TC-28.2`
        *   **Test Method Signature:** N/A (Tool Acceptance Test)
        *   **Test Logic:** Run the compiled `cc-cli` binary against a sample C# repository. The tool must produce the full set of artifacts, conforming exactly to the Artifact Specification defined in Phase 9.
        *   **Required Proof of Passing:** A directory listing and content sample of the generated artifacts.
    *   **Test Case ID:** `TC-28.3`
        *   **Test Method Signature:** `func TestCSharpCCAgent_E2E_WithRealTool(t *testing.T)`
        *   **Test Logic:** This test re-runs `TC-27.2`, but the `csharp-cc-cli.agent.yaml` is updated to point to the newly compiled, *real* **Claude Code CLI (`cc-cli`)** binary instead of the mock script.
        *   **Required Proof of Passing:** `go test` output confirming the entire, real end-to-end workflow succeeds.

---

### **3. Implementation Plan (The Execution)**

**Note:** This phase involves creating a new, standalone Go CLI project in `/tools/claude-code-cli`.

#### [ ] STORY-10.1: `cc-cli` Scaffolding and Core Logic

1.  **Task:** Scaffold the **Claude Code CLI (`cc-cli`)** Go project.
    *   **Instruction:** `In a new directory /tools/claude-code-cli, initialize a new Go module. Use Cobra to create a CLI application that accepts arguments (--repo-path, --output-path) and optional parameters.`
    *   **Fulfills:** `PROD-018` (foundation).
    *   **Verification via Test Cases:** N/A (Structural).
    *   **Documentation:**
        *   [ ] **Documentation Updated:** **Instruction:** `Create a comprehensive README.md within the /tools/claude-code-cli directory. This is the official user manual for the standalone Claude Code CLI. Document its purpose, installation, and usage.` **Evidence:** The new `README.md` for `cc-cli`.
2.  **Task:** Implement repository scanning and prompt generation.
    *   **Instruction:** `Implement the core logic for the Claude Code CLI. This includes scanning the target repository for key files (READMEs, .sln, etc.), reading their content, and assembling a sophisticated initial prompt for the Claude LLM to perform high-level analysis.`
    *   **Fulfills:** `PROD-018`.
    *   **Verification via Test Cases:** (Verified by `TC-28.1`).
    *   **Documentation:** [ ] No documentation updates required for this task.

---
> ### **Story Completion: STORY-10.1**
> 1.  **Run Full Regression Test:**
>     *   [ ] **All Prior Tests Passed:** **Instruction:** `Execute 'go test ./...' within the cc-cli module.` **Evidence:** Test runner output.
> 2.  **Create Git Commit:**
>     *   [ ] **Work Committed:** **Instruction:** `git commit -m "feat(cc-cli): scaffold claude code cli and implement core scanning"`. **Evidence:** Commit hash.
> 3.  **Finalize Story:**
>     *   **Instruction:** Update this story's main checkbox to `[x]`.

---

#### [ ] STORY-10.2: Claude LLM Interaction and Artifact Generation

1.  **Task:** Implement the Claude API client and interaction logic.
    *   **Instruction:** `Add an LLM client to the cc-cli tool to communicate with the Claude API. Implement the logic to send the initial analysis prompt and parse the structured response from the LLM.`
    *   **Fulfills:** `PROD-018`.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-28.1`:**
            *   [ ] **Test Method Created:** **Evidence:** Provide the Go unit tests for the Claude API interaction and parsing logic.
            *   [ ] **Test Method Passed:** **Evidence:** Console output from `go test`.
    *   **Documentation:** [ ] No documentation updates required for this task.
2.  **Task:** Implement artifact writer conforming to the specification.
    *   **Instruction:** `Create a component that takes the structured analysis from the Claude LLM and writes it to the disk, creating the exact directory and file structure defined in the Artifact Specification from Phase 9.`
    *   **Fulfills:** `PROD-018`.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-28.2`:**
            *   [ ] **Test Method Created:** **Evidence:** Provide the acceptance test script or Go test code.
            *   [ ] **Test Method Passed:** **Evidence:** The generated artifacts from the test run.
    *   **Documentation:**
        *   [ ] **Documentation Updated:** **Instruction:** `Update the cc-cli README.md with an "Output" section that links to the main docloom Artifact Specification document, ensuring a single source of truth.` **Evidence:** Diff of `cc-cli`'s `README.md`.
3.  **Task:** Final integration and release.
    *   **Instruction:** `Update the build process (Makefile, GoReleaser) to build and package the cc-cli binary. Update the csharp-cc-cli.agent.yaml to point to the real cc-cli command. Run the final E2E test to verify the entire system works with the real tool.`
    *   **Fulfills:** `PROD-018`.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-28.3`:**
            *   [ ] **Test Method Created:** **Evidence:** The E2E test is already written; this task is to run it against the real binary.
            *   [ ] **Test Method Passed:** **Evidence:** Console output from `go test` in the main `docloom` repo.
    *   **Documentation:**
        *   [ ] **Documentation Updated:** **Instruction:** `In the main docloom README.md, update the section on the csharp-cc-cli agent. Add instructions for how to install the required **Claude Code CLI (`cc-cli`)** binary, or note that it's bundled with the docloom release.` **Evidence:** Diff of the main `README.md`.

---
> ### **Story Completion: STORY-10.2**
> 1.  **Run Full Regression Test:**
>     *   [ ] **All Prior Tests Passed:** **Instruction:** `Execute 'go test ./...' in both the main docloom repo and the cc-cli repo.` **Evidence:** Summary outputs from both test runners.
> 2.  **Create Git Commit:**
>     *   [ ] **Work Committed:** **Instruction:** `git commit -m "feat(cc-cli): complete claude code cli and final integration"`. **Evidence:** Commit hash.
> 3.  **Finalize Story:**
>     *   **Instruction:** Update this story's main checkbox to `[x]`.

---

### **4. Definition of Done**

This Phase is officially complete **only when all `STORY-10.x` checkboxes in Section 3 are marked `[x]` AND the Final Acceptance Gate below is passed.**

#### Final Acceptance Gate

*   [ ] **Final Full Regression Test Passed:**
    *   **Instruction:** `Execute 'go test ./...'.`
    *   **Evidence:** Provide the full, final summary output from the test runner.

*   **Final Instruction:** Once the `Final Full Regression Test Passed` checkbox is marked `[x]`, modify the main title to `[x] PHASE-10`.