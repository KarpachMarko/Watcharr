package activity

import (
	"errors"
	"log/slog"
	"time"

	"github.com/sbondCo/Watcharr/database/entity"
	"gorm.io/gorm"
)

type ActivityAddRequest struct {
	WatchedID  uint                `json:"watchedId" binding:"required"`
	Type       entity.ActivityType `json:"type" binding:"required"`
	Data       string              `json:"data" binding:"required"`
	CustomDate *time.Time          `json:"customDate,omitempty"`
}

type ActivityUpdateRequest struct {
	CustomDate time.Time `json:"customDate" binding:"required"`
}

func getActivity(db *gorm.DB, userId uint, watchedId uint) ([]entity.Activity, error) {
	activity := new([]entity.Activity)
	res := db.Model(&entity.Activity{}).Where("user_id = ? AND watched_id = ?", userId, watchedId).Find(&activity)
	if res.Error != nil {
		slog.Error("Failed getting activity from database", "error", res.Error.Error())
		return []entity.Activity{}, errors.New("failed getting activity")
	}
	return *activity, nil
}

func AddActivity(db *gorm.DB, userId uint, ar ActivityAddRequest) (entity.Activity, error) {
	if ar.WatchedID == 0 {
		return entity.Activity{}, errors.New("watchedId must be set to add an activity")
	}
	activity := entity.Activity{UserID: userId, WatchedID: ar.WatchedID, Type: ar.Type, Data: ar.Data, CustomDate: ar.CustomDate}
	res := db.Create(&activity)
	if res.Error != nil {
		slog.Error("Error adding activity to database", "error", res.Error.Error())
		return entity.Activity{}, errors.New("failed adding new activity to database")
	}
	slog.Debug("Adding activity", "added_activity", activity)
	return activity, nil
}

func updateActivity(db *gorm.DB, userId uint, id uint, activityUpdateRequest ActivityUpdateRequest) error {
	if id == 0 {
		return errors.New("id must be set to update an activity")
	}
	if activityUpdateRequest.CustomDate.IsZero() {
		return errors.New("customDate must be set to update an activity")
	}
	res := db.Model(&entity.Activity{}).Where("user_id = ? AND id = ?", userId, id).Update("custom_date", activityUpdateRequest.CustomDate)
	if res.Error != nil {
		slog.Error("Error updating activity in database", "error", res.Error.Error())
		return errors.New("failed updating activity in database")
	}
	if res.RowsAffected < 1 {
		slog.Error("No activities were updated. This may be because the activity doesn't exist or is not owned by the calling user.")
		return errors.New("failed updating activity in database")
	}
	slog.Debug("Updating activity", "updated_activity", id)
	return nil
}

func deleteActivity(db *gorm.DB, userId uint, id uint) error {
	if id == 0 {
		return errors.New("an id must be provided to delete an activity")
	}
	res := db.Where("user_id = ?", userId).Delete(&entity.Activity{}, id)
	if res.Error != nil {
		slog.Error("Error deleting activity in database", "error", res.Error.Error())
		return errors.New("failed deleting activity in database")
	}
	if res.RowsAffected < 1 {
		slog.Error("No activities were deleted. This may be because the activity doesn't exist or is not owned by the calling user.")
		return errors.New("failed deleting activity from database")
	}
	return nil
}
