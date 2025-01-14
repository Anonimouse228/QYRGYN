package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/gorilla/sessions"
	"net/http"
)

func RequireAdmin(c *gin.Context) {
	session, err := store.Get(c.Request, "session")
	if err != nil {

		fmt.Println("Error retrieving session:", err)
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	userRole := session.Values["role"]
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		c.Abort()
		return
	}
	c.Next()
}
