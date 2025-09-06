### **Phase 1: Project Scaffolding & Deterministic Rendering**
**Focus:** Establish the foundational project structure, CLI, configuration, and a deterministic rendering pipeline. This phase deliberately omits the AI component to build and test the static parts of the system in isolation.

### **Phase 2: Core AI Integration & Validation**
**Focus:** Introduce the "brain" of the system. This phase integrates the AI client, prompt generation, and the critical JSON schema validation and repair loop, ensuring the system can reliably produce structured data from an AI model.

### **Phase 3: Advanced Input Processing & User Experience**
**Focus:** Enhance the quality and relevance of the AI's input. This involves implementing sophisticated source ingestion (PDFs), content chunking, and providing user-facing tools for transparency and debugging, like `--dry-run`.

### **Phase 4: Productionization, Distribution & Polish**
**Focus:** Prepare the application for real-world use. This phase centers on creating robust CI/CD pipelines for multi-platform binary and container releases, implementing final features, and ensuring operational readiness.

***

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

## [ ] PHASE-1: Project Scaffolding & Deterministic Rendering

---

### **1. Phase Context (What & Why)**

| ID | Title |
| :--- | :--- |
| PHASE-1 | Project Scaffolding & Deterministic Rendering |

> **As a** Lead Developer, **I want** a runnable CLI skeleton with configuration loading and a deterministic HTML rendering pipeline, **so that** I can establish the core project structure and test the output generation workflow with mock data before introducing AI dependencies.

---

### **2. Phase Scope & Test Case Definitions (The Contract)**

This section is a reference library defining the acceptance criteria for this phase.

