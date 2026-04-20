# GoDE Design System

A design system for **GoDE** (Go Differential Evolution) — a production-ready Differential Evolution optimization framework that lets researchers and engineers run multi-objective algorithms (primarily GDE3) against benchmark problems and visualize Pareto fronts through the browser.

## Product context

GoDE is a gRPC/HTTP server + React web client + CLI. The purpose of this project — the design system — is to support the **web UI** that makes it easy to configure DE runs with different variants and problems, monitor progress in real time, and explore results visually.

**Core audience:** optimization researchers, grad students, and engineers. They care about correctness, parameter control, and clean data visualization — not marketing flourish.

**Tone:** technical, dense, plain. No emoji. Neutral chrome (white / near-black / slate) that stays out of the way of scientific plots.

### Products represented
There is **one surface**: the **GoDE Web Console** (`web/`). CLI (`decli`) and server are out of scope for visual design. The console has these primary screens:

1. **Login / Register** — centered card, minimal
2. **Dashboard** — quick-action cards + recent executions
3. **Executions list** — history with status badges
4. **New Execution** — form to configure algorithm + variant + problem + DE parameters
5. **Execution detail** — metadata, live progress, and Pareto plots (2D/3D) with crowding-distance coloring

## Sources

All visual and behavioral ground truth comes from the GoDE codebase:

- **Repo:** https://github.com/nicholaspcr/GoDE (default branch `master`)
- **Parent research project:** https://github.com/nicholaspcr/GDE3 (branch `legacy`) — original academic code from CEFET-MG
- **Implementation plan:** `docs/FRONTEND_IMPLEMENTATION_PLAN.md` in the repo
- **Tech stack (from `web/package.json`):** React 19, TypeScript 5.9, Vite 7, Tailwind CSS 4, TanStack Query, Zustand, React Hook Form + Zod, React Router 7, Plotly.js, Lucide React, shadcn/ui conventions (Radix Slot, class-variance-authority, clsx, tailwind-merge)
- **Theme tokens:** `web/src/index.css` (shadcn default light palette, HSL tokens, 0.5rem radius scale)

## Index

- `README.md` — this file. Product context + content + visual foundations + iconography.
- `colors_and_type.css` — CSS variables for colors, typography, spacing, radii, shadows. Mirrors the HSL tokens from `web/src/index.css`.
- `fonts/` — Inter font files (Google Fonts substitution — see note below).
- `assets/` — logo marks, icons, and the only visual asset in the product (`vite.svg` — Vite default favicon, likely a placeholder for a future GoDE logo).
- `preview/` — individual HTML preview cards that populate the Design System tab.
- `ui_kits/web_console/` — React recreations of the full web UI: login, dashboard, execution form, execution detail with Pareto plot.
- `SKILL.md` — skill manifest so this system can be loaded as an Agent Skill.

### UI kits
- `ui_kits/web_console/` — the GoDE Web Console. Click-thru prototype covering all 5 core screens (login → dashboard → list → new execution → detail with Pareto plot).

## CONTENT FUNDAMENTALS

### Voice
Plain technical English. The product is a lab tool; copy reads like docstrings and form labels, not marketing.

- **Person:** second-person imperative for actions (“Enter your credentials”, “Configure and run a new DE optimization”, “Create your first optimization run!”). No “we”, no “let’s”.
- **Casing:** Title Case for page headings and card titles (“New Execution”, “Algorithm Configuration”, “GDE3 Parameters”). Sentence case for descriptions and helper text. ALL-CAPS reserved for enum labels coming off the API (`EXECUTION_STATUS_RUNNING`) — the UI strips the prefix before display (`Running`, `Completed`).
- **Punctuation:** minimal. Periods in descriptions, not in titles or labels. Ampersand not used.
- **Numbers:** raw numbers, no thousands separators in forms. Percentages rounded to integer (`{Math.round(progress.overallPercent)}%`).

### Vocabulary
Match the domain exactly. Never invent friendlier names for jargon — researchers want the real terms.

- **Execution** (not “job”, “run”, “task”) — one DE optimization invocation.
- **Algorithm** (`gde3`), **Variant** (`rand/1`, `best/2`, `current-to-best/1`, etc.), **Problem** (`zdt1`, `dtlz2`, `wfg4`) — always lowercase in the API, rendered as-is in dropdowns.
- **Parameters:** `CR` (Crossover Rate), `F` (Scaling Factor), `P` (Selection Parameter), `Population Size`, `Dimensions`, `Objectives`, `Generations`, `Floor Limiter`, `Ceil Limiter`, `Crowding Distance`.
- **Status:** Pending, Running, Completed, Failed, Cancelled.

