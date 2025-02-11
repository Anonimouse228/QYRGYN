package routes

import (
	"QYRGYN/controllers"
	"QYRGYN/middleware"
	"QYRGYN/task1"
	"QYRGYN/util"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
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
	router.Use(middleware.RateLimitMiddleware(rl))

	// Admin stuff
	adminRoutes := router.Group("/admin")
	adminRoutes.Use(middleware.RequireAdmin)
	adminRoutes.Use(middleware.AuthRequired)
	adminRoutes.Use(middleware.RateLimitMiddleware(rl))
	{
		adminRoutes.GET("/users", controllers.GetUsers)
		adminRoutes.POST("/users", controllers.CreateUser)
		adminRoutes.GET("/users/new", controllers.CreateUserHTML)
		adminRoutes.PATCH("/users/:id", controllers.UpdateUser)
		adminRoutes.GET("/users/edit/:id", controllers.AdminUpdateUserHTML)
		//adminRoutes.GET("/users/:id", controllers.AdminGetUser)
		adminRoutes.POST("/users/:id/delete", controllers.DeleteUser)
		//adminRoutes.GET("/execute-query", controllers.ExecuteQueryHTML)
		//adminRoutes.POST("/execute-query", controllers.ExecuteQuery)
		adminRoutes.GET("/execute-query", controllers.ExecuteQueryHTML)
		adminRoutes.POST("/execute-query", controllers.ExecuteQuery)

		adminRoutes.GET("/chats", controllers.AdminChatList)
		adminRoutes.GET("/chat/:chatID", controllers.AdminChat)
		adminRoutes.POST("/chat/:chatID/close", controllers.AdminCloseChat)
		adminRoutes.POST("/chat/:chatID/send", controllers.AdminSendMessage)

		adminRoutes.GET("/", controllers.AdminDashboard)

	}

	router.GET("/execute-query", controllers.ExecuteQueryHTML)
	router.POST("/execute-query", controllers.ExecuteQuery)

	// Register and login
	router.GET("/", func(c *gin.Context) { c.Redirect(http.StatusFound, "/login") })
	router.GET("/register", controllers.RegisterHTML)
	router.GET("/login", controllers.LoginHTML)
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)
	router.POST("/logout", controllers.Logout)
	router.GET("/verify", controllers.VerifyEmail)

	// Protected routes
	auth := router.Group("/")
	auth.Use(middleware.AuthRequired)

	// Post routes
	auth.GET("/posts", controllers.GetPosts)
	auth.GET("/posts/new", controllers.NewPostHTML)
	auth.POST("/posts", controllers.CreatePost)
	auth.GET("/posts/:id/edit", controllers.UpdatePostHTML)
	auth.GET("/posts/:id", controllers.GetPost)
	auth.POST("/posts/:id", controllers.UpdatePost)
	auth.DELETE("/posts/:id", controllers.DeletePost)
	auth.POST("/posts/:id/like", controllers.ToggleLike)

	// User thingies
	auth.GET("/profile", controllers.GetUserProfile)
	auth.GET("/users/:id/edit", controllers.UpdateUserHTML)
	auth.POST("/users/:id", controllers.UpdateUserProfile)

	// Email system
	auth.GET("/helpdesk", controllers.HelpdeskPageHTML)
	auth.POST("/helpdesk", controllers.HelpdeskController)

	// WebSocket
	auth.GET("/wss", controllers.HandleConnections)
	go controllers.HandleMessages()
	auth.POST("/chat/start/:id", controllers.StartChat)
	auth.GET("/chat/:chatID", controllers.ChatPage)
	auth.POST("/chat/:chatID/send", controllers.SendMessage)

	// First assignment, first task
	router.GET("task1", task1.Get)
	router.POST("task1", task1.Post)

	// Payment routes
	auth.GET("/payment", controllers.PaymentPage) // Show payment page
	//auth.POST("/payment", controllers.Payment)              // Handle payment submission
	auth.POST("/payment", controllers.ProcessPayment) // Process payment through external service

}
