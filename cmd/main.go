package main

import (
	"log"
	"net/http"
	conf "shadowguard/pkg/config"
	"shadowguard/pkg/database"
	"shadowguard/pkg/middleware"

	"github.com/gorilla/mux"
)

var version = "undefined"

func main() {
	r := mux.NewRouter()
	client := &http.Client{}
	config := conf.Init()

	conn, err := database.New(config.Database)
	if err != nil {
		panic(err)
	}

	for _, endpoint := range config.Endpoints {
		uri := config.Host + config.Port + endpoint.Internal
		for _, method := range endpoint.Methods {
			interceptFunc := middleware.Intercept(client, method, uri, endpoint.Plugins, conn)
			r.HandleFunc(endpoint.External, interceptFunc).Methods(method)
		}
	}

	// represents the API of the customer using ShadowGuard, this could be hosted anywhere. It just needs to be able to be reached.
	r.HandleFunc("/our-customers-api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("reached"))
	})

	log.Printf("Starting ShadowGuard v %s\n", version)
	log.Printf("Listening on port %s\n", config.Port)
	log.Fatal(http.ListenAndServe(config.Port, r))
}
