package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScanner_Scan(t *testing.T) {
	// Create a temporary test repository
	tmpDir := t.TempDir()

	// Create test file structure
	files := map[string]string{
		"TestProject.sln": `Microsoft Visual Studio Solution File`,
		"src/TestProject.csproj": `<Project Sdk="Microsoft.NET.Sdk.Web">
  <PropertyGroup>
    <TargetFramework>net6.0</TargetFramework>
  </PropertyGroup>
</Project>`,
		"README.md": `# Test Project
This is a test project for scanner testing.`,
		"src/Program.cs": `using System;
namespace TestProject {
    class Program {
        static void Main() {
            Console.WriteLine("Hello World");
        }
    }
}`,
		"src/Controllers/TestController.cs": `using Microsoft.AspNetCore.Mvc;
namespace TestProject.Controllers {
    [ApiController]
    public class TestController : ControllerBase {
        [HttpGet("test")]
        public IActionResult Get() => Ok("test");
    }
}`,
		"Dockerfile": `FROM mcr.microsoft.com/dotnet/sdk:6.0
WORKDIR /app
COPY . .
RUN dotnet build`,
		"appsettings.json": `{
  "Logging": {
    "LogLevel": {
      "Default": "Information"
    }
  }
}`,
	}

	// Create test files
	for path, content := range files {
		fullPath := filepath.Join(tmpDir, path)
		dir := filepath.Dir(fullPath)
		require.NoError(t, os.MkdirAll(dir, 0755))
		require.NoError(t, os.WriteFile(fullPath, []byte(content), 0644))
	}

	// Create scanner and scan
	scanner := NewScanner(tmpDir)
	result, err := scanner.Scan()

	// Assertions
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, tmpDir, result.RootPath)

	// Check solution file was found
	assert.Equal(t, "TestProject.sln", result.SolutionFile)

	// Check project files were found
	assert.Contains(t, result.ProjectFiles, filepath.Join("src", "TestProject.csproj"))

	// Check README was found
	assert.Contains(t, result.ReadmeFiles, "README.md")

	// Check that files were collected
	assert.Greater(t, len(result.Files), 0)

	// Verify specific files were scanned
	var foundSolution, foundProject, foundReadme, foundProgram, foundController, foundConfig bool
	for _, file := range result.Files {
		switch file.RelPath {
		case "TestProject.sln":
			foundSolution = true
			assert.Equal(t, "solution", file.FileType)
		case filepath.Join("src", "TestProject.csproj"):
			foundProject = true
			assert.Equal(t, "project", file.FileType)
		case "README.md":
			foundReadme = true
			assert.Equal(t, "readme", file.FileType)
		case filepath.Join("src", "Program.cs"):
			foundProgram = true
			assert.Equal(t, "source", file.FileType)
		case filepath.Join("src", "Controllers", "TestController.cs"):
			foundController = true
			assert.Equal(t, "source", file.FileType)
		case "Dockerfile", "appsettings.json":
			foundConfig = true
		}
	}

	assert.True(t, foundSolution, "Solution file should be found")
	assert.True(t, foundProject, "Project file should be found")
	assert.True(t, foundReadme, "README file should be found")
	assert.True(t, foundProgram, "Program.cs should be found")
	assert.True(t, foundController, "Controller file should be found")
	assert.True(t, foundConfig, "Config files should be found")
}

func TestScanner_SkipsHiddenAndBuildDirectories(t *testing.T) {
	tmpDir := t.TempDir()

	// Create directories that should be skipped
	skipDirs := []string{
		".git/config",
		"bin/Debug/app.dll",
		"obj/Debug/app.pdb",
		"packages/Newtonsoft.Json/lib.dll",
		"node_modules/package/index.js",
	}

	for _, path := range skipDirs {
		fullPath := filepath.Join(tmpDir, path)
		dir := filepath.Dir(fullPath)
		require.NoError(t, os.MkdirAll(dir, 0755))
		require.NoError(t, os.WriteFile(fullPath, []byte("should not scan"), 0644))
	}

	// Create a file that should be scanned
	validFile := filepath.Join(tmpDir, "Program.cs")
	require.NoError(t, os.WriteFile(validFile, []byte("class Program {}"), 0644))

	scanner := NewScanner(tmpDir)
	result, err := scanner.Scan()

	require.NoError(t, err)

	// Verify skipped files are not in results
	for _, file := range result.Files {
		assert.NotContains(t, file.RelPath, ".git")
		assert.NotContains(t, file.RelPath, "bin")
		assert.NotContains(t, file.RelPath, "obj")
		assert.NotContains(t, file.RelPath, "packages")
		assert.NotContains(t, file.RelPath, "node_modules")
	}
}
