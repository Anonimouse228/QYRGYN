package main

import (
	"QYRGYN/config"
	"QYRGYN/database"
	"QYRGYN/middleware"
	"QYRGYN/routes"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {

	database.InitDatabase(config.GetDatabaseURL())

	router := gin.Default()

	middleware.SetupLogger()
	router.Use(middleware.Logger())

	log.Println("Server running on port localhost:8080/posts")

	routes.InitRoutes(router)

	router.Use(func(c *gin.Context) {
		if method := c.Request.FormValue("_method"); method != "" {
			c.Request.Method = method
		}
		c.Next()
	})

	router.Run(":8080")

}
