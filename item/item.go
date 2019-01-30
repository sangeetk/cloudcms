package item

const (
	// WidgetInput is input field
	WidgetInput = "input"
	// WidgetFile is file field
	WidgetFile = "file"
	// WidgetTextarea is textarea field
	WidgetTextarea = "textarea"
	// WidgetRichtext is richtext editor field
	WidgetRichtext = "richtext"
	// WidgetCheckbox is checkbox field
	WidgetCheckbox = "checkbox"
	// WidgetRadio is radio field
	WidgetRadio = "radio"
	// WidgetSelect is select field
	WidgetSelect = "select"
	// WidgetSelectMultiple is select field with multiple values
	WidgetSelectMultiple = "selectmultiple"
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
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Label    string      `json:"label"`
	Widget   string      `json:"widget"`
	Value    interface{} `json:"value"`
	Editable bool        `json:"editable"`
}

// Types is a map used to reference a type name to its actual Editable type
// mainly for lookups in /admin route based utilities
var Types map[string]func() interface{}

// Fields contains the definition of item
var Fields map[string][]Field

// HeaderFields -
var HeaderFields = []Field{
	{Name: "ID", Type: "integer", Label: "ID", Widget: WidgetInput, Value: "", Editable: false},
	{Name: "Title", Type: "string", Label: "Enter the Title here", Widget: WidgetInput, Value: "", Editable: true},
	{Name: "Slug", Type: "string", Label: "Slug", Widget: WidgetInput, Value: "", Editable: false},
	{Name: "Status", Type: "string", Label: "Status", Widget: WidgetInput, Value: "", Editable: true},
	{Name: "CreatedAt", Type: "date", Label: "Created At", Widget: WidgetInput, Value: "", Editable: false},
	{Name: "UpdatedAt", Type: "date", Label: "Updated At", Widget: WidgetInput, Value: "", Editable: false},
	{Name: "DeletedAt", Type: "date", Label: "Deleted At", Widget: WidgetInput, Value: "", Editable: false},
}

func init() {
	Types = make(map[string]func() interface{})
	Fields = make(map[string][]Field)
}

// IndexContent enables Searching
func (h *Header) IndexContent() bool {
	return false
}
