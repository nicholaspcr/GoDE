import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useNavigate } from 'react-router-dom'
import { useSupportedOptions } from '@/api/hooks/useSupportedOptions'
import { useRunAsync } from '@/api/hooks/useExecutions'
import { Button, Input, Label, Card, Select, useToast } from '@/components/ui'

const deConfigSchema = z.object({
  algorithm: z.string().min(1, 'Algorithm is required'),
  variant: z.string().min(1, 'Variant is required'),
  problem: z.string().min(1, 'Problem is required'),
  executions: z.number().int().positive(),
  generations: z.number().int().positive(),
  populationSize: z.number().int().positive(),
  dimensionsSize: z.number().int().positive(),
  objectivesSize: z.number().int().min(2),
  floorLimiter: z.number(),
  ceilLimiter: z.number(),
  gde3Cr: z.number().min(0).max(1),
  gde3F: z.number().min(0).max(2),
  gde3P: z.number().min(0).max(1),
})

type DEConfigFormData = z.infer<typeof deConfigSchema>

interface ExecutionFormProps {
  onSuccess?: (executionId: string) => void
}

export function ExecutionForm({ onSuccess }: ExecutionFormProps) {
  const navigate = useNavigate()
  const { algorithms, variants, problems, isLoading: optionsLoading } = useSupportedOptions()
  const runAsync = useRunAsync()
  const { addToast } = useToast()

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<DEConfigFormData>({
    resolver: zodResolver(deConfigSchema),
    defaultValues: {
      executions: 1,
      generations: 100,
      populationSize: 100,
      dimensionsSize: 30,
      objectivesSize: 2,
      floorLimiter: 0,
      ceilLimiter: 1,
      gde3Cr: 0.9,
      gde3F: 0.5,
      gde3P: 0.1,
    },
  })

  const onSubmit = async (data: DEConfigFormData) => {
    try {
      const response = await runAsync.mutateAsync({
        algorithm: data.algorithm,
        variant: data.variant,
        problem: data.problem,
        deConfig: {
          executions: String(data.executions),
          generations: String(data.generations),
          populationSize: String(data.populationSize),
          dimensionsSize: String(data.dimensionsSize),
          objectivesSize: String(data.objectivesSize),
          floorLimiter: data.floorLimiter,
          ceilLimiter: data.ceilLimiter,
          gde3: {
            cr: data.gde3Cr,
            f: data.gde3F,
            p: data.gde3P,
          },
        },
      })
      if (response.executionId) {
        addToast('Execution started successfully!', 'success')
        onSuccess?.(response.executionId)
        navigate(`/executions/${response.executionId}`)
      }
    } catch {
      addToast('Failed to start execution', 'error')
    }
  }

  if (optionsLoading) {
    return <div className="text-muted-foreground">Loading options...</div>
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      <Card className="p-6">
        <h3 className="text-lg font-semibold mb-4">Algorithm Configuration</h3>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="space-y-2">
            <Label htmlFor="algorithm">Algorithm</Label>
            <Select {...register('algorithm')}>
              <option value="">Select algorithm</option>
              {algorithms.map((alg) => (
                <option key={alg} value={alg}>
                  {alg}
                </option>
              ))}
            </Select>
            {errors.algorithm && (
              <p className="text-sm text-destructive">{errors.algorithm.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="variant">Variant</Label>
            <Select {...register('variant')}>
              <option value="">Select variant</option>
              {variants.map((v) => (
                <option key={v.name} value={v.name}>
                  {v.name}
                </option>
              ))}
            </Select>
            {errors.variant && (
              <p className="text-sm text-destructive">{errors.variant.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="problem">Problem</Label>
            <Select {...register('problem')}>
              <option value="">Select problem</option>
              {problems.map((p) => (
                <option key={p.name} value={p.name}>
                  {p.name}
                </option>
              ))}
            </Select>
            {errors.problem && (
              <p className="text-sm text-destructive">{errors.problem.message}</p>
            )}
          </div>
        </div>
      </Card>

      <Card className="p-6">
        <h3 className="text-lg font-semibold mb-4">Execution Parameters</h3>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div className="space-y-2">
            <Label htmlFor="executions">Executions</Label>
            <Input
              type="number"
              {...register('executions', { valueAsNumber: true })}
              min={1}
            />
            {errors.executions && (
              <p className="text-sm text-destructive">{errors.executions.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="generations">Generations</Label>
            <Input
              type="number"
              {...register('generations', { valueAsNumber: true })}
              min={1}
            />
            {errors.generations && (
              <p className="text-sm text-destructive">{errors.generations.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="populationSize">Population Size</Label>
            <Input
              type="number"
              {...register('populationSize', { valueAsNumber: true })}
              min={1}
            />
            {errors.populationSize && (
              <p className="text-sm text-destructive">{errors.populationSize.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="dimensionsSize">Dimensions</Label>
            <Input
              type="number"
              {...register('dimensionsSize', { valueAsNumber: true })}
              min={1}
            />
            {errors.dimensionsSize && (
              <p className="text-sm text-destructive">{errors.dimensionsSize.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="objectivesSize">Objectives</Label>
            <Input
              type="number"
              {...register('objectivesSize', { valueAsNumber: true })}
              min={2}
            />
            {errors.objectivesSize && (
              <p className="text-sm text-destructive">{errors.objectivesSize.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="floorLimiter">Floor Limiter</Label>
            <Input
              type="number"
              step="0.1"
              {...register('floorLimiter', { valueAsNumber: true })}
            />
            {errors.floorLimiter && (
              <p className="text-sm text-destructive">{errors.floorLimiter.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="ceilLimiter">Ceil Limiter</Label>
            <Input
              type="number"
              step="0.1"
              {...register('ceilLimiter', { valueAsNumber: true })}
            />
            {errors.ceilLimiter && (
              <p className="text-sm text-destructive">{errors.ceilLimiter.message}</p>
            )}
          </div>
        </div>
      </Card>

      <Card className="p-6">
        <h3 className="text-lg font-semibold mb-4">GDE3 Parameters</h3>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="space-y-2">
            <Label htmlFor="gde3Cr">CR (Crossover Rate)</Label>
            <Input
              type="number"
              step="0.01"
              {...register('gde3Cr', { valueAsNumber: true })}
              min={0}
              max={1}
            />
            {errors.gde3Cr && (
              <p className="text-sm text-destructive">{errors.gde3Cr.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="gde3F">F (Scaling Factor)</Label>
            <Input
              type="number"
              step="0.01"
              {...register('gde3F', { valueAsNumber: true })}
              min={0}
              max={2}
            />
            {errors.gde3F && (
              <p className="text-sm text-destructive">{errors.gde3F.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="gde3P">P (Selection Parameter)</Label>
            <Input
              type="number"
              step="0.01"
              {...register('gde3P', { valueAsNumber: true })}
              min={0}
              max={1}
            />
            {errors.gde3P && (
              <p className="text-sm text-destructive">{errors.gde3P.message}</p>
            )}
          </div>
        </div>
      </Card>

      {runAsync.error && (
        <div className="p-4 bg-destructive/10 text-destructive rounded-md">
          Failed to start execution. Please try again.
        </div>
      )}

      <div className="flex justify-end gap-4">
        <Button
          type="button"
          variant="outline"
          onClick={() => navigate('/executions')}
        >
          Cancel
        </Button>
        <Button type="submit" disabled={isSubmitting || runAsync.isPending}>
          {runAsync.isPending ? 'Starting...' : 'Start Execution'}
        </Button>
      </div>
    </form>
  )
}
