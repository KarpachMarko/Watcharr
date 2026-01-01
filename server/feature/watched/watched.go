package watched

import (
	"encoding/json"
	"errors"
	"log/slog"
	"strconv"
	"time"

	"github.com/sbondCo/Watcharr/database/entity"
	"github.com/sbondCo/Watcharr/domain"
	"github.com/sbondCo/Watcharr/util"
	"gorm.io/gorm"
)

type WatchedAddRequest struct {
	Status      entity.WatchedStatus `json:"status"`
	Rating      float64              `json:"rating" binding:"max=10"`
	Thoughts    string               `json:"thoughts"`
	ContentID   int                  `json:"contentId" binding:"required"`
	ContentType entity.ContentType   `json:"contentType" binding:"required,oneof=movie tv"`
	// Pass a watched date and we will set the CreatedAt (and initial UpdatedAt)
	// properties for this watched entry to this specific date.
	WatchedDate time.Time `json:"watchedDate,omitempty"`
}

type WatchedUpdateRequest struct {
	Status         entity.WatchedStatus `json:"status" binding:"required_without_all=Rating Thoughts RemoveThoughts Pinned"`
	Rating         float64              `json:"rating" binding:"max=10,required_without_all=Status Thoughts RemoveThoughts Pinned"`
	Thoughts       string               `json:"thoughts" binding:"required_without_all=Status Rating RemoveThoughts Pinned"`
	RemoveThoughts bool                 `json:"removeThoughts"`
	Pinned         *bool                `json:"pinned" binding:"required_without_all=Status Rating Thoughts RemoveThoughts"`
}

type WatchedUpdateResponse struct {
	NewActivity entity.Activity `json:"newActivity"`
}

type WatchedRemoveResponse struct {
	NewActivity entity.Activity `json:"newActivity"`
}

type ContentProvider interface {
	GetOrCacheContent(contentType entity.ContentType, tmdbId int) (entity.Content, error)
}

type Service struct {
	db               *gorm.DB
	cp               ContentProvider
	activityProvider domain.ActivityAddProvider
}

func NewService(db *gorm.DB, cp ContentProvider, activityProvider domain.ActivityAddProvider) *Service {
	return &Service{
		db,
		cp,
		activityProvider,
	}
}

// Get entire watched list
func (s *Service) getWatched(userId uint) ([]entity.Watched, error) {
	watched := new([]entity.Watched)
	res := s.db.Model(&entity.Watched{}).
		Preload("Content").
		Preload("Game").
		Preload("Game.Poster").
		Preload("Activity").
		Preload("WatchedSeasons").
		Preload("WatchedEpisodes").
		Preload("Tags").
		Where("user_id = ?", userId).
		Find(&watched)
	if res.Error != nil {
		slog.Error("getWatched: Failed!", "error", res.Error)
		return []entity.Watched{}, res.Error
	}
	return *watched, nil
}

// Returns a page of users watched list.
func (s *Service) getWatchedPage(
	userId uint,
	pp util.PaginationParams,
	wr domain.WatchedGetPageRequest,
) (util.PaginationResponse[entity.Watched], error) {
	slog.Debug("getWatchedPage: A page was requested.", "user_id", userId, "pagination_params", pp, "wr", wr)
	watched := new([]entity.Watched)
	pRes := &util.PaginationResponse[entity.Watched]{}
	res := s.db.Debug().
		Model(&entity.Watched{}).
		Where(&entity.Watched{UserID: userId}).
		Joins("Content").
		Joins("Game").
		Preload("Game.Poster").
		Preload("Tags").
		// Refine our results first (filters, sort);
		Scopes(
			watchedRefine(wr),
		).
		// Then count results (after filter);
		Count(&pRes.TotalResults).
		// Now calculate pagination properties with a TotalResults
		// that takes filtered out items into account.
		Scopes(
			util.Paginate(pp, pRes),
		).
		Find(&watched)
	if res.Error != nil {
		slog.Error("getWatchedPage: Failed!", "error", res.Error)
		return util.PaginationResponse[entity.Watched]{}, res.Error
	}
	pRes.Results = *watched
	pRes.Finished(pp)
	return *pRes, nil
}

// Get a watched list item by id (must be for `userId`).
func (s *Service) GetWatchedItemById(userId uint, id uint) (entity.Watched, error) {
	watched := new(entity.Watched)
	res := s.db.Model(&entity.Watched{}).Preload("Content").Where("user_id = ? AND id = ?", userId, id).Find(&watched)
	if res.Error != nil {
		slog.Error("GetWatchedItemById: Failed!", "error", res.Error)
		return entity.Watched{}, res.Error
	}
	return *watched, nil
}

