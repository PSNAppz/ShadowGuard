package main

import (
	"log"
	"net/http"
	conf "shadowguard/pkg/config"
	"shadowguard/pkg/middleware"

	"github.com/gorilla/mux"
)

var version = "undefined"

func main() {
	r := mux.NewRouter()
	client := &http.Client{}
	config := conf.Init()

	for _, endpoint := range config.Endpoints {
		uri := config.Host + config.Port + endpoint.Internal
		for _, method := range endpoint.Methods {
			interceptFunc := middleware.Intercept(client, method, uri, endpoint.Plugins)
			r.HandleFunc(endpoint.External, interceptFunc).Methods(method)
		}
	}

	log.Printf("Starting shadowguard version %s\n", version)
	log.Printf("Listening on port %s\n", config.Port)
	log.Fatal(http.ListenAndServe(config.Port, r))
}
