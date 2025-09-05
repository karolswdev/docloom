# docloom - Software Requirements Specification

**Version:** 1.0  
**Status:** Baseline

## Introduction

This document outlines the software requirements for **docloom**. It serves as the single source of truth for what the system must do, the constraints under which it must operate, and the rules governing its development and deployment.

Each requirement has a **unique, stable ID** (e.g., `PROD-001`). These IDs **MUST** be used to link implementation stories and test cases back to these foundational requirements, ensuring complete traceability.

The requirement keywords (`MUST`, `MUST NOT`, `SHOULD`, `SHOULD NOT`, `MAY`) are used as defined in RFC 2119.

---

## 1. Product & Functional Requirements

*Defines what the system does; its core features and capabilities.*

| ID | Title | Description | Rationale |
| :--- | :--- | :--- | :--- |
| <a name="PROD-001"></a>**PROD-001** | Generate Documents from Sources | The system **MUST** generate a complete document by combining provided source materials with a selected template type (e.g., Architecture Vision) and AI-generated content mapped to that template’s field schema. | Enables the core value proposition: automated, template‑driven document creation.
| <a name="PROD-002"></a>**PROD-002** | Template Registry | The system **MUST** support a registry of templates (initially: `architecture-vision`, `technical-debt-summary`, `reference-architecture`) each with a defined JSON field schema and prompt. | Ensures extensibility and consistent structure across document types.
| <a name="PROD-003"></a>**PROD-003** | HTML Output | For each template, the system **MUST** render a filled HTML document based on an HTML skeleton with `data-field` placeholders (e.g., `template/architecture-vision.html`). | Produces high‑fidelity, printable outputs aligned with branding.
| <a name="PROD-004"></a>**PROD-004** | Sidecar Field JSON | The system **MUST** output a JSON file containing the structured fields used to render the HTML. | Enables traceability, review, and re-rendering without re-calling the model.
| <a name="PROD-005"></a>**PROD-005** | Source Ingestion | The system **MUST** ingest local sources: directories and files of types `.md`, `.txt`, and `.pdf` (text extracted). | Covers common engineering inputs with minimal friction.
| <a name="PROD-006"></a>**PROD-006** | Chunking & Selection | The system **SHOULD** chunk and rank source content to fit model context; it **MAY** use heuristic ranking first, with embeddings added later. | Controls token cost and improves relevance.
| <a name="PROD-007"></a>**PROD-007** | Model Integration | The system **MUST** call an OpenAI‑compatible chat endpoint to produce structured JSON that conforms to the selected template’s schema. | Provides vendor‑agnostic AI integration.
| <a name="PROD-008"></a>**PROD-008** | Schema Validation & Repair | The system **MUST** validate model output against the JSON schema and **MUST** perform up to N automated repair attempts on failure. | Ensures reliable, machine‑processable outputs.
| <a name="PROD-009"></a>**PROD-009** | Determinism Controls | The system **SHOULD** support `temperature`, `seed`, and `max_tokens` controls to improve reproducibility. | Supports consistent outputs for review/audit.
| <a name="PROD-010"></a>**PROD-010** | Dry‑Run & Explainability | The system **SHOULD** provide a `--dry-run` mode that previews selected sources, prompts, and schema without calling the model. | Improves transparency and debuggability.
| <a name="PROD-011"></a>**PROD-011** | Citations (Optional) | The system **MAY** add source citations and an appendix when enabled via a flag. | Supports evidence‑based documents when needed.
| <a name="PROD-012"></a>**PROD-012** | Configurable Defaults | The system **MUST** load defaults from a config file (e.g., `docloom.yaml`) and environment variables, overridden by CLI flags. | Enables team‑wide standardization with local overrides.
| <a name="PROD-013"></a>**PROD-013** | Exit Codes & Errors | The system **MUST** return non‑zero exit codes on failure and **MUST** emit actionable error messages. | Integrates with CI and provides clear remediation.

---

## 2. User Interaction Requirements

*Defines how a user interacts with the system. Focuses on usability and user‑facing workflows.*

