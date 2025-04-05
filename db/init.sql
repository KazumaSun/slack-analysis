-- db/init.sql
-- ユーザーテーブル
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    user_key VARCHAR(255) UNIQUE NOT NULL,  -- SlackのユーザーID
    user_name VARCHAR(255) NOT NULL,        -- ユーザー名（表示名または実名）
    grade INTEGER NOT NULL DEFAULT 1,       -- ユーザーのグレード
    team_key INTEGER NOT NULL DEFAULT 1,    -- チームキー
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- チームテーブル（Slackチャンネルとの対応）
CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,                 -- team_keyとして使用
    channel_id VARCHAR(255) UNIQUE NOT NULL, -- SlackのチャンネルID
    channel_name VARCHAR(255) NOT NULL,      -- チャンネル名
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- インデックス（パフォーマンス向上のため）
CREATE INDEX IF NOT EXISTS idx_users_user_key ON users(user_key);
CREATE INDEX IF NOT EXISTS idx_users_team_key ON users(team_key);
CREATE INDEX IF NOT EXISTS idx_teams_channel_id ON teams(channel_id);

-- 元のactivity_logsテーブルを残す場合（必要に応じて）
CREATE TABLE IF NOT EXISTS activity_logs (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users(id),
  timestamp TIMESTAMP,
  status TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);