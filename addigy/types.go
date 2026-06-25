package addigy

import "encoding/json"

// ErrorResponse is the standard error body returned by the Addigy API
// (response_entities.ErrorResponse).
type ErrorResponse struct {
	Code       int             `json:"code,omitempty"`
	Message    string          `json:"message,omitempty"`
	ErrorChain json.RawMessage `json:"error_chain,omitempty"`
}

// Metadata is the pagination block returned by list/search endpoints
// (response_entities.Metadata).
type Metadata struct {
	Page        int `json:"page"`
	PageCount   int `json:"page_count"`
	PerPage     int `json:"per_page"`
	ResultCount int `json:"result_count"`
	Total       int `json:"total"`
}
