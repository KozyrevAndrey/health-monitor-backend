# Health Monitor Design System

Design system for the redesigned dashboard of **Health Monitor** — a self-hosted uptime-monitoring tool (single Go binary + SQLite) that checks HTTP endpoints, TCP ports and DNS records, opens/resolves incidents, and alerts via Telegram / Email / Gmail / webhooks. Audience: a solo developer or small team; the dashboard stays open on a second monitor for hours, so **glanceability and density beat marketing polish**.

Source: redesign brief pasted in chat (July 2026). No Figma/codebase was attached; components were authored from the brief's deliverables list. Constraint respected throughout: everything must be re-implementable as a **static page with vanilla JS** — so `tokens/components.css` carries the full visual contract as plain CSS classes (`hm-*`), and the React components in `components/` are thin wrappers over those classes.

## Content fundamentals

- Tone: terse, factual, operator-to-operator. "All systems operational", "2 targets down", "Checked 12s ago · every 30s". No exclamation marks, no marketing.
- Sentence case everywhere (buttons, titles, labels). Uppercase only for tiny overline labels (`.overline`).
- Direct address, imperative for actions: "Add your first target", "Clear filters".
- Errors shown verbatim in monospace: `dial tcp 10.0.4.2:5432: i/o timeout`.
- Numbers always tabular (`.num`): `99.982%`, `187 ms`. Units spaced, lowercase (`ms`, `23m`, `1h 05m`).
- **No emoji, ever** — the redesign explicitly removes the old 🔍/➕ emoji-as-icons.

## Visual foundations

- **Dark-first.** Cool blue-black ramp `--bg-inset → --bg-0 (page #0d1117) → --bg-1 (card) → --bg-2 (hover) → --bg-3 (active)`. Light theme via `[data-theme="light"]`, same token names.
- **Status colors do the talking**: green `--status-up`, red `--status-down`, amber `--status-degraded`, gray `--status-unknown`, each with `-subtle` bg and `-border` variants. One palette drives badges, tick bars, charts, toasts, heroes.
- Accent blue `--accent` (#3b82f6) is for **interaction only** (buttons, focus, links, chart line) — never for status.
- Type: Inter (UI) + JetBrains Mono (metrics, URLs, errors). Dense 11–40px scale; base UI text 13px. Google Fonts stand-in — production ships self-hosted woff2.
- Spacing: strict 4px scale (`--space-1…12`). Radii 4/6/10/14. Elevation = 1px border first, soft shadow second; overlays get `--shadow-overlay`.
- Motion: 120–320ms, one gentle ease. Pulse animation is **reserved for ongoing-down** — motion means "needs attention". `prefers-reduced-motion` kills all animation.
- Cards: `--bg-1`, 1px `--border-1`, 10px radius, `--shadow-1`. Hover = border brightens. A DOWN card gets a red border; a DOWN table row gets a red tint. Alarm escalation is border + pulse + tint, never just a recolored pill.
- Accessibility: focus ring token `--focus-ring`, status = icon + text always, 44px touch targets on phones, `aria-live` hero, dialogs with Esc/backdrop close.

## Iconography

Inline SVG stroke icons (Lucide-style: 24px grid, 2px stroke, round caps), authored in `components/core/Icon.jsx` — no icon font, no CDN at runtime, no emoji. Conventions: http → `globe`, tcp → `server`, dns → `at-sign`; telegram → `send`, email/gmail → `mail`, webhook → `link`; down → `alert-triangle`, up → `check`, degraded → `zap`, paused → `pause`.

## Index

- `styles.css` — entry point (imports everything below)
- `tokens/` — `colors.css`, `typography.css`, `spacing.css`, `fonts.css`, `base.css`, `components.css` (the vanilla `hm-*` class layer)
- `guidelines/` — foundation specimen cards
- `components/` — React wrappers, one folder per concern:
  - `core/` — Icon
  - `status/` — StatusDot, StatusBadge, TickBar, Sparkline, LiveIndicator
  - `data/` — StatCard, TargetCard, TargetTable, IncidentItem, NotifierRow
  - `forms/` — Button, IconButton, TextField, SelectField, Switch, Segmented, TagInput, DurationField
  - `feedback/` — Toast, ConfirmDialog, SlideOver, EmptyState, Skeleton
- `ui_kits/dashboard/` — full screens (interactive): `index.html` (incident state), `ok.html` (all-OK), `table.html` (dense variant), `detail.html` (slide-over), `detail-page.html` (full-page variant), `target-form.html`, `notifiers.html`, `notifier-form.html`, `mobile.html`, `login.html` (single-password auth, fits the solo-operator model). Shared: `app.jsx`, `detail.jsx`, `forms.jsx`, `data.js`, `app.css`.

### Components

Icon, StatusDot, StatusBadge, TickBar, Sparkline, LiveIndicator, StatCard, TargetCard, TargetTable, IncidentItem, NotifierRow, Button, IconButton, TextField, SelectField, Switch, Segmented, TagInput, DurationField, Toast, ConfirmDialog, SlideOver, EmptyState, Skeleton.

**Intentional additions** (not in the brief's component list, needed to compose screens): Icon (glyph set), LiveIndicator (brief §1 asks for an SSE indicator), TargetTable (brief's "target table row", delivered as the full table).

## Design decisions to know

- **Themes:** dark (primary), light, and **"Makima" (v2 vibe)** via `[data-theme="makima"]` (`tokens/theme-makima.css`) — warm near-black surfaces, crimson "control" accent, olive-sage up / vermillion down / ringed-eye gold degraded, concentric-ring brand mark, hairline crimson top thread. Screens: `ui_kits/dashboard/v2-makima.html`, `login-v2.html`. Original-design homage to a "calm total control" mood — no character imagery.
- Detail view: **slide-over is the default** (`detail.html`); a full-page drill-down variant exists (`detail-page.html`) — both share `TargetDetailBody`, where the period tabs drive stats *and* chart.
- Overview density: card grid default with a table toggle; `table.html` shows table-first.
- Notifiers are demoted to their own tab — config never competes with live status.
- Tick bar = last 45 checks (UptimeRobot-style), newest right, hover tooltips with time + latency/error.
- Charts in mocks are dependency-free SVG (`ResponseChart` in `ui_kits/dashboard/detail.jsx`); production keeps its bundled Chart.js, colored with `--chart-*` tokens.
