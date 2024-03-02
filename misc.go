package payarc

type Metadata struct {
	Include    []string   `json:"include"`
	Custom     []any      `json:"custom"`
	Pagination Pagination `json:"pagination"`
}

type PaginationLinks struct {
	Previous string `json:"previous,omitempty"`
	Next     string `json:"next,omitempty"`
}

type Pagination struct {
	Total       int             `json:"total"`
	Count       int             `json:"count"`
	PerPage     int             `json:"per_page"`
	CurrentPage int             `json:"current_page"`
	TotalPages  int             `json:"total_pages"`
	Links       PaginationLinks `json:"links"`
}
