package main

import (
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
)

// CORSMiddleware: 簡易なCORS対応。OPTIONSは204で即時応答。
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        allowOrigin := os.Getenv("CORS_ALLOW_ORIGIN")
        if allowOrigin == "" {
            allowOrigin = "*"
        }
        c.Header("Access-Control-Allow-Origin", allowOrigin)
        c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if c.Request.Method == http.MethodOptions {
            c.Status(http.StatusNoContent)
            c.Abort()
            return
        }
        c.Next()
    }
}
