package router

import (
	"time"

	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/sbondCo/Watcharr/config"
	"gorm.io/gorm"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

// TODO don't use generic ValueRequest and KeyValueRequest..
// each handler/service should define their own struct.

type ValueRequest struct {
	Value any `json:"value"`
}

type KeyValueRequest struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type BaseRouter struct {
	// Our database.
	DB *gorm.DB
	// Our base router group.
	Router *gin.RouterGroup
	// Our in-memory store used for cache.
	MemStore *persistence.InMemoryStore
	// Our server config.
	Cfg *config.ServerConfig
}

func NewBaseRouter(db *gorm.DB, rg *gin.RouterGroup, cfg *config.ServerConfig) *BaseRouter {
	return &BaseRouter{
		DB:       db,
		Router:   rg,
		MemStore: persistence.NewInMemoryStore(time.Hour * 24),
		Cfg:      cfg,
	}
}
