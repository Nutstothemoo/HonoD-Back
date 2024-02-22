package middleware

import (
    token "ginapp/api/tokens"
    "net/http"

    "github.com/gin-gonic/gin"
)

func authenticateWithRole(c *gin.Context, role string) {
    cookie, err := c.Request.Cookie("auth_token")
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "no auth_token cookie"})
        c.Abort()
        return
    }

    ClientToken := cookie.Value
    if ClientToken == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "empty auth_token cookie"})
        c.Abort()
        return
    }
    claims, err := token.ValidateToken(ClientToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err})
        c.Abort()
        return
    }

    if claims.Role != role && role != "" {
        c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to access this resource"})
        c.Abort()
        return
    }

    c.Set("email", claims.Email)
    c.Set("userId", claims.Uid)
    c.Next()
}

func Authentication() gin.HandlerFunc {
    return func(c *gin.Context) {
        authenticateWithRole(c, "")
    }
}

func DealerAuthentication() gin.HandlerFunc {
    return func(c *gin.Context) {
        authenticateWithRole(c, "dealer")
    }
}

func AdminAuthentication() gin.HandlerFunc {
    return func(c *gin.Context) {
        authenticateWithRole(c, "admin")
    }
}