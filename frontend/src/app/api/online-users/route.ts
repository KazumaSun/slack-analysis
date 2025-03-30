// src/app/api/online-users/route.ts

import { NextResponse } from 'next/server'

const mockOnlineUsers = [
  { id: 1, name: "佐藤", status: "online" },
  { id: 2, name: "鈴木", status: "online" },
  { id: 3, name: "田中", status: "away" }
]

export async function GET() {
  return NextResponse.json(mockOnlineUsers)
}
