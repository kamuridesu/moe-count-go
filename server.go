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
	r.Run()
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
		user, err = searchForUser(DB, username)
		if err != nil {
			user, err = insertUserIntoDB(DB, username)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
			}
		}
		if err == nil {
			c.Data(200, "text/html; charset=utf-8", generateSVG(user.counter, IMAGES).Bytes())
			updateUserViewCount(DB, user)
		}
	}

}
