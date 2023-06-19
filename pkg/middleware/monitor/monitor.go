package monitor

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Read performs monitoring operation, contacts internal server and returns response to client
func New(client *http.Client, method, url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("monitoring", url)
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
		fmt.Printf("%+v\n", resp)
		w.Write(respBody)
	}
}
