package domain

import (
	"time"

	"github.com/sbondCo/Watcharr/database/entity"
)

type ImportResponseType string

var (
	// Successful import
	IMPORT_SUCCESS ImportResponseType = "IMPORT_SUCCESS"
	// Import failed for reasons user cant fix
	IMPORT_FAILED ImportResponseType = "IMPORT_FAILED"
	// Import query returned multiple results, user must decide
	IMPORT_MULTI ImportResponseType = "IMPORT_MULTI"
	// Import query returned zero results, user must provide more info
	IMPORT_NOTFOUND ImportResponseType = "IMPORT_NOTFOUND"
	// Item already exists so couldn't import (unique constraint hit when adding)
	IMPORT_EXISTS ImportResponseType = "IMPORT_EXISTS"
)

type ImportRequest struct {
	Name             string                  `json:"name"`
	Year             int                     `json:"year"`
	TmdbID           int                     `json:"tmdbId"`
	Type             entity.ContentType      `json:"type"`
	Rating           float64                 `json:"rating" binding:"max=10"`
	RatingCustomDate *time.Time              `json:"ratingCustomDate"`
	Status           entity.WatchedStatus    `json:"status"`
	Thoughts         string                  `json:"thoughts"`
	DatesWatched     []time.Time             `json:"datesWatched"`
	Activity         []entity.Activity       `json:"activity"`
	WatchedEpisodes  []entity.WatchedEpisode `json:"watchedEpisodes"`
	WatchedSeason    []entity.WatchedSeason  `json:"watchedSeasons"`
	Tags             []TagAddRequest         `json:"tags"`
	ImdbID           string                  `json:"imdbId"`
}

type ImportResponse struct {
	Type    ImportResponseType `json:"type"`
	Results []Media            `json:"results,omitempty"`
	Match   Media              `json:"match,omitzero"`
	// On success this will be filled with the new watched entry
	WatchedEntry entity.Watched `json:"watchedEntry,omitzero"`
}
