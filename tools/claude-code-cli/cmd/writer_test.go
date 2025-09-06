package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArtifactWriter_Write(t *testing.T) {
	// Create temporary output directory
	tmpDir := t.TempDir()

	// Create test response
	response := &AnalysisResponse{
		ProjectName: "TestProject",
		Description: "A comprehensive test project",
		ProjectType: "Web API",
		Framework:   ".NET 6",
		Architecture: Architecture{
			Pattern:       "Clean Architecture",
			Layers:        []string{"API", "Application", "Domain"},
			KeyComponents: []string{"Controllers", "Services"},
		},
		Dependencies: Dependencies{
			NuGet:    []string{"Microsoft.AspNetCore"},
			External: []string{"PostgreSQL"},
		},
		Features: []Feature{
			{Name: "User Management", Description: "CRUD for users"},
		},
		APIs: []API{
			{Endpoint: "/api/users", Method: "GET", Description: "Get users"},
		},
		DataModel: DataModel{
			Entities: []string{"User", "Role"},
			Database: "PostgreSQL",
		},
		Testing: Testing{
			Framework: "xUnit",
			Coverage:  "80%",
			Types:     []string{"Unit"},
		},
		Deployment: Deployment{
			Containerized: true,
			CICD:          "GitHub Actions",
			Hosting:       "Docker",
		},
		Security: Security{
			Authentication: "JWT",
			Authorization:  "RBAC",
			Considerations: []string{"HTTPS"},
		},
		TechnicalDebt: []TechnicalDebtItem{
			{
				Area:           "Testing",
				Description:    "Need more integration tests",
				Impact:         "Medium",
				Recommendation: "Add integration test suite",
			},
		},
		Recommendations: []string{
			"Add caching layer",
		},
	}

	// Write artifacts
	writer := NewArtifactWriter(tmpDir)
	err := writer.Write(response)
	require.NoError(t, err)

	// Verify directory structure
	assert.DirExists(t, filepath.Join(tmpDir, "analysis"))
	assert.DirExists(t, filepath.Join(tmpDir, "repository-context"))
	assert.DirExists(t, filepath.Join(tmpDir, "technical-insights"))

	// Verify files exist
	assert.FileExists(t, filepath.Join(tmpDir, "metadata.json"))
	assert.FileExists(t, filepath.Join(tmpDir, "analysis", "summary.json"))
	assert.FileExists(t, filepath.Join(tmpDir, "analysis", "project-overview.md"))
	assert.FileExists(t, filepath.Join(tmpDir, "repository-context", "architecture.md"))
	assert.FileExists(t, filepath.Join(tmpDir, "repository-context", "api-endpoints.md"))
	assert.FileExists(t, filepath.Join(tmpDir, "technical-insights", "technical-debt.md"))
	assert.FileExists(t, filepath.Join(tmpDir, "technical-insights", "recommendations.md"))

	// Verify JSON content
	summaryPath := filepath.Join(tmpDir, "analysis", "summary.json")
	summaryData, err := os.ReadFile(summaryPath)
	require.NoError(t, err)

	var loadedResponse AnalysisResponse
	err = json.Unmarshal(summaryData, &loadedResponse)
	require.NoError(t, err)
	assert.Equal(t, response.ProjectName, loadedResponse.ProjectName)

	// Verify metadata
	metadataPath := filepath.Join(tmpDir, "metadata.json")
	metadataData, err := os.ReadFile(metadataPath)
	require.NoError(t, err)

	var metadata map[string]interface{}
	err = json.Unmarshal(metadataData, &metadata)
	require.NoError(t, err)
	assert.Equal(t, "csharp-cc-cli", metadata["agent"])
	assert.Equal(t, "1.0.0", metadata["version"])

	// Verify markdown content
	overviewPath := filepath.Join(tmpDir, "analysis", "project-overview.md")
	overviewContent, err := os.ReadFile(overviewPath)
	require.NoError(t, err)
	assert.Contains(t, string(overviewContent), "TestProject")
	assert.Contains(t, string(overviewContent), "Web API")
	assert.Contains(t, string(overviewContent), ".NET 6")
	assert.Contains(t, string(overviewContent), "User Management")
}

func TestArtifactWriter_EmptyOptionalSections(t *testing.T) {
	// Test with minimal response (no APIs, technical debt, or recommendations)
	tmpDir := t.TempDir()

	response := &AnalysisResponse{
		ProjectName: "MinimalProject",
		Description: "Minimal test",
		ProjectType: "Library",
		Framework:   ".NET Standard",
		Architecture: Architecture{
			Pattern: "Simple",
			Layers:  []string{"Core"},
		},
		Dependencies: Dependencies{
			NuGet: []string{},
		},
		Features: []Feature{},
		// Empty APIs, TechnicalDebt, and Recommendations
	}

	writer := NewArtifactWriter(tmpDir)
	err := writer.Write(response)
	require.NoError(t, err)

	// Verify required files exist
	assert.FileExists(t, filepath.Join(tmpDir, "metadata.json"))
	assert.FileExists(t, filepath.Join(tmpDir, "analysis", "summary.json"))
	assert.FileExists(t, filepath.Join(tmpDir, "analysis", "project-overview.md"))
	assert.FileExists(t, filepath.Join(tmpDir, "repository-context", "architecture.md"))

	// Verify optional files don't exist when data is empty
	assert.NoFileExists(t, filepath.Join(tmpDir, "repository-context", "api-endpoints.md"))
	assert.NoFileExists(t, filepath.Join(tmpDir, "technical-insights", "technical-debt.md"))
	assert.NoFileExists(t, filepath.Join(tmpDir, "technical-insights", "recommendations.md"))
}
