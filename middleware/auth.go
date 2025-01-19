package middleware

import (
	"QYRGYN/database"
	"QYRGYN/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("your-secret-key"))

func AuthRequired(c *gin.Context) {
	// Retrieve session
	session, err := store.Get(c.Request, "session")
	if err != nil {

		fmt.Println("Error retrieving session:", err)
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	// Fetch user ID from session
	userID, ok := session.Values["userID"].(uint)

	if !ok || userID == 0 {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	var user models.User
	database.DB.First(&user, userID)
	if !user.Verified {
		c.JSON(http.StatusForbidden, gin.H{"error": "Email is not verified"})
		c.Abort()
		return
	}

	// Store user ID in the context
	c.Set("userID", uint(userID))

	c.Next()
}
