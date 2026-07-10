# Health Monitor — Dashboard Redesign Brief

## What this app is

A self-hosted **uptime monitoring** tool (like a lightweight UptimeRobot / BetterStack): a single Go binary with SQLite that checks HTTP endpoints, TCP ports and DNS records on a schedule, records check results, opens/resolves **incidents** on failures, and sends alerts through configurable **notifiers** (Telegram, Email/SMTP, Gmail, generic webhook). The web UI is a single-page dashboard served by the same binary.

**Audience:** a solo developer / small team watching their own services. Density and glanceability matter more than marketing polish. The dashboard often stays open on a second monitor for hours.

## Current UI (what to improve)

One HTML page (`web/static/index.html`, vanilla JS + Chart.js, no build step) containing, top to bottom:

1. **Header** — purple gradient banner (`#667eea → #764ba2`), emoji "🔍" in the title.
2. **4 stat cards** — Total Targets / Healthy / Failing / Active Incidents (plain numbers).
3. **"Monitored Targets" section** — a vertical list of bordered cards: name, status badge (OK / FAIL / UNKNOWN pill), type, interval, last response time, Edit/Delete/Details buttons.
4. **"Notifiers" section** — list of cards with a type badge (telegram/gmail/email/webhook), enabled/disabled toggle pill, Edit/Delete.
5. **"Recent Incidents" section** — left-border colored cards (red ongoing / green resolved) with target name, duration, failure count, last error.
6. **Modals** — Add/Edit Target form (type-dependent fields for HTTP/TCP/DNS), Add/Edit Notifier form (5 notifier types with very different fields), and a Target Detail modal: period switcher (24h/7d/30d), stat cards, response-time line chart (last 100 checks), recent-checks list, incidents, and a raw JSON `<pre>` config dump.

Live updates come via SSE (`/api/v1/events`, `check` and `alert` events) with a 30s polling fallback; "Last updated" timestamp at the bottom.

### Known weaknesses to fix

- **No dark mode** — painful for an always-on monitoring screen; monitoring tools are dark-first by convention.
- **Weak status hierarchy** — when something is DOWN, nothing screams. A failing target looks almost identical to a healthy one (small pill changes color). The #1 job of this screen — "is anything broken right now?" — takes scanning.
- **No uptime visualization on cards** — no 90-day / 24h tick bars (the UptimeRobot-style green/red segment strip), no sparklines. Response-time data and uptime % exist in the API but are invisible until you open the detail modal.
- **Everything is one long page** — targets, notifiers and incidents stacked vertically; notifiers (rarely touched config) get the same visual weight as live status.
- **Emoji as icons** (🔍 ➕ 🔄 💾) — inconsistent rendering, unprofessional; needs an inline SVG icon set.
- **Browser `confirm()`/`alert()` dialogs** for delete confirmation and errors; no toast notifications for SSE alerts or CRUD results.
- **Detail view is a modal** — cramped for a chart + table + incidents + config; deserves a proper drill-down panel or page-like view.
- **Raw JSON config dump** in detail view instead of formatted key-value display.
- **No search / filter / sort** for targets; fine at 5 targets, unusable at 30.
- **No table/list density option** — cards only.
- **Forms are long unstructured stacks** of inputs inside modals; type-dependent fields appear/disappear with no grouping or steps.
- Hardcoded colors everywhere, no CSS custom properties / design tokens; many inline `style="..."` attributes in markup.
- **Zero `@media` queries** — responsiveness relies solely on `auto-fit/minmax` grids; the target-card action row (Details/Edit/Delete + badge) overflows on phones; touch targets are below 44px.
- **Detail chart ignores the selected period** — the 24h/7d/30d switcher only updates the stat cards; the chart always shows the last 100 checks.
- **Accessibility gaps:** status conveyed by color alone, modals without `role="dialog"`/focus trap/Esc-to-close, no `aria-live` for live updates, some grey text (`#9ca3af` on white) borderline for WCAG AA.

## Data available from the backend (design around real data)

REST API (`/api/v1/...`) + SSE:

