import { useQuery } from '@tanstack/react-query'
import { deApi } from '../client'

export const useAlgorithms = () => {
  return useQuery({
    queryKey: ['supported', 'algorithms'],
    queryFn: async () => {
      const api = deApi()
      const response = await api.differentialEvolutionServiceListSupportedAlgorithms()
      return response.algorithms ?? []
    },
    staleTime: Infinity, // These don't change during a session
  })
}

export const useVariants = () => {
  return useQuery({
    queryKey: ['supported', 'variants'],
    queryFn: async () => {
      const api = deApi()
      const response = await api.differentialEvolutionServiceListSupportedVariants()
      return response.variants ?? []
    },
    staleTime: Infinity,
  })
}

export const useProblems = () => {
  return useQuery({
    queryKey: ['supported', 'problems'],
    queryFn: async () => {
      const api = deApi()
      const response = await api.differentialEvolutionServiceListSupportedProblems()
      return response.problems ?? []
    },
    staleTime: Infinity,
  })
}

export const useSupportedOptions = () => {
  const algorithms = useAlgorithms()
  const variants = useVariants()
  const problems = useProblems()

  return {
    algorithms: algorithms.data ?? [],
    variants: variants.data ?? [],
    problems: problems.data ?? [],
    isLoading: algorithms.isLoading || variants.isLoading || problems.isLoading,
    isError: algorithms.isError || variants.isError || problems.isError,
  }
}
