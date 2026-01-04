import { Link } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { useAuthStore } from '@/stores/authStore'
import { useLogout, useExecutions } from '@/api/hooks'

export function DashboardPage() {
  const username = useAuthStore((s) => s.username)
  const logout = useLogout()
  const { data: executionsData, isLoading } = useExecutions()

  const handleLogout = () => {
    logout.mutate()
  }

  const recentExecutions = executionsData?.executions?.slice(0, 5) ?? []

  return (
    <div className="min-h-screen bg-background">
      <header className="border-b">
        <div className="container mx-auto px-4 py-4 flex items-center justify-between">
          <h1 className="text-2xl font-bold">GoDE</h1>
          <div className="flex items-center gap-4">
            <span className="text-sm text-muted-foreground">
              Welcome, {username}
            </span>
            <Button variant="outline" size="sm" onClick={handleLogout}>
              Sign Out
            </Button>
          </div>
        </div>
      </header>

      <main className="container mx-auto px-4 py-8">
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          <Card>
            <CardHeader>
              <CardTitle>New Execution</CardTitle>
              <CardDescription>
                Configure and run a new DE optimization
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Button asChild className="w-full">
                <Link to="/executions/new">Create Execution</Link>
              </Button>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>All Executions</CardTitle>
              <CardDescription>
                View and manage your optimization runs
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Button asChild variant="outline" className="w-full">
                <Link to="/executions">View All</Link>
              </Button>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Quick Stats</CardTitle>
              <CardDescription>
                Your optimization activity
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold">
                {isLoading ? '...' : executionsData?.totalCount ?? 0}
              </div>
              <p className="text-sm text-muted-foreground">Total executions</p>
            </CardContent>
          </Card>
        </div>

        <div className="mt-8">
          <h2 className="text-xl font-semibold mb-4">Recent Executions</h2>
          {isLoading ? (
            <p className="text-muted-foreground">Loading...</p>
          ) : recentExecutions.length === 0 ? (
            <Card>
              <CardContent className="py-8 text-center text-muted-foreground">
                No executions yet. Create your first optimization run!
              </CardContent>
            </Card>
          ) : (
            <div className="space-y-2">
              {recentExecutions.map((execution) => (
                <Card key={execution.id}>
                  <CardContent className="py-4 flex items-center justify-between">
                    <div>
                      <p className="font-medium">{execution.id}</p>
                      <p className="text-sm text-muted-foreground">
                        Status: {execution.status?.replace('EXECUTION_STATUS_', '')}
                      </p>
                    </div>
                    <Button asChild variant="outline" size="sm">
                      <Link to={`/executions/${execution.id}`}>View</Link>
                    </Button>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </div>
      </main>
    </div>
  )
}
