package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

var db *sql.DB

// InitRoutes initializes all API routes
func InitRoutes(router *gin.Engine) {
	// Create a group for API routes, e.g., /api/v1
	apiGroup := router.Group("/api/v1")

	apiGroup.GET("/posts", func(c *gin.Context) {
		getBlogPosts(c, db)
	})

	apiGroup.GET("/posts/:id", func(c *gin.Context) {
		getBlogPostByID(c, db)
	})

	apiGroup.POST("/posts", func(c *gin.Context) {
		createBlogPost(c, db)
	})

	apiGroup.PUT("/posts/:id", func(c *gin.Context) {
		updateBlogPost(c, db)
	})

	apiGroup.DELETE("/posts/:id", func(c *gin.Context) {
		deleteBlogPost(c, db)
	})
}
