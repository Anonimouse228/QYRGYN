package controllers

import (
	"QYRGYN/database"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ExecuteQuery(c *gin.Context) {
	query := c.DefaultPostForm("sqlQuery", "")

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SQL query cannot be empty"})
		return
	}

	// Execute the query
	rows, err := database.DB.Raw(query).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to execute query: %v", err)})
		return
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		columns, _ := rows.Columns()
		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(interface{})
		}

		err := rows.Scan(values...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to scan result: %v", err)})
			return
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			row[col] = *(values[i].(*interface{}))
		}
		result = append(result, row)
	}

	// Return the query results
	c.JSON(http.StatusOK, gin.H{"result": result})
}