// Get a watched list item by content (tmdb) id (must be for `userId`).
func (s *Service) GetWatchedItemByTmdbId(userId uint, tmdbId uint, contentType entity.ContentType) (entity.Watched, error) {
	slog.Debug("GetWatchedItemByTmdbId: Running.", "userId", userId, "tmdbId", tmdbId)
	watched := new(entity.Watched)
	res := s.db.Model(&entity.Watched{}).
		Joins("Content").
		Preload("Activity").
		Preload("WatchedSeasons").
		Preload("WatchedEpisodes").
		Preload("Tags").
		Where("user_id = ? AND Content.tmdb_id = ? AND Content.type = ?", userId, tmdbId, contentType).
		Take(&watched)
	if res.Error != nil {
		slog.Error("GetWatchedItemByTmdbId: Failed!", "error", res.Error)
		return entity.Watched{}, res.Error
	}
	slog.Debug("GetWatchedItemByTmdbId: Done.", "userId", userId, "tmdbId", tmdbId, "watched_item", watched)
	return *watched, nil
}

// Same as `getWatchedItemByTmdbId` except for getting in bulk (multiple content ids).
// `c` entries should be in format: [tmdb_id, ContentType] (Note: Couldn't figure out
// if it's possible to type this to enforce [int, ContentType] type for entries)
func (s *Service) GetWatchedItemsByTmdbIds(userId uint, c [][]any) ([]entity.Watched, error) {
	slog.Debug("GetWatchedItemsByTmdbIds: Running.", "userId", userId, "c", c)
	watched := new([]entity.Watched)
	res := s.db.Model(&entity.Watched{}).
		Joins("Content").
		Preload("Activity").
		Preload("WatchedSeasons").
		Preload("WatchedEpisodes").
		Preload("Tags").
		Where("user_id = ?", userId).
		Where("(Content.tmdb_id, Content.type) IN ?", c).
		Find(&watched)
	if res.Error != nil {
		slog.Error("GetWatchedItemsByTmdbIds: Failed!", "error", res.Error)
		return []entity.Watched{}, res.Error
	}
	slog.Debug(
		"GetWatchedItemsByTmdbIds: Done.",
		"userId", userId,
		"watcheds_found", len(*watched),
		// "wdev", *watched,
	)
	return *watched, nil
}

// Get a watched list item by game (igdb) id (must be for `userId`).
// TODO update var names soon
func (s *Service) GetWatchedItemByIgdbId(userId uint, tmdbId uint) (entity.Watched, error) {
	slog.Debug("getWatchedItemByIgdbId: Running.", "userId", userId, "tmdbId", tmdbId)
	watched := new(entity.Watched)
	res := s.db.Model(&entity.Watched{}).
		Joins("Game").
		Preload("Game.Poster").
		Preload("Activity").
		Preload("Tags").
		Where("user_id = ? AND Game.igdb_id = ?", userId, tmdbId).
		Take(&watched)
	if res.Error != nil {
		slog.Error("getWatchedItemByIgdbId: Failed!", "error", res.Error)
		return entity.Watched{}, res.Error
	}
	slog.Debug("getWatchedItemByIgdbId: Done.", "userId", userId, "tmdbId", tmdbId, "watched_item", watched)
	return *watched, nil
}

// Same as `getWatchedItemByIgdbId` except for getting in bulk (multiple content ids).
// `c` should be a slice of igdb ids.
func (s *Service) GetWatchedItemsByIgdbIds(userId uint, c []int) ([]entity.Watched, error) {
	slog.Debug("getWatchedItemsByIgdbIds: Running.", "userId", userId, "c", c)
	watched := new([]entity.Watched)
	res := s.db.Model(&entity.Watched{}).
		Joins("Game").
		Preload("Game.Poster").
		Preload("Activity").
		Preload("Tags").
		Where("user_id = ?", userId).
		Where("(Game.igdb_id) IN ?", c).
		Find(&watched)
	if res.Error != nil {
		slog.Error("getWatchedItemsByIgdbIds: Failed!", "error", res.Error)
		return []entity.Watched{}, res.Error
	}
	slog.Debug(
		"getWatchedItemsByIgdbIds: Done.",
		"userId", userId,
		"watcheds_found", len(*watched),
		// "wdev", *watched,
	)
	return *watched, nil
}

