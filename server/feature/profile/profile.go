package profile

import (
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/sbondCo/Watcharr/database/entity"
	"gorm.io/gorm"
)

type Profile struct {
	Joined               time.Time `json:"joined"`
	ShowsWatched         int32     `json:"showsWatched"`
	MoviesWatched        int32     `json:"moviesWatched"`
	MoviesWatchedRuntime uint32    `json:"moviesWatchedRuntime"`
	ShowsWatchedRuntime  uint32    `json:"showsWatchedRuntime"`
}

// Check if content has been previsouly watched by looking for related activity.
func hasBeenPreviouslyWatched(a *[]entity.Activity) bool {
	wp := false
	var relatedActivity []entity.Activity
	for _, v := range *a {
		if v.Type == entity.ADDED_WATCHED ||
			v.Type == entity.IMPORTED_ADDED_WATCHED ||
			v.Type == entity.IMPORTED_WATCHED ||
			v.Type == entity.STATUS_CHANGED {
			relatedActivity = append(relatedActivity, v)
		}
	}
	if len(relatedActivity) <= 0 {
		return false
	}
	for _, ra := range relatedActivity {
		if ra.Type == entity.IMPORTED_ADDED_WATCHED {
			wp = true
			break
		} else if ra.Type == entity.ADDED_WATCHED || ra.Type == entity.IMPORTED_WATCHED {
			if ra.Data == "" {
				continue
			}
			var v map[string]any
			err := json.Unmarshal([]byte(ra.Data), &v)
			if err != nil {
				slog.Error("Checking ADDED_WATCHED or IMPORTED_WATCHED.. failed to parse json data", "error", err)
				continue
			}
			if status, ok := v["status"]; ok {
				if status == "FINISHED" {
					wp = true
					break
				}
			}
		} else if ra.Type == entity.STATUS_CHANGED {
			if ra.Data == "FINISHED" {
				wp = true
				break
			}
		}
	}
	return wp
}

// Gets any data required for profile page
func getProfile(db *gorm.DB, userId uint) (Profile, error) {
	user := new(entity.User)
	res := db.Model(&entity.User{}).Where("id = ?", userId).Take(&user)
	if res.Error != nil {
		slog.Error("Failed to get profile:", "error", res.Error.Error())
		return Profile{}, errors.New("failed to get profile")
	}
	watched := new([]entity.Watched)
	res = db.Model(&entity.Watched{}).Preload("Content").Preload("Activity").Where("user_id = ?", userId).Find(&watched)
	if res.Error != nil {
		slog.Error("Profile: Failed to get watched for processing:", "error", res.Error.Error())
		return Profile{}, errors.New("failed to get watched for processing")
	}
	var (
		showsWatched         int32
		moviesWatched        int32
		moviesWatchedRuntime uint32
		showsWatchedRuntime  uint32
	)
	for _, w := range *watched {
		isFinished := false
		if w.Status == entity.FINISHED {
			isFinished = true
		} else if *user.IncludePreviouslyWatched && hasBeenPreviouslyWatched(&w.Activity) {
			// If status is not finished and user has IncludePreviouslyWatched enabled,
			// then we can also check if content hasBeenPreviouslyWatched.
			isFinished = true
		}
		if isFinished {
			if w.Content == nil {
				continue
			}
			c := *w.Content
			if c.Type == entity.SHOW {
				showsWatched++
				// This aint a science, just a very inaccurate guesstimate.
				if c.NumberOfEpisodes != 0 {
					var showRuntime uint32 = 30
					if c.Runtime != 0 {
						showRuntime = c.Runtime
					}
					showsWatchedRuntime += showRuntime * c.NumberOfEpisodes
					slog.Debug("calcualted", "show", c.Title, "runti", showRuntime*c.NumberOfEpisodes)
				}
			} else if c.Type == entity.MOVIE {
				moviesWatched++
				moviesWatchedRuntime += c.Runtime
			}
		}
	}
	profile := Profile{
		Joined:               user.CreatedAt,
		ShowsWatched:         showsWatched,
		MoviesWatched:        moviesWatched,
		MoviesWatchedRuntime: moviesWatchedRuntime,
		ShowsWatchedRuntime:  showsWatchedRuntime,
	}
	return profile, nil
}
