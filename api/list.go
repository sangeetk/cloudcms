package api

// ListRequest - list request
type ListRequest struct {
	Type     string `json:"type"`
	Language string `json:"language"`
	Status   string `json:"status"`
	Limit    int    `json:"limit"`
	Skip     int    `json:"skip"`
}

// ListResults - list results
type ListResults struct {
	Type     string        `json:"type"`
	Language string        `json:"language"`
	List     []interface{} `json:"list"`
	Total    int           `json:"total"`
	Limit    int           `json:"limit"`
	Skip     int           `json:"skip"`
	Err      string        `json:"err,omitempty"`
}
