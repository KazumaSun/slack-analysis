package usecase

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"backend/repository"

	"github.com/slack-go/slack"
)

// ConversationUsecase は会話に関するユースケースを提供します
type ConversationUsecase struct {
	repo          *repository.Repository
	slackTokenBot string
}

// 初期化関数
func NewConversationUsecase(repo *repository.Repository, slackTokenBot string) *ConversationUsecase {
	return &ConversationUsecase{
		repo:          repo,
		slackTokenBot: slackTokenBot,
	}
}

// InitializeChannelConversations は指定したチャンネルの会話履歴の初期化を行います
func (u *ConversationUsecase) InitializeChannelConversations(channelID string) ([]repository.SlackConversation, error) {
	api := slack.New(u.slackTokenBot)
	allMessages := []slack.Message{}
	allConversations := []repository.SlackConversation{}

	// チャンネルにボットを参加させる
	_, _, _, err := api.JoinConversation(channelID)
	if err != nil {
		if strings.Contains(err.Error(), "missing_scope") {
			log.Printf("スコープが不足しています: %v", err)
			return nil, fmt.Errorf("missing required scope: %w", err)
		}
		log.Printf("チャンネルへの参加に失敗しました: %v", err)
		return nil, fmt.Errorf("failed to join channel: %w", err)
	}

	historyParams := slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Limit:     1000,
	}

	for {
		history, err := api.GetConversationHistory(&historyParams)
		if err != nil {
			log.Printf("会話履歴の取得に失敗しました: %v", err)
			return nil, fmt.Errorf("failed to fetch conversation history: %w", err)
		}

		allMessages = append(allMessages, history.Messages...)
		if history.ResponseMetaData.NextCursor == "" {
			break
		}

		historyParams.Cursor = history.ResponseMetaData.NextCursor
		time.Sleep(1200 * time.Millisecond)
	}

	for _, message := range allMessages {
		ts, err := FormatSlackTimestamp(message.Timestamp)
		if err != nil {
			log.Printf("タイムスタンプのフォーマットに失敗しました: %v", err)
			continue
		}
		allConversations = append(allConversations, repository.SlackConversation{
			ChannelID:   channelID,
			UserID:      message.User,
			WorkspaceID: message.Team,
			Text:        message.Text,
			Timestamp:   ts,
		})
	}

	return allConversations, nil
}

// FormatSlackTimestamp は Slack API から取得したタイムスタンプ文字列
// (例: "1601055549.000100") を受け取り、
// "YYYY/MM/DD hh:mm:ss" (24時間表記) の文字列にフォーマットします
// エラーが発生した場合は、0 とエラーを返します。
func FormatSlackTimestamp(slackTs string) (string, error) {
	if slackTs == "" {
		return "0", fmt.Errorf("input timestamp string is empty")
	}

	// "." で文字列を分割
	parts := strings.Split(slackTs, ".")
	if len(parts) == 0 {
		// 通常はありえないが念のため
		return "0", fmt.Errorf("invalid timestamp format: splitting resulted in zero parts for '%s'", slackTs)
	}

	// 最初の部分（秒の部分）を取得
	secondsStr := parts[0]

	// 秒の部分の文字列を int64 に変換
	// 第2引数は基数 (10進数なので10)
	// 第3引数はビット数 (64ビット整数なので64)
	unixTimeSeconds, err := strconv.ParseInt(secondsStr, 10, 64)
	if err != nil {
		// 変換に失敗した場合（数字以外の文字が含まれているなど）
		return "0", fmt.Errorf("failed to parse timestamp '%s' to int64: %v", secondsStr, err)
	}

	// int64のUnixタイムスタンプ(秒)を time.Time 型に変換
	// 第2引数はナノ秒部分。秒単位のタイムスタンプなら0でOK
	t := time.Unix(unixTimeSeconds, 0)

	// time.Time 型の値を指定したレイアウト文字列でフォーマット
	// "YYYY/MM/DD hh:mm:ss" (24時間表記) に対応するGoのレイアウトは "2006/01/02 15:04:05"
	formattedTime := t.Format("2006/01/02 15:04:05")

	// 正常に変換できたらUnixタイムスタンプ(秒)を返す
	// return unixTimeSeconds, nil
	return formattedTime, nil
}
