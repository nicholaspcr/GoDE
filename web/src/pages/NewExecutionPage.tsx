import { Link } from 'react-router-dom'
import { ExecutionForm } from '@/components/execution'

export function NewExecutionPage() {
  return (
    <div className="container mx-auto py-8 px-4">
      <div className="mb-6">
        <Link to="/executions" className="text-sm text-muted-foreground hover:underline">
          &larr; Back to Executions
        </Link>
      </div>

      <div className="max-w-4xl">
        <h1 className="text-2xl font-bold mb-6">New Execution</h1>
        <p className="text-muted-foreground mb-8">
          Configure and start a new differential evolution optimization run.
        </p>
        <ExecutionForm />
      </div>
    </div>
  )
}
