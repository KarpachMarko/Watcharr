package domain

import (
	"github.com/sbondCo/Watcharr/feature/watched/addedtocontent"
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

type SearchResultType string

const (
	SearchResultTypeTMDB = "tmdb"
	SearchResultTypeIGDB = "igdb"
)

type SearchRequest struct {
	// The type of content we are searching for.
	// SearchTypeMulti encompasses all types of media in the results.
	Type SearchType `form:"type"`
	// The search term.
	Query string `form:"query"`
}

type SearchResult struct {
	// The type of result.
	// Since there are cases where multiple data types can be returned
	// from different external apis, each result will be wrapped in an
	// identifiable type so the client knows what to expect in each result.
	Type SearchResultType `json:"type"`
	// The actual result.
	Data addedtocontent.Addable `json:"data"`
}

type SearchResponse struct {
	util.PaginationResponse[SearchResult]
}
