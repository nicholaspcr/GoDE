import { useNavigate, useLocation } from 'react-router-dom'
import { LoginForm } from '@/components/auth'

export function LoginPage() {
  const navigate = useNavigate()
  const location = useLocation()

  // Get the page they were trying to visit before being redirected
  const from = location.state?.from?.pathname || '/dashboard'

  const handleSuccess = () => {
    navigate(from, { replace: true })
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-background p-4">
      <div className="w-full max-w-md space-y-6">
        <div className="text-center">
          <h1 className="text-3xl font-bold">GoDE</h1>
          <p className="text-muted-foreground">
            Differential Evolution Optimization Framework
          </p>
        </div>
        <LoginForm onSuccess={handleSuccess} />
      </div>
    </div>
  )
}
