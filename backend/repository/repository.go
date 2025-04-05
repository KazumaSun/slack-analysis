// backend/repository/repository.go
package repository

import (
	"database/sql"
	"log"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// SaveUser はユーザー情報をDBに保存します
func (r *Repository) SaveUser(user User) error {
	query := `
		INSERT INTO users (user_key, user_name, grade, team_key)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_key) DO UPDATE
		SET user_name = $2, grade = $3, team_key = $4
	`
	
	_, err := r.db.Exec(query, user.UserKey, user.UserName, user.Grade, user.TeamKey)
	if err != nil {
		log.Printf("Failed to save user: %v", err)
		return err
	}
	
	return nil
}

// SaveTeam はチームとチャンネルの対応をDBに保存します
func (r *Repository) SaveTeam(team Team) error {
	query := `
		INSERT INTO teams (id, channel_id, channel_name)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE
		SET channel_id = $2, channel_name = $3
	`
	
	_, err := r.db.Exec(query, team.ID, team.ChannelID, team.ChannelName)
	if err != nil {
		log.Printf("Failed to save team: %v", err)
		return err
	}
	
	return nil
}

// GetAllUsers はすべてのユーザー情報を取得します
func (r *Repository) GetAllUsers() ([]User, error) {
	query := `SELECT id, user_key, user_name, grade, team_key FROM users`
	
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("Failed to get users: %v", err)
		return nil, err
	}
	defer rows.Close()
	
	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.UserKey, &user.UserName, &user.Grade, &user.TeamKey); err != nil {
			log.Printf("Failed to scan user: %v", err)
			return nil, err
		}
		users = append(users, user)
	}
	
	return users, nil
}