# GoDE TypeScript Frontend Implementation Plan

This document outlines the implementation plan for adding a TypeScript frontend to the GoDE Differential Evolution optimization framework.

## Overview

The frontend will provide:
- User authentication (login/register)
- DE optimization job configuration and execution
- Real-time execution progress monitoring
- Pareto front visualization (2D/3D plots)
- Execution history management

## Technology Stack

| Layer | Technology | Version | Rationale |
|-------|------------|---------|-----------|
| Framework | React | 18.x | Component-based, strong TypeScript support |
| Language | TypeScript | 5.x | Type safety, better DX |
| Build Tool | Vite | 5.x | Fast HMR, optimized builds |
| Server State | TanStack Query | 5.x | Caching, background refetch, mutations |
| Client State | Zustand | 4.x | Minimal, TypeScript-first |
| API Client | OpenAPI Generator | 7.x | Auto-generated from proto specs |
| Visualization | Plotly.js | 2.x | 2D/3D scatter plots |
| UI Components | Shadcn/UI | latest | Accessible, customizable |
| Styling | Tailwind CSS | 3.x | Utility-first CSS |
| Form Validation | Zod + React Hook Form | - | Type-safe validation |
| Testing | Vitest + Playwright | - | Unit + E2E tests |

---

## Phase 1: Project Setup

### 1.1 Initialize Vite + React + TypeScript Project

```bash
cd /home/nick/gh/nicholaspcr/GoDE.git/fork/master
npm create vite@latest web -- --template react-ts
cd web
npm install
```

### 1.2 Install Core Dependencies

```bash
# UI and Styling
npm install -D tailwindcss postcss autoprefixer
npm install clsx tailwind-merge class-variance-authority
npm install lucide-react  # Icons

# State Management
npm install @tanstack/react-query zustand

# Routing
npm install react-router-dom

# Form Handling
npm install react-hook-form @hookform/resolvers zod

# Visualization
npm install plotly.js-dist-min react-plotly.js
npm install -D @types/react-plotly.js

# Development
npm install -D @types/node
```

### 1.3 Configure Tailwind CSS

Create `tailwind.config.js`:
```javascript
/** @type {import('tailwindcss').Config} */
export default {
  darkMode: ["class"],
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: {
        border: "hsl(var(--border))",
        background: "hsl(var(--background))",
        foreground: "hsl(var(--foreground))",
        primary: {
          DEFAULT: "hsl(var(--primary))",
          foreground: "hsl(var(--primary-foreground))",
        },
        // ... additional shadcn colors
      },
    },
  },
  plugins: [],
}
```

### 1.4 Configure Vite

Update `vite.config.ts`:
```typescript
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 5173,
    proxy: {
      '/v1': {
        target: 'http://localhost:8081',
        changeOrigin: true,
      },
    },
  },
})
```

### 1.5 Setup Directory Structure

```
web/
├── src/
│   ├── main.tsx
│   ├── App.tsx
│   ├── index.css
│   ├── api/
│   │   ├── generated/      # OpenAPI Generator output
│   │   ├── client.ts       # API configuration
│   │   └── hooks/          # React Query hooks
│   ├── components/
│   │   ├── ui/             # Shadcn base components
│   │   ├── layout/         # Layout components
│   │   ├── auth/           # Auth components
│   │   ├── execution/      # Execution components
│   │   └── visualization/  # Plot components
│   ├── pages/              # Page components
│   ├── stores/             # Zustand stores
│   ├── lib/                # Utilities
│   └── types/              # TypeScript types
├── public/
├── package.json
├── tsconfig.json
├── vite.config.ts
├── tailwind.config.js
├── postcss.config.js
└── .env.example
```

### 1.6 Add Makefile Targets

Add to root `Makefile`:
```makefile
##@ Frontend

.PHONY: web-deps
web-deps: ## Install frontend dependencies
	cd web && npm install

.PHONY: web-dev
web-dev: ## Run frontend development server
	cd web && npm run dev

.PHONY: web-build
web-build: ## Build frontend for production
	cd web && npm run build

.PHONY: web-test
web-test: ## Run frontend tests
	cd web && npm run test

.PHONY: web-lint
web-lint: ## Lint frontend code
	cd web && npm run lint

.PHONY: web-api
web-api: openapi ## Generate TypeScript API client from OpenAPI spec
	cd web && npx @openapitools/openapi-generator-cli generate \
		-i ../docs/openapi/api.swagger.json \
		-g typescript-fetch \
		-o src/api/generated \
		--additional-properties=typescriptThreePlus=true,supportsES6=true

.PHONY: dev-full
dev-full: ## Run full stack development (backend + frontend)
	@echo "Starting PostgreSQL and Redis..."
	docker compose up -d postgres redis
	@sleep 3
	@echo "Starting backend and frontend..."
	$(MAKE) -j2 run web-dev
```