| ID | Title | Description | Rationale |
| :--- | :--- | :--- | :--- |
| <a name="USER-001"></a>**USER-001** | CLI `generate` Command | The user **MUST** be able to run `docloom generate --type <template> --source <paths...> --out <file>` with standard flags for model, base URL, API key, temperature, seed, and retries. | Provides a simple, consistent entry point.
| <a name="USER-002"></a>**USER-002** | Help & Usage | The CLI **MUST** offer `--help` with examples and template listings; `docloom templates list` **SHOULD** enumerate supported templates. | Improves discoverability.
| <a name="USER-003"></a>**USER-003** | Verbose Logging | A `--verbose` flag **MUST** print detailed steps (selected chunks, token estimates, validation results) without leaking secrets. | Aids troubleshooting while protecting credentials.
| <a name="USER-004"></a>**USER-004** | Dry‑Run Mode | A `--dry-run` flag **MUST** avoid network calls and print the prepared prompt, selected chunks, and target schema. | Enables safe inspection before cost is incurred.
| <a name="USER-005"></a>**USER-005** | Safe Writes | The CLI **MUST** avoid overwriting existing outputs unless `--force` is provided. | Prevents accidental data loss.
| <a name="USER-006"></a>**USER-006** | Config File | The user **MUST** be able to supply a config file via `--config` and set `API key` via env (`OPENAI_API_KEY`) or config secrets. | Simplifies setup.
| <a name="USER-007"></a>**USER-007** | Container Usage | The user **SHOULD** be able to run the tool via container (e.g., `docker run ghcr.io/karolswdev/docloom:TAG generate ...`) with bind-mounted sources and outputs. | Supports Docker‑first usage without local installs.

---

## 3. Architectural Requirements

*Defines high‑level, non‑negotiable design principles and structural constraints.*

| ID | Title | Description | Rationale |
| :--- | :--- | :--- | :--- |
| <a name="ARCH-001"></a>**ARCH-001** | Modular Package Layout | The system architecture **MUST** separate concerns: `ingest`, `chunk`, `prompt`, `ai`, `templates`, `render`, `validate`, `config`, and `cli`. | Improves maintainability and testability.
| <a name="ARCH-002"></a>**ARCH-002** | Provider‑Agnostic AI Client | The AI client **MUST** be behind an interface and **MUST** allow base URL and model selection to support OpenAI‑compatible providers. | Avoids vendor lock‑in.
| <a name="ARCH-003"></a>**ARCH-003** | Strict Schema Contracts | Template outputs **MUST** be defined by JSON Schema and validated before render. | Guarantees structural correctness.
| <a name="ARCH-004"></a>**ARCH-004** | Idempotent Rendering | Rendering **MUST** be pure given `fields.json` and template assets; no network calls during render. | Deterministic builds and reproducibility.
| <a name="ARCH-005"></a>**ARCH-005** | Secrets Handling | API keys **MUST** be read from env or config and **MUST NOT** be logged or persisted. Redaction **SHOULD** be applied to request/response logs. | Security best practices.
| <a name="ARCH-006"></a>**ARCH-006** | Extensible Template Registry | New templates **MUST** be addable without modifying core logic (registration pattern + asset folder + schema). | Facilitates growth.
| <a name="ARCH-007"></a>**ARCH-007** | Offline Friendly | The CLI **SHOULD** function in `--dry-run` and `--render-only --fields fields.json` modes without network. | Supports constrained environments.

---

## 4. Non-Functional Requirements (NFRs)

*Defines the quality attributes and operational characteristics of the system. The "-ilities".*

| ID | Title | Description | Rationale |
| :--- | :--- | :--- | :--- |
| <a name="NFR-001"></a>**NFR-001** | Performance: Local Steps | Ingestion, chunking, validation, and render steps (excluding model latency) **MUST** complete within 5 seconds for 200 pages of plain text on a typical developer laptop. | Maintains responsiveness independent of model performance.
| <a name="NFR-002"></a>**NFR-002** | Security: Credentials | API keys **MUST** be provided via env or encrypted config and **MUST NOT** be logged or included in crash reports. | Protects secrets.
| <a name="NFR-003"></a>**NFR-003** | Reliability: Retries | The system **MUST** retry model calls and schema repairs up to a configurable limit with exponential backoff. | Increases success rates under transient errors.
| <a name="NFR-004"></a>**NFR-004** | Portability | Binaries **MUST** run on Linux, macOS, and Windows on x86_64 and arm64 with Go 1.22+. | Broad developer reach.
| <a name="NFR-005"></a>**NFR-005** | Observability | Logs **MUST** be structured (JSON when `--json-logs` is set) and include request IDs and timings. | Supports debugging and CI integration.
| <a name="NFR-006"></a>**NFR-006** | Maintainability | Core packages (excluding networked AI calls) **SHOULD** have ≥80% unit test coverage; golden tests **SHOULD** validate prompts and render output. | Ensures code quality over time.

---

## 5. Technology & Platform Requirements

*Defines the specific technologies, frameworks, and platforms that are mandated for use.*

