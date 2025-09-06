## [ ] PHASE-8: The Interactive Research Agent & AI Analysis Loop

---

### **1. Phase Context (What & Why)**

| ID | Title |
| :--- | :--- |
| PHASE-8 | The Interactive Research Agent & AI Analysis Loop |

> **As a** Senior Engineer, **I want** to provide `docloom` with a repository and have the AI autonomously direct a research agent to explore the codebase, read relevant files, and synthesize a deep understanding of the project, **so that** I can generate a comprehensive and accurate technical document with near-zero manual pre-analysis.

---

### **2. Phase Scope & Test Case Definitions (The Contract)**

*   **Requirement:** **(Implied by Vision)** - **Agent-as-Toolkit**
    *   **Test Case ID:** `TC-24.1`
        *   **Test Method Signature:** `func TestAgentExecutor_RunTool(t *testing.T)`
        *   **Test Logic:** (Arrange) Define an agent YAML with multiple distinct `tools` (e.g., `list_projects`, `get_file_content`). Create a mock agent binary that responds to these tool names as subcommands. (Act) Invoke the Agent Executor, specifying a single tool to run. (Assert) The executor must successfully call the agent binary with the correct subcommand and arguments, and capture its stdout.
        *   **Required Proof of Passing:** Test runner output confirming the correct tool was invoked and its specific output was returned.
*   **Requirement:** **(Implied by Vision)** - **LLM-Orchestrated Analysis Loop**
    *   **Test Case ID:** `TC-25.1`
        *   **Test Method Signature:** `func TestOrchestrator_AnalysisLoop(t *testing.T)`
        *   **Test Logic:** (Integration Test) (Arrange) Mock the AI Client to return a sequence of responses: first, a request to call the `list_projects` tool; second, a request to call `get_file_content` on a specific project file. Mock the Agent Executor to return predefined data for these calls. (Act) Run the `generate` command with the agent. (Assert) The orchestrator must correctly interpret the AI's tool-use requests, call the Agent Executor twice with the correct parameters, send the results back to the AI, and finally receive the final JSON for rendering.
        *   **Required Proof of Passing:** Debug logs showing the full, multi-turn conversation flow: `(User Prompt) -> (AI Tool Call 1) -> (Agent Result 1) -> (AI Tool Call 2) -> (Agent Result 2) -> (Final AI JSON)`.
*   **Requirement:** **(Implied by Vision)** - **Template-Defined Intelligence**
    *   **Test Case ID:** `TC-26.1`
        *   **Test Method Signature:** `func TestGenerateCmd_UsesTemplateAnalysisPrompt(t *testing.T)`
        *   **Test Logic:** (E2E Test) (Arrange) Create two templates, A and B, each with a different `analysis_prompt`. Mock the AI Client to capture the initial prompt. (Act) Run `docloom generate --agent ... --type template-a` and then `... --type template-b`. (Assert) The initial prompt sent to the AI must exactly match the `analysis_prompt` from the corresponding template file, proving that the analysis is driven by the template's specific goals.
        *   **Required Proof of Passing:** The captured initial prompts sent to the mocked AI client for both runs.

---

### **3. Implementation Plan (The Execution)**

#### [x] STORY-8.1: Evolving the Agent into a Toolkit

1.  **Task:** Redesign the `agent.agent.yaml` specification to support tools.
    *   **Instruction:** `Update the Agent Definition Go structs and the YAML schema. Replace the single 'runner' field with a 'tools' array. Each tool MUST have a 'name', a 'description' (for the LLM), and a 'command'. The 'description' is critical—it's how the LLM will know when to use the tool.`
    *   **Fulfills:** `Agent-as-Toolkit` architecture.
    *   **Verification via Test Cases:** N/A (Schema change, verified by executor).
    *   **Documentation:**
        *   [x] **Documentation Updated:** **Instruction:** `Rewrite docs/agents/schema.md to reflect the new tool-based architecture. Provide a detailed example of an agent.agent.yaml with multiple tools. Emphasize the importance of writing clear, actionable descriptions for each tool, as this is the primary interface for the AI.` **Evidence:** The updated `schema.md` content.
