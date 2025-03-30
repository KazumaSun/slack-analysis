// src/app/api/users/route.ts

import { NextRequest, NextResponse } from "next/server"
import { users } from "../_data"

export async function GET() {
  return NextResponse.json(users)
}

export async function DELETE(req: NextRequest) {
  const { searchParams } = new URL(req.url)
  const id = Number(searchParams.get("id"))
  const index = users.findIndex((u) => u.id === id)
  if (index !== -1) users.splice(index, 1)
  return NextResponse.json({ message: "削除完了" })
}

export async function PATCH(req: NextRequest) {
  const { searchParams } = new URL(req.url)
  const id = Number(searchParams.get("id"))
  const body = await req.json()

  const user = users.find((u) => u.id === id)
  if (!user) return NextResponse.json({ error: "ユーザーが見つかりません" }, { status: 404 })

  user.name = body.name
  user.grade = body.grade
  user.channels = body.channels

  return NextResponse.json({ message: "更新完了" })
}
