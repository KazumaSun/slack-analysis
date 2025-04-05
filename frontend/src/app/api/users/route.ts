// src/app/api/users/route.ts

import { NextResponse } from "next/server"
import { users } from "../_data"

export async function GET() {
  return NextResponse.json(users)
}
