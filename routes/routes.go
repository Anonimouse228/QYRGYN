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

	// Create an instance of a rate limiter
	rl := middleware.NewRateLimiter(2, 5)

	// Load templates
	router.LoadHTMLGlob("views/templates/user/*")

	// Register and login
	router.GET("/register", middleware.RateLimitMiddleware(rl), controllers.RegisterHTML)
	router.GET("/login", middleware.RateLimitMiddleware(rl), controllers.LoginHTML)
	router.POST("/register", middleware.RateLimitMiddleware(rl), controllers.Register)
	router.POST("/login", middleware.RateLimitMiddleware(rl), controllers.Login)
	router.POST("/logout", middleware.RateLimitMiddleware(rl), controllers.Logout)
	router.GET("/verify", middleware.RateLimitMiddleware(rl), controllers.VerifyEmail)
	// Protected routes
	auth := router.Group("/")
	auth.Use(middleware.AuthRequired)

	// Post routes
	auth.GET("/posts", middleware.RateLimitMiddleware(rl), controllers.GetPosts)
	auth.GET("/posts/new", middleware.RateLimitMiddleware(rl), controllers.NewPost)
	auth.POST("/posts", middleware.RateLimitMiddleware(rl), controllers.CreatePost)
	auth.GET("/posts/:id/edit", controllers.EditPost)
	auth.GET("/posts/:id", controllers.GetPost)
	auth.POST("/posts/:id", middleware.RateLimitMiddleware(rl), controllers.UpdatePost)
	auth.DELETE("/posts/:id", middleware.RateLimitMiddleware(rl), controllers.DeletePost)

	// User thingies
	auth.GET("/users/:id", controllers.GetUserProfile)
	auth.GET("/users/:id/edit", controllers.UpdateUserHTML)
	auth.POST("/users/:id", controllers.UpdateUserProfile)

	// Email system
	auth.GET("/helpdesk", controllers.HelpdeskPageHTML)
	auth.POST("/helpdesk", controllers.HelpdeskController)

	// First assignment, first task
	router.GET("task1", middleware.RateLimitMiddleware(rl), task1.Get)
	router.POST("task1", middleware.RateLimitMiddleware(rl), task1.Post)

	// User routes with rate limiter
	//router.GET("/users", middleware.RateLimitMiddleware(rl), controllers.GetUsers)
	//router.GET("/users/new", middleware.RateLimitMiddleware(rl), controllers.NewUserForm)
	//router.POST("/users", middleware.RateLimitMiddleware(rl), controllers.CreateUser)
	//router.GET("/users/:id", middleware.RateLimitMiddleware(rl), controllers.GetUser)
	//router.GET("/users/:id/edit", controllers.EditUser)
	//router.PATCH("/users/:id", middleware.RateLimitMiddleware(rl), controllers.UpdateUser)
	//router.DELETE("/users/:id", middleware.RateLimitMiddleware(rl), controllers.DeleteUser)

}
