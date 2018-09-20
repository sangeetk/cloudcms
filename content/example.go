package content

import (
	"fmt"
	"strings"
	"time"

	"git.urantiatech.com/cloudcms/cloudcms/item"
)

// Example content type
type Example struct {
	item.Header

	Body    string    `json:"body"`
	Date    time.Time `json:"date"`
	Array   []string  `json:"array"`
	Integer int       `json:"integer"`
	Bool    bool      `json:"bool"`
}

func init() {
	item.Types[strings.ToLower("Example")] = func() interface{} { return new(Example) }
}

// String defines how a Example is printed. Update it using more descriptive
// fields from the Example struct type
func (e *Example) String() string {
	return fmt.Sprintf("This is Example %s", e.Title)
}
