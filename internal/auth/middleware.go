package auth

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" || len(authHeader) <= 7 || authHeader[:7] != "Bearer " {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or missing authorization header"})
            c.Abort()
            return
        }

        tokenString := authHeader[7:]
        token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })
        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
            c.Abort()
            return
        }

        claims, ok := token.Claims.(*Claims)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
            c.Abort()
            return
        }

        c.Set("user_id", claims.UserID)
        c.Set("email", claims.Email)
        c.Set("role", claims.Role)

        c.Next()
    }
}

func ProfileHandler(c *gin.Context) {
    email := c.GetString("email")
    role := c.GetString("role")

    c.JSON(http.StatusOK, gin.H{
        "email":   email,
        "role":    role,
        "message": "This is a protected profile endpoint",
    })
}
