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
	item.Fields[strings.ToLower("Example")] = append(item.HeaderFields, []item.Field{
		{Name: "Body", Type: "string", Label: "Enter the Body here", Widget: "input", Value: "This is body text"},
		{Name: "Date", Type: "date", Label: "Enter the Date here", Widget: "input", Value: time.Now()},
		{Name: "Array", Type: "[]string", Label: "Enter Array here", Widget: "select", Value: []string{"first", "second", "third"}},
		{Name: "Integer", Type: "int", Label: "Enter the Integer here", Widget: "input", Value: 100},
		{Name: "Bool", Type: "bool", Label: "Please enter true or false", Widget: "checkbox", Value: true},
	}...)
}

// String defines how a Example is printed. Update it using more descriptive
// fields from the Example struct type
func (e *Example) String() string {
	return fmt.Sprintf("This is Example %s", e.Title)
}
