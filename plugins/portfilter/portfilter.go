package portfilter

import (
	"errors"
	"log"
	"net/http"
	"shadowguard/pkg/database"
	"shadowguard/pkg/plugin"
	"shadowguard/pkg/publisher"
	"strconv"
	"strings"
)

var Type string = "portfilter"

func init() {
	plugin.RegisterPlugin(Type, NewPortFilterPlugin)
}

type PortFilterPlugin struct {
	db            database.DB
	Settings      map[string]interface{}
	activeMode    bool
	portBlacklist []interface{}
	portWhitelist []interface{}
	publishers    []publisher.Publisher
}

// NewPortFilterPlugin initializes the PortFilterPlugin.
func NewPortFilterPlugin(settings map[string]interface{}, db database.DB) plugin.Plugin {
	publishers, err := publisher.CreatePublishers(settings)
	if err != nil {
		panic(err)
	}

	portBlacklist, _ := settings["port-blacklist"].([]interface{})
	portWhitelist, _ := settings["port-whitelist"].([]interface{})

	return &PortFilterPlugin{
		db:            db,
		Settings:      settings,
		activeMode:    settings["active_mode"].(bool),
		portBlacklist: portBlacklist,
		portWhitelist: portWhitelist,
		publishers:    publishers,
	}
}

func (p *PortFilterPlugin) Type() string {
	return Type
}

func (p *PortFilterPlugin) IsActiveMode() bool {
	return p.activeMode
}

func (p *PortFilterPlugin) Notify(message string) {
	for _, pub := range p.publishers {
		err := pub.Publish(message)
		if err != nil {
			log.Printf("unable to notify publisher. message %s - error: %v", message, err)
		}
	}
}

func (p *PortFilterPlugin) Handle(r *http.Request) error {
	// Extract port from the request's host field
	hostPort := r.Host
	port := 80 // Default HTTP port

	if parts := strings.Split(hostPort, ":"); len(parts) == 2 {
		var err error
		port, err = strconv.Atoi(parts[1])
		if err != nil {
			return errors.New("invalid port in request")
		}
	}

	// Check port against blacklist
	for _, blacklistedPort := range p.portBlacklist {
		if int(blacklistedPort.(int)) == port {
			req, err := database.NewRequest(r, "portblacklist")
			if err != nil {
				print("ERROR")
				println(err)
				return err
			}
			p.db.Insert(req)
			p.Notify("Port is blacklisted: " + strconv.Itoa(port))
			return errors.New("port is blacklisted")
		}
	}

	// Check port against whitelist if defined
	if len(p.portWhitelist) > 0 {
		isWhitelisted := false
		for _, whitelistedPort := range p.portWhitelist {
			if int(whitelistedPort.(int)) == port {
				req, err := database.NewRequest(r, "portwhitelist")

				if err != nil {
					return err
				}
				p.db.Insert(req)
				isWhitelisted = true
				break
			}
		}
		if !isWhitelisted {
			p.Notify("Port is not whitelisted: " + strconv.Itoa(port))
			return errors.New("port is not whitelisted")
		}
	}

	return nil
}