- **Target:** `id`, `name`, `type` (http | tcp | dns), type config (URL / host:port / domain+record), `interval`, `timeout`, `enabled`, `tags[]`, `description`.
- **Check history** per target (up to N results): `status` (success | failure | warning | unknown), `response_time_ms`, `status_code`, `error`, `message`, `checked_at`. Enough for sparklines, tick bars, and response-time charts.
- **Stats** per target per period (24h / 7d / 30d): `uptime_percentage`, `avg/min/max_response_time_ms`, `total/successful/failed_checks`, `current_status`, `consecutive_failures`, `last_check_at`, `last_success_at`, `last_failure_at`.
- **Incidents:** `status` (ongoing | resolved), `started_at`, `resolved_at`, `duration`, `failure_count`, `last_error`, `severity` (info | warning | critical); global list, ongoing-only, and per-target lists.
- **Notifiers:** `id`, `type` (telegram | email | gmail | gmail_oauth | webhook), `enabled`, type-specific config (secrets masked as `***`).
- **SSE events:** every completed check (`target_id`, `status`, `response_time_ms`, `checked_at`) and every alert (down / up / slow_response / ssl_expiring) in real time.

## What to design

A redesigned dashboard, keeping the single-page/no-build-step spirit, with:

1. **Status-first overview.** A hero strip that answers "is everything OK?" instantly: e.g., a big "All systems operational" state vs. an alarm state listing what's down. Stat cards can support it but shouldn't be the headline. Show a subtle live indicator (SSE connected / reconnecting / polling fallback).
2. **Target cards/rows with embedded history:** status dot + name + type icon, an UptimeRobot-style tick bar of recent checks (green/red/yellow segments with hover tooltips), current uptime % and last response time, sparkline optional. Provide both a comfortable card grid and a compact table view toggle. Search box + status filter (All / Down / Degraded / Paused) + sort.
3. **Target detail as a slide-over panel or full drill-down view** (not a cramped modal): large status header, period tabs (24h/7d/30d) that drive **both** the stats **and** the chart, uptime + response-time stat row, response-time chart with success/failure point coloring, incident timeline, recent checks table, and a human-readable configuration block with an Edit action.
4. **Incidents:** ongoing incidents surfaced at the top of the page (alarm state), resolved history in a collapsible/timeline list with severity color coding and durations.
5. **Notifiers demoted to a settings area** (tab or secondary section), with cleaner type badges/icons and enable toggles as proper switches.
6. **Forms:** modal or slide-over forms with grouped sections (Identity → Check settings → Advanced), segmented control for target type, inline validation, duration inputs (30s / 1m / 5m presets), tag input; notifier form with type picker cards instead of a `<select>`.
7. **Feedback:** toast notifications (target down/up alerts from SSE, CRUD success/failure), styled confirm dialog for deletes, skeleton loading states, designed empty states ("No targets yet → Add your first target").
8. **Design system:** CSS custom properties for all tokens; semantic status palette (success/warning/danger/neutral) consistent across badges, tick bars, charts, toasts; **dark theme as the primary theme with a light option** (or `prefers-color-scheme` + toggle); system font stack or a single modern sans (e.g., Inter); tabular numerals for metrics; consistent inline SVG icon set (Lucide/Heroicons style); 4px spacing scale; visible focus states and WCAG AA contrast.
9. **Accessibility built in:** proper dialog semantics (focus trap, Esc-to-close, focus return), status never conveyed by color alone (icon + text), `aria-live` region for real-time status changes, ≥44px touch targets, explicit responsive breakpoints (mobile / tablet / desktop).

**Aesthetic direction:** modern observability tool — think BetterStack / Checkly / Grafana Cloud: calm dark surfaces, restrained accent color, status colors doing the talking, dense but breathable layout. Avoid the current "purple gradient + emoji" look.

**Constraints:** must remain implementable as a static page with vanilla JS (no React/build step); Chart.js is already bundled and self-hosted; no external fonts/CDNs required at runtime; desktop-first but sensibly responsive down to mobile.

## Deliverables

- Design tokens (colors incl. dark/light themes, typography, spacing, radii, shadows)
- Components: status badge/dot, uptime tick bar, sparkline, stat card, target card + table row, incident item, notifier row, toast, confirm dialog, modal/slide-over, form controls, tabs/segmented control, empty & skeleton states
- Screens: dashboard overview (all-OK state and incident/alarm state), target detail panel, add/edit target form, notifier settings, mobile layout of the overview
