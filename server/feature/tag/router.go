package tag

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sbondCo/Watcharr/feature/auth/authmiddleware"
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
	tag := r.br.Router.Group("/tag").Use(authmiddleware.AuthRequired(nil, r.br.Cfg))

	tag.GET("", r.GetTags)
	// TODO implement getting a tag (with pagination support)
	// tag.GET(":id", r.GetTag)
	tag.POST("", r.CreateTag)
	tag.PUT(":id", r.UpdateTag)
	tag.DELETE(":id", r.DeleteTag)
}

// Get all of our tags.
func (r *Router) GetTags(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	tags, err := r.service.GetTags(r.br.DB, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, router.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, tags)
}

// Get all items within one of our tags.
func (r *Router) GetTag(c *gin.Context) {
	// id, err := strconv.Atoi(c.Param("id"))
	// if err != nil {
	// 	slog.Error("getTag route failed to convert id param to int", "error", err)
	// 	c.Status(http.StatusBadRequest)
	// 	return
	// }
	// userId := c.MustGet("userId").(uint)
	// tags, err := getTag(r.br.DB, userId, uint(id))
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, router.ErrorResponse{Error: err.Error()})
	// 	return
	// }
	// c.JSON(http.StatusOK, tags)
}

// Create a tag.
func (r *Router) CreateTag(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	var tr TagAddRequest
	err := c.ShouldBindJSON(&tr)
	if err == nil {
		response, err := r.service.AddTag(r.br.DB, userId, tr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, router.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, response)
		return
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
}

func (r *Router) UpdateTag(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Status(400)
		slog.Error("tag update rote: failed to process tag id.", "error", err.Error(), "id", c.Param("id"))
		return
	}
	var tr TagAddRequest
	err = c.ShouldBindJSON(&tr)
	if err == nil {
		err := r.service.UpdateTag(r.br.DB, userId, uint(id), tr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, router.ErrorResponse{Error: err.Error()})
			return
		}
		c.Status(http.StatusOK)
		return
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
}

func (r *Router) DeleteTag(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Status(400)
		slog.Error("tag delete rote: failed to process tag id.", "error", err.Error(), "id", c.Param("id"))
		return
	}
	err = r.service.DeleteTag(r.br.DB, userId, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, router.ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
