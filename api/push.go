package api

// PushRequest structure
type PushRequest struct {
	SeqNum  int64       `json:"seq_num"`
	Op      string      `json:"op"`
	Request interface{} `json:"request"`
}

// PushResponse structure
type PushResponse struct {
	SeqNum int64  `json:"seq_num"`
	Err    string `json:"err,omitempty"`
}
