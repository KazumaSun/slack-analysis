// src/app/api/channels/route.ts

import { NextResponse } from "next/server"
import { channels } from "../_data"

export async function GET() {
  return NextResponse.json(channels)
}
