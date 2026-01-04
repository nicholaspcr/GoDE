import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { ProtectedRoute, AuthProvider } from '@/components/auth'
import { LoginPage, RegisterPage, DashboardPage } from '@/pages'

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
          <Route path="/executions" element={<DashboardPage />} />
          <Route path="/executions/new" element={<DashboardPage />} />
          <Route path="/executions/:id" element={<DashboardPage />} />
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
