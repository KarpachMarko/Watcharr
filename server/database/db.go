package database

import (
	"path"

	"github.com/sbondCo/Watcharr/config"
	"github.com/sbondCo/Watcharr/database/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Create a new database connection.
func New() (*gorm.DB, error) {
	// Open the database.
	db, err := gorm.Open(
		sqlite.Open(path.Join(config.DataPath, "watcharr.db")),
		&gorm.Config{TranslateError: true},
	)
	if err != nil {
		return nil, err
	}
	// Perform auto migration.
	err = db.AutoMigrate(
		&entity.User{},
		&entity.UserServices{},
		&entity.Content{},
		&entity.Watched{},
		&entity.WatchedSeason{},
		&entity.WatchedEpisode{},
		&entity.Activity{},
		&entity.Token{},
		&entity.Follow{},
		&entity.Image{},
		&entity.Game{},
		&entity.ArrRequest{},
		&entity.Tag{},
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}
