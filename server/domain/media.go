// Types that we can use for all content types (movie, tv, game, everything).
// Data responses to the client can use these "uniform" types to make access
// easier.

package domain

import (
	"time"

	"github.com/sbondCo/Watcharr/database/entity"
	"github.com/sbondCo/Watcharr/util"
)

type MediaType string

const (
	MediaTypeTMDBMovie  = "tmdb_movie"
	MediaTypeTMDBShow   = "tmdb_tv"
	MediaTypeTMDBPerson = "tmdb_person"

	MediaTypeIGDBGame = "igdb_game"
)

type Media struct {
	// The type of media.
	Type MediaType `json:"type,omitempty"`
	// The ids associated with this media.
	IDs MediaIDs `json:"ids"`
	// The name of the media.
	Name string `json:"name,omitempty"`
	// A description.
	Summary string `json:"summary,omitempty"`
	// The poster.
	Poster *entity.Image `json:"poster,omitempty"`
	// The external poster path.
	ExtPosterPath string `json:"extPosterPath,omitempty"`
	// The rating.
	Rating uint `json:"rating,omitempty"`
	// The amount of votes that made up the rating.
	RatingCount uint `json:"ratingCount,omitempty"`
	// Watched data.
	Watched *entity.Watched `json:"watched,omitempty"`
	// Similar media.
	Similar []Media `json:"similar,omitempty"`
	// Release date / first air date.
	ReleaseDate time.Time `json:"releaseDate,omitzero"`

	//
	// Properties that are less important (not used for all responses).
	//

	// Backdrop path.
	ExtBackdropPath string `json:"extBackdropPath,omitempty"`
	// Genres.
	Genres []MediaGenre `json:"genres,omitempty"`
	// Media website.
	Homepage string `json:"homepage,omitempty"`
	// Trailer video.
	Trailer string `json:"trailer,omitempty"`

	//
	// Properties only for movies/tv.
	//

	// Runtime.
	Runtime        uint                `json:"runtime,omitempty"`
	WatchProviders *MediaWatchProvider `json:"watchProviders,omitempty"`

	//
	// Properties only for Games
	//

	// Game modes.
	GameModes []MediaGenre `json:"gameModes,omitempty"`
}

func (t Media) GetId() int {
	switch t.Type {
	case MediaTypeTMDBMovie,
		MediaTypeTMDBShow:
		return t.IDs.TMDB
	case MediaTypeIGDBGame:
		return t.IDs.IGDB
	}
	return -99
}

func (t Media) GetMediaType() util.SupportedMedia {
	switch t.Type {
	case MediaTypeTMDBMovie:
		return util.SupportedMediaMovie
	case MediaTypeTMDBShow:
		return util.SupportedMediaShow
	case MediaTypeIGDBGame:
		return util.SupportedMediaGame
	}
	// Unsupported...
	return ""
}

type MediaIDs struct {
	// The internal ID
	// Watcharr uint

	// For tmdb data
	TMDB int `json:"tmdb,omitempty"`
	IMDB int `json:"imdb,omitempty"`

	// For igdb data
	IGDB int `json:"igdb,omitempty"`
}

type MediaGenre struct {
	// ID of the genre on the external database.
	ID uint `json:"id,omitempty"`
	// Name of genre.
	Name uint `json:"name,omitempty"`
}

type MediaWatchProvider struct {
	// Name of the provider.
	Name string `json:"name,omitempty"`

// Converter for Content (tv/movie) entity to Media
func NewMediaFromContent(c *entity.Content) Media {
	m := Media{
		IDs: MediaIDs{
			TMDB: c.TmdbID,
		},
		Name:          c.Title,
		Summary:       c.Overview,
		ExtPosterPath: c.PosterPath,
		Rating:        uint(c.VoteAverage),
		RatingCount:   uint(c.VoteCount),
		Runtime:       uint(c.Runtime),
	}
	switch c.Type {
	case entity.MOVIE:
		m.Type = MediaTypeTMDBMovie
	case entity.SHOW:
		m.Type = MediaTypeTMDBShow
	}
	if c.ReleaseDate != nil {
		m.ReleaseDate = *c.ReleaseDate
	}
	return m
}

// Converter for Game entity to Media
func NewMediaFromGame(c *entity.Game) Media {
	m := Media{
		IDs: MediaIDs{
			IGDB: c.IgdbID,
		},
		Type:          MediaTypeIGDBGame,
		Name:          c.Name,
		Summary:       c.Summary,
		Poster:        c.Poster,
		ExtPosterPath: c.CoverID,
		Rating:        uint(c.Rating),
		RatingCount:   uint(c.RatingCount),
	}
	if c.ReleaseDate != nil {
		m.ReleaseDate = *c.ReleaseDate
	}
	return m
}
