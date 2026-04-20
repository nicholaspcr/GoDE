import type { ApiV1ExecutionStatus } from '@/api/generated'

// UI-friendly label for each ApiV1ExecutionStatus value. Per the GoDE design
// system, enum prefixes are stripped and the result is Title Case
// (Running / Completed / Failed / Cancelled), never ALL-CAPS.
export const executionStatusLabel: Record<ApiV1ExecutionStatus, string> = {
  EXECUTION_STATUS_UNSPECIFIED: 'Unknown',
  EXECUTION_STATUS_PENDING: 'Pending',
  EXECUTION_STATUS_RUNNING: 'Running',
  EXECUTION_STATUS_COMPLETED: 'Completed',
  EXECUTION_STATUS_FAILED: 'Failed',
  EXECUTION_STATUS_CANCELLED: 'Cancelled',
}

// Badge variant per status, matching the design-system UI kit
// (ExecutionDetail.jsx statusMap).
export const executionStatusVariant: Record<
  ApiV1ExecutionStatus,
  'default' | 'secondary' | 'destructive' | 'outline'
> = {
  EXECUTION_STATUS_UNSPECIFIED: 'outline',
  EXECUTION_STATUS_PENDING: 'secondary',
  EXECUTION_STATUS_RUNNING: 'default',
  EXECUTION_STATUS_COMPLETED: 'outline',
  EXECUTION_STATUS_FAILED: 'destructive',
  EXECUTION_STATUS_CANCELLED: 'secondary',
}