*   **Requirement:** **TECH-P-001, TECH-P-002, TECH-L-001, TECH-L-004, USER-002** - [CLI Foundation](./SRS.md#TECH-P-001)
    *   **Test Case ID:** `TC-1.1`
        *   **Test Method Signature:** `func TestRootCmd_HelpFlag(t *testing.T)`
        *   **Test Logic:** Execute `docloom --help`. Assert that the output contains usage information for the `generate` command and exits with code 0.
        *   **Required Proof of Passing:** Console output showing the help text and successful test run.
*   **Requirement:** **PROD-012, USER-006** - [Configurable Defaults](./SRS.md#PROD-012)
    *   **Test Case ID:** `TC-2.1`
        *   **Test Method Signature:** `func TestConfig_LoadWithPrecedence(t *testing.T)`
        *   **Test Logic:** (Arrange) Create a config file, set an environment variable, and provide a CLI flag for the same value (e.g., `model`). (Act) Load the configuration. (Assert) The value from the CLI flag MUST be the final, effective value. Repeat for env over file.
        *   **Required Proof of Passing:** Unit test output confirming the correct precedence is applied in all cases.
*   **Requirement:** **PROD-002, ARCH-006** - [Template Registry](./SRS.md#PROD-002)
    *   **Test Case ID:** `TC-3.1`
        *   **Test Method Signature:** `func TestTemplateRegistry_Load(t *testing.T)`
        *   **Test Logic:** (Arrange) Point the registry to a directory containing valid template assets (`architecture-vision.html`, `schema.json`, `prompt.txt`). (Act) Load the templates. (Assert) The registry should correctly load and contain the `architecture-vision` template with its assets populated.
        *   **Required Proof of Passing:** Test runner output showing the test passed. Debug logs may show the loaded template data.
*   **Requirement:** **PROD-003, PROD-004, ARCH-004** - [Idempotent Rendering](./SRS.md#PROD-003)
    *   **Test Case ID:** `TC-4.1`
        *   **Test Method Signature:** `func TestRenderer_RenderHTML_Golden(t *testing.T)`
        *   **Test Logic:** This will be a golden file test. (Arrange) Provide a `testdata/architecture-vision.html` template and a `testdata/fields.json` input file. (Act) Run the render function. (Assert) The generated HTML output MUST exactly match a pre-approved `testdata/expected.html` golden file. The `fields.json` must also be saved to the output directory.
        *   **Required Proof of Passing:** Test runner output confirming the test passed. The golden file mechanism will prove correctness.
*   **Requirement:** **DEV-001, DEV-004** - [Containerized Workflow & CI](./SRS.md#DEV-001)
    *   **Test Case ID:** `TC-5.1`
        *   **Test Method Signature:** N/A (CI Workflow Test)
        *   **Test Logic:** A GitHub Actions workflow is created. On push, it builds the Docker image, then runs `go build ./...` and `go test ./...` inside the container.
        *   **Required Proof of Passing:** A screenshot or link to a successful GitHub Actions run for the `ci.yml` workflow.

---

### **3. Implementation Plan (The Execution)**

#### [x] STORY-1.1: Project Setup & CLI Foundation

1.  **Task:** Initialize Go module and project structure.
    *   **Instruction:** `Execute 'go mod init github.com/karolswdev/docloom'. Create the modular package layout as specified in ARCH-001: /cmd, /internal/config, /internal/cli, /internal/render, /internal/templates.`
    *   **Fulfills:** **ARCH-001, TECH-P-001, TECH-P-006**.
    *   **Verification via Test Cases:** N/A (Structural setup).
    *   **Documentation:**
        *   [x] **Documentation Updated:** **Instruction:** `Create a README.md with the project title and a brief description.` **Evidence:** `README.md` content.
2.  **Task:** Implement the basic CLI structure with Cobra.
    *   **Instruction:** `Use spf13/cobra to create the root command and a 'generate' subcommand in the /internal/cli package. Add a persistent --verbose flag. Implement a basic --help flag handler.`
    *   **Fulfills:** **USER-001, USER-002, TECH-P-002, TECH-L-001**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-1.1`:**
            *   [x] **Test Method Created:** **Evidence:** Provided in `/internal/cli/root_test.go`.
            *   [x] **Test Method Passed:** **Evidence:** `evidence/PHASE-1/story-1.1/task-2/test-output/TC-1.1.log` showing test passed.
    *   **Documentation:**
        *   [x] **Documentation Updated:** **Instruction:** `Update README.md with a "Usage" section showing the basic 'docloom generate --help' command.` **Evidence:** Added Usage section to README.md.
3.  **Task:** Integrate structured logging with Zerolog.
    *   **Instruction:** `Add rs/zerolog to the project. Configure a global logger that respects the --verbose flag to switch between INFO and DEBUG levels. Ensure logs are human-readable by default.`
    *   **Fulfills:** **NFR-005 (partial), TECH-L-004**.
    *   **Verification via Test Cases:** N/A (Manual verification during other tests).
    *   **Documentation:** [x] No documentation updates required for this task.

---
> ### **Story Completion: STORY-1.1**
>
> You may only proceed once all checkboxes for all tasks within this story are marked `[x]`. Then, you **MUST** complete the following steps in order:
>
> 1.  **Run Full Regression Test:**
>     *   [x] **All Prior Tests Passed:** Checked after running all tests created in the project up to this point.
>     *   **Instruction:** `Execute 'go test ./...'.`
>     *   **Evidence:** `evidence/PHASE-1/story-1.1/regression-test.log`
> 2.  **Create Git Commit:**
>     *   [x] **Work Committed:** Checked after creating the Git commit.
>     *   **Instruction:** `Execute 'git add .' followed by 'git commit -m "feat(cli): setup project structure and basic CLI with cobra"'.`
>     *   **Evidence:** Commit hash: 854a91c
> 3.  **Finalize Story:**
>     *   **Instruction:** Once the two checkboxes above are complete, you **MUST** update this story's main checkbox from `[ ]` to `[x]`.

---

#### [x] STORY-1.2: Configuration Loading & Template Management

1.  **Task:** Implement configuration loading logic.
    *   **Instruction:** `In the /internal/config package, implement logic to load settings from a 'docloom.yaml' file, environment variables (e.g., DOCLOOM_MODEL), and CLI flags. Ensure the precedence is CLI > ENV > FILE > DEFAULTS.`
    *   **Fulfills:** **PROD-012, USER-006**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-2.1`:**
            *   [x] **Test Method Created:** **Evidence:** Provided in `/internal/config/config_test.go`.
            *   [x] **Test Method Passed:** **Evidence:** `evidence/PHASE-1/story-1.2/task-1/test-output/TC-2.1.log` showing test passed.
    *   **Documentation:**
        *   [x] **Documentation Updated:** **Instruction:** `Add a 'Configuration' section to the README explaining the precedence order and showing a sample docloom.yaml.` **Evidence:** Added Configuration section with precedence and sample YAML to README.md.
2.  **Task:** Implement the Template Registry.
    *   **Instruction:** `Create a template registry in /internal/templates that can discover and load templates from a given directory. Each template consists of an HTML file, a JSON schema, and a prompt file. Embed the default templates into the binary using 'embed'.`
    *   **Fulfills:** **PROD-002, ARCH-006**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-3.1`:**
            *   [x] **Test Method Created:** **Evidence:** Provided in `/internal/templates/registry_test.go`.
            *   [x] **Test Method Passed:** **Evidence:** `evidence/PHASE-1/story-1.2/task-2/test-output/TC-3.1.log` showing test passed.
    *   **Documentation:** [x] No documentation updates required for this task.

---
> ### **Story Completion: STORY-1.2**
>
> You may only proceed once all checkboxes for all tasks within this story are marked `[x]`. Then, you **MUST** complete the following steps in order:
>
> 1.  **Run Full Regression Test:**
>     *   [x] **All Prior Tests Passed:** **Instruction:** `Execute 'go test ./...'.` **Evidence:** `evidence/PHASE-1/story-1.2/regression-test.log`
> 2.  **Create Git Commit:**
>     *   [x] **Work Committed:** **Instruction:** `Execute 'git add .' followed by 'git commit -m "feat(config): implement configuration and template registry"'.` **Evidence:** Commit hash: b1bcb52
> 3.  **Finalize Story:**
>     *   **Instruction:** Once the two checkboxes above are complete, you **MUST** update this story's main checkbox from `[ ]` to `[x]`.

---

#### [ ] STORY-1.3: Deterministic Rendering & DevOps Foundation

1.  **Task:** Implement the HTML renderer.
    *   **Instruction:** `In the /internal/render package, create a function that takes a map of field data (from a parsed JSON) and an HTML template. It MUST replace placeholders like '<!-- data-field="document.title" -->' with the corresponding data. This function must be pure and have no side effects other than returning the rendered string.`
    *   **Fulfills:** **PROD-003, PROD-004, ARCH-004, TECH-P-005**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-4.1`:**
            *   [x] **Test Method Created:** **Evidence:** Provided in `/internal/render/renderer_test.go`.
            *   [x] **Test Method Passed:** **Evidence:** `evidence/PHASE-1/story-1.3/task-1/test-output/TC-4.1.log` showing test passed.
    *   **Documentation:** [x] No documentation updates required for this task.
2.  **Task:** Create Dockerfile and basic CI workflow.
    *   **Instruction:** `Create a multi-stage Dockerfile that builds the Go binary and creates a minimal final image. Create a .github/workflows/ci.yml file that triggers on push, builds the container, and runs 'go build' and 'go test' inside it.`
    *   **Fulfills:** **DEV-001, DEV-004**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-5.1`:**
            *   [x] **Test Method Created:** **Evidence:** Created `.github/workflows/ci.yml` and saved to `evidence/PHASE-1/story-1.3/task-2/logs/ci.yml`.
            *   [x] **Test Method Passed:** **Evidence:** Docker build and test successful locally - `evidence/PHASE-1/story-1.3/task-2/logs/docker-build.log` and `docker-run.log`.
    *   **Documentation:**
        *   [x] **Documentation Updated:** **Instruction:** `Add a "Container Usage" section to the README explaining how to build and run the Docker image.` **Evidence:** Added Container Usage section to README.md.

---
> ### **Story Completion: STORY-1.3**
>
> You may only proceed once all checkboxes for all tasks within this story are marked `[x]`. Then, you **MUST** complete the following steps in order:
>
> 1.  **Run Full Regression Test:**
>     *   [ ] **All Prior Tests Passed:** **Instruction:** `Execute 'go test ./...'.` **Evidence:** Full summary output from the test runner.
> 2.  **Create Git Commit:**
>     *   [ ] **Work Committed:** **Instruction:** `Execute 'git add .' followed by 'git commit -m "feat(render): add deterministic renderer and initial CI pipeline"'.` **Evidence:** The commit hash.
> 3.  **Finalize Story:**
>     *   **Instruction:** Once the two checkboxes above are complete, you **MUST** update this story's main checkbox from `[ ]` to `[x]`.

---

### **4. Definition of Done**

This Phase is officially complete **only when all `STORY-1.x` checkboxes in Section 3 are marked `[x]` AND the Final Acceptance Gate below is passed.**

#### Final Acceptance Gate

*   **Instruction:** You are at the final gate for this phase. Before marking the entire phase as done, you must perform one last, full regression test to ensure nothing was broken by the final commits.
*   [ ] **Final Full Regression Test Passed:**
    *   **Instruction:** `Execute 'go test ./...'.`
    *   **Evidence:** Provide the full, final summary output from the test runner, showing the grand total of tests for this phase and confirming that 100% have passed.

*   **Final Instruction:** Once the `Final Full Regression Test Passed` checkbox above is marked `[x]`, your final action for this phase is to modify the main title of this document, changing `[ ] PHASE-1` to `[x] PHASE-1`. This concludes your work on this phase file.
