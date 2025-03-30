// backend/usecase/slack_usecase.go
package usecase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"your-project/backend/repository"
)

type SlackUsecase struct {
	repo       *repository.Repository
	slackToken string
}

func NewSlackUsecase(repo *repository.Repository, slackToken string) *SlackUsecase {
	return &SlackUsecase{
		repo:       repo,
		slackToken: slackToken,
	}
}

// InitializeUsers はSlack APIからユーザーリストを取得し、DBに保存します
func (u *SlackUsecase) InitializeUsers() error {
	// Slack APIからユーザーリストを取得
	users, err := u.fetchSlackUsers()
	if err != nil {
		return err
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
			return err
		}
	}

	// チャンネル一覧を取得
	channels, err := u.fetchSlackChannels()
	if err != nil {
		return err
	}

	// フィルタリングとDBへの保存
	teamKey := 1
	for _, channel := range channels {
		// "develop" または "team" を含むチャンネルのみ保存
		if strings.Contains(channel.Name, "develop") || strings.Contains(channel.Name, "team") {
			team := repository.Team{
				ID:          teamKey,
				ChannelID:   channel.ID,
				ChannelName: channel.Name,
			}
			
			if err := u.repo.SaveTeam(team); err != nil {
				return err
			}
			
			teamKey++
		}
	}

	return nil
}

// GetAllUsers はDBからすべてのユーザー情報を取得します
func (u *SlackUsecase) GetAllUsers() ([]repository.User, error) {
	return u.repo.GetAllUsers()
}

// fetchSlackUsers はSlack APIからユーザーリストを取得します
func (u *SlackUsecase) fetchSlackUsers() ([]repository.SlackUser, error) {
	req, err := http.NewRequest("GET", "https://slack.com/api/users.list", nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Add("Authorization", "Bearer "+u.slackToken)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
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
	
	req.Header.Add("Authorization", "Bearer "+u.slackToken)
	
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
	
	body, err := ioutil.ReadAll(resp.Body)
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