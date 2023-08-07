package ipfilter

import (
	"errors"
	"log"
	"net"
	"net/http"
	"shadowguard/pkg/plugin"
	"shadowguard/pkg/receiver"
)

type IPFilterPlugin struct {
	Settings   map[string]interface{}
	ActiveMode bool
	Receivers  []receiver.NotificationReceiver
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

func (p *IPFilterPlugin) IsActiveMode() bool {
	return p.ActiveMode
}

func (p *IPFilterPlugin) Notify(message string) {
	for _, r := range p.Receivers {
		err := r.Notify(message)
		if err != nil {
			log.Printf("unable to notify receiver. message %s - error: %v", message, err)
		}
	}
}

func NewIPFilterPlugin(settings map[string]interface{}, activeMode bool) plugin.Plugin {
	receivers, err := receiver.CreateReceivers(settings)
	if err != nil {
		panic(err)
	}

	return &IPFilterPlugin{
		Settings:   settings,
		ActiveMode: activeMode,
		Receivers:  receivers,
	}
}

func init() {
	plugin.RegisterPlugin("ipfilter", NewIPFilterPlugin)
}
