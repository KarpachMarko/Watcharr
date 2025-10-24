package entity

import (
	"time"
)

type ContentType string

const (
	MOVIE ContentType = "movie"
	SHOW  ContentType = "tv"
	// Show episode
	SHOW_EPISODE ContentType = "tv_episode"
)

// For storing cached content, so we can serve the basic local data for watched list to work
type Content struct {
	ID               int         `json:"id" gorm:"primaryKey;autoIncrement"`
	TmdbID           int         `json:"tmdbId" gorm:"uniqueIndex:contentidtotypeidx;not null"`
	Title            string      `json:"title"`
	PosterPath       string      `json:"poster_path"`
	Overview         string      `json:"overview"`
	Type             ContentType `json:"type" gorm:"uniqueIndex:contentidtotypeidx;not null"`
	ReleaseDate      *time.Time  `json:"release_date,omitempty"`
	Popularity       float32     `json:"popularity"`
	VoteAverage      float32     `json:"vote_average"`
	VoteCount        uint32      `json:"vote_count"`
	ImdbID           string      `json:"imdb_id"`
	Status           string      `json:"status"`
	Budget           uint32      `json:"budget"`
	Revenue          uint32      `json:"revenue"`
	Runtime          uint32      `json:"runtime"`
	NumberOfEpisodes uint32      `json:"numberOfEpisodes"`
	NumberOfSeasons  uint32      `json:"numberOfSeasons"`
}
