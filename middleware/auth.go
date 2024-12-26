package middleware

import (
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
		// Handle session retrieval failure
		fmt.Println("Error retrieving session:", err) // Log the error for debugging
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	// Fetch user ID from session
	userID, ok := session.Values["userID"].(uint) // Use 'int' instead of 'uint'
	if !ok || userID == 0 {                       // Ensure user ID is valid
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	// Store user ID in the context
	c.Set("userID", uint(userID)) // Cast to 'uint' if needed later

	// Continue processing
	c.Next()
}
