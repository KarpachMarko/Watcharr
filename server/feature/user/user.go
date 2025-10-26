package user

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"log/slog"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/sbondCo/Watcharr/config"
	"github.com/sbondCo/Watcharr/database/entity"
	"github.com/sbondCo/Watcharr/feature/image"
	"gorm.io/gorm"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

// Update user settings
func (s *Service) UserUpdate(db *gorm.DB, userId uint, ur entity.UserSettings) (entity.UserSettings, error) {
	slog.Debug("user update request running", "user_id", userId, "ur", ur)
	user := new(entity.User)
	res := db.Where("id = ?", userId).Take(&user)
	if res.Error != nil {
		slog.Error("user update failed", "user_id", userId, "error", res.Error)
		return entity.UserSettings{}, errors.New("failed to retrieve user")
	}
	if ur.HideSpoilers != nil {
		user.HideSpoilers = ur.HideSpoilers
	}
	if ur.Private != nil {
		user.Private = ur.Private
	}
	if ur.PrivateThoughts != nil {
		user.PrivateThoughts = ur.PrivateThoughts
	}
	if ur.IncludePreviouslyWatched != nil {
		user.IncludePreviouslyWatched = ur.IncludePreviouslyWatched
	}
	if ur.AutomateShowStatuses != nil {
		user.AutomateShowStatuses = ur.AutomateShowStatuses
	}
	if ur.Country != nil {
		user.Country = ur.Country
	}
	if ur.RatingSystem != nil {
		user.RatingSystem = ur.RatingSystem
	}
	if ur.RatingStep != nil {
		user.RatingStep = ur.RatingStep
	}
	db.Save(&user)
	return entity.UserSettings{
		Private:                  user.Private,
		PrivateThoughts:          user.PrivateThoughts,
		HideSpoilers:             user.HideSpoilers,
		IncludePreviouslyWatched: user.IncludePreviouslyWatched,
		AutomateShowStatuses:     user.AutomateShowStatuses,
		Country:                  user.Country,
	}, nil
}

func (s *Service) UserGetSettings(db *gorm.DB, userId uint) (entity.UserSettings, error) {
	slog.Debug("user update request running", "user_id", userId)
	user := new(entity.User)
	res := db.Where("id = ?", userId).Take(&user)
	if res.Error != nil {
		slog.Error("user get failed", "user_id", userId, "error", res.Error)
		return entity.UserSettings{}, errors.New("failed to retrieve user")
	}
	return entity.UserSettings{
		Private:                  user.Private,
		PrivateThoughts:          user.PrivateThoughts,
		HideSpoilers:             user.HideSpoilers,
		IncludePreviouslyWatched: user.IncludePreviouslyWatched,
		AutomateShowStatuses:     user.AutomateShowStatuses,
		Country:                  user.Country,
		RatingSystem:             user.RatingSystem,
		RatingStep:               user.RatingStep,
	}, nil
}

func (s *Service) UserSearch(db *gorm.DB, currentUsersId uint, q string) ([]entity.PublicUser, error) {
	slog.Debug("user search request running", "query", q)
	users := new([]entity.PublicUser)
	res := db.Where("private = 0 AND username LIKE ? AND id != ?", "%"+q+"%", currentUsersId).Table("users").Find(&users)
	if res.Error != nil {
		slog.Error("user search failed", "error", res.Error)
		return []entity.PublicUser{}, errors.New("failed to find users")
	}
	return *users, nil
}

func (s *Service) GetUserInfo(db *gorm.DB, currentUsersId uint) (entity.PrivateUser, error) {
	slog.Debug("user get info request running")
	user := new(entity.PrivateUser)
	res := db.Where("id = ?", currentUsersId).Table("users").Preload("Avatar").Take(&user)
	if res.Error != nil {
		slog.Error("user get info failed", "error", res.Error)
		return entity.PrivateUser{}, errors.New("failed to find current user")
	}
	return *user, nil
}

// For getting a public user's info, when viewing their list for example
func (s *Service) GetUserPublicInfo(db *gorm.DB, userId uint, username string) (entity.PublicUser, error) {
	slog.Debug("user get info request running")
	user := new(entity.PublicUser)
	res := db.Where("private = 0 AND id = ? AND username = ?", userId, username).Table("users").Preload("Avatar").Take(&user)
	if res.Error != nil {
		slog.Error("public user get info failed", "error", res.Error)
		return entity.PublicUser{}, errors.New("failed to find user")
	}
	return *user, nil
}

func (s *Service) UserUpdateBio(db *gorm.DB, userId uint, newBio string) error {
	slog.Debug("userUpdateBio request running", "user_id", userId, "newBio", newBio)
	if res := db.Model(&entity.User{}).Where("id = ?", userId).Update("bio", newBio); res.Error != nil {
		slog.Error("userUpdateBio failed", "user_id", userId, "error", res.Error)
		return errors.New("failed to update bio")
	}
	return nil
}

func (s *Service) UploadUserAvatar(c *gin.Context, db *gorm.DB, userId uint) (entity.Image, error) {
	file, err := c.FormFile("avatar")
	if err != nil {
		slog.Error("failed to get file", "error", err)
		return entity.Image{}, errors.New("no file found")
	}

	slog.Debug("an avatar is being uploaded", "name", file.Filename)

	f, _ := file.Open()
	if err := image.IsValidImageType(f); err != nil {
		return entity.Image{}, errors.New("invalid image type")
	}
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err) // TODO nu le fatal
	}
	hs := hex.EncodeToString(h.Sum(nil))

	slog.Debug("image hash calculated", "hash", hs, "first_letter", hs[0:1])

	// Upload the file to specific dst.
	outp := path.Join("img/up/", hs[0:1], hs+filepath.Ext(file.Filename))
	c.SaveUploadedFile(file, path.Join(config.DataPath, outp))

	_, err = f.Seek(0, 0)
	if err != nil {
		slog.Error("uploadUserAvatar seeking back to start of image failed", "error", err)
	}

	// No need to remove old image, the daily cleanup task will handle removing unused ones.
	var img entity.Image
	err = db.Transaction(func(tx *gorm.DB) error {
		// Insert avatar into db
		img, err = image.InsertImage(db, hs, outp, f)
		if err != nil {
			return err
		}
		if img.ID == 0 {
			return errors.New("image has no id")
		}
		// Update users avatar to newly inserted
		if err := tx.Where("id = ?", userId).Updates(&entity.User{AvatarID: img.ID}).Error; err != nil {
			return err
		}
		// commit transaction if no errors
		return nil
	})
	if err != nil {
		slog.Error("uploadUserAvatar failed!", "error", err)
		return entity.Image{}, errors.New("uploadUserAvatar transaction failed")
	}
	return img, nil
}
