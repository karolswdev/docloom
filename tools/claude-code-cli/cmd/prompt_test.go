package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPromptGenerator_Generate(t *testing.T) {
	// Create test scan result
	scanResult := &ScanResult{
		RootPath:     "/test/repo",
		SolutionFile: "TestApp.sln",
		ProjectFiles: []string{"src/TestApp.csproj", "tests/TestApp.Tests.csproj"},
		ReadmeFiles:  []string{"README.md"},
		Files: []FileInfo{
			{
				RelPath:  "TestApp.sln",
				Content:  "Solution file content",
				FileType: "solution",
			},
			{
				RelPath:  "src/TestApp.csproj",
				Content:  "<Project Sdk=\"Microsoft.NET.Sdk.Web\"></Project>",
				FileType: "project",
			},
			{
				RelPath:  "README.md",
				Content:  "# Test Application\nThis is a test.",
				FileType: "readme",
			},
			{
				RelPath:  "src/Program.cs",
				Content:  "class Program { }",
				FileType: "source",
			},
		},
	}
	
	generator := NewPromptGenerator()
	prompt := generator.Generate(scanResult)
	
	// Verify prompt contains expected sections
	assert.Contains(t, prompt, "## Repository Overview")
	assert.Contains(t, prompt, "Solution File: TestApp.sln")
	assert.Contains(t, prompt, "Project Files: src/TestApp.csproj, tests/TestApp.Tests.csproj")
	assert.Contains(t, prompt, "README Files: README.md")
	
	// Verify key files section
	assert.Contains(t, prompt, "## Key Files")
	assert.Contains(t, prompt, "### Solution Files")
	assert.Contains(t, prompt, "### Project Files")
	assert.Contains(t, prompt, "### Readme Files")
	assert.Contains(t, prompt, "### Source Files")
	
	// Verify file contents are included
	assert.Contains(t, prompt, "Solution file content")
	assert.Contains(t, prompt, "<Project Sdk=\"Microsoft.NET.Sdk.Web\"></Project>")
	assert.Contains(t, prompt, "# Test Application")
	assert.Contains(t, prompt, "class Program { }")
	
	// Verify analysis requirements section
	assert.Contains(t, prompt, "## Analysis Requirements")
	assert.Contains(t, prompt, "\"projectName\"")
	assert.Contains(t, prompt, "\"architecture\"")
	assert.Contains(t, prompt, "\"dependencies\"")
	assert.Contains(t, prompt, "\"technicalDebt\"")
	
	// Verify JSON structure is included
	assert.Contains(t, prompt, "```json")
	assert.Contains(t, prompt, "\"features\"")
	assert.Contains(t, prompt, "\"apis\"")
	assert.Contains(t, prompt, "\"security\"")
}

func TestPromptGenerator_EmptyRepository(t *testing.T) {
	// Test with minimal/empty scan result
	scanResult := &ScanResult{
		RootPath: "/empty/repo",
		Files:    []FileInfo{},
	}
	
	generator := NewPromptGenerator()
	prompt := generator.Generate(scanResult)
	
	// Should still generate a valid prompt structure
	assert.Contains(t, prompt, "## Repository Overview")
	assert.Contains(t, prompt, "## Key Files")
	assert.Contains(t, prompt, "## Analysis Requirements")
	
	// Should not have solution/project info
	assert.NotContains(t, prompt, "Solution File:")
}

func TestPromptGenerator_FileTypePrioritization(t *testing.T) {
	// Create scan result with mixed file types
	scanResult := &ScanResult{
		RootPath: "/test/repo",
		Files: []FileInfo{
			{RelPath: "src/Service.cs", Content: "service", FileType: "source"},
			{RelPath: "README.md", Content: "readme", FileType: "readme"},
			{RelPath: "Test.sln", Content: "solution", FileType: "solution"},
			{RelPath: "app.config", Content: "config", FileType: "config"},
			{RelPath: "Test.csproj", Content: "project", FileType: "project"},
		},
	}
	
	generator := NewPromptGenerator()
	prompt := generator.Generate(scanResult)
	
	// Find positions of each file type section in the prompt
	solutionPos := strings.Index(prompt, "### Solution Files")
	projectPos := strings.Index(prompt, "### Project Files")
	readmePos := strings.Index(prompt, "### Readme Files")
	configPos := strings.Index(prompt, "### Config Files")
	sourcePos := strings.Index(prompt, "### Source Files")
	
	// Verify files appear in priority order
	if solutionPos > 0 && projectPos > 0 {
		assert.Less(t, solutionPos, projectPos, "Solution should come before project")
	}
	if projectPos > 0 && readmePos > 0 {
		assert.Less(t, projectPos, readmePos, "Project should come before readme")
	}
	if readmePos > 0 && configPos > 0 {
		assert.Less(t, readmePos, configPos, "Readme should come before config")
	}
	if configPos > 0 && sourcePos > 0 {
		assert.Less(t, configPos, sourcePos, "Config should come before source")
	}
}