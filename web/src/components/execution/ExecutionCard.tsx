import { Link } from 'react-router-dom'
import { Card, Badge, Progress, Button } from '@/components/ui'
import { useCancelExecution, useDeleteExecution } from '@/api/hooks/useExecutions'
import { useExecutionProgressValue } from '@/api/hooks/useProgress'
import type { ApiV1Execution, ApiV1ExecutionStatus } from '@/api/generated'

interface ExecutionCardProps {
  execution: ApiV1Execution
}

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
    dateStyle: 'short',
    timeStyle: 'short',
  }).format(date)
}

export function ExecutionCard({ execution }: ExecutionCardProps) {
  const cancelExecution = useCancelExecution()
  const deleteExecution = useDeleteExecution()
  const progress = useExecutionProgressValue(execution.id)

  const status = execution.status ?? 'EXECUTION_STATUS_UNSPECIFIED'
  const { label, variant } = statusConfig[status]
  const isRunning = status === 'EXECUTION_STATUS_RUNNING' || status === 'EXECUTION_STATUS_PENDING'
  const canDelete = status === 'EXECUTION_STATUS_COMPLETED' ||
                    status === 'EXECUTION_STATUS_FAILED' ||
                    status === 'EXECUTION_STATUS_CANCELLED'

  const handleCancel = () => {
    if (execution.id) {
      cancelExecution.mutate(execution.id)
    }
  }

  const handleDelete = () => {
    if (execution.id && confirm('Are you sure you want to delete this execution?')) {
      deleteExecution.mutate(execution.id)
    }
  }

  return (
    <Card className="p-4">
      <div className="flex items-start justify-between">
        <div className="space-y-1">
          <div className="flex items-center gap-2">
            <Link
              to={`/executions/${execution.id}`}
              className="font-medium hover:underline"
            >
              Execution {execution.id?.slice(0, 8)}...
            </Link>
            <Badge variant={variant}>{label}</Badge>
          </div>
          <p className="text-sm text-muted-foreground">
            Created: {formatDate(execution.createdAt)}
          </p>
          {execution.completedAt && (
            <p className="text-sm text-muted-foreground">
              Completed: {formatDate(execution.completedAt)}
            </p>
          )}
        </div>

        <div className="flex gap-2">
          {isRunning && (
            <Button
              size="sm"
              variant="outline"
              onClick={handleCancel}
              disabled={cancelExecution.isPending}
            >
              Cancel
            </Button>
          )}
          {canDelete && (
            <Button
              size="sm"
              variant="destructive"
              onClick={handleDelete}
              disabled={deleteExecution.isPending}
            >
              Delete
            </Button>
          )}
          <Button size="sm" variant="outline" asChild>
            <Link to={`/executions/${execution.id}`}>View</Link>
          </Button>
        </div>
      </div>

      {isRunning && (
        <div className="mt-4 space-y-2">
          <div className="flex justify-between text-sm">
            <span>Progress</span>
            <span>{Math.round(progress.overallPercent)}%</span>
          </div>
          <Progress value={progress.overallPercent} />
          <div className="flex justify-between text-xs text-muted-foreground">
            <span>
              Generation {progress.currentGeneration}/{progress.totalGenerations}
            </span>
            <span>
              Execution {progress.completedExecutions}/{progress.totalExecutions}
            </span>
          </div>
        </div>
      )}

      {execution.config && (
        <div className="mt-4 grid grid-cols-2 md:grid-cols-4 gap-2 text-sm">
          <div>
            <span className="text-muted-foreground">Generations:</span>{' '}
            {execution.config.generations}
          </div>
          <div>
            <span className="text-muted-foreground">Population:</span>{' '}
            {execution.config.populationSize}
          </div>
          <div>
            <span className="text-muted-foreground">Dimensions:</span>{' '}
            {execution.config.dimensionsSize}
          </div>
          <div>
            <span className="text-muted-foreground">Objectives:</span>{' '}
            {execution.config.objectivesSize}
          </div>
        </div>
      )}

      {execution.error && (
        <div className="mt-4 p-2 bg-destructive/10 text-destructive text-sm rounded">
          {execution.error}
        </div>
      )}
    </Card>
  )
}
