// Package main implements the C# analyzer agent as a multi-tool executable.
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
	"github.com/spf13/cobra"
)

// AgentOutput represents the structured output from the agent
type AgentOutput struct {
	ProjectSummary        ProjectSummary        `json:"projectSummary"`
	APISurface            *parser.APISurface    `json:"apiSurface"`
	ArchitecturalInsights ArchitecturalInsights `json:"architecturalInsights"`
}

// ProjectSummary provides high-level project statistics
type ProjectSummary struct {
	TotalNamespaces int      `json:"totalNamespaces"`
	TotalClasses    int      `json:"totalClasses"`
	TotalInterfaces int      `json:"totalInterfaces"`
	TotalMethods    int      `json:"totalMethods"`
	TotalProperties int      `json:"totalProperties"`
	PublicAPIs      int      `json:"publicAPIs"`
	Projects        []string `json:"projects"`
}

// ArchitecturalInsights contains detected patterns and insights
type ArchitecturalInsights struct {
	DetectedPatterns []string `json:"detectedPatterns"`
	HasAsyncMethods  bool     `json:"hasAsyncMethods"`
	UsesInterfaces   bool     `json:"usesInterfaces"`
	HasAbstractions  bool     `json:"hasAbstractions"`
}

// FileInfo represents basic file information
type FileInfo struct {
	Path    string `json:"path"`
	Size    int64  `json:"size"`
	Content string `json:"content"`
}

var rootCmd = &cobra.Command{
	Use:   "docloom-agent-csharp",
	Short: "C# analyzer agent for DocLoom",
	Long:  `A multi-tool C# analyzer agent that provides various code analysis capabilities.`,
}

var summarizeReadmeCmd = &cobra.Command{
	Use:   "summarize_readme [path]",
	Short: "Summarizes README files in the repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sourcePath := args[0]
		
		// Find README files
		readmeFiles := []string{}
		filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			name := strings.ToLower(info.Name())
			if name == "readme.md" || name == "readme.txt" || name == "readme" {
				readmeFiles = append(readmeFiles, path)
			}
			return nil
		})

		output := map[string]interface{}{
			"readmeCount": len(readmeFiles),
			"readmePaths": readmeFiles,
			"summary":     fmt.Sprintf("Found %d README files in the repository", len(readmeFiles)),
		}

		json.NewEncoder(os.Stdout).Encode(output)
	},
}

var listProjectsCmd = &cobra.Command{
	Use:   "list_projects [path]",
	Short: "Lists all C# projects in the repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sourcePath := args[0]
		
		// Find .csproj files
		projects := []string{}
		filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if strings.HasSuffix(path, ".csproj") {
				relPath, _ := filepath.Rel(sourcePath, path)
				projects = append(projects, relPath)
			}
			return nil
		})

		output := map[string]interface{}{
			"projectCount": len(projects),
			"projects":     projects,
		}

		json.NewEncoder(os.Stdout).Encode(output)
	},
}

var getDependenciesCmd = &cobra.Command{
	Use:   "get_dependencies [path]",
	Short: "Analyzes project dependencies",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sourcePath := args[0]
		
		// Find and analyze .csproj files for dependencies
		dependencies := map[string][]string{}
		
		filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
			if err != nil || !strings.HasSuffix(path, ".csproj") {
				return nil
			}
			
			content, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			
			// Simple extraction of PackageReference elements
			deps := []string{}
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if strings.Contains(line, "<PackageReference") {
					// Extract Include attribute
					start := strings.Index(line, "Include=\"")
					if start > 0 {
						start += 9
						end := strings.Index(line[start:], "\"")
						if end > 0 {
							deps = append(deps, line[start:start+end])
						}
					}
				}
			}
			
			relPath, _ := filepath.Rel(sourcePath, path)
			dependencies[relPath] = deps
			return nil
		})

		output := map[string]interface{}{
			"dependencies": dependencies,
			"summary":      fmt.Sprintf("Analyzed %d projects", len(dependencies)),
		}

		json.NewEncoder(os.Stdout).Encode(output)
	},
}

