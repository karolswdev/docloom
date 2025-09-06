// Package chunk provides functionality for chunking and selecting content.
package chunk

import (
	"strings"
	"unicode"

	"github.com/rs/zerolog/log"
)

// Chunker handles content chunking and selection based on token limits.
type Chunker struct {
	// MaxTokens defines the maximum number of tokens allowed in the output.
	MaxTokens int
	// TokensPerChar is an approximation of tokens per character (default: 0.25 = ~4 chars per token).
	TokensPerChar float64
}

// NewChunker creates a new Chunker with default settings.
func NewChunker(maxTokens int) *Chunker {
	return &Chunker{
		MaxTokens:     maxTokens,
		TokensPerChar: 0.25, // Approximation: 1 token ≈ 4 characters
	}
}

// ChunkAndSelect takes input text and returns a truncated version that fits within token limits.
// This implements a simple heuristic: estimate token count and truncate if necessary.
func (c *Chunker) ChunkAndSelect(content string) string {
	if content == "" {
		return ""
	}

	// Estimate current token count
	estimatedTokens := c.EstimateTokens(content)

	log.Debug().
		Int("estimated_tokens", estimatedTokens).
		Int("max_tokens", c.MaxTokens).
		Int("content_length", len(content)).
		Msg("Processing content for chunking")

	// If content fits within limit, return as-is
	if estimatedTokens <= c.MaxTokens {
		log.Debug().Msg("Content fits within token limit, returning as-is")
		return content
	}

	// Calculate approximate character limit based on token limit
	maxChars := int(float64(c.MaxTokens) / c.TokensPerChar)

	// Apply smart truncation
	truncated := c.smartTruncate(content, maxChars)

	log.Info().
		Int("original_length", len(content)).
		Int("truncated_length", len(truncated)).
		Int("original_tokens", estimatedTokens).
		Int("truncated_tokens", c.EstimateTokens(truncated)).
		Msg("Content truncated to fit token limit")

	return truncated
}

// EstimateTokens estimates the number of tokens in the given text.
func (c *Chunker) EstimateTokens(text string) int {
	if text == "" {
		return 0
	}

	// Simple heuristic: count words and punctuation as rough token estimate
	// This is a simplified version; real tokenizers are more complex

	wordCount := 0
	inWord := false

	for _, r := range text {
		if unicode.IsSpace(r) {
			if inWord {
				wordCount++
				inWord = false
			}
		} else {
			inWord = true
		}
	}

	if inWord {
		wordCount++
	}

	// Alternative estimation using character count
	charEstimate := int(float64(len(text)) * c.TokensPerChar)

	// Use the average of word count and character-based estimate
	// This provides a more balanced estimation
	estimate := (wordCount + charEstimate) / 2

	return estimate
}

// smartTruncate truncates content intelligently at a reasonable boundary.
func (c *Chunker) smartTruncate(content string, maxChars int) string {
	if len(content) <= maxChars {
		return content
	}

	// First, try to truncate at a paragraph boundary
	truncated := content[:maxChars]

	// Look for the last paragraph break within the limit
	lastParagraph := strings.LastIndex(truncated, "\n\n")
	if lastParagraph > maxChars/2 { // Only use if we keep at least half the content
		return strings.TrimSpace(content[:lastParagraph]) + "\n\n[Content truncated due to token limit]"
	}

	// Otherwise, try to break at the last sentence
	lastSentence := -1
	sentenceEnds := []string{". ", "! ", "? ", ".\n", "!\n", "?\n"}

	for _, end := range sentenceEnds {
		idx := strings.LastIndex(truncated, end)
		if idx > lastSentence {
			lastSentence = idx
		}
	}

	if lastSentence > maxChars/2 {
		// Include the sentence-ending punctuation
		return strings.TrimSpace(content[:lastSentence+1]) + "\n\n[Content truncated due to token limit]"
	}

	// Last resort: break at the last word boundary
	lastSpace := strings.LastIndexFunc(truncated, unicode.IsSpace)
	if lastSpace > 0 {
		return strings.TrimSpace(content[:lastSpace]) + "...\n\n[Content truncated due to token limit]"
	}

	// Absolute last resort: hard truncate
	return truncated + "...\n\n[Content truncated due to token limit]"
}

// ChunkByParagraphs splits content into paragraph-based chunks.
// This is useful for more sophisticated chunking strategies.
func (c *Chunker) ChunkByParagraphs(content string) []string {
	if content == "" {
		return []string{}
	}

	// Split by double newlines (paragraph boundaries)
	paragraphs := strings.Split(content, "\n\n")

	var chunks []string
	currentChunk := ""
	currentTokens := 0

	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}

		paraTokens := c.EstimateTokens(para)

		// If adding this paragraph would exceed the limit, start a new chunk
		if currentTokens+paraTokens > c.MaxTokens && currentChunk != "" {
			chunks = append(chunks, strings.TrimSpace(currentChunk))
			currentChunk = para
			currentTokens = paraTokens
		} else {
			if currentChunk != "" {
				currentChunk += "\n\n"
			}
			currentChunk += para
			currentTokens += paraTokens
		}
	}

	// Add the last chunk if it's not empty
	if currentChunk != "" {
		chunks = append(chunks, strings.TrimSpace(currentChunk))
	}

	log.Debug().
		Int("input_length", len(content)).
		Int("chunk_count", len(chunks)).
		Msg("Content split into chunks")

	return chunks
}
