import { NextResponse } from "next/server"
import { history } from "../../_data"

export async function GET(req: Request, { params }: { params: { channel_id: string } }) {
  const { channel_id } = params
  const channelHistory = history.filter((item) => item.channel_id === channel_id)
  return NextResponse.json(channelHistory)
}