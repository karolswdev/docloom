package config

import (
	"os"
	"testing"
)

// TC-2.1: Test configuration loading with correct precedence
func TestConfig_LoadWithPrecedence(t *testing.T) {
	// Test case 1: CLI flag should override everything
	t.Run("CLI overrides all", func(t *testing.T) {
		// Arrange
		fileValue := "file-model"
		envValue := "env-model"
		cliValue := "cli-model"
		
		// Act
		result := LoadWithPrecedence(fileValue, envValue, cliValue, "model")
		
		// Assert
		if result != cliValue {
			t.Errorf("Expected CLI value '%s' to have highest precedence, got '%s'", cliValue, result)
		}
	})
	
	// Test case 2: ENV should override file when no CLI value
	t.Run("ENV overrides file", func(t *testing.T) {
		// Arrange
		fileValue := "file-model"
		envValue := "env-model"
		cliValue := "" // No CLI value
		
		// Act
		result := LoadWithPrecedence(fileValue, envValue, cliValue, "model")
		
		// Assert
		if result != envValue {
			t.Errorf("Expected ENV value '%s' to override file value, got '%s'", envValue, result)
		}
	})
	
	// Test case 3: File value should be used when no ENV or CLI
	t.Run("File value when no overrides", func(t *testing.T) {
		// Arrange
		fileValue := "file-model"
		envValue := "" // No ENV value
		cliValue := "" // No CLI value
		
		// Act
		result := LoadWithPrecedence(fileValue, envValue, cliValue, "model")
		
		// Assert
		if result != fileValue {
			t.Errorf("Expected file value '%s' when no overrides, got '%s'", fileValue, result)
		}
	})
	
	// Test case 4: Full integration test with actual Config struct
	t.Run("Full config precedence", func(t *testing.T) {
		// Arrange
		// Set environment variable
		os.Setenv("DOCLOOM_MODEL", "env-gpt-4")
		defer os.Unsetenv("DOCLOOM_MODEL")
		
		// CLI overrides
		cliOverrides := map[string]interface{}{
			"model": "cli-gpt-4-turbo",
		}
		
		// Act
		cfg, err := Load("", cliOverrides)
		
		// Assert
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}
		
		// CLI value should win
		if cfg.Model != "cli-gpt-4-turbo" {
			t.Errorf("Expected CLI override 'cli-gpt-4-turbo', got '%s'", cfg.Model)
		}
	})
	
	// Test case 5: ENV precedence without CLI override
	t.Run("ENV precedence without CLI", func(t *testing.T) {
		// Arrange
		os.Setenv("DOCLOOM_MODEL", "env-model-value")
		defer os.Unsetenv("DOCLOOM_MODEL")
		
		// No CLI overrides
		cliOverrides := map[string]interface{}{}
		
		// Act
		cfg, err := Load("", cliOverrides)
		
		// Assert
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}
		
		// ENV value should be used
		if cfg.Model != "env-model-value" {
			t.Errorf("Expected ENV value 'env-model-value', got '%s'", cfg.Model)
		}
	})
	
	// Test case 6: Default values when nothing is set
	t.Run("Default values", func(t *testing.T) {
		// Arrange - ensure no env vars are set
		os.Unsetenv("DOCLOOM_MODEL")
		
		// Act
		cfg, err := Load("", nil)
		
		// Assert
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}
		
		// Should use default
		if cfg.Model != "gpt-4" {
			t.Errorf("Expected default model 'gpt-4', got '%s'", cfg.Model)
		}
		
		if cfg.Temperature != 0.7 {
			t.Errorf("Expected default temperature 0.7, got %f", cfg.Temperature)
		}
		
		if cfg.MaxRetries != 3 {
			t.Errorf("Expected default max_retries 3, got %d", cfg.MaxRetries)
		}
	})
	
	// Test case 7: Multiple environment variables
	t.Run("Multiple env vars", func(t *testing.T) {
		// Arrange
		os.Setenv("DOCLOOM_MODEL", "env-model")
		os.Setenv("DOCLOOM_BASE_URL", "https://api.example.com")
		os.Setenv("OPENAI_API_KEY", "test-key")
		defer func() {
			os.Unsetenv("DOCLOOM_MODEL")
			os.Unsetenv("DOCLOOM_BASE_URL")
			os.Unsetenv("OPENAI_API_KEY")
		}()
		
		// Act
		cfg, err := Load("", nil)
		
		// Assert
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}
		
		if cfg.Model != "env-model" {
			t.Errorf("Expected model from env 'env-model', got '%s'", cfg.Model)
		}
		
		if cfg.BaseURL != "https://api.example.com" {
			t.Errorf("Expected base URL from env, got '%s'", cfg.BaseURL)
		}
		
		if cfg.APIKey != "test-key" {
			t.Errorf("Expected API key from env, got '%s'", cfg.APIKey)
		}
	})
	
	// Test case 8: BaseURL configuration precedence
	t.Run("BaseURL precedence", func(t *testing.T) {
		// Test default BaseURL
		cfg := DefaultConfig()
		if cfg.BaseURL != "https://api.openai.com/v1" {
			t.Errorf("Expected default BaseURL 'https://api.openai.com/v1', got '%s'", cfg.BaseURL)
		}
		
		// Test environment override
		os.Setenv("DOCLOOM_BASE_URL", "https://custom.api.com/v1")
		defer os.Unsetenv("DOCLOOM_BASE_URL")
		
		cfg, err := Load("", nil)
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}
		if cfg.BaseURL != "https://custom.api.com/v1" {
			t.Errorf("Expected BaseURL from env 'https://custom.api.com/v1', got '%s'", cfg.BaseURL)
		}
		
		// Test CLI override
		cliOverrides := map[string]interface{}{
			"base_url": "https://cli.api.com/v1",
		}
		cfg, err = Load("", cliOverrides)
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}
		if cfg.BaseURL != "https://cli.api.com/v1" {
			t.Errorf("Expected BaseURL from CLI 'https://cli.api.com/v1', got '%s'", cfg.BaseURL)
		}
	})
	
	// Test case 9: Model selection for different providers
	t.Run("Model selection for different providers", func(t *testing.T) {
		testCases := []struct {
			name     string
			model    string
			baseURL  string
			valid    bool
		}{
			{"OpenAI GPT-4", "gpt-4", "https://api.openai.com/v1", true},
			{"OpenAI GPT-3.5", "gpt-3.5-turbo", "https://api.openai.com/v1", true},
			{"Azure OpenAI", "gpt-4", "https://myinstance.openai.azure.com", true},
			{"Local LLM", "llama2", "http://localhost:8080/v1", true},
			{"Claude via OpenAI-compatible API", "claude-3-opus", "https://api.anthropic.com/v1", true},
		}
		
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				cliOverrides := map[string]interface{}{
					"model":    tc.model,
					"base_url": tc.baseURL,
				}
				
				cfg, err := Load("", cliOverrides)
				if err != nil && tc.valid {
					t.Fatalf("Failed to load config for %s: %v", tc.name, err)
				}
				
				if tc.valid {
					if cfg.Model != tc.model {
						t.Errorf("Expected model '%s', got '%s'", tc.model, cfg.Model)
					}
					if cfg.BaseURL != tc.baseURL {
						t.Errorf("Expected BaseURL '%s', got '%s'", tc.baseURL, cfg.BaseURL)
					}
				}
			})
		}
	})
}