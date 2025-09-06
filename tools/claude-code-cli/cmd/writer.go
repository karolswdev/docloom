package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
	
	"github.com/rs/zerolog/log"
)

// ArtifactWriter writes analysis artifacts to disk
type ArtifactWriter struct {
	outputPath string
}

// NewArtifactWriter creates a new artifact writer
func NewArtifactWriter(outputPath string) *ArtifactWriter {
	return &ArtifactWriter{
		outputPath: outputPath,
	}
}

// Write writes the analysis response as artifacts according to the specification
func (w *ArtifactWriter) Write(response *AnalysisResponse) error {
	// Create output directory structure as per Phase 9 specification
	dirs := []string{
		filepath.Join(w.outputPath, "analysis"),
		filepath.Join(w.outputPath, "repository-context"),
		filepath.Join(w.outputPath, "technical-insights"),
	}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		log.Debug().Str("dir", dir).Msg("Created artifact directory")
	}
	
	// Write main analysis summary
	summaryPath := filepath.Join(w.outputPath, "analysis", "summary.json")
	if err := w.writeJSON(summaryPath, response); err != nil {
		return fmt.Errorf("failed to write summary: %w", err)
	}
	log.Debug().Str("file", summaryPath).Msg("Wrote analysis summary")
	
	// Write project overview
	overviewPath := filepath.Join(w.outputPath, "analysis", "project-overview.md")
	if err := w.writeProjectOverview(overviewPath, response); err != nil {
		return fmt.Errorf("failed to write project overview: %w", err)
	}
	log.Debug().Str("file", overviewPath).Msg("Wrote project overview")
	
	// Write architecture details
	archPath := filepath.Join(w.outputPath, "repository-context", "architecture.md")
	if err := w.writeArchitecture(archPath, response); err != nil {
		return fmt.Errorf("failed to write architecture: %w", err)
	}
	log.Debug().Str("file", archPath).Msg("Wrote architecture details")
	
	// Write API documentation
	if len(response.APIs) > 0 {
		apiPath := filepath.Join(w.outputPath, "repository-context", "api-endpoints.md")
		if err := w.writeAPIs(apiPath, response); err != nil {
			return fmt.Errorf("failed to write API documentation: %w", err)
		}
		log.Debug().Str("file", apiPath).Msg("Wrote API documentation")
	}
	
	// Write technical debt analysis
	if len(response.TechnicalDebt) > 0 {
		debtPath := filepath.Join(w.outputPath, "technical-insights", "technical-debt.md")
		if err := w.writeTechnicalDebt(debtPath, response); err != nil {
			return fmt.Errorf("failed to write technical debt: %w", err)
		}
		log.Debug().Str("file", debtPath).Msg("Wrote technical debt analysis")
	}
	
	// Write recommendations
	if len(response.Recommendations) > 0 {
		recPath := filepath.Join(w.outputPath, "technical-insights", "recommendations.md")
		if err := w.writeRecommendations(recPath, response); err != nil {
			return fmt.Errorf("failed to write recommendations: %w", err)
		}
		log.Debug().Str("file", recPath).Msg("Wrote recommendations")
	}
	
	// Write metadata file as per specification
	metadataPath := filepath.Join(w.outputPath, "metadata.json")
	timestamp := os.Getenv("SOURCE_DATE_EPOCH")
	if timestamp == "" {
		timestamp = fmt.Sprintf("%d", time.Now().Unix())
	}
	metadata := map[string]interface{}{
		"agent":      "csharp-cc-cli",
		"version":    "1.0.0",
		"timestamp":  timestamp,
		"model":      "claude-3-opus",
		"files_analyzed": len(response.Features),
	}
	if err := w.writeJSON(metadataPath, metadata); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}
	log.Debug().Str("file", metadataPath).Msg("Wrote metadata")
	
	return nil
}

func (w *ArtifactWriter) writeJSON(path string, data interface{}) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (w *ArtifactWriter) writeProjectOverview(path string, response *AnalysisResponse) error {
	content := fmt.Sprintf(`# Project Overview

## %s

%s

**Project Type:** %s  
**Framework:** %s

## Features

`, response.ProjectName, response.Description, response.ProjectType, response.Framework)
	
	for _, feature := range response.Features {
		content += fmt.Sprintf("- **%s**: %s\n", feature.Name, feature.Description)
	}
	
	return os.WriteFile(path, []byte(content), 0644)
}

func (w *ArtifactWriter) writeArchitecture(path string, response *AnalysisResponse) error {
	content := fmt.Sprintf(`# Architecture

## Pattern
%s

## Layers
`, response.Architecture.Pattern)
	
	for _, layer := range response.Architecture.Layers {
		content += fmt.Sprintf("- %s\n", layer)
	}
	
	content += "\n## Key Components\n"
	for _, component := range response.Architecture.KeyComponents {
		content += fmt.Sprintf("- %s\n", component)
	}
	
	content += fmt.Sprintf(`
## Dependencies

### NuGet Packages
`)
	for _, pkg := range response.Dependencies.NuGet {
		content += fmt.Sprintf("- %s\n", pkg)
	}
	
	if len(response.Dependencies.External) > 0 {
		content += "\n### External Services\n"
		for _, svc := range response.Dependencies.External {
			content += fmt.Sprintf("- %s\n", svc)
		}
	}
	
	return os.WriteFile(path, []byte(content), 0644)
}

func (w *ArtifactWriter) writeAPIs(path string, response *AnalysisResponse) error {
	content := "# API Endpoints\n\n"
	
	for _, api := range response.APIs {
		content += fmt.Sprintf("## %s %s\n%s\n\n", api.Method, api.Endpoint, api.Description)
	}
	
	return os.WriteFile(path, []byte(content), 0644)
}

func (w *ArtifactWriter) writeTechnicalDebt(path string, response *AnalysisResponse) error {
	content := "# Technical Debt Analysis\n\n"
	
	for _, debt := range response.TechnicalDebt {
		content += fmt.Sprintf(`## %s
**Impact:** %s

%s

**Recommendation:** %s

---

`, debt.Area, debt.Impact, debt.Description, debt.Recommendation)
	}
	
	return os.WriteFile(path, []byte(content), 0644)
}

func (w *ArtifactWriter) writeRecommendations(path string, response *AnalysisResponse) error {
	content := "# Recommendations\n\n"
	
	for i, rec := range response.Recommendations {
		content += fmt.Sprintf("%d. %s\n", i+1, rec)
	}
	
	return os.WriteFile(path, []byte(content), 0644)
}