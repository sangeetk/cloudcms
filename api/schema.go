package api

import (
	"git.urantiatech.com/cloudcms/cloudcms/item"
)

// SchemaRequest - schema request
type SchemaRequest struct {
	Type string `json:"type"`
}

// ContentType is the type of content stored in cms
type ContentType struct {
	Fields []item.Field `json:"fields"`
}

// SchemaResponse - schema response
type SchemaResponse struct {
	Schema map[string]ContentType `json:"schema,omitempty"`
	Err    string                 `json:"err,omitempty"`
}
