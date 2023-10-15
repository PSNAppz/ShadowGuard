package monitor

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"shadowguard/pkg/database"
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

	var httpBuf bytes.Buffer
	httpBuf.WriteString("this is test data")
	req, err := http.NewRequest("GET", "http://example.com", &httpBuf)
	if err != nil {
		t.Fatal(err)
	}

	settings := map[string]interface{}{
		"verbose": true,
	}

	mockDB := database.NewMock()
	m := New(settings, mockDB)
	m.Handle(req)

	logged := buf.String()

	// Check if log output contains the expected strings
	expectedStrings := []string{
		`Incoming Request Details`,
		`"Method": "GET"`,
		`"URL": "http://example.com"`,
	}

	for _, s := range expectedStrings {
		if !strings.Contains(logged, s) {
			t.Errorf("expected log to contain %q, got:\n%s", s, logged)
		}
	}
}
