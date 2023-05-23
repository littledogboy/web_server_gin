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
	spiders.MRTHomeSpider(spiders.Meirentu.Doman, page, func(h spiders.Home, err error) {
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

// eg: /menuItems
func MenuItems(c *gin.Context) {
	fmt.Println("pong")
	spiders.MenuItems(func(data spiders.MenuData, err error) {
		c.JSON(http.StatusOK, data)
	})
}

// eg: /group?href=xxx&page=x
func GroupPage(c *gin.Context) {
	// href
	href := c.Query("href")
	// page
	page := c.Query("page")
	spiders.MRTGroupSpider(href, page, func(h spiders.Home, err error) {
		c.JSON(http.StatusOK, h)
	})
}
