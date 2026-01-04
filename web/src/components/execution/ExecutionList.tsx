import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useExecutions } from '@/api/hooks/useExecutions'
import { ExecutionCard } from './ExecutionCard'
import { Button, Select } from '@/components/ui'
import type { ApiV1ExecutionStatus } from '@/api/generated'

const statusOptions: { value: ApiV1ExecutionStatus | ''; label: string }[] = [
  { value: '', label: 'All' },
  { value: 'EXECUTION_STATUS_PENDING', label: 'Pending' },
  { value: 'EXECUTION_STATUS_RUNNING', label: 'Running' },
  { value: 'EXECUTION_STATUS_COMPLETED', label: 'Completed' },
  { value: 'EXECUTION_STATUS_FAILED', label: 'Failed' },
  { value: 'EXECUTION_STATUS_CANCELLED', label: 'Cancelled' },
]

export function ExecutionList() {
  const [statusFilter, setStatusFilter] = useState<ApiV1ExecutionStatus | undefined>()
  const { data, isLoading, isError, refetch } = useExecutions(statusFilter)

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <div className="flex items-center gap-4">
          <h2 className="text-xl font-semibold">Executions</h2>
          <Select
            value={statusFilter ?? ''}
            onChange={(e) =>
              setStatusFilter(e.target.value as ApiV1ExecutionStatus | undefined || undefined)
            }
            className="w-40"
          >
            {statusOptions.map((opt) => (
              <option key={opt.value} value={opt.value}>
                {opt.label}
              </option>
            ))}
          </Select>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" onClick={() => refetch()}>
            Refresh
          </Button>
          <Button asChild>
            <Link to="/executions/new">New Execution</Link>
          </Button>
        </div>
      </div>

      {isLoading && (
        <div className="text-center py-8 text-muted-foreground">
          Loading executions...
        </div>
      )}

      {isError && (
        <div className="text-center py-8 text-destructive">
          Failed to load executions. Please try again.
        </div>
      )}

      {!isLoading && !isError && data?.executions?.length === 0 && (
        <div className="text-center py-8 text-muted-foreground">
          No executions found.{' '}
          <Link to="/executions/new" className="text-primary hover:underline">
            Create your first execution
          </Link>
        </div>
      )}

      {data?.executions?.map((execution) => (
        <ExecutionCard key={execution.id} execution={execution} />
      ))}
    </div>
  )
}
