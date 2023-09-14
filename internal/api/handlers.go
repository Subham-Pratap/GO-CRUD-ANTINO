package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BlogPost struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Define your handlers here (getBlogPosts, getBlogPostByID, etc.)

func getBlogPosts(c *gin.Context, db *sql.DB) {
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

func getBlogPostByID(c *gin.Context, db *sql.DB) {
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

func createBlogPost(c *gin.Context, db *sql.DB) {
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

func updateBlogPost(c *gin.Context, db *sql.DB) {
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

func deleteBlogPost(c *gin.Context, db *sql.DB) {
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
