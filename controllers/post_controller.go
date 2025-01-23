package controllers

import (
	"QYRGYN/database"
	"QYRGYN/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"time"
)

func CreatePost(c *gin.Context) {
	var post models.Post

	userID, exists := c.Get("userID")
	if !exists {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "In case somehow the userID is not set in the context (shouldn't happen lol)"})
		return
	}

	post.UserId = int(userID.(uint))

	if len(post.Content) > 228 {
		c.HTML(http.StatusConflict, "error.html", gin.H{"error": "Post content too long. Should be < 228"})
		return
	}

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

func NewPostHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "new_post.html", nil)
}

func UpdatePostHTML(c *gin.Context) {
	var post models.Post
	id := c.Param("id")
	if err := database.DB.First(&post, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Post not found"})
		return
	}
	c.HTML(http.StatusOK, "edit_post.html", gin.H{"post": post})
}

func GetPost(c *gin.Context) {
	// A custom struct with username
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

	// Join to display not authorId but his username
	err := database.DB.Table("posts").
		Select("posts.id, posts.user_id, posts.content, posts.likes, posts.created_at, posts.updated_at, users.username").
		Joins("LEFT JOIN users ON users.id = posts.user_id").
		Where("posts.id = ?", id).
		First(&post).Error
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

	if len(post.Content) > 228 {
		c.HTML(http.StatusConflict, "error.html", gin.H{"error": "Content too large. Should be < 228"})
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

func ToggleLike(c *gin.Context) {
	userID, _ := c.Get("userID")

	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid post ID"})
		return
	}

	// Check if the user already liked the post
	var like models.Like
	err = database.DB.Where("post_id = ? AND user_id = ?", postID, userID).First(&like).Error
	if err == nil {
		// If the like exists, remove it
		if err := database.DB.Delete(&like).Error; err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Failed to remove like"})
			return
		}

		// Decrement the likes count in the post table
		if err := database.DB.Model(&models.Post{}).
			Where("id = ?", postID).
			Update("likes", gorm.Expr("likes - 1")).Error; err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Failed to update post likes"})
			return
		}

		c.Redirect(http.StatusFound, "/posts")
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Database error"})
		return
	}

	// If the like does not exist, add a new one
	newLike := models.Like{
		PostId: uint(postID),
		UserId: userID.(uint),
	}
	if err := database.DB.Create(&newLike).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Failed to add like"})
		return
	}

	// Increment the likes count in the post table
	if err := database.DB.Model(&models.Post{}).
		Where("id = ?", postID).
		Update("likes", gorm.Expr("likes + 1")).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Failed to update post likes"})
		return
	}

	c.Redirect(http.StatusFound, "/posts")
}
