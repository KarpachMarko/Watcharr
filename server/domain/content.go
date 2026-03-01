package domain

import "time"

type ContentDetailsResponse struct {
	Media
}

type PersonCreditsResponse struct {
	Credits []Media `json:"credits,omitempty"`
}

type PersonDetailsResponse struct {
	Name               string    `json:"name,omitempty"`
	Birthday           time.Time `json:"birthday,omitzero"`
	Deathday           time.Time `json:"deathday,omitzero"`
	Age                int       `json:"age,omitempty"`
	PlaceOfBirth       string    `json:"placeOfBirth,omitempty"`
	KnownForDepartment string    `json:"knownForDepartment,omitempty"`
	Biography          string    `json:"biography,omitempty"`
	ExtPosterPath      string    `json:"extPosterPath,omitempty"`
	Homepage           string    `json:"homepage,omitempty"`
}
