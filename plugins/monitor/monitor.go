package monitor

import (
	"log"
	"net/http"
	"shadowguard/pkg/database"
	"shadowguard/pkg/plugin"
	"shadowguard/pkg/publisher"
)

var Type string = "monitor"

func init() {
	plugin.RegisterPlugin(Type, New)
}

// New generates a new monitor plugin object
func New(settings map[string]interface{}, db database.DB) plugin.Plugin {
	publishers, err := publisher.CreatePublishers(settings)
	if err != nil {
		panic(err)
	}
	return &MonitorPlugin{Settings: settings, publishers: publishers, db: db}
}

// MonitorPlugin records incoming request and notifies various publishers that have been configured.
type MonitorPlugin struct {
	db         database.DB
	Settings   map[string]interface{}
	publishers []publisher.Publisher
}

func (m *MonitorPlugin) Type() string {
	return Type
}

func (m *MonitorPlugin) IsActiveMode() bool {
	return true
}

// Notify publishes the incoming request message as a string to all interested receivers.
func (m *MonitorPlugin) Notify(message string) {
	for _, p := range m.publishers {
		err := p.Publish(message)
		if err != nil {
			log.Printf("unable to notify publisher. message %s - error: %v", message, err)
		}
	}
}

// Handle formats the incoming request, optionally logs it, inserts into the database and notifies publishers.
func (m *MonitorPlugin) Handle(r *http.Request) error {
	requestModel, err := database.NewRequest(r, Type)
	if err != nil {
		return err
	}
	// handle verbose logging
	if verboseInterface, ok := m.Settings["verbose"]; ok {
		if verbose, ok := verboseInterface.(bool); ok && verbose {
			log.Println("Incoming Request Details")
			log.Println(requestModel)
		}
	}

	m.db.Insert(requestModel)
	m.Notify(requestModel.String())
	return nil
}
