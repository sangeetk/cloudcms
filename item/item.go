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

// Item is a generic content type
type Item map[string]interface{}
