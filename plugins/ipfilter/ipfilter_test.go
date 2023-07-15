package ipfilter

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIPFilterPlugin(t *testing.T) {
	// Create a new IPFilterPlugin with a blacklist
	plugin := NewIPFilterPlugin(map[string]interface{}{
		"blacklist": []interface{}{"127.0.0.1"},
		"whitelist": []interface{}{},
	}, true)

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	req.RemoteAddr = "127.0.0.1:80"
	// Test that the IPFilterPlugin blocks the request when the IP is blacklisted
	err := plugin.Handle(req)
	if err == nil || err.Error() != "IP address is blacklisted" {
		t.Errorf("IPFilterPlugin did not block blacklisted IP. Error: %v", err)
	}

	// Update the plugin settings to include a whitelist
	plugin = NewIPFilterPlugin(map[string]interface{}{
		"blacklist": []interface{}{},
		"whitelist": []interface{}{"127.0.0.1"},
	}, true)

	// Test that the IPFilterPlugin allows the request when the IP is whitelisted
	err = plugin.Handle(req)
	if err != nil {
		t.Errorf("IPFilterPlugin blocked whitelisted IP. Error: %v", err)
	}

	// Update the plugin settings to include neither a whitelist nor a blacklist
	plugin = NewIPFilterPlugin(map[string]interface{}{
		"blacklist": []string{},
		"whitelist": []string{},
	}, true)

	// Test that the IPFilterPlugin allows the request when the IP is neither whitelisted nor blacklisted
	err = plugin.Handle(req)
	if err != nil {
		t.Errorf("IPFilterPlugin blocked IP not on whitelist or blacklist. Error: %v", err)
	}
}
