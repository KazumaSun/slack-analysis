// backend/handler/slack_handler.go
package handler

import (
	"backend/usecase"
	"fmt"
	"log"
	"net/http"
	"strconv" // 文字列を数値に変換するためにインポート

	"backend/repository" // repository.User を使うためにインポート
	"github.com/gin-gonic/gin"
)

type SlackHandler struct {
	slackUsecase *usecase.SlackUsecase
}

func NewSlackHandler(slackUsecase *usecase.SlackUsecase) *SlackHandler {
	return &SlackHandler{
		slackUsecase: slackUsecase,
	}
}

// InitializeUsersHandler はユーザー初期化APIのハンドラー
func (h *SlackHandler) InitializeUsersHandler(c *gin.Context) {
	err := h.slackUsecase.InitializeUsers()
	if err != nil {
		// エラーメッセージにエンドポイント情報を加えるなどしても良い
		log.Printf("Error in InitializeUsersHandler: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			// エラーラップされている場合、元のエラーも含めて返すことができる
			"error": fmt.Sprintf("Failed to initialize users: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Users initialized successfully",
	})
}

// InitializeChannelsHandler はチャンネル初期化APIのハンドラー (新規追加)
func (h *SlackHandler) InitializeChannelsHandler(c *gin.Context) {
	err := h.slackUsecase.InitializeChannels()
	if err != nil {
		log.Printf("Error in InitializeChannelsHandler: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to initialize channels: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Channels initialized successfully",
	})
}

// GetAllUsersHandler はユーザー情報取得APIのハンドラー
func (h *SlackHandler) GetAllUsersHandler(c *gin.Context) {
	users, err := h.slackUsecase.GetAllUsers()
	if err != nil {
		log.Printf("Error in GetAllUsersHandler: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to get users: %v", err),
		})
		return
	}
	// users が空スライスの場合もそのまま返す (JSONでは [] となる)
	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

// GetAllChannelsHandler はチャンネル情報取得APIのハンドラー (新規追加)
func (h *SlackHandler) GetAllChannelsHandler(c *gin.Context) {
	channels, err := h.slackUsecase.GetAllChannels()
	if err != nil {
		log.Printf("Error in GetAllChannelsHandler: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to get channels: %v", err),
		})
		return
	}

	// channels が空スライスの場合もそのまま返す (JSONでは [] となる)
	c.JSON(http.StatusOK, gin.H{
		"channels": channels,
	})
}

// UpdateUserHandler はユーザー情報更新APIのハンドラー (新規追加)
func (h *SlackHandler) UpdateUserHandler(c *gin.Context) {
	// 1. URL パラメータから ID を取得
	idStr := c.Param("id")
	// 文字列の ID を int 型に変換
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// ID が数値でない場合は 400 Bad Request を返す
		log.Printf("Error converting id parameter '%s' to int: %v", idStr, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid user ID format: %s", idStr),
		})
		return
	}

	// 2. リクエストボディの JSON を User 構造体にバインド
	var user repository.User
	// c.ShouldBindJSON はリクエストボディをパースし、バリデーションも行う（User構造体のタグに基づく）
	if err := c.ShouldBindJSON(&user); err != nil {
		// JSON の形式が不正、または必須フィールドが欠けている場合は 400 Bad Request を返す
		log.Printf("Error binding JSON for update user (id: %d): %v", id, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request body: %v", err),
		})
		return
	}

	// 3. Usecase層の更新メソッドを呼び出す
	err = h.slackUsecase.UpdateUser(id, user)
	if err != nil {
		// Usecase/Repository でエラーが発生した場合
		log.Printf("Error updating user (id: %d): %v", id, err)

		// Repository で "no user found" エラーを返した場合、404 Not Found を返すなどの分岐も可能
		// 例: if errors.Is(err, ...) or strings.Contains(err.Error(), "no user found") { ... }
		// ここでは単純に 500 Internal Server Error を返す
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to update user: %v", err),
		})
		return
	}

	// 4. 成功レスポンスを返す
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("User with id %d updated successfully", id),
		// 更新後のユーザー情報を返すこともできる (Usecaseが返すように変更が必要)
		// "user": updatedUser,
	})
}