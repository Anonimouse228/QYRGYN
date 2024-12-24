package controllers

import (
	"QYRGYN/database"
	"QYRGYN/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CreatePost(c *gin.Context) {
	var post models.Post

	// Retrieve user ID directly from the context
	userID, exists := c.Get("userID")
	if !exists {
		// In case somehow the userID is not set in the context (shouldn't happen)
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "In case somehow the userID is not set in the context (shouldn't happen)"})
		return
	}

	// Assign the userID to the post
	post.UserId = int(userID.(uint))

	// Bind the form data to the post struct
	if err := c.ShouldBind(&post); err != nil {
		fmt.Println("Binding error:", err)
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Post received: %+v\n", post)

	// Save the post to the database
	if err := database.DB.Create(&post).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	// Redirect to posts page after successful creation
	c.Redirect(http.StatusFound, "/posts")
}

func GetPosts(c *gin.Context) {
	var posts []models.Post
	query := database.DB

	// Filtering
	if content := c.Query("content"); content != "" {
		query = query.Where("content LIKE ?", "%"+content+"%")
	}

	// Sorting
	sortBy := c.DefaultQuery("sort", "created_at") // Default sort by date
	order := c.DefaultQuery("order", "asc")
	query = query.Order(fmt.Sprintf("%s %s", sortBy, order))

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "5"))
	offset := (page - 1) * pageSize

	var total int64
	query.Model(&models.Post{}).Count(&total) // Count total records
	query = query.Offset(offset).Limit(pageSize)

	if err := query.Find(&posts).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Failed to retrieve posts: " + err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "posts.html", gin.H{
		"posts":      posts,
		"page":       page,
		"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
		"content":    c.Query("content"),
		"sort":       c.DefaultQuery("sort", "id"),
		"order":      c.DefaultQuery("order", "asc"),
	})

}

//func GetPosts(c *gin.Context) {
//	var posts []models.Post
//
//	if err := database.DB.Find(&posts).Error; err != nil {
//
//		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
//			"error": "Failed to retrieve posts: " + err.Error(),
//		})
//		return
//	}
//
//	c.HTML(http.StatusOK, "posts.html", gin.H{"posts": posts})
//}

func NewPost(c *gin.Context) {
	c.HTML(http.StatusOK, "new_post.html", nil)
}

func EditPost(c *gin.Context) {
	var post models.Post
	id := c.Param("id")
	if err := database.DB.First(&post, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Post not found"})
		return
	}
	c.HTML(http.StatusOK, "edit_post.html", gin.H{"post": post})
}

func GetPost(c *gin.Context) {
	var post models.Post
	id := c.Param("id")
	if err := database.DB.First(&post, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Post not found"})
		return
	}
	c.HTML(http.StatusOK, "post.html", gin.H{"post": post})
}

func UpdatePost(c *gin.Context) {
	var post models.Post
	id := c.Param("id")
	if err := database.DB.First(&post, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBind(&post); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": err.Error()})
		return
	}

	database.DB.Save(&post)
	c.Redirect(http.StatusFound, "/posts")
}

func DeletePost(c *gin.Context) {
	id := c.Param("id")
	fmt.Println("ID IS: ", id)
	if err := database.DB.Delete(&models.Post{}, id).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Could not delete post"})
		return
	}
	c.Redirect(http.StatusFound, "/posts")
}
