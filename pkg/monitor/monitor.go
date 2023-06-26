package monitor

import (
	"io"
	"log"
	"net/http"
)

// Performs monitoring operation and prints out the details of the request
// TODO: Add more details to the log and store the same for future analysis
func Handle(r *http.Request, settings map[string]interface{}) {
	log.Println("Incoming Request Details")
	log.Println("Settings", settings)
	requestDetails := newRequestDetails(r)
	log.Printf("%+v\n\n", requestDetails)
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
