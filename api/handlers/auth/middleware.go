package handlerAPIAuth

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthSessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		userID := session.Get("user_id")
		if userID == nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "authentication required"},
			)
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
