package addedtocontent

import (
	"log/slog"

	"github.com/sbondCo/Watcharr/database/entity"
)

// This struct is for embedding inside content response structs.
// This holds the watched entry response data that will go along
// with the content responses.
type WatchedAddedToContent struct {
	// The related watched entry.
	Watched *entity.Watched `json:"watched,omitempty"`
	// If we failed to get the watched entry,
	// set this to true, so the frontend can
	// notify the user of why there is possibly
	// missing watched list data.
	FailedToGetWatched bool `json:"failedToGetWatched,omitempty"`
}

type Addable interface {
	AddWatched(w *entity.Watched)
	GetId() int
	GetMediaType() string
}

type WatchedProvider interface {
	GetWatchedItemsByTmdbIds(userId uint, c [][]any) ([]entity.Watched, error)
}

func AddWAC[S Addable](s []S, wp WatchedProvider, userId uint) error {
	contentIdAndTypePairs := [][]any{}
	for _, v := range s {
		contentIdAndTypePairs = append(contentIdAndTypePairs, []any{
			v.GetId(),
			entity.ContentType(v.GetMediaType()),
		})
	}
	if ws, err := wp.GetWatchedItemsByTmdbIds(userId, contentIdAndTypePairs); err == nil {
		for _, v := range ws {
			for _, vv := range s {
				if vv.GetId() == v.Content.TmdbID && vv.GetMediaType() == string(v.Content.Type) {
					vv.AddWatched(&v)
				}
			}
		}
	} else {
		// TODO Set 'FailedToGetWatched' to `true` for the whole response obj when supported in structs
		slog.Error("Getting watched items by tmdbIds failed!")
	}
	return nil
}
