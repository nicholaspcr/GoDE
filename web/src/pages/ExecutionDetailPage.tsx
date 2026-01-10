import { useParams, Link, useNavigate } from 'react-router-dom'
import {
  useExecution,
  useExecutionResults,
  useCancelExecution,
  useDeleteExecution,
} from '@/api/hooks/useExecutions'
import { useExecutionProgressValue } from '@/api/hooks/useProgress'
import { Card, Badge, Progress, Button } from '@/components/ui'
import { ParetoVisualization } from '@/components/visualization'
import type { ApiV1ExecutionStatus } from '@/api/generated'

const statusConfig: Record<
  ApiV1ExecutionStatus,
  { label: string; variant: 'default' | 'secondary' | 'destructive' | 'outline' }
> = {
  EXECUTION_STATUS_UNSPECIFIED: { label: 'Unknown', variant: 'outline' },
  EXECUTION_STATUS_PENDING: { label: 'Pending', variant: 'secondary' },
  EXECUTION_STATUS_RUNNING: { label: 'Running', variant: 'default' },
  EXECUTION_STATUS_COMPLETED: { label: 'Completed', variant: 'outline' },
  EXECUTION_STATUS_FAILED: { label: 'Failed', variant: 'destructive' },
  EXECUTION_STATUS_CANCELLED: { label: 'Cancelled', variant: 'secondary' },
}

function formatDate(date: Date | undefined): string {
  if (!date) return '-'
  return new Intl.DateTimeFormat('en-US', {
    dateStyle: 'medium',
    timeStyle: 'medium',
  }).format(date)
}