var getAPISurfaceCmd = &cobra.Command{
	Use:   "get_api_surface [path]",
	Short: "Extracts the public API surface",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sourcePath := args[0]
		
		// Find all C# files
		csFiles, _ := findCSharpFiles(sourcePath)
		
		// Parse all files
		p := parser.New()
		var allAPIs parser.APISurface
		namespaceMap := make(map[string]*parser.Namespace)

		for _, file := range csFiles {
			content, err := os.ReadFile(file)
			if err != nil {
				continue
			}

			api, err := p.ExtractAPISurface(context.Background(), string(content))
			if err != nil {
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
			allAPIs.Namespaces = append(allAPIs.Namespaces, *ns)
		}

		// Generate summary
		summary := ProjectSummary{}
		for _, ns := range allAPIs.Namespaces {
			summary.TotalNamespaces++
			for _, class := range ns.Classes {
				if class.IsInterface {
					summary.TotalInterfaces++
				} else {
					summary.TotalClasses++
				}
				if class.IsPublic {
					summary.PublicAPIs++
				}
				summary.TotalMethods += len(class.Methods)
				summary.TotalProperties += len(class.Properties)
			}
		}

		output := map[string]interface{}{
			"summary":    summary,
			"apiSurface": allAPIs,
		}

		json.NewEncoder(os.Stdout).Encode(output)
	},
}

var getFileContentCmd = &cobra.Command{
	Use:   "get_file_content [path]",
	Short: "Gets the content of a specific file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		
		// Check if file exists
		info, err := os.Stat(filePath)
		if err != nil {
			output := map[string]interface{}{
				"error": fmt.Sprintf("File not found: %s", filePath),
			}
			json.NewEncoder(os.Stdout).Encode(output)
			return
		}

		// Read file content
		content, err := os.ReadFile(filePath)
		if err != nil {
			output := map[string]interface{}{
				"error": fmt.Sprintf("Error reading file: %v", err),
			}
			json.NewEncoder(os.Stdout).Encode(output)
			return
		}

		output := FileInfo{
			Path:    filePath,
			Size:    info.Size(),
			Content: string(content),
		}

		json.NewEncoder(os.Stdout).Encode(output)
	},
}

// Legacy mode for backward compatibility
var legacyCmd = &cobra.Command{
	Use:   "analyze [source_path] [output_path]",
	Short: "Legacy analysis mode (deprecated)",
	Args:  cobra.ExactArgs(2),
	Run:   runLegacyAnalysis,
}

func init() {
	rootCmd.AddCommand(summarizeReadmeCmd)
	rootCmd.AddCommand(listProjectsCmd)
	rootCmd.AddCommand(getDependenciesCmd)
	rootCmd.AddCommand(getAPISurfaceCmd)
	rootCmd.AddCommand(getFileContentCmd)
	rootCmd.AddCommand(legacyCmd)
}

func main() {
	// Check if running in legacy mode (for backward compatibility)
	if len(os.Args) >= 3 && !strings.Contains(os.Args[1], "_") {
		// Legacy mode: docloom-agent-csharp <source> <output>
		runLegacyAnalysis(nil, os.Args[1:3])
		return
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runLegacyAnalysis(cmd *cobra.Command, args []string) {
	sourcePath := args[0]
	outputPath := args[1]

	// Read parameters from environment
	includeInternal := parseBoolParam("PARAM_INCLUDE_INTERNAL", false)
	maxDepth := parseIntParam("PARAM_MAX_DEPTH", 10)
	extractMetrics := parseBoolParam("PARAM_EXTRACT_METRICS", true)

	fmt.Fprintf(os.Stderr, "C# Analyzer Agent starting (legacy mode)...\n")
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

	// Write output files
	writeProjectSummary(outputPath, &output.ProjectSummary)
	writeAPISurface(outputPath, &allAPIs)
	writeArchitecturalInsights(outputPath, &output.ArchitecturalInsights)

	// Write JSON output
	jsonPath := filepath.Join(outputPath, "analysis.json")
	jsonData, _ := json.MarshalIndent(output, "", "  ")
	os.WriteFile(jsonPath, jsonData, 0600)

	fmt.Fprintf(os.Stderr, "Analysis complete. Output written to %s\n", outputPath)
}

// Helper functions (same as before)
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