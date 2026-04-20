import { Link } from 'react-router-dom'
import { AppShell } from '@/components/layout'
import { ExecutionForm } from '@/components/execution'

export function NewExecutionPage() {
  return (
    <AppShell>
      <div className="mb-6">
        <Link
          to="/executions"
          className="text-muted-foreground text-sm hover:underline"
        >
          &larr; Back to Executions
        </Link>
      </div>

      <div className="max-w-4xl">
        <h1 className="mb-6 text-2xl font-bold tracking-tight">
          New Execution
        </h1>
        <p className="text-muted-foreground mb-8">
          Configure and start a new differential evolution optimization run.
        </p>
        <ExecutionForm />
      </div>
    </AppShell>
  )
}
