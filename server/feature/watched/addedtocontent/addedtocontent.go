package addedtocontent

import (
	"log/slog"

	"github.com/sbondCo/Watcharr/database/entity"
	"gorm.io/gorm"
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
	GetId() int
	GetMediaType() entity.ContentType
}

type WatchedProvider interface {
	GetWatchedItemsByTmdbIds(userId uint, c [][]any) ([]entity.Watched, error)
	GetWatchedItemByTmdbId(userId uint, tmdbId uint, contentType entity.ContentType) (entity.Watched, error)
}

type AddListCall[S Addable] struct {
	s     []S
	addCb func(i int, w *entity.Watched)
}

func NewAddListCall[S Addable](s []S, addCb func(i int, w *entity.Watched)) *AddListCall[S] {
	return &AddListCall[S]{
		s,
		addCb,
	}
}

// A helper for adding `Watched` data to structs of generic data.
// The actual adding of watched data will be handled by the caller
// through the `addCb` callback function.
func AddList[S Addable](
	wp WatchedProvider,
	userId uint,
	s []S,
	addCb func(i int, w *entity.Watched),
) error {
	// TODO Check len of s
	contentIdAndTypePairs := [][]any{}
	for _, v := range s {
		contentIdAndTypePairs = append(contentIdAndTypePairs, []any{
			v.GetId(),
			v.GetMediaType(),
		})
	}
	if ws, err := wp.GetWatchedItemsByTmdbIds(userId, contentIdAndTypePairs); err == nil {
		for _, v := range ws {
			for i, vv := range s {
				if vv.GetId() == v.Content.TmdbID && vv.GetMediaType() == v.Content.Type {
					addCb(i, &v)
				}
			}
		}
	} else {
		// TODO Set 'FailedToGetWatched' to `true` for the whole response obj when supported in structs
		slog.Error("Getting watched items by tmdbIds failed!")
	}
	return nil
}

func Add[S Addable](
	wp WatchedProvider,
	userId uint,
	s S,
	addCb func(w *entity.Watched),
) {
	watchedEntry, err := wp.GetWatchedItemByTmdbId(
		userId,
		uint(s.GetId()),
		s.GetMediaType(),
	)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			// withWatchedResp.FailedToGetWatched = true
		}
		return
	}
	addCb(&watchedEntry)
}

func AddSingularAndList[S Addable, S2 Addable](
	wp WatchedProvider,
	userId uint,
	s S,
	addCb func(w *entity.Watched),
	list []*AddListCall[S2],
) {
	Add(wp, userId, s, addCb)
	for i := range list {
		AddList(wp, userId, list[i].s, list[i].addCb)
	}
}
