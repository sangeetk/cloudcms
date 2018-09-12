package content

import (
	"time"
)

// Example content type
type Example struct {
	String  string    `json:"string"`
	Date    time.Time `json:"date"`
	Array   []string  `json:"array"`
	Integer int       `json:"integer"`
	Bool    bool      `json:"bool"`
}
