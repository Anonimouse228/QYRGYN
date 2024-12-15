package routes

import (
	"QYRGYN/controllers"
	"QYRGYN/task1"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		if c.Request.Method == "POST" {
			if c.Request.FormValue("_method") == "DELETE" {
				c.Request.Method = "DELETE"

			}
			if c.Request.FormValue("_method") == "UPDATE" {
				c.Request.Method = "UPDATE"
			}

		}
		c.Next()
	})

	router.GET("task1", task1.Get)
	router.POST("task1", task1.Post)

	router.LoadHTMLGlob("views/templates/user/*")
	router.GET("/posts", controllers.GetPosts)
	router.GET("/posts/new", controllers.NewPost)
	router.POST("/posts", controllers.CreatePost)
	router.GET("/posts/:id/edit", controllers.EditPost)
	router.GET("/posts/:id", controllers.GetPost)
	router.POST("/posts/:id", controllers.UpdatePost)
	router.DELETE("/posts/:id", controllers.DeletePost)

	router.GET("/users", controllers.GetUsers)
	router.GET("/users/new", controllers.NewUserForm)
	router.POST("/users", controllers.CreateUser)
	router.GET("/users/:id", controllers.GetUser)
	router.GET("/users/:id/edit", controllers.EditUser)
	router.PATCH("/users/:id", controllers.UpdateUser)
	router.DELETE("/users/:id", controllers.DeleteUser)

}
