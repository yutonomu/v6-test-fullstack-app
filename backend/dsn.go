package main

import (
    "fmt"
    "net/url"
    "os"
    "strconv"
    "strings"
)

// inferDSNFromEnv: 環境変数からPostgresのDSNを推測して生成する。
// 優先度: APP_DB_* > POSTGRES_*。見つからなければ空文字を返す。
func inferDSNFromEnv() string {
    // devcontainer内のアプリからサービス名`db`で接続するケースを想定
    appUser := os.Getenv("APP_DB_USER")
    appPass := os.Getenv("APP_DB_PASSWORD")
    appName := os.Getenv("APP_DB_NAME")
    appHost := os.Getenv("APP_DB_HOST")
    appPort := os.Getenv("APP_DB_PORT")
    if appUser != "" && appName != "" {
        if appHost == "" {
            appHost = "db"
        }
        port := parsePort(appPort, 5432)
        return buildPostgresURL(appUser, appPass, appHost, port, appName, "disable")
    }

    // ホスト側からポートフォワードで接続するケース
    pgUser := os.Getenv("POSTGRES_USER")
    pgPass := os.Getenv("POSTGRES_PASSWORD")
    pgName := os.Getenv("POSTGRES_DB")
    pgPort := os.Getenv("POSTGRES_PORT")
    if pgUser != "" && pgName != "" {
        port := parsePort(pgPort, 5432)
        return buildPostgresURL(pgUser, pgPass, "localhost", port, pgName, "disable")
    }
    return ""
}

// parsePort: 文字列ポートをintに変換。失敗時はデフォルト値を返す。
func parsePort(s string, def int) int {
    if s == "" {
        return def
    }
    if p, err := strconv.Atoi(s); err == nil && p > 0 {
        return p
    }
    return def
}

// buildPostgresURL: Postgres用のURLを構築する（postgres://... 形式）。
func buildPostgresURL(user, pass, host string, port int, dbname, sslmode string) string {
    u := &url.URL{Scheme: "postgres"}
    if user != "" {
        u.User = url.UserPassword(user, pass)
    }
    if port > 0 {
        u.Host = fmt.Sprintf("%s:%d", host, port)
    } else {
        u.Host = host
    }
    if dbname != "" {
        u.Path = "/" + strings.TrimPrefix(dbname, "/")
    }
    q := url.Values{}
    if sslmode != "" {
        q.Set("sslmode", sslmode)
    }
    u.RawQuery = q.Encode()
    return u.String()
}

// redactPassword: DSNのパスワードをマスク（ログ出力用）。
func redactPassword(dsn string) string {
    parsed, err := url.Parse(dsn)
    if err != nil {
        return dsn
    }
    if parsed.User != nil {
        username := parsed.User.Username()
        if _, set := parsed.User.Password(); set {
            parsed.User = url.UserPassword(username, "****")
        }
    }
    return parsed.String()
}

