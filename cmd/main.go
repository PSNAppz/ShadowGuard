package main

import (
	conf "AegisGuard/pkg/config"
	"AegisGuard/pkg/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	client := &http.Client{}
	handlers := make(map[string]http.HandlerFunc)
	config := conf.Init()

	for _, endpoint := range config.Endpoints {
		uri := config.Host + config.Port + endpoint.Internal
		monitorFunc := middleware.Intercept(client, endpoint.Method, uri, endpoint.Tasks)
		r.HandleFunc(endpoint.External, monitorFunc)
		handlers[endpoint.External] = monitorFunc
	}
	log.Printf("Listening on port %s\n", config.Port)
	http.ListenAndServe(config.Port, r)
}
