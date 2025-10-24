package profile

import (
	"net/http"

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
	profile := r.br.Router.Group("/profile").Use(authmiddleware.AuthRequired(nil, r.br.Cfg))

	// Get user profile details
	profile.GET("", r.GetProfile)
}

// Get user profile details
func (r *Router) GetProfile(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	response, err := getProfile(r.br.DB, userId)
	if err != nil {
		c.JSON(http.StatusForbidden, router.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}
