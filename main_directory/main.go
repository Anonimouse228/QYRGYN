package main_directory

import (
	"QYRGYN/config"
	"QYRGYN/database"
	"QYRGYN/middleware"
	"QYRGYN/routes"
	"github.com/gin-gonic/gin"
	"log"
)

func SetupRouter() *gin.Engine {
	database.InitDatabase(config.GetDatabaseURL())

	router := gin.Default()

	middleware.SetupLogger()
	router.Use(middleware.Logger())

	routes.InitRoutes(router)

	router.Use(func(c *gin.Context) {
		if method := c.Request.FormValue("_method"); method != "" {
			c.Request.Method = method
		}
		c.Next()
	})
	return router
}

//func main() {
//	router := SetupRouter()
//
//	log.Println("Server running on port localhost:8080/posts")
//	err := router.Run(config.GetPort())
//	if err != nil {
//		log.Fatal(err)
//	}
//}

func main() {
	port := config.GetPort()
	if port == "" {
		port = "8080" // Default port for local development
	}

	router := SetupRouter()
	log.Println("Server running on port " + port)
	err := router.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}
