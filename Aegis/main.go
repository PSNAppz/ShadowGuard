package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Task represents an individual Aegis operation for an endpoint
type Task struct {
	Type     string `json:"type"`
	Method   string `json:"method"`
	External string `json:"external"`
	Internal string `json:"internal"`
}

// Config represents the general configuration of Aegis
type Config struct {
	Host  string `json:"host"`
	Port  string `json:"port"`
	Tasks []Task `json:"tasks"`
}

// Read performs monitoring operation, contacts internal server and returns response to client
func newMonitorFunc(client *http.Client, method, url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("monitoring", url)
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			panic(err)
		}
		resp, err := client.Do(req)
		defer resp.Body.Close()
		if err != nil {
			panic(err)
		}

		w.WriteHeader(resp.StatusCode)
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", resp)
		w.Write(respBody)
	}
}

func main() {
	// Config file path can be set dynamically using environment variables.
	// The default is assumed to be `aegis.json` in the same directory.
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

	// Handle tasks
	r := mux.NewRouter()
	client := &http.Client{}
	handlers := make(map[string]http.HandlerFunc)
	for _, task := range config.Tasks {
		if task.Type == "monitor" {
			monitorFunc := newMonitorFunc(client, task.Method, config.Host+config.Port+task.Internal)
			r.HandleFunc(task.External, monitorFunc)
			handlers[task.External] = monitorFunc
		}
	}

	log.Printf("Listening on port %s\n", config.Port)
	http.ListenAndServe(config.Port, r)
}
