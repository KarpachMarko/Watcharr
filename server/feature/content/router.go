package content

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-gonic/gin"
	"github.com/sbondCo/Watcharr/database/entity"
	"github.com/sbondCo/Watcharr/feature/auth/authmiddleware"
	"github.com/sbondCo/Watcharr/feature/watched/addedtocontent"
	"github.com/sbondCo/Watcharr/router"
	"github.com/sbondCo/Watcharr/util"
)

type WatchedProvider interface {
	UpdateWatchedLastViewedSeason(userId uint, id uint, seasonNum int) error
	GetWatchedItemsByTmdbIds(userId uint, c [][]any) ([]entity.Watched, error)
}

type Router struct {
	br *router.BaseRouter
	cs *Service
	wp WatchedProvider
}

func NewRouter(br *router.BaseRouter, cs *Service, wp WatchedProvider) *Router {
	return &Router{
		br: br,
		cs: cs,
		wp: wp,
	}
}

func (r *Router) AddRoutes() {
	content := r.br.Router.Group("/content").Use(authmiddleware.AuthRequired(nil, r.br.Cfg))
	exp := time.Hour * 24

	// TODO verify the routes that use cache here actually need it
	// (because watched data will be added to most)

	// Search for content
	content.GET("/search/multi", router.PaginatedRequest(true), r.GetSearchMulti)
	// Search for movies
	content.GET("/search/movie", router.PaginatedRequest(true), r.GetSearchMovie)
	// Search for shows
	content.GET("/search/tv", router.PaginatedRequest(true), r.GetSearchTv)
	// Search for people
	content.GET("/search/person", router.PaginatedRequest(true), cache.CachePage(r.br.MemStore, exp, r.GetSearchPerson))
	// Search for content with external id
	content.GET("/search/ext/:id/:source", cache.CachePage(r.br.MemStore, exp, r.GetSearchByExternalId))
	// Get movie details (for movie page)
	content.GET("/movie/:id", router.WhereaboutsRequired(r.br.Cfg), cache.CachePage(r.br.MemStore, exp, r.GetMovieDetails))
	// Get movie cast
	content.GET("/movie/:id/credits", cache.CachePage(r.br.MemStore, exp, r.GetMovieCredits))
	// Get tv details (for tv page)
	content.GET("/tv/:id", router.WhereaboutsRequired(r.br.Cfg), r.GetTvDetails)
	// Get tv cast
	content.GET("/tv/:id/credits", cache.CachePage(r.br.MemStore, exp, r.GetTvCredits))
	// Get season details
	// Supports `watchedId` query parameter for saving the requested season as `LastViewedSeason`.
	content.GET("/tv/:id/season/:num", r.GetSeasonDetails)
	// Get person details
	content.GET("/person/:id", cache.CachePage(r.br.MemStore, exp, r.GetPerson))
	// Get person credits
	content.GET("/person/:id/credits", cache.CachePage(r.br.MemStore, exp, r.GetPersonCredits))
	// Discover movies
	content.GET("/discover/movies", r.GetDiscoverMovies)
	// Discover shows
	content.GET("/discover/tv", r.GetDiscoverTv)
	// Get all trending (movies, tv, people)
	content.GET("/trending", r.GetTrending)
	// Upcoming Movies
	content.GET("/upcoming/movies", r.GetUpcomingMovies)
	// Upcoming Tv
	content.GET("/upcoming/tv", r.GetUpcomingTv)
	// Available regions for watch providers
	content.GET("/regions", r.GetRegions)
}

func (r *Router) GetSearchMulti(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: "a query was not provided"})
		return
	}
	pp := c.MustGet("paginationParams").(util.PaginationParams)
	content, err := r.cs.SearchContent(query, pp.Page)
	if err != nil {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
		return
	}
	// withWatchedResp := searchContentAddWatched(r.br.DB, userId, content)
	// c.JSON(http.StatusOK, withWatchedResp)

	addedtocontent.AddWAC(content.Results, r.wp, userId)
	c.JSON(http.StatusOK, content)
}

