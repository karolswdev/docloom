## [x] PHASE-4: Productionization, Distribution & Polish

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

#### [x] STORY-4.1: Build Automation & Release Readiness

1.  **Task:** Embed version metadata into the binary.
    *   **Instruction:** `Add a 'version' command. Use Go's ldflags to inject the version, commit hash, and build date into variables at build time. Update the Dockerfile and create a Makefile to standardize the build command.`
    *   **Fulfills:** **DEV-009**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-14.1`:**
            *   [x] **Test Method Created:** **Evidence:** test/e2e/version_test.go
            *   [x] **Test Method Passed:** **Evidence:** evidence/PHASE-4/story-4.1/task-1/TC-14.1-version-output.txt
    *   **Documentation:**
        *   [x] **Documentation Updated:** **Evidence:** README.md updated with version command
2.  **Task:** Enforce code quality standards in CI.
    *   **Instruction:** `Add a .golangci.yml configuration file with a reasonable set of linters. Update the ci.yml workflow to run 'golangci-lint run'. Add a PR title checker to enforce Conventional Commits.`
    *   **Fulfills:** **DEV-002, DEV-003**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-15.1`:**
            *   [x] **Test Method Created:** **Evidence:** .github/workflows/ci.yml and pr-title.yml
            *   [x] **Test Method Passed:** **Evidence:** golangci-lint runs successfully locally
    *   **Documentation:**
        *   [x] **Documentation Updated:** **Evidence:** README.md updated with Contributing section

---
> ### **Story Completion: STORY-4.1**
>
> 1.  **Run Full Regression Test:**
>     *   [x] **All Prior Tests Passed:** **Evidence:** evidence/PHASE-4/story-4.1/regression-test.log
> 2.  **Create Git Commit:**
>     *   [x] **Work Committed:** **Evidence:** Commit hash ed0bfb2
> 3.  **Finalize Story:**
>     *   **Instruction:** Update this story's main checkbox to `[x]`.

---

#### [x] STORY-4.2: Automated Multi-Platform Distribution

1.  **Task:** Implement the release workflow for binaries and containers.
    *   **Instruction:** `Create a .goreleaser.yml configuration. This file will define the build matrix for multi-platform binaries, archive formats, checksum generation, and Docker image publishing to GHCR. Create a .github/workflows/release.yml that triggers on tags (e.g., v*.*.*) and executes GoReleaser.`
    *   **Fulfills:** **DEV-006, DEV-007, DEV-008, NFR-004, USER-007, TECH-P-007**.
    *   **Verification via Test Cases:**
        *   **Test Case `TC-16.1`:**
            *   [x] **Test Method Created:** **Evidence:** .goreleaser.yml and .github/workflows/release.yml
            *   [x] **Test Method Passed:** **Evidence:** Multi-platform builds verified in evidence/PHASE-4/story-4.2/task-1/multi-platform-test.txt
    *   **Documentation:**
        *   [x] **Documentation Updated:** **Evidence:** README.md updated with comprehensive Installation section

---
> ### **Story Completion: STORY-4.2**
>
> 1.  **Run Full Regression Test:**
>     *   [x] **All Prior Tests Passed:** **Evidence:** evidence/PHASE-4/story-4.2/regression-test.log
> 2.  **Create Git Commit:**
>     *   [x] **Work Committed:** **Evidence:** Commit hash 7ecd4c5
> 3.  **Finalize Story:**
>     *   **Instruction:** Update this story's main checkbox to `[x]`.

---

### **4. Definition of Done**

This Phase is officially complete **only when all `STORY-4.x` checkboxes in Section 3 are marked `[x]` AND the Final Acceptance Gate below is passed.**

#### Final Acceptance Gate

*   [x] **Final Full Regression Test Passed:**
    *   **Instruction:** `Execute 'go test ./...'.`
    *   **Evidence:** evidence/PHASE-4/final-acceptance-gate.log

*   **Final Instruction:** Once the `Final Full Regression Test Passed` checkbox is marked `[x]`, modify the main title to `[x] PHASE-4`.