package task1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "GET request received successfully",
	})
}

func Post(c *gin.Context) {
	var requestData map[string]interface{}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Invalid JSON format",
		})
		return
	}

	if value, exists := requestData["message"]; exists {

		if messageValue, ok := value.(string); ok {

			if messageValue == "" {
				c.JSON(http.StatusOK, gin.H{
					"status":  "success",
					"message": "Data successfully received",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status":  "success",
				"message": "Data successfully received",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "'message' field must be a string",
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"status":  "fail",
		"message": "The 'message' key is required",
	})
}
