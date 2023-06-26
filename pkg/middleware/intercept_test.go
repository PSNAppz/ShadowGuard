package middleware

import (
	"AegisGuard/pkg/config"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIntercept(t *testing.T) {
	// create internal endpoint (or API endpoint)
	// this will be contacted by the Intercept function
	internalServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		w.Write([]byte("Internal server reached."))
	}))
	defer internalServer.Close()

	internalClient := internalServer.Client()
	defer internalServer.CloseClientConnections()

	method := "GET"
	tasks := []config.TaskType{
		config.MONITOR,
	}
	// create external endpoint (or client facing endpoint)
	// this will be contacted using the client
	externalServer := httptest.NewServer(Intercept(internalClient, method, internalServer.URL, tasks))
	defer externalServer.Close()

	externalClient := externalServer.Client()
	defer externalServer.CloseClientConnections()

	// Initiate request to external client
	resp, err := externalClient.Get(externalServer.URL)
	if err != nil {
		t.Errorf("External client encountered an error while making request. %+v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Unable to read response body from external client. %+v", err)
	}

	if string(respBody) != "Internal server reached." {
		t.Errorf("Incorrect response returned. Found %s", respBody)
	}

	if resp.StatusCode != 201 {
		t.Errorf("Incorrect response code. Found %d", resp.StatusCode)
	}
}