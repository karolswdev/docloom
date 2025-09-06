// Package ingest provides functionality for ingesting source documents.
package ingest

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

// Ingester handles the ingestion of source files.
type Ingester struct {
	// SupportedExtensions defines the file extensions that will be ingested.
	SupportedExtensions []string
}

// NewIngester creates a new Ingester with default supported extensions.
func NewIngester() *Ingester {
	return &Ingester{
		SupportedExtensions: []string{".md", ".txt"},
	}
}

// IngestSources recursively walks the provided paths and reads the content
// of all supported files into a single concatenated string.
func (i *Ingester) IngestSources(paths []string) (string, error) {
	var contentBuilder strings.Builder
	filesProcessed := 0

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			return "", fmt.Errorf("failed to stat path %s: %w", path, err)
		}

		if info.IsDir() {
			// Recursively walk the directory
			err = filepath.Walk(path, func(filePath string, fileInfo os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if fileInfo.IsDir() {
					return nil
				}

				if i.isSupportedFile(filePath) {
					content, err := i.readFile(filePath)
					if err != nil {
						log.Warn().Err(err).Str("file", filePath).Msg("Failed to read file, skipping")
						return nil // Continue processing other files
					}

					if contentBuilder.Len() > 0 {
						contentBuilder.WriteString("\n\n")
					}
					contentBuilder.WriteString(fmt.Sprintf("--- File: %s ---\n", filePath))
					contentBuilder.WriteString(content)
					filesProcessed++
					log.Debug().Str("file", filePath).Int("bytes", len(content)).Msg("Ingested file")
				}

				return nil
			})

			if err != nil {
				return "", fmt.Errorf("failed to walk directory %s: %w", path, err)
			}
		} else {
			// Single file
			if i.isSupportedFile(path) {
				content, err := i.readFile(path)
				if err != nil {
					return "", fmt.Errorf("failed to read file %s: %w", path, err)
				}

				if contentBuilder.Len() > 0 {
					contentBuilder.WriteString("\n\n")
				}
				contentBuilder.WriteString(fmt.Sprintf("--- File: %s ---\n", path))
				contentBuilder.WriteString(content)
				filesProcessed++
				log.Debug().Str("file", path).Int("bytes", len(content)).Msg("Ingested file")
			} else {
				log.Warn().Str("file", path).Msg("File type not supported for ingestion")
			}
		}
	}

	if filesProcessed == 0 {
		return "", fmt.Errorf("no supported files found in the provided paths")
	}

	log.Info().Int("files", filesProcessed).Int("total_bytes", contentBuilder.Len()).Msg("Ingestion complete")
	return contentBuilder.String(), nil
}

// isSupportedFile checks if a file has a supported extension.
func (i *Ingester) isSupportedFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, supportedExt := range i.SupportedExtensions {
		if ext == supportedExt {
			return true
		}
	}
	return false
}

// readFile reads the entire content of a file.
func (i *Ingester) readFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// AddSupportedExtension adds a new supported file extension.
func (i *Ingester) AddSupportedExtension(ext string) {
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	ext = strings.ToLower(ext)
	
	// Check if already exists
	for _, existing := range i.SupportedExtensions {
		if existing == ext {
			return
		}
	}
	
	i.SupportedExtensions = append(i.SupportedExtensions, ext)
}