// Types that we can use for all content types (movie, tv, game, everything).
// Data responses to the client can use these "uniform" types to make access
// easier.

package domain

import "github.com/sbondCo/Watcharr/database/entity"

type MediaType string

const (
	MediaTypeTMDBMovie  = "tmdb_movie"
	MediaTypeTMDBShow   = "tmdb_tv"
	MediaTypeTMDBPerson = "tmdb_person"

	MediaTypeIGDBGame = "igdb_game"
)

type Media struct {
	// The type of media.
	Type MediaType
	// The ids associated with this media.
	IDs MediaIDs
	// The title/name of the media.
	Title string
	// A description.
	Summary string
	// The poster.
	Poster *entity.Image `json:"poster,omitempty"`
	// The external poster path.
	ExtPosterPath string
	// The rating.
	Rating uint
	// The amount of votes that made up the rating.
	RatingCount uint
}

type MediaIDs struct {
	// The internal ID
	// ID uint

	// For tmdb data
	TMDBID int
	IMDBID int

	// For igdb data
	IGDBID int
}