### 1.7 Setup ESLint and Prettier

ESLint config (`.eslintrc.cjs`):
```javascript
module.exports = {
  root: true,
  env: { browser: true, es2020: true },
  extends: [
    'eslint:recommended',
    'plugin:@typescript-eslint/recommended',
    'plugin:react-hooks/recommended',
  ],
  ignorePatterns: ['dist', '.eslintrc.cjs', 'src/api/generated'],
  parser: '@typescript-eslint/parser',
  plugins: ['react-refresh'],
  rules: {
    'react-refresh/only-export-components': [
      'warn',
      { allowConstantExport: true },
    ],
  },
}
```

Prettier config (`.prettierrc`):
```json
{
  "semi": false,
  "singleQuote": true,
  "tabWidth": 2,
  "trailingComma": "es5"
}
```

---

## Phase 2: API Integration

### 2.1 Generate OpenAPI Spec

Ensure `buf.gen.openapi.yaml` generates the spec:
```bash
make openapi
```

### 2.2 Generate TypeScript Client

```bash
make web-api
```

This generates:
- `src/api/generated/apis/*.ts` - Service API classes
- `src/api/generated/models/*.ts` - TypeScript interfaces

### 2.3 Configure API Client

Create `src/api/client.ts`:
```typescript
import { Configuration } from './generated'
import { useAuthStore } from '@/stores/authStore'

export const createApiConfig = (): Configuration => {
  const token = useAuthStore.getState().accessToken

  return new Configuration({
    basePath: import.meta.env.VITE_API_BASE_URL || '',
    headers: token ? { Authorization: `Bearer ${token}` } : {},
  })
}
```

### 2.4 Create React Query Hooks

Example `src/api/hooks/useAuth.ts`:
```typescript
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { AuthServiceApi } from '../generated'
import { createApiConfig } from '../client'
import { useAuthStore } from '@/stores/authStore'

export const useLogin = () => {
  const setAuth = useAuthStore((s) => s.setAuth)

  return useMutation({
    mutationFn: async (credentials: { username: string; password: string }) => {
      const api = new AuthServiceApi(createApiConfig())
      return api.authServiceLogin({
        body: credentials,
      })
    },
    onSuccess: (data) => {
      setAuth({
        accessToken: data.accessToken!,
        refreshToken: data.refreshToken!,
        expiresIn: data.expiresIn!,
      })
    },
  })
}
```

---

## Phase 3: Authentication UI

### 3.1 Auth Store

Create `src/stores/authStore.ts`:
```typescript
import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface AuthState {
  accessToken: string | null
  refreshToken: string | null
  expiresAt: number | null
  user: { username: string; email: string } | null
  isAuthenticated: boolean
  setAuth: (auth: { accessToken: string; refreshToken: string; expiresIn: number }) => void
  logout: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      accessToken: null,
      refreshToken: null,
      expiresAt: null,
      user: null,
      isAuthenticated: false,
      setAuth: ({ accessToken, refreshToken, expiresIn }) =>
        set({
          accessToken,
          refreshToken,
          expiresAt: Date.now() + expiresIn * 1000,
          isAuthenticated: true,
        }),
      logout: () =>
        set({
          accessToken: null,
          refreshToken: null,
          expiresAt: null,
          user: null,
          isAuthenticated: false,
        }),
    }),
    { name: 'gode-auth' }
  )
)
```

### 3.2 Login Form Component

Create `src/components/auth/LoginForm.tsx`:
```typescript
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useLogin } from '@/api/hooks/useAuth'
import { useNavigate } from 'react-router-dom'

const loginSchema = z.object({
  username: z.string().min(1, 'Username is required'),
  password: z.string().min(1, 'Password is required'),
})

type LoginFormData = z.infer<typeof loginSchema>

export function LoginForm() {
  const navigate = useNavigate()
  const login = useLogin()

  const form = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
  })

  const onSubmit = async (data: LoginFormData) => {
    try {
      await login.mutateAsync(data)
      navigate('/dashboard')
    } catch (error) {
      // Handle error
    }
  }

  return (
    <form onSubmit={form.handleSubmit(onSubmit)}>
      {/* Form fields */}
    </form>
  )
}
```

### 3.3 Protected Route

