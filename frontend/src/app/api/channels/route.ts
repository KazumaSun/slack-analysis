// src/app/api/channels/route.ts

import { NextRequest, NextResponse } from "next/server"
import { channels } from "../_data"

export async function GET() {
  return NextResponse.json(channels)
}

export async function PATCH(req: NextRequest) {
  const { name, users: newUsers } = await req.json()

  const channel = channels.find((c) => c.name === name)
  if (!channel) return NextResponse.json({ error: "チャンネルが見つかりません" }, { status: 404 })

  channel.users = newUsers
  return NextResponse.json({ message: "チャンネル更新完了" })
}
