package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAIClient_GenerateJSON_Success tests successful JSON generation from the AI client.
func TestAIClient_GenerateJSON_Success(t *testing.T) {
	// Arrange: Create a mock server that returns a valid OpenAI-compatible JSON response
	expectedContent := `{"title": "Test Document", "summary": "This is a test"}`
	
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		assert.Equal(t, "/v1/chat/completions", r.URL.Path)
		assert.Equal(t, "Bearer test-api-key", r.Header.Get("Authorization"))
		
		// Return a valid OpenAI response
		response := map[string]interface{}{
			"id":      "test-id",
			"object":  "chat.completion",
			"created": time.Now().Unix(),
			"model":   "gpt-3.5-turbo",
			"choices": []map[string]interface{}{
				{
					"index": 0,
					"message": map[string]interface{}{
						"role":    "assistant",
						"content": expectedContent,
					},
					"finish_reason": "stop",
				},
			},
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Create AI client pointing to mock server
	config := Config{
		BaseURL:    mockServer.URL + "/v1",
		APIKey:     "test-api-key",
		Model:      "gpt-3.5-turbo",
		MaxRetries: 0, // No retries for this test
	}
	
	client, err := NewOpenAIClient(config)
	require.NoError(t, err)

	// Act: Call the client's generation method
	ctx := context.Background()
	result, err := client.GenerateJSON(ctx, "Generate a test document")

	// Assert: The client correctly parses the response and returns the content string without error
	require.NoError(t, err)
	assert.Equal(t, expectedContent, result)
	
	// Verify the returned content is valid JSON
	var jsonCheck interface{}
	err = json.Unmarshal([]byte(result), &jsonCheck)
	assert.NoError(t, err)
}

// TestAIClient_GenerateJSON_RetriesOn503 tests that the client retries on 503 Service Unavailable.
func TestAIClient_GenerateJSON_RetriesOn503(t *testing.T) {
	// Arrange: Create a mock server that returns 503 twice, then 200 OK on the third call
	requestCount := int32(0)
	expectedContent := `{"status": "success"}`
	
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&requestCount, 1)
		
		if count <= 2 {
			// First two requests return 503
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]interface{}{
					"message": "Service temporarily unavailable",
					"type":    "service_unavailable",
					"code":    "503",
				},
			})
			return
		}
		
		// Third request succeeds
		response := map[string]interface{}{
			"id":      "test-id",
			"object":  "chat.completion",
			"created": time.Now().Unix(),
			"model":   "gpt-3.5-turbo",
			"choices": []map[string]interface{}{
				{
					"index": 0,
					"message": map[string]interface{}{
						"role":    "assistant",
						"content": expectedContent,
					},
					"finish_reason": "stop",
				},
			},
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Configure the client with 3 retries and short retry delay
	config := Config{
		BaseURL:    mockServer.URL + "/v1",
		APIKey:     "test-api-key",
		Model:      "gpt-3.5-turbo",
		MaxRetries: 2, // Allow up to 2 retries (3 total attempts)
		RetryDelay: 10 * time.Millisecond, // Short delay for testing
	}
	
	client, err := NewOpenAIClient(config)
	require.NoError(t, err)

	// Act: Call the generation method
	ctx := context.Background()
	result, err := client.GenerateJSON(ctx, "Test prompt")

	// Assert: The client should make three requests and ultimately succeed
	require.NoError(t, err)
	assert.Equal(t, expectedContent, result)
	assert.Equal(t, int32(3), atomic.LoadInt32(&requestCount), "Expected exactly 3 requests")
}

// TestNewOpenAIClient_Validation tests client creation validation.
func TestNewOpenAIClient_Validation(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectError string
	}{
		{
			name: "missing API key",
			config: Config{
				Model: "gpt-3.5-turbo",
			},
			expectError: "API key is required",
		},
		{
			name: "missing model",
			config: Config{
				APIKey: "test-key",
			},
			expectError: "model is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewOpenAIClient(tt.config)
			assert.Nil(t, client)
			assert.EqualError(t, err, tt.expectError)
		})
	}
}

// TestOpenAIClient_GenerateJSON_InvalidJSON tests handling of invalid JSON response.
func TestOpenAIClient_GenerateJSON_InvalidJSON(t *testing.T) {
	// Arrange: Mock server returns invalid JSON
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"id":      "test-id",
			"object":  "chat.completion",
			"created": time.Now().Unix(),
			"model":   "gpt-3.5-turbo",
			"choices": []map[string]interface{}{
				{
					"index": 0,
					"message": map[string]interface{}{
						"role":    "assistant",
						"content": "This is not valid JSON", // Invalid JSON
					},
					"finish_reason": "stop",
				},
			},
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	config := Config{
		BaseURL:    mockServer.URL + "/v1",
		APIKey:     "test-api-key",
		Model:      "gpt-3.5-turbo",
		MaxRetries: 0,
	}
	
	client, err := NewOpenAIClient(config)
	require.NoError(t, err)

	// Act
	ctx := context.Background()
	result, err := client.GenerateJSON(ctx, "Test prompt")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "AI response is not valid JSON")
	assert.Empty(t, result)
}

