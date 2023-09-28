package requestfilter

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestFilterPlugin(t *testing.T) {
	// Mock the function to always return "US" for our tests

	// Test 1: IP Blacklist
	plugin := NewRequestFilterPlugin(map[string]interface{}{
		"active_mode":      true,
		"ip-blacklist":     []interface{}{"127.0.0.1"},
		"ip-whitelist":     []interface{}{},
		"region-whitelist": []interface{}{},
		"region-blacklist": []interface{}{},
	})

	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	req.RemoteAddr = "127.0.0.1:80"

	err := plugin.Handle(req)
	if err == nil || err.Error() != "IP address is blacklisted" {
		t.Errorf("RequestFilterPlugin did not block blacklisted IP. Error: %v", err)
	}

	// Test 2: IP Whitelist
	plugin = NewRequestFilterPlugin(map[string]interface{}{
		"active_mode":      true,
		"ip-blacklist":     []interface{}{},
		"ip-whitelist":     []interface{}{"127.0.0.1"},
		"region-whitelist": []interface{}{},
		"region-blacklist": []interface{}{},
	})

	err = plugin.Handle(req)
	if err != nil {
		t.Errorf("RequestFilterPlugin blocked whitelisted IP. Error: %v", err)
	}

	// Test 3: Region Whitelist (should block since "US" is not in whitelist)
	plugin = NewRequestFilterPlugin(map[string]interface{}{
		"active_mode":      true,
		"ip-blacklist":     []interface{}{},
		"ip-whitelist":     []interface{}{},
		"region-whitelist": []interface{}{"US"},
		"region-blacklist": []interface{}{},
	})
	// Mock the IP to return "US" for our tests
	req.RemoteAddr = "115.240.90.163:80" // Indian IP sample
	err = plugin.Handle(req)
	if err == nil || err.Error() != "access from this region is not whitelisted" {
		t.Errorf("RequestFilterPlugin did not block non-whitelisted region. Error: %v", err)
	}

	// Test 4: No Restrictions
	plugin = NewRequestFilterPlugin(map[string]interface{}{
		"active_mode":      true,
		"ip-blacklist":     []interface{}{},
		"ip-whitelist":     []interface{}{},
		"region-whitelist": []interface{}{},
		"region-blacklist": []interface{}{},
	})

	err = plugin.Handle(req)
	if err != nil {
		t.Errorf("RequestFilterPlugin blocked IP without restrictions. Error: %v", err)
	}
}
