package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

// Config represents the application configuration
type Config struct {
	Model       string  `yaml:"model" env:"DOCLOOM_MODEL"`
	BaseURL     string  `yaml:"base_url" env:"DOCLOOM_BASE_URL"`
	APIKey      string  `yaml:"api_key" env:"DOCLOOM_API_KEY,OPENAI_API_KEY"`
	TemplateDir string  `yaml:"template_dir" env:"DOCLOOM_TEMPLATE_DIR"`
	Temperature float64 `yaml:"temperature" env:"DOCLOOM_TEMPERATURE"`
	Seed        int     `yaml:"seed" env:"DOCLOOM_SEED"`
	MaxRetries  int     `yaml:"max_retries" env:"DOCLOOM_MAX_RETRIES"`
	Force       bool    `yaml:"force" env:"DOCLOOM_FORCE"`
	Verbose     bool    `yaml:"verbose" env:"DOCLOOM_VERBOSE"`
	DryRun      bool    `yaml:"dry_run" env:"DOCLOOM_DRY_RUN"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Model:       "gpt-4",
		BaseURL:     "https://api.openai.com/v1",
		Temperature: 0.7,
		MaxRetries:  3,
		TemplateDir: "templates",
		Force:       false,
		Verbose:     false,
		DryRun:      false,
	}
}

// Load loads configuration with proper precedence: CLI flags > ENV > File > Defaults
// This function demonstrates the precedence but actual implementation will need
// proper YAML parsing and flag integration
func Load(configFile string, cliOverrides map[string]interface{}) (*Config, error) {
	// Start with defaults
	cfg := DefaultConfig()

	// Load from file if provided
	if configFile != "" {
		if err := loadFromFile(cfg, configFile); err != nil {
			log.Warn().Err(err).Str("file", configFile).Msg("Failed to load config file, using defaults")
		}
	}

	// Override with environment variables
	loadFromEnv(cfg)

	// Override with CLI flags (highest precedence)
	applyCliOverrides(cfg, cliOverrides)

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// LoadWithPrecedence is a simplified version for testing that demonstrates precedence
func LoadWithPrecedence(fileValue, envValue, cliValue string, field string) string {
	// Start with default
	result := ""

	// File has lowest precedence (after defaults)
	if fileValue != "" {
		result = fileValue
	}

	// Environment overrides file
	if envValue != "" {
		result = envValue
	}

	// CLI has highest precedence
	if cliValue != "" {
		result = cliValue
	}

	return result
}

// loadFromFile loads configuration from a YAML file
func loadFromFile(cfg *Config, path string) error {
	// For now, we'll stub this out as YAML parsing will be added when needed
	// In a real implementation, this would use a YAML library to unmarshal the file
	log.Debug().Str("path", path).Msg("Loading config from file")
	// TODO: Implement YAML parsing and populate cfg
	_ = cfg // cfg will be used when YAML parsing is implemented
	return nil
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv(cfg *Config) {
	// Check for model override
	if val := os.Getenv("DOCLOOM_MODEL"); val != "" {
		cfg.Model = val
	}

	// Check for base URL override
	if val := os.Getenv("DOCLOOM_BASE_URL"); val != "" {
		cfg.BaseURL = val
	}

	// Check for API key from multiple sources
	if val := os.Getenv("DOCLOOM_API_KEY"); val != "" {
		cfg.APIKey = val
	} else if val := os.Getenv("OPENAI_API_KEY"); val != "" {
		cfg.APIKey = val
	}

	// Check for temperature override
	if val := os.Getenv("DOCLOOM_TEMPERATURE"); val != "" {
		// In production, parse to float64
		log.Debug().Str("temperature", val).Msg("Temperature override from env")
	}

	// Check for template directory override
	if val := os.Getenv("DOCLOOM_TEMPLATE_DIR"); val != "" {
		cfg.TemplateDir = val
	}
}

// applyStringOverride applies a string override if valid
func applyStringOverride(target *string, value interface{}) {
	if v, ok := value.(string); ok && v != "" {
		*target = v
	}
}

// applyIntOverride applies an int override if valid
func applyIntOverride(target *int, value interface{}) {
	if v, ok := value.(int); ok {
		*target = v
	}
}

// applyFloatOverride applies a float override if valid
func applyFloatOverride(target *float64, value interface{}) {
	if v, ok := value.(float64); ok {
		*target = v
	}
}

// applyBoolOverride applies a bool override if valid
func applyBoolOverride(target *bool, value interface{}) {
	if v, ok := value.(bool); ok {
		*target = v
	}
}

// applyCliOverrides applies CLI flag overrides to the configuration
func applyCliOverrides(cfg *Config, overrides map[string]interface{}) {
	if overrides == nil {
		return
	}

	// Apply each override
	for key, value := range overrides {
		switch key {
		case "model":
			applyStringOverride(&cfg.Model, value)
		case "base_url":
			applyStringOverride(&cfg.BaseURL, value)
		case "api_key":
			applyStringOverride(&cfg.APIKey, value)
		case "temperature":
			applyFloatOverride(&cfg.Temperature, value)
		case "seed":
			applyIntOverride(&cfg.Seed, value)
		case "max_retries":
			applyIntOverride(&cfg.MaxRetries, value)
		case "template_dir":
			applyStringOverride(&cfg.TemplateDir, value)
		case "force":
			applyBoolOverride(&cfg.Force, value)
		case "verbose":
			applyBoolOverride(&cfg.Verbose, value)
		case "dry_run":
			applyBoolOverride(&cfg.DryRun, value)
		}
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Basic validation - can be expanded as needed
	if c.MaxRetries < 0 {
		c.MaxRetries = 0
	}

	if c.Temperature < 0 || c.Temperature > 2 {
		c.Temperature = 0.7 // Reset to default if out of range
	}

	// Ensure template directory is absolute or relative to working directory
	if c.TemplateDir != "" && !filepath.IsAbs(c.TemplateDir) {
		if wd, err := os.Getwd(); err == nil {
			c.TemplateDir = filepath.Join(wd, c.TemplateDir)
		}
	}

	return nil
}

// String returns a string representation of the config (hiding sensitive values)
func (c *Config) String() string {
	apiKeyDisplay := "<not set>"
	if c.APIKey != "" {
		if len(c.APIKey) > 8 {
			apiKeyDisplay = c.APIKey[:4] + "..." + c.APIKey[len(c.APIKey)-4:]
		} else {
			apiKeyDisplay = strings.Repeat("*", len(c.APIKey))
		}
	}

	return strings.Join([]string{
		"Config{",
		"  Model: " + c.Model,
		"  BaseURL: " + c.BaseURL,
		"  APIKey: " + apiKeyDisplay,
		"  Temperature: " + string(rune(c.Temperature)),
		"  TemplateDir: " + c.TemplateDir,
		"}",
	}, "\n")
}