// Get another users **public** watchlist.
func (s *Service) getPublicWatched(userId uint, username string) ([]entity.Watched, error) {
	slog.Debug("getPublicWatched running", "user_id", userId, "username", username)
	// First we need to make sure the users list is public
	user := new(entity.User)
	// Figure we require knowlege of the users id and name to make it
	// harder to just type in random ids and see someones list.. dunno
	// if this is a thing we need but its here.. for now at least.
	res := s.db.Where("id = ? AND username = ?", userId, username).Take(&user)
	if res.Error != nil {
		slog.Error("Failed to get user for getPublicWatched request", "user_id", userId)
		return []entity.Watched{}, errors.New("failed to check privacy settings")
	}
	if user.Private != nil && *user.Private {
		slog.Error("getPublicWatched attempted to get a private list", "user_id", userId)
		return []entity.Watched{}, errors.New("this watched list is private")
	}
	// Now we know the user is public, return their list
	watched := new([]entity.Watched)
	res = s.db.Model(&entity.Watched{}).Preload("Content").Preload("Game").Preload("Game.Poster").Preload("Activity").Where("user_id = ?", userId).Find(&watched)
	if res.Error != nil {
		panic(res.Error)
	}
	return *watched, nil
}

func (s *Service) AddWatched(userId uint, ar WatchedAddRequest, at entity.ActivityType) (entity.Watched, error) {
	slog.Debug("Adding watched item", "userId", userId, "contentType", ar.ContentType, "contentId", ar.ContentID)
	// Get content cache (or cache it if we don't have it locally)
	content, err := s.cp.GetOrCacheContent(ar.ContentType, ar.ContentID)
	if err != nil {
		return entity.Watched{}, err
	}
	// Error if content has no id
	if content.ID == 0 {
		return entity.Watched{}, errors.New("failed to find content id")
	}
	// Create watched entry in db
	if ar.Status == "" {
		// Set default status for when content is added by
		// rating it instead of giving status first.
		if ar.ContentType == "movie" {
			ar.Status = entity.FINISHED
		} else {
			ar.Status = entity.WATCHING
		}
	}
	watched := entity.Watched{Status: ar.Status, Rating: ar.Rating, UserID: userId, ContentID: &content.ID}
	if ar.Thoughts != "" {
		watched.Thoughts = ar.Thoughts
	}
	// If custom WatchedDate passed, set CreatedAt and UpdatedAt fields to it.
	if !ar.WatchedDate.IsZero() {
		slog.Debug("Adding watched item: The provided WatchedDate is valid.", "watched_date", ar.WatchedDate, "userId", userId, "contentType", ar.ContentType, "contentId", ar.ContentID)
		watched.CreatedAt = ar.WatchedDate
		watched.UpdatedAt = ar.WatchedDate
	}
	res := s.db.Create(&watched)
	if res.Error != nil {
		if res.Error == gorm.ErrDuplicatedKey {
			res = s.db.Model(&entity.Watched{}).Unscoped().Preload("Activity").Where("user_id = ? AND content_id = ?", userId, watched.ContentID).Take(&watched)
			if res.Error != nil {
				return entity.Watched{}, errors.New("content already on watched list. errored checking for soft deleted record")
			}
			if watched.DeletedAt.Time.IsZero() {
				return watched, errors.New("content already on watched list")
			} else {
				slog.Info("addWatched: Watched list item for this content exists as soft deleted record.. attempting to restore")
				res = s.db.Model(&entity.Watched{}).Unscoped().Where("user_id = ? AND content_id = ?", userId, watched.ContentID).Updates(map[string]interface{}{"status": ar.Status, "rating": ar.Rating, "deleted_at": nil})
				watched.Status = ar.Status
				watched.Rating = ar.Rating
				watched.Thoughts = ar.Thoughts
				if res.Error != nil {
					slog.Error("addWatched: Failed to restore soft deleted watch list item", "error", res.Error)
					return entity.Watched{}, errors.New("content already on watched list. errored removing soft delete timestamp")
				}
			}
		} else {
			slog.Error("Error adding watched content to database", "error", res.Error.Error())
			return entity.Watched{}, errors.New("failed adding content to database")
		}
	}
	slog.Debug("Added watched list item", "item", watched)

	var act entity.Activity
	activityJson, err := json.Marshal(map[string]interface{}{"status": ar.Status, "rating": ar.Rating})
	if err != nil {
		slog.Error("Failed to marshal json for data in ADD_WATCHED activity request, adding without data", "error", err.Error())
		act, _ = s.activityProvider.AddActivity(userId, domain.ActivityAddRequest{WatchedID: watched.ID, Type: at})
	} else {
		act, _ = s.activityProvider.AddActivity(userId, domain.ActivityAddRequest{WatchedID: watched.ID, Type: at, Data: string(activityJson)})
	}
	watched.Activity = append(watched.Activity, act)
	watched.Content = &content
	return watched, nil
}