Create `src/components/auth/ProtectedRoute.tsx`:
```typescript
import { Navigate, Outlet } from 'react-router-dom'
import { useAuthStore } from '@/stores/authStore'

export function ProtectedRoute() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  return <Outlet />
}
```

---

## Phase 4: Execution Management

### 4.1 Execution Form

Create `src/components/execution/ExecutionForm.tsx`:
```typescript
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useSupportedOptions } from '@/api/hooks/useSupportedOptions'
import { useRunAsync } from '@/api/hooks/useExecutions'

const deConfigSchema = z.object({
  algorithm: z.string().min(1),
  variant: z.string().min(1),
  problem: z.string().min(1),
  executions: z.number().int().positive().default(1),
  generations: z.number().int().positive().default(100),
  populationSize: z.number().int().positive().default(100),
  dimensionsSize: z.number().int().positive().default(30),
  objectivesSize: z.number().int().min(2).default(2),
  floorLimiter: z.number().default(0),
  ceilLimiter: z.number().default(1),
  gde3: z.object({
    cr: z.number().min(0).max(1).default(0.9),
    f: z.number().min(0).max(2).default(0.5),
    p: z.number().min(0).max(1).default(0.1),
  }).optional(),
})

export function ExecutionForm() {
  const { algorithms, variants, problems } = useSupportedOptions()
  const runAsync = useRunAsync()

  // Form implementation
}
```

### 4.2 Execution Progress Hook

Create `src/api/hooks/useProgress.ts`:
```typescript
import { useQuery } from '@tanstack/react-query'
import { DifferentialEvolutionServiceApi } from '../generated'
import { createApiConfig } from '../client'

export const useExecutionProgress = (executionId: string) => {
  return useQuery({
    queryKey: ['execution', executionId],
    queryFn: async () => {
      const api = new DifferentialEvolutionServiceApi(createApiConfig())
      return api.differentialEvolutionServiceGetExecutionStatus({
        executionId,
      })
    },
    refetchInterval: (query) => {
      const status = query.state.data?.execution?.status
      if (status === 'EXECUTION_STATUS_RUNNING') return 2000
      if (status === 'EXECUTION_STATUS_PENDING') return 5000
      return false
    },
    enabled: !!executionId,
  })
}
```

### 4.3 Execution List

Create `src/components/execution/ExecutionList.tsx`:
```typescript
import { useExecutions } from '@/api/hooks/useExecutions'
import { ExecutionCard } from './ExecutionCard'

export function ExecutionList() {
  const { data, isLoading } = useExecutions()

  if (isLoading) return <div>Loading...</div>

  return (
    <div className="space-y-4">
      {data?.executions?.map((execution) => (
        <ExecutionCard key={execution.id} execution={execution} />
      ))}
    </div>
  )
}
```

---

## Phase 5: Visualization

### 5.1 2D Pareto Plot

Create `src/components/visualization/ParetoPlot2D.tsx`:
```typescript
import Plot from 'react-plotly.js'
import { Vector } from '@/api/generated'

interface ParetoPlot2DProps {
  vectors: Vector[]
  xAxis: number
  yAxis: number
}

export function ParetoPlot2D({ vectors, xAxis, yAxis }: ParetoPlot2DProps) {
  const data = [{
    x: vectors.map((v) => v.objectives?.[xAxis] ?? 0),
    y: vectors.map((v) => v.objectives?.[yAxis] ?? 0),
    mode: 'markers',
    type: 'scatter',
    marker: {
      color: vectors.map((v) => v.crowdingDistance ?? 0),
      colorscale: 'Viridis',
      showscale: true,
    },
    hovertemplate:
      `Objective ${xAxis + 1}: %{x}<br>` +
      `Objective ${yAxis + 1}: %{y}<br>` +
      `<extra></extra>`,
  }]

  const layout = {
    title: 'Pareto Front',
    xaxis: { title: `Objective ${xAxis + 1}` },
    yaxis: { title: `Objective ${yAxis + 1}` },
  }

  return <Plot data={data} layout={layout} />
}
```

### 5.2 3D Pareto Plot

