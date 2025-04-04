package main

import (
	"database/sql" // データベース操作のためのパッケージ
	"fmt"          // フォーマット操作のためのパッケージ
	"log"          // ログ出力のためのパッケージ
	"os"           // 環境変数操作のためのパッケージ

	"github.com/gin-contrib/cors" // CORS設定のためのパッケージ
	"github.com/gin-gonic/gin"    // Ginフレームワークのためのパッケージ
	_ "github.com/lib/pq"         // PostgreSQLドライバのためのパッケージ(初期化関数のみ使用)
)

func main() {
	// 環境変数から設定を取得
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// データベース接続文字列の構築
	dbConnectionString := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbName)

	// Slack APIトークンの取得
	slackTokenBot := os.Getenv("SLACK_API_TOKEN_BOT")
	slackTokenUser := os.Getenv("SLACK_API_TOKEN_USER")
	// Slack APIトークンの確認
	if slackTokenBot == "" || slackTokenUser == "" {
		log.Fatal("SLACK_API_TOKEN_BOT or SLACK_API_TOKEN_USER is not set")
	}

	// DB接続
	db, err := sql.Open("postgres", dbConnectionString)
	// 接続の確認
	if err != nil {
		println("DB接続失敗: " + err.Error())
		panic(err) // panicは、エラーが発生した場合にプログラムを強制終了させる
	} else {
		println("DB接続成功")
	}

	// db.CloseでDB接続を閉じる
	defer db.Close() // defarで、関数が終了する際にdb.Close()を実行する

	// Ginのルーターを作成
	router := gin.Default()

	// CORS設定
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},  // アクセスを許可するオリジン
		AllowMethods:     []string{"GET", "POST"},            // アクセスを許可するHTTPメソッド
		AllowHeaders:     []string{"Origin", "Content-Type"}, // アクセスを許可するHTTPヘッダー
		ExposeHeaders:    []string{"Content-Length"},         // レスポンスで公開するHTTPヘッダー
		AllowCredentials: true,                               // 認証情報を許可する
	}))

	// GETリクエストのルートを定義
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// データを挿入するクエリ
	insertSQL := "INSERT INTO users (slack_id, name) VALUES ($1, $2) RETURNING id"
	// データを挿入する
	var newId int
	err = db.QueryRow(insertSQL, "U12345678", "John Doe").Scan(&newId) // データを挿入
	if err != nil {
		// UNIQUE制約違反などの可能性も考慮
		log.Printf("データ挿入に失敗しました: %v", err)
		// 実際のアプリケーションではエラーの種類によって処理を分岐する
	} else {
		fmt.Printf("新しいユーザーを挿入しました ID: %d\n", newId)
	}

	// データを取得するクエリ
	router.GET("/users", func(c *gin.Context) {
		var json string

		// データを取得する
		rows, err := db.Query("SELECT * FROM users")
		if err != nil {
			log.Fatalf("データ取得失敗: %v", err)
		} else {
			println("データ取得成功")
		}
		defer rows.Close() // rowsを閉じる
		// データを表示
		for rows.Next() {
			var id int
			var slack_id string
			var name string
			var attributes any
			var created_at string

			err := rows.Scan(&id, &slack_id, &name, &attributes, &created_at) // データを変数に格納
			if err != nil {
				log.Fatalf("データスキャン失敗: %v", err)
			}
			fmt.Printf("ID: %d, Name: %s\n", id, name)          // データを表示
			json += fmt.Sprintf("ID: %d, Name: %s\n", id, name) // JSON形式でデータを格納
		}
		c.JSON(200, gin.H{ // レスポンスをJSON形式で返す
			"message": "データ取得成功",
			"data":    json,
		})
	})
	// _, err = db.Exec(createTableSQL) // Execは結果セットを返さないクエリ用

	// ポート番号の取得
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// サーバー起動
	log.Printf("%s番ポートでサーバーを起動", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("サーバー起動失敗: %v", err)
	}
}
