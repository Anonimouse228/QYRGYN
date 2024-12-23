package routes

import (
	"QYRGYN/controllers"
	"QYRGYN/middleware"
	"QYRGYN/task1"
	"QYRGYN/util"
	"github.com/gin-gonic/gin"
	"html/template"
)

func InitRoutes(router *gin.Engine) {
	r := gin.Default()

	router.SetFuncMap(template.FuncMap{
		"add": util.Add,
		"sub": util.Sub,
	})

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
	rl := middleware.NewRateLimiter(2, 5)

	// Task routes with rate limiter
	router.GET("task1", middleware.RateLimitMiddleware(rl), task1.Get)
	router.POST("task1", middleware.RateLimitMiddleware(rl), task1.Post)

	// Load templates
	router.LoadHTMLGlob("views/templates/user/*")

	// Post routes with rate limiter
	router.GET("/posts", middleware.RateLimitMiddleware(rl), controllers.GetPosts)
	router.GET("/posts/new", middleware.RateLimitMiddleware(rl), controllers.NewPost)
	router.POST("/posts", middleware.RateLimitMiddleware(rl), controllers.CreatePost)
	router.GET("/posts/:id/edit", controllers.EditPost)
	router.GET("/posts/:id", controllers.GetPost)
	router.POST("/posts/:id", middleware.RateLimitMiddleware(rl), controllers.UpdatePost)
	router.DELETE("/posts/:id", middleware.RateLimitMiddleware(rl), controllers.DeletePost)

	// User routes with rate limiter
	router.GET("/users", middleware.RateLimitMiddleware(rl), controllers.GetUsers)
	router.GET("/users/new", middleware.RateLimitMiddleware(rl), controllers.NewUserForm)
	router.POST("/users", middleware.RateLimitMiddleware(rl), controllers.CreateUser)
	router.GET("/users/:id", middleware.RateLimitMiddleware(rl), controllers.GetUser)
	router.GET("/users/:id/edit", controllers.EditUser)
	router.PATCH("/users/:id", middleware.RateLimitMiddleware(rl), controllers.UpdateUser)
	router.DELETE("/users/:id", middleware.RateLimitMiddleware(rl), controllers.DeleteUser)

}
