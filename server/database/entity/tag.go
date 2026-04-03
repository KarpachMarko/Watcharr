package entity

import "github.com/sbondCo/Watcharr/database/dbmodel"

type Tag struct {
	dbmodel.GormModel
	// ID of user that own this tag.
	UserID uint `json:"-" gorm:"not null"`
	// Name of the tag.
	Name string `json:"name" gorm:"not null"`
	// Hex of text color.
	Color string `json:"color"`
	// Hex of background color.
	BgColor string `json:"bgColor"`
	// All watched items.
	Watched []Watched `json:"watched,omitempty" gorm:"many2many:watched_tags;"`
}
