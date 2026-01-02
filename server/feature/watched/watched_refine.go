// Watched sorting & filtering.

package watched

import (
	"log/slog"
	"strings"

	"github.com/sbondCo/Watcharr/database/entity"
	"github.com/sbondCo/Watcharr/domain"
	"github.com/sbondCo/Watcharr/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func refineFilterType(db *gorm.DB, ft []util.SupportedMedia) {
	if len(ft) <= 0 {
		return
	}
	qstr := []string{}
	astr := []any{}
	for _, v := range ft {
		switch v {
		case util.SupportedMediaMovie:
			qstr = append(qstr, "Content.type = ?")
			astr = append(astr, "movie")
		case util.SupportedMediaShow:
			qstr = append(qstr, "Content.type = ?")
			astr = append(astr, "tv")
		case util.SupportedMediaGame:
			qstr = append(qstr, "watcheds.game_id IS NOT NULL")
		}
	}
	q := strings.Join(qstr, " OR ")
	slog.Debug("watchedRefine: Filter type.", "q", q, "astr", astr)
	db.Where(q, astr...)
}

func refineFilterStatus(db *gorm.DB, f []entity.WatchedStatus) {
	if len(f) <= 0 {
		return
	}
	for i := range f {
		// Ensure string **case** is valid WatchedStatus by converting to uppercase.
		f[i] = entity.WatchedStatus(strings.ToUpper(string(f[i])))
	}
	db.Where("watcheds.status IN ?", f)
}

func refineSort(db *gorm.DB, sort domain.WatchedSort, dir domain.SortDirection) {
	if sort == "" {
		return
	}
	obc := func(cn string) clause.OrderByColumn {
		o := clause.OrderByColumn{}
		if dir == domain.WatchedSortDirAsc {
			o.Desc = false
		} else {
			o.Desc = true
		}
		o.Column = clause.Column{Name: cn}
		return o
	}
	switch sort {
	case domain.WatchedSortDateAdded:
		db.Order(obc("watcheds.created_at"))
	case domain.WatchedSortLastChanged:
		db.Order(obc("watcheds.updated_at"))
	case domain.WatchedSortLastFinished:
		db.
			// This join looks for the latest activity for each watched entry
			// that indiciates a 'FINISHED' status. The date of these is used
			// in the sort below.
			// This seems the best way to support this sort with how our current
			// activity data is structured.
			Joins(`LEFT JOIN (
					SELECT
						watched_id AS a_watched_id,
						MAX(COALESCE(custom_date, created_at)) AS a_sort_by_date
					FROM activities
					WHERE data LIKE "%FINISHED%" AND deleted_at IS NULL
					GROUP BY watched_id
				) q ON q.a_watched_id = watcheds.id`).
			Order(obc("q.a_sort_by_date"))
	case domain.WatchedSortRating:
		db.Order(obc("watcheds.rating"))
	case domain.WatchedSortAlphabetical:
		db.
			Order(obc("Content__title")).
			Order(obc("Game__name"))
	}
}

// list data.
// gorm scope for applying sort and filters to watched
func watchedRefine(wr domain.WatchedGetPageRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Apply filters
		refineFilterType(db, wr.FilterType)
		refineFilterStatus(db, wr.FilterStatus)
		// Apply sort
		refineSort(db, wr.Sort, wr.SortDir)
		return db
	}
}
