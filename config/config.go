package config

import (
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
)

type AgentsConfig struct {
	Model  string `mapstructure:"model"`
	Agents Agents `mapstructure:"agents"`
	MCP    MCP    `mapstructure:"mcp"`

	OpenAI OpenAI `mapstructure:"openai"`
	// Anthropic Anthropic `mapstructure:"anthropic"`
	Azure Azure `mapstructure:"azure"`
	// DeepSeek DeepSeek `mapstructure:"deepseek"`
}

type MCP struct {
	Servers map[string]MCPServer `mapstructure:"servers"`
}

type MCPServer struct {
	Transport    Transport         `mapstructure:"transport"`
	Url          string            `mapstructure:"url"`
	Command      string            `mapstructure:"command"`
	Args         []string          `mapstructure:"args"`
	Headers      map[string]string `mapstructure:"headers"`
	Environments map[string]string `mapstructure:"env"`
}

type Agents struct {
	Agents map[string]Agent
}

type Agent struct {
	Url string `mapstructure:"url"`
}

type Anthropic struct {
	ApiKey  string `mapstructure:"api_key"`
	BaseUrl string `mapstructure:"base_url"`
}

type OpenAI struct {
	ApiKey  string `mapstructure:"api_key"`
	BaseUrl string `mapstructure:"base_url"`
}

type Azure struct {
	UseDefaultAzureCredential bool   `mapstructure:"use_default_azure_credential"`
	ApiKey                    string `mapstructure:"api_key"`
	BaseUrl                   string `mapstructure:"base_url"`
	ApiVersion                string `mapstructure:"api_version"`
}

type DeepSeek struct {
	ApiKey  string `mapstructure:"api_key"`
	BaseUrl string `mapstructure:"base_url"`
}

type Google struct {
	ApiKey  string `mapstructure:"api_key"`
	BaseUrl string `mapstructure:"base_url"`
}

type Generic struct {
	ApiKey  string `mapstructure:"api_key"`
	BaseUrl string `mapstructure:"base_url"`
}

type OpenRouter struct {
	ApiKey  string `mapstructure:"api_key"`
	BaseUrl string `mapstructure:"base_url"`
}

type TensorZero struct {
	BaseUrl string `mapstructure:"base_url"`
}

func LoadConfig() (*AgentsConfig, error) {
	var agentsConfig AgentsConfig

	config := viper.New()
	config.SetConfigName("agents.config")
	config.SetConfigType("yaml")
	config.AddConfigPath(".")

	config.SetEnvPrefix("agents")
	// secrets.AllowEmptyEnv(true)
	config.AutomaticEnv()

	if err := config.ReadInConfig(); err != nil {
		slog.Info("ERROR", "err", err)
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file not found; ignore error if desired
			return nil, fmt.Errorf("error loading configs. %w", err)
		}
	}

	secrets := viper.New()
	secrets.SetConfigName("agents.secrets")
	secrets.SetConfigType("yaml")
	secrets.AddConfigPath(".")

	secrets.SetEnvPrefix("agents")
	// secrets.AllowEmptyEnv(true)
	secrets.AutomaticEnv()

	if err := secrets.ReadInConfig(); err != nil {
		slog.Info("ERROR", "err", err)
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file not found; ignore error if desired
			return nil, fmt.Errorf("error loading secrets. %w", err)
		}
	}

	config.MergeConfigMap(secrets.AllSettings())

	if err := config.Unmarshal(&agentsConfig); err != nil {
		return nil, fmt.Errorf("error load agents.config.yaml. %w", err)
	}

	return &agentsConfig, nil
}
