package activity

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sbondCo/Watcharr/feature/auth/authmiddleware"
	"github.com/sbondCo/Watcharr/router"
)

type Router struct {
	br *router.BaseRouter
}

func NewRouter(br *router.BaseRouter) *Router {
	return &Router{br: br}
}

func (r *Router) AddRoutes() {
	activity := r.br.Router.Group("/activity").Use(authmiddleware.AuthRequired(nil, r.br.Cfg))

	activity.GET(":watchedId", r.GetActivity)
	activity.POST("", r.AddActivity)
	activity.PUT(":id", r.UpdateActivity)
	activity.DELETE(":id", r.DeleteActivity)
}

func (r *Router) GetActivity(c *gin.Context) {
	watchedId, err := strconv.ParseUint(c.Param("watchedId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, router.ErrorResponse{Error: "check watched id route param"})
		return
	}
	userId := c.MustGet("userId").(uint)
	activity, err := getActivity(r.br.DB, userId, uint(watchedId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, router.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, activity)
}

func (r *Router) AddActivity(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	var ar ActivityAddRequest
	err := c.ShouldBindJSON(&ar)
	if err == nil {
		response, err := AddActivity(r.br.DB, userId, ar)
		if err != nil {
			c.JSON(http.StatusForbidden, router.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, response)
		return
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
}

func (r *Router) UpdateActivity(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Status(400)
		return
	}
	var activityUpdateRequest ActivityUpdateRequest
	err = c.ShouldBindJSON(&activityUpdateRequest)
	if err == nil {
		err = updateActivity(r.br.DB, userId, uint(id), activityUpdateRequest)
		if err != nil {
			c.JSON(http.StatusForbidden, router.ErrorResponse{Error: err.Error()})
			return
		}
		c.Status(http.StatusOK)
		return
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
}

func (r *Router) DeleteActivity(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Status(400)
		slog.Error("Could not process activity id when attempting a deletion", "error", err.Error(), "id", c.Param("id"))
		return
	}
	err = deleteActivity(r.br.DB, userId, uint(id))
	if err != nil {
		c.JSON(http.StatusForbidden, router.ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