### Microcopy examples (from repo)
- Page header: `GoDE` / `Differential Evolution Optimization Framework`
- Card description: `Configure and run a new DE optimization`
- Empty state: `No executions yet. Create your first optimization run!`
- Button label: `Create Execution`, `Start Execution`, `Cancel`, `Delete`, `View`, `Sign Out`
- Error toast: `Invalid username or password`, `Failed to start execution`
- Success toast: `Successfully logged in!`, `Execution started successfully!`
- Helper: `Welcome, {username}` (no exclamation point)

### No emoji
Zero emoji anywhere in the product. Iconography is Lucide SVG only. Unicode arrows are used sparingly in inline copy (e.g. `← Back to Executions`).

## VISUAL FOUNDATIONS

The aesthetic is **shadcn/ui default light** — a near-white canvas with near-black primary, thin 1px borders, soft shadows, and generous whitespace. Nothing is branded loudly. The Pareto plots are the only saturated color in view (Viridis scale: deep purple → teal → yellow-green).

### Color
All tokens are HSL, defined once in `web/src/index.css` and mirrored in `colors_and_type.css`.

| Token | HSL | Role |
|---|---|---|
| `--background` | `0 0% 100%` | Page canvas |
| `--foreground` | `222.2 84% 4.9%` | Body text, headings |
| `--card` | `0 0% 100%` | Card fill (same as background, separated by border + shadow) |
| `--primary` | `222.2 47.4% 11.2%` | Primary buttons, progress fill, active links (near-black slate) |
| `--primary-foreground` | `210 40% 98%` | Text on primary |
| `--secondary` | `210 40% 96.1%` | Badge secondary, progress track |
| `--muted` | `210 40% 96.1%` | Hover surfaces |
| `--muted-foreground` | `215.4 16.3% 46.9%` | Labels, descriptions, metadata |
| `--accent` | `210 40% 96.1%` | Button ghost/outline hover |
| `--destructive` | `0 84.2% 60.2%` | Delete, errors, failed badge |
| `--border` / `--input` | `214.3 31.8% 91.4%` | 1px borders everywhere |
| `--ring` | `222.2 84% 4.9%` | Focus ring |
| Success | `green-500` (Tailwind) | Completed status semantic |
| Warning | `yellow-500` (Tailwind) | Warning status semantic |
| Viridis colorscale | Plotly default | Crowding-distance coloring on Pareto plots |

There is **no dark theme** in the current build, though `darkMode: ["class"]` is configured in the Tailwind plan (`FRONTEND_IMPLEMENTATION_PLAN.md`). Treat light as canonical.

### Type
- **Font family:** system stack via Tailwind default (`ui-sans-serif, system-ui, -apple-system, "Segoe UI", Roboto, ...`). **Inter is substituted here** as the nearest Google Fonts match for consistent rendering across previews — see note below.
- **Feature settings:** `"rlig" 1, "calt" 1` (required ligatures + contextual alternates) on body. Monospaced digits are used for IDs and numeric values (Tailwind `font-mono`).
- **Scale (in use in the repo):**
  - `text-3xl` 1.875rem / bold — app wordmark on login
  - `text-2xl` 1.5rem / bold — page title (`Execution Details`), card title default
  - `text-xl` 1.25rem / semibold — section heading (`Recent Executions`)
  - `text-lg` 1.125rem / semibold — card section heading (`Algorithm Configuration`)
  - `text-sm` 0.875rem — body default, button text, labels, badges (xs for badges)
  - `text-xs` 0.75rem — metadata under progress bars
- **Weights:** `font-medium` (500) for buttons/labels, `font-semibold` (600) for card titles, `font-bold` (700) for page titles, `font-mono` for IDs.
- **Line-height:** Tailwind defaults. Card titles add `leading-none tracking-tight`.

### Spacing & layout
- Base unit: 4px (Tailwind default).
- Common rhythm: `space-y-2` between form field parts (label → input → error), `space-y-4` between form fields, `space-y-6` between form sections, `gap-6` for card grids, `py-8 px-4` for page padding.
- Page shell: `container mx-auto px-4` + fixed header `py-4` + main `py-8`.
- Forms cap at `max-w-md` for auth, full-width inside cards for execution config.
- Grid breakpoints: single-column mobile → `md:grid-cols-2` → `lg:grid-cols-3`.

