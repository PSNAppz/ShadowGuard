package monitor

import (
	"io"
	"log"
	"net/http"
	"shadowguard/pkg/plugin"
	"shadowguard/pkg/receiver"
)

// Register this plugin in the plugin package
func init() {
	plugin.RegisterPlugin("monitor", NewMonitorPlugin)
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

type MonitorPlugin struct {
	Settings   map[string]interface{}
	ActiveMode bool
	Receivers  []receiver.NotificationReceiver
}

func (m *MonitorPlugin) newRequestDetails(r *http.Request) RequestDetails {
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

func (m *MonitorPlugin) Handle(r *http.Request) error {
	log.Println("Incoming Request Details")
	log.Println("Settings", m.Settings)
	requestDetails := m.newRequestDetails(r)
	log.Printf("%+v\n\n", requestDetails)
	return nil
}

func (m *MonitorPlugin) GetType() string {
	return "monitor"
}

func (m *MonitorPlugin) GetSettings() map[string]interface{} {
	return m.Settings
}

func (m *MonitorPlugin) IsActiveMode() bool {
	return m.ActiveMode
}

func (m *MonitorPlugin) Notify(message string) {
	for _, r := range m.Receivers {
		err := r.Notify(message)
		if err != nil {
			log.Printf("unable to notify receiver. message %s - error: %v", message, err)
		}
	}
}

func NewMonitorPlugin(settings map[string]interface{}, activeMode bool) plugin.Plugin {
	receivers, err := receiver.CreateReceivers(settings)
	if err != nil {
		panic(err)
	}

	m := &MonitorPlugin{Settings: settings, ActiveMode: activeMode, Receivers: receivers}

	return m
}
