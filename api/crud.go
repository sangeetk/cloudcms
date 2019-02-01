package api

const (
	// CreateOp opration
	CreateOp = "create"
	// ReadOp operation
	ReadOp = "read"
	// UpdateOp operation
	UpdateOp = "update"
	// DeleteOp operation
	DeleteOp = "delete"
)

// CreateRequest structure
type CreateRequest struct {
	Type     string      `json:"type"`
	Language string      `json:"language"`
	Slug     string      `json:"slug"`
	Content  interface{} `json:"content"`
}

// UpdateRequest structure
type UpdateRequest struct {
	Type     string      `json:"type"`
	Language string      `json:"language"`
	Slug     string      `json:"slug"`
	Content  interface{} `json:"content"`
}

// DeleteRequest structure
type DeleteRequest struct {
	Type     string `json:"type"`
	Language string `json:"language"`
	Slug     string `json:"slug"`
}

// ReadRequest structure
type ReadRequest struct {
	Type     string `json:"type"`
	Language string `json:"language"`
	Slug     string `json:"slug"`
}

// Response structure
type Response struct {
	Type     string      `json:"type"`
	Language string      `json:"language"`
	Content  interface{} `json:"content"`
	Err      string      `json:"err,omitempty"`
}
