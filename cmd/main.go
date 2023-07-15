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
	config := conf.Init()

	for _, endpoint := range config.Endpoints {
		uri := config.Host + ":" + config.Port + endpoint.Internal
		for _, method := range endpoint.Methods {
			interceptFunc := middleware.Intercept(client, method, uri, endpoint.Plugins)
			r.HandleFunc(endpoint.External, interceptFunc).Methods(method)
		}
	}

	log.Printf("Listening on port %s\n", config.Port)
	http.ListenAndServe(config.Port, r)
}
