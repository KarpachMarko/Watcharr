package entity

import "time"

// For storing cached games, so we can serve the basic local data for watched list to work
type Game struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UpdatedAt time.Time `json:"updatedAt"`
	IgdbID    int       `json:"igdbId" gorm:"uniqueIndex;not null"`
	Name      string    `json:"name"`
	CoverID   string    `json:"coverId"`
	Summary   string    `json:"summary"`
	Storyline string    `json:"storyline"`
	// First release date
	ReleaseDate *time.Time `json:"releaseDate,omitempty"`
	Rating      float64    `json:"rating"`
	RatingCount int        `json:"ratingCount"`
	Status      int        `json:"status"`
	Category    int        `json:"category"`
	// Arrays turned to strings that may be useful
	GameModes string `json:"gameModes"`
	Genres    string `json:"genres"`
	Platforms string `json:"platforms"`
	// Id to poster image row (cached game cover)
	PosterID *uint  `json:"-"`
	Poster   *Image `json:"poster,omitempty"`
}
