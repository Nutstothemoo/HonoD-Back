package middleware

import (
	token "ginapp/tokens"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authentification() gin.HandlerFunc {
	return func(c *gin.Context) {
		ClientToken := c.Request.Header.Get("token")
		if ClientToken == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no authtoken header"})
			c.Abort()
			return
		}

		claims, err := token.ValidateToken(ClientToken)
		if err !=""  {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("id", claims.Uid)
		c.Next()
	}
}