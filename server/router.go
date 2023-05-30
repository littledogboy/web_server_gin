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
	router.GET("/menuItems", cache.CachePage(store, time.Hour, MenuItems))
	router.GET("/group", cache.CachePage(store, time.Hour, GroupPage))
	router.GET("/tags", cache.CachePage(store, time.Hour, Tags))
	router.GET("/tag", cache.CachePage(store, time.Hour, TagPage))

	// ping
	router.GET("/ping", Ping)

	router.Run("0.0.0.0:8080")
}
