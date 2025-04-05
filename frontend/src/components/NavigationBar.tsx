// src/components/NavigationBar.tsx
'use client'

import { AppBar, Toolbar, Typography, Button, Box } from '@mui/material'
import Link from 'next/link'
import { usePathname } from 'next/navigation'

const navItems = [
  { label: 'トップ', href: '/' },
  { label: 'ダッシュボード', href: '/dashboard' },
  { label: 'チャンネル/ユーザー一覧', href: '/users' },
]

export default function NavigationBar() {
  const pathname = usePathname()

  return (
    <AppBar position="static" color="primary">
      <Toolbar>
        <Typography variant="h6" sx={{ flexGrow: 1 }}>
          <Link href="/" style={{ textDecoration: 'none', color: 'inherit' }}>
            SeeLACK
          </Link>
        </Typography>
        <Box>
          {navItems.map((item) => (
            <Link key={item.href} href={item.href} passHref>
              <Button
                color={pathname === item.href ? 'inherit' : 'secondary'}
                sx={{ color: '#fff' }}
              >
                {item.label}
              </Button>
            </Link>
          ))}
        </Box>
      </Toolbar>
    </AppBar>
  )
}
