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

// Add a watched entry request
type WatchedAddRequest struct {
	// Type of content we are adding to watched.
	ContentType util.SupportedMedia `json:"contentType" binding:"required,oneof=movie tv game"`
	// ID of content from tmdb (if ContentType is movie or tv).
	TMDBID int `json:"tmdbId"`
	// ID of content from igdb (if ContentType is game).
	IGDBID int `json:"igdbId"`

	Status   entity.WatchedStatus `json:"status"`
	Rating   float64              `json:"rating" binding:"max=10"`
	Thoughts string               `json:"thoughts"`
	// Pass a watched date and we will set the CreatedAt (and initial UpdatedAt)
	// properties for this watched entry to this specific date.
	WatchedDate time.Time `json:"watchedDate,omitempty"`
}

// Update watched entry request
type WatchedUpdateRequest struct {
	Status         entity.WatchedStatus `json:"status" binding:"required_without_all=Rating Thoughts RemoveThoughts Pinned"`
	Rating         float64              `json:"rating" binding:"max=10,required_without_all=Status Thoughts RemoveThoughts Pinned"`
	Thoughts       string               `json:"thoughts" binding:"required_without_all=Status Rating RemoveThoughts Pinned"`
	RemoveThoughts bool                 `json:"removeThoughts"`
	Pinned         *bool                `json:"pinned" binding:"required_without_all=Status Rating Thoughts RemoveThoughts"`
}

// Update response.
type WatchedUpdateResponse struct {
	NewActivity entity.Activity `json:"newActivity"`
}

// Removal response.
type WatchedRemoveResponse struct {
	NewActivity entity.Activity `json:"newActivity"`
}