func (r *Router) GetSearchMovie(c *gin.Context) {
	// userId := c.MustGet("userId").(uint)
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: "a query was not provided"})
		return
	}
	pp := c.MustGet("paginationParams").(util.PaginationParams)
	content, err := r.cs.SearchMovies(query, pp.Page)
	if err != nil {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
		return
	}
	// withWatchedResp := searchMoviesAddWatched(r.br.DB, userId, content)
	// c.JSON(http.StatusOK, withWatchedResp)
	// HACK TEST
	c.JSON(http.StatusOK, content)
}

func (r *Router) GetSearchTv(c *gin.Context) {
	// userId := c.MustGet("userId").(uint)
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: "a query was not provided"})
		return
	}
	pp := c.MustGet("paginationParams").(util.PaginationParams)
	content, err := r.cs.SearchTv(query, pp.Page)
	if err != nil {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
		return
	}
	// withWatchedResp := searchTvAddWatched(r.br.DB, userId, content)
	// c.JSON(http.StatusOK, withWatchedResp)
	// HACK TEST
	c.JSON(http.StatusOK, content)
}

func (r *Router) GetSearchPerson(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: "a query was not provided"})
		return
	}
	pp := c.MustGet("paginationParams").(util.PaginationParams)
	content, err := r.cs.SearchPeople(query, pp.Page)
	if err != nil {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, content)
}

func (r *Router) GetSearchByExternalId(c *gin.Context) {
	if c.Param("id") == "" {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: "an id was not provided"})
		return
	}
	content, err := r.cs.SearchByExternalId(c.Param("id"), c.Param("source"))
	if err != nil {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, content)
}

func (r *Router) GetMovieDetails(c *gin.Context) {
	// userId := c.MustGet("userId").(uint)
	if c.Param("id") == "" {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: "an id was not provided"})
		return
	}
	content, err := r.cs.MovieDetails(
		c.Param("id"),
		c.MustGet("userCountry").(string),
		map[string]string{
			"append_to_response": "videos,watch/providers,similar",
		},
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
		return
	}
	// withWatchedResp := movieDetailsAddWatched(r.br.DB, userId, content)
	// c.JSON(http.StatusOK, withWatchedResp)
	// HACK TEST
	c.JSON(http.StatusOK, content)
}

func (r *Router) GetMovieCredits(c *gin.Context) {
	if c.Param("id") == "" {
		c.Status(400)
		return
	}
	content, err := r.cs.MovieCredits(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, content)
}

func (r *Router) GetTvDetails(c *gin.Context) {
	// userId := c.MustGet("userId").(uint)
	if c.Param("id") == "" {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: "an id was not provided"})
		return
	}
	// 1. Get details
	content, err := r.cs.TvDetails(
		c.Param("id"),
		c.MustGet("userCountry").(string),
		map[string]string{
			"append_to_response": "videos,watch/providers,similar,external_ids,keywords",
		},
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
		return
	}
	// withWatchedResp := tvDetailsAddWatched(r.br.DB, userId, content)
	// c.JSON(http.StatusOK, withWatchedResp)
	// HACK TEST
	c.JSON(http.StatusOK, content)
}

func (r *Router) GetTvCredits(c *gin.Context) {
	if c.Param("id") == "" {
		c.Status(400)
		return
	}
	content, err := r.cs.TvCredits(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, content)
}

