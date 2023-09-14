package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"strconv"
)

var db *sql.DB

type BlogPost struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func main() {

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dsn := "root:@tcp(localhost:3306)/myblogdb"
	dsnn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	fmt.Println(dsnn)
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Initialize the Gin router
	router := gin.Default()

	// Define routes
	router.GET("/api/v1/posts", getBlogPosts)
	router.GET("/api/v1/posts/:id", getBlogPostByID)
	router.POST("/api/v1/posts", createBlogPost)
	router.PUT("/api/v1/posts/:id", updateBlogPost)
	router.DELETE("/api/v1/posts/:id", deleteBlogPost)

	// Start the server
	port := 8080
	router.Run(fmt.Sprintf(":%d", port))
}

func getBlogPosts(c *gin.Context) {
	rows, err := db.Query("SELECT id, title, content FROM blog_posts")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve blog posts"})
		return
	}
	defer rows.Close()

	var posts []BlogPost
	for rows.Next() {
		var post BlogPost
		if err := rows.Scan(&post.ID, &post.Title, &post.Content); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan rows"})
			return
		}
		posts = append(posts, post)
	}

	c.JSON(http.StatusOK, posts)
}

func getBlogPostByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	var post BlogPost
	err = db.QueryRow("SELECT id, title, content FROM blog_posts WHERE id=?", id).Scan(&post.ID, &post.Title, &post.Content)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve the post"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func createBlogPost(c *gin.Context) {
	var newPost BlogPost
	if err := c.BindJSON(&newPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	//c.JSON(http.StatusOK, gin.H{"Value": newPost})
	_, err := db.Exec("INSERT INTO blog_posts (title, content) VALUES (?, ?)", newPost.Title, newPost.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create the post"})
		return
	}

	c.JSON(http.StatusCreated, newPost)
}

func updateBlogPost(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var updatedPost BlogPost
	if err := c.BindJSON(&updatedPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	_, err = db.Exec("UPDATE blog_posts SET title=?, content=? WHERE id=?", updatedPost.Title, updatedPost.Content, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the post"})
		return
	}

	c.JSON(http.StatusOK, updatedPost)
}

func deleteBlogPost(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	_, err = db.Exec("DELETE FROM blog_posts WHERE id=?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the post"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
