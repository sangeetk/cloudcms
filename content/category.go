package content

import (
	"fmt"
	"strings"

	"git.urantiatech.com/cloudcms/cloudcms/item"
)

// Category content type
type Category struct {
	item.Header
}

func init() {
	item.Types[strings.ToLower("Category")] = func() interface{} { return new(Category) }
	item.Fields[strings.ToLower("Category")] = append(item.HeaderFields, []item.Field{}...)
}

// String defines how a Category is printed. Update it using more descriptive
// fields from the Category struct type
func (c *Category) String() string {
	return fmt.Sprintf("This is Category %s", c.Title)
}
