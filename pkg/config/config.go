package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

// PluginConfig represents the configuration for a single plugin
type PluginConfig struct {
	Type       string                 `json:"type"`
	Settings   map[string]interface{} `json:"settings"`
	ActiveMode bool                   `json:"active_mode"`
}

// Endpoint represents an external API Endpoint and its corresponding internal endpoint
type Endpoint struct {
	Plugins  []PluginConfig `json:"plugins"`
	Methods  []string       `json:"methods"`
	External string         `json:"external"`
	Internal string         `json:"internal"`
}

// Config represents the general configuration of ShadowGuard
type Config struct {
	Host      string     `json:"host"`
	Port      string     `json:"port"`
	Endpoints []Endpoint `json:"endpoints"`
}

// Init initializes the configuration from a file.
// The config file path can be set dynamically using environment variables.
// The default is assumed to be `config.json` in the same directory.
func Init() *Config {
	configFilePath := os.Getenv("SHADOW_CONFIG")
	if configFilePath == "" {
		configFilePath = "config.json"
	}

	log.Printf("Reading configuration file %s\n", configFilePath)
	configJsonFile, err := os.Open(configFilePath)
	if err != nil {
		panic(err)
	}
	defer configJsonFile.Close()
	byteData, err := io.ReadAll(configJsonFile)
	if err != nil {
		panic(err)
	}

	var config Config
	json.Unmarshal(byteData, &config)
	log.Printf("Configuration file loaded. %+v\n", config)
	return &config
}
