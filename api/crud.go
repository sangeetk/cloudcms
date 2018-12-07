package api

const (
	// Create opration
	Create = "create"
	// Read operation
	Read = "read"
	// Update operation
	Update = "update"
	// Delete operation
	Delete = "delete"
)

// CreateRequest structure
type CreateRequest struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

// UpdateRequest structure
type UpdateRequest struct {
	Type    string      `json:"type"`
	Slug    string      `json:"slug"`
	Content interface{} `json:"content"`
}

// DeleteRequest structure
type DeleteRequest struct {
	Type string `json:"type"`
	Slug string `json:"slug"`
}

// ReadRequest structure
type ReadRequest struct {
	Type string `json:"type"`
	Slug string `json:"slug"`
}

// Response structure
type Response struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
	Err     string      `json:"err,omitempty"`
}