2.  **Task:** Refactor the C# Analyzer into a multi-tool binary.
    *   **Instruction:** `Modify the Go-native C# agent from Phase 7. Instead of one main function, use subcommands (e.g., with Cobra for Go) to expose different functionalities: 'summarize_readme', 'list_projects', 'get_dependencies', 'get_api_surface', and 'get_file_content'.`
    *   **Fulfills:** Go-native agent implementation of the toolkit pattern.
    *   **Verification via Test Cases:** N/A (Internal agent change, verified E2E).
    *   **Documentation:**
        *   [x] **Documentation Updated:** **Instruction:** `Update docs/agents/csharp-analyzer.md to document each new subcommand (tool), its purpose, and the format of its output. This serves as the developer documentation for the agent itself.` **Evidence:** Diff of `csharp-analyzer.md`.
3.  **Task:** Upgrade the Agent Executor to invoke specific tools.
    *   **Instruction:** `Update the /internal/agent/Executor to accept a tool name as an argument. The executor will now be responsible for looking up the specific tool's command from the agent definition and executing it.`
    *   **Fulfills:** `Agent-as-Toolkit` architecture.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-24.1`:**
            *   [x] **Test Method Created:** **Evidence:** Provide the Go test code `TestAgentExecutor_RunTool`.
            *   [x] **Test Method Passed:** **Evidence:** Console output from `go test`.
    *   **Documentation:**
        *   [x] **Documentation Updated:** **Instruction:** `Update docs/agents/authoring-guide.md to reflect the new, more powerful tool-based paradigm. Provide a clear example of how an agent binary should parse subcommands to function as a toolkit.` **Evidence:** The updated authoring guide.

---
> ### **Story Completion: STORY-8.1**
> 1.  **Run Full Regression Test:**
>     *   [x] **All Prior Tests Passed:** **Instruction:** `Execute 'go test ./...'.` **Evidence:** Full summary output.
> 2.  **Create Git Commit:**
>     *   [x] **Work Committed:** **Instruction:** `git commit -m "refactor(agent): evolve agents to a tool-based architecture"`. **Evidence:** Commit hash: 88f0d4a.
> 3.  **Finalize Story:**
>     *   **Instruction:** Update this story's main checkbox to `[x]`.

---

#### [ ] STORY-8.2: The AI Analysis Loop

1.  **Task:** Upgrade the AI client to support Tool Calling / Function Calling.
    *   **Instruction:** `Update the /internal/ai/Client interface and the OpenAI implementation to support the tool-calling feature. This involves passing the agent's tool definitions to the API and correctly parsing the AI's response, which may be a final message or a request to call a tool.`
    *   **Fulfills:** Core requirement for the analysis loop.
    *   **Verification via Test Cases:** N/A (Internal change, verified in integration).
    *   **Documentation:** [ ] No documentation updates required for this task.
2.  **Task:** Implement the multi-turn Analysis Loop in the Orchestrator.
    *   **Instruction:** `In the /internal/generate/Orchestrator, create a new private method, runAnalysisLoop. This method will contain a loop that: 1. Sends the current conversation history and available tools to the LLM. 2. Checks if the response is a tool call or a final answer. 3. If a tool call, uses the Agent Executor to run the tool. 4. Appends the tool's output to the conversation history. 5. Repeats until the LLM provides a final answer or a max turn limit is reached.`
    *   **Fulfills:** `LLM-Orchestrated Analysis Loop`.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-25.1`:**
            *   [ ] **Test Method Created:** **Evidence:** Provide the Go integration test code `TestOrchestrator_AnalysisLoop`.
            *   [ ] **Test Method Passed:** **Evidence:** Console output from `go test`.
    *   **Documentation:**
        *   [ ] **Documentation Updated:** **Instruction:** `Create a new architecture document at docs/architecture/analysis-loop.md. Fully document this new, sophisticated workflow. Include a Mermaid sequence diagram showing the interaction between the User, Orchestrator, AI Client, and Agent Executor.` **Evidence:** The new `analysis-loop.md` file with the diagram.

---
> ### **Story Completion: STORY-8.2**
> 1.  **Run Full Regression Test:**
>     *   [ ] **All Prior Tests Passed:** **Instruction:** `Execute 'go test ./...'.` **Evidence:** Full summary output.
> 2.  **Create Git Commit:**
>     *   [ ] **Work Committed:** **Instruction:** `feat(core): implement LLM-orchestrated analysis loop for agents"`. **Evidence:** Commit hash.
> 3.  **Finalize Story:**
>     *   **Instruction:** Update this story's main checkbox to `[x]`.

