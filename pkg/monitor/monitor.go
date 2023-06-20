package monitor

import (
	"io"
	"log"
	"net/http"
)

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

// Performs monitoring operation and prints out the details of the request
// TODO: Add more details to the log and store the same for future analysis
func Listen(r *http.Request) {
	log.Println("Incoming Request Details")

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

	log.Printf("%+v\n\n", requestDetails)

}