Create `src/components/visualization/ParetoPlot3D.tsx`:
```typescript
import Plot from 'react-plotly.js'
import { Vector } from '@/api/generated'

interface ParetoPlot3DProps {
  vectors: Vector[]
  xAxis: number
  yAxis: number
  zAxis: number
}

export function ParetoPlot3D({ vectors, xAxis, yAxis, zAxis }: ParetoPlot3DProps) {
  const data = [{
    x: vectors.map((v) => v.objectives?.[xAxis] ?? 0),
    y: vectors.map((v) => v.objectives?.[yAxis] ?? 0),
    z: vectors.map((v) => v.objectives?.[zAxis] ?? 0),
    mode: 'markers',
    type: 'scatter3d',
    marker: {
      size: 4,
      color: vectors.map((v) => v.crowdingDistance ?? 0),
      colorscale: 'Viridis',
    },
  }]

  const layout = {
    title: 'Pareto Front (3D)',
    scene: {
      xaxis: { title: `Obj ${xAxis + 1}` },
      yaxis: { title: `Obj ${yAxis + 1}` },
      zaxis: { title: `Obj ${zAxis + 1}` },
    },
  }

  return <Plot data={data} layout={layout} />
}
```

---

## Phase 6: Docker Integration

### 6.1 Frontend Dockerfile

Create `web/Dockerfile`:
```dockerfile
# Build stage
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
ARG VITE_API_BASE_URL
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL
RUN npm run build

# Production stage
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

### 6.2 Nginx Configuration

Create `web/nginx.conf`:
```nginx
server {
    listen 80;
    root /usr/share/nginx/html;
    index index.html;

    # API proxy
    location /v1/ {
        proxy_pass http://deserver:8081;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # SPA fallback
    location / {
        try_files $uri $uri/ /index.html;
    }

    # Gzip compression
    gzip on;
    gzip_types text/plain text/css application/json application/javascript;
}
```

### 6.3 Docker Compose Service

Add to `docker-compose.yml`:
```yaml
  frontend:
    container_name: gode-frontend
    build:
      context: ./web
      args:
        VITE_API_BASE_URL: ""
    restart: unless-stopped
    depends_on:
      - deserver
    ports:
      - "3001:80"
    networks:
      - gode-network
```

---

## Backend Changes Required

### CORS Configuration

Update `internal/server/middleware/cors.go` for development:
```go
// Add localhost:5173 (Vite dev server) to allowed origins
AllowedOrigins: []string{
    "http://localhost:5173",
    "http://localhost:3001",
}
AllowCredentials: true
```

---

## Testing Strategy

### Unit Tests (Vitest)

```bash
npm install -D vitest @testing-library/react @testing-library/jest-dom jsdom
```

### E2E Tests (Playwright)

```bash
npm install -D @playwright/test
```

---

## Environment Variables

Create `web/.env.example`:
```bash
# API Configuration
VITE_API_BASE_URL=http://localhost:8081

# Feature Flags (optional)
VITE_ENABLE_3D_PLOTS=true
```

---

## Implementation Checklist

### Phase 1: Project Setup
- [ ] Initialize Vite + React + TypeScript project
- [ ] Install and configure Tailwind CSS
- [ ] Set up Shadcn/UI base components
- [ ] Configure path aliases in tsconfig.json
- [ ] Set up ESLint and Prettier
- [ ] Add Makefile targets for frontend
- [ ] Create directory structure

### Phase 2: API Integration
- [ ] Generate OpenAPI spec from protos
- [ ] Generate TypeScript API client
- [ ] Configure TanStack Query provider
- [ ] Create API client configuration
- [ ] Build useAuth hook
- [ ] Build useSupportedOptions hook
- [ ] Build useExecutions hook
- [ ] Build useProgress hook

### Phase 3: Authentication UI
- [ ] Create Zustand auth store with persistence
- [ ] Build LoginForm component
- [ ] Build RegisterForm component
- [ ] Implement ProtectedRoute wrapper
- [ ] Add token refresh logic
- [ ] Test auth flow end-to-end

### Phase 4: Execution Management
- [ ] Build ExecutionForm with validation
- [ ] Create algorithm-specific form fields
- [ ] Implement ExecutionList with pagination
- [ ] Build ExecutionCard component
- [ ] Add ExecutionStatus badge
- [ ] Implement ProgressBar component
- [ ] Add cancel/delete execution actions

### Phase 5: Visualization
- [ ] Integrate Plotly.js
- [ ] Build ParetoPlot2D component
- [ ] Build ParetoPlot3D component
- [ ] Create axis selector for N-D objectives
- [ ] Add ObjectiveTable for tabular view
- [ ] Implement export to PNG/SVG

### Phase 6: Polish & Testing
- [ ] Add loading states and skeletons
- [ ] Implement error boundaries
- [ ] Add toast notifications
- [ ] Make responsive design
- [ ] Write unit tests for key components
- [ ] Add E2E tests with Playwright
- [ ] Create Docker configuration
- [ ] Update docker-compose.yml
