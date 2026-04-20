import type { ReactNode } from 'react'
import { Link } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { useAuthStore } from '@/stores/authStore'
import { useLogout } from '@/api/hooks'

interface AppShellProps {
  children: ReactNode
}

// Top-level page shell for authenticated routes. Renders the border-b header
// with the GoDE wordmark, user greeting, and Sign Out control, then the page
// content in a centered container. Matches design-system/ui_kits/web_console
// /AppShell.jsx exactly.
export function AppShell({ children }: AppShellProps) {
  const username = useAuthStore((s) => s.username)
  const logout = useLogout()

  return (
    <div className="bg-background text-foreground min-h-screen">
      <header className="border-b">
        <div className="container mx-auto flex items-center justify-between px-4 py-4">
          <Link
            to="/dashboard"
            className="text-2xl font-bold tracking-tight hover:no-underline"
          >
            GoDE
          </Link>
          <div className="flex items-center gap-4">
            {username && (
              <span className="text-muted-foreground text-sm">
                Welcome, {username}
              </span>
            )}
            <Button
              variant="outline"
              size="sm"
              onClick={() => logout.mutate()}
              disabled={logout.isPending}
            >
              Sign Out
            </Button>
          </div>
        </div>
      </header>
      <main className="container mx-auto px-4 py-8">{children}</main>
    </div>
  )
}
