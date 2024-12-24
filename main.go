package main

import (
	"QYRGYN/database"
	"QYRGYN/routes"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {

	//database.InitDatabase(config.GetDatabaseURL())
	database.InitDatabase("host=localhost port=5432 user=postgres dbname=qyrgyn sslmode=disable")
	router := gin.Default()

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
