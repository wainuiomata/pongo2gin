package main

import (
	"log"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"

	pongo2gin "gitlab.com/go-box/pongo2gin/v1"
)

// GetAllData returns all posts
func GetAllData(c *gin.Context) {
	posts := []string{
		"Rob van der Linde",
		"John Curley",
		"Andrejs Cainikovs",
	}

	// Call the HTML method of the Context to render a template
	c.HTML(http.StatusOK, "index.html",
		pongo2.Context{
			"title": "hello pongo",
			"posts": posts,
		},
	)
}

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.Use(gin.Recovery())
	r.HTMLRender = pongo2gin.Default()
	r.GET("/", GetAllData)
	log.Fatal(r.Run(":8888"))
}
