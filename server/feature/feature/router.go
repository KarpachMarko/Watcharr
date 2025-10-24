package feature

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sbondCo/Watcharr/feature/auth/authmiddleware"
	"github.com/sbondCo/Watcharr/router"
)

type Router struct {
	*router.BaseRouter
	service *Service
}

func NewRouter(br *router.BaseRouter, service *Service) *Router {
	return &Router{
		BaseRouter: br,
		service:    service,
	}
}

func (r *Router) AddRoutes() {
	feature := r.Router.Group("/features").Use(authmiddleware.AuthRequired(r.DB, r.Cfg))

	// Get enabled features (aka functionality)
	feature.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, r.service.GetEnabledFeatures(c.GetInt("userPermissions")))
	})
}
