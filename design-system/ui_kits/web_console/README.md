# GoDE Web Console — UI Kit

Recreation of the GoDE web console UI. Built with React 18 + Babel Standalone + Tailwind-like inline styles (no build step) for easy inclusion in design artifacts. Mirrors the shadcn/ui + Tailwind patterns from `web/src/` in the source repo.

## Screens
- **Login** — centered card, wordmark above
- **Dashboard** — 3 quick-action cards + recent executions list
- **New Execution** — stacked form cards (Algorithm, Parameters, GDE3)
- **Execution Detail** — info + config grid, live progress bar, Pareto plot

## Components
- `Primitives.jsx` — Button, Badge, Card, Input, Select, Label, Progress
- `AppShell.jsx` — top header with wordmark + user menu
- `LoginForm.jsx`, `Dashboard.jsx`, `ExecutionForm.jsx`, `ExecutionCard.jsx`, `ExecutionDetail.jsx`, `ParetoPlot.jsx`

## Usage
Open `index.html` — the kit runs as a click-thru prototype: login → dashboard → new execution → detail.
