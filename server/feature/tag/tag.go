package tag

import (
	"errors"
	"log/slog"

	"github.com/sbondCo/Watcharr/database/dbmodel"
	"github.com/sbondCo/Watcharr/database/entity"
	"github.com/sbondCo/Watcharr/domain"
	"gorm.io/gorm"
)

// I think tags will be private for the user.
// If the user wants to make a public list, they should make a custom view.

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetTags(db *gorm.DB, userId uint) ([]entity.Tag, error) {
	tags := new([]entity.Tag)
	res := db.Model(&entity.Tag{}).Where("user_id = ?", userId).Find(&tags)
	if res.Error != nil {
		slog.Error("getTags: Failed getting tags from database", "error", res.Error.Error())
		return []entity.Tag{}, errors.New("failed getting tags")
	}
	return *tags, nil
}

// func GetTag(db *gorm.DB, userId uint, tagId uint) (Tag, error) {
// 	tag := new(Tag)
// 	res := db.Model(&Tag{}).Where("id = ? AND user_id = ?", tagId, userId).Preload("Watched").Find(&tag)
// 	if res.Error != nil {
// 		slog.Error("getTag: Failed getting tag from database", "error", res.Error.Error())
// 		return Tag{}, errors.New("failed getting tag")
// 	}
// 	if tag.ID == 0 {
// 		slog.Error("getTag: Tag does not exist for this user.", "user_id", userId)
// 		return Tag{}, errors.New("tag does not exist")
// 	}
// 	return *tag, nil
// }

// This method should only be used when we don't have the tagId
// (eg: when we are importing data) because this is not technically
// reliable, since users can have multiple tags with the same name/colors
// (realistically they probably won't, but...).
func (s *Service) GetTagByNameAndColor(db *gorm.DB, userId uint, tagName string, tagColor string, tagBgColor string) (entity.Tag, error) {
	tag := new(entity.Tag)
	res := db.Model(&entity.Tag{}).Where("name = ? AND user_id = ? AND color = ? AND bg_color = ?", tagName, userId, tagColor, tagBgColor).Preload("Watched").Find(&tag)
	if res.Error != nil {
		slog.Error("getTagByNameAndColor: Failed getting tag from database", "error", res.Error.Error())
		return entity.Tag{}, errors.New("failed getting tag")
	}
	if tag.ID == 0 {
		slog.Error("getTagByNameAndColor: Tag does not exist for this user.", "user_id", userId)
		return entity.Tag{}, errors.New("tag does not exist")
	}
	return *tag, nil
}

// Let user create a tag.
func (s *Service) AddTag(db *gorm.DB, userId uint, tr domain.TagAddRequest) (entity.Tag, error) {
	if tr.Name == "" {
		return entity.Tag{}, errors.New("tag must have a name")
	}
	tag := entity.Tag{UserID: userId, Name: tr.Name, Color: tr.Color, BgColor: tr.BgColor}
	res := db.Create(&tag)
	if res.Error != nil {
		slog.Error("Error adding tag to database", "error", res.Error.Error())
		return entity.Tag{}, errors.New("failed adding new tag to database")
	}
	slog.Debug("Adding tag", "added_tag", tag)
	return tag, nil
}

// Let user update one of their tags (replaces).
func (s *Service) UpdateTag(db *gorm.DB, userId uint, tagId uint, tr domain.TagAddRequest) error {
	if tr.Name == "" {
		return errors.New("tag must have a name")
	}
	tag := entity.Tag{Name: tr.Name, Color: tr.Color, BgColor: tr.BgColor}
	res := db.Where("id = ? AND user_id = ?", tagId, userId).Updates(&tag)
	if res.Error != nil {
		slog.Error("Error updating tag in database", "error", res.Error.Error())
		return errors.New("failed updating tag in database")
	}
	if res.RowsAffected == 0 {
		slog.Error("updateTag: Zero rows affected.. tag likely does not exist", "tag_id", tagId, "user_id", userId)
		return errors.New("tag does not exist")
	}
	slog.Debug("updateTag:", "updated_tag", tag)
	return nil
}

// Let user delete their own tag.
func (s *Service) DeleteTag(db *gorm.DB, userId uint, tagId uint) error {
	if tagId == 0 {
		return errors.New("no tag id provided")
	}
	slog.Debug("deleteTag:", "tag_id", tagId, "user_id", userId)
	// Select("Watched") so relations in watched_tags table are removed too.
	// ID is passed in the .Delete param so the .Select call can do it's job (relies on the primary key).
	res := db.Unscoped().Where("id = ? AND user_id = ?", tagId, userId).Select("Watched").Delete(&entity.Tag{GormModel: dbmodel.GormModel{ID: tagId}})
	if res.Error != nil {
		slog.Error("deleteTag: Error deleting tag from database", "error", res.Error.Error(), "tag_id", tagId, "user_id", userId)
		return errors.New("failed deleting tag from database")
	}
	if res.RowsAffected == 0 {
		slog.Error("deleteTag: Zero rows affected.. tag must not exist for user", "tag_id", tagId, "user_id", userId)
		return errors.New("tag does not exist")
	}
	return nil
}
