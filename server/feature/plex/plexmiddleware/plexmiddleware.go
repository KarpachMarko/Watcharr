package plexmiddleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/sbondCo/Watcharr/config"
	"github.com/sbondCo/Watcharr/database/entity"
	"gorm.io/gorm"
)

// Plex access middleware, ensures user is a Plex user.
// To be ran after AuthRequired middleware with extra data.
func PlexAccessRequired(db *gorm.DB, cfg *config.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.MustGet("userId").(uint)
		slog.Debug("PlexAccessRequired middleware hit", "user_id", userId)
		userType := c.MustGet("userType").(entity.UserType)
		if cfg.PLEX_HOST == "" || cfg.PLEX_MACHINE_ID == "" {
			slog.Error("PlexAccessRequired: Plex has not been configured.", "user_id", userId)
			c.AbortWithStatus(401)
			return
		}
		if userType != entity.PLEX_USER {
			slog.Error("PlexAccessRequired: User is not a Plex user..", "user_id", userId, "user_type", userType)
			c.AbortWithStatus(401)
			return
		}
		userPlexService := new(entity.UserServices)
		if res := db.Where("user_id = ? AND name = ?", userId, "plex").Take(&userPlexService); res.Error != nil {
			slog.Error("PlexAccessRequired: Failed when attempting to get users plex service integration..", "user_id", userId, "user_type", userType)
			c.AbortWithStatus(401)
			return
		}
		if userPlexService.ClientID == "" || userPlexService.AuthToken == "" || userPlexService.AuthToken2 == "" {
			slog.Error("PlexAccessRequired: User has missing details from service (clientId, authToken or authToken2)..", "user_id", userId, "client_id", userPlexService.ClientID)
			c.AbortWithStatus(401)
			return
		}
		c.Set("plexAuthToken", userPlexService.AuthToken)
		c.Set("plexLocalAuthToken", userPlexService.AuthToken2)
	}
}
