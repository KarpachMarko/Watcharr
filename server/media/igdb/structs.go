package igdb

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/sbondCo/Watcharr/domain"
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

// Website Type Enum
type WebsiteType int

const (
	WebsiteTypeOfficial WebsiteType = 1
	WebsiteTypeSteam    WebsiteType = 13
	WebsiteTypeGOG      WebsiteType = 17
)

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

func (t *GameSearchResponseResult) AsMedia() domain.Media {
	m := domain.Media{
		Type: domain.MediaTypeIGDBGame,
		IDs: domain.MediaIDs{
			IGDB: t.ID,
		},
		Name:          t.Name,
		Summary:       t.Summary,
		ExtPosterPath: t.Cover.ImageID,
		ReleaseDate:   t.FirstReleaseDate.Time,
	}
	return m
}

// This type is used for quite a few of our service funcs, maybe be a bit
// careful when updating it and separate out stuff if need be.
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

func (t *GameSimilar) AsMedia() domain.Media {
	m := domain.Media{
		Type: domain.MediaTypeIGDBGame,
		IDs: domain.MediaIDs{
			IGDB: t.ID,
		},
		Name:          t.Name,
		Summary:       t.Summary,
		ExtPosterPath: t.Cover.ImageID,
		ReleaseDate:   t.FirstReleaseDate.Time,
	}
	return m
}

// Details

type GameDetailsResponse struct {
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
		ID      int         `json:"id"`
		Type    WebsiteType `json:"type"`
		Trusted bool        `json:"trusted"`
		URL     string      `json:"url"`
	} `json:"websites"`

	SimilarGame []GameSimilar `json:"similar_games"`
}

func (t *GameDetailsResponse) AsMedia() domain.Media {
	m := domain.Media{
		Type: domain.MediaTypeIGDBGame,
		IDs: domain.MediaIDs{
			IGDB: t.ID,
		},
		Name:          t.Name,
		Summary:       t.Summary,
		ExtPosterPath: t.Cover.ImageID,
		ReleaseDate:   t.FirstReleaseDate.Time,
		Rating:        uint(t.Rating),
		RatingCount:   uint(t.RatingCount),
	}
	if len(t.Artworks) > 0 {
		m.ExtBackdropPath = t.Artworks[0].ImageID
	}
	// Process websites
	for _, v := range t.Websites {
		switch v.Type {
		case WebsiteTypeOfficial:
			m.Homepage = v.URL
		case WebsiteTypeSteam:
			m.Providers = append(m.Providers, domain.MediaProvider{
				Name: "Steam",
				Link: v.URL,
			})
		case WebsiteTypeGOG:
			m.Providers = append(m.Providers, domain.MediaProvider{
				Name: "GOG",
				Link: v.URL,
			})
		}
	}
	// Genres
	for _, v := range t.Genres {
		m.Genres = append(m.Genres, domain.MediaGenre{
			ID:   uint(v.ID),
			Name: v.Name,
		})
	}
	// Game modes
	for _, v := range t.GameModes {
		m.GameModes = append(m.GameModes, domain.MediaGenre{
			ID:   uint(v.ID),
			Name: v.Name,
		})
	}
	// Videos
	for _, v := range t.Videos {
		nameLower := strings.ToLower(v.Name)
		if !strings.Contains(nameLower, "trailer") {
			// Currently we only care about trailers
			continue
		}
		// Is best?
		isBest := false
		if nameLower == "trailer" || nameLower == "launch trailer" {
			isBest = true
		}
		m.Videos = append(m.Videos, domain.MediaVideo{
			ID:   v.VideoID,
			Name: v.Name,
			// Currently we only care about trailers
			Type: domain.MediaVideoTypeTrailer,
			Best: isBest,
		})
	}
	// Convert similar items to media too.
	for i := range t.SimilarGame {
		m.Similar = append(m.Similar, t.SimilarGame[i].AsMedia())
	}
	return m
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

// Popularity primitives where we just want the game id.

type PopularityPrimitivesGameIdsResponseResult struct {
	GameID int `json:"game_id"`
}

type PopularityPrimitivesGameIdsResponse []PopularityPrimitivesGameIdsResponseResult
