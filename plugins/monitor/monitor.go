package monitor

import (
	"AegisGuard/pkg/plugin"
	"io"
	"log"
	"net/http"
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
	settings map[string]interface{}
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
	log.Println("Settings", m.settings)
	requestDetails := m.newRequestDetails(r)
	log.Printf("%+v\n\n", requestDetails)
	return nil
}

func (m *MonitorPlugin) GetType() string {
	return "monitor"
}

func (m *MonitorPlugin) GetSettings() map[string]interface{} {
	return m.settings
}

func (m *MonitorPlugin) GetMode() plugin.Mode {
	return plugin.Passive
}

func NewMonitorPlugin(settings map[string]interface{}) plugin.Plugin {
	return &MonitorPlugin{settings: settings}
}
