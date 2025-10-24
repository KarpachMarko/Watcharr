package jellyfin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sbondCo/Watcharr/database/entity"
	"github.com/sbondCo/Watcharr/feature/auth/authmiddleware"
	"github.com/sbondCo/Watcharr/router"
)

type Router struct {
	br          *router.BaseRouter
	s           *Service
	syncService *SyncService
}

func NewRouter(br *router.BaseRouter, s *Service, syncService *SyncService) *Router {
	return &Router{
		br:          br,
		s:           s,
		syncService: syncService,
	}
}

func (r *Router) AddRoutes() {
	jf := r.br.Router.Group("/jellyfin").Use(authmiddleware.AuthRequired(r.br.DB, r.br.Cfg), r.s.JellyfinAccessRequired(r.br.Cfg))

	// Check if jf has item
	jf.GET("/:type/:name/:tmdbId", r.GetFindContent)
	// Sync users jellyfin watched items to watchlist
	jf.GET("/sync", r.GetSync)
}

// Check if jf has item
func (r *Router) GetFindContent(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	userType := c.MustGet("userType").(entity.UserType)
	username := c.MustGet("username").(string)
	userThirdPartyId := c.MustGet("userThirdPartyId").(string)
	userThirdPartyAuth := c.MustGet("userThirdPartyAuth").(string)
	response, err := r.s.JellyfinContentFind(
		userId,
		userType,
		username,
		userThirdPartyId,
		userThirdPartyAuth,
		c.Param("type"),
		c.Param("name"),
		c.Param("tmdbId"),
	)
	if err != nil {
		c.JSON(http.StatusForbidden, router.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// Sync users jellyfin watched items to watchlist
func (r *Router) GetSync(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	userType := c.MustGet("userType").(entity.UserType)
	username := c.MustGet("username").(string)
	userThirdPartyId := c.MustGet("userThirdPartyId").(string)
	userThirdPartyAuth := c.MustGet("userThirdPartyAuth").(string)
	response, err := r.syncService.jellyfinSyncWatched(r.br.DB, userId, userType, username, userThirdPartyId, userThirdPartyAuth)
	if err != nil {
		c.JSON(http.StatusForbidden, router.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}
