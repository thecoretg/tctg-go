package huntress

// ErrorResponse is the standard error body returned by the Huntress API.
type ErrorResponse struct {
	Message string `json:"message,omitempty"`
}

// Pagination is the cursor block included on every paginated list response.
// An empty NextPageToken indicates the last page.
type Pagination struct {
	NextPageURL   string `json:"next_page_url,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
}

// Address is a mailing address used for account billing and shipping.
type Address struct {
	Line1      string `json:"line1,omitempty"`
	Line2      string `json:"line2,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	Country    string `json:"country,omitempty"`
}

// User is a minimal representation of a Huntress user, embedded in other
// resources such as remediation approvals.
type User struct {
	ID    int64  `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
}
