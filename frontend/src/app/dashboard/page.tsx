'use client'

import { useEffect, useState } from "react"
import {
  Box,
  Typography,
  useTheme,
  useMediaQuery,
  Paper,
  List,
  ListItem,
  ListItemText,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  OutlinedInput,
  Chip,
  Stack,
  ToggleButton,
  ToggleButtonGroup,
} from "@mui/material"
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts"
import dayjs from "dayjs"

const allMockUsers = [
  { id: 1, name: "佐藤" },
  { id: 2, name: "鈴木" },
  { id: 3, name: "田中" },
  { id: 4, name: "山本" },
  { id: 5, name: "中村" },
]

const mockAvailability: Record<string, string[]> = {
  "佐藤": ["月-10:00", "水-14:00", "金-15:00"],
  "鈴木": ["水-14:00", "木-10:00", "金-15:00"],
  "田中": ["水-14:00", "金-15:00"],
  "山本": ["月-10:00", "火-13:00"],
  "中村": ["金-15:00", "水-14:00"],
}

type ActivityScale = "day" | "week" | "month"

type User = {
  id: number
  name: string
  status: string
}

export default function DashboardPage() {
  const theme = useTheme()
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'))

  const [onlineData, setOnlineData] = useState<{ time: string; count: number }[]>([])
  const [users, setUsers] = useState<User[]>([])
  const [selectedUserNames, setSelectedUserNames] = useState<string[]>([])
  const [suggestions, setSuggestions] = useState<string[]>([])
  const [scale, setScale] = useState<ActivityScale>("day")

  const generateActivityData = (scale: ActivityScale) => {
    const now = dayjs()
    if (scale === "day") {
      return [...Array(24)].map((_, i) => {
        const hour = now.subtract(23 - i, 'hour')
        return {
          time: hour.format("HH:00"),
          count: Math.floor(Math.random() * 10 + 1),
        }
      })
    } else if (scale === "week") {
      return [...Array(7)].map((_, i) => {
        const day = now.subtract(6 - i, 'day')
        return {
          time: day.format("M/D"),
          count: Math.floor(Math.random() * 30 + 10),
        }
      })
    } else {
      return [...Array(6)].map((_, i) => {
        const day = now.subtract(25 - i * 5, 'day')
        return {
          time: day.format("M/D"),
          count: Math.floor(Math.random() * 100 + 50),
        }
      })
    }
  }

  const suggestBestMTTime = () => {
    const availabilityCount: Record<string, number> = {}
    selectedUserNames.forEach((name) => {
      mockAvailability[name]?.forEach((slot) => {
        availabilityCount[slot] = (availabilityCount[slot] || 0) + 1
      })
    })

    const sorted = Object.entries(availabilityCount)
      .sort((a, b) => b[1] - a[1])
      .slice(0, 3)

    if (sorted.length > 0) {
      setSuggestions(sorted.map(([slot, count]) => `${slot}（${count}人が可能）`))
    } else {
      setSuggestions(["共通で空いている時間が見つかりません"])
    }
  }

  useEffect(() => {
    setOnlineData(generateActivityData(scale))
    fetch("/api/online-users").then((res) => res.json()).then(setUsers)
  }, [scale])

  return (
    <Box p={isMobile ? 2 : 4}>
      <Typography variant={isMobile ? "h5" : "h4"} gutterBottom>
        ダッシュボード
      </Typography>

      <Stack spacing={3}>
        {/* グラフ */}
        <Box>
          <Paper sx={{ p: 2 }}>
            <Typography variant="h6" gutterBottom>
              アクティビティ推移
            </Typography>
            <ToggleButtonGroup
              value={scale}
              exclusive
              onChange={(_, val) => val && setScale(val)}
              size="small"
              sx={{ mb: 2 }}
            >
              <ToggleButton value="day">直近24時間</ToggleButton>
              <ToggleButton value="week">直近7日</ToggleButton>
              <ToggleButton value="month">直近30日</ToggleButton>
            </ToggleButtonGroup>
            <Box sx={{ width: '100%', height: isMobile ? 300 : 400 }}>
              <ResponsiveContainer>
                <LineChart data={onlineData}>
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

        {/* オンラインユーザー & MT時間提案 */}
        <Stack direction={isMobile ? 'column' : 'row'} spacing={3}>
          <Box flex={1}>
            <Paper sx={{ p: 2 }}>
              <Typography variant="h6" gutterBottom>
                現在オンラインのユーザー
              </Typography>
              <List>
                {users.map((user) => (
                  <ListItem key={user.id} disablePadding>
                    <ListItemText primary={user.name} secondary={user.status} />
                  </ListItem>
                ))}
                {users.length === 0 && (
                  <Typography color="text.secondary">ユーザーが見つかりません</Typography>
                )}
              </List>
            </Paper>
          </Box>

          <Box flex={1}>
            <Paper sx={{ p: 2 }}>
              <Typography variant="h6" gutterBottom>
                MT時間を提案
              </Typography>
              <FormControl fullWidth size="small" sx={{ mb: 2 }}>
                <InputLabel>ユーザーを選択</InputLabel>
                <Select
                  multiple
                  value={selectedUserNames}
                  onChange={(e) => setSelectedUserNames(e.target.value as string[])}
                  input={<OutlinedInput label="ユーザーを選択" />}
                  renderValue={(selected) => (
                    <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                      {selected.map((value) => (
                        <Chip key={value} label={value} />
                      ))}
                    </Box>
                  )}
                >
                  {allMockUsers.map((user) => (
                    <MenuItem key={user.name} value={user.name}>{user.name}</MenuItem>
                  ))}
                </Select>
              </FormControl>
              <Box>
                <Typography variant="body1" gutterBottom>
                  最も多くのメンバーが参加できそうな時間候補:
                </Typography>
                {suggestions.map((s, idx) => (
                  <Typography key={idx} variant="h6" color="primary">
                    {s}
                  </Typography>
                ))}
              </Box>
              <Box mt={2}>
                <Chip
                  label="提案を更新"
                  color="primary"
                  clickable
                  onClick={suggestBestMTTime}
                />
              </Box>
            </Paper>
          </Box>
        </Stack>

        {/* その他カード */}
        <Box>
          <Paper sx={{ p: 2 }}>
            <Typography variant="h6" gutterBottom>
              アクティビティ減少アラート（準備中）
            </Typography>
            <Typography color="text.secondary">
              投稿数の減少や離脱傾向を検知します。
            </Typography>
          </Paper>
        </Box>
      </Stack>
    </Box>
  )
}
