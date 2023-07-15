package ipfilter

import (
	"AegisGuard/pkg/plugin"
	"errors"
	"log"
	"net"
	"net/http"
)

type IPFilterPlugin struct {
	Settings map[string]interface{}
}

func (p *IPFilterPlugin) Handle(r *http.Request) error {
	// Extract the IP from the request
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	// Get the blacklist and whitelist from the settings
	blacklist, ok := p.Settings["blacklist"].([]interface{})
	if !ok {
		log.Println("Error: expected 'blacklist' to be a slice of interfaces")
	}
	whitelist, ok := p.Settings["whitelist"].([]interface{})
	if !ok {
		log.Println("Error: expected 'whitelist' to be a slice of interfaces")
	}

	// Check if the IP is on the blacklist
	for _, blacklistedIP := range blacklist {
		if ip == blacklistedIP {
			return errors.New("IP address is blacklisted")
		}
	}

	// Check if the IP is on the whitelist
	for _, whitelistedIP := range whitelist {
		if ip == whitelistedIP {
			// If the IP is on the whitelist, do nothing and let the request continue
			return nil
		}
	}

	// If the IP is not on either list, allow the request to continue
	return nil
}

func (p *IPFilterPlugin) GetType() string {
	return "ipfilter"
}

func (p *IPFilterPlugin) GetSettings() map[string]interface{} {
	return p.Settings
}

func (p *IPFilterPlugin) GetMode() plugin.Mode {
	return plugin.Active
}

func NewIPFilterPlugin(settings map[string]interface{}) plugin.Plugin {
	return &IPFilterPlugin{Settings: settings}
}

func init() {
	plugin.RegisterPlugin("ipfilter", NewIPFilterPlugin)
}