package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

// ScanResult contains the results of scanning a repository
type ScanResult struct {
	RootPath     string
	Files        []FileInfo
	SolutionFile string
	ProjectFiles []string
	ReadmeFiles  []string
}

// FileInfo contains information about a scanned file
type FileInfo struct {
	Path     string
	RelPath  string
	Content  string
	FileType string
}

// Scanner scans a C# repository for key files
type Scanner struct {
	rootPath string
}

// NewScanner creates a new repository scanner
func NewScanner(rootPath string) *Scanner {
	return &Scanner{
		rootPath: rootPath,
	}
}

// Scan performs the repository scan
func (s *Scanner) Scan() (*ScanResult, error) {
	result := &ScanResult{
		RootPath: s.rootPath,
		Files:    []FileInfo{},
	}
	
	// Key file patterns to look for (currently using direct checks in the walk function)
	// Commented out to avoid unused variable error - patterns are checked inline below
	// patterns := map[string][]string{
	// 	"solution": {"*.sln"},
	// 	"project":  {"*.csproj", "*.fsproj", "*.vbproj"},
	// 	"readme":   {"README.md", "README.txt", "README"},
	// 	"config":   {"appsettings.json", "appsettings.*.json", "web.config", "app.config"},
	// 	"docker":   {"Dockerfile", "docker-compose.yml", "docker-compose.yaml"},
	// 	"ci":       {".github/workflows/*.yml", ".gitlab-ci.yml", "azure-pipelines.yml"},
	// }
	
	// Walk the repository
	err := filepath.Walk(s.rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue walking despite errors
		}
		
		// Skip hidden directories and common non-source directories
		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") && name != "." && name != ".github" {
				return filepath.SkipDir
			}
			if name == "bin" || name == "obj" || name == "packages" || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		
		relPath, _ := filepath.Rel(s.rootPath, path)
		fileName := filepath.Base(path)
		
		// Check for solution files
		if strings.HasSuffix(fileName, ".sln") {
			result.SolutionFile = relPath
			content, _ := s.readFile(path)
			result.Files = append(result.Files, FileInfo{
				Path:     path,
				RelPath:  relPath,
				Content:  content,
				FileType: "solution",
			})
			log.Debug().Str("file", relPath).Msg("Found solution file")
		}
		
		// Check for project files
		if strings.HasSuffix(fileName, ".csproj") || strings.HasSuffix(fileName, ".fsproj") {
			result.ProjectFiles = append(result.ProjectFiles, relPath)
			content, _ := s.readFile(path)
			result.Files = append(result.Files, FileInfo{
				Path:     path,
				RelPath:  relPath,
				Content:  content,
				FileType: "project",
			})
			log.Debug().Str("file", relPath).Msg("Found project file")
		}
		
		// Check for README files
		if strings.HasPrefix(strings.ToUpper(fileName), "README") {
			result.ReadmeFiles = append(result.ReadmeFiles, relPath)
			content, _ := s.readFile(path)
			result.Files = append(result.Files, FileInfo{
				Path:     path,
				RelPath:  relPath,
				Content:  content,
				FileType: "readme",
			})
			log.Debug().Str("file", relPath).Msg("Found README file")
		}
		
		// Check for important config files
		if fileName == "appsettings.json" || fileName == "Dockerfile" || fileName == "docker-compose.yml" {
			content, _ := s.readFile(path)
			result.Files = append(result.Files, FileInfo{
				Path:     path,
				RelPath:  relPath,
				Content:  content,
				FileType: "config",
			})
			log.Debug().Str("file", relPath).Msg("Found config file")
		}
		
		// Sample some C# source files (limit to avoid token overflow)
		if strings.HasSuffix(fileName, ".cs") && len(result.Files) < 50 {
			// Only include key source files
			if strings.Contains(fileName, "Program.cs") || 
			   strings.Contains(fileName, "Startup.cs") ||
			   strings.Contains(fileName, "Controller") ||
			   strings.Contains(fileName, "Service") ||
			   strings.Contains(path, "/Models/") ||
			   strings.Contains(path, "/Interfaces/") {
				content, _ := s.readFile(path)
				result.Files = append(result.Files, FileInfo{
					Path:     path,
					RelPath:  relPath,
					Content:  content,
					FileType: "source",
				})
				log.Debug().Str("file", relPath).Msg("Found key source file")
			}
		}
		
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to walk repository: %w", err)
	}
	
	return result, nil
}

func (s *Scanner) readFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	
	// Limit file size to prevent memory issues
	const maxSize = 1024 * 1024 // 1MB
	limited := io.LimitReader(file, maxSize)
	
	content, err := io.ReadAll(limited)
	if err != nil {
		return "", err
	}
	
	return string(content), nil
}