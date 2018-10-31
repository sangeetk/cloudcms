package content

import (
	"fmt"
	"strings"

	"git.urantiatech.com/cloudcms/cloudcms/item"
)

// Blog content type
type Blog struct {
	item.Header

	Body string `json:"body"`
}

func init() {
	item.Types[strings.ToLower("Blog")] = func() interface{} { return new(Blog) }
	item.Fields[strings.ToLower("Blog")] = append(item.HeaderFields, []item.Field{
		{Name: "Body", Type: "string", Label: "Enter the Body here", Widget: "input", Value: "This is body text"},
	}...)
}

// String defines how a Blog is printed. Update it using more descriptive
// fields from the Blog struct type
func (e *Blog) String() string {
	return fmt.Sprintf("This is Blog %s", e.Title)
}
