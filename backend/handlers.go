package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NotesGetHandler: /notes GET（直近のノートをlimit件取得）
func NotesGetHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := 50
		if l := c.Query("limit"); l != "" {
			if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 500 {
				limit = v
			}
		}
		var notes []Note
		if err := db.Order("id desc").Limit(limit).Find(&notes).Error; err != nil {
			log.Printf("db error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch"})
			return
		}
		c.JSON(http.StatusOK, notes)
	}
}

// NotesPostHandler: /notes POST（ノート作成）
func NotesPostHandler(db *gorm.DB) gin.HandlerFunc {
	type reqBody struct {
		Content string `json:"content"`
	}
	return func(c *gin.Context) {
		var body reqBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}
		if len(body.Content) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "content is required"})
			return
		}
		note := Note{Content: body.Content}
		if err := db.Create(&note).Error; err != nil {
			log.Printf("db error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create"})
			return
		}
		c.JSON(http.StatusOK, note)
	}
}
