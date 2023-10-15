package database

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Request struct {
	gorm.Model
	Type             string
	Method           string
	URL              string
	Header           string
	Body             string
	Host             string
	RemoteAddr       string
	ContentLength    int64
	TransferEncoding pq.StringArray `gorm:"type:varchar(255)[];default:null"`
}

func (r Request) String() string {
	requestDetailsBytes, err := json.MarshalIndent(r, "", " ")
	if err != nil {
		panic(err)
	}

	return string(requestDetailsBytes)
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

func NewRequest(r *http.Request, pluginType string) (*Request, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return &Request{
		Type:             pluginType,
		Method:           r.Method,
		URL:              r.URL.String(),
		Header:           headerToString(r.Header),
		Body:             string(bodyBytes),
		Host:             r.Host,
		RemoteAddr:       r.RemoteAddr,
		ContentLength:    r.ContentLength,
		TransferEncoding: pq.StringArray(r.TransferEncoding),
	}, nil
}
