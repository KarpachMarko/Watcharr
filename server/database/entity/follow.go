package entity

import "time"

// Database struct, only internal.
type Follow struct {
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"-"`
	UserID         uint      `gorm:"primaryKey:usr_id_to_followed_id;not null;check:user_id != followed_user_id" json:"-"`
	User           User      `json:"-"`
	FollowedUserID uint      `gorm:"primaryKey:usr_id_to_followed_id;not null" json:"-"`
	FollowedUser   User      `json:"-"`
}
