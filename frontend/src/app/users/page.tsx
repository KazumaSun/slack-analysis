'use client'

import { Box, Typography, List, ListItem, ListItemText, Collapse, IconButton, Button, Stack } from "@mui/material"
import { ExpandLess, ExpandMore } from "@mui/icons-material"
import { useState, useEffect } from "react"
import { Channel, History, User } from "@/type"
import { API_BASE_URL } from "@/constants"

export default function UsersPage() {
  const [openChannels, setOpenChannels] = useState<{ [key: string]: boolean }>({})
  const [users, setUsers] = useState<User[]>([])
  const [channels, setChannels] = useState<Channel[]>([])
  const [history, setHistory] = useState<{ [key: string]: History[] }>({})

  // チャンネルの展開状態をトグル
  const toggleChannel = (channelId: string) => {
    setOpenChannels((prev) => ({
      ...prev,
      [channelId]: !prev[channelId],
    }))
  }

  // データを取得する
  useEffect(() => {
    const fetchData = async () => {
      try {
        // ユーザーリストを取得
        const usersResponse = await fetch(`${API_BASE_URL}/api/users`)
        const usersData = await usersResponse.json()
        setUsers(usersData)

        // チャンネルリストを取得
        const channelsResponse = await fetch(`${API_BASE_URL}/channels`)
        const channelsData = await channelsResponse.json()
        setChannels(channelsData)

        // 各チャンネルの投稿履歴を取得
        const historyData: { [key: string]: History[] } = {}
        for (const channel of channelsData) {
          const historyResponse = await fetch(`${API_BASE_URL}/history/${channel.channel_id}`)
          const channelHistory = await historyResponse.json()
          historyData[channel.channel_id] = channelHistory
        }
        setHistory(historyData)
      } catch (error) {
        console.error("データの取得に失敗しました:", error)
      }
    }

    fetchData()
  }, [])

  // チャンネルごとのユーザー一覧を生成
  const channelUsers = channels.map((channel) => {
    const usersInChannel = history[channel.channel_id]?.map((historyItem) => {
      const user = users.find((u) => u.user_key === historyItem.user_id)
      return user ? user.user_name : null
    }).filter((userName): userName is string => userName !== null)
    return { ...channel, users: Array.from(new Set(usersInChannel)) }
  })

  // ユーザーリスト初期化
  const handleUserInit = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/api/users/init`, { method: "POST" })
      if (response.ok) {
        alert("ユーザー情報が初期化されました")
        window.location.reload();
      } else {
        console.error("ユーザー情報初期化に失敗しました:", await response.text())
      }
    } catch (error) {
      console.error("ユーザー情報初期化に失敗しました:", error)
    }
  }

  // チャンネルリスト初期化
  const handleChannelInit = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/api/channels/init`, { method: "POST" })
      if (response.ok) {
        alert("チャンネル情報が初期化されました")
        window.location.reload();
      } else {
        console.error("チャンネル情報初期化に失敗しました:", await response.text())
      }
    } catch (error) {
      console.error("チャンネル情報初期化に失敗しました:", error)
    }
  }

  return (
    <Box p={4}>
      <Typography variant="h4" mb={2}>
        チャンネル/ユーザー管理
      </Typography>

      {/* 初期化ボタン */}
      <Stack direction="row" spacing={2} mb={4}>
        <Button variant="contained" color="primary" onClick={handleUserInit}>
          ユーザー初期化
        </Button>
        <Button variant="contained" color="secondary" onClick={handleChannelInit}>
          チャンネル初期化
        </Button>
      </Stack>

      {/* チャンネル一覧 */}
      <Box>
        <Typography variant="h6" mb={2}>
          チャンネル一覧
        </Typography>
        {channelUsers.map((channel) => (
          <Box key={channel.channel_id} mb={4}>
            <Box display="flex" alignItems="center" justifyContent="space-between">
              <Typography fontWeight="bold" mb={1}>
                #{channel.name}
              </Typography>
              <IconButton onClick={() => toggleChannel(channel.channel_id)}>
                {openChannels[channel.channel_id] ? <ExpandLess /> : <ExpandMore />}
              </IconButton>
            </Box>
            <Collapse in={openChannels[channel.channel_id]} timeout="auto" unmountOnExit>
              <List>
                {channel.users.map((userName) => (
                  <ListItem key={userName}>
                    <ListItemText primary={userName} />
                  </ListItem>
                ))}
                {channel.users.length === 0 && (
                  <Typography color="text.secondary">
                    このチャンネルにはユーザーがいません
                  </Typography>
                )}
              </List>
            </Collapse>
          </Box>
        ))}
      </Box>
    </Box>
  )
}