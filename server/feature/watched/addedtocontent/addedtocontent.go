package addedtocontent

import (
	"log/slog"

	"github.com/sbondCo/Watcharr/database/entity"
	"gorm.io/gorm"
)

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
		slog.Error("AddList: Getting watched items by tmdbIds failed!", "error", err)
		return err
	}
	return nil
}

func Add[S Addable](
	wp WatchedProvider,
	userId uint,
	s S,
	addCb func(w *entity.Watched),
) error {
	watchedEntry, err := wp.GetWatchedItemByTmdbId(
		userId,
		uint(s.GetId()),
		s.GetMediaType(),
	)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	addCb(&watchedEntry)
	return nil
}

func AddSingularAndList[S Addable, S2 Addable](
	wp WatchedProvider,
	userId uint,
	s S,
	addCb func(w *entity.Watched),
	list []*AddListCall[S2],
) error {
	err := Add(wp, userId, s, addCb)
	if err != nil {
		return err
	}
	for i := range list {
		if err := AddList(wp, userId, list[i].s, list[i].addCb); err != nil {
			return err
		}
	}
	return nil
}
