package domain

import (
	"github.com/sbondCo/Watcharr/database/entity"
	"github.com/sbondCo/Watcharr/util"
)

// Get watched page request extra (GET) options.
// Since this is user input, validity of string types cannot be guaranteed.
type WatchedGetPageRequest struct {
	// Sorting type.
	Sort WatchedSort `form:"sort"`
	// Sorting direction (asc or desc).
	SortDir SortDirection `form:"sortDir,default=desc"`
	// Filtering options.
	FilterType   []util.SupportedMedia  `form:"type" collection_format:"csv"`
	FilterStatus []entity.WatchedStatus `form:"status" collection_format:"csv"`
}

type WatchedSort string

const (
	WatchedSortDateAdded    WatchedSort = "DATEADDED"
	WatchedSortLastChanged  WatchedSort = "LASTCHANGED"
	WatchedSortLastFinished WatchedSort = "LASTFIN"
	WatchedSortRating       WatchedSort = "RATING"
	WatchedSortAlphabetical WatchedSort = "ALPHA"
)

type SortDirection string

const (
	WatchedSortDirAsc  SortDirection = "asc"
	WatchedSortDirDesc SortDirection = "desc"
)
