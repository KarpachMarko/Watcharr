// Watched sorting & filtering.

package watched

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// gorm scope for applying sort and filters to watched
// list data.
func watchedRefine(wr WatchedGetPageRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Apply filters

		// Apply sort
		if wr.Sort != "" {
			obc := func(cn string) clause.OrderByColumn {
				o := clause.OrderByColumn{}
				if wr.SortDir == sortAscending {
					o.Desc = false
				} else {
					o.Desc = true
				}
				o.Column = clause.Column{Name: cn}
				return o
			}
			switch wr.Sort {
			case watchedSortDateAdded:
				db.Order(obc("watcheds.created_at"))
			case watchedSortLastChanged:
				db.Order(obc("watcheds.updated_at"))
			case watchedSortLastFinished:
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
			case watchedSortRating:
				db.Order(obc("watcheds.rating"))
			case watchedSortAlphabetical:
				db.
					Order(obc("Content__title")).
					Order(obc("Game__name"))
			}
		}
		return db
	}
}
