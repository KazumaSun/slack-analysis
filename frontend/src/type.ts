// ユーザーの型
export interface User {
  user_key: string; // SlackのユーザーID
  user_name: string; // ユーザー名（表示名または実名）
  grade: number; // ユーザーのグレード
  team_key: number; // チームキー
}

// チャンネルの型
export interface Channel {
  channel_id: string; // SlackのチャンネルID
  channel_name: string; // チャンネル名
}

// 投稿履歴の型
export interface History {
  user_id: string;
  workspace_id: string;
  channel_id: string;
  message?: string;
  ts: string;
}