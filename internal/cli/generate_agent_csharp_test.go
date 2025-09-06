package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/karolswdev/docloom/internal/agent"
)

// TestCSharpAgent_E2E_Integration tests the full integration of the C# analyzer agent
func TestCSharpAgent_E2E_Integration(t *testing.T) {
	// Skip if not in CI or if explicitly disabled
	if os.Getenv("SKIP_E2E") == "true" {
		t.Skip("Skipping E2E test")
	}

	// Arrange - Create a sample C# project
	tmpDir := t.TempDir()

	// Create sample C# source files
	sampleProjectDir := filepath.Join(tmpDir, "sample-csharp-project")
	require.NoError(t, os.MkdirAll(sampleProjectDir, 0755))

	// Sample C# code with various constructs
	sampleCode1 := `using System;
using System.Collections.Generic;

namespace SampleProject.Services
{
	/// <summary>
	/// Service for managing products
	/// </summary>
	public class ProductService
	{
		private readonly IProductRepository _repository;
		
		public ProductService(IProductRepository repository)
		{
			_repository = repository;
		}
		
		/// <summary>
		/// Gets all products
		/// </summary>
		public async Task<List<Product>> GetAllProductsAsync()
		{
			return await _repository.GetAllAsync();
		}
		
		/// <summary>
		/// Finds a product by ID
		/// </summary>
		public Product FindById(int id)
		{
			return _repository.FindById(id);
		}
	}
	
	public interface IProductRepository
	{
		Task<List<Product>> GetAllAsync();
		Product FindById(int id);
	}
	
	public class Product
	{
		public int Id { get; set; }
		public string Name { get; set; }
		public decimal Price { get; set; }
	}
}`

	sampleCode2 := `using System;

namespace SampleProject.Controllers
{
	/// <summary>
	/// API controller for products
	/// </summary>
	public class ProductController
	{
		private readonly ProductService _service;
		
		public ProductController(ProductService service)
		{
			_service = service;
		}
		
		/// <summary>
		/// GET endpoint for all products
		/// </summary>
		public async Task<IActionResult> GetProducts()
		{
			var products = await _service.GetAllProductsAsync();
			return Ok(products);
		}
	}
}`

	// Write sample files
	require.NoError(t, os.WriteFile(
		filepath.Join(sampleProjectDir, "ProductService.cs"),
		[]byte(sampleCode1),
		0644,
	))

	require.NoError(t, os.WriteFile(
		filepath.Join(sampleProjectDir, "ProductController.cs"),
		[]byte(sampleCode2),
		0644,
	))

	// Create agent definition in temp location
	agentDir := filepath.Join(tmpDir, "agents")
	require.NoError(t, os.MkdirAll(agentDir, 0755))

	agentDef := agent.Definition{
		APIVersion: "v1",
		Kind:       "Agent",
		Metadata: agent.Metadata{
			Name:        "csharp-analyzer",
			Description: "Test C# analyzer",
		},
		Spec: agent.Spec{
			Runner: agent.Runner{
				Command: filepath.Join(".", "build", "docloom-agent-csharp"),
			},
			Parameters: []agent.Parameter{
				{
					Name:        "include-internal",
					Type:        "boolean",
					Description: "Include internal classes",
					Default:     false,
				},
			},
		},
	}

	agentDefPath := filepath.Join(agentDir, "csharp-analyzer.agent.yaml")
	agentDefData, err := yaml.Marshal(&agentDef)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(agentDefPath, agentDefData, 0644))

	// Build the agent binary if it doesn't exist
	agentBinaryPath := filepath.Join(".", "build", "docloom-agent-csharp")
	if _, statErr := os.Stat(agentBinaryPath); os.IsNotExist(statErr) {
		// Try to build it
		t.Logf("Building C# agent binary...")
		// Note: In a real CI environment, this would be built by the Makefile
		t.Skip("Agent binary not found, skipping E2E test")
	}

	// Act - Run the agent through the generate command workflow
	// Note: This is a simplified test that verifies the agent can be executed
	// In a full E2E test, we would run the complete generate command

	outputDir := filepath.Join(tmpDir, "output")
	require.NoError(t, os.MkdirAll(outputDir, 0755))

	// Create a registry and load the agent
	registry := agent.NewRegistry()
	registry.AddSearchPath(agentDir)
	require.NoError(t, registry.Discover())

	agents := registry.List()
	require.Len(t, agents, 1, "Should find one agent")
	assert.Equal(t, "csharp-analyzer", agents[0].Metadata.Name)

	// Execute the agent
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	cache, err := agent.NewArtifactCache()
	require.NoError(t, err)
	executor := agent.NewExecutor(registry, cache, logger)

	runOpts := agent.RunOptions{
		AgentName:  "csharp-analyzer",
		SourcePath: sampleProjectDir,
		Parameters: map[string]string{
			"include-internal": "false",
		},
	}

	result, err := executor.Run(runOpts)

	// If execution fails, it might be because the binary isn't built
	if err != nil {
		if strings.Contains(err.Error(), "executable file not found") ||
			strings.Contains(err.Error(), "no such file") {
			t.Skip("Agent binary not available, skipping execution test")
		}
		require.NoError(t, err, "Agent execution should succeed")
	}

	// Assert - Verify the agent produced expected artifacts
	require.NotNil(t, result, "Should have execution result")
	require.Equal(t, 0, result.ExitCode, "Agent should exit successfully")

	// List artifacts in the output directory
	artifacts, err := filepath.Glob(filepath.Join(result.OutputPath, "*.md"))
	require.NoError(t, err)
	require.NotEmpty(t, artifacts, "Should produce markdown artifacts")

	// Check for expected markdown files
	foundProjectSummary := false
	foundAPISurface := false
	foundInsights := false

	// Also check for JSON file
	jsonFile := filepath.Join(result.OutputPath, "analysis.json")

	for _, artifact := range artifacts {
		baseName := filepath.Base(artifact)
		switch baseName {
		case "ProjectSummary.md":
			foundProjectSummary = true
			// Read and verify content
			content, err := os.ReadFile(artifact)
			require.NoError(t, err)
			assert.Contains(t, string(content), "Project Summary")
			assert.Contains(t, string(content), "Total Classes")
		case "ApiSurface.md":
			foundAPISurface = true
			// Verify it contains our sample classes
			content, err := os.ReadFile(artifact)
			require.NoError(t, err)
			contentStr := string(content)
			assert.Contains(t, contentStr, "ProductService", "Should contain ProductService class")
			assert.Contains(t, contentStr, "ProductController", "Should contain ProductController class")
			assert.Contains(t, contentStr, "IProductRepository", "Should contain IProductRepository interface")
			assert.Contains(t, contentStr, "GetAllProductsAsync", "Should contain async method")
		case "ArchitecturalInsights.md":
			foundInsights = true
			// Verify architectural patterns
			content, err := os.ReadFile(artifact)
			require.NoError(t, err)
			contentStr := string(content)
			assert.Contains(t, contentStr, "Architectural Insights")
			assert.Contains(t, contentStr, "Repository Pattern", "Should detect repository pattern")
			assert.Contains(t, contentStr, "Asynchronous Programming", "Should detect async methods")
		}
	}

	// Check JSON file separately
	if _, err := os.Stat(jsonFile); err == nil {
		content, err := os.ReadFile(jsonFile)
		require.NoError(t, err)

		var analysis map[string]interface{}
		require.NoError(t, json.Unmarshal(content, &analysis))

		// Check for expected fields
		assert.Contains(t, analysis, "projectSummary")
		assert.Contains(t, analysis, "apiSurface")
		assert.Contains(t, analysis, "architecturalInsights")
	}

	assert.True(t, foundProjectSummary, "Should produce ProjectSummary.md")
	assert.True(t, foundAPISurface, "Should produce ApiSurface.md")
	assert.True(t, foundInsights, "Should produce ArchitecturalInsights.md")

	// Final verification: The complete workflow would render HTML
	// This is verified by checking that the artifacts can be used as sources
	t.Log("C# Analyzer Agent E2E test completed successfully")
}
