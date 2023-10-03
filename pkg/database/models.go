package database

import (
	"encoding/json"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Request struct {
	gorm.Model
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
