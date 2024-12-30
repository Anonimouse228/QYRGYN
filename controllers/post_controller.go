package controllers

import (
	"QYRGYN/database"
	"QYRGYN/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func CreatePost(c *gin.Context) {
	var post models.Post

	userID, exists := c.Get("userID")
	if !exists {

		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "In case somehow the userID is not set in the context (shouldn't happen)"})
		return
	}

	post.UserId = int(userID.(uint))

	if err := c.ShouldBind(&post); err != nil {
		fmt.Println("Binding error:", err)
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Post received: %+v\n", post)

	if err := database.DB.Create(&post).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/posts")
}

//	func GetPosts(c *gin.Context) {
//		//TODO: Make pageSize customizable through dropdown menu
//		var posts []models.Post
//		query := database.DB
//
//		// Filtering
//		if content := c.Query("content"); content != "" {
//			query = query.Where("content LIKE ?", "%"+content+"%")
//		}
//
//		// Sorting
//		sortBy := c.DefaultQuery("sort", "created_at") // Default sort by date
//		order := c.DefaultQuery("order", "asc")
//		query = query.Order(fmt.Sprintf("%s %s", sortBy, order))
//
//		// Pagination
//		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
//		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "5"))
//		offset := (page - 1) * pageSize
//
//		var total int64
//		query.Model(&models.Post{}).Count(&total) // Count total records
//		query = query.Offset(offset).Limit(pageSize)
//
//		if err := query.Find(&posts).Error; err != nil {
//			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
//				"error": "Failed to retrieve posts: " + err.Error(),
//			})
//			return
//		}
//
//		c.HTML(http.StatusOK, "posts.html", gin.H{
//			"posts":      posts,
//			"page":       page,
//			"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
//			"content":    c.Query("content"),
//			"sort":       c.DefaultQuery("sort", "id"),
//			"order":      c.DefaultQuery("order", "asc"),
//		})
//
// }
func GetPosts(c *gin.Context) {
	//TODO: Make pageSize customizable through dropdown menu
	var posts []struct {
		ID        int       `json:"id"`
		UserId    int       `json:"userid"`
		Content   string    `json:"content"`
		Likes     int       `json:"likes"`
		CreatedAt time.Time `json:"createdat"`
		UpdatedAt time.Time `json:"updatedat"`
		Username  string    `json:"username"`
	}

	query := database.DB.Table("posts").
		Select("posts.id, posts.user_id, posts.content, posts.likes, posts.created_at, posts.updated_at, users.username").
		Joins("LEFT JOIN users ON users.id = posts.user_id")

	// Filtering
	if content := c.Query("content"); content != "" {
		query = query.Where("posts.content LIKE ?", "%"+content+"%")
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
	query.Count(&total) // Count total records
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
	// Define a custom struct to include username
	var post struct {
		ID        int       `json:"id"`
		UserId    int       `json:"userid"`
		Username  string    `json:"username"`
		Content   string    `json:"content"`
		Likes     int       `json:"likes"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	id := c.Param("id")

	// Perform join query with GORM
	err := database.DB.Table("posts").
		Select("posts.id, posts.user_id, posts.content, posts.likes, posts.created_at, posts.updated_at, users.username").
		Joins("LEFT JOIN users ON users.id = posts.user_id").
		Where("posts.id = ?", id).
		First(&post).Error // Pass a pointer to the struct
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Post not found"})
		return
	}

	// Render the post with the username included
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
