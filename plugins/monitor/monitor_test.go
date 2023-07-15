package monitor

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestMonitor(t *testing.T) {
	// Redirect log output to a buffer
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(io.Discard)
	}()

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	settings := map[string]interface{}{
		"verbose": true,
	}

	m := NewMonitorPlugin(settings, false)
	m.Handle(req)

	logged := buf.String()

	// Check if log output contains the expected strings
	expectedStrings := []string{
		"Incoming Request Details",
		"Settings map[verbose:true]",
		"Method:GET",
		"URL:http://example.com",
	}

	for _, s := range expectedStrings {
		if !strings.Contains(logged, s) {
			t.Errorf("expected log to contain %q, got:\n%s", s, logged)
		}
	}
}
