import { useQuery } from '@tanstack/react-query'
import { deApi } from '../client'
import type { ApiV1ExecutionStatus } from '../client'

const RUNNING_STATUSES: ApiV1ExecutionStatus[] = [
  'EXECUTION_STATUS_PENDING',
  'EXECUTION_STATUS_RUNNING',
]

export const useExecutionProgress = (executionId: string | undefined) => {
  return useQuery({
    queryKey: ['execution', executionId, 'progress'],
    queryFn: async () => {
      if (!executionId) throw new Error('No execution ID')
      const api = deApi()
      return api.differentialEvolutionServiceGetExecutionStatus({ executionId })
    },
    enabled: !!executionId,
    refetchInterval: (query) => {
      const status = query.state.data?.execution?.status
      if (status && RUNNING_STATUSES.includes(status)) {
        // Poll every 2 seconds while running
        return 2000
      }
      // Stop polling when completed/failed/cancelled
      return false
    },
  })
}

export const useExecutionProgressValue = (executionId: string | undefined) => {
  const { data, isLoading, isError } = useExecutionProgress(executionId)

  const progress = data?.progress
  const execution = data?.execution

  const currentGeneration = progress?.currentGeneration
    ? Number(progress.currentGeneration)
    : 0
  const totalGenerations = progress?.totalGenerations
    ? Number(progress.totalGenerations)
    : 1
  const completedExecutions = progress?.completedExecutions
    ? Number(progress.completedExecutions)
    : 0
  const totalExecutions = progress?.totalExecutions
    ? Number(progress.totalExecutions)
    : 1

  const generationPercent =
    totalGenerations > 0 ? (currentGeneration / totalGenerations) * 100 : 0
  const executionPercent =
    totalExecutions > 0 ? (completedExecutions / totalExecutions) * 100 : 0

  // Overall progress combines both generation and execution progress
  const overallPercent =
    totalExecutions > 0
      ? ((completedExecutions + currentGeneration / totalGenerations) /
          totalExecutions) *
        100
      : generationPercent

  return {
    isLoading,
    isError,
    status: execution?.status,
    currentGeneration,
    totalGenerations,
    completedExecutions,
    totalExecutions,
    generationPercent,
    executionPercent,
    overallPercent: Math.min(100, overallPercent),
    partialPareto: progress?.partialPareto,
    isRunning: execution?.status
      ? RUNNING_STATUSES.includes(execution.status)
      : false,
    isCompleted: execution?.status === 'EXECUTION_STATUS_COMPLETED',
    isFailed: execution?.status === 'EXECUTION_STATUS_FAILED',
    isCancelled: execution?.status === 'EXECUTION_STATUS_CANCELLED',
  }
}
