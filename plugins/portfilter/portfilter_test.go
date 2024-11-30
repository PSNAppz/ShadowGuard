package portfilter

import (
	"net/http"
	"net/http/httptest"
	"shadowguard/pkg/database"
	"testing"
)

func TestPortFilterPlugin(t *testing.T) {
	// Test 1: Port Blacklist
	settings := map[string]interface{}{
		"port-blacklist": []interface{}{22, 3306},
		"port-whitelist": []interface{}{},
		"active_mode":    true,
	}
	plugin := NewPortFilterPlugin(settings, database.NewMock()).(*PortFilterPlugin)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "localhost:22"

	err := plugin.Handle(req)
	if err == nil || err.Error() != "port is blacklisted" {
		t.Errorf("PortFilterPlugin did not block blacklisted port. Error: %v", err)
	}

	// Test 2: Port Whitelist
	settings = map[string]interface{}{
		"port-blacklist": []interface{}{},
		"port-whitelist": []interface{}{80, 443},
		"active_mode":    true,
	}
	plugin = NewPortFilterPlugin(settings, database.NewMock()).(*PortFilterPlugin)

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "localhost:80"

	err = plugin.Handle(req)
	if err != nil {
		t.Errorf("PortFilterPlugin blocked whitelisted port. Error: %v", err)
	}

	// Test 3: Default Allow Behavior (No Whitelist or Blacklist)
	settings = map[string]interface{}{
		"port-blacklist": []interface{}{},
		"port-whitelist": []interface{}{},
		"active_mode":    true,
	}
	plugin = NewPortFilterPlugin(settings, database.NewMock()).(*PortFilterPlugin)

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "localhost:8080"

	err = plugin.Handle(req)
	if err != nil {
		t.Errorf("PortFilterPlugin blocked port when no restrictions were configured. Error: %v", err)
	}

	// Test 4: Whitelist Restriction (Non-whitelisted port)
	settings = map[string]interface{}{
		"port-blacklist": []interface{}{},
		"port-whitelist": []interface{}{80, 443},
		"active_mode":    true,
	}
	plugin = NewPortFilterPlugin(settings, database.NewMock()).(*PortFilterPlugin)

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "localhost:8080"

	err = plugin.Handle(req)
	if err == nil || err.Error() != "port is not whitelisted" {
		t.Errorf("PortFilterPlugin did not block non-whitelisted port. Error: %v", err)
	}
}
