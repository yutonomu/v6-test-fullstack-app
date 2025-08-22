package main

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

// NewRouter: アプリのHTTPルーター（Gin）を生成する。
// ここでミドルウェア適用とエンドポイント配線を行う。
func NewRouter(db *gorm.DB) *gin.Engine {
    r := gin.New()
    r.Use(gin.Logger(), gin.Recovery())
    r.Use(CORSMiddleware())

    // ヘルスチェック
    r.GET("/healthz", func(c *gin.Context) {
        c.String(200, "ok")
    })

    // ノート一覧/作成
    r.GET("/notes", NotesGetHandler(db))
    r.POST("/notes", NotesPostHandler(db))

    return r
}
