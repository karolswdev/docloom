# Creating Intelligent Templates

## Overview

Intelligent templates in DocLoom combine structured document schemas with AI analysis prompts to create goal-oriented, context-aware documentation. This guide teaches you how to craft analysis prompts that guide the AI's research process to align with your document's objectives.

## The Power of Template-Driven Intelligence

Traditional documentation tools require you to manually gather and organize information. DocLoom's intelligent templates flip this paradigm: you define the document's goal and structure, and the AI autonomously explores your codebase to fulfill that goal.

## Anatomy of an Intelligent Template

An intelligent template consists of:

1. **Schema**: Defines the document structure and data fields
2. **HTML Template**: Provides the visual layout
3. **Analysis Prompts**: Guide the AI's research strategy
4. **Generation Prompt**: Structures the final output

### The Analysis Section

```yaml
analysis:
  system_prompt: |
    Define the AI's role, expertise, and approach.
    This sets the context for the entire analysis.
  
  initial_user_prompt: |
    Provide specific instructions on what to analyze
    and how to approach the exploration.
```

## Crafting Effective System Prompts

The system prompt establishes the AI's persona and expertise. It should:

### 1. Define the Role

```yaml
system_prompt: |
  You are a senior security architect conducting a comprehensive security audit.
  Your expertise includes OWASP Top 10, secure coding practices, and threat modeling.
```

### 2. Set the Analytical Framework

```yaml
system_prompt: |
  You are a performance engineer analyzing system bottlenecks.
  Apply the USE method (Utilization, Saturation, Errors) and focus on:
  - Resource consumption patterns
  - Concurrency and synchronization issues
  - Database query optimization
  - Network latency factors
```

### 3. Establish Tool Usage Strategy

```yaml
system_prompt: |
  You are a documentation specialist creating API references.
  Use the available tools systematically:
  1. First, discover all API endpoints
  2. For each endpoint, analyze parameters and responses
  3. Identify authentication and authorization patterns
  4. Extract example usage from tests
```

## Designing Initial User Prompts

The initial user prompt provides specific instructions for the analysis:

### 1. Define Clear Objectives

```yaml
initial_user_prompt: |
  Analyze this repository to create a Migration Guide from version 2.x to 3.x.
  
  Focus on:
  - Breaking API changes
  - Deprecated features and their replacements
  - New capabilities and how to adopt them
  - Step-by-step migration procedures
```

### 2. Specify Analysis Sequence

```yaml
initial_user_prompt: |
  Perform a comprehensive architectural review following these steps:
  
  1. Map the overall system structure
     - Identify main components and services
     - Trace data flow between components
  
  2. Analyze architectural patterns
     - Identify design patterns used
     - Evaluate consistency of implementation
  
  3. Assess quality attributes
     - Scalability mechanisms
     - Fault tolerance strategies
     - Security boundaries
```

### 3. Provide Evaluation Criteria

```yaml
initial_user_prompt: |
  Evaluate this codebase for production readiness:
  
  Critical Requirements:
  - Test coverage must exceed 80%
  - All public APIs must have documentation
  - No high-severity security vulnerabilities
  
  Important Factors:
  - Logging and monitoring implementation
  - Error handling consistency
  - Configuration management approach
```

## Template Examples by Document Type

### Architecture Documentation

```yaml
analysis:
  system_prompt: |
    You are a principal architect documenting system architecture.
    Your expertise covers distributed systems, microservices, and cloud patterns.
    Approach the analysis as if conducting an architecture review for stakeholders.
  
  initial_user_prompt: |
    Create an Architecture Decision Record (ADR) by:
    1. Identifying key architectural decisions made in the codebase
    2. Inferring the context and drivers for each decision
    3. Analyzing the consequences and trade-offs
    4. Documenting alternatives that might have been considered
```

### Security Audit

```yaml
analysis:
  system_prompt: |
    You are a security researcher performing a security assessment.
    Apply STRIDE threat modeling and check for CWE/SANS Top 25 vulnerabilities.
    Consider both code-level and architectural security concerns.
  
  initial_user_prompt: |
    Conduct a security audit focusing on:
    1. Authentication and authorization implementation
    2. Input validation and sanitization
    3. Cryptographic practices
    4. Dependency vulnerabilities
    5. Secrets management
    
    For each finding, assess severity and provide remediation guidance.
```

### Performance Analysis

```yaml
analysis:
  system_prompt: |
    You are a performance engineer optimizing system efficiency.
    Focus on algorithmic complexity, resource usage, and scalability limits.
    Use Big-O analysis and identify bottlenecks.
  
  initial_user_prompt: |
    Analyze performance characteristics:
    1. Identify computationally expensive operations
    2. Find database query patterns and potential N+1 problems
    3. Analyze concurrent code for lock contention
    4. Check for memory leaks or excessive allocations
    5. Evaluate caching strategies
```

### API Documentation

```yaml
analysis:
  system_prompt: |
    You are an API designer creating developer-friendly documentation.
    Focus on clarity, completeness, and practical examples.
    Consider both REST principles and developer experience.
  
  initial_user_prompt: |
    Document the API comprehensively:
    1. Catalog all endpoints with methods and paths
    2. Document request/response schemas with examples
    3. Identify authentication requirements
    4. Extract rate limiting and quota information
    5. Find integration examples from tests
    6. Note any versioning or deprecation patterns
```

## Advanced Prompt Techniques

### 1. Conditional Analysis

