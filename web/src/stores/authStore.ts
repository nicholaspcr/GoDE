import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface AuthState {
  accessToken: string | null
  refreshToken: string | null
  expiresAt: number | null
  username: string | null
  isAuthenticated: boolean

  // Actions
  setAuth: (auth: {
    accessToken: string
    refreshToken: string
    expiresIn: number
    username?: string
  }) => void
  setUsername: (username: string) => void
  logout: () => void
  isTokenExpired: () => boolean
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      accessToken: null,
      refreshToken: null,
      expiresAt: null,
      username: null,
      isAuthenticated: false,

      setAuth: ({ accessToken, refreshToken, expiresIn, username }) =>
        set({
          accessToken,
          refreshToken,
          expiresAt: Date.now() + expiresIn * 1000,
          username: username ?? get().username,
          isAuthenticated: true,
        }),

      setUsername: (username) => set({ username }),

      logout: () =>
        set({
          accessToken: null,
          refreshToken: null,
          expiresAt: null,
          username: null,
          isAuthenticated: false,
        }),

      isTokenExpired: () => {
        const { expiresAt } = get()
        if (!expiresAt) return true
        // Consider expired 30 seconds before actual expiry
        return Date.now() > expiresAt - 30000
      },
    }),
    {
      name: 'gode-auth',
      partialize: (state) => ({
        accessToken: state.accessToken,
        refreshToken: state.refreshToken,
        expiresAt: state.expiresAt,
        username: state.username,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
)
