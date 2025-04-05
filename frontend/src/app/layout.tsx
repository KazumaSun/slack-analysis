// src/app/layout.tsx

import './globals.css'
import type { Metadata } from 'next'
import NavigationBar from '@/components/NavigationBar'
import { CssBaseline, Container } from '@mui/material'

export const metadata: Metadata = {
  title: 'SeeLACK',
  description: 'Slackのオンライン傾向・アクティビティを可視化するツール',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  const isHomePage = typeof window !== 'undefined' && window.location.pathname === '/'

  return (
    <html lang="ja">
      <body style={{ overflow: isHomePage ? 'hidden' : 'auto' }}>
        <CssBaseline />
        <NavigationBar />
        <Container
          maxWidth="lg"
          sx={{
            mt: isHomePage ? 0 : 4, // トップページでは余白をなくす
            minHeight: isHomePage ? '100vh' : 'auto', // 高さを柔軟に調整
          }}
        >
          {children}
        </Container>
      </body>
    </html>
  )
}
