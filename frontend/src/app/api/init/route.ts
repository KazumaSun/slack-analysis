// src/app/api/init/route.ts

import { NextResponse } from "next/server"
import { resetMockData } from "../_data"

export async function POST() {
  resetMockData()
  return NextResponse.json({ message: "初期化完了" })
}
