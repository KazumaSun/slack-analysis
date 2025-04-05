// backend/repository/models.go
package repository

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