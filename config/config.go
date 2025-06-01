package config

import (
	"fmt"
	"log/slog"

	"github.com/jlrosende/go-agents/mcp"
	"github.com/spf13/viper"
)

type AgentsConfig struct {
	Agents map[string]Agent `mapstructure:"agents"`
	MCP    MCP              `mapstructure:"mcp"`

	OpenAI     OpenAI     `mapstructure:"openai"`
	Anthropic  Anthropic  `mapstructure:"anthropic"`
	Azure      Azure      `mapstructure:"azure"`
	Generic    Generic    `mapstructure:"geeneric"`
	Google     Google     `mapstructure:"google"`
	DeepSeek   DeepSeek   `mapstructure:"deepseek"`
	OpenRouter OpenRouter `mapstructure:"openrouter"`
	TensorZero TensorZero `mapstructure:"tensorzero"`

	Logger Logger
}

type MCP struct {
	Servers map[string]MCPServer `mapstructure:"servers"`
}

type MCPServer struct {
	Transport    mcp.Transport     `mapstructure:"transport"`
	Url          string            `mapstructure:"url"`
	Command      string            `mapstructure:"command"`
	Args         []string          `mapstructure:"args"`
	Headers      map[string]string `mapstructure:"headers"`
	Environments map[string]string `mapstructure:"env"`
}

type Agent struct {
	Url          string   `mapstructure:"url"`
	Model        string   `mapstructure:"model"`
	Instructions string   `mapstructure:"instructions"`
	Servers      []string `mapstructure:"servers"`
	IncludeTools []string `mapstructure:"include_tools"`
	ExcludeTools []string `mapstructure:"exclude_tools"`
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

type Logger struct {
	Type  string `mapstructure:"type"`
	Level string `mapstructure:"level"`
	Path  string `mapstructure:"path"`
}

func LoadConfig() (*AgentsConfig, error) {
	var agentsConfig AgentsConfig

	config := viper.New()
	config.SetConfigName("agents.config")
	config.SetConfigType("yaml")
	config.AddConfigPath(".")

	// LLM defaults
	config.SetDefault("openai.base_url", "")
	config.SetDefault("generic.base_url", "http://localhost:11434/v1")
	config.SetDefault("deepseek.base_url", "https://api.deepseek.com/v1")
	config.SetDefault("anthropic.base_url", "https://api.anthropic.com/v1/")
	config.SetDefault("openrouter.base_url", "https://openrouter.ai/api/v1/")
	config.SetDefault("google.base_url", "https://generativelanguage.googleapis.com/v1beta/openai/")

	// logger defaults
	config.SetDefault("logger.type", "console")
	config.SetDefault("logger.level", "warn")
	config.SetDefault("logger.path", "agent.jsonl")

	config.SetEnvPrefix("agents")
	// secrets.AllowEmptyEnv(true)
	config.AutomaticEnv()

	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file not found; ignore error if desired
			return nil, fmt.Errorf("error loading configs. %w", err)
		}
	}

	slog.Debug(fmt.Sprintf("loaded config file %s", config.ConfigFileUsed()))

	secrets := viper.New()
	secrets.SetConfigName("agents.secrets")
	secrets.SetConfigType("yaml")
	secrets.AddConfigPath(".")

	secrets.SetEnvPrefix("agents")
	// secrets.AllowEmptyEnv(true)
	secrets.AutomaticEnv()

	if err := secrets.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file not found; ignore error if desired
			return nil, fmt.Errorf("error loading secrets. %w", err)
		}
	}

	slog.Debug(fmt.Sprintf("loaded secrets file %s", secrets.ConfigFileUsed()))

	config.MergeConfigMap(secrets.AllSettings())

	if err := config.Unmarshal(&agentsConfig); err != nil {
		return nil, fmt.Errorf("error load agents.config.yaml. %w", err)
	}

	return &agentsConfig, nil
}
