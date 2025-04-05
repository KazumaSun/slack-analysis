'use client'

import { Box, Typography, List, ListItem, ListItemText, FormControl, InputLabel, Select, MenuItem, SelectChangeEvent, Divider, Stack, Button } from "@mui/material"
import { useState, useEffect } from "react"
import { Channel, History, User } from "@/type"
import { API_BASE_URL } from "@/constants"

export default function UsersPage() {
  const [users, setUsers] = useState<User[]>([])
  const [channels, setChannels] = useState<Channel[]>([])
  const [selectedChannel, setSelectedChannel] = useState<string>("")
  const [channelUsers, setChannelUsers] = useState<string[]>([])

  // ユーザー情報とチャンネル情報を取得
  useEffect(() => {
    const fetchInitialData = async () => {
      try {
        // ユーザー情報を取得
        const usersResponse = await fetch(`${API_BASE_URL}/users`)
        if (!usersResponse.ok) throw new Error("Failed to fetch users")
        const usersData = await usersResponse.json()
        setUsers(usersData.users)

        // チャンネル情報を取得
        const channelsResponse = await fetch(`${API_BASE_URL}/channels`)
        if (!channelsResponse.ok) throw new Error("Failed to fetch channels")
        const channelsData = await channelsResponse.json()
        const channelsArray = Array.isArray(channelsData.channels) ? channelsData.channels : []
        setChannels(channelsArray)

        // デフォルトで最初のチャンネルを選択
        if (channelsArray.length > 0) {
          setSelectedChannel(channelsArray[0].channel_id)
        }
      } catch (error) {
        console.error("データの取得に失敗しました:", error)
      }
    }

    fetchInitialData()
  }, [])

  // 選択されたチャンネルの履歴を取得
  useEffect(() => {
    if (!selectedChannel) return

    const fetchChannelHistory = async () => {
      try {
        const historyResponse = await fetch(`${API_BASE_URL}/history/${selectedChannel}`)
        if (!historyResponse.ok) throw new Error(`Failed to fetch history for channel ${selectedChannel}`)
        const channelHistory = await historyResponse.json()

        // チャンネルのメンバー一覧を生成
        const usersInChannel = channelHistory.messages.map((historyItem: History) => {
          const user = users.find((u) => u.user_key === historyItem.user_id)
          return user ? user.user_name : "不明なユーザー"
        }).filter((userName: string | null): userName is string => userName !== null)

        setChannelUsers(Array.from(new Set(usersInChannel))) // 重複を排除
      } catch (error) {
        console.error("履歴の取得に失敗しました:", error)
      }
    }

    fetchChannelHistory()
  }, [selectedChannel, users])

  // チャンネル選択時の処理
  const handleChannelChange = (event: SelectChangeEvent<string>) => {
    setSelectedChannel(event.target.value as string)
  }

  // ユーザー初期化
  const initializeUsers = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/users/init`, { method: "POST" })
      if (!response.ok) throw new Error("Failed to initialize users")
      alert("ユーザーの初期化が成功しました")
      window.location.reload() // ページを再読み込み
    } catch (error) {
      console.error("ユーザーの初期化に失敗しました:", error)
      alert("ユーザーの初期化に失敗しました")
    }
  }

  // チャンネル初期化
  const initializeChannels = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/channels/init`, { method: "POST" })
      if (!response.ok) throw new Error("Failed to initialize channels")
      alert("チャンネルの初期化が成功しました")
      window.location.reload() // ページを再読み込み
    } catch (error) {
      console.error("チャンネルの初期化に失敗しました:", error)
      alert("チャンネルの初期化に失敗しました")
    }
  }

  return (
    <Box p={4}>
      <Typography variant="h4" mb={2}>
        チャンネル/ユーザー一覧
      </Typography>

      {/* 初期化ボタン */}
      <Stack direction="row" spacing={2} mb={4}>
        <Button variant="contained" color="primary" onClick={initializeUsers}>
          ユーザー初期化
        </Button>
        <Button variant="contained" color="secondary" onClick={initializeChannels}>
          チャンネル初期化
        </Button>
      </Stack>

      {/* チャンネル選択セレクタ */}
      <FormControl fullWidth margin="normal">
        <InputLabel>チャンネルを選択</InputLabel>
        <Select value={selectedChannel} onChange={handleChannelChange}>
          {channels.map((channel) => (
            <MenuItem key={channel.channel_id} value={channel.channel_id}>
              {channel.channel_name}
            </MenuItem>
          ))}
        </Select>
      </FormControl>

      {/* 選択されたチャンネルのメンバー一覧 */}
      <Box mt={4}>
        <Typography variant="h6" mb={2}>
          メンバー一覧
        </Typography>
        <List>
          {channelUsers.map((userName) => (
            <ListItem key={userName}>
              <ListItemText primary={userName} />
            </ListItem>
          ))}
          {channelUsers.length === 0 && (
            <Typography color="text.secondary">このチャンネルにはユーザーがいません</Typography>
          )}
        </List>
      </Box>

      <Divider sx={{ my: 4 }} />

      {/* チャンネル一覧 */}
      <Box mt={4}>
        <Typography variant="h6" mb={2}>
          チャンネル一覧
        </Typography>
        <List>
          {channels.map((channel) => (
            <ListItem key={channel.channel_id}>
              <ListItemText
                primary={channel.channel_name}
              />
            </ListItem>
          ))}
        </List>
      </Box>
    </Box>
  )
}