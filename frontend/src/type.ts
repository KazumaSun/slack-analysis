// ユーザーの型
export interface User {
  user_id: string; // ユーザID
  name: string; // ユーザ名（profileのdisplay_name）
}

// チャンネルの型
export interface Channel {
  channel_id: string; // チャンネルID
  name: string; // チャンネル名
}

// 投稿履歴の型
export interface History {
  user_id: string;
  workspace_id: string;
  channel_id: string;
  message?: string;
  ts: string;
}