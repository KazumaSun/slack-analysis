'use client'

import { useEffect, useState } from "react"
import {
  Box,
  Typography,
  Paper,
  List,
  ListItem,
  ListItemText,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  OutlinedInput,
  Stack,
  Chip,
  ToggleButton,
  ToggleButtonGroup,
} from "@mui/material"
import dayjs from "dayjs"
import isBetween from "dayjs/plugin/isBetween"
import { Channel, History, User } from "@/type"

dayjs.extend(isBetween)
import { SelectChangeEvent } from "@mui/material/Select"
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts"

type ActivityScale = "day" | "week" | "month"

export default function DashboardPage() {
  const [users, setUsers] = useState<User[]>([])
  const [channels, setChannels] = useState<Channel[]>([])
  const [selectedChannel, setSelectedChannel] = useState<string>("")
  const [filteredHistory, setFilteredHistory] = useState<History[]>([])
  const [selectedUsers, setSelectedUsers] = useState<string[]>([])
  const [suggestedTimes, setSuggestedTimes] = useState<string[]>([])
  const [activityData, setActivityData] = useState<{ time: string; count: number }[]>([])
  const [scale, setScale] = useState<ActivityScale>("day")

  // ユーザー情報とチャンネル情報を取得
  useEffect(() => {
    const fetchInitialData = async () => {
      try {
        // ユーザー情報を取得
        const usersResponse = await fetch("http://localhost:8080/users")
        if (!usersResponse.ok) throw new Error("Failed to fetch users")
        const usersData = await usersResponse.json()
        setUsers(usersData.users)

        // チャンネル情報を取得
        const channelsResponse = await fetch("http://localhost:8080/channels")
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

  // 選択されたチャンネルの投稿履歴を取得
  useEffect(() => {
    if (!selectedChannel) return

    const fetchChannelHistory = async () => {
      try {
        const historyResponse = await fetch(`http://localhost:8080/history/${selectedChannel}`)
        if (!historyResponse.ok) throw new Error(`Failed to fetch history for channel ${selectedChannel}`)
        const channelHistory = await historyResponse.json()
        setFilteredHistory(channelHistory.messages)
      } catch (error) {
        console.error("履歴の取得に失敗しました:", error)
      }
    }

    fetchChannelHistory()
  }, [selectedChannel])

  // ユーザーIDをユーザー名に変換
  const getUserName = (userId: string) => {
    const user = users.find((user) => user.user_key === userId)
    return user ? user.user_name : "不明なユーザー"
  }

  // チャンネル変更時の処理
  const handleChannelChange = (event: SelectChangeEvent<string>) => {
    const channelId = event.target.value as string
    setSelectedChannel(channelId)
  }

  // 投稿履歴を基にチャンネル参加ユーザーを取得
  const getChannelUsers = () => {
    const userIds = new Set(filteredHistory.map((entry) => entry.user_id))
    return users.filter((user) => userIds.has(user.user_key))
  }

  // ユーザー選択時の処理
  const handleUserChange = (event: SelectChangeEvent<string[]>) => {
    const users = event.target.value as string[]
    setSelectedUsers(users)
  }

  // ミーティング時間候補を計算
  const calculateMeetingTimes = () => {
    const userActivity: Record<string, number> = {}

    // 選択されたユーザーの投稿履歴を集計
    filteredHistory.forEach((entry) => {
      if (selectedUsers.includes(entry.user_id)) {
        const time = dayjs(entry.timestamp).format("HH:00") // 時刻のみ
        userActivity[time] = (userActivity[time] || 0) + 1
      }
    })

    // 全ユーザーがアクティブな時間を抽出
    let commonTimes = Object.entries(userActivity)
      .filter(([, count]) => count === selectedUsers.length)
      .sort(([timeA], [timeB]) => (timeA > timeB ? 1 : -1))
      .slice(0, 5) // 上位5つを取得

    // 全ユーザーがアクティブな時間がない場合、最も活動している時間帯を抽出
    if (commonTimes.length === 0) {
      commonTimes = Object.entries(userActivity)
        .sort(([, countA], [, countB]) => countB - countA) // 活動が多い順にソート
        .slice(0, 5)
    }

    setSuggestedTimes(commonTimes.map(([time]) => time))
  }

  // アクティビティデータを生成
  const generateActivityData = () => {
    const now = dayjs()
    const activity: Record<string, number> = {}

    // 投稿履歴を基にアクティビティデータを集計
    filteredHistory.forEach((entry) => {
      const time =
        scale === "day"
          ? dayjs(entry.timestamp).format("YYYY/MM/DD HH:00")
          : scale === "week"
            ? dayjs(entry.timestamp).format("YYYY/MM/DD")
            : dayjs(entry.timestamp).format("YYYY/MM/DD")
      activity[time] = (activity[time] || 0) + 1
    })

    const sortedActivity: { time: string; count: number }[] = []

    if (scale === "day") {
      // 直近24時間（1時間単位）
      for (let i = 0; i < 24; i++) {
        const hour = now.subtract(23 - i, "hour").format("YYYY/MM/DD HH:00")
        sortedActivity.push({ time: hour, count: activity[hour] || 0 })
      }
    } else if (scale === "week") {
      // 直近7日間（1日単位）
      for (let i = 0; i < 7; i++) {
        const day = now.subtract(6 - i, "day").format("YYYY/MM/DD")
        sortedActivity.push({ time: day, count: activity[day] || 0 })
      }
    } else {
      // 直近30日間（5日単位）
      for (let i = 0; i < 6; i++) {
        const start = now.subtract(25 - i * 5, "day").format("YYYY/MM/DD")
        const end = now.subtract(20 - i * 5, "day").format("YYYY/MM/DD")
        const range = `${start}`
        const count = Object.keys(activity)
          .filter((key) => dayjs(key).isBetween(start, end, null, "[]"))
          .reduce((sum, key) => sum + activity[key], 0)
        sortedActivity.push({ time: range, count })
      }
    }

    setActivityData(sortedActivity)
  }

  // アクティビティデータの更新
  useEffect(() => {
    generateActivityData()
  }, [filteredHistory, scale])

  return (
    <Box p={4}>
      <Typography variant="h4" gutterBottom>
        ダッシュボード
      </Typography>

      <Stack spacing={3}>
        {/* チャンネル選択 */}
        <Box>
          <Paper sx={{ p: 2 }}>
            <Typography variant="h6" gutterBottom>
              チャンネルを選択
            </Typography>
            <FormControl fullWidth size="small">
              <InputLabel>チャンネル</InputLabel>
              <Select
                value={selectedChannel}
                onChange={handleChannelChange}
                input={<OutlinedInput label="チャンネル" />}
              >
                {channels.map((channel) => (
                  <MenuItem key={channel.channel_id} value={channel.channel_id}>
                    {channel.channel_name}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Paper>
        </Box>

        {/* ユーザー選択 */}
        <Box>
          <Paper sx={{ p: 2 }}>
            <Typography variant="h6" gutterBottom>
              ユーザーを選択
            </Typography>
            <FormControl fullWidth size="small">
              <InputLabel>ユーザー</InputLabel>
              <Select
                multiple
                value={selectedUsers}
                onChange={handleUserChange}
                input={<OutlinedInput label="ユーザー" />}
                renderValue={(selected) => (
                  <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                    {selected.map((value) => (
                      <Chip key={value} label={getUserName(value)} />
                    ))}
                  </Box>
                )}
              >
                {getChannelUsers().map((user) => (
                  <MenuItem key={user.user_key} value={user.user_key}>
                    {user.user_name}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
            <Box mt={2}>
              <Typography variant="h6" gutterBottom>
                ミーティング時間候補
              </Typography>
              <List>
                {suggestedTimes.map((time, index) => (
                  <ListItem key={index}>
                    <ListItemText primary={time} />
                  </ListItem>
                ))}
                {suggestedTimes.length === 0 && (
                  <Typography color="text.secondary">候補がありません</Typography>
                )}
              </List>
              <Box mt={2}>
                <Chip
                  label="時間候補を計算"
                  color="primary"
                  clickable
                  onClick={calculateMeetingTimes}
                />
              </Box>
            </Box>
          </Paper>
        </Box>

        {/* アクティビティグラフ */}
        <Box>
          <Paper sx={{ p: 2 }}>
            <Typography variant="h6" gutterBottom>
              アクティビティグラフ
            </Typography>
            <ToggleButtonGroup
              value={scale}
              exclusive
              onChange={(_, value) => value && setScale(value)}
              size="small"
              sx={{ mb: 2 }}
            >
              <ToggleButton value="day">直近24時間</ToggleButton>
              <ToggleButton value="week">直近7日間</ToggleButton>
              <ToggleButton value="month">直近30日間</ToggleButton>
            </ToggleButtonGroup>
            <Box sx={{ width: '100%', height: 400 }}>
              <ResponsiveContainer>
                <LineChart data={activityData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="time" />
                  <YAxis />
                  <Tooltip />
                  <Line type="monotone" dataKey="count" stroke="#1976d2" />
                </LineChart>
              </ResponsiveContainer>
            </Box>
          </Paper>
        </Box>

        {/* 投稿履歴 */}
        <Box>
          <Paper sx={{ p: 2 }}>
            <Typography variant="h6" gutterBottom>
              投稿履歴
            </Typography>
            <List sx={{ maxHeight: 300, overflow: 'auto' }}>
              {filteredHistory.slice(0, 10).map((entry, index) => (
                <ListItem key={index} disablePadding>
                  <ListItemText
                    primary={`${entry.text}`}
                    secondary={`ユーザー: ${getUserName(entry.user_id)}, 時刻: ${entry.timestamp}`}
                  />
                </ListItem>
              ))}
              {filteredHistory.length === 0 && (
                <Typography color="text.secondary">投稿履歴がありません</Typography>
              )}
            </List>
          </Paper>
        </Box>
      </Stack>
    </Box>
  )
}
