package templates

// Default template content with analysis prompts

const (
	// Architecture Vision template with analysis
	architectureVisionAnalysisSystem = `You are an expert software architect analyzing a codebase to create an Architecture Vision document. 
Your goal is to understand the system's structure, design patterns, and architectural decisions.
Use the available tools to explore the repository systematically, starting with high-level structure and drilling down into details as needed.`

	architectureVisionAnalysisUser = `Please analyze this repository to create a comprehensive Architecture Vision document. Follow these steps:
1. First, use tools to understand the overall repository structure
2. Identify key architectural patterns and design decisions
3. Analyze the technology stack and dependencies
4. Examine the system's components and their relationships
5. Generate a complete Architecture Vision document according to the schema

Focus on:
- System purpose and business goals
- Key architectural decisions and rationale
- Component structure and interactions
- Technology choices and trade-offs
- Quality attributes and constraints`

	// Technical Debt Summary template with analysis
	technicalDebtAnalysisSystem = `You are a senior engineer conducting a technical debt assessment.
Your role is to identify areas of technical debt, code quality issues, and improvement opportunities.
Use the available tools to analyze code quality, identify anti-patterns, and assess maintainability.`

	technicalDebtAnalysisUser = `Please analyze this repository to create a Technical Debt Summary. Follow these steps:
1. Examine the codebase structure for complexity and organization issues
2. Identify duplicated code, long methods, and large classes
3. Check for outdated dependencies and security vulnerabilities
4. Analyze test coverage and quality
5. Generate a prioritized technical debt report

Focus on:
- Code complexity and maintainability issues
- Missing or inadequate tests
- Outdated or vulnerable dependencies
- Architectural anti-patterns
- Recommended refactoring priorities`

	// Reference Architecture template with analysis
	referenceArchAnalysisSystem = `You are a principal architect creating a reference architecture document.
Your goal is to extract reusable patterns, best practices, and architectural guidelines from the codebase.
Use the available tools to identify exemplary implementations and patterns worth documenting.`

	referenceArchAnalysisUser = `Please analyze this repository to create a Reference Architecture document. Follow these steps:
1. Identify and document architectural patterns used
2. Extract reusable components and frameworks
3. Document best practices and conventions
4. Analyze cross-cutting concerns (security, logging, error handling)
5. Generate a comprehensive reference architecture guide

Focus on:
- Reusable architectural patterns
- Component templates and frameworks
- Development guidelines and standards
- Cross-cutting concern implementations
- Example implementations and usage patterns`
)

// UpdateDefaultTemplatesWithAnalysis adds analysis prompts to the default templates
func (r *Registry) UpdateDefaultTemplatesWithAnalysis() {
	// Update Architecture Vision template
	if tmpl, exists := r.templates["architecture-vision"]; exists {
		tmpl.Analysis = &Analysis{
			SystemPrompt:      architectureVisionAnalysisSystem,
			InitialUserPrompt: architectureVisionAnalysisUser,
		}
	}

	// Update Technical Debt Summary template
	if tmpl, exists := r.templates["technical-debt-summary"]; exists {
		tmpl.Analysis = &Analysis{
			SystemPrompt:      technicalDebtAnalysisSystem,
			InitialUserPrompt: technicalDebtAnalysisUser,
		}
	}

	// Update Reference Architecture template
	if tmpl, exists := r.templates["reference-architecture"]; exists {
		tmpl.Analysis = &Analysis{
			SystemPrompt:      referenceArchAnalysisSystem,
			InitialUserPrompt: referenceArchAnalysisUser,
		}
	}
}