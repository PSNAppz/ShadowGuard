package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type TaskType string

const (
	MONITOR TaskType = "monitor"
	CAPTURE TaskType = "capture"
)

// Task represents an individual Aegis operation for an endpoint
type Task struct {
	Type     TaskType `json:"type"`
	Method   string   `json:"method"`
	External string   `json:"external"`
	Internal string   `json:"internal"`
}

// Config represents the general configuration of Aegis
type Config struct {
	Host  string `json:"host"`
	Port  string `json:"port"`
	Tasks []Task `json:"tasks"`
}

// Config file path can be set dynamically using environment variables.
// The default is assumed to be `aegis.json` in the same directory.
func Init() *Config {
	configFilePath := os.Getenv("AEGIS_CONFIG")
	if configFilePath == "" {
		configFilePath = "aegis.json"
	}

	log.Printf("Reading configuration file %s\n", configFilePath)
	configJsonFile, err := os.Open(configFilePath)
	if err != nil {
		panic(err)
	}
	defer configJsonFile.Close()
	byteData, err := ioutil.ReadAll(configJsonFile)
	if err != nil {
		panic(err)
	}

	var config Config
	json.Unmarshal(byteData, &config)
	log.Printf("Configuration file loaded. %+v\n", config)
	return &config
}
