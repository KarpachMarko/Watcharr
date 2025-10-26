package domain

import (
	"time"

	"github.com/sbondCo/Watcharr/database/entity"
	"gorm.io/gorm"
)

type (
	ActivityAddRequest struct {
		WatchedID  uint                `json:"watchedId" binding:"required"`
		Type       entity.ActivityType `json:"type" binding:"required"`
		Data       string              `json:"data" binding:"required"`
		CustomDate *time.Time          `json:"customDate,omitempty"`
	}

	ActivityUpdateRequest struct {
		CustomDate time.Time `json:"customDate" binding:"required"`
	}

	ActivityAddProvider interface {
		AddActivity(db *gorm.DB, userId uint, ar ActivityAddRequest) (entity.Activity, error)
	}
)
