import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

function App() {
  return (
    <div className="min-h-screen bg-background p-8">
      <div className="mx-auto max-w-4xl space-y-8">
        <div className="text-center">
          <h1 className="text-4xl font-bold tracking-tight">GoDE</h1>
          <p className="mt-2 text-muted-foreground">
            Differential Evolution Optimization Framework
          </p>
        </div>

        <div className="grid gap-4 md:grid-cols-2">
          <Card>
            <CardHeader>
              <CardTitle>Run Optimization</CardTitle>
              <CardDescription>
                Configure and execute DE algorithms
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Button className="w-full">New Execution</Button>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>View Results</CardTitle>
              <CardDescription>
                Visualize Pareto fronts and execution history
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Button variant="outline" className="w-full">
                View Executions
              </Button>
            </CardContent>
          </Card>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Getting Started</CardTitle>
          </CardHeader>
          <CardContent className="space-y-2 text-sm text-muted-foreground">
            <p>1. Register or login to your account</p>
            <p>2. Configure a new DE execution with your desired algorithm and problem</p>
            <p>3. Monitor progress in real-time</p>
            <p>4. Visualize the resulting Pareto front</p>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

export default App
