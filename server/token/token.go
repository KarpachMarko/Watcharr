package token

import (
	"errors"
	"log/slog"
	"time"

	"github.com/sbondCo/Watcharr/database/entity"
	"github.com/sbondCo/Watcharr/util"
	"gorm.io/gorm"
)

const TokenMaxAge = 2 * time.Minute

func CreateOneUseToken(db *gorm.DB, t entity.TokenType, userId uint) (string, error) {
	token, err := util.GenerateString(8)
	if err != nil {
		slog.Error("createOneUseToken: Failed to generate string!", "error", err)
		return "", errors.New("failed to generate token")
	}
	res := db.Create(&entity.Token{Type: t, Value: token, UserID: userId})
	if res.Error != nil {
		slog.Error("createOneUseToken: Failed to insert token into db!", "error", res.Error)
		return "", errors.New("failed to generate token")
	}
	return token, nil
}

// Cleans up tokens older than 2m.
func CleanupTokens(db *gorm.DB) {
	slog.Debug("cleanupTokens: Cleaning up old tokens from db")
	twoMinsAgo := time.Now().Add(-TokenMaxAge)
	resp := db.Where("created_at < ?", twoMinsAgo).Delete(&entity.Token{})
	if resp.Error != nil {
		slog.Error("cleanupTokens: Failed to run DELETE on old tokens!", "error", resp.Error)
	}
}
