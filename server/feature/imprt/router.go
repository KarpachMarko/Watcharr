package imprt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sbondCo/Watcharr/feature/auth/authmiddleware"
	"github.com/sbondCo/Watcharr/router"
)

type Router struct {
	br *router.BaseRouter
	s  *Service
	ts *TraktService
}

func NewRouter(br *router.BaseRouter) *Router {
	return &Router{br: br}
}

func (r *Router) AddRoutes() {
	imprt := r.br.Router.Group("/import").Use(authmiddleware.AuthRequired(nil, r.br.Cfg))

	imprt.POST("", r.ImportContent)
	imprt.POST("/trakt", r.ImportTrakt)
}

// Import content (the client handle processing data and sends it to us in a uniform way).
func (r *Router) ImportContent(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	var ar ImportRequest
	err := c.ShouldBindJSON(&ar)
	if err == nil {
		response, err := r.s.ImportContent(r.br.DB, userId, ar)
		if err != nil {
			c.JSON(http.StatusForbidden, router.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, response)
		return
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
}

// Import Trakt.
func (r *Router) ImportTrakt(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	var ar TraktImportRequest
	err := c.ShouldBindJSON(&ar)
	if err == nil {
		response, err := r.ts.TraktImportWatched(r.br.DB, userId, ar.Username)
		if err != nil {
			c.JSON(http.StatusForbidden, router.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, response)
		return
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, router.ErrorResponse{Error: err.Error()})
}