| ID | Title | Description | Rationale |
| :--- | :--- | :--- | :--- |
| <a name="TECH-P-001"></a>**TECH-P-001** | Primary Language | The application's backend **MUST** be implemented in Go 1.22+ with modules. | Performance, portability, and single-binary distribution.
| <a name="TECH-P-002"></a>**TECH-P-002** | CLI Framework | The CLI **MUST** use a proven framework (e.g., `spf13/cobra`). | Ergonomics and subcommand support.
| <a name="TECH-P-003"></a>**TECH-P-003** | AI Client | The AI client **MUST** support OpenAI‑compatible APIs with configurable `baseURL`, `model`, and timeouts. | Provider flexibility.
| <a name="TECH-P-004"></a>**TECH-P-004** | PDF Text Extraction | The system **SHOULD** use `pdftotext` when available; it **MAY** fall back to a Go library for PDF extraction. | Robust ingestion.
| <a name="TECH-P-005"></a>**TECH-P-005** | HTML Rendering | Rendering **MUST** operate on local HTML assets (e.g., `architecture-vision.html`, `style.css`, `terumo.css`, `logo.svg`) without network. | Deterministic output.
| <a name="TECH-P-006"></a>**TECH-P-006** | Module Path | The Go module path **MUST** be `github.com/karolswdev/docloom` and remain stable for importers. | Enables reproducible builds and third‑party integration.
| <a name="TECH-P-007"></a>**TECH-P-007** | Distribution Channels | The project **MUST** support both Docker image distribution and installable CLI binaries (and `go install`), documented with examples. | Ensures Docker‑first and local CLI workflows.

#### 5.1. Required Libraries & 3rd Party Components

This table defines the specific, approved bill of materials for the project. Additions or deviations from this list **MUST** be formally approved.

| Component / Library Name | Required Version | Purpose | Requirement ID |
| :--- | :--- | :--- | :--- |
| <a name="TECH-L-001"></a>`spf13/cobra` | ≥ 1.8.0 | CLI framework | **TECH-L-001** |
| <a name="TECH-L-002"></a>`github.com/sashabaranov/go-openai` | ≥ 1.22.0 | OpenAI‑compatible client | **TECH-L-002** |
| <a name="TECH-L-003"></a>`github.com/santhosh-tekuri/jsonschema/v5` | ≥ 5.3.0 | JSON Schema validation | **TECH-L-003** |
| <a name="TECH-L-004"></a>`github.com/rs/zerolog` | ≥ 1.33.0 | Structured logging | **TECH-L-004** |
| <a name="TECH-L-005"></a>`github.com/stretchr/testify` | ≥ 1.8.4 | Testing assertions | **TECH-L-005** |
| <a name="TECH-L-006"></a>`pdftotext` (poppler-utils) | System package | PDF text extraction (preferred) | **TECH-L-006** |

---

## 6. Operational & DevOps Requirements

*Defines the rules and constraints governing the development workflow, build process, and execution environment. These requirements are paramount and **MUST** be followed in all development and CI/CD activities.*

| ID | Title | Description | Rationale |
| :--- | :--- | :--- | :--- |
| <a name="DEV-001"></a>**DEV-001** | Containerized Workflow | All build, test, and run commands **MUST** be runnable in a containerized environment defined by the project's `Dockerfile` (local execution is permitted, but CI **MUST** use the container). | Ensures reproducibility and lowers onboarding friction.
| <a name="DEV-002"></a>**DEV-002** | Git Commit Standard | All Git commit messages **MUST** adhere to the Conventional Commits specification. | Enables automated changelogs and semantic versioning.
| <a name="DEV-003"></a>**DEV-003** | Static Code Analysis | `golangci-lint` **MUST** be integrated; CI **MUST** fail on new issues at or above "medium" severity. | Improves code quality and security.
| <a name="DEV-004"></a>**DEV-004** | CI/CD Pipeline | The project **MUST** include a CI workflow (e.g., `.github/workflows/ci.yml`) that runs build, tests, and static analysis on every push and PR. | Guarantees continuous validation.
| <a name="DEV-005"></a>**DEV-005** | Secrets Hygiene | CI **MUST** source secrets from the platform’s secret store; secrets **MUST NOT** be committed. Optional `gitleaks` **SHOULD** scan for leaks. | Prevents credential exposure.
| <a name="DEV-006"></a>**DEV-006** | Release Artifacts | CI **SHOULD** publish versioned, multi‑platform binaries (Linux/macOS/Windows; x86_64/arm64) on tagged releases. | Streamlines distribution.
| <a name="DEV-007"></a>**DEV-007** | Container Image Publish | CI **MUST** build and publish a versioned OCI image (e.g., `ghcr.io/karolswdev/docloom:<version>`) on tagged releases and `:main` for the default branch. | Provides Docker‑first distribution.
| <a name="DEV-008"></a>**DEV-008** | Installable CLI | Releases **MUST** include checksummed binaries and support installation via `go install github.com/karolswdev/docloom/cmd/docloom@latest`; a Homebrew tap **MAY** be provided. | Supports local CLI use without Docker.
| <a name="DEV-009"></a>**DEV-009** | Version Metadata | The binary **MUST** embed version, commit, and build date and expose them via `docloom --version`. | Aids support and reproducibility.
