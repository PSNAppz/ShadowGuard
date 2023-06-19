package main

import (
	conf "AegisGuard/pkg/config"
	"AegisGuard/pkg/middleware/monitor"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	client := &http.Client{}
	handlers := make(map[string]http.HandlerFunc)
	config := conf.Init()

	for _, task := range config.Tasks {
		if task.Type == conf.MONITOR {
			uri := config.Host + config.Port + task.Internal
			monitorFunc := monitor.New(client, task.Method, uri)
			r.HandleFunc(task.External, monitorFunc)
			handlers[task.External] = monitorFunc
		}
	}

	log.Printf("Listening on port %s\n", config.Port)
	http.ListenAndServe(config.Port, r)
}
