import { useEffect, useRef } from 'react'
import { useAuthStore } from '@/stores/authStore'
import { useRefreshToken } from './useAuth'

const REFRESH_THRESHOLD_MS = 60 * 1000 // Refresh 60 seconds before expiry

export function useTokenRefresh() {
  const expiresAt = useAuthStore((s) => s.expiresAt)
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)
  const refreshToken = useRefreshToken()
  const refreshTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  useEffect(() => {
    // Clear any existing timeout
    if (refreshTimeoutRef.current) {
      clearTimeout(refreshTimeoutRef.current)
      refreshTimeoutRef.current = null
    }

    // Don't schedule refresh if not authenticated or no expiry time
    if (!isAuthenticated || !expiresAt) {
      return
    }

    const timeUntilExpiry = expiresAt - Date.now()
    const timeUntilRefresh = timeUntilExpiry - REFRESH_THRESHOLD_MS

    // If token is already expired or about to expire, refresh immediately
    if (timeUntilRefresh <= 0) {
      refreshToken.mutate()
      return
    }

    // Schedule refresh
    refreshTimeoutRef.current = setTimeout(() => {
      refreshToken.mutate()
    }, timeUntilRefresh)

    return () => {
      if (refreshTimeoutRef.current) {
        clearTimeout(refreshTimeoutRef.current)
      }
    }
  }, [expiresAt, isAuthenticated, refreshToken])
}
