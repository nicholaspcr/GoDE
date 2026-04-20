import { z } from 'zod'

// Schema for all Vite env vars the app reads. Validated once at module load
// so a missing or malformed var fails fast with a clear message instead of
// producing a runtime 404 deep in an API call.
const envSchema = z.object({
  // Base URL for the backend. Empty in dev (Vite proxy handles /v1/*). In
  // production builds this must be either a fully-qualified URL or an
  // empty string for same-origin deployments behind a reverse proxy.
  VITE_API_BASE_URL: z
    .string()
    .default('')
    .refine(
      (value) => value === '' || /^https?:\/\//.test(value),
      'VITE_API_BASE_URL must be empty or start with http:// or https://',
    ),

  // Feature flag for the 3D Pareto plot. Accepts "true"/"false" strings
  // because Vite injects env vars as strings.
  VITE_ENABLE_3D_PLOTS: z
    .enum(['true', 'false'])
    .default('true')
    .transform((value) => value === 'true'),
})

const parsed = envSchema.safeParse(import.meta.env)

if (!parsed.success) {
  const issues = parsed.error.issues
    .map((issue) => `  - ${issue.path.join('.')}: ${issue.message}`)
    .join('\n')
  throw new Error(`Invalid environment configuration:\n${issues}`)
}

export const env = {
  API_BASE_URL: parsed.data.VITE_API_BASE_URL,
  ENABLE_3D_PLOTS: parsed.data.VITE_ENABLE_3D_PLOTS,
} as const
