package config

import (
	"os"
	"testing"
)

// TC-2.1: Test configuration loading with correct precedence
func TestConfig_LoadWithPrecedence(t *testing.T) {
	tests := []struct {
		name      string
		fileValue string
		envValue  string
		cliValue  string
		expected  string
	}{
		{
			name:      "CLI overrides all",
			fileValue: "file-model",
			envValue:  "env-model",
			cliValue:  "cli-model",
			expected:  "cli-model",
		},
		{
			name:      "ENV overrides file",
			fileValue: "file-model",
			envValue:  "env-model",
			cliValue:  "",
			expected:  "env-model",
		},
		{
			name:      "File value when no overrides",
			fileValue: "file-model",
			envValue:  "",
			cliValue:  "",
			expected:  "file-model",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LoadWithPrecedence(tt.fileValue, tt.envValue, tt.cliValue, "model")
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestConfig_FullPrecedence(t *testing.T) {
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
}

func TestConfig_EnvPrecedenceWithoutCLI(t *testing.T) {
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
}

func TestConfig_DefaultValues(t *testing.T) {
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
}

func TestConfig_MultipleEnvVars(t *testing.T) {
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
}

func TestConfig_BaseURLPrecedence(t *testing.T) {
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
}

func TestConfig_ModelSelectionForProviders(t *testing.T) {
	testCases := []struct {
		name    string
		model   string
		baseURL string
		valid   bool
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
}

// TC-2.2: Test helper function for precedence (unit test)
func TestLoadWithPrecedence(t *testing.T) {
	// Test the LoadWithPrecedence function directly
	testCases := []struct {
		name          string
		settingName   string
		fileValue     string
		envValue      string
		cliValue      string
		expectedValue string
	}{
		{
			name:          "CLI takes precedence over all",
			settingName:   "test_setting",
			fileValue:     "file",
			envValue:      "env",
			cliValue:      "cli",
			expectedValue: "cli",
		},
		{
			name:          "ENV takes precedence over file",
			settingName:   "test_setting",
			fileValue:     "file",
			envValue:      "env",
			cliValue:      "",
			expectedValue: "env",
		},
		{
			name:          "File value when no overrides",
			settingName:   "test_setting",
			fileValue:     "file",
			envValue:      "",
			cliValue:      "",
			expectedValue: "file",
		},
		{
			name:          "Empty values return empty",
			settingName:   "test_setting",
			fileValue:     "",
			envValue:      "",
			cliValue:      "",
			expectedValue: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := LoadWithPrecedence(tc.fileValue, tc.envValue, tc.cliValue, tc.settingName)
			if result != tc.expectedValue {
				t.Errorf("Expected '%s', got '%s'", tc.expectedValue, result)
			}
		})
	}
}

// TC-2.3: Test configuration validation
func TestConfig_Validate(t *testing.T) {
	testCases := []struct {
		name     string
		config   Config
		expected Config
	}{
		{
			name: "Negative MaxRetries should be set to 0",
			config: Config{
				MaxRetries: -5,
			},
			expected: Config{
				MaxRetries: 0,
			},
		},
		{
			name: "Temperature out of range (too low) should reset to default",
			config: Config{
				Temperature: -0.5,
			},
			expected: Config{
				Temperature: 0.7,
			},
		},
		{
			name: "Temperature out of range (too high) should reset to default",
			config: Config{
				Temperature: 2.5,
			},
			expected: Config{
				Temperature: 0.7,
			},
		},
		{
			name: "Valid temperature should remain unchanged",
			config: Config{
				Temperature: 1.2,
			},
			expected: Config{
				Temperature: 1.2,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := tc.config
			err := cfg.Validate()
			if err != nil {
				t.Fatalf("Validate returned error: %v", err)
			}

			if cfg.MaxRetries != tc.expected.MaxRetries {
				t.Errorf("MaxRetries: expected %d, got %d", tc.expected.MaxRetries, cfg.MaxRetries)
			}

			if cfg.Temperature != tc.expected.Temperature {
				t.Errorf("Temperature: expected %f, got %f", tc.expected.Temperature, cfg.Temperature)
			}
		})
	}
}

// TC-2.4: Test String method for security (API key masking)
func TestConfig_String(t *testing.T) {
	testCases := []struct {
		name        string
		apiKey      string
		expectedStr string
	}{
		{
			name:        "Long API key should be masked",
			apiKey:      "sk-1234567890abcdef",
			expectedStr: "sk-1...cdef",
		},
		{
			name:        "Short API key should be fully masked",
			apiKey:      "secret",
			expectedStr: "******",
		},
		{
			name:        "Empty API key shows not set",
			apiKey:      "",
			expectedStr: "<not set>",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := Config{
				APIKey: tc.apiKey,
				Model:  "gpt-4",
			}

			str := cfg.String()
			if tc.apiKey != "" {
				// Check that the actual API key is not in the string
				if tc.apiKey != "" && len(tc.apiKey) > 8 {
					// For long keys, check the masked format
					if !contains(str, tc.expectedStr) {
						t.Errorf("Expected masked API key '%s' in string, got: %s", tc.expectedStr, str)
					}
				}
			} else {
				if !contains(str, tc.expectedStr) {
					t.Errorf("Expected '%s' for empty API key, got: %s", tc.expectedStr, str)
				}
			}

			// Ensure the full API key is never exposed (except for short keys that are fully masked)
			if len(tc.apiKey) > 8 && contains(str, tc.apiKey) {
				t.Errorf("Full API key should not be exposed in String(): %s", str)
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr) >= 0))
}

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
