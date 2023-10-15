package middleware

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"shadowguard/pkg/config"
	"shadowguard/pkg/database"
	"strings"
	"testing"
)

func TestInterceptWithActivePlugins(t *testing.T) {
	// create internal endpoint (or API endpoint)
	internalServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		w.Write([]byte("Internal server reached."))
	}))
	defer internalServer.Close()

	internalClient := internalServer.Client()
	defer internalServer.CloseClientConnections()

	method := "GET"
	pluginConfigs := []config.PluginConfig{
		{
			Type: "requestfilter",
			Settings: map[string]interface{}{
				"active_mode":      true,
				"ip-blacklist":     []interface{}{"127.0.0.1"},
				"ip-whitelist":     []interface{}{},
				"region-whitelist": []interface{}{},
				"region-blacklist": []interface{}{},
			},
		},
	}
	// create external endpoint (or client facing endpoint)
	externalServer := httptest.NewServer(Intercept(internalClient, method, internalServer.URL, pluginConfigs, database.NewMock()))
	defer externalServer.Close()

	externalClient := externalServer.Client()
	defer externalServer.CloseClientConnections()

	// Initiate request to external client
	buf := bytes.NewBufferString("this is test data")
	req := httptest.NewRequest(http.MethodGet, externalServer.URL, buf)
	req.RemoteAddr = "127.0.0.1:80" // httptest creates 127.0.0.1:random_port
	req.RequestURI = ""
	req.Header["Connection"] = []string{"Keep-Alive"}
	resp, err := externalClient.Do(req)
	if err != nil {
		t.Errorf("External client encountered an error while making request. %+v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Unable to read response body from external client. %+v", err)
	}

	respBodyStr := strings.TrimSpace(string(respBody))
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Incorrect response code for blocked IP. Found %d", resp.StatusCode)
	} else if respBodyStr != "Request blocked by plugin: IP address is blacklisted" {
		t.Errorf("Unexpected response body when IP is blacklisted. Found %s", respBodyStr)
	}

	if connectionHeaders, found := resp.Request.Header["Connection"]; found {
		if connectionHeaders[0] != "Keep-Alive" {
			t.Errorf("Invalid Connection header found. Found %s, expected %s", connectionHeaders[0], "Keep-Alive")
		}
	} else {
		t.Errorf("Connection header is missing from the request")
	}
}

func TestInterceptWithPassivePlugins(t *testing.T) {
	// create internal endpoint (or API endpoint)
	internalServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		w.Write([]byte("Internal server reached."))
	}))
	defer internalServer.Close()

	internalClient := internalServer.Client()
	defer internalServer.CloseClientConnections()

	method := "GET"
	pluginConfigs := []config.PluginConfig{
		{
			Type:     "monitor",
			Settings: map[string]interface{}{"verbose": true},
		},
	}
	// create external endpoint (or client facing endpoint)
	externalServer := httptest.NewServer(Intercept(internalClient, method, internalServer.URL, pluginConfigs, database.NewMock()))
	defer externalServer.Close()

	externalClient := externalServer.Client()
	defer externalServer.CloseClientConnections()

	buf := bytes.NewBufferString("this is test data")
	// Initiate request to external client
	req := httptest.NewRequest(http.MethodGet, externalServer.URL, buf)
	req.RequestURI = ""
	req.Header["Connection"] = []string{"Keep-Alive"}
	resp, err := externalClient.Do(req)
	if err != nil {
		t.Errorf("External client encountered an error while making request. %+v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Unable to read response body from external client. %+v", err)
	}

	respBodyStr := strings.TrimSpace(string(respBody))

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Incorrect response code. Found %d", resp.StatusCode)
	} else if respBodyStr != "Internal server reached." {
		t.Errorf("Incorrect response returned. Found %s", respBodyStr)
	}

	if connectionHeaders, found := resp.Request.Header["Connection"]; found {
		if connectionHeaders[0] != "Keep-Alive" {
			t.Errorf("Invalid Connection header found. Found %s, expected %s", connectionHeaders[0], "Keep-Alive")
		}
	} else {
		t.Errorf("Connection header is missing from the request")
	}
}
