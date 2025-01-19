package main

import (
	"QYRGYN/config"
	"QYRGYN/database"
	"QYRGYN/middleware"
	"QYRGYN/routes"
	"github.com/gin-gonic/gin"
	"log"
)

func SetupRoutes() *gin.Engine {
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

func main() {
	router := SetupRoutes()

	port := config.GetPort()
	println("Port: ", port)
	if port == "" {
		log.Println("No port configured")
		port = "8080"
	}
	err := router.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}

}
