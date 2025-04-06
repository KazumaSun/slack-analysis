// backend/repository/repository.go
package repository

import (
	"database/sql"
	"fmt"
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
		INSERT INTO teams (channel_id, channel_name) -- id を INSERT 文から除外
		VALUES ($1, $2)                             -- 引数も $1, $2 に
		ON CONFLICT (channel_id) DO UPDATE          -- コンフリクトは channel_id でチェック (UNIQUE 制約が必要)
		SET channel_name = $2                       -- channel_name を更新
	`

	// Exec に渡す引数から team.ID を削除
	_, err := r.db.Exec(query, team.ChannelID, team.ChannelName)
	if err != nil {
		log.Printf("Failed to save team (channel_id: %s): %v", team.ChannelID, err) // ログに詳細追加
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

	// rows.Err() をチェックして、ループ中のエラーを確認 (重要)
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating team rows: %v", err)
		return nil, err
	}

	// データがない場合、空のスライスを返す（nil ではなく）
	if users == nil {
		return []User{}, nil // nil ではなく空スライスを返すのが一般的
	}
	
	return users, nil
}

// GetAllTeams はすべてのチーム情報を取得します (新規追加)
func (r *Repository) GetAllTeams() ([]Team, error) {
	query := `SELECT id, channel_id, channel_name FROM teams ORDER BY id ASC` // ORDER BY を追加すると良いかも

	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("Failed to get teams: %v", err)
		return nil, err // エラーを返す
	}
	defer rows.Close()

	var teams []Team
	for rows.Next() {
		var team Team
		// Scan するカラムの順番を SELECT 文に合わせる
		if err := rows.Scan(&team.ID, &team.ChannelID, &team.ChannelName); err != nil {
			log.Printf("Failed to scan team: %v", err)
			return nil, err // エラーを返す
		}
		teams = append(teams, team)
	}

	// rows.Err() をチェックして、ループ中のエラーを確認 (重要)
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating team rows: %v", err)
		return nil, err
	}


	// データがない場合、空のスライスを返す（nil ではなく）
	if teams == nil {
		return []Team{}, nil // nil ではなく空スライスを返すのが一般的
	}


	return teams, nil
}

// UpdateUser は指定されたIDのユーザー情報を更新します (新規追加)
func (r *Repository) UpdateUser(id int, user User) error {
	// user_key は通常更新しないことが多いが、リクエストに含まれるなら更新対象に入れる
	// もし user_key を更新したくない場合は SET 句から user_key = $2 を削除し、引数の順番も調整する
	query := `
		UPDATE users
		SET user_key = $2, user_name = $3, grade = $4, team_key = $5
		WHERE id = $1
	`

	// Exec を実行し、結果を取得
	result, err := r.db.Exec(query, id, user.UserKey, user.UserName, user.Grade, user.TeamKey)
	if err != nil {
		log.Printf("Failed to execute update user query (id: %d): %v", id, err)
		return fmt.Errorf("database error executing update query for id %d: %w", id, err) // エラーラップ
	}

	// 更新された行数を取得 (任意だが、更新が成功したか確認できる)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to get rows affected for update user (id: %d): %v", id, err)
		// RowsAffected のエラーは続行しても良い場合もあるが、ここではエラーとして返す
		return fmt.Errorf("database error getting rows affected for id %d: %w", id, err)
	}

	// 更新対象のIDが存在せず、更新が行われなかった場合
	if rowsAffected == 0 {
		log.Printf("No user found with id %d to update", id)
		// ここでエラーを返すか、成功として扱うかは要件による
		// 例: 存在しないIDの場合はエラーとする
		return fmt.Errorf("no user found with id %d", id) // エラーとして返す例
	}

	log.Printf("Successfully updated user with id %d", id)
	return nil
}