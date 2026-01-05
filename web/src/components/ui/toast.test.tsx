import { describe, it, expect, vi } from 'vitest'
import { render, screen, act } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { ToastProvider, useToast } from './toast'

function TestComponent() {
  const { addToast } = useToast()
  return (
    <div>
      <button onClick={() => addToast('Success message', 'success')}>
        Show Success
      </button>
      <button onClick={() => addToast('Error message', 'error')}>
        Show Error
      </button>
    </div>
  )
}

describe('Toast', () => {
  it('renders toast when addToast is called', async () => {
    const user = userEvent.setup()
    render(
      <ToastProvider>
        <TestComponent />
      </ToastProvider>
    )

    await user.click(screen.getByText('Show Success'))
    expect(screen.getByText('Success message')).toBeInTheDocument()
  })

  it('renders toast with correct type styling', async () => {
    const user = userEvent.setup()
    render(
      <ToastProvider>
        <TestComponent />
      </ToastProvider>
    )

    await user.click(screen.getByText('Show Error'))
    const toast = screen.getByRole('alert')
    expect(toast).toHaveClass('bg-destructive')
  })

  it('removes toast when dismiss button is clicked', async () => {
    const user = userEvent.setup()
    render(
      <ToastProvider>
        <TestComponent />
      </ToastProvider>
    )

    await user.click(screen.getByText('Show Success'))
    expect(screen.getByText('Success message')).toBeInTheDocument()

    await user.click(screen.getByLabelText('Dismiss'))
    expect(screen.queryByText('Success message')).not.toBeInTheDocument()
  })

  it('auto-removes toast after timeout', async () => {
    vi.useFakeTimers()
    render(
      <ToastProvider>
        <TestComponent />
      </ToastProvider>
    )

    await act(async () => {
      screen.getByText('Show Success').click()
    })
    expect(screen.getByText('Success message')).toBeInTheDocument()

    await act(async () => {
      vi.advanceTimersByTime(5000)
    })
    expect(screen.queryByText('Success message')).not.toBeInTheDocument()

    vi.useRealTimers()
  })
})
