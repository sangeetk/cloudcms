package item

import (
	"golang.org/x/text/language"
)

const (
	// WidgetInput is input field
	WidgetInput = "input"
	// WidgetDate is date field
	WidgetDate = "date"
	// WidgetFile is file field
	WidgetFile = "file"
	// WidgetTextarea is textarea field
	WidgetTextarea = "textarea"
	// WidgetRichtext is richtext editor field
	WidgetRichtext = "richtext"
	// WidgetTags is tags field
	WidgetTags = "tags"

/*
	// WidgetCheckbox is checkbox field
	WidgetCheckbox = "checkbox"
	// WidgetRadio is radio field
	WidgetRadio = "radio"
	// WidgetSelect is select field
	WidgetSelect = "select"
	// WidgetSelectMultiple is select field with multiple values
	WidgetSelectMultiple = "selectmultiple"
*/
)

// Header contains some common header fields for content type
type Header struct {
	ID        uint64 `json:"id"`
	Language  string `json:"language"`
	Slug      string `json:"slug"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	DeletedAt int64  `json:"deleted_at"`
}

// Field represent a single field of content type
type Field struct {
	Name       string      `json:"name"`
	Heading    string      `json:"heading"`
	Widget     string      `json:"widget"`
	Helptext   string      `json:"helptext"`
	Value      interface{} `json:"value"`
	Editable   bool        `json:"editable"`
	UseAsSlug  bool        `json:"useasslug"`
	UseForSlug bool        `json:"useforslug"`
	FileType   string      `json:"filetype"`
	HasLabel   bool        `json:"has_label"`
	SkipHeader bool        `json:"skip_header"`
	SkipFooter bool        `json:"skip_footer"`
}

// Languages keep mapping between Types & Languages
var Languages map[string]language.Tag

// Types is a map used to reference a type name to its actual Editable type
// mainly for lookups in /admin route based utilities
var Types map[string]func() interface{}

// Fields contains the definition of item
var Fields map[string][]Field

// HeaderFields -
var HeaderFields = []Field{
	/*
		{Name: "ID", Widget: WidgetInput, Helptext: "ID", Value: "", Editable: false, UseForSlug: false},
		{Name: "Language", Widget: WidgetInput, Helptext: "Enter the Language here", Value: "", Editable: true, UseForSlug: false},
		{Name: "Slug", Widget: WidgetInput, Helptext: "Slug", Value: "", Editable: false, UseForSlug: false},
		{Name: "Status", Helptext: "Status", Widget: WidgetInput, Value: "", Editable: true, UseForSlug: false},
		{Name: "CreatedAt", Widget: WidgetInput, Helptext: "Created At", Value: "", Editable: false, UseForSlug: false},
		{Name: "UpdatedAt", Widget: WidgetInput, Helptext: "Updated At", Value: "", Editable: false, UseForSlug: false},
		{Name: "DeletedAt", Widget: WidgetInput, Helptext: "Deleted At", Value: "", Editable: false, UseForSlug: false},
	*/
}

func init() {
	Languages = make(map[string]language.Tag)
	Types = make(map[string]func() interface{})
	Fields = make(map[string][]Field)
}

// IndexContent enables Searching
func (h *Header) IndexContent() bool {
	return false
}
