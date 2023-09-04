package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func serve() {
	r := gin.Default()
	r.LoadHTMLGlob("template/*")
	r.Static("/static", "./static/fonts")
	r.GET("/", getUserCounter)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	r.Run("0.0.0.0:80")
}

func getUserCounter(c *gin.Context) {
	username := c.Query("username")
	var user User
	var err error
	if username == "" {
		// c.JSON(404, gin.H{
		// 	"message": "Error! Missing username param",
		// })
		c.HTML(http.StatusOK, "404.tmpl", gin.H{
			"message": "Error! Missing username param",
		})
	} else {
		user, err = searchForUser(mainDatabase, username)
		if err != nil {
			user, err = insertUserIntoDB(mainDatabase, username)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
			}
		}
		if err == nil {
			c.Data(200, "text/html; charset=utf-8", generateSVG(user.counter, IMAGES).Bytes())
			updateUserViewCount(mainDatabase, user)
		}
	}

}
