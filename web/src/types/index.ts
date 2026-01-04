// Re-export generated API types
// export * from '@/api/generated/models'

// Additional frontend-specific types

export type ExecutionStatus =
  | 'EXECUTION_STATUS_UNSPECIFIED'
  | 'EXECUTION_STATUS_PENDING'
  | 'EXECUTION_STATUS_RUNNING'
  | 'EXECUTION_STATUS_COMPLETED'
  | 'EXECUTION_STATUS_FAILED'
  | 'EXECUTION_STATUS_CANCELLED'

export interface AuthTokens {
  accessToken: string
  refreshToken: string
  expiresIn: number
}

export interface User {
  username: string
  email: string
}
