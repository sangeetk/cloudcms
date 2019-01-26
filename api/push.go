package api

// PushRequest structure
type PushRequest struct {
	Seq     uint64      `json:"seq"`
	Op      string      `json:"op"`
	Request interface{} `json:"request"`
}

// PushResponse structure
type PushResponse struct {
	Seq uint64 `json:"seq"` // Sends the next seq num
	Err string `json:"err,omitempty"`
}