export function ExecutionDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { data: executionData, isLoading: executionLoading } = useExecution(id)
  const { data: resultsData, isLoading: resultsLoading } = useExecutionResults(id)
  const progress = useExecutionProgressValue(id)
  const cancelExecution = useCancelExecution()
  const deleteExecution = useDeleteExecution()

  if (executionLoading) {
    return (
      <div className="container mx-auto py-8 px-4">
        <div className="text-center text-muted-foreground">Loading...</div>
      </div>
    )
  }

  const execution = executionData?.execution
  if (!execution) {
    return (
      <div className="container mx-auto py-8 px-4">
        <div className="text-center text-destructive">Execution not found</div>
      </div>
    )
  }

  const status = execution.status ?? 'EXECUTION_STATUS_UNSPECIFIED'
  const { label, variant } = statusConfig[status]
  const isRunning = status === 'EXECUTION_STATUS_RUNNING' || status === 'EXECUTION_STATUS_PENDING'
  const isCompleted = status === 'EXECUTION_STATUS_COMPLETED'
  const canDelete = isCompleted ||
                    status === 'EXECUTION_STATUS_FAILED' ||
                    status === 'EXECUTION_STATUS_CANCELLED'

  const handleCancel = () => {
    if (id) {
      cancelExecution.mutate(id)
    }
  }

  const handleDelete = () => {
    if (id && confirm('Are you sure you want to delete this execution?')) {
      deleteExecution.mutate(id, {
        onSuccess: () => navigate('/executions'),
      })
    }
  }

  return (
    <div className="container mx-auto py-8 px-4">
      <div className="mb-6">
        <Link to="/executions" className="text-sm text-muted-foreground hover:underline">
          &larr; Back to Executions
        </Link>
      </div>

      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-4">
          <h1 className="text-2xl font-bold">Execution Details</h1>
          <Badge variant={variant}>{label}</Badge>
        </div>
        <div className="flex gap-2">
          {isRunning && (
            <Button
              variant="outline"
              onClick={handleCancel}
              disabled={cancelExecution.isPending}
            >
              Cancel
            </Button>
          )}
          {canDelete && (
            <Button
              variant="destructive"
              onClick={handleDelete}
              disabled={deleteExecution.isPending}
            >
              Delete
            </Button>
          )}
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card className="p-6">
          <h2 className="text-lg font-semibold mb-4">Information</h2>
          <dl className="space-y-2 text-sm">
            <div className="flex justify-between">
              <dt className="text-muted-foreground">ID</dt>
              <dd className="font-mono">{execution.id}</dd>
            </div>
            <div className="flex justify-between">
              <dt className="text-muted-foreground">Status</dt>
              <dd>{label}</dd>
            </div>
            <div className="flex justify-between">
              <dt className="text-muted-foreground">Problem</dt>
              <dd>{execution.problem || '-'}</dd>
            </div>
            <div className="flex justify-between">
              <dt className="text-muted-foreground">Algorithm</dt>
              <dd>{execution.algorithm || '-'}</dd>
            </div>
            <div className="flex justify-between">
              <dt className="text-muted-foreground">Variant</dt>
              <dd>{execution.variant || '-'}</dd>
            </div>
            <div className="flex justify-between">
              <dt className="text-muted-foreground">Created</dt>
              <dd>{formatDate(execution.createdAt)}</dd>
            </div>
            {execution.completedAt && (
              <div className="flex justify-between">
                <dt className="text-muted-foreground">Completed</dt>
                <dd>{formatDate(execution.completedAt)}</dd>
              </div>
            )}
            {execution.paretoId && (
              <div className="flex justify-between">
                <dt className="text-muted-foreground">Pareto ID</dt>
                <dd className="font-mono">{execution.paretoId}</dd>
              </div>
            )}
          </dl>
        </Card>

        <Card className="p-6">
          <h2 className="text-lg font-semibold mb-4">Configuration</h2>
          {execution.config && (
            <dl className="grid grid-cols-2 gap-2 text-sm">
              <div className="flex justify-between">
                <dt className="text-muted-foreground">Executions</dt>
                <dd>{execution.config.executions}</dd>
              </div>
              <div className="flex justify-between">
                <dt className="text-muted-foreground">Generations</dt>
                <dd>{execution.config.generations}</dd>
              </div>
              <div className="flex justify-between">
                <dt className="text-muted-foreground">Population</dt>
                <dd>{execution.config.populationSize}</dd>
              </div>
              <div className="flex justify-between">
                <dt className="text-muted-foreground">Dimensions</dt>
                <dd>{execution.config.dimensionsSize}</dd>
              </div>
              <div className="flex justify-between">
                <dt className="text-muted-foreground">Objectives</dt>
                <dd>{execution.config.objectivesSize}</dd>
              </div>
              <div className="flex justify-between">
                <dt className="text-muted-foreground">Floor</dt>
                <dd>{execution.config.floorLimiter}</dd>
              </div>
              <div className="flex justify-between">
                <dt className="text-muted-foreground">Ceiling</dt>
                <dd>{execution.config.ceilLimiter}</dd>
              </div>
              {execution.config.gde3 && (
                <>
                  <div className="flex justify-between">
                    <dt className="text-muted-foreground">CR</dt>
                    <dd>{execution.config.gde3.cr}</dd>
                  </div>
                  <div className="flex justify-between">
                    <dt className="text-muted-foreground">F</dt>
                    <dd>{execution.config.gde3.f}</dd>
                  </div>
                  <div className="flex justify-between">
                    <dt className="text-muted-foreground">P</dt>
                    <dd>{execution.config.gde3.p}</dd>
                  </div>
                </>
              )}
            </dl>
          )}
        </Card>
      </div>

      {isRunning && (
        <Card className="p-6 mt-6">
          <h2 className="text-lg font-semibold mb-4">Progress</h2>
          <div className="space-y-4">
            <div>
              <div className="flex justify-between text-sm mb-2">
                <span>Overall Progress</span>
                <span>{Math.round(progress.overallPercent)}%</span>
              </div>
              <Progress value={progress.overallPercent} />
            </div>
            <div className="grid grid-cols-2 gap-4 text-sm">
              <div>
                <span className="text-muted-foreground">Generation: </span>
                {progress.currentGeneration} / {progress.totalGenerations}
              </div>
              <div>
                <span className="text-muted-foreground">Execution: </span>
                {progress.completedExecutions} / {progress.totalExecutions}
              </div>
            </div>
          </div>
        </Card>
      )}

      {execution.error && (
        <Card className="p-6 mt-6 border-destructive">
          <h2 className="text-lg font-semibold mb-4 text-destructive">Error</h2>
          <p className="text-sm text-destructive">{execution.error}</p>
        </Card>
      )}

      {isCompleted && (
        <div className="mt-6">
          <h2 className="text-lg font-semibold mb-4">Results</h2>
          {resultsLoading ? (
            <Card className="p-6">
              <div className="text-muted-foreground">Loading results...</div>
            </Card>
          ) : resultsData?.pareto?.vectors?.length ? (
            <ParetoVisualization vectors={resultsData.pareto.vectors} />
          ) : (
            <Card className="p-6">
              <p className="text-muted-foreground">No results available</p>
            </Card>
          )}
        </div>
      )}
    </div>
  )
}
