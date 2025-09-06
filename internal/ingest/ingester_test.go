package ingest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIngester_IngestSources tests the source ingestion functionality.
func TestIngester_IngestSources(t *testing.T) {
	// Arrange: Create a test directory with .md and .txt files, and a subdirectory
	tempDir := t.TempDir()

	// Create test files
	testFiles := map[string]string{
		"doc1.md":                  "# Document 1\n\nThis is markdown content.",
		"doc2.txt":                 "This is plain text content.",
		"subdir/doc3.md":          "# Document 3\n\nNested markdown content.",
		"subdir/nested/doc4.txt":  "Deeply nested text content.",
		"ignore.pdf":              "This should be ignored.",
		"subdir/ignore.docx":      "This should also be ignored.",
	}

	// Create the directory structure and files
	for path, content := range testFiles {
		fullPath := filepath.Join(tempDir, path)
		dir := filepath.Dir(fullPath)
		
		err := os.MkdirAll(dir, 0755)
		require.NoError(t, err)
		
		err = os.WriteFile(fullPath, []byte(content), 0644)
		require.NoError(t, err)
	}

	// Create ingester
	ingester := NewIngester()

	// Act: Call the ingest function on the parent directory
	result, err := ingester.IngestSources([]string{tempDir})

	// Assert: The function should return the concatenated content of all valid files
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// Check that all expected files are included
	assert.Contains(t, result, "Document 1")
	assert.Contains(t, result, "This is markdown content.")
	assert.Contains(t, result, "This is plain text content.")
	assert.Contains(t, result, "Document 3")
	assert.Contains(t, result, "Nested markdown content.")
	assert.Contains(t, result, "Deeply nested text content.")

	// Check that ignored files are not included
	assert.NotContains(t, result, "This should be ignored.")
	assert.NotContains(t, result, "This should also be ignored.")

	// Check that file paths are included as separators
	assert.Contains(t, result, "--- File:")
	assert.Contains(t, result, "doc1.md")
	assert.Contains(t, result, "doc2.txt")
	assert.Contains(t, result, "doc3.md")
	assert.Contains(t, result, "doc4.txt")

	// Verify the files are separated properly
	sections := strings.Split(result, "--- File:")
	// First element will be empty, so we expect 5 sections (4 files + 1 empty)
	assert.Equal(t, 5, len(sections))
}

// TestIngester_IngestSources_SingleFile tests ingesting a single file.
func TestIngester_IngestSources_SingleFile(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "single.md")
	content := "# Single Document\n\nThis is the content."
	
	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	ingester := NewIngester()

	// Act
	result, err := ingester.IngestSources([]string{testFile})

	// Assert
	require.NoError(t, err)
	assert.Contains(t, result, "Single Document")
	assert.Contains(t, result, "This is the content.")
	assert.Contains(t, result, "single.md")
}

// TestIngester_IngestSources_MultiplePaths tests ingesting from multiple paths.
func TestIngester_IngestSources_MultiplePaths(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	
	// Create files in different locations
	file1 := filepath.Join(tempDir, "file1.txt")
	err := os.WriteFile(file1, []byte("Content from file 1"), 0644)
	require.NoError(t, err)

	subdir := filepath.Join(tempDir, "subdir")
	err = os.MkdirAll(subdir, 0755)
	require.NoError(t, err)
	
	file2 := filepath.Join(subdir, "file2.md")
	err = os.WriteFile(file2, []byte("# Content from file 2"), 0644)
	require.NoError(t, err)

	ingester := NewIngester()

	// Act: Provide both the file and the directory as separate paths
	result, err := ingester.IngestSources([]string{file1, subdir})

	// Assert
	require.NoError(t, err)
	assert.Contains(t, result, "Content from file 1")
	assert.Contains(t, result, "Content from file 2")
}

// TestIngester_IngestSources_NoSupportedFiles tests behavior when no supported files are found.
func TestIngester_IngestSources_NoSupportedFiles(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	
	// Create only unsupported files
	unsupportedFile := filepath.Join(tempDir, "document.pdf")
	err := os.WriteFile(unsupportedFile, []byte("PDF content"), 0644)
	require.NoError(t, err)

	ingester := NewIngester()

	// Act
	result, err := ingester.IngestSources([]string{tempDir})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no supported files found")
	assert.Empty(t, result)
}

// TestIngester_IngestSources_NonExistentPath tests handling of non-existent paths.
func TestIngester_IngestSources_NonExistentPath(t *testing.T) {
	ingester := NewIngester()

	// Act
	result, err := ingester.IngestSources([]string{"/non/existent/path"})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to stat path")
	assert.Empty(t, result)
}

// TestIngester_IngestSources_EmptyFile tests handling of empty files.
func TestIngester_IngestSources_EmptyFile(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	emptyFile := filepath.Join(tempDir, "empty.md")
	
	err := os.WriteFile(emptyFile, []byte(""), 0644)
	require.NoError(t, err)

	ingester := NewIngester()

	// Act
	result, err := ingester.IngestSources([]string{emptyFile})

	// Assert
	require.NoError(t, err)
	assert.Contains(t, result, "empty.md")
	// The result should contain the file header and the empty content
	assert.Contains(t, result, "--- File:")
	// Even with empty content, the file is still processed
}

// TestIngester_AddSupportedExtension tests adding custom file extensions.
func TestIngester_AddSupportedExtension(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	customFile := filepath.Join(tempDir, "custom.log")
	err := os.WriteFile(customFile, []byte("Log file content"), 0644)
	require.NoError(t, err)

	ingester := NewIngester()

	// Initially, .log files should not be supported
	result, err := ingester.IngestSources([]string{customFile})
	assert.Error(t, err)

	// Act: Add .log as a supported extension
	ingester.AddSupportedExtension(".log")

	// Now it should work
	result, err = ingester.IngestSources([]string{customFile})

	// Assert
	require.NoError(t, err)
	assert.Contains(t, result, "Log file content")
	assert.Contains(t, result, "custom.log")
}

// TestIngester_AddSupportedExtension_NoDot tests adding extension without leading dot.
func TestIngester_AddSupportedExtension_NoDot(t *testing.T) {
	ingester := NewIngester()

	// Act: Add extension without dot
	ingester.AddSupportedExtension("log")

	// Assert: Should be added with dot
	found := false
	for _, ext := range ingester.SupportedExtensions {
		if ext == ".log" {
			found = true
			break
		}
	}
	assert.True(t, found, "Extension should be added with leading dot")
}

// TestIngester_AddSupportedExtension_Duplicate tests that duplicate extensions are not added.
func TestIngester_AddSupportedExtension_Duplicate(t *testing.T) {
	ingester := NewIngester()
	initialCount := len(ingester.SupportedExtensions)

	// Act: Add the same extension multiple times
	ingester.AddSupportedExtension(".md")  // Already exists
	ingester.AddSupportedExtension("md")   // Same, without dot
	ingester.AddSupportedExtension(".MD")  // Same, different case

	// Assert: Count should not change
	assert.Equal(t, initialCount, len(ingester.SupportedExtensions))
}