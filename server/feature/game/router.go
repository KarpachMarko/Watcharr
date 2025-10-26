package game

import (
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-gonic/gin"
	"github.com/sbondCo/Watcharr/database/entity"
	"github.com/sbondCo/Watcharr/feature/auth/authmiddleware"
	"github.com/sbondCo/Watcharr/media/igdb"
	"github.com/sbondCo/Watcharr/router"
)

type Router struct {
	br      *router.BaseRouter
	service *Service
}

func NewRouter(br *router.BaseRouter, service *Service) *Router {
	return &Router{
		br,
		service,
	}
}

func (r *Router) AddRoutes() {
	gamer := r.br.Router.Group("/game").Use(authmiddleware.AuthRequired(nil, r.br.Cfg))
	exp := time.Hour * 24

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

	// Search for games
	gamer.GET("/search", cache.CachePage(r.br.MemStore, exp, r.GetSearch))
	// Search for game by id (for search page, same minimal details as /search returned)
	gamer.GET("/search/:id", cache.CachePage(r.br.MemStore, exp, r.GetSearchById))
	// Game details for game page
	gamer.GET("/:id", cache.CachePage(r.br.MemStore, exp, r.GetGameDetails))
	// Add game to played(watched) list
	gamer.POST("/played", r.AddPlayed)

	// IMPORTANT: Routes below only for admins!
	gamer.Use(authmiddleware.AuthRequired(r.br.DB, r.br.Cfg), authmiddleware.AdminRequired())
	{
		gamer.POST("/config", r.UpdateConfig)
	}
}

func (r *Router) GetSearch(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: "a query was not provided"})
		return
	}
	decodedQuery, err := url.QueryUnescape(query)
	if err != nil {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: "query parameter invalid"})
		return
	}
	games, err := r.br.Cfg.TWITCH.Search(decodedQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, games)
}

func (r *Router) GetSearchById(c *gin.Context) {
	if c.Param("id") == "" {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: "an id was not provided"})
		return
	}
	games, err := r.br.Cfg.TWITCH.SearchById(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, games)
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
	withWatchedResp := gameDetailsAddWatched(r.br.DB, userId, content)
	c.JSON(http.StatusOK, withWatchedResp)
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
