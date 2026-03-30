package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/devsh/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Interactively setup devsh configuration",
	Long:  "Launch an interactive terminal wizard to configure LLM providers, API keys, models, and safety levels.",
	Run: func(cmd *cobra.Command, args []string) {
		err := config.LoadConfig()
		if err != nil {
			// It's okay if loading fails initially (e.g., config file doesn't exist)
			config.AppConfig = config.DefaultConfig()
		}

		fmt.Println("🛠️  Welcome to the devsh interactive setup wizard!")

		defaultProvider := config.AppConfig.LLMProvider
		if defaultProvider == "openai" && config.AppConfig.OpenAIBaseURL != "" {
			defaultProvider = "other (openai-compatible)"
		}

		// 1. Ask for LLM Provider
		providerPrompt := &survey.Select{
			Message: "Select your LLM Provider:",
			Options: []string{"openai", "anthropic", "gemini", "ollama", "other (openai-compatible)"},
			Default: defaultProvider,
		}
		var provider string
		err = survey.AskOne(providerPrompt, &provider)
		if err != nil {
			fmt.Println("Setup cancelled.")
			os.Exit(0)
		}
		if provider == "other (openai-compatible)" {
			config.AppConfig.LLMProvider = "openai"
		} else {
			config.AppConfig.LLMProvider = provider
		}

		// 2. Ask for API Key based on provider
		switch provider {
		case "openai":
			keyPrompt := &survey.Password{
				Message: "Enter your OpenAI API Key (leave blank to keep current):",
			}
			var key string
			survey.AskOne(keyPrompt, &key)
			if key != "" {
				config.AppConfig.OpenAPIKey = key
			}
			config.AppConfig.OpenAIBaseURL = "" // Reset to official OpenAI URL

		case "other (openai-compatible)":
			keyPrompt := &survey.Password{
				Message: "Enter your API Key (leave blank if none required):",
			}
			var key string
			survey.AskOne(keyPrompt, &key)
			if key != "" {
				config.AppConfig.OpenAPIKey = key
			}
			
			urlPrompt := &survey.Input{
				Message: "Enter Remote Base URL (e.g., https://api.groq.com/openai/v1):",
				Default: config.AppConfig.OpenAIBaseURL,
			}
			survey.AskOne(urlPrompt, &config.AppConfig.OpenAIBaseURL)

		case "anthropic":
			keyPrompt := &survey.Password{
				Message: "Enter your Anthropic (Claude) API Key (leave blank to keep current):",
			}
			var key string
			survey.AskOne(keyPrompt, &key)
			if key != "" {
				config.AppConfig.AnthropicAPIKey = key
			}

		case "gemini":
			keyPrompt := &survey.Password{
				Message: "Enter your Gemini (Google) API Key (leave blank to keep current):",
			}
			var key string
			survey.AskOne(keyPrompt, &key)
			if key != "" {
				config.AppConfig.GeminiAPIKey = key
			}

		case "ollama":
			hostPrompt := &survey.Input{
				Message: "Enter Ollama Host URL:",
				Default: config.AppConfig.OllamaHost,
			}
			survey.AskOne(hostPrompt, &config.AppConfig.OllamaHost)
		}

		// 3. Ask for Model Name
		modelPrompt := &survey.Input{
			Message: fmt.Sprintf("Enter the Model name to use for %s:", provider),
			Default: config.AppConfig.Model,
		}
		survey.AskOne(modelPrompt, &config.AppConfig.Model)

		// 4. Ask for Safety Level
		safetyPrompt := &survey.Select{
			Message: "Select your Safety Level (How strictly should devsh ask for confirmations?):",
			Options: []string{"low", "medium", "high"},
			Default: config.AppConfig.SafetyLevel,
			Description: func(value string, index int) string {
				switch value {
				case "low":
					return "Auto-execute most non-destructive commands. Prompt rarely."
				case "medium":
					return "Auto-execute benign commands. Prompt for deletes/writes."
				case "high":
					return "Strict mode. Almost always prompt for confirmation."
				}
				return ""
			},
		}
		survey.AskOne(safetyPrompt, &config.AppConfig.SafetyLevel)

		// Save the finalized configuration
		err = config.SaveConfig()
		if err != nil {
			fmt.Printf("❌ Failed to save config: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ Configuration saved successfully!")
	},
}
