package middleware

import (
    "github.com/gin-gonic/gin"
)

// CORS enables cross-origin requests and handles OPTIONS preflight with 204.
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.GetHeader("Origin")
        if origin == "" {
            origin = "*"
        }
        c.Header("Access-Control-Allow-Origin", origin)
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Requested-With, Accept, Origin")
        c.Header("Vary", "Origin")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    }
}

