// src/app/layout.tsx

import './globals.css'
import type { Metadata } from 'next'
import NavigationBar from '@/components/NavigationBar'
import { CssBaseline, Container } from '@mui/material'

export const metadata: Metadata = {
  title: 'Slack Activity Insight',
  description: 'Slackのオンライン傾向・アクティビティを可視化するツール',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="ja">
      <body>
        <CssBaseline />
        <NavigationBar />
        <Container maxWidth="lg" sx={{ mt: 4 }}>
          {children}
        </Container>
      </body>
    </html>
  )
}
