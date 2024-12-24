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

	router.LoadHTMLGlob("views/templates/user/*")

	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)

	// Protected routes
	auth := router.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/posts", middleware.RateLimitMiddleware(rl), controllers.GetPosts)
		auth.GET("/posts/new", middleware.RateLimitMiddleware(rl), controllers.NewPost)
		auth.POST("/posts", middleware.RateLimitMiddleware(rl), controllers.CreatePost)
		auth.GET("/posts/:id/edit", controllers.EditPost)
		auth.GET("/posts/:id", controllers.GetPost)
		auth.POST("/posts/:id", middleware.RateLimitMiddleware(rl), controllers.UpdatePost)
		auth.DELETE("/posts/:id", middleware.RateLimitMiddleware(rl), controllers.DeletePost)
	}

	// Task routes with rate limiter
	router.GET("task1", middleware.RateLimitMiddleware(rl), task1.Get)
	router.POST("task1", middleware.RateLimitMiddleware(rl), task1.Post)

	// Load templates

	// Post routes with rate limiter

	// User routes with rate limiter
	router.GET("/users", middleware.RateLimitMiddleware(rl), controllers.GetUsers)
	router.GET("/users/new", middleware.RateLimitMiddleware(rl), controllers.NewUserForm)
	router.POST("/users", middleware.RateLimitMiddleware(rl), controllers.CreateUser)
	router.GET("/users/:id", middleware.RateLimitMiddleware(rl), controllers.GetUser)
	router.GET("/users/:id/edit", controllers.EditUser)
	router.PATCH("/users/:id", middleware.RateLimitMiddleware(rl), controllers.UpdateUser)
	router.DELETE("/users/:id", middleware.RateLimitMiddleware(rl), controllers.DeleteUser)

}