// Get season details
// Supports `watchedId` query parameter for saving the requested season as `LastViewedSeason`.
func (r *Router) GetSeasonDetails(c *gin.Context) {
	if c.Param("id") == "" || c.Param("num") == "" {
		c.Status(400)
		return
	}
	content, err := r.cs.SeasonDetails(c.Param("id"), c.Param("num"))
	if err != nil {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
		return
	}
	// If a `watchedId` is passed, we should update it with this season
	// number, so the LastViewedSeason field is up to date (this seemed
	// better than making a new request for just saving this).
	// We will attach a `watcharr-lastviewedseason-saved` header if
	// this part succeeds so the client can decide on showing an error.
	if watchedIdQ := c.Query("watchedId"); watchedIdQ != "" {
		userId := c.MustGet("userId").(uint)
		watchedId, err := strconv.ParseUint(watchedIdQ, 10, 64)
		if err != nil {
			slog.Error("get season details route: Processing watchedId param failed", "error", err.Error(), "id", watchedIdQ)
		} else {
			if seasonNum, err := strconv.ParseInt(c.Param("num"), 10, 64); err == nil {
				if err = r.wp.UpdateWatchedLastViewedSeason(userId, uint(watchedId), int(seasonNum)); err == nil {
					c.Header("watcharr-lastviewedseason-saved", "1")
				}
			} else {
				slog.Error("get season details route: Parsing season number as int failed", "error", err.Error(), "season_num", c.Param("num"))
			}
		}
	} else {
		slog.Debug("get season details route: No watchedId parameter found.. not doing anything.")
	}
	c.JSON(http.StatusOK, content)
}

func (r *Router) GetPerson(c *gin.Context) {
	if c.Param("id") == "" {
		c.Status(400)
		return
	}
	content, err := r.cs.PersonDetails(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, content)
}

func (r *Router) GetPersonCredits(c *gin.Context) {
	if c.Param("id") == "" {
		c.Status(400)
		return
	}
	content, err := r.cs.PersonCredits(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, content)
}

func (r *Router) GetDiscoverMovies(c *gin.Context) {
	// userId := c.MustGet("userId").(uint)
	content, err := r.cs.DiscoverMovies()
	if err != nil {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
		return
	}
	// withWatchedResp := discoverMoviesAddWatched(r.br.DB, userId, content)
	// c.JSON(http.StatusOK, withWatchedResp)
	// HACK TEST
	c.JSON(http.StatusOK, content)
}

func (r *Router) GetDiscoverTv(c *gin.Context) {
	// userId := c.MustGet("userId").(uint)
	content, err := r.cs.DiscoverTv()
	if err != nil {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
		return
	}
	// withWatchedResp := discoverTvAddWatched(r.br.DB, userId, content)
	// c.JSON(http.StatusOK, withWatchedResp)
	// HACK TEST
	c.JSON(http.StatusOK, content)
}

func (r *Router) GetTrending(c *gin.Context) {
	// userId := c.MustGet("userId").(uint)
	content, err := r.cs.AllTrending()
	if err != nil {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
		return
	}
	// withWatchedResp := allTrendingAddWatched(r.br.DB, userId, content)
	// c.JSON(http.StatusOK, withWatchedResp)
	// HACK TEST
	c.JSON(http.StatusOK, content)
}

func (r *Router) GetUpcomingMovies(c *gin.Context) {
	// userId := c.MustGet("userId").(uint)
	content, err := r.cs.UpcomingMovies()
	if err != nil {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
		return
	}
	// withWatchedResp := upcomingMoviesAddWatched(r.br.DB, userId, content)
	// c.JSON(http.StatusOK, withWatchedResp)
	// HACK TEST
	c.JSON(http.StatusOK, content)
}

func (r *Router) GetUpcomingTv(c *gin.Context) {
	// userId := c.MustGet("userId").(uint)
	content, err := r.cs.UpcomingTv()
	if err != nil {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
		return
	}
	// withWatchedResp := upcomingTvAddWatched(r.br.DB, userId, content)
	// c.JSON(http.StatusOK, withWatchedResp)
	// HACK TEST
	c.JSON(http.StatusOK, content)
}

func (r *Router) GetRegions(c *gin.Context) {
	re, err := r.cs.Regions()
	if err != nil {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, re)
}
