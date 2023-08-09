package ipfilter

import (
	"errors"
	"log"
	"net"
	"net/http"
	"shadowguard/pkg/plugin"
	"shadowguard/pkg/publisher"
)

var Type string = "ipfilter"

func init() {
	plugin.RegisterPlugin(Type, NewIPFilterPlugin)
}

type IPFilterPlugin struct {
	Settings   map[string]interface{}
	activeMode bool
	blacklist  []interface{}
	whitelist  []interface{}
	publishers []publisher.Publisher
}

func NewIPFilterPlugin(settings map[string]interface{}) plugin.Plugin {
	publishers, err := publisher.CreatePublishers(settings)
	if err != nil {
		panic(err)
	}

	// Get the blacklist and whitelist from the settings
	blacklist, ok := settings["blacklist"].([]interface{})
	if !ok {
		panic("expected 'blacklist' to be a slice of interfaces")
	}

	whitelist, ok := settings["whitelist"].([]interface{})
	if !ok {
		panic("expected 'whitelist' to be a slice of interfaces")
	}

	return &IPFilterPlugin{
		Settings:   settings,
		activeMode: settings["active_mode"].(bool),
		blacklist:  blacklist,
		whitelist:  whitelist,
		publishers: publishers,
	}
}

func (p *IPFilterPlugin) Type() string {
	return Type
}

func (p *IPFilterPlugin) IsActiveMode() bool {
	return p.activeMode
}

func (p *IPFilterPlugin) Notify(message string) {
	for _, p := range p.publishers {
		err := p.Publish(message)
		if err != nil {
			log.Printf("unable to notify publisher. message %s - error: %v", message, err)
		}
	}
}

func (p *IPFilterPlugin) Handle(r *http.Request) error {
	// Extract the IP from the request
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return err
	}

	// Check if the IP is on the blacklist
	for _, blacklistedIP := range p.blacklist {
		if ip == blacklistedIP {
			return errors.New("IP address is blacklisted")
		}
	}

	// Return nil if the IP is on the whitelist and if it is not
	return nil
}
