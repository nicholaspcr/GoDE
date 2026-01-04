import {
  Configuration,
  ApiV1AuthServiceApi,
  ApiV1DifferentialEvolutionServiceApi,
  ApiV1UserServiceApi,
  ApiV1ParetoServiceApi,
} from './generated'
import { useAuthStore } from '@/stores/authStore'

const BASE_PATH = import.meta.env.VITE_API_BASE_URL || ''

export const createApiConfig = (): Configuration => {
  const token = useAuthStore.getState().accessToken

  return new Configuration({
    basePath: BASE_PATH,
    headers: token ? { Authorization: `Bearer ${token}` } : {},
  })
}

// API client factories
export const authApi = () => new ApiV1AuthServiceApi(createApiConfig())
export const deApi = () => new ApiV1DifferentialEvolutionServiceApi(createApiConfig())
export const userApi = () => new ApiV1UserServiceApi(createApiConfig())
export const paretoApi = () => new ApiV1ParetoServiceApi(createApiConfig())

// Re-export types for convenience
export * from './generated/models'
