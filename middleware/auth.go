package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("your-secret-key"))

func AuthRequired(c *gin.Context) {
	session, _ := store.Get(c.Request, "session")

	userID, ok := session.Values["user_id"].(uint)
	if !ok || userID == 0 {

		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	// Store user ID in the context
	c.Set("userID", userID)

	c.Next()
}
