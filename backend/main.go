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

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"your-project/backend/handler"
	"your-project/backend/repository"
	"your-project/backend/usecase"
)

func main() {
	// 環境変数から設定を取得
	dbConnectionString := os.Getenv("DATABASE_URL")
	if dbConnectionString == "" {
		dbConnectionString = "postgres://postgres:postgres@db:5432/slack_analysis?sslmode=disable"
	}
	
	slackToken := os.Getenv("SLACK_API_TOKEN")
	if slackToken == "" {
		log.Fatal("SLACK_API_TOKEN environment variable is required")
	}
	
	// データベース接続
	db, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	
	// 依存関係の初期化
	repo := repository.NewRepository(db)
	slackUsecase := usecase.NewSlackUsecase(repo, slackToken)
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
	router.POST("/api/initialize-users", slackHandler.InitializeUsersHandler)
	router.GET("/api/users", slackHandler.GetAllUsersHandler)
	
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