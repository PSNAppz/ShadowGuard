package main

import (
	configsPkg "AegisGuard/pkg/config"
	"AegisGuard/pkg/middleware/monitor"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	// Handle tasks
	r := mux.NewRouter()
	client := &http.Client{}
	handlers := make(map[string]http.HandlerFunc)
	config := configsPkg.InitConfig()

	for _, task := range config.Tasks {
		if task.Type == configsPkg.CAPTURE {
			monitorFunc := monitor.NewMonitorFunc(client, task.Method, task.Internal)
			r.HandleFunc(task.External, monitorFunc)
			handlers[task.External] = monitorFunc
		}
	}

	log.Printf("Listening on port %s\n", config.Port)
	http.ListenAndServe(config.Port, r)
}
