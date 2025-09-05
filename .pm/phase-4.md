## [ ] PHASE-4: Productionization, Distribution & Polish

---

### **1. Phase Context (What & Why)**

| ID | Title |
| :--- | :--- |
| PHASE-4 | Productionization, Distribution & Polish |

> **As a** DevOps Engineer and End User, **I want** a fully automated build, test, and release pipeline that distributes multi-platform binaries and Docker images, **so that** anyone can easily and reliably install and use `docloom` in their preferred environment.

---

### **2. Phase Scope & Test Case Definitions (The Contract)**

*   **Requirement:** **DEV-009** - [Version Metadata](./SRS.md#DEV-009)
    *   **Test Case ID:** `TC-14.1`
        *   **Test Method Signature:** `func TestVersionCmd(t *testing.T)`
        *   **Test Logic:** (E2E Test) Run the compiled binary with `docloom --version`. The test will require build-time variables to be set. (Assert) The output must contain a semantic version, commit hash, and build date, and they must not be empty.
        *   **Required Proof of Passing:** Console output from running the command.
*   **Requirement:** **DEV-002, DEV-003** - [Code Quality & Standards](./SRS.md#DEV-002)
    *   **Test Case ID:** `TC-15.1`
        *   **Test Method Signature:** N/A (CI Workflow Test)
        *   **Test Logic:** Update the `ci.yml` workflow to include a job that runs `golangci-lint run ./...`. Separately, add a commit-lint check for pull request titles.
        *   **Required Proof of Passing:** Link to a successful CI run that shows the linter passing.
*   **Requirement:** **DEV-006, DEV-007, DEV-008, NFR-004, USER-007** - [Release & Distribution](./SRS.md#DEV-006)
    *   **Test Case ID:** `TC-16.1`
        *   **Test Method Signature:** N/A (Release Workflow Test)
        *   **Test Logic:** Create a separate `release.yml` GitHub Actions workflow that triggers on Git tags. This workflow will use a tool like GoReleaser to build multi-platform binaries (Linux, macOS, Windows for amd64/arm64), create checksums, generate a changelog from Conventional Commits, and publish a GitHub Release. It must also build and push a version-tagged Docker image to GHCR.
        *   **Required Proof of Passing:** A link to a tagged release on GitHub showing the attached binaries and checksums, and a link to the corresponding container image on GHCR.

---

### **3. Implementation Plan (The Execution)**

#### [ ] STORY-4.1: Build Automation & Release Readiness

1.  **Task:** Embed version metadata into the binary.
    *   **Instruction:** `Add a 'version' command. Use Go's ldflags to inject the version, commit hash, and build date into variables at build time. Update the Dockerfile and create a Makefile to standardize the build command.`
    *   **Fulfills:** **DEV-009**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-14.1`:**
            *   [ ] **Test Method Created:** **Evidence:** Code for the end-to-end version command test.
            *   [ ] **Test Method Passed:** **Evidence:** Test runner output showing the command works as expected.
    *   **Documentation:**
        *   [ ] **Documentation Updated:** **Instruction:** `Add 'docloom --version' to the README's usage section.` **Evidence:** Diff of `README.md`.
2.  **Task:** Enforce code quality standards in CI.
    *   **Instruction:** `Add a .golangci.yml configuration file with a reasonable set of linters. Update the ci.yml workflow to run 'golangci-lint run'. Add a PR title checker to enforce Conventional Commits.`
    *   **Fulfills:** **DEV-002, DEV-003**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-15.1`:**
            *   [ ] **Test Method Created:** **Evidence:** The updated YAML content for `ci.yml`.
            *   [ ] **Test Method Passed:** **Evidence:** Link to a successful CI run showing the lint job passed.
    *   **Documentation:**
        *   [ ] **Documentation Updated:** **Instruction:** `Add a "Contributing" section to the README that specifies the Conventional Commits standard.` **Evidence:** Diff of `README.md`.

---
> ### **Story Completion: STORY-4.1**
>
> 1.  **Run Full Regression Test:**
>     *   [ ] **All Prior Tests Passed:** **Instruction:** `Execute 'go test ./...'.` **Evidence:** Full summary output.
> 2.  **Create Git Commit:**
>     *   [ ] **Work Committed:** **Instruction:** `Execute 'git commit -m "feat(devops): embed version info and enforce linting in CI"'.` **Evidence:** Commit hash.
> 3.  **Finalize Story:**
>     *   **Instruction:** Update this story's main checkbox to `[x]`.

---

#### [ ] STORY-4.2: Automated Multi-Platform Distribution

1.  **Task:** Implement the release workflow for binaries and containers.
    *   **Instruction:** `Create a .goreleaser.yml configuration. This file will define the build matrix for multi-platform binaries, archive formats, checksum generation, and Docker image publishing to GHCR. Create a .github/workflows/release.yml that triggers on tags (e.g., v*.*.*) and executes GoReleaser.`
    *   **Fulfills:** **DEV-006, DEV-007, DEV-008, NFR-004, USER-007, TECH-P-007**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-16.1`:**
            *   [ ] **Test Method Created:** **Evidence:** Content of `.goreleaser.yml` and `release.yml`.
            *   [ ] **Test Method Passed:** **Evidence:** Link to a test release (e.g., `v0.1.0-test`) on GitHub and GHCR.
    *   **Documentation:**
        *   [ ] **Documentation Updated:** **Instruction:** `Create a comprehensive "Installation" section in the README. Provide instructions for downloading binaries from GitHub Releases, using 'go install', and using the Docker image.` **Evidence:** Final `README.md` content.

---
> ### **Story Completion: STORY-4.2**
>
> 1.  **Run Full Regression Test:**
>     *   [ ] **All Prior Tests Passed:** **Instruction:** `Execute 'go test ./...'.` **Evidence:** Full summary output.
> 2.  **Create Git Commit:**
>     *   [ ] **Work Committed:** **Instruction:** `Execute 'git commit -m "feat(release): implement automated releases with goreleaser"'.` **Evidence:** Commit hash.
> 3.  **Finalize Story:**
>     *   **Instruction:** Update this story's main checkbox to `[x]`.

---

### **4. Definition of Done**

This Phase is officially complete **only when all `STORY-4.x` checkboxes in Section 3 are marked `[x]` AND the Final Acceptance Gate below is passed.**

#### Final Acceptance Gate

*   [ ] **Final Full Regression Test Passed:**
    *   **Instruction:** `Execute 'go test ./...'.`
    *   **Evidence:** Provide the full, final summary output from the test runner.

*   **Final Instruction:** Once the `Final Full Regression Test Passed` checkbox is marked `[x]`, modify the main title to `[x] PHASE-4`.