package main

import "time"

// Note: メモエンティティ（ID, 内容, 作成日時）
type Note struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"createdAt"`
}

