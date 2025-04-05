import { NextResponse } from 'next/server'

export async function POST() {
  // ユーザーの初期化処理を実装
  console.log("ユーザー初期化処理を実行")
  return NextResponse.json({ message: "ユーザー初期化が完了しました" })
}