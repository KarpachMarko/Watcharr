package domain

type ContentDetailsResponse struct {
	Media
}

type PersonCreditsResponse struct {
	Credits []Media `json:"credits,omitempty"`
}
