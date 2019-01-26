package api

// Event corresponds to a single creat, update or delete request
type Event struct {
	Seq     uint64      `json:"seq"`
	Op      string      `json:"op"`
	Request interface{} `json:"request"`
}

// PullRequest structure
type PullRequest struct {
	Seq   uint64 `json:"seq"` // Sends the next seq number
	Count uint64 `json:"count"`
}

// PullResponse structure
type PullResponse struct {
	Events []Event `json:"events"`
	Err    string  `json:"err,omitempty"`
}
