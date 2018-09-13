package content

import (
	"time"

	"git.urantiatech.com/cloudcms/cloudcms/item"
)

// Example content type
type Example struct {
	item.Header

	String  string    `json:"string"`
	Date    time.Time `json:"date"`
	Array   []string  `json:"array"`
	Integer int       `json:"integer"`
	Bool    bool      `json:"bool"`
}
