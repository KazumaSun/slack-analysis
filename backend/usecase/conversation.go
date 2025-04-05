package usecase

import (
	"fmt"
	"log"
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
func (u *ConversationUsecase) InitializeChannelConversations(channelID string) ([]slack.Message, error) {
	api := slack.New(u.slackTokenBot) // Slack APIの初期化
	// allConversations := []repository.SlackConversation{}	// 全ての会話を格納するスライス
	allMessages := []slack.Message{} // 全ての会話を格納するスライス

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

	// 動作確認用
	for i, message := range allMessages {
		fmt.Printf("メッセージ %d: %s\n", i+1, message.Text)
	}
	return allMessages, nil

	// フロントに返す処理とDBに格納する処理は並列処理でやると良い挑戦になるかも
	// // 取得した会話履歴をDBに保存
	// for _, message := range allMessages {
	// 	if err := u.repo.SaveConversation(message); err != nil {
	// 		return err
	// 	}
	// }
	// conversationsを作成して、それをSaveして、errっていうのもありか

	// return nil
}
