package main

import (
	"log/slog"

	"github.com/sbondCo/Watcharr/game"
	"gorm.io/gorm"
)

type GameDetailsResponseWithPlayed struct {
	game.GameDetailsResponseBase
	SimilarGame []GameSimilarWithWatched `json:"similar_games"`
	WatchedAddedToContent
}

type GameSimilarWithWatched struct {
	game.GameSimilar
	WatchedAddedToContent
}

func gameDetailsAddWatched(
	db *gorm.DB,
	userId uint,
	content game.GameDetailsResponse,
) GameDetailsResponseWithPlayed {
	withWatchedResp := GameDetailsResponseWithPlayed{}
	withWatchedResp.GameDetailsResponseBase = content.GameDetailsResponseBase
	// Append watched list entry if exists
	if watchedEntry, err := getWatchedItemByIgdbId(db, userId, uint(content.ID)); err != nil {
		if err != gorm.ErrRecordNotFound {
			withWatchedResp.FailedToGetWatched = true
		}
	} else {
		withWatchedResp.Watched = &watchedEntry
	}
	// Add similar content with any watched entries
	similarContentIds := []int{}
	for _, v := range content.SimilarGame {
		withWatchedResp.SimilarGame = append(
			withWatchedResp.SimilarGame,
			GameSimilarWithWatched{
				GameSimilar: v,
			},
		)
		similarContentIds = append(similarContentIds, v.ID)
	}
	if ws, err := getWatchedItemsByIgdbIds(db, userId, similarContentIds); err == nil {
		for _, v := range ws {
			for i, vv := range withWatchedResp.SimilarGame {
				if vv.ID == v.Game.IgdbID {
					withWatchedResp.SimilarGame[i].WatchedAddedToContent.Watched = &v
				}
			}
		}
	} else {
		// TODO Set 'FailedToGetWatched' to `true` for the whole response obj when supported in structs
		slog.Error("Getting watched items by igdbIds failed!")
	}
	return withWatchedResp
}
