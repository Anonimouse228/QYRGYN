//package controllers
//
//import (
//	"QYRGYN/database"
//	"QYRGYN/models"
//	"net/http"
//
//	"github.com/gin-gonic/gin"
//)
//
//func CreateUser(c *gin.Context) {
//	var user models.User
//	if err := c.ShouldBind(&user); err != nil {
//		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": err.Error()})
//		return
//	}
//
//	if err := database.DB.Create(&user).Error; err != nil {
//		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
//		return
//	}
//
//	c.Redirect(http.StatusFound, "/users")
//}
//
//func GetAllUsersHtml(c *gin.Context) {
//	var users []models.User
//	database.DB.Preload("User").Find(&users)
//	c.HTML(http.StatusOK, "users.html", gin.H{
//		"users": users,
//	})
//}
//
//func NewUser(c *gin.Context) {
//	c.HTML(http.StatusOK, "new_user.html", nil)
//}
//
//func EditUser(c *gin.Context) {
//	var user models.User
//	id := c.Param("id")
//	if err := database.DB.First(&user, id).Error; err != nil {
//		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "User not found"})
//		return
//	}
//	c.HTML(http.StatusOK, "edit_user.html", gin.H{"user": user})
//}
//
//func GetUserHtml(c *gin.Context) {
//	var user models.User
//	id := c.Param("id")
//	if err := database.DB.First(&user, id).Error; err != nil {
//		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "User not found"})
//		return
//	}
//	c.HTML(http.StatusOK, "user_detail.html", gin.H{"user": user})
//}
//
//func UpdateUser(c *gin.Context) {
//	var user models.User
//	id := c.Param("id")
//	if err := database.DB.First(&user, id).Error; err != nil {
//		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "User not found"})
//		return
//	}
//
//	if err := c.ShouldBind(&user); err != nil {
//		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": err.Error()})
//		return
//	}
//
//	database.DB.Save(&user)
//	c.Redirect(http.StatusFound, "/users")
//}
//
//func DeleteUser(c *gin.Context) {
//	id := c.Param("id")
//	if err := database.DB.Delete(&models.User{}, id).Error; err != nil {
//		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Could not delete user"})
//		return
//	}
//	c.Redirect(http.StatusFound, "/users")
//}

package controllers

import (
	"QYRGYN/database"
	"QYRGYN/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "users.html", gin.H{"users": users})
}

func NewUserForm(c *gin.Context) {
	c.HTML(http.StatusOK, "new_user.html", nil)
}

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBind(&user); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/users")
}

func GetUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	if err := database.DB.First(&user, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "User not found"})
		return
	}
	c.HTML(http.StatusOK, "user.html", gin.H{"user": user})
}

func EditUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	if err := database.DB.First(&user, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Post not found"})
		return
	}
	c.HTML(http.StatusOK, "edit_user.html", gin.H{"user": user})
}

func UpdateUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	if err := database.DB.First(&user, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBind(&user); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": err.Error()})
		return
	}

	database.DB.Save(&user)
	c.Redirect(http.StatusFound, "/users")
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	fmt.Println("ID IS: ", id)
	if err := database.DB.Delete(&models.User{}, id).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Could not delete user"})
		return
	}
	c.Redirect(http.StatusFound, "/users")
}
