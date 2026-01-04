import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { deApi, type ApiV1DEConfig, type ApiV1ExecutionStatus } from '../client'

interface RunAsyncParams {
  algorithm: string
  variant: string
  problem: string
  deConfig: ApiV1DEConfig
}

export const useExecutions = (status?: ApiV1ExecutionStatus, limit = 50, offset = 0) => {
  return useQuery({
    queryKey: ['executions', { status, limit, offset }],
    queryFn: async () => {
      const api = deApi()
      const response = await api.differentialEvolutionServiceListExecutions({
        status,
        limit,
        offset,
      })
      return response
    },
  })
}

export const useExecution = (executionId: string | undefined) => {
  return useQuery({
    queryKey: ['execution', executionId],
    queryFn: async () => {
      if (!executionId) throw new Error('No execution ID')
      const api = deApi()
      return api.differentialEvolutionServiceGetExecutionStatus({ executionId })
    },
    enabled: !!executionId,
  })
}

export const useExecutionResults = (executionId: string | undefined) => {
  return useQuery({
    queryKey: ['execution', executionId, 'results'],
    queryFn: async () => {
      if (!executionId) throw new Error('No execution ID')
      const api = deApi()
      return api.differentialEvolutionServiceGetExecutionResults({ executionId })
    },
    enabled: !!executionId,
  })
}

export const useRunAsync = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (params: RunAsyncParams) => {
      const api = deApi()
      return api.differentialEvolutionServiceRunAsync({
        body: {
          algorithm: params.algorithm,
          variant: params.variant,
          problem: params.problem,
          deConfig: params.deConfig,
        },
      })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['executions'] })
    },
  })
}

export const useCancelExecution = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (executionId: string) => {
      const api = deApi()
      return api.differentialEvolutionServiceCancelExecution({
        executionId,
        body: {},
      })
    },
    onSuccess: (_, executionId) => {
      queryClient.invalidateQueries({ queryKey: ['execution', executionId] })
      queryClient.invalidateQueries({ queryKey: ['executions'] })
    },
  })
}

export const useDeleteExecution = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (executionId: string) => {
      const api = deApi()
      return api.differentialEvolutionServiceDeleteExecution({ executionId })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['executions'] })
    },
  })
}
