package game

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/sbondCo/Watcharr/database/entity"
	"github.com/sbondCo/Watcharr/feature/auth/authmiddleware"
	"github.com/sbondCo/Watcharr/feature/watched/addedtocontent"
	"github.com/sbondCo/Watcharr/media/igdb"
	"github.com/sbondCo/Watcharr/router"
	"github.com/sbondCo/Watcharr/util"
)

type WatchedProvider interface {
	UpdateWatchedLastViewedSeason(userId uint, id uint, seasonNum int) error
	GetWatchedItemBySupportedMediaId(userId uint, id uint, t util.SupportedMedia) (entity.Watched, error)
	GetWatchedItemsBySupportedMediaIds(userId uint, c []addedtocontent.IdToTypePair) ([]entity.Watched, error)
}

type Router struct {
	br              *router.BaseRouter
	service         *Service
	watchedProvider WatchedProvider
}

func NewRouter(br *router.BaseRouter, service *Service, watchedProvider WatchedProvider) *Router {
	return &Router{
		br,
		service,
		watchedProvider,
	}
}

func (r *Router) AddRoutes() {
	gamer := r.br.Router.Group("/game").Use(authmiddleware.AuthRequired(nil, r.br.Cfg))

	// TODO This config init can be moved to NewRouter, then `gdb` can be accessible in Router for all service funcs.
	r.br.Cfg.TWITCH.OnTokenRefreshed(func() {
		// Save new token to config when we refresh it.
		slog.Debug("GameRoutes: token refreshed.. saving to config.")
		if err := r.br.Cfg.Write(); err != nil {
			slog.Error("GameRoutes: failed to save refreshed token to config.", "error", err)
		}
	})
	err := r.br.Cfg.TWITCH.Init()
	// Save cfg if init succeeded, this will save our access token
	if err != nil {
		slog.Error("GameRoutes: Twitch init failed!", "error", err)
	}

	// Game details for game page
	gamer.GET("/:id", r.GetGameDetails)
	// Add game to played(watched) list
	gamer.POST("/played", r.AddPlayed)

	// IMPORTANT: Routes below only for admins!
	gamer.Use(authmiddleware.AuthRequired(r.br.DB, r.br.Cfg), authmiddleware.AdminRequired())
	{
		gamer.POST("/config", r.UpdateConfig)
	}
}

func (r *Router) GetGameDetails(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	if c.Param("id") == "" {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: "an id was not provided"})
		return
	}
	content, err := r.br.Cfg.TWITCH.GameDetails(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
		return
	}
	ww := igdb.GameDetailsResponseWithWatched{}
	if err := copier.Copy(&ww, &content); err != nil {
		slog.Error("GetGameDetails: Failed to copy content to with watched struct", "error", err)
		c.JSON(
			http.StatusInternalServerError,
			router.ErrorResponse{Error: "failed to prepare response"},
		)
		return
	}
	if err := addedtocontent.AddSingularAndList(
		r.watchedProvider,
		userId,
		ww,
		func(w *entity.Watched) {
			ww.Watched = w
		},
		[]*addedtocontent.AddListCall[igdb.GameSimilarWithWatched]{
			addedtocontent.NewAddListCall(
				ww.SimilarGame,
				func(i int, w *entity.Watched) {
					ww.SimilarGame[i].Watched = w
				},
			),
		},
	); err != nil {
		slog.Error("GetGameDetails: Failed to add watched to content!", "error", err)
		c.JSON(
			http.StatusInternalServerError,
			router.ErrorResponse{Error: "failed to add watched data to response"},
		)
		return
	}
	c.JSON(http.StatusOK, ww)
}

func (r *Router) AddPlayed(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	var ar PlayedAddRequest
	err := c.ShouldBindJSON(&ar)
	if err == nil {
		response, err := r.service.addPlayed(r.br.DB, &r.br.Cfg.TWITCH, userId, ar, entity.ADDED_WATCHED)
		if err != nil {
			c.JSON(http.StatusForbidden, router.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, response)
		return
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
}

func (r *Router) UpdateConfig(c *gin.Context) {
	var ar igdb.IGDB
	err := c.ShouldBindJSON(&ar)
	if err == nil {
		err := r.br.Cfg.SaveTwitchConfig(ar)
		if err != nil {
			c.JSON(http.StatusForbidden, router.ErrorResponse{Error: err.Error()})
			return
		}
		// gdb = &b.cfg.TWITCH
		c.Status(http.StatusOK)
		return
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
}
