package api

// Event corresponds to a single creat, update or delete request
type Event struct {
	SeqNum  int64       `json:"seq_num"`
	Op      string      `json:"op"`
	Request interface{} `json:"request"`
}

// PullRequest structure
type PullRequest struct {
	SeqnNum int64 `json:"seq_num"`
	Count   int64 `json:"count"`
}

// PullResponse structure
type PullResponse struct {
	Events []Event `json:"events"`
	Err    string  `json:"err,omitempty"`
}
