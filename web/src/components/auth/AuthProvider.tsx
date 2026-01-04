import { useTokenRefresh } from '@/api/hooks'

interface AuthProviderProps {
  children: React.ReactNode
}

export function AuthProvider({ children }: AuthProviderProps) {
  // Set up automatic token refresh
  useTokenRefresh()

  return <>{children}</>
}
