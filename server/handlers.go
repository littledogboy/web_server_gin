package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"web_server_gin/spiders"
)

// eg:/ping
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// eg: /home?page=1
func HomePage(c *gin.Context) {
	fmt.Println("pong")
	page := c.Query("page")
	spiders.RecommendSpider(page, func(h spiders.Home, err error) {
		c.JSON(http.StatusOK, h)
	})
}

// eg: /detail?href=https://meirentu.cc/pic/392051089446.html
func DetailPage(c *gin.Context) {
	fmt.Println("pong")
	// get href
	href := c.Query("href")
	spiders.DetailViewSpider(href, func(ad spiders.AlbumDetail) {
		c.JSON(http.StatusOK, ad)
	})
}
