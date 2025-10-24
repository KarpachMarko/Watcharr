package setup

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sbondCo/Watcharr/feature/auth"
	"github.com/sbondCo/Watcharr/feature/setup/setupglob"
	"github.com/sbondCo/Watcharr/router"
	"gorm.io/gorm"
)

type AuthProvider interface {
	RegisterFirstUser(urr *auth.UserRegisterRequest, db *gorm.DB) (auth.AuthResponse, error)
}

type Router struct {
	br           *router.BaseRouter
	authProvider AuthProvider
}

func NewRouter(br *router.BaseRouter) *Router {
	return &Router{br: br}
}

// Since we cannot remove these setup routes after they are registered,
// each route/service should ensure we are still in setup before continuing.
// After server restart, these routes shouldn't exist if setup finished
// (currently it is finished if a user is created).
//
// Each controller can check ServerInSetup var first, then each service
// can double check what it needs to (eg create_admin service, registerFirstUser,
// will check that no users exist).
func (r *Router) AddRoutes() {
	setup := r.br.Router.Group("/setup")

	// Server setup routes are being added, so we are in setup now.
	setupglob.ServerInSetup = true

	setup.POST("/create_admin", r.CreateAdmin)
}

// Create first user (which will be an admin).
func (r *Router) CreateAdmin(c *gin.Context) {
	if !setupglob.ServerInSetup {
		c.JSON(http.StatusForbidden, router.ErrorResponse{Error: "not in setup"})
		return
	}
	var user auth.UserRegisterRequest
	if c.ShouldBindJSON(&user) == nil {
		response, err := r.authProvider.RegisterFirstUser(&user, r.br.DB)
		if err != nil {
			c.JSON(http.StatusForbidden, router.ErrorResponse{Error: err.Error()})
			return
		} else {
			// Set in setup to false after first user registered successfully
			setupglob.ServerInSetup = false
		}
		c.JSON(http.StatusOK, response)
		return
	}
	c.Status(400)
}
