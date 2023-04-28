package server

import (
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
)

func Start() {
	router := gin.Default()

	store := persistence.NewInMemoryStore(time.Minute)
	// Cached Page
	router.GET("/home", cache.CachePage(store, time.Minute, HomePage))
	router.GET("/detail", cache.CachePage(store, time.Hour, DetailPage))

	// ping
	router.Any("/ping", Ping)

	router.Run("0.0.0.0:8080")
}
