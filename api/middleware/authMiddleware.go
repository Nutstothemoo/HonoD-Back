package middleware

import (
	token "ginapp/api/tokens"
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
		c.Set("userId", claims.Uid)
		c.Next()
	}
}

func DealerAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
			ClientToken := c.Request.Header.Get("token")
			if ClientToken == "" {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "no authtoken header"})
					c.Abort()
					return
			}

			claims, err := token.ValidateToken(ClientToken)
			if err != "" {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
					c.Abort()
					return
			}

			if claims.Role != "dealer" {
					c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to access this resource"})
					c.Abort()
					return
			}

			c.Set("email", claims.Email)
			c.Set("userId", claims.Uid)
			c.Next()
	}
}

func AdminAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
			ClientToken := c.Request.Header.Get("token")
			if ClientToken == "" {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "no authtoken header"})
					c.Abort()
					return
			}

			claims, err := token.ValidateToken(ClientToken)
			if err != "" {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
					c.Abort()
					return
			}

			if claims.Role != "admin" {
					c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to access this resource"})
					c.Abort()
					return
			}

			c.Set("email", claims.Email)
			c.Set("userId", claims.Uid)
			c.Next()
	}
}