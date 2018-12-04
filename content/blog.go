package content

import (
	"fmt"
	"strings"

	"git.urantiatech.com/cloudcms/cloudcms/item"
)

// Blog content type
type Blog struct {
	item.Header

	Body     string `json:"body"`
	Category string `json:"category"`
}

func init() {
	item.Types[strings.ToLower("Blog")] = func() interface{} { return new(Blog) }
	item.Fields[strings.ToLower("Blog")] = append(item.HeaderFields, []item.Field{
		{Name: "Body", Type: "string", Label: "Enter the Body here", Widget: item.WidgetTextarea, Value: "This is body text"},
		{Name: "Category", Type: "@category", Label: "Select the Category here", Widget: item.WidgetSelect, Value: ""},
	}...)
}

// String defines how a Blog is printed. Update it using more descriptive
// fields from the Blog struct type
func (b *Blog) String() string {
	return fmt.Sprintf("This is Blog %s", b.Title)
}
