// src/app/api/_data.ts
import { User, Channel, History } from '@/type';

// モック状態（開発用）
export const users: User[] = [
  { user_id: "U12345", name: "佐藤" },
  { user_id: "U67890", name: "鈴木" },
  { user_id: "U54321", name: "田中" },
  { user_id: "U98765", name: "山田" },
];

export const channels: Channel[] = [
  { channel_id: "C12345", name: "general" },
  { channel_id: "C67890", name: "random" },
  { channel_id: "C98765", name: "development" },
];

export const history: History[] = [
  { user_id: "U12345", workspace_id: "W12345", channel_id: "C12345", message: "Hello!", ts: "2025/04/05 12:00:00" },
  { user_id: "U67890", workspace_id: "W12345", channel_id: "C12345", message: "Good morning!", ts: "2025/04/05 08:00:00" },
  { user_id: "U54321", workspace_id: "W12345", channel_id: "C67890", message: "Random message", ts: "2025/04/04 15:20:00" },
  { user_id: "U98765", workspace_id: "W12345", channel_id: "C98765", message: "Development update", ts: "2025/04/03 10:00:00" },
  { user_id: "U12345", workspace_id: "W12345", channel_id: "C12345", message: "Meeting at 3 PM", ts: "2025/04/02 14:00:00" },
  { user_id: "U67890", workspace_id: "W12345", channel_id: "C67890", message: "Lunch break", ts: "2025/04/01 12:00:00" },
  { user_id: "U54321", workspace_id: "W12345", channel_id: "C98765", message: "Code review completed", ts: "2025/03/31 16:30:00" },
  { user_id: "U98765", workspace_id: "W12345", channel_id: "C12345", message: "Project deadline extended", ts: "2025/03/30 09:00:00" },
  { user_id: "U12345", workspace_id: "W12345", channel_id: "C67890", message: "Random thoughts", ts: "2025/03/29 18:45:00" },
  { user_id: "U67890", workspace_id: "W12345", channel_id: "C98765", message: "New feature deployed", ts: "2025/03/28 11:15:00" },
  { user_id: "U54321", workspace_id: "W12345", channel_id: "C12345", message: "Testing", ts: "2025/04/05 14:00:00" },
  { user_id: "U98765", workspace_id: "W12345", channel_id: "C12345", message: "Debugging", ts: "2025/04/05 16:00:00" },
]
