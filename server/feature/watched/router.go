package watched

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sbondCo/Watcharr/database/entity"
	"github.com/sbondCo/Watcharr/domain"
	"github.com/sbondCo/Watcharr/feature/auth/authmiddleware"
	"github.com/sbondCo/Watcharr/media/tmdb"
	"github.com/sbondCo/Watcharr/router"
	"github.com/sbondCo/Watcharr/util"
)

type Router struct {
	br *router.BaseRouter
	t  *tmdb.TMDB
	s  *Service
}

func NewRouter(
	br *router.BaseRouter,
	t *tmdb.TMDB,
	service *Service,
) *Router {
	return &Router{
		br: br,
		t:  t,
		s:  service,
	}
}

// TODO all handlers moving here, then the base router will become the initializer of all handlers and
// initial creator of all services and passes them down to services that want them

func (r *Router) AddRoutes() {
	watched := r.br.Router.Group("/watched").Use(authmiddleware.AuthRequired(nil, r.br.Cfg))

	watched.GET("", router.PaginatedRequest(false), r.GetWatchedList)
	watched.GET(":id/:username", r.GetPublicWatchedList)
	watched.POST("", r.AddWatched)
	watched.PUT(":id", r.UpdateWatched)
	watched.DELETE(":id", r.DeleteWatched)
	// TODO Move add/delete watched from tag to the `tag` package (the service code is there so the route may as well be under there, also avoids a circular dep).
	watched.POST(":id/tag/:tagId", r.AddWatchedToTag)
	watched.DELETE(":id/tag/:tagId", r.DeleteWatchedFromTag)
}

// Get our (logged in user) watched list.
func (r *Router) GetWatchedList(c *gin.Context) {
	isPaginated := c.MustGet("paginationEnabled").(bool)
	userId := c.MustGet("userId").(uint)
	if isPaginated {
		pp := c.MustGet("paginationParams").(util.PaginationParams)
		wp := domain.WatchedGetPageRequest{
			// Defaults..
			Sort:    domain.WatchedSortDateAdded,
			SortDir: domain.WatchedSortDirAsc,
		}
		if err := c.ShouldBind(&wp); err != nil {
			c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: "failed to get request parameters"})
			return
		}
		if wp, err := r.s.getWatchedPage(userId, pp, wp); err == nil {
			c.JSON(http.StatusOK, wp)
		} else {
			c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: "failed to get page"})
		}
		return
	}
	// Non paginated response (doesn't support sorting/filtering atm)
	if w, err := r.s.getWatched(userId); err == nil {
		c.JSON(http.StatusOK, w)
	} else {
		c.JSON(http.StatusForbidden, router.ErrorResponse{Error: "failed"})
	}
}

// Get another users watched list (if its public).
func (r *Router) GetPublicWatchedList(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		slog.Error("getPublicWatched route failed to convert id param to uint", "id", id)
		c.Status(400)
		return
	}
	response, err := r.s.getPublicWatched(uint(id), c.Param("username"))
	if err != nil {
		c.JSON(http.StatusForbidden, router.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (r *Router) AddWatched(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	var ar WatchedAddRequest
	err := c.ShouldBindJSON(&ar)
	if err == nil {
		response, err := r.s.AddWatched(userId, ar, entity.ADDED_WATCHED)
		if err != nil {
			c.JSON(http.StatusForbidden, router.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, response)
		return
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
}

func (r *Router) UpdateWatched(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(400)
		return
	}
	userId := c.MustGet("userId").(uint)
	var ur WatchedUpdateRequest
	err = c.ShouldBindJSON(&ur)
	if err == nil {
		response, err := r.s.updateWatched(userId, uint(id), ur)
		if err != nil {
			c.JSON(http.StatusForbidden, router.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, response)
		return
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
}

func (r *Router) DeleteWatched(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err == nil {
		userId := c.MustGet("userId").(uint)
		response, err := r.s.removeWatched(userId, uint(id))
		if err != nil {
			c.JSON(http.StatusForbidden, router.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, response)
		return
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
}

func (r *Router) AddWatchedToTag(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		slog.Error("tag watched route failed to convert id param to int", "error", err)
		c.Status(http.StatusBadRequest)
		return
	}
	tagId, err := strconv.Atoi(c.Param("tagId"))
	if err != nil {
		slog.Error("tag watched route failed to convert tagId param to int", "error", err)
		c.Status(http.StatusBadRequest)
		return
	}
	userId := c.MustGet("userId").(uint)
	err = AddWatchedToTag(r.br.DB, userId, uint(tagId), uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (r *Router) DeleteWatchedFromTag(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		slog.Error("tag watched route failed to convert id param to int", "error", err)
		c.Status(http.StatusBadRequest)
		return
	}
	tagId, err := strconv.Atoi(c.Param("tagId"))
	if err != nil {
		slog.Error("tag watched route failed to convert tagId param to int", "error", err)
		c.Status(http.StatusBadRequest)
		return
	}
	userId := c.MustGet("userId").(uint)
	err = RmWatchedFromTag(r.br.DB, userId, uint(tagId), uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
