package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

// ClaudeClient handles communication with the Claude API
type ClaudeClient struct {
	apiKey    string
	model     string
	maxTokens int
}

// NewClaudeClient creates a new Claude API client
func NewClaudeClient(apiKey, model string, maxTokens int) *ClaudeClient {
	return &ClaudeClient{
		apiKey:    apiKey,
		model:     model,
		maxTokens: maxTokens,
	}
}

// AnalysisResponse represents the structured response from Claude
type AnalysisResponse struct {
	ProjectName     string              `json:"projectName"`
	Description     string              `json:"description"`
	ProjectType     string              `json:"projectType"`
	Framework       string              `json:"framework"`
	Architecture    Architecture        `json:"architecture"`
	Dependencies    Dependencies        `json:"dependencies"`
	Features        []Feature           `json:"features"`
	APIs            []API               `json:"apis"`
	DataModel       DataModel           `json:"dataModel"`
	Testing         Testing             `json:"testing"`
	Deployment      Deployment          `json:"deployment"`
	Security        Security            `json:"security"`
	TechnicalDebt   []TechnicalDebtItem `json:"technicalDebt"`
	Recommendations []string            `json:"recommendations"`
	RawResponse     string              `json:"-"` // Store raw response for debugging
}

type Architecture struct {
	Pattern       string   `json:"pattern"`
	Layers        []string `json:"layers"`
	KeyComponents []string `json:"keyComponents"`
}

type Dependencies struct {
	NuGet    []string `json:"nuget"`
	External []string `json:"external"`
}

type Feature struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type API struct {
	Endpoint    string `json:"endpoint"`
	Method      string `json:"method"`
	Description string `json:"description"`
}

type DataModel struct {
	Entities []string `json:"entities"`
	Database string   `json:"database"`
}

type Testing struct {
	Framework string   `json:"framework"`
	Coverage  string   `json:"coverage"`
	Types     []string `json:"types"`
}

type Deployment struct {
	Containerized bool   `json:"containerized"`
	CICD          string `json:"cicd"`
	Hosting       string `json:"hosting"`
}

type Security struct {
	Authentication string   `json:"authentication"`
	Authorization  string   `json:"authorization"`
	Considerations []string `json:"considerations"`
}

type TechnicalDebtItem struct {
	Area           string `json:"area"`
	Description    string `json:"description"`
	Impact         string `json:"impact"`
	Recommendation string `json:"recommendation"`
}

// Analyze sends the prompt to Claude and returns the structured response
func (c *ClaudeClient) Analyze(prompt string) (*AnalysisResponse, error) {
	// Use OpenAI-compatible client for Claude API
	client := openai.NewClient(c.apiKey)

	// Configure for Claude (Anthropic) API if needed
	config := openai.DefaultConfig(c.apiKey)
	if strings.Contains(c.model, "claude") {
		// Set Anthropic API endpoint if using Claude directly
		config.BaseURL = "https://api.anthropic.com/v1"
	}
	client = openai.NewClientWithConfig(config)

	// Create the chat completion request
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: c.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are an expert C# software architect. Provide your analysis as valid JSON matching the provided structure.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens:   c.maxTokens,
			Temperature: 0.3,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to call Claude API: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from Claude API")
	}

	// Extract JSON from response
	rawResponse := resp.Choices[0].Message.Content

	// Parse JSON response
	var analysis AnalysisResponse
	if err := json.Unmarshal([]byte(rawResponse), &analysis); err != nil {
		// Try to extract JSON if it's wrapped in markdown code blocks
		jsonStart := strings.Index(rawResponse, "{")
		jsonEnd := strings.LastIndex(rawResponse, "}")
		if jsonStart >= 0 && jsonEnd > jsonStart {
			jsonContent := rawResponse[jsonStart : jsonEnd+1]
			if err := json.Unmarshal([]byte(jsonContent), &analysis); err != nil {
				return nil, fmt.Errorf("failed to parse Claude response as JSON: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to parse Claude response as JSON: %w", err)
		}
	}

	analysis.RawResponse = rawResponse
	return &analysis, nil
}
