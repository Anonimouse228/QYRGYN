package task1

import (
	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "GET request received successfully",
	})
}

func Post(c *gin.Context) {
	var requestData struct {
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(400, gin.H{
			"status":  "fail",
			"message": "Invalid JSON format",
		})
		return
	}

	if requestData.Message == "" {
		c.JSON(400, gin.H{
			"status":  "fail",
			"message": "Invalid JSON message",
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Data successfully received",
	})
}
