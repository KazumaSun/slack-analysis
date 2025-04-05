import { NextResponse } from 'next/server'

export async function POST() {
  // チャンネルの初期化処理を実装
  console.log("チャンネル初期化処理を実行")
  return NextResponse.json({ message: "チャンネル初期化が完了しました" })
}