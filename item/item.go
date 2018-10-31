package item

// Header contains some common header fields for content type
type Header struct {
	ID        uint64 `json:"id"`
	Title     string `json:"title"`
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
	Editable bool        `json:"bool"`
}

// Types is a map used to reference a type name to its actual Editable type
// mainly for lookups in /admin route based utilities
var Types map[string]func() interface{}

// Fields contains the definition of item
var Fields map[string][]Field

// HeaderFields -
var HeaderFields = []Field{
	{Name: "ID", Type: "integer", Label: "ID", Widget: "input", Value: "", Editable: false},
	{Name: "Title", Type: "string", Label: "Enter the Title here", Widget: "input", Value: "", Editable: true},
	{Name: "Slug", Type: "string", Label: "Slug", Widget: "input", Value: "", Editable: false},
	{Name: "Status", Type: "string", Label: "Status", Widget: "input", Value: "", Editable: true},
	{Name: "CreatedAt", Type: "date", Label: "Created At", Widget: "input", Value: "", Editable: false},
	{Name: "UpdatedAt", Type: "date", Label: "Updated At", Widget: "input", Value: "", Editable: false},
	{Name: "DeletedAt", Type: "date", Label: "Deleted At", Widget: "input", Value: "", Editable: false},
}

func init() {
	Types = make(map[string]func() interface{})
	Fields = make(map[string][]Field)
}

// IndexContent enables Searching
func (h *Header) IndexContent() bool {
	return false
}
