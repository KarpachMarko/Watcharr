package domain

import (
	"time"

	"github.com/sbondCo/Watcharr/database/entity"
	"github.com/sbondCo/Watcharr/util"
)

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

type WatchedGetPageExtraProps struct {
	// Only get these watched ids.
	WatchedIds []int
}

type WatchedGetPageResponseResult struct {
	ID        uint                 `json:"id"`
	CreatedAt time.Time            `json:"createdAt"`
	UpdatedAt time.Time            `json:"updatedAt"`
	Status    entity.WatchedStatus `json:"status"`
	Rating    float64              `json:"rating"`
	Pinned    bool                 `json:"pinned"`
	Media     Media                `json:"media"`
}

func NewWatchedGetPageResponseResult(w *entity.Watched) WatchedGetPageResponseResult {
	r := WatchedGetPageResponseResult{
		ID:        w.ID,
		CreatedAt: w.CreatedAt,
		UpdatedAt: w.UpdatedAt,
		Status:    w.Status,
		Rating:    w.Rating,
		Pinned:    w.Pinned,
	}
	if w.Content != nil {
		r.Media = NewMediaFromContent(w.Content)
	} else if w.Game != nil {
		r.Media = NewMediaFromGame(w.Game)
	}
	return r
}

type WatchedGetPageResponse []WatchedGetPageResponseResult

func NewWatchedGetPageResponse(w []entity.Watched) WatchedGetPageResponse {
	r := WatchedGetPageResponse{}
	for _, v := range w {
		r = append(r, NewWatchedGetPageResponseResult(&v))
	}
	return r
}
