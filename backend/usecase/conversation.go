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
	api := slack.New(u.slackTokenBot)                    // Slack APIの初期化
	allMessages := []slack.Message{}                     // 全ての会話を格納するスライス
	allConversations := []repository.SlackConversation{} // 全ての会話を格納するスライス

	// チャンネルの会話履歴を取得するパラメータ
	historyParams := slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Limit:     1000,
	}

	// ページネーションのためのループ
	for {
		// APIを呼び出して履歴を取得
		history, err := api.GetConversationHistory(&historyParams)
		if err != nil {
			log.Fatalf("会話履歴の取得に失敗しました: %v", err)
		}

		// 取得したメッセージを allMessages に追加
		allMessages = append(allMessages, history.Messages...)
		fmt.Printf("%d 件のメッセージを取得しました (合計: %d 件)\n", len(history.Messages), len(allMessages))

		// 次のカーソルがあるか確認
		if history.ResponseMetaData.NextCursor == "" {
			fmt.Println("全てのメッセージを取得しました。")
			break // ループを抜ける
		}

		// 次のページのカーソルをパラメータにセット
		historyParams.Cursor = history.ResponseMetaData.NextCursor
		fmt.Printf("次のカーソル: %s\n", historyParams.Cursor)

		// --- レート制限への配慮 (重要！) ---
		// Slack APIにはレート制限があります。連続で呼び出しすぎるとエラーになります。
		// Tier 3 (conversations.history) は比較的要求数が多い (約50+/min) ですが、
		// 念のため、ループの間に短い待機時間を設けることを推奨します。
		time.Sleep(1200 * time.Millisecond) // 例: 1.2秒待機
	}

	// 取得したメッセージを表示
	for i, message := range allMessages {
		fmt.Printf("メッセージ %d: %s\n", i+1, message.Text)    // メッセージの内容を表示
		ts, err := FormatSlackTimestamp(message.Timestamp) // タイムスタンプをフォーマット
		if err != nil {
			log.Printf("タイムスタンプのフォーマットに失敗しました: %v", err)
			continue
		}
		// メッセージを全ての会話のスライスに追加
		allConversations = append(allConversations, repository.SlackConversation{
			ChannelID:   channelID,
			UserID:      message.User,
			WorkspaceID: message.Team,
			Text:        message.Text,
			Timestamp:   ts,
		})
	}

	// 取得した会話履歴をDBに保存

	// フロントに返す処理とDBに格納する処理は並列処理でやると良い挑戦になるかも
	return allConversations, nil

	// // 取得した会話履歴をDBに保存
	// for _, message := range allMessages {
	// 	if err := u.repo.SaveConversation(message); err != nil {
	// 		return err
	// 	}
	// }
	// conversationsを作成して、それをSaveして、errっていうのもありか
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
