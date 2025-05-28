package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Transport string

const (
	TRANSPORT_HTTP  Transport = "http"
	TRANSPORT_STDIO Transport = "stdio"
	TRANSPORT_SSE   Transport = "sse"
)

type AgentsConfig struct {
	Model     string
	Agents    Agents
	MCPServer MCPServer

	OpenAI    OpenAI
	Anthropic Anthropic
	DeepSeek  Azure
}

type MCP struct {
	Servers map[string]MCPServer
}

type MCPServer struct {
	Transport    Transport
	url          string
	command      string
	args         []string
	headers      map[string]string
	environments map[string]string
}

type Agents struct {
	Agents map[string]Agent
}

type Agent struct {
}

type Anthropic struct {
	ApiKey   string
	BasePath string
}

type OpenAI struct {
	ApiKey   string
	BasePath string
}

type Azure struct {
	UseDefaultAzureCredential bool
	ApiKey                    string
	BasePath                  string
	ApiVersion                string
}

type DeepSeek struct {
	ApiKey   string
	BasePath string
}

type Google struct {
	ApiKey   string
	BasePath string
}

type Generic struct {
	ApiKey   string
	BasePath string
}

type OpenRouter struct {
	ApiKey   string
	BasePath string
}
type TensorZero struct {
	BasePath string
}

func LoadConfig() (*AgentsConfig, error) {
	var agentsConfig AgentsConfig
	config := viper.New()
	config.SetConfigName("agents.config")
	config.SetConfigType("yaml")
	config.AddConfigPath(".")

	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file not found; ignore error if desired
			return nil, fmt.Errorf("error loading configs. %w", err)
		}
	}
	config.Unmarshal(&agentsConfig)

	fmt.Printf("Config: %+v", agentsConfig)

	var agentsSecrets AgentsConfig

	secrets := viper.New()
	secrets.SetConfigName("agents.secrets")
	secrets.SetConfigType("yaml")
	secrets.AddConfigPath(".")

	if err := secrets.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file not found; ignore error if desired
			return nil, fmt.Errorf("error loading secrets. %w", err)
		}
	}

	secrets.AutomaticEnv()
	secrets.SetEnvPrefix("agents")
	secrets.AllowEmptyEnv(true)

	secrets.Unmarshal(&agentsSecrets)

	fmt.Printf("Secrets: %+v", agentsSecrets)

	config.MergeConfigMap(secrets.AllSettings())

	return &agentsConfig, nil
}
