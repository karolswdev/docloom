// Package main implements the C# analyzer agent executable.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/karolswdev/docloom/internal/agents/csharp/parser"
)

// AgentOutput represents the structured output from the agent
type AgentOutput struct {
	ProjectSummary        ProjectSummary        `json:"projectSummary"`
	APISurface            *parser.APISurface    `json:"apiSurface"`
	ArchitecturalInsights ArchitecturalInsights `json:"architecturalInsights"`
}

// ProjectSummary provides high-level project statistics
type ProjectSummary struct {
	TotalNamespaces int `json:"totalNamespaces"`
	TotalClasses    int `json:"totalClasses"`
	TotalInterfaces int `json:"totalInterfaces"`
	TotalMethods    int `json:"totalMethods"`
	TotalProperties int `json:"totalProperties"`
	PublicAPIs      int `json:"publicAPIs"`
}

// ArchitecturalInsights contains detected patterns and insights
type ArchitecturalInsights struct {
	DetectedPatterns []string `json:"detectedPatterns"`
	HasAsyncMethods  bool     `json:"hasAsyncMethods"`
	UsesInterfaces   bool     `json:"usesInterfaces"`
	HasAbstractions  bool     `json:"hasAbstractions"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <source_path> <output_path>\n", os.Args[0])
		os.Exit(1)
	}

	sourcePath := os.Args[1]
	outputPath := os.Args[2]

	// Read parameters from environment
	includeInternal := parseBoolParam("PARAM_INCLUDE_INTERNAL", false)
	maxDepth := parseIntParam("PARAM_MAX_DEPTH", 10)
	extractMetrics := parseBoolParam("PARAM_EXTRACT_METRICS", true)

	fmt.Fprintf(os.Stderr, "C# Analyzer Agent starting...\n")
	fmt.Fprintf(os.Stderr, "Source: %s\n", sourcePath)
	fmt.Fprintf(os.Stderr, "Output: %s\n", outputPath)
	fmt.Fprintf(os.Stderr, "Parameters: includeInternal=%v, maxDepth=%d, extractMetrics=%v\n",
		includeInternal, maxDepth, extractMetrics)

	// Ensure output directory exists
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Find all C# files
	csFiles, err := findCSharpFiles(sourcePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding C# files: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Found %d C# files\n", len(csFiles))

	// Parse all files
	p := parser.New()
	var allAPIs parser.APISurface
	namespaceMap := make(map[string]*parser.Namespace)

	for _, file := range csFiles {
		fmt.Fprintf(os.Stderr, "Analyzing: %s\n", file)

		content, readErr := os.ReadFile(file)
		if readErr != nil {
			fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", file, readErr)
			continue
		}

		api, parseErr := p.ExtractAPISurface(context.Background(), string(content))
		if parseErr != nil {
			fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", file, parseErr)
			continue
		}

		// Merge namespaces
		for _, ns := range api.Namespaces {
			if existing, ok := namespaceMap[ns.Name]; ok {
				existing.Classes = append(existing.Classes, ns.Classes...)
			} else {
				nsCopy := ns
				namespaceMap[ns.Name] = &nsCopy
			}
		}
	}

	// Convert map back to slice
	for _, ns := range namespaceMap {
		if !includeInternal || shouldIncludeNamespace(ns) {
			allAPIs.Namespaces = append(allAPIs.Namespaces, *ns)
		}
	}

	// Generate output
	output := generateOutput(&allAPIs, extractMetrics)

	// Write ProjectSummary.md
	if writeErr := writeProjectSummary(outputPath, &output.ProjectSummary); writeErr != nil {
		fmt.Fprintf(os.Stderr, "Error writing ProjectSummary.md: %v\n", writeErr)
		os.Exit(1)
	}

	// Write ApiSurface.md
	if writeErr := writeAPISurface(outputPath, &allAPIs); writeErr != nil {
		fmt.Fprintf(os.Stderr, "Error writing ApiSurface.md: %v\n", writeErr)
		os.Exit(1)
	}

	// Write ArchitecturalInsights.md
	if writeErr := writeArchitecturalInsights(outputPath, &output.ArchitecturalInsights); writeErr != nil {
		fmt.Fprintf(os.Stderr, "Error writing ArchitecturalInsights.md: %v\n", writeErr)
		os.Exit(1)
	}

	// Write JSON output for further processing
	jsonPath := filepath.Join(outputPath, "analysis.json")
	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(jsonPath, jsonData, 0600); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Analysis complete. Output written to %s\n", outputPath)
}

func findCSharpFiles(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip common non-source directories
		if info.IsDir() {
			name := info.Name()
			if name == "bin" || name == "obj" || name == ".git" || name == "packages" {
				return filepath.SkipDir
			}
		}

		if !info.IsDir() && strings.HasSuffix(path, ".cs") {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func generateOutput(api *parser.APISurface, _ bool) *AgentOutput {
	output := &AgentOutput{
		APISurface: api,
	}

	// Calculate summary statistics
	for _, ns := range api.Namespaces {
		output.ProjectSummary.TotalNamespaces++

		for _, class := range ns.Classes {
			if class.IsInterface {
				output.ProjectSummary.TotalInterfaces++
				output.ArchitecturalInsights.UsesInterfaces = true
			} else {
				output.ProjectSummary.TotalClasses++
			}

			if class.IsAbstract {
				output.ArchitecturalInsights.HasAbstractions = true
			}

			if class.IsPublic {
				output.ProjectSummary.PublicAPIs++
			}

			output.ProjectSummary.TotalMethods += len(class.Methods)
			output.ProjectSummary.TotalProperties += len(class.Properties)

			// Check for async methods
			for _, method := range class.Methods {
				if strings.Contains(method.ReturnType, "Task") ||
					strings.Contains(method.ReturnType, "async") {
					output.ArchitecturalInsights.HasAsyncMethods = true
				}
			}
		}
	}

	// Detect architectural patterns
	output.ArchitecturalInsights.DetectedPatterns = detectPatterns(api)

	return output
}

func detectPatterns(api *parser.APISurface) []string {
	patterns := []string{}

	hasRepository := false
	hasFactory := false
	hasService := false
	hasController := false

	for _, ns := range api.Namespaces {
		for _, class := range ns.Classes {
			className := strings.ToLower(class.Name)

			if strings.Contains(className, "repository") {
				hasRepository = true
			}
			if strings.Contains(className, "factory") {
				hasFactory = true
			}
			if strings.Contains(className, "service") {
				hasService = true
			}
			if strings.Contains(className, "controller") {
				hasController = true
			}
		}
	}

	if hasRepository {
		patterns = append(patterns, "Repository Pattern")
	}
	if hasFactory {
		patterns = append(patterns, "Factory Pattern")
	}
	if hasService {
		patterns = append(patterns, "Service Layer")
	}
	if hasController {
		patterns = append(patterns, "MVC/API Controllers")
	}

	return patterns
}

func shouldIncludeNamespace(ns *parser.Namespace) bool {
	// Check if namespace has any public types
	for _, class := range ns.Classes {
		if class.IsPublic {
			return true
		}
	}
	return false
}

func writeProjectSummary(outputPath string, summary *ProjectSummary) error {
	content := fmt.Sprintf(`# Project Summary

## Statistics

- **Total Namespaces**: %d
- **Total Classes**: %d
- **Total Interfaces**: %d
- **Total Methods**: %d
- **Total Properties**: %d
- **Public APIs**: %d

## Overview

This C# project contains %d namespaces with %d classes and %d interfaces. 
The codebase exposes %d public APIs with a total of %d methods and %d properties.
`,
		summary.TotalNamespaces,
		summary.TotalClasses,
		summary.TotalInterfaces,
		summary.TotalMethods,
		summary.TotalProperties,
		summary.PublicAPIs,
		summary.TotalNamespaces,
		summary.TotalClasses,
		summary.TotalInterfaces,
		summary.PublicAPIs,
		summary.TotalMethods,
		summary.TotalProperties,
	)

	path := filepath.Join(outputPath, "ProjectSummary.md")
	return os.WriteFile(path, []byte(content), 0600)
}

func writeAPISurface(outputPath string, api *parser.APISurface) error {
	var sb strings.Builder

	sb.WriteString("# API Surface\n\n")
	sb.WriteString("## Namespaces\n\n")

	for _, ns := range api.Namespaces {
		if ns.Name == "<global>" {
			sb.WriteString("### Global Namespace\n\n")
		} else {
			sb.WriteString(fmt.Sprintf("### %s\n\n", ns.Name))
		}

		for _, class := range ns.Classes {
			visibility := "internal"
			if class.IsPublic {
				visibility = "public"
			}

			typeKind := "class"
			if class.IsInterface {
				typeKind = "interface"
			} else if class.IsAbstract {
				typeKind = "abstract class"
			}

			sb.WriteString(fmt.Sprintf("#### %s %s %s\n\n", visibility, typeKind, class.Name))

			if class.DocComment != "" {
				sb.WriteString(fmt.Sprintf("_%s_\n\n", class.DocComment))
			}

			if len(class.Properties) > 0 {
				sb.WriteString("**Properties:**\n\n")
				for _, prop := range class.Properties {
					if prop.IsPublic {
						static := ""
						if prop.IsStatic {
							static = "static "
						}
						sb.WriteString(fmt.Sprintf("- %s`%s %s`", static, prop.Type, prop.Name))
						if prop.DocComment != "" {
							sb.WriteString(fmt.Sprintf(" - %s", prop.DocComment))
						}
						sb.WriteString("\n")
					}
				}
				sb.WriteString("\n")
			}

			if len(class.Methods) > 0 {
				sb.WriteString("**Methods:**\n\n")
				for _, method := range class.Methods {
					if method.IsPublic {
						static := ""
						if method.IsStatic {
							static = "static "
						}
						sb.WriteString(fmt.Sprintf("- %s`%s`", static, method.Signature))
						if method.DocComment != "" {
							sb.WriteString(fmt.Sprintf(" - %s", method.DocComment))
						}
						sb.WriteString("\n")
					}
				}
				sb.WriteString("\n")
			}
		}
	}

	path := filepath.Join(outputPath, "ApiSurface.md")
	return os.WriteFile(path, []byte(sb.String()), 0600)
}

func writeArchitecturalInsights(outputPath string, insights *ArchitecturalInsights) error {
	var sb strings.Builder

	sb.WriteString("# Architectural Insights\n\n")

	sb.WriteString("## Detected Patterns\n\n")
	if len(insights.DetectedPatterns) > 0 {
		for _, pattern := range insights.DetectedPatterns {
			sb.WriteString(fmt.Sprintf("- %s\n", pattern))
		}
	} else {
		sb.WriteString("No common architectural patterns detected.\n")
	}
	sb.WriteString("\n")

	sb.WriteString("## Design Characteristics\n\n")

	if insights.HasAsyncMethods {
		sb.WriteString("- **Asynchronous Programming**: The codebase uses async/await patterns\n")
	}

	if insights.UsesInterfaces {
		sb.WriteString("- **Interface-Based Design**: Interfaces are used for abstraction\n")
	}

	if insights.HasAbstractions {
		sb.WriteString("- **Abstract Classes**: Abstract base classes provide shared functionality\n")
	}

	if !insights.HasAsyncMethods && !insights.UsesInterfaces && !insights.HasAbstractions {
		sb.WriteString("The codebase follows a straightforward implementation approach.\n")
	}

	sb.WriteString("\n## Recommendations\n\n")
	sb.WriteString("Based on the analysis:\n\n")

	if !insights.UsesInterfaces {
		sb.WriteString("- Consider introducing interfaces to improve testability and flexibility\n")
	}

	if !insights.HasAsyncMethods {
		sb.WriteString("- Consider using async/await for I/O-bound operations\n")
	}

	if len(insights.DetectedPatterns) == 0 {
		sb.WriteString("- Consider implementing common design patterns where appropriate\n")
	}

	path := filepath.Join(outputPath, "ArchitecturalInsights.md")
	return os.WriteFile(path, []byte(sb.String()), 0600)
}

func parseBoolParam(name string, defaultValue bool) bool {
	val := os.Getenv(name)
	if val == "" {
		return defaultValue
	}
	return strings.ToLower(val) == "true"
}

func parseIntParam(name string, defaultValue int) int {
	val := os.Getenv(name)
	if val == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return i
}
