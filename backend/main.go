// package main

// import (
// 	"github.com/gin-gonic/gin"
// )

// func main() {
//     r := gin.Default()

//     r.GET("/ping", func(c *gin.Context) {
//         c.JSON(200, gin.H{"message": "pong"})
//     })

//     r.Run(":8080")
// }

// backend/main.go
package main

import (
	"database/sql"
	"log"
	"os"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"backend/handler"
	"backend/repository"
	"backend/usecase"
)

func main() {
	// 環境変数から設定を取得
	// 環境変数から設定を取得
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// データベース接続文字列の構築
	dbConnectionString := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", 
									dbUser, dbPassword, dbHost, dbName)

	slackTokenUser := os.Getenv("SLACK_API_TOKEN_USER")
	if slackTokenUser == "" {
		log.Fatal("SLACK_API_TOKEN_USER environment variable is required")
	}

	slackTokenBot := os.Getenv("SLACK_API_TOKEN_BOT")
	if slackTokenBot == "" {
		log.Fatal("SLACK_API_TOKEN_BOT environment variable is required")
	}
	
	// データベース接続
	db, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	
	// 依存関係の初期化
	repo := repository.NewRepository(db)
	slackUsecase := usecase.NewSlackUsecase(repo, slackTokenUser, slackTokenBot)
	slackHandler := handler.NewSlackHandler(slackUsecase)
	
	// Ginルーターの設定
	router := gin.Default()
	
	// CORS設定
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	
	// ルート定義
	router.GET("/users", slackHandler.GetAllUsersHandler)           // GET /users
	router.POST("/users/init", slackHandler.InitializeUsersHandler) // POST /users/init
	router.GET("/channels", slackHandler.GetAllChannelsHandler)     // GET /channels
	router.POST("/channels/init", slackHandler.InitializeChannelsHandler) // POST /channels/init
	
	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Server started on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}