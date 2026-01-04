import { useMutation, useQueryClient } from '@tanstack/react-query'
import { authApi } from '../client'
import { useAuthStore } from '@/stores/authStore'

interface LoginCredentials {
  username: string
  password: string
}

interface RegisterData {
  username: string
  email: string
  password: string
}

export const useLogin = () => {
  const setAuth = useAuthStore((s) => s.setAuth)
  const setUsername = useAuthStore((s) => s.setUsername)

  return useMutation({
    mutationFn: async (credentials: LoginCredentials) => {
      const api = authApi()
      return api.authServiceLogin({
        body: {
          username: credentials.username,
          password: credentials.password,
        },
      })
    },
    onSuccess: (data, variables) => {
      if (data.accessToken && data.refreshToken && data.expiresIn) {
        setAuth({
          accessToken: data.accessToken,
          refreshToken: data.refreshToken,
          expiresIn: Number(data.expiresIn),
          username: variables.username,
        })
        setUsername(variables.username)
      }
    },
  })
}

export const useRegister = () => {
  return useMutation({
    mutationFn: async (data: RegisterData) => {
      const api = authApi()
      return api.authServiceRegister({
        body: {
          user: {
            ids: { username: data.username },
            email: data.email,
            password: data.password,
          },
        },
      })
    },
  })
}

export const useLogout = () => {
  const logout = useAuthStore((s) => s.logout)
  const username = useAuthStore((s) => s.username)
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async () => {
      if (username) {
        const api = authApi()
        await api.authServiceLogout({
          body: { username },
        })
      }
    },
    onSettled: () => {
      logout()
      queryClient.clear()
    },
  })
}

export const useRefreshToken = () => {
  const refreshToken = useAuthStore((s) => s.refreshToken)
  const setAuth = useAuthStore((s) => s.setAuth)
  const logout = useAuthStore((s) => s.logout)

  return useMutation({
    mutationFn: async () => {
      if (!refreshToken) {
        throw new Error('No refresh token available')
      }
      const api = authApi()
      return api.authServiceRefreshToken({
        body: { refreshToken },
      })
    },
    onSuccess: (data) => {
      if (data.accessToken && data.refreshToken && data.expiresIn) {
        setAuth({
          accessToken: data.accessToken,
          refreshToken: data.refreshToken,
          expiresIn: Number(data.expiresIn),
        })
      }
    },
    onError: () => {
      logout()
    },
  })
}
