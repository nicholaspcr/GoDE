import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { ProtectedRoute, AuthProvider } from '@/components/auth'
import {
  LoginPage,
  RegisterPage,
  DashboardPage,
  ExecutionsPage,
  NewExecutionPage,
  ExecutionDetailPage,
} from '@/pages'

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
        {/* Public routes */}
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />

        {/* Protected routes */}
        <Route element={<ProtectedRoute />}>
          <Route path="/dashboard" element={<DashboardPage />} />
          <Route path="/executions" element={<ExecutionsPage />} />
          <Route path="/executions/new" element={<NewExecutionPage />} />
          <Route path="/executions/:id" element={<ExecutionDetailPage />} />
        </Route>

        {/* Redirect root to dashboard (will redirect to login if not authenticated) */}
        <Route path="/" element={<Navigate to="/dashboard" replace />} />

        {/* Catch all - redirect to dashboard */}
        <Route path="*" element={<Navigate to="/dashboard" replace />} />
      </Routes>
      </AuthProvider>
    </BrowserRouter>
  )
}

export default App
