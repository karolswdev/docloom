package chunk

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestChunker_SimpleChunking tests the simple heuristic chunking functionality (TC-11.1).
func TestChunker_SimpleChunking(t *testing.T) {
	// Arrange: Provide a long text string that exceeds a small, predefined token limit
	longText := strings.Repeat("This is a sample sentence that contains multiple words. ", 100)
	
	// Set a small token limit (e.g., 50 tokens)
	maxTokens := 50
	chunker := NewChunker(maxTokens)
	
	// Act: Run the chunking and selection logic
	result := chunker.ChunkAndSelect(longText)
	
	// Assert: The returned string must be shorter than the original and within the token limit
	require.NotEmpty(t, result, "Result should not be empty")
	assert.Less(t, len(result), len(longText), "Result should be shorter than original")
	
	// Verify that the estimated tokens are within the limit
	estimatedTokens := chunker.EstimateTokens(result)
	assert.LessOrEqual(t, estimatedTokens, maxTokens+10, "Estimated tokens should be close to the limit (with small margin)")
	
	// Check that truncation indicator is added
	assert.Contains(t, result, "[Content truncated due to token limit]", "Should include truncation indicator")
}

// TestChunker_ContentWithinLimit tests that content within limit is returned as-is.
func TestChunker_ContentWithinLimit(t *testing.T) {
	// Arrange
	shortText := "This is a short text that fits within the token limit."
	maxTokens := 100 // Generous limit
	chunker := NewChunker(maxTokens)
	
	// Act
	result := chunker.ChunkAndSelect(shortText)
	
	// Assert
	assert.Equal(t, shortText, result, "Content within limit should be returned unchanged")
	assert.NotContains(t, result, "[Content truncated", "Should not include truncation indicator")
}

// TestChunker_EmptyContent tests handling of empty content.
func TestChunker_EmptyContent(t *testing.T) {
	// Arrange
	chunker := NewChunker(100)
	
	// Act
	result := chunker.ChunkAndSelect("")
	
	// Assert
	assert.Empty(t, result, "Empty input should return empty output")
}

// TestChunker_EstimateTokens tests token estimation.
func TestChunker_EstimateTokens(t *testing.T) {
	chunker := NewChunker(100)
	
	tests := []struct {
		name     string
		input    string
		minTokens int
		maxTokens int
	}{
		{
			name:     "Empty string",
			input:    "",
			minTokens: 0,
			maxTokens: 0,
		},
		{
			name:     "Single word",
			input:    "Hello",
			minTokens: 1,
			maxTokens: 3,
		},
		{
			name:     "Short sentence",
			input:    "This is a test sentence.",
			minTokens: 4,
			maxTokens: 10,
		},
		{
			name:     "Longer text",
			input:    "The quick brown fox jumps over the lazy dog. This is a pangram.",
			minTokens: 10,
			maxTokens: 25,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := chunker.EstimateTokens(tt.input)
			assert.GreaterOrEqual(t, tokens, tt.minTokens, "Token count should be at least %d", tt.minTokens)
			assert.LessOrEqual(t, tokens, tt.maxTokens, "Token count should be at most %d", tt.maxTokens)
		})
	}
}

// TestChunker_SmartTruncate tests intelligent truncation at boundaries.
func TestChunker_SmartTruncate(t *testing.T) {
	// Arrange
	content := `First paragraph with some content.

Second paragraph with more information. This is important.

Third paragraph that might get truncated. It contains additional details that may not fit.`
	
	chunker := NewChunker(30) // Small limit to force truncation
	
	// Act
	result := chunker.ChunkAndSelect(content)
	
	// Assert
	assert.NotEmpty(t, result)
	assert.Less(t, len(result), len(content), "Content should be truncated")
	// Should truncate at a reasonable boundary (paragraph or sentence)
	assert.NotContains(t, result, "Third paragraph", "Third paragraph should be truncated")
}

// TestChunker_ChunkByParagraphs tests paragraph-based chunking.
func TestChunker_ChunkByParagraphs(t *testing.T) {
	// Arrange
	content := `First paragraph with more text to exceed the token limit.

Second paragraph with additional content that needs to be chunked properly.

Third paragraph containing important information that should be preserved.

Fourth paragraph with concluding remarks and final thoughts.`
	
	chunker := NewChunker(20) // Small limit to create multiple chunks
	
	// Act
	chunks := chunker.ChunkByParagraphs(content)
	
	// Assert
	assert.GreaterOrEqual(t, len(chunks), 1, "Should create at least one chunk")
	
	// Each chunk should be within the token limit
	for i, chunk := range chunks {
		tokens := chunker.EstimateTokens(chunk)
		assert.LessOrEqual(t, tokens, chunker.MaxTokens+5, "Chunk %d should be within token limit (with margin)", i)
	}
	
	// Verify all paragraphs are included
	combined := strings.Join(chunks, "\n\n")
	assert.Contains(t, combined, "First paragraph")
	assert.Contains(t, combined, "Second paragraph")
	assert.Contains(t, combined, "Third paragraph")
	assert.Contains(t, combined, "Fourth paragraph")
}

// TestChunker_VeryLongSingleParagraph tests handling of very long single paragraphs.
func TestChunker_VeryLongSingleParagraph(t *testing.T) {
	// Arrange
	longParagraph := strings.Repeat("This is a very long sentence without paragraph breaks. ", 50)
	chunker := NewChunker(50)
	
	// Act
	result := chunker.ChunkAndSelect(longParagraph)
	
	// Assert
	assert.Less(t, len(result), len(longParagraph), "Should truncate long paragraph")
	assert.Contains(t, result, "[Content truncated due to token limit]")
	
	// Should break at a word boundary
	assert.NotEqual(t, result[len(result)-1], ' ', "Should not end with a space before truncation message")
}

// TestChunker_SpecialCharacters tests handling of special characters and unicode.
func TestChunker_SpecialCharacters(t *testing.T) {
	// Arrange
	content := "Special chars: & < > \" ' © ™ émojis: 🚀 🎉 中文字符"
	chunker := NewChunker(100)
	
	// Act
	result := chunker.ChunkAndSelect(content)
	
	// Assert
	assert.Equal(t, content, result, "Special characters should be preserved")
	
	// Test estimation doesn't crash with special chars
	tokens := chunker.EstimateTokens(content)
	assert.Greater(t, tokens, 0, "Should estimate tokens for special characters")
}