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

	log.Println("Server running on port localhost:8080/posts")

	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}

}
