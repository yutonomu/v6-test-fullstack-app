package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// main: エントリポイント。設定読み込み、DB初期化、ルーター配線、サーバ起動のみを担当します。
func main() {
	// .env読み込み（backend/.env → プロジェクト直下../.env の順で試行）
	_ = godotenv.Load("../.env")

	// DB接続文字列の決定（環境変数優先、なければ推測）
	// dsn := os.Getenv("DATABASE_URL")
	// if dsn == "" {
	// 	dsn = inferDSNFromEnv()
	// 	if dsn != "" {
	// 		log.Printf("DATABASE_URL を環境から推測: %s", redactPassword(dsn))
	// 	}
	// }
	// if dsn == "" {
	// 	log.Fatal("DATABASE_URL が必要です (例: postgres://user:pass@host:5432/db?sslmode=disable)")
	// }

	dbHost := os.Getenv("APP_DB_HOST")
	dbUser := os.Getenv("APP_DB_USER")
	dbPassword := os.Getenv("APP_DB_PASSWORD")
	dbName := os.Getenv("APP_DB_NAME")
	dbPort := os.Getenv("APP_DB_PORT")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// DB接続
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("DB接続失敗: %v", err)
	}

	// マイグレーション（Noteテーブル）
	if err := db.AutoMigrate(&Note{}); err != nil {
		log.Fatalf("マイグレーション失敗: %v", err)
	}

	// ルーター生成（Gin + CORS + ハンドラ配線）
	router := NewRouter(db)

	// ポート決定と起動
	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Printf("listening on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
