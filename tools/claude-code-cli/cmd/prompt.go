package cmd

import (
	"fmt"
	"strings"
)

// PromptGenerator generates analysis prompts for Claude
type PromptGenerator struct{}

// NewPromptGenerator creates a new prompt generator
func NewPromptGenerator() *PromptGenerator {
	return &PromptGenerator{}
}

// Generate creates an analysis prompt from scan results
func (p *PromptGenerator) Generate(scan *ScanResult) string {
	var prompt strings.Builder
	
	prompt.WriteString("You are an expert C# software architect analyzing a repository to generate comprehensive documentation.\n\n")
	prompt.WriteString("Please analyze the following C# repository and provide a structured JSON response with detailed insights.\n\n")
	
	// Add repository structure overview
	prompt.WriteString("## Repository Overview\n\n")
	if scan.SolutionFile != "" {
		prompt.WriteString(fmt.Sprintf("Solution File: %s\n", scan.SolutionFile))
	}
	if len(scan.ProjectFiles) > 0 {
		prompt.WriteString(fmt.Sprintf("Project Files: %s\n", strings.Join(scan.ProjectFiles, ", ")))
	}
	if len(scan.ReadmeFiles) > 0 {
		prompt.WriteString(fmt.Sprintf("README Files: %s\n", strings.Join(scan.ReadmeFiles, ", ")))
	}
	prompt.WriteString("\n")
	
	// Add file contents
	prompt.WriteString("## Key Files\n\n")
	
	// Prioritize files by type
	fileTypes := []string{"solution", "project", "readme", "config", "source"}
	for _, fileType := range fileTypes {
		hasType := false
		for _, file := range scan.Files {
			if file.FileType == fileType {
				if !hasType {
					prompt.WriteString(fmt.Sprintf("### %s Files\n\n", strings.Title(fileType)))
					hasType = true
				}
				prompt.WriteString(fmt.Sprintf("#### File: %s\n```\n%s\n```\n\n", file.RelPath, file.Content))
			}
		}
	}
	
	// Add analysis instructions
	prompt.WriteString("## Analysis Requirements\n\n")
	prompt.WriteString("Based on the repository contents above, provide a JSON response with the following structure:\n\n")
	prompt.WriteString("```json\n")
	prompt.WriteString(`{
  "projectName": "Name of the project",
  "description": "Brief description of what this project does",
  "projectType": "Type of project (Web API, Console App, Library, etc.)",
  "framework": "Target framework (e.g., .NET 6, .NET Core 3.1)",
  "architecture": {
    "pattern": "Architectural pattern used (e.g., MVC, Clean Architecture, etc.)",
    "layers": ["List of architectural layers"],
    "keyComponents": ["List of key components or services"]
  },
  "dependencies": {
    "nuget": ["List of key NuGet packages"],
    "external": ["External services or APIs"]
  },
  "features": [
    {
      "name": "Feature name",
      "description": "What this feature does"
    }
  ],
  "apis": [
    {
      "endpoint": "API endpoint path",
      "method": "HTTP method",
      "description": "What this endpoint does"
    }
  ],
  "dataModel": {
    "entities": ["List of main domain entities"],
    "database": "Database technology if any"
  },
  "testing": {
    "framework": "Testing framework used",
    "coverage": "Estimated test coverage",
    "types": ["Types of tests (unit, integration, etc.)"]
  },
  "deployment": {
    "containerized": true/false,
    "cicd": "CI/CD platform if any",
    "hosting": "Target hosting environment"
  },
  "security": {
    "authentication": "Authentication method",
    "authorization": "Authorization approach",
    "considerations": ["Security considerations or features"]
  },
  "technicalDebt": [
    {
      "area": "Area of concern",
      "description": "Description of the technical debt",
      "impact": "High/Medium/Low",
      "recommendation": "Suggested improvement"
    }
  ],
  "recommendations": [
    "List of recommendations for improvement"
  ]
}`)
	prompt.WriteString("\n```\n\n")
	prompt.WriteString("Provide your analysis as valid JSON that matches this structure. Be thorough and specific based on the actual code you've seen.\n")
	
	return prompt.String()
}