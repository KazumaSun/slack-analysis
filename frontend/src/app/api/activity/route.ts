// src/app/api/activity/route.ts

import { NextResponse } from 'next/server'

const mockActivityData = [
  { hour: "09:00", count: 2 },
  { hour: "10:00", count: 4 },
  { hour: "11:00", count: 6 },
  { hour: "12:00", count: 3 },
  { hour: "13:00", count: 5 },
  { hour: "14:00", count: 8 },
  { hour: "15:00", count: 7 },
  { hour: "16:00", count: 4 },
]

export async function GET() {
  return NextResponse.json(mockActivityData)
}
