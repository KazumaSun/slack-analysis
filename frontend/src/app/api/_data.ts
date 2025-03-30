// src/app/api/_data.ts

export type User = {
  id: number
  name: string
  grade: number
  channels: string[]
}

export type Channel = {
  name: string
  users: string[]
}

// モック状態（開発用）
export let users: User[] = []
export let channels: Channel[] = []

export const resetMockData = () => {
  users = [
    { id: 1, name: "佐藤", grade: 3, channels: ["general", "dev"] },
    { id: 2, name: "鈴木", grade: 2, channels: ["general"] },
    { id: 3, name: "田中", grade: 4, channels: ["dev"] },
  ]

  channels = [
    { name: "general", users: ["佐藤", "鈴木"] },
    { name: "dev", users: ["佐藤", "田中"] },
  ]
}
