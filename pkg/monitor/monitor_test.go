package monitor

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestMonitor(t *testing.T) {
	var str bytes.Buffer
	log.SetOutput(&str)

	req, err := http.NewRequest("GET", "http://test.com", nil)
	if err != nil {
		t.Errorf("Unable to create new request. %+v", err)
	}

	Listen(req)

	details := newRequestDetails(req)
	if !strings.Contains(str.String(), fmt.Sprintf("%+v", details)) {
		t.Errorf("Unable to find formatted request. Found %s", str.String())
	}
}
