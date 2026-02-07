package igdb

import (
	"encoding/json"
	"time"

	"github.com/sbondCo/Watcharr/database/entity"
	"github.com/sbondCo/Watcharr/domain"
	"github.com/sbondCo/Watcharr/util"
)

type TwitchTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// So we can unmarshall the unix timestamps returned from igdb into time.Time.
type UnixTime struct {
	time.Time
}

func (u *UnixTime) UnmarshalJSON(b []byte) error {
	var timestamp int64
	err := json.Unmarshal(b, &timestamp)
	if err != nil {
		return err
	}
	u.Time = time.Unix(timestamp, 0)
	return nil
}

// Only the fields we request included in each struct

// Search

type GameSearchResponseResult struct {
	ID    int `json:"id"`
	Cover struct {
		ID      int    `json:"id"`
		ImageID string `json:"image_id"`
	} `json:"cover"`
	FirstReleaseDate UnixTime `json:"first_release_date"`
	Name             string   `json:"name"`
	Summary          string   `json:"summary,omitempty"`
	VersionTitle     string   `json:"version_title,omitempty"`
}

func (t GameSearchResponseResult) GetId() int {
	return t.ID
}

func (t GameSearchResponseResult) GetMediaType() util.SupportedMedia {
	return util.SupportedMediaGame
}

func (t *GameSearchResponseResult) AsMedia() domain.Media {
	m := domain.Media{
		Type: domain.MediaTypeIGDBGame,
		IDs: domain.MediaIDs{
			IGDB: t.ID,
		},
		Name:          t.Name,
		Summary:       t.Summary,
		ExtPosterPath: t.Cover.ImageID,
	}
	return m
}

type GameSearchResponse []GameSearchResponseResult

// Similar

type GameSimilar struct {
	ID               int      `json:"id"`
	Name             string   `json:"name"`
	Summary          string   `json:"summary"`
	FirstReleaseDate UnixTime `json:"first_release_date"`
	Cover            struct {
		ID      int    `json:"id"`
		ImageID string `json:"image_id"`
	} `json:"cover"`
}

func (t GameSimilar) GetId() int {
	return t.ID
}

func (t GameSimilar) GetMediaType() util.SupportedMedia {
	return util.SupportedMediaGame
}

type GameSimilarWithWatched struct {
	GameSimilar
	Watched *entity.Watched `json:"watched,omitempty"`
}

// Details

type GameDetailsResponseBase struct {
	ID       int `json:"id"`
	Artworks []struct {
		Width   int    `json:"width"`
		Height  int    `json:"height"`
		ImageID string `json:"image_id"`
	} `json:"artworks"`
	Category int `json:"category"`
	Cover    struct {
		ID      int    `json:"id"`
		ImageID string `json:"image_id"`
	} `json:"cover"`
	FirstReleaseDate UnixTime `json:"first_release_date"`
	GameModes        []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"game_modes"`
	Genres []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
	InvolvedCompanies []struct {
		ID      int `json:"id"`
		Company struct {
			ID          int    `json:"id"`
			Description string `json:"description"`
			Name        string `json:"name"`
			Slug        string `json:"slug"`
			Websites    []struct {
				ID       int    `json:"id"`
				Category int    `json:"category"`
				Trusted  bool   `json:"trusted"`
				URL      string `json:"url"`
			} `json:"websites"`
		} `json:"company"`
		Developer  bool `json:"developer"`
		Porting    bool `json:"porting"`
		Publisher  bool `json:"publisher"`
		Supporting bool `json:"supporting"`
	} `json:"involved_companies"`
	Name      string `json:"name"`
	Platforms []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"platforms"`
	Rating      float64 `json:"rating"`
	RatingCount int     `json:"rating_count"`
	Summary     string  `json:"summary"`
	Storyline   string  `json:"storyline"`
	Status      int     `json:"status"`
	URL         string  `json:"url"`
	Videos      []struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		VideoID string `json:"video_id"`
	} `json:"videos"`
	Websites []struct {
		ID       int    `json:"id"`
		Category int    `json:"category"`
		Trusted  bool   `json:"trusted"`
		URL      string `json:"url"`
	} `json:"websites"`
}

type GameDetailsResponse struct {
	GameDetailsResponseBase
	SimilarGame []GameSimilar `json:"similar_games"`
}

func (t GameDetailsResponseBase) GetId() int {
	return t.ID
}

func (t GameDetailsResponseBase) GetMediaType() util.SupportedMedia {
	return util.SupportedMediaGame
}

type GameDetailsResponseWithWatched struct {
	GameDetailsResponseBase
	SimilarGame []GameSimilarWithWatched `json:"similar_games"`
	Watched     *entity.Watched          `json:"watched,omitempty"`
}

// Basic Details

type GameDetailsBasicResponse struct {
	ID       int `json:"id"`
	Category int `json:"category"`
	Cover    struct {
		ImageID string `json:"image_id"`
	} `json:"cover"`
	FirstReleaseDate UnixTime `json:"first_release_date"`
	GameModes        []struct {
		Name string `json:"name"`
	} `json:"game_modes"`
	Genres []struct {
		Name string `json:"name"`
	} `json:"genres"`
	Name      string `json:"name"`
	Platforms []struct {
		Name string `json:"name"`
	} `json:"platforms"`
	Rating      float64 `json:"rating"`
	RatingCount int     `json:"rating_count"`
	Summary     string  `json:"summary"`
	Storyline   string  `json:"storyline"`
	Status      int     `json:"status"`
}