// TestOpenAIClient_GenerateJSON_ContextCancellation tests that context cancellation stops retries.
func TestOpenAIClient_GenerateJSON_ContextCancellation(t *testing.T) {
	requestCount := int32(0)
	
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requestCount, 1)
		// Always return 503 to trigger retries
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{
				"message": "Service temporarily unavailable",
				"type":    "service_unavailable",
				"code":    "503",
			},
		})
	}))
	defer mockServer.Close()

	config := Config{
		BaseURL:    mockServer.URL + "/v1",
		APIKey:     "test-api-key",
		Model:      "gpt-3.5-turbo",
		MaxRetries: 5,
		RetryDelay: 100 * time.Millisecond,
	}
	
	client, err := NewOpenAIClient(config)
	require.NoError(t, err)

	// Create a context that will be cancelled after first request
	ctx, cancel := context.WithCancel(context.Background())
	
	// Cancel context after a short delay
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	// Act
	result, err := client.GenerateJSON(ctx, "Test prompt")

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
	assert.Empty(t, result)
	// Should have made only 1 or 2 requests before cancellation
	assert.LessOrEqual(t, atomic.LoadInt32(&requestCount), int32(2))
}

// TestOpenAIClient_ConfigIntegration tests that model and base URL configuration are properly used.
func TestOpenAIClient_ConfigIntegration(t *testing.T) {
	tests := []struct {
		name           string
		config         Config
		expectedModel  string
		expectedPath   string
		expectedAuth   string
	}{
		{
			name: "OpenAI default configuration",
			config: Config{
				BaseURL: "https://api.openai.com/v1",
				APIKey:  "sk-test123",
				Model:   "gpt-4",
			},
			expectedModel: "gpt-4",
			expectedPath:  "/v1/chat/completions",
			expectedAuth:  "Bearer sk-test123",
		},
		{
			name: "Azure OpenAI configuration",
			config: Config{
				BaseURL: "https://myinstance.openai.azure.com",
				APIKey:  "azure-key-456",
				Model:   "gpt-35-turbo",
			},
			expectedModel: "gpt-35-turbo",
			expectedPath:  "/chat/completions",
			expectedAuth:  "Bearer azure-key-456",
		},
		{
			name: "Local LLM configuration",
			config: Config{
				BaseURL: "http://localhost:8080/v1",
				APIKey:  "local-key",
				Model:   "llama2-7b",
			},
			expectedModel: "llama2-7b",
			expectedPath:  "/v1/chat/completions",
			expectedAuth:  "Bearer local-key",
		},
		{
			name: "Claude via OpenAI-compatible API",
			config: Config{
				BaseURL: "https://api.anthropic.com/v1",
				APIKey:  "claude-key",
				Model:   "claude-3-opus",
			},
			expectedModel: "claude-3-opus",
			expectedPath:  "/v1/chat/completions",
			expectedAuth:  "Bearer claude-key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock server to capture the request details
			var capturedModel string
			var capturedPath string
			var capturedAuth string
			
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedPath = r.URL.Path
				capturedAuth = r.Header.Get("Authorization")
				
				// Parse request body to get the model
				var reqBody map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err == nil {
					if model, ok := reqBody["model"].(string); ok {
						capturedModel = model
					}
				}
				
				// Return a valid response
				response := map[string]interface{}{
					"id":      "test-id",
					"object":  "chat.completion",
					"created": time.Now().Unix(),
					"model":   tt.expectedModel,
					"choices": []map[string]interface{}{
						{
							"index": 0,
							"message": map[string]interface{}{
								"role":    "assistant",
								"content": `{"test": "response"}`,
							},
							"finish_reason": "stop",
						},
					},
				}
				
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(response)
			}))
			defer mockServer.Close()
			
			// Update config to use mock server, preserving the path structure
			tt.config.BaseURL = mockServer.URL + strings.TrimPrefix(tt.config.BaseURL, strings.Split(tt.config.BaseURL, "/")[0]+"//"+strings.Split(tt.config.BaseURL, "/")[2])
			
			// Create client
			client, err := NewOpenAIClient(tt.config)
			require.NoError(t, err)
			
			// Make a request
			ctx := context.Background()
			_, err = client.GenerateJSON(ctx, "Test prompt")
			require.NoError(t, err)
			
			// Verify the request used the correct configuration
			assert.Equal(t, tt.expectedModel, capturedModel, "Model mismatch")
			assert.Equal(t, tt.expectedPath, capturedPath, "Path mismatch")
			assert.Equal(t, tt.expectedAuth, capturedAuth, "Authorization header mismatch")
		})
	}
}

// TestOpenAIClient_DefaultBaseURL tests that the default base URL is set when not provided.
func TestOpenAIClient_DefaultBaseURL(t *testing.T) {
	// Create a client without specifying BaseURL
	config := Config{
		APIKey: "test-key",
		Model:  "gpt-4",
	}
	
	client, err := NewOpenAIClient(config)
	require.NoError(t, err)
	require.NotNil(t, client)
	
	// Verify the default BaseURL was set
	assert.Equal(t, "https://api.openai.com/v1", client.config.BaseURL)
}