// backend/handler/slack_handler.go
package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"your-project/backend/usecase"
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
		log.Printf("Failed to initialize users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Users initialized successfully",
	})
}

// GetAllUsersHandler はユーザー情報取得APIのハンドラー
func (h *SlackHandler) GetAllUsersHandler(c *gin.Context) {
	users, err := h.slackUsecase.GetAllUsers()
	if err != nil {
		log.Printf("Failed to get users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}