### Backgrounds
- Single flat `--background` (white). No gradients, no textures, no hero images, no decorative illustrations. Pareto plots use `paper_bgcolor: 'transparent'` so they inherit card white.

### Corner radii
- `--radius-lg: 0.5rem` (cards, dropdowns, large surfaces)
- `--radius-md: calc(0.5rem - 2px)` (buttons, inputs)
- `--radius-sm: calc(0.5rem - 4px)` (small chips)
- `rounded-full` for badges and progress bars only.

### Borders
- Every card, input, button-outline has a **1px solid border** in `--border` (hsl 214.3 31.8% 91.4%). The border is the primary surface separator — there is no background-tint trick.
- `border-destructive` on error cards.

### Shadows
- `shadow-sm` on cards: the shadcn default — `0 1px 2px 0 rgb(0 0 0 / 0.05)`. Very subtle.
- No other shadow is used. No elevation-on-hover. No inner shadows.

### Animation & transitions
- Buttons/links: `transition-colors` only — fast hover swaps between background tints, no scale, no translate.
- Progress bar fill: `transition-all duration-300 ease-in-out` on width.
- Focus ring: `focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2` — 2px ring, 2px offset. Accessibility-first.
- No bounce, no spring, no entry animations.

### Hover & press states
- **Primary button:** `hover:bg-primary/90` (10% lighter alpha on primary).
- **Destructive button:** `hover:bg-destructive/90`.
- **Outline / ghost button:** `hover:bg-accent hover:text-accent-foreground` (subtle slate-50 fill).
- **Link:** `hover:underline` with `underline-offset-4`.
- **Badge:** `hover:bg-<variant>/80`.
- **Disabled:** `disabled:pointer-events-none disabled:opacity-50`.
- Press states reuse hover colors — no extra shrink or depth change.

### Transparency & blur
- Not used in chrome. The only transparency is `bg-destructive/10` for inline error containers and the `/80` / `/90` hover alpha tints. No backdrop-blur anywhere.

### Cards
A card is: `rounded-lg border bg-card text-card-foreground shadow-sm` with `p-6` default. Header + content are vertically stacked with `space-y-1.5`. Nothing more. Cards do not carry colored left borders, icon headers, or ribbons.

### Imagery
None shipped. Plotly plots are the only visual content beyond type and chrome. Color vibe of plots: **cool-biased scientific** — Viridis (purple/teal/green/yellow), transparent backgrounds, gray gridlines at `rgba(128,128,128,0.2)`.

### Fixed elements
- Top header: `border-b` with the `GoDE` wordmark left-aligned, user greeting + Sign Out right-aligned.
- No sidebar.
- No sticky footer.

## ICONOGRAPHY

- **Library:** [`lucide-react`](https://lucide.dev/) 0.562 (from `web/package.json`). 24×24 default, 2px stroke, rounded linejoin/linecap. Monochrome — inherits `currentColor`.
- **Usage pattern:** icons appear inline with button text at `h-4 w-4` and stand alone inside `size="icon"` buttons at their intrinsic 24×24. Always paired with a visible label or `aria-label`.
- **No icon font.** Everything is tree-shaken SVG from the Lucide React package. For these design-system previews we load Lucide via CDN (`unpkg.com/lucide@latest`) so icons match the real app 1:1.
- **No emoji.** Zero occurrences in the codebase.
- **No custom SVGs** beyond the Vite favicon placeholder (`web/public/vite.svg`), which is not used inside the app chrome. A future GoDE logo is expected but has not been designed — previews use a typographic wordmark (`GoDE`) as the placeholder logo.
- **Unicode chars as icons:** one case only — `←` arrow in the `Back to Executions` link (`&larr;`).
- **Plotly toolbar** ships its own icons (pan/zoom/download). Left as Plotly defaults; `displaylogo: false` removes the Plotly wordmark; `lasso2d` and `select2d` are removed from the 2D modebar.

## Font substitution note

**The repo does not bundle custom font files.** It relies on the OS system font stack. For consistent design-system previews across machines, `fonts/` contains **Inter** (SIL OFL) as a neutral neo-grotesque substitute — Inter’s metrics and humanist-geometric feel are the closest freely available match to the shadcn/Tailwind baseline.

> **Flag for the user:** If you want a different typeface (e.g. Geist, IBM Plex Sans, or a paid type family), drop the files into `fonts/` and update `--font-sans` in `colors_and_type.css`. If GoDE ever ships a logotype, add it to `assets/` and swap the wordmark in the UI kit.
