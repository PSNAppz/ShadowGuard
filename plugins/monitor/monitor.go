package monitor

import (
	"io"
	"log"
	"net/http"
	"shadowguard/pkg/database"
	"shadowguard/pkg/plugin"
	"shadowguard/pkg/publisher"

	"github.com/lib/pq"
)

var Type string = "monitor"

func init() {
	plugin.RegisterPlugin(Type, New)
}

// New generates a new monitor plugin object
func New(settings map[string]interface{}, db *database.Database) plugin.Plugin {
	publishers, err := publisher.CreatePublishers(settings)
	if err != nil {
		panic(err)
	}
	return &MonitorPlugin{Settings: settings, publishers: publishers, db: db}
}

// MonitorPlugin records incoming request and notifies various publishers that have been configured.
type MonitorPlugin struct {
	db         *database.Database
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

func headerToString(header http.Header) string {
	var result string

	// Iterate through the header fields
	for key, values := range header {
		for _, value := range values {
			// Concatenate the key and value into the result string
			result += key + ": " + value + "\n"
		}
	}

	return result
}

// Handle formats the incoming request, optionally logs it, inserts into the database and notifies publishers.
func (m *MonitorPlugin) Handle(r *http.Request) error {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	reqData := database.Request{
		Method:           r.Method,
		URL:              r.URL.String(),
		Header:           headerToString(r.Header),
		Body:             string(bodyBytes),
		Host:             r.Host,
		RemoteAddr:       r.RemoteAddr,
		ContentLength:    r.ContentLength,
		TransferEncoding: pq.StringArray(r.TransferEncoding),
	}

	// handle verbose logging
	if verboseInterface, ok := m.Settings["verbose"]; ok {
		if verbose, ok := verboseInterface.(bool); ok && verbose {
			log.Println("Incoming Request Details")
			log.Println(reqData)
		}
	}

	m.db.Insert(&reqData)
	m.Notify(reqData.String())
	return nil
}
