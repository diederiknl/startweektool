// cmd/server/main.go
package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()

	// Serve static files
	r.Static("/static", "./web/static")
	r.LoadHTMLGlob("web/templates/*")

	// Base route
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "Startweek Tool - Rotterdam Academy",
		})
	})

	// Start server
	log.Println("Server starting on http://localhost:8080")
	r.Run(":8080")
}
