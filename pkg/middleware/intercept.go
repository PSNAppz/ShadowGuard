package middleware

import (
	"AegisGuard/pkg/config"
	"AegisGuard/pkg/monitor"
	"io"
	"log"
	"net/http"
)

// Read performs intercept operation, contacts internal server and returns response to client
func Intercept(client *http.Client, method, url string, tasks []config.TaskType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Execute tasks
		go executeTasks(tasks, r)
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
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		log.Printf("OUTGOING RESPONSE: %+v\n\n", resp)
		w.Write(respBody)
	}
}

func executeTasks(tasks []config.TaskType, r *http.Request) {
	for _, task := range tasks {
		switch task {
		case config.MONITOR:
			log.Println("Monitoring Started...")
			monitor.Listen(r)
		case config.CAPTURE:
			log.Println("Capturing Started...")
		default:
			log.Println("No Tasks to execute")
		}
	}
}
