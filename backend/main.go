package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// ローカル開発用に.envファイルを読み込む（本番環境では無視される）
	_ = godotenv.Load("../.env")

	// --- 接続文字列 (DSN) の決定 ---
	// 優先度1: DATABASE_URL (本番環境: Supabaseの接続文字列)
	dsn := os.Getenv("DATABASE_URL")

	// 優先度2: ローカル開発用の環境変数 (APP_DB_*)
	if dsn == "" {
		log.Println("DATABASE_URL is not set, falling back to local Docker setup...")
		dbHost := os.Getenv("APP_DB_HOST")
		dbUser := os.Getenv("APP_DB_USER")
		dbPassword := os.Getenv("APP_DB_PASSWORD")
		dbName := os.Getenv("APP_DB_NAME")
		dbPort := os.Getenv("APP_DB_PORT")

		if dbUser != "" && dbName != "" && dbHost != "" {
			dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
				dbHost, dbPort, dbUser, dbPassword, dbName)
		}
	}

	if dsn == "" {
		log.Fatal("Database connection string is not configured. Set DATABASE_URL or APP_DB_* variables.")
	}

	// DB接続
	var db *gorm.DB
	var err error

	// 接続リトライ処理を追加
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("Successfully connected to the database.")
			break
		}
		log.Printf("Failed to connect to database (attempt %d/5): %v", i+1, err)
		log.Println("Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
	}

	// 5回リトライしても接続できなかった場合は終了する
	if err != nil {
		log.Fatalf("Could not connect to the database after several retries: %v", err)
	}

	// マイグレーション
	if err := db.AutoMigrate(&Note{}); err != nil {
		log.Fatalf("マイグレーション失敗: %v", err)
	}

	// ルーター生成
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
