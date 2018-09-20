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

// Types is a map used to reference a type name to its actual Editable type
// mainly for lookups in /admin route based utilities
var Types map[string]func() interface{}

func init() {
	Types = make(map[string]func() interface{})
}

// IndexContent enables Searching
func (h *Header) IndexContent() bool {
	return false
}
