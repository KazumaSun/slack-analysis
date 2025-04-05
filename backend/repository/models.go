// backend/repository/models.go
package repository

import (
	"time"
)

type User struct {
	ID       int    `json:"id" db:"id"`
	UserKey  string `json:"user_key" db:"user_key"`
	UserName string `json:"user_name" db:"user_name"`
	Grade    int    `json:"grade" db:"grade"`
	TeamKey  int    `json:"team_key" db:"team_key"`
}

type Team struct {
	ID          int    `json:"id" db:"id"`
	ChannelID   string `json:"channel_id" db:"channel_id"`
	ChannelName string `json:"channel_name" db:"channel_name"`
}

type ActivityLog struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	TeamID      int       `json:"team_id"`
	WorkspaceID int       `json:"workspace_id"`
	Message     string    `json:"message"`
	Timestamp   time.Time `json:"timestamp"`
}

type SlackUser struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Profile struct {
		DisplayName string `json:"display_name"`
		RealName    string `json:"real_name"`
	} `json:"profile"`
}

type SlackChannel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SlackConversation struct {
	Messages []struct {
		User      string `json:"user"`
		Text      string `json:"text"`
		Timestamp string `json:"ts"`
	} `json:"messages"`
	HasMore          bool `json:"has_more"`
	ResponseMetadata struct {
		NextCursor string `json:"next_cursor"`
	} `json:"response_metadata"`
}
