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
	r.Run("0.0.0.0:8080")
}

func getUserCounter(c *gin.Context) {
	username := c.Query("username")
	var user User
	var err error
	if username == "" {
		c.HTML(http.StatusBadRequest, "404.tmpl", gin.H{
			"message": "Error! Missing username param",
		})
	} else {
		user, err = mainDatabase.searchForUser(username)
		if err != nil {
			user, err = mainDatabase.insertUserIntoDB(username)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
			}
		}
		if err == nil {
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Header("Vary", "Accept-Encoding")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
			c.Data(200, "image/svg+xml", generateSVG(user.counter, loadedImages).Bytes())
			mainDatabase.updateUserViewCount(user)
		}
	}

}