// this method is too ugly to look at please make him look better, future irhm
func (s *Service) updateWatched(userId uint, id uint, ar WatchedUpdateRequest) (WatchedUpdateResponse, error) {
	slog.Debug("UpdateWatched", "request_data", ar)
	upwat := entity.Watched{}
	res := s.db.Model(&entity.Watched{}).Where("id = ? AND user_id = ?", id, userId).Take(&upwat)
	if res.Error != nil {
		slog.Error("Watched entry update failed:", "id", id, "error", res.Error.Error())
		return WatchedUpdateResponse{}, errors.New("failed to update watched entry")
	}
	originalThoughts := upwat.Thoughts
	if ar.Rating != 0 {
		upwat.Rating = ar.Rating
	}
	if ar.Status != "" {
		upwat.Status = ar.Status
	}
	if ar.Thoughts != "" {
		upwat.Thoughts = ar.Thoughts
	}
	if ar.RemoveThoughts {
		upwat.Thoughts = ""
	}
	if ar.Pinned != nil {
		upwat.Pinned = *ar.Pinned
	}
	res = s.db.Save(upwat)
	if res.RowsAffected <= 0 {
		return WatchedUpdateResponse{}, errors.New("no watched entry found")
	}
	addedActivity := entity.Activity{}
	if ar.Rating != 0 {
		addedActivity, _ = s.activityProvider.AddActivity(userId, domain.ActivityAddRequest{WatchedID: id, Type: entity.RATING_CHANGED, Data: strconv.Itoa(int(ar.Rating))})
	}
	if ar.Status != "" {
		addedActivity, _ = s.activityProvider.AddActivity(userId, domain.ActivityAddRequest{WatchedID: id, Type: entity.STATUS_CHANGED, Data: string(ar.Status)})
	}
	if ar.Thoughts != "" {
		addedActivity, _ = s.activityProvider.AddActivity(userId, domain.ActivityAddRequest{WatchedID: id, Type: entity.THOUGHTS_CHANGED})
	}
	if ar.RemoveThoughts {
		addedActivity, _ = s.activityProvider.AddActivity(userId, domain.ActivityAddRequest{WatchedID: id, Type: entity.THOUGHTS_REMOVED, Data: originalThoughts})
	}
	return WatchedUpdateResponse{NewActivity: addedActivity}, nil
}

func (s *Service) UpdateWatchedLastViewedSeason(userId uint, id uint, seasonNum int) error {
	slog.Debug("UpdateWatchedLastViewedSeason", "user_id", userId, "id", id, "season_num", seasonNum)
	res := s.db.
		Model(&entity.Watched{}).
		Where("id = ? AND user_id = ?", id, userId).
		Update("last_viewed_season", seasonNum)
	if res.Error != nil {
		slog.Error("updateWatchedLastViewedSeason: Failed when updating.", "error", res.Error)
		return errors.New("failed to update db")
	}
	if res.RowsAffected == 0 {
		// likely the watched entry does not exist or is not owned by this `userId`.
		slog.Error("updateWatchedLastViewedSeason: Watched entry does not exist.")
		return errors.New("watched entry does not exist")
	}
	return nil
}

func (s *Service) removeWatched(userId uint, id uint) (WatchedRemoveResponse, error) {
	slog.Debug("Removing watched item:", "id", id, "user_id", userId)
	// Our model has a deleted_at field, which will make gorm do a soft delete.
	// Since other tables (eg activities) will link their rows to a watched_id, it's best to soft
	// delete, so if user restores watched item they still have activity for example (also so
	// someone else wont get other users activity if auto increment gives them the same watched id).
	res := s.db.Model(&entity.Watched{}).Where("id = ? AND user_id = ?", id, userId).Delete(&entity.Watched{})
	if res.Error != nil {
		slog.Error("Removing watched entry failed", "id", id, "error", res.Error.Error())
		return WatchedRemoveResponse{}, errors.New("failed to remove watched entry")
	}
	if res.RowsAffected <= 0 {
		return WatchedRemoveResponse{}, errors.New("no watched entry found")
	}
	addedActivity, _ := s.activityProvider.AddActivity(userId, domain.ActivityAddRequest{WatchedID: id, Type: entity.REMOVED_WATCHED})
	return WatchedRemoveResponse{NewActivity: addedActivity}, nil
}
