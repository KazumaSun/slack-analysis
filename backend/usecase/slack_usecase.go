// backend/usecase/slack_usecase.go
package usecase

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"backend/repository"
)

type SlackUsecase struct {
	repo       *repository.Repository
	slackTokenUser string
	slackTokenBot string
}

func NewSlackUsecase(repo *repository.Repository, slackTokenUser string, slackTokenBot string) *SlackUsecase {
	return &SlackUsecase{
		repo:       repo,
		slackTokenUser: slackTokenUser,
		slackTokenBot: slackTokenBot,
	}
}

// InitializeUsers はSlack APIからユーザーリストを取得し、DBに保存します
func (u *SlackUsecase) InitializeUsers() error {
	// Slack APIからユーザーリストを取得
	users, err := u.fetchSlackUsers()
	if err != nil {
		return fmt.Errorf("InitializeUsers: failed to fetch slack users: %w", err)
	}

	// ユーザーをDBに保存
	for _, slackUser := range users {
		// 表示名が空の場合は実名を使用
		userName := slackUser.Profile.DisplayName
		if userName == "" {
			userName = slackUser.Profile.RealName
		}
		
		user := repository.User{
			UserKey:  slackUser.ID,
			UserName: userName,
			Grade:    1,         // 初期値
			TeamKey:  1,         // 初期値
		}
		
		if err := u.repo.SaveUser(user); err != nil {
			return fmt.Errorf("InitializeUsers: failed to save user %s (%s): %w", userName, slackUser.ID, err)
		}
	}

	return nil
}

// InitializeChannels は Slack API からチャンネルリストを取得し、フィルタリングしてDBに保存します (新規追加)
func (u *SlackUsecase) InitializeChannels() error {
	// チャンネル一覧を取得
	channels, err := u.fetchSlackChannels()
	if err != nil {
		return fmt.Errorf("InitializeChannels: failed to fetch slack channels: %w", err)
	}

	// フィルタリングとDBへの保存
	for _, channel := range channels {
		// "develop" または "team" を含むチャンネルのみ保存
		if strings.Contains(channel.Name, "develop") || strings.Contains(channel.Name, "team") {
			team := repository.Team{
				// ID:          teamKey, // DB が SERIAL で自動生成するなら不要
				ChannelID:   channel.ID,
				ChannelName: channel.Name,
			}
			// SaveTeam メソッドも ID を引数に取らないように修正が必要かも
			if err := u.repo.SaveTeam(team); err != nil {
				return fmt.Errorf("InitializeChannels: failed to save team %s (%s): %w", channel.Name, channel.ID, err)
			}
			// teamKey++ // ID を連番で振る場合
		}
	}
	return nil
}

// GetAllUsers はDBからすべてのユーザー情報を取得します (変更なし)
func (u *SlackUsecase) GetAllUsers() ([]repository.User, error) {
	users, err := u.repo.GetAllUsers()
	if err != nil {
		// Usecase層でもエラーをラップするとトレースしやすい
		return nil, fmt.Errorf("GetAllUsers: failed to get users from repository: %w", err)
	}
	return users, nil
}

// GetAllChannels はDBからすべてのチーム（チャンネル）情報を取得します (新規追加)
func (u *SlackUsecase) GetAllChannels() ([]repository.Team, error) {
	teams, err := u.repo.GetAllTeams()
	if err != nil {
		return nil, fmt.Errorf("GetAllChannels: failed to get teams from repository: %w", err)
	}
	return teams, nil
}

// fetchSlackUsers はSlack APIからユーザーリストを取得します
func (u *SlackUsecase) fetchSlackUsers() ([]repository.SlackUser, error) {
	req, err := http.NewRequest("GET", "https://slack.com/api/users.list", nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Add("Authorization", "Bearer "+u.slackTokenUser)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var result struct {
		Ok    bool                   `json:"ok"`
		Error string                 `json:"error"`
		Users []repository.SlackUser `json:"members"`
	}
	
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	
	if !result.Ok {
		return nil, fmt.Errorf("Slack API error: %s", result.Error)
	}
	
	return result.Users, nil
}

// fetchSlackChannels はSlack APIからチャンネル一覧を取得します
func (u *SlackUsecase) fetchSlackChannels() ([]repository.SlackChannel, error) {
	req, err := http.NewRequest("GET", "https://slack.com/api/conversations.list", nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Add("Authorization", "Bearer "+u.slackTokenBot)
	
	// パブリックチャンネルのみ取得
	q := req.URL.Query()
	q.Add("types", "public_channel")
	req.URL.RawQuery = q.Encode()
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var result struct {
		Ok       bool                     `json:"ok"`
		Error    string                   `json:"error"`
		Channels []repository.SlackChannel `json:"channels"`
	}
	
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	
	if !result.Ok {
		return nil, fmt.Errorf("Slack API error: %s", result.Error)
	}
	
	return result.Channels, nil
}

// UpdateUser は指定されたIDのユーザー情報を更新します (新規追加)
func (u *SlackUsecase) UpdateUser(id int, user repository.User) error {
	// ここでビジネスロジック（バリデーションなど）を追加することも可能

	// RepositoryのUpdateUserメソッドを呼び出す
	err := u.repo.UpdateUser(id, user)
	if err != nil {
		// エラーをラップして返す
		return fmt.Errorf("failed to update user in repository (id: %d): %w", id, err)
	}
	return nil
}