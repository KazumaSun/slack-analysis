'use client'

import { useEffect, useState } from "react"
import {
  TextField,
  Select,
  MenuItem,
  InputLabel,
  FormControl,
  Button,
  Box,
  OutlinedInput,
  Chip,
  Typography,
} from "@mui/material"

type User = {
  id: number
  name: string
  grade: number
  channels: string[]
}

type Channel = {
  name: string
  users: string[]
}

export default function UsersPage() {
  const [users, setUsers] = useState<User[]>([])
  const [channels, setChannels] = useState<Channel[]>([])
  const [editId, setEditId] = useState<number | null>(null)
  const [editedUser, setEditedUser] = useState<Partial<User>>({})

  const fetchAll = async () => {
    const [u, c] = await Promise.all([
      fetch("/api/users").then(res => res.json()),
      fetch("/api/channels").then(res => res.json()),
    ])
    setUsers(u)
    setChannels(c)
  }

  const resetData = async () => {
    await fetch("/api/init", { method: "POST" })
    fetchAll()
  }

  const deleteUser = async (id: number) => {
    await fetch(`/api/users?id=${id}`, { method: "DELETE" })
    fetchAll()
  }

  const saveUser = async (id: number) => {
    await fetch(`/api/users?id=${id}`, {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(editedUser),
    })
    setEditId(null)
    setEditedUser({})
    fetchAll()
  }

  const updateChannelUsers = async (name: string, userList: string[]) => {
    await fetch("/api/channels", {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name, users: userList }),
    })
    fetchAll()
  }

  useEffect(() => {
    fetchAll()
  }, [])

  return (
    <Box p={4}>
      <Typography variant="h4" mb={2}>ユーザー / チャンネル管理</Typography>

      <Button variant="outlined" onClick={resetData} sx={{ mb: 4 }}>
        データ初期化
      </Button>

      <Box mb={6}>
        <Typography variant="h6" mb={2}>ユーザー一覧</Typography>
        {users.map(user => (
          <Box key={user.id} mb={2}>
            {editId === user.id ? (
              <Box display="flex" alignItems="center" gap={2}>
                <TextField
                  label="名前"
                  defaultValue={user.name}
                  onChange={(e) => setEditedUser({ ...editedUser, name: e.target.value })}
                  size="small"
                />

                <FormControl size="small" sx={{ minWidth: 100 }}>
                  <InputLabel>学年</InputLabel>
                  <Select
                    defaultValue={user.grade}
                    label="学年"
                    onChange={(e) =>
                      setEditedUser({ ...editedUser, grade: Number(e.target.value) })
                    }
                  >
                    {[1, 2, 3, 4].map((year) => (
                      <MenuItem key={year} value={year}>{year}年</MenuItem>
                    ))}
                  </Select>
                </FormControl>

                <FormControl size="small" sx={{ minWidth: 200 }}>
                  <InputLabel>チャンネル</InputLabel>
                  <Select
                    multiple
                    defaultValue={user.channels}
                    onChange={(e) =>
                      setEditedUser({
                        ...editedUser,
                        channels: e.target.value as string[],
                      })
                    }
                    input={<OutlinedInput label="チャンネル" />}
                    renderValue={(selected) => (
                      <Box sx={{ display: "flex", flexWrap: "wrap", gap: 0.5 }}>
                        {(selected as string[]).map((value) => (
                          <Chip key={value} label={value} />
                        ))}
                      </Box>
                    )}
                  >
                    {channels.map((c) => (
                      <MenuItem key={c.name} value={c.name}>{c.name}</MenuItem>
                    ))}
                  </Select>
                </FormControl>

                <Button variant="contained" onClick={() => saveUser(user.id)}>保存</Button>
              </Box>
            ) : (
              <Box display="flex" alignItems="center" gap={2}>
                <Typography>{user.name}（{user.grade}年） - [{user.channels.join(", ")}]</Typography>
                <Button size="small" onClick={() => deleteUser(user.id)}>削除</Button>
                <Button size="small" onClick={() => {
                  setEditId(user.id)
                  setEditedUser(user)
                }}>編集</Button>
              </Box>
            )}
          </Box>
        ))}
      </Box>

      <Box>
        <Typography variant="h6" mb={2}>チャンネル一覧</Typography>
        {channels.map((channel) => (
          <Box key={channel.name} mb={2}>
            <Typography fontWeight="bold">#{channel.name}</Typography>
            <FormControl size="small" sx={{ minWidth: 300, mt: 1 }}>
              <InputLabel>参加ユーザー</InputLabel>
              <Select
                multiple
                defaultValue={channel.users}
                onChange={(e) =>
                  updateChannelUsers(
                    channel.name,
                    e.target.value as string[]
                  )
                }
                input={<OutlinedInput label="参加ユーザー" />}
                renderValue={(selected) => (
                  <Box sx={{ display: "flex", flexWrap: "wrap", gap: 0.5 }}>
                    {(selected as string[]).map((value) => (
                      <Chip key={value} label={value} />
                    ))}
                  </Box>
                )}
              >
                {users.map((u) => (
                  <MenuItem key={u.name} value={u.name}>{u.name}</MenuItem>
                ))}
              </Select>
            </FormControl>
          </Box>
        ))}
      </Box>
    </Box>
  )
}
