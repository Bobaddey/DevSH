package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the DEVSH configuration
type Config struct {
	LLMProvider     string   `yaml:"llm_provider"`
	Model           string   `yaml:"model"`
	OpenAPIKey      string   `yaml:"openai_api_key,omitempty"` // Can also be set via OPENAI_API_KEY
	OpenAIBaseURL   string   `yaml:"openai_base_url,omitempty"` // For OpenAI-compatible APIs (e.g. Groq, LMStudio)
	AnthropicAPIKey string   `yaml:"anthropic_api_key,omitempty"` // Can also be set via ANTHROPIC_API_KEY
	GeminiAPIKey    string   `yaml:"gemini_api_key,omitempty"` // Can also be set via GEMINI_API_KEY
	OllamaHost      string   `yaml:"ollama_host,omitempty"`    // e.g., http://localhost:11434
	SafetyLevel     string   `yaml:"safety_level"`             // low, medium, high
	DefaultTools    []string `yaml:"default_tools"`
}

var AppConfig *Config

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		LLMProvider:  "openai",
		Model:        "gpt-4",
		SafetyLevel:  "high",
		DefaultTools: []string{"bash", "docker", "kubectl", "aws", "terraform", "git"},
		OllamaHost:   "http://localhost:11434",
	}
}

// LoadConfig loads the configuration from ~/.devsh/config.yaml
// If the file does not exist, it creates it with default values.
func LoadConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not find user home directory: %v", err)
	}

	configDir := filepath.Join(homeDir, ".devsh")
	configFile := filepath.Join(configDir, "config.yaml")

	// Create config dir if not exists
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err = os.MkdirAll(configDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create config directory: %v", err)
		}
	}

	// Create default config file if it doesn't exist
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		AppConfig = DefaultConfig()
		if err := SaveConfig(); err != nil {
			return fmt.Errorf("failed to save default config: %v", err)
		}
		return nil
	}

	// Read existing config
	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	AppConfig = &Config{}
	err = yaml.Unmarshal(data, AppConfig)
	if err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}

	// Environment variable overrides
	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		AppConfig.OpenAPIKey = key
	}
	if url := os.Getenv("OPENAI_BASE_URL"); url != "" {
		AppConfig.OpenAIBaseURL = url
	}
	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		AppConfig.AnthropicAPIKey = key
	}
	if key := os.Getenv("GEMINI_API_KEY"); key != "" {
		AppConfig.GeminiAPIKey = key
	}
	if host := os.Getenv("OLLAMA_HOST"); host != "" {
		AppConfig.OllamaHost = host
	}

	return nil
}

// SaveConfig saves the current configuration to file
func SaveConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configFile := filepath.Join(homeDir, ".devsh", "config.yaml")

	data, err := yaml.Marshal(AppConfig)
	if err != nil {
		return err
	}

	return os.WriteFile(configFile, data, 0644)
}
