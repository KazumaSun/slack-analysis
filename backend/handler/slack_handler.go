// backend/handler/slack_handler.go
package handler

import (
	"log"
	"net/http"
	"fmt"
	"github.com/gin-gonic/gin"
	"backend/usecase"
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