```yaml
initial_user_prompt: |
  Analyze the testing strategy:
  
  If unit tests exist:
    - Calculate coverage percentages
    - Identify untested critical paths
    - Evaluate test quality and assertions
  
  If integration tests exist:
    - Map test scenarios to user journeys
    - Check for environment dependencies
    - Assess test data management
  
  If no tests exist:
    - Identify highest-risk areas needing tests
    - Suggest testing strategy based on architecture
```

### 2. Comparative Analysis

```yaml
initial_user_prompt: |
  Compare this implementation with industry best practices:
  
  For each major component:
  1. Identify the pattern or approach used
  2. Compare with standard implementations
  3. Note deviations and assess if justified
  4. Suggest improvements where applicable
```

### 3. Progressive Refinement

```yaml
initial_user_prompt: |
  Build understanding progressively:
  
  Level 1 - Overview:
    - Repository structure and organization
    - Main technologies and frameworks
  
  Level 2 - Components:
    - Individual service responsibilities
    - Inter-component communication
  
  Level 3 - Implementation:
    - Core algorithms and data structures
    - Critical business logic
  
  Level 4 - Quality:
    - Code quality metrics
    - Technical debt assessment
```

## Optimizing for Tool Usage

### Guide Tool Selection

```yaml
system_prompt: |
  When analyzing, use tools in this priority:
  1. list_projects - Always start here for repository overview
  2. get_dependencies - Understand the technology stack
  3. get_api_surface - For public interface analysis
  4. get_file_content - For detailed implementation review
  
  Minimize redundant tool calls by caching insights mentally.
```

### Specify Information Extraction

```yaml
initial_user_prompt: |
  When using the get_file_content tool, focus on:
  - Configuration files (*.config, *.yaml, *.json)
  - Entry points (main.*, index.*, app.*)
  - Interface definitions (*.proto, *.graphql, swagger.*)
  - Test files matching critical components
  
  Extract patterns rather than memorizing entire files.
```

## Testing Your Prompts

### 1. Clarity Test

- Can another engineer understand the analysis goal?
- Are success criteria clearly defined?
- Is the sequence logical and complete?

### 2. Completeness Test

- Does the prompt cover all aspects of the document schema?
- Are edge cases considered?
- Will the analysis gather sufficient information?

### 3. Efficiency Test

- Is the analysis sequence optimized?
- Are tool calls minimized while maintaining thoroughness?
- Can the AI complete analysis within token limits?

## Common Pitfalls to Avoid

### 1. Vague Instructions

❌ **Poor:**
```yaml
initial_user_prompt: |
  Analyze the code and create documentation.
```

✅ **Better:**
```yaml
initial_user_prompt: |
  Analyze the codebase to create a Developer Onboarding Guide:
  1. Map project structure and key directories
  2. Document setup and build procedures
  3. Identify main architectural components
  4. Extract coding conventions from existing code
  5. Find examples of common development tasks
```

### 2. Assumption of Knowledge

❌ **Poor:**
```yaml
system_prompt: |
  You understand this system's architecture.
```

✅ **Better:**
```yaml
system_prompt: |
  You are analyzing an unfamiliar codebase. Start with no assumptions.
  Discover the architecture through systematic exploration using available tools.
  Build understanding incrementally from structure to implementation details.
```

### 3. Missing Context

❌ **Poor:**
```yaml
initial_user_prompt: |
  Find all the security issues.
```

✅ **Better:**
```yaml
initial_user_prompt: |
  Perform security analysis for a web application that handles financial data:
  - Compliance requirements: PCI DSS, SOX
  - Threat model: External attackers, insider threats
  - Critical assets: Payment data, user credentials, transaction logs
  - Check for: OWASP Top 10, CWE Top 25, business logic flaws
```

## Prompt Templates Library

### For Greenfield Projects

```yaml
analysis:
  system_prompt: |
    You are a technical advisor evaluating a new codebase for adoption.
    Assess maturity, quality, and fitness for purpose.
  
  initial_user_prompt: |
    Evaluate this project for potential adoption:
    1. Assess project maturity (docs, tests, CI/CD)
    2. Evaluate code quality and maintainability
    3. Check dependency health and licensing
    4. Identify integration points and requirements
    5. Estimate adoption effort and risks
```

### For Legacy Systems

```yaml
analysis:
  system_prompt: |
    You are a modernization consultant analyzing legacy systems.
    Focus on understanding current state and modernization opportunities.
  
  initial_user_prompt: |
    Analyze this legacy system for modernization:
    1. Document current architecture and technology stack
    2. Identify technical debt and obsolete patterns
    3. Find tightly coupled components
    4. Assess testability and test coverage
    5. Recommend incremental modernization strategy
```

### For Compliance Documentation

```yaml
analysis:
  system_prompt: |
    You are a compliance officer documenting regulatory adherence.
    Apply relevant standards (GDPR, HIPAA, SOC2) to the analysis.
  
  initial_user_prompt: |
    Document compliance posture:
    1. Identify data types and classifications
    2. Trace data flow and storage
    3. Document access controls and authentication
    4. Find audit logging implementation
    5. Assess encryption and data protection
    6. Map findings to compliance requirements
```

## Conclusion

Creating intelligent templates is both an art and a science. The key is to:

1. **Think like a researcher**: How would you manually explore the codebase?
2. **Be specific**: Clear instructions yield better results
3. **Consider the audience**: Who will read the final document?
4. **Iterate and refine**: Test with different codebases and adjust

Well-crafted analysis prompts transform DocLoom from a documentation tool into an intelligent analysis platform that understands your specific documentation needs and autonomously fulfills them.