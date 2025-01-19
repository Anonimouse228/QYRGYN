package controllers

import (
	"QYRGYN/database"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ExecuteQueryHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "SQL_ENTRY.html", nil)
}

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

	// Fetch columns
	columns, err := rows.Columns()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch columns: %v", err)})
		return
	}

	// Prepare result data
	var result []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePointers := make([]interface{}, len(columns))
		for i := range values {
			valuePointers[i] = &values[i]
		}

		if err := rows.Scan(valuePointers...); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to scan result: %v", err)})
			return
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			row[col] = values[i]
		}
		result = append(result, row)
	}

	// Return readable query results
	c.JSON(http.StatusOK, gin.H{
		"query":   query,
		"status":  "success",
		"columns": columns,
		"rows":    result,
	})
}