---

#### [ ] STORY-8.3: Template-Driven Intelligence

1.  **Task:** Update the Template definition to include analysis prompts.
    *   **Instruction:** `Add a new 'Analysis' struct to the /internal/templates/Template definition. This struct will contain fields for 'SystemPrompt' and 'InitialUserPrompt'. These will be loaded from a new 'analysis' section in the template's definition file.`
    *   **Fulfills:** `Template-Defined Intelligence`.
    *   **Verification via Test Cases:** N/A (Schema change, verified by next task).
    *   **Documentation:**
        *   [ ] **Documentation Updated:** **Instruction:** `Update docs/templates/schema.md to include the new 'analysis' section. Provide a powerful example of a system prompt that instructs the AI on how to act as a research director and use the available tools.` **Evidence:** Diff of the template schema documentation.
2.  **Task:** Integrate template analysis prompts into the Orchestrator.
    *   **Instruction:** `Modify the Orchestrator's Generate method. When an agent is used, it must now load the analysis prompts from the selected template and use them to begin the conversation in the analysis loop.`
    *   **Fulfills:** `Template-Defined Intelligence`.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-26.1`:**
            *   [ ] **Test Method Created:** **Evidence:** Provide the E2E test code `TestGenerateCmd_UsesTemplateAnalysisPrompt`.
            *   [ ] **Test Method Passed:** **Evidence:** Console output from `go test`.
    *   **Documentation:**
        *   [ ] **Documentation Updated:** **Instruction:** `Write a new guide: docs/guides/creating-intelligent-templates.md. This is a world-class document that teaches users *how to think* about writing analysis prompts. Explain how to guide the AI's research process to align with the document's goal (e.g., an architecture vision prompt focuses on high-level structure, while a security audit prompt asks the AI to look for specific vulnerability patterns).` **Evidence:** The new guide's content.

---
> ### **Story Completion: STORY-8.3**
> 1.  **Run Full Regression Test:**
>     *   [ ] **All Prior Tests Passed:** **Instruction:** `Execute 'go test ./...'.` **Evidence:** Full summary output.
> 2.  **Create Git Commit:**
>     *   [ ] **Work Committed:** **Instruction:** `feat(templates): add analysis prompts to templates for goal-oriented research"`. **Evidence:** Commit hash.
> 3.  **Finalize Story:**
>     *   **Instruction:** Update this story's main checkbox to `[x]`.

---

### **4. Definition of Done**

This Phase is officially complete **only when all `STORY-8.x` checkboxes in Section 3 are marked `[x]` AND the Final Acceptance Gate below is passed.**

#### Final Acceptance Gate

*   [ ] **Final Full Regression Test Passed:**
    *   **Instruction:** `Execute 'go test ./...'.`
    *   **Evidence:** Provide the full, final summary output from the test runner.

*   **Final Instruction:** Once the `Final Full Regression Test Passed` checkbox is marked `[x]`, modify the main title to `[x] PHASE-8`.