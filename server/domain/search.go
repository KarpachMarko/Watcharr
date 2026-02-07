package domain

import (
	"github.com/go-playground/validator/v10"
	"github.com/sbondCo/Watcharr/util"
)

type SearchType string

const (
	// Search for **all available media types**.
	SearchTypeMulti = "multi"
	// Search for a **movie**.
	SearchTypeMovie = "movie"
	// Search for a **show**.
	SearchTypeShow = "show"
	// Search for a **person** (actor).
	SearchTypePerson = "person"
	// Search for a **game**.
	SearchTypeGame = "game"
)

type SearchRequest struct {
	// The type of content we are searching for.
	// SearchTypeMulti encompasses all types of media in the results.
	Type SearchType `form:"type" binding:"validsearchtype"`
	// The search term.
	Query string `form:"query"`
}

type SearchResponse struct {
	util.PaginationResponse[Media]
}

var ValidSearchType validator.Func = func(fl validator.FieldLevel) bool {
	st, ok := fl.Field().Interface().(SearchType)
	if ok {
		switch st {
		case SearchTypeMulti,
			SearchTypeMovie,
			SearchTypeShow,
			SearchTypePerson,
			SearchTypeGame:
			return true
		}
	}
	return false
}
