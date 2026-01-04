import { useNavigate } from 'react-router-dom'
import { RegisterForm } from '@/components/auth'

export function RegisterPage() {
  const navigate = useNavigate()

  const handleSuccess = () => {
    // After registration, redirect to login
    navigate('/login', {
      state: { message: 'Account created successfully. Please sign in.' }
    })
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
        <RegisterForm onSuccess={handleSuccess} />
      </div>
    </div>
  )
}
