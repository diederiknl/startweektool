// cmd/server/main.go
package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"startweektool/internal/database"
	"startweektool/internal/handlers"
)

func main() {
	// Configure logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Initialize database
	db, err := database.InitDB("startweek.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize Gin
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

	// Admin routes
	r.GET("/admin", func(c *gin.Context) {
		c.HTML(200, "admin.html", gin.H{
			"title": "Admin - Startweek Tool",
		})
	})

	adminHandler := handlers.NewAdminHandler(db)
	r.POST("/api/admin/upload", adminHandler.UploadCSV)
	r.GET("/api/admin/preview", adminHandler.GetPreview)
	r.GET("/api/search", adminHandler.SearchStudents)

	// Start server
	log.Println("Server starting on http://localhost:8080")
	r.Run(":8080")
}
