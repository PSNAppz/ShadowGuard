package monitor

import (
	"io/ioutil"
	"log"
	"net/http"
)

// Read performs monitoring operation, contacts internal server and returns response to client
func New(client *http.Client, method, url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming request: %+v\n", r)
		defer r.Body.Close()
		req, err := http.NewRequest(method, url, r.Body)
		if err != nil {
			panic(err)
		}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		w.WriteHeader(resp.StatusCode)
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		log.Printf("Outgoing response: %+v\n", resp)
		w.Write(respBody)
	}
}
