## [ ] PHASE-2: Core AI Integration & Validation

---

### **1. Phase Context (What & Why)**

| ID | Title |
| :--- | :--- |
| PHASE-2 | Core AI Integration & Validation |

> **As a** Document Author, **I want** the system to call a language model with my source documents and a template prompt, **so that** I can receive structured, validated JSON content that is ready for rendering, forming the core automated value of the tool.

---

### **2. Phase Scope & Test Case Definitions (The Contract)**

*   **Requirement:** **PROD-007, ARCH-002, TECH-P-003, TECH-L-002** - [Provider-Agnostic AI Client](./SRS.md#ARCH-002)
    *   **Test Case ID:** `TC-6.1`
        *   **Test Method Signature:** `func TestAIClient_GenerateJSON_Success(t *testing.T)`
        *   **Test Logic:** (Arrange) Use an `httptest` mock server that returns a valid OpenAI-compatible JSON response. Instantiate the AI client to point to this mock server. (Act) Call the client's generation method. (Assert) The client correctly parses the response and returns the content string without error.
        *   **Required Proof of Passing:** Test runner output showing the test passed.
    *   **Test Case ID:** `TC-6.2`
        *   **Test Method Signature:** `func TestAIClient_GenerateJSON_RetriesOn503(t *testing.T)`
        *   **Test Logic:** (Arrange) Use a mock server that returns a 503 Service Unavailable error twice, then a 200 OK on the third call. Configure the client with 3 retries. (Act) Call the generation method. (Assert) The client should make three requests and ultimately succeed.
        *   **Required Proof of Passing:** Test runner output showing the test passed, with logs indicating retries occurred.
*   **Requirement:** **PROD-005** - [Source Ingestion](./SRS.md#PROD-005)
    *   **Test Case ID:** `TC-7.1`
        *   **Test Method Signature:** `func TestIngester_IngestSources(t *testing.T)`
        *   **Test Logic:** (Arrange) Create a test directory with `.md` and `.txt` files, and a subdirectory. (Act) Call the ingest function on the parent directory. (Assert) The function should return the concatenated content of all valid files, correctly traversing the subdirectory.
        *   **Required Proof of Passing:** Test runner output.
*   **Requirement:** **PROD-008, ARCH-003, TECH-L-003** - [Schema Validation & Repair](./SRS.md#PROD-008)
    *   **Test Case ID:** `TC-8.1`
        *   **Test Method Signature:** `func TestValidator_Validate_ValidJSON(t *testing.T)`
        *   **Test Logic:** (Arrange) Provide a valid JSON string and a corresponding JSON schema. (Act) Run the validation. (Assert) The function returns no error.
        *   **Required Proof of Passing:** Test runner output.
    *   **Test Case ID:** `TC-8.2`
        *   **Test Method Signature:** `func TestValidator_Validate_InvalidJSON(t *testing.T)`
        *   **Test Logic:** (Arrange) Provide a JSON string with a type error (e.g., a number where a string is expected) and its schema. (Act) Run validation. (Assert) The function returns a specific validation error detailing the incorrect field.
        *   **Required Proof of Passing:** Test runner output showing the expected error was returned.
    *   **Test Case ID:** `TC-8.3`
        *   **Test Method Signature:** `func TestGenerationFlow_RepairLoop(t *testing.T)`
        *   **Test Logic:** (Integration Test) (Arrange) Mock the AI client to first return invalid JSON, then valid JSON on the second call (the "repair" call). (Act) Run the main generation orchestrator. (Assert) The orchestrator should call the AI client twice and successfully return the valid JSON, logging the repair attempt.
        *   **Required Proof of Passing:** Test runner output with logs showing "Validation failed, attempting repair..." followed by success.
*   **Requirement:** **ARCH-005, NFR-002** - [Secrets Handling](./SRS.md#ARCH-005)
    *   **Test Case ID:** `TC-9.1`
        *   **Test Method Signature:** `func TestConfig_SecretRedactionInLogs(t *testing.T)`
        *   **Test Logic:** (Arrange) Set an API key via env var. Configure logging to a buffer. (Act) Run a command that would log the configuration, such as a verbose dry-run. (Assert) The API key in the log buffer MUST be redacted (e.g., `OPENAI_API_KEY: "sk-****"`)
        *   **Required Proof of Passing:** Log output from the test buffer showing the redacted key.

---

### **3. Implementation Plan (The Execution)**

#### [ ] STORY-2.1: AI Client & Source Ingestion

1.  **Task:** Implement a provider-agnostic AI client.
    *   **Instruction:** `In a new /internal/ai package, define an AIClient interface. Implement a concrete client using the go-openai library that communicates with an OpenAI-compatible endpoint. The client must be configurable for base URL, model, and API key.`
    *   **Fulfills:** **PROD-007, ARCH-002, TECH-P-003, TECH-L-002**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-6.1`:**
            *   [ ] **Test Method Created:** **Evidence:** Code for `TestAIClient_GenerateJSON_Success`.
            *   [ ] **Test Method Passed:** **Evidence:** `go test` output.
    *   **Documentation:** [ ] No documentation updates required for this task.
2.  **Task:** Implement retry logic for AI calls.
    *   **Instruction:** `Wrap the AI client's network calls with a retry mechanism that performs up to N retries with exponential backoff on transient errors (e.g., 5xx status codes).`
    *   **Fulfills:** **NFR-003**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-6.2`:**
            *   [ ] **Test Method Created:** **Evidence:** Code for `TestAIClient_GenerateJSON_RetriesOn503`.
            *   [ ] **Test Method Passed:** **Evidence:** `go test` output.
    *   **Documentation:** [ ] No documentation updates required for this task.
3.  **Task:** Implement basic source file ingestion.
    *   **Instruction:** `Create an /internal/ingest package. Implement a function that recursively walks a list of source paths, reading the content of all .md and .txt files into a single string.`
    *   **Fulfills:** **PROD-005 (partial)**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-7.1`:**
            *   [ ] **Test Method Created:** **Evidence:** Code for `TestIngester_IngestSources`.
            *   [ ] **Test Method Passed:** **Evidence:** `go test` output.
    *   **Documentation:** [ ] No documentation updates required for this task.

---
> ### **Story Completion: STORY-2.1**
>
> 1.  **Run Full Regression Test:**
>     *   [ ] **All Prior Tests Passed:** **Instruction:** `Execute 'go test ./...'.` **Evidence:** Full summary output.
> 2.  **Create Git Commit:**
>     *   [ ] **Work Committed:** **Instruction:** `Execute 'git commit -m "feat(ai): implement ai client and basic source ingestion"'.` **Evidence:** Commit hash.
> 3.  **Finalize Story:**
>     *   **Instruction:** Update this story's main checkbox to `[x]`.

---

#### [ ] STORY-2.2: Generation Workflow with Validation and Repair

1.  **Task:** Implement prompt assembly logic.
    *   **Instruction:** `Create an /internal/prompt package. Implement a builder that takes the ingested source content and a template's prompt text and assembles the final prompt to be sent to the AI model.`
    *   **Fulfills:** **PROD-001**.
    *   **Verification via Test Cases:** N/A (Will be tested via integration in next task).
    *   **Documentation:** [ ] No documentation updates required for this task.
2.  **Task:** Implement JSON Schema validation.
    *   **Instruction:** `Create an /internal/validate package. Using the jsonschema library, implement a function that validates a JSON string against a given schema string.`
    *   **Fulfills:** **PROD-008 (partial), ARCH-003, TECH-L-003**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-8.1`:**
            *   [ ] **Test Method Created:** **Evidence:** Code for `TestValidator_Validate_ValidJSON`.
            *   [ ] **Test Method Passed:** **Evidence:** `go test` output.
        *   **Test Case `TC-8.2`:**
            *   [ ] **Test Method Created:** **Evidence:** Code for `TestValidator_Validate_InvalidJSON`.
            *   [ ] **Test Method Passed:** **Evidence:** `go test` output.
    *   **Documentation:** [ ] No documentation updates required for this task.
3.  **Task:** Orchestrate the generate -> validate -> repair loop.
    *   **Instruction:** `In the 'generate' command's logic, tie everything together: Ingest sources -> Assemble prompt -> Call AI -> Validate response. If validation fails, construct a repair prompt (including the error) and re-call the AI up to N times. Ensure API keys are handled securely and not logged.`
    *   **Fulfills:** **PROD-008, PROD-013, ARCH-005, NFR-002**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-8.3`:**
            *   [ ] **Test Method Created:** **Evidence:** Code for `TestGenerationFlow_RepairLoop`.
            *   [ ] **Test Method Passed:** **Evidence:** `go test` output.
        *   **Test Case `TC-9.1`:**
            *   [ ] **Test Method Created:** **Evidence:** Code for `TestConfig_SecretRedactionInLogs`.
            *   [ ] **Test Method Passed:** **Evidence:** Log output from test.
    *   **Documentation:**
        *   [ ] **Documentation Updated:** **Instruction:** `Update the "Usage" section in README.md to show how to provide an API key via environment variable.` **Evidence:** Diff of `README.md`.

---
> ### **Story Completion: STORY-2.2**
>
> 1.  **Run Full Regression Test:**
>     *   [ ] **All Prior Tests Passed:** **Instruction:** `Execute 'go test ./...'.` **Evidence:** Full summary output.
> 2.  **Create Git Commit:**
>     *   [ ] **Work Committed:** **Instruction:** `Execute 'git commit -m "feat(core): implement full generation workflow with validation and repair"'.` **Evidence:** Commit hash.
> 3.  **Finalize Story:**
>     *   **Instruction:** Update this story's main checkbox to `[x]`.

---

### **4. Definition of Done**

This Phase is officially complete **only when all `STORY-2.x` checkboxes in Section 3 are marked `[x]` AND the Final Acceptance Gate below is passed.**

#### Final Acceptance Gate

*   [ ] **Final Full Regression Test Passed:**
    *   **Instruction:** `Execute 'go test ./...'.`
    *   **Evidence:** Provide the full, final summary output from the test runner.

*   **Final Instruction:** Once the `Final Full Regression Test Passed` checkbox is marked `[x]`, modify the main title to `[x] PHASE-2`.
