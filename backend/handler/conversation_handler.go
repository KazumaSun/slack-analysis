// backend/handler/slack_handler.go
package handler

import (
	"log"
	"net/http"

	"backend/usecase"

	"github.com/gin-gonic/gin"
)

type ConversationHandler struct {
	conversationUsecase *usecase.ConversationUsecase
}

func NewConversationHandler(conversationUsecase *usecase.ConversationUsecase) *ConversationHandler {
	return &ConversationHandler{
		conversationUsecase: conversationUsecase,
	}
}

func (h *ConversationHandler) InitializeChannelConversationsHandler(c *gin.Context) {

	// コンテキストからチャンネルIDを取得
	channelID := c.Param("channel_id") // URLパラメータから取得
	if channelID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "channel_id is required",
		})
		return
	}
	// 動作確認用
	log.Printf("channelID: %s", channelID)

	// チャンネルの会話履歴を取得
	allMessages, err := h.conversationUsecase.InitializeChannelConversations(channelID)
	if err != nil {
		log.Printf("Failed to initialize channel conversations: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 成功した場合のレスポンス
	c.JSON(http.StatusOK, gin.H{
		"message":    "Users initialized successfully",
		"channel_id": channelID,
		"status":     "success",
		"messages":   allMessages,
	})
}
