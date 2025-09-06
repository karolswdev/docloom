package cmd

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalysisResponse_JSONMarshaling(t *testing.T) {
	// Test that our response structure can properly marshal/unmarshal JSON
	response := &AnalysisResponse{
		ProjectName: "TestProject",
		Description: "A test project for validation",
		ProjectType: "Web API",
		Framework:   ".NET 6",
		Architecture: Architecture{
			Pattern:       "Clean Architecture",
			Layers:        []string{"API", "Application", "Domain", "Infrastructure"},
			KeyComponents: []string{"Controllers", "Services", "Repositories"},
		},
		Dependencies: Dependencies{
			NuGet:    []string{"Microsoft.AspNetCore", "EntityFramework.Core"},
			External: []string{"Redis", "PostgreSQL"},
		},
		Features: []Feature{
			{Name: "User Management", Description: "CRUD operations for users"},
			{Name: "Authentication", Description: "JWT-based authentication"},
		},
		APIs: []API{
			{Endpoint: "/api/users", Method: "GET", Description: "Get all users"},
			{Endpoint: "/api/users/{id}", Method: "GET", Description: "Get user by ID"},
		},
		DataModel: DataModel{
			Entities: []string{"User", "Role", "Permission"},
			Database: "PostgreSQL",
		},
		Testing: Testing{
			Framework: "xUnit",
			Coverage:  "85%",
			Types:     []string{"Unit", "Integration"},
		},
		Deployment: Deployment{
			Containerized: true,
			CICD:          "GitHub Actions",
			Hosting:       "Kubernetes",
		},
		Security: Security{
			Authentication: "JWT",
			Authorization:  "Role-based",
			Considerations: []string{"HTTPS only", "Rate limiting"},
		},
		TechnicalDebt: []TechnicalDebtItem{
			{
				Area:           "Database",
				Description:    "Missing indexes on frequently queried columns",
				Impact:         "Medium",
				Recommendation: "Add indexes to improve query performance",
			},
		},
		Recommendations: []string{
			"Implement caching for frequently accessed data",
			"Add comprehensive logging",
		},
	}

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(response, "", "  ")
	require.NoError(t, err)

	// Unmarshal back
	var decoded AnalysisResponse
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)

	// Verify key fields
	assert.Equal(t, response.ProjectName, decoded.ProjectName)
	assert.Equal(t, response.ProjectType, decoded.ProjectType)
	assert.Equal(t, response.Framework, decoded.Framework)
	assert.Equal(t, len(response.Features), len(decoded.Features))
	assert.Equal(t, len(response.APIs), len(decoded.APIs))
	assert.Equal(t, response.Architecture.Pattern, decoded.Architecture.Pattern)
	assert.Equal(t, len(response.Architecture.Layers), len(decoded.Architecture.Layers))
	assert.Equal(t, response.Deployment.Containerized, decoded.Deployment.Containerized)
	assert.Equal(t, len(response.TechnicalDebt), len(decoded.TechnicalDebt))
}

func TestClaudeClient_ParseJSONFromMarkdown(t *testing.T) {
	// Test that we can extract JSON from markdown-wrapped responses
	markdownResponse := "Here is the analysis of your C# repository:\n\n```json\n{\n  \"projectName\": \"TestAPI\",\n  \"description\": \"Test API project\",\n  \"projectType\": \"Web API\",\n  \"framework\": \".NET 6\"\n}\n```\n\nThe analysis is complete."

	// Extract JSON
	jsonStart := strings.Index(markdownResponse, "{")
	jsonEnd := strings.LastIndex(markdownResponse, "}")
	require.True(t, jsonStart >= 0)
	require.True(t, jsonEnd > jsonStart)

	jsonContent := markdownResponse[jsonStart : jsonEnd+1]

	// Parse JSON
	var analysis AnalysisResponse
	err := json.Unmarshal([]byte(jsonContent), &analysis)
	require.NoError(t, err)

	assert.Equal(t, "TestAPI", analysis.ProjectName)
	assert.Equal(t, "Test API project", analysis.Description)
	assert.Equal(t, "Web API", analysis.ProjectType)
	assert.Equal(t, ".NET 6", analysis.Framework)
}

func TestClaudeClient_NewClient(t *testing.T) {
	client := NewClaudeClient("test-key", "claude-3-opus", 4096)

	assert.NotNil(t, client)
	assert.Equal(t, "test-key", client.apiKey)
	assert.Equal(t, "claude-3-opus", client.model)
	assert.Equal(t, 4096, client.maxTokens)
}
