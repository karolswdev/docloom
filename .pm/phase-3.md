## [x] PHASE-3: Advanced Input Processing & User Experience

---

### **1. Phase Context (What & Why)**

| ID | Title |
| :--- | :--- |
| PHASE-3 | Advanced Input Processing & User Experience |

> **As a** Power User, **I want** the system to intelligently select content from large and diverse source files like PDFs, and I want tools to inspect the system's process before incurring costs, **so that** I can get more relevant results and debug my inputs effectively.

---

### **2. Phase Scope & Test Case Definitions (The Contract)**

*   **Requirement:** **PROD-005, TECH-P-004, TECH-L-006** - [PDF Text Extraction](./SRS.md#PROD-005)
    *   **Test Case ID:** `TC-10.1`
        *   **Test Method Signature:** `func TestIngester_ExtractPDFText(t *testing.T)`
        *   **Test Logic:** (Requires `pdftotext` in the test environment). (Arrange) Provide a sample PDF file with known text. (Act) Call the ingestion logic on this file. (Assert) The returned string must contain the expected text from the PDF.
        *   **Required Proof of Passing:** Test runner output.
*   **Requirement:** **PROD-006** - [Chunking & Selection](./SRS.md#PROD-006)
    *   **Test Case ID:** `TC-11.1`
        *   **Test Method Signature:** `func TestChunker_SimpleChunking(t *testing.T)`
        *   **Test Logic:** (Arrange) Provide a long text string that exceeds a small, predefined token limit. (Act) Run the chunking and selection logic. (Assert) The returned string must be shorter than the original and within the token limit. A simple heuristic might just truncate the text.
        *   **Required Proof of Passing:** Test runner output showing the length assertion passed.
*   **Requirement:** **PROD-010, USER-004, ARCH-007** - [Dry-Run & Explainability](./SRS.md#PROD-010)
    *   **Test Case ID:** `TC-12.1`
        *   **Test Method Signature:** `func TestGenerateCmd_DryRun(t *testing.T)`
        *   **Test Logic:** (E2E Test) (Arrange) Mock the AI Client to fail the test if it is ever called. (Act) Run the `generate` command with the `--dry-run` flag. (Assert) The command should exit successfully (code 0) and print the assembled prompt and selected source chunks to its standard output. The AI client must not have been called.
        *   **Required Proof of Passing:** Console output from the CLI showing the prompt, and test runner confirming success.
*   **Requirement:** **USER-003, USER-005** - [Verbose Logging & Safe Writes](./SRS.md#USER-003)
    *   **Test Case ID:** `TC-13.1`
        *   **Test Method Signature:** `func TestGenerateCmd_VerboseLogging(t *testing.T)`
        *   **Test Logic:** (E2E Test) (Arrange) Run the CLI against a mock AI server. Capture stdout/stderr. (Act) Run `generate --verbose`. (Assert) The captured output must contain detailed debug messages (e.g., "Ingesting file X", "Selected N chunks", "Model call successful").
        *   **Required Proof of Passing:** Captured log output.
    *   **Test Case ID:** `TC-13.2`
        *   **Test Method Signature:** `func TestGenerateCmd_SafeWrites(t *testing.T)`
        *   **Test Logic:** (E2E Test) (Arrange) Create a dummy output file (e.g., `output.html`). (Act) Run the `generate` command pointing to that output file. (Assert) The command MUST fail with a non-zero exit code and an error message about the file existing. (Act 2) Run again with `--force`. (Assert 2) The command MUST succeed.
        *   **Required Proof of Passing:** Console output from both CLI runs showing the expected behavior.

---

### **3. Implementation Plan (The Execution)**

#### [x] STORY-3.1: Enhanced Source Ingestion and Processing

1.  **Task:** Implement PDF text extraction.
    *   **Instruction:** `Update the /internal/ingest package to handle .pdf files. Implement a strategy that uses the 'pdftotext' command-line tool when available. The implementation must gracefully handle the case where the tool is not in the system's PATH.`
    *   **Fulfills:** **PROD-005, TECH-P-004, TECH-L-006**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-10.1`:**
            *   [x] **Test Method Created:** **Evidence:** Code for `TestIngester_ExtractPDFText`.
            *   [x] **Test Method Passed:** **Evidence:** `go test` output. See: ./evidence/PHASE-3/story-3.1/task-1/test-output/TC-10.1.log
    *   **Documentation:**
        *   [x] **Documentation Updated:** **Instruction:** `Update README to note the dependency on 'poppler-utils' for PDF processing.` **Evidence:** Diff of `README.md`. See: ./evidence/PHASE-3/story-3.1/task-1/README_pdf_docs.diff
2.  **Task:** Implement heuristic content chunking.
    *   **Instruction:** `Create a new /internal/chunk package. Implement a simple heuristic chunking strategy that takes the ingested text, estimates token count, and truncates it to fit within a configured context limit (e.g., max_tokens).`
    *   **Fulfills:** **PROD-006, ARCH-001**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-11.1`:**
            *   [x] **Test Method Created:** **Evidence:** Code for `TestChunker_SimpleChunking`.
            *   [x] **Test Method Passed:** **Evidence:** `go test` output. See: ./evidence/PHASE-3/story-3.1/task-2/test-output/TC-11.1.log
    *   **Documentation:** [x] No documentation updates required for this task.

---
> ### **Story Completion: STORY-3.1**
>
> 1.  **Run Full Regression Test:**
>     *   [x] **All Prior Tests Passed:** **Instruction:** `Execute 'go test ./...'.` **Evidence:** Full summary output. See: ./evidence/PHASE-3/story-3.1/regression-test.log
> 2.  **Create Git Commit:**
>     *   [x] **Work Committed:** **Instruction:** `Execute 'git commit -m "feat(ingest): add pdf extraction and content chunking"'.` **Evidence:** Commit hash: 3ce4e6d. See: ./evidence/PHASE-3/story-3.1/commit.txt
> 3.  **Finalize Story:**
>     *   **Instruction:** Update this story's main checkbox to `[x]`.

---

#### [ ] STORY-3.2: User Experience and Transparency Features

1.  **Task:** Implement `--dry-run` mode.
    *   **Instruction:** `Modify the 'generate' command logic. If the --dry-run flag is present, the process should stop before calling the AI client. It must print the final assembled prompt, the list of selected source chunks, and the target JSON schema to the console.`
    *   **Fulfills:** **PROD-010, USER-004, ARCH-007**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-12.1`:**
            *   [x] **Test Method Created:** **Evidence:** Code for `TestGenerateCmd_DryRun`.
            *   [x] **Test Method Passed:** **Evidence:** `go test` output. See: ./evidence/PHASE-3/story-3.2/task-1/test-output/TC-12.1.log
    *   **Documentation:**
        *   [x] **Documentation Updated:** **Instruction:** `Add '--dry-run' to the usage examples in the README, explaining its purpose.` **Evidence:** Diff of `README.md`. See: ./evidence/PHASE-3/story-3.2/README_flags_docs.diff
2.  **Task:** Implement `--verbose` logging and safe file writes.
    *   **Instruction:** `Throughout the generation pipeline, add detailed debug-level logs that are only shown when --verbose is active. Implement a check before writing output files; if a file exists, exit with an error unless the --force flag is provided.`
    *   **Fulfills:** **USER-003, USER-005**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-13.1`:**
            *   [x] **Test Method Created:** **Evidence:** Code for `TestGenerateCmd_VerboseLogging`.
            *   [x] **Test Method Passed:** **Evidence:** `go test` output. See: ./evidence/PHASE-3/story-3.2/task-2/test-output/TC-13.1.log
        *   **Test Case `TC-13.2`:**
            *   [x] **Test Method Created:** **Evidence:** Code for `TestGenerateCmd_SafeWrites`.
            *   [x] **Test Method Passed:** **Evidence:** `go test` output. See: ./evidence/PHASE-3/story-3.2/task-2/test-output/TC-13.2.log
    *   **Documentation:**
        *   [x] **Documentation Updated:** **Instruction:** `Add '--force' and '--verbose' flags to the README usage guide.` **Evidence:** Diff of `README.md`. See: ./evidence/PHASE-3/story-3.2/README_flags_docs.diff

---
> ### **Story Completion: STORY-3.2**
>
> 1.  **Run Full Regression Test:**
>     *   [x] **All Prior Tests Passed:** **Instruction:** `Execute 'go test ./...'.` **Evidence:** Full summary output. See: ./evidence/PHASE-3/story-3.2/regression-test.log
> 2.  **Create Git Commit:**
>     *   [x] **Work Committed:** **Instruction:** `Execute 'git commit -m "feat(ux): implement dry-run, verbose logging, and safe writes"'.` **Evidence:** Commit hash: 23c82c6. See: ./evidence/PHASE-3/story-3.2/commit.txt
> 3.  **Finalize Story:**
>     *   **Instruction:** Update this story's main checkbox to `[x]`.

---

### **4. Definition of Done**

This Phase is officially complete **only when all `STORY-3.x` checkboxes in Section 3 are marked `[x]` AND the Final Acceptance Gate below is passed.**

#### Final Acceptance Gate

*   [x] **Final Full Regression Test Passed:**
    *   **Instruction:** `Execute 'go test ./...'.`
    *   **Evidence:** Provide the full, final summary output from the test runner. See: ./evidence/PHASE-3/final-acceptance-gate.log

*   **Final Instruction:** Once the `Final Full Regression Test Passed` checkbox is marked `[x]`, modify the main title to `[x] PHASE-3`.
