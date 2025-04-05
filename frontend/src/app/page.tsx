'use client'

import Link from "next/link"
import Image from "next/image"
import {
  Box,
  Button,
  Typography,
  useTheme,
  useMediaQuery,
  Stack,
} from "@mui/material"

export default function Home() {
  const theme = useTheme()

  // 画面サイズの判定（レスポンシブ対応）
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'))

  return (
    <div
      style={{
        height: '100vh', // ビューポート全体の高さ
        overflow: 'hidden', // スクロールを無効化
      }}
    >
      <Box
        component="main"
        sx={{
          height: '100vh',
          display: 'flex',
          justifyContent: 'flex-start',
          alignItems: 'center',
          px: 2,
          pt: 8,
        }}
      >
        <Stack
          spacing={4}
          alignItems="center"
          textAlign="center"
          maxWidth="md"
          sx={{
            maxHeight: '100vh', // スタック全体の高さを制限
          }}
        >
          <Typography
            variant={isMobile ? "h4" : "h3"}
            fontWeight="bold"
          >
            SeeLACK
          </Typography>

          <Typography variant={isMobile ? "body1" : "h6"} color="text.secondary">
            チームの空気を、グラフにしよう。
          </Typography>

          <Box
            sx={{
              animation: 'float 3s ease-in-out infinite',
              '@keyframes float': {
                '0%': { transform: 'translateY(0px)' },
                '50%': { transform: 'translateY(-10px)' },
                '100%': { transform: 'translateY(0px)' },
              },
            }}
          >
            <Image
              src="/globe.svg"
              alt="Gopher mascot"
              width={isMobile ? 100 : 140}
              height={isMobile ? 100 : 140}
            />
          </Box>

          <Link href="/dashboard" passHref>
            <Button
              variant="contained"
              size={isMobile ? "medium" : "large"}
            >
              ▶ 使ってみる
            </Button>
          </Link>
        </Stack>
      </Box>
    </div>
  )
}
