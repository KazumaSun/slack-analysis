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

  // AppBar 高さを動的に取得
  const appBarHeight = isMobile
    ? (theme.mixins.toolbar?.['@media (min-width:0px)'] as { minHeight?: number })?.minHeight ?? 56
    : (theme.mixins.toolbar?.['@media (min-width:600px)'] as { minHeight?: number })?.minHeight ?? 64

  return (
    <Box
      component="main"
      sx={{
        height: `calc(100vh - ${appBarHeight}px)`,
        overflow: 'hidden',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        px: 2,
        background: 'linear-gradient(to bottom right, #eef2ff, #f3e8ff)',
      }}
    >
      <Stack
        spacing={4}
        alignItems="center"
        textAlign="center"
        maxWidth="md"
      >
        <Typography
          variant={isMobile ? "h4" : "h3"}
          fontWeight="bold"
        >
          Slack Activity Insight
        </Typography>

        <Typography variant={isMobile ? "body1" : "h6"} color="text.secondary">
          チームの空気を、グラフにしよう。
        </Typography>

        <Box>
          <Image
            src="/globe.svg"
            alt="Gopher mascot"
            width={isMobile ? 120 : 160}
            height={isMobile ? 120 : 160}
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
  )
}
