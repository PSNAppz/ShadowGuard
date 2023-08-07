package monitor

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"shadowguard/pkg/plugin"
	"shadowguard/pkg/receiver"
)

var Type string = "monitor"

type MonitorPlugin struct {
	Settings  map[string]interface{}
	receivers []receiver.NotificationReceiver
}

// Register this plugin in the plugin package
func init() {
	plugin.RegisterPlugin(Type, NewMonitorPlugin)
}

func (m *MonitorPlugin) Type() string {
	return Type
}

func (m *MonitorPlugin) IsActiveMode() bool {
	return true
}

func (m *MonitorPlugin) Notify(message string) {
	for _, r := range m.receivers {
		err := r.Notify(message)
		if err != nil {
			log.Printf("unable to notify receiver. message %s - error: %v", message, err)
		}
	}
}

func NewMonitorPlugin(settings map[string]interface{}) plugin.Plugin {
	receivers, err := receiver.CreateReceivers(settings)
	if err != nil {
		panic(err)
	}

	return &MonitorPlugin{Settings: settings, receivers: receivers}
}

func (m *MonitorPlugin) Handle(r *http.Request) error {
	log.Println("Incoming Request Details")
	requestDetails := newRequestDetails(r)
	log.Println(requestDetails)
	return nil
}

type RequestDetails struct {
	Method           string
	URL              string
	Header           http.Header
	Body             io.ReadCloser
	Host             string
	RemoteAddr       string
	ContentLength    int64
	TransferEncoding []string
}

func (r RequestDetails) String() string {
	requestDetailsBytes, err := json.MarshalIndent(r, "", " ")
	if err != nil {
		panic(err)
	}

	return string(requestDetailsBytes)
}

func newRequestDetails(r *http.Request) RequestDetails {
	requestDetails := RequestDetails{
		Method:           r.Method,
		URL:              r.URL.String(),
		Header:           r.Header,
		Body:             r.Body,
		Host:             r.Host,
		RemoteAddr:       r.RemoteAddr,
		ContentLength:    r.ContentLength,
		TransferEncoding: r.TransferEncoding,
	}
	return requestDetails
}
