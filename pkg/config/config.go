package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

// PluginConfig represents the configuration for a single plugin
type PluginConfig struct {
	Type     string                 `json:"type"`
	Settings map[string]interface{} `json:"settings"`
}

// Endpoint represents an external API Endpoint and its corresponding internal endpoint
type Endpoint struct {
	Plugins  []PluginConfig `json:"plugins"`
	Methods  []string       `json:"methods"`
	External string         `json:"external"`
	Internal string         `json:"internal"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

// Config represents the general configuration of ShadowGuard
type Config struct {
	Database  DatabaseConfig `json:"database"`
	Host      string         `json:"host"`
	Port      string         `json:"port"`
	Endpoints []Endpoint     `json:"endpoints"`
}

var (
	current *Config
	mu      sync.RWMutex
)

func loadConfigFromFile(path string) (*Config, error) {
	configJsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer configJsonFile.Close()

	byteData, err := io.ReadAll(configJsonFile)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(byteData, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func watchConfigFile(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Error creating watcher: %v", err)
		return
	}
	defer watcher.Close()

	dir := filepath.Dir(path)
	if err := watcher.Add(dir); err != nil {
		log.Printf("Error watching config directory: %v", err)
		return
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 && filepath.Clean(event.Name) == filepath.Clean(path) {
				log.Printf("Configuration file changed. Reloading\n")
				if cfg, err := loadConfigFromFile(path); err != nil {
					log.Printf("Failed to reload configuration: %v", err)
				} else {
					mu.Lock()
					*current = *cfg
					mu.Unlock()
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
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
	cfg, err := loadConfigFromFile(configFilePath)
	if err != nil {
		panic(err)
	}

	mu.Lock()
	current = cfg
	mu.Unlock()

	go watchConfigFile(configFilePath)

	log.Printf("Configuration file loaded.\n")
	return current
}
