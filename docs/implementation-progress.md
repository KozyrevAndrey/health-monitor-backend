# Health Monitor - Implementation Progress

> The document reflects the **actual** state of the code (cross-checked against git and the sources).

---

## ✅ Completed

### Phase 0-1: Infrastructure and basic architecture
- [x] Go module, project structure, `.gitignore`/`.dockerignore`
- [x] Makefile (20+ commands), Dockerfile (multi-stage, CGO for SQLite), docker-compose
- [x] `.golangci.yml`, README.md
- [x] Domain Layer (`internal/domain/`): Target, CheckResult, Alert, Incident, NotifierConfig + interfaces (Checker, *Repository, Notifier, Scheduler, AlertManager, CheckerRegistry, NotifierRepository)
- [x] Configuration (`pkg/config/`): Viper, YAML, env override (explicit binding of env vars — `cc21b15`)
- [x] Logging (`pkg/logger/`): zerolog, levels, json/console, stdout/stderr/file
- [x] Entry point (`cmd/server/main.go`): graceful shutdown, signal handling, lifecycle

### Phase 2: Storage Layer
- [x] Migrations (`migrations/001_initial_schema.sql`), tables targets/check_results/incidents/alerts + 14 indexes
- [x] GORM + SQLite, AutoMigrate, connection pooling, auto-create db dir
- [x] Repositories: TargetRepository, CheckResultRepository (history, stats, DeleteOlderThan), IncidentRepository
- [x] **NotifierRepository** + model/migration (`8155946`) — notifier configs in the DB

### Phase 3-5: Checkers
- [x] Checker registry (thread-safe), resolving the checker by `target.Type`
- [x] **HTTP/HTTPS** checker: status code, response time, SSL validation + expiry warning, custom headers/methods, redirects, skip-verify
- [x] **TCP** checker (`tcp.go`): port availability check + connection time, host/port/address metadata ✨
- [x] **DNS** checker (`dns.go`): resolving A/AAAA/CNAME/MX/TXT/NS, resolution time, custom DNS server, validation against expected values, record dedup ✨
- [x] Tests (`http_test.go`, `tcp_test.go`, `dns_test.go`)
- [x] UI: **type-aware** target form (HTTP url/status · TCP host/port · DNS domain/record/server/expected), endpoint in the card

### Phase 4: Scheduler
- [x] Ticker-based, independent per-target intervals, concurrent execution (goroutine per target)
- [x] AddTarget/RemoveTarget/UpdateTarget, graceful shutdown (WaitGroup)
- [x] Loading targets, saving results, processing through AlertManager

### Phase 5-6: Alert System
- [x] AlertManager (`internal/alerting/manager.go`): state per target, consecutive fail/success, rule engine
- [x] Rules: consecutive failures, response-time threshold, SSL expiry, DOWN/UP transitions
- [x] Incident management (creation/resolution, last error)
- [x] Notifier registry + broadcast, integration with the scheduler
- [x] **ClearNotifiers** for hot-reload (`f9cb8d2`)
- [x] Tests (`internal/alerting/manager_test.go`)

### Phase 7: Notification System
- [x] **Telegram** (`telegram.go`): Bot API, multiple chat IDs, Markdown, icons/severity, **proxy (http/https/socks5)** ✨
- [x] **Email/SMTP** (`email.go`): auth, TLS/SSL (incl. implicit TLS on 465 — `f573485`), HTML+plaintext multipart
- [x] **Gmail Service Account** (`gmail.go`): domain-wide delegation, impersonation
- [x] **Gmail OAuth2** (`gmail_oauth.go`): personal account via credentials/token files
- [x] **Webhook (generic)** (`webhook.go`): url/method/custom headers/payload template (text/template), retry+backoff on 5xx/429, optional proxy ✨
- [x] Tests for telegram/email/webhook
- [x] **Notifier CRUD via REST API** (`1034d6f`) + hot-reload + secret masking; notifiers are managed from the DB/UI, not from YAML (`8518538`)

### Phase 8: HTTP API
- [x] Chi router v5, configurable timeouts, graceful shutdown
- [x] Middleware: Request ID, Real IP, logging, recovery, timeout, CORS, **Basic Auth** (`fbdbe5a`)
- [x] Targets CRUD + results + stats; Incidents (list/get/ongoing + per-target); Notifiers CRUD
- [x] **`GET /health`** with real probes: ping DB + scheduler, healthy/degraded/unhealthy statuses (503 on unhealthy) ✨
- [x] **OpenAPI/Swagger**: `oapi-codegen` (`internal/generated/api.gen.go`), swagger handlers (`swagger_handlers.go`, `openapi_adapter.go`)

### Phase 9: Web Dashboard
- [x] SPA (`web/static/index.html`), vanilla JS, responsive layout
- [x] Stats cards, target cards with statuses, latest incidents
- [x] **Full CRUD of targets and notifiers** via UI modals (`c274a8d`, `c776017`); per-type notifier fields
- [x] **Real-time via SSE** (`GET /api/v1/events`): live updates of statuses and alerts, fallback polling 30s ✨
- [x] **Target detail modal** with a response time chart (Chart.js, self-hosted), stat cards with period switching (24h/7d/30d), a table of the latest checks, target incidents, and config ✨
- [x] Static file server via Chi

### Real-time (SSE)
- [x] In-process event broker (`internal/events/broker.go`): pub/sub, non-blocking, drop-on-full ✨
- [x] Publishing `check`/`alert` events from the scheduler and alert manager (nil-safe `SetEventPublisher`)
- [x] SSE endpoint outside the logging/Timeout group; removing the write deadline via `ResponseController`, heartbeat 25s
- [x] UI on `EventSource` with **incremental updates** from the payload: `check` updates the card/counters without backend requests, `alert` — only the incidents panel; full refresh only as a fallback
- [x] Verified: the stream stays alive >90s (passed WriteTimeout 15s and chi Timeout 60s)

### Data Retention
- [x] Background cleanup (`internal/retention/cleaner.go`): deletion of old check_results and resolved incidents per `RetentionConfig` (`cleanup_interval`), initial sweep on start, graceful stop via ctx ✨
- [x] `DeleteResolvedOlderThan` in IncidentRepository, cleaner tests

### Deployment & CI
- [x] **Production docker-compose** (`docker-compose.prod.yml`) + Traefik labels + `DOMAIN`/basic-auth env (`a3a9f9f`, `5d4d98a`, `578b240`)
- [x] `.env` / `.env.example`, mount secrets volume (`4bd3014`)
- [x] **CI** (`.github/workflows/test.yml`): go vet + build + `go test -short -race` (blocking) + golangci-lint (advisory, ~200 stylistic findings on non-linted code) on push/PR ✨
- [x] **Release** (`.github/workflows/release.yml`): on tag `v*` — push the Docker image to GHCR + binary to the GitHub Release ✨
- [x] External network tests (HTTP checker) gated by `testing.Short()` for stable CI

---

## ⏳ Not Done (per plan)

- [ ] **Additional checkers**: ICMP/Ping (http.go + tcp.go + dns.go exist)
- [ ] **Worker pool** in the scheduler (currently goroutine-per-target without a bounded pool/queue)
- [ ] **Prometheus `/metrics`**
- [ ] **CLI**: `validate`, `backup`/`restore`
- [ ] Load/security tests (CI/CD — done)
- [ ] PostgreSQL adapter, RBAC/multi-user, maintenance windows, status page, alert routing/escalation, quiet hours

---

## 📊 Current Status

| Component | % |
|---|---|
| Infrastructure / Domain / Config / Logging | 100% |
| Storage (SQLite + repos + notifier configs) | 100% |
| Checkers: HTTP + TCP + DNS | 100% |
| Scheduler (without worker pool) | 90% |
| Alert Manager (+ hot-reload) | 100% |
| Notifiers: Telegram(+proxy)/Email/Gmail SA/Gmail OAuth/Webhook | 100% |
| HTTP API (+ OpenAPI/Swagger, session login) | 100% |
| Web Dashboard (full CRUD, SSE real-time) | 100% |
| Prod deploy (Traefik) | 100% |
| Data Retention (cleanup job) | 100% |
| Real-time (SSE) | 100% |
| Additional checkers (ICMP) | 0% |
| Charts (target detail, Chart.js) | 100% |
| Prometheus metrics | 0% |

**Last Updated:** 2026-06-27

---

## 🚀 Next Steps (priority)

1. **Worker pool** in the scheduler
2. **Prometheus `/metrics`**
3. **ICMP/Ping checker** ([backlog](backlog.md) — requires CAP_NET_RAW privileges)
4. CLI `validate`/`backup`

---

## 📚 Docs

- [`docs/checkers.md`](checkers.md) — how to use HTTP/TCP/DNS checks (config fields, API/UI examples, retention)
- [`docs/backlog.md`](backlog.md) — deferred tasks (ICMP/Ping and others)

## 📝 Notes

- Clean Architecture, interface-based DI, graceful shutdown everywhere
- Targets and notifiers are managed **only via API/UI** (removed from YAML)
- A new notifier type is wired up in 3 layers: `internal/notifier/<type>.go` → two switches (`cmd/server/main.go`, `internal/api/notifier_handlers.go`) → UI (`web/static/index.html`)
- Notifier configs are validated via the constructor on create/update (broken url/template/proxy → 400)
- Known masking limitations: values inside webhook `headers` (e.g. `Authorization`) and credentials in `proxy_url` (`socks5://user:pass@...`) are **not masked** in API responses — only top-level sensitive keys are masked (`bot_token`, `password`, `token`, `secret`, `smtp_password`)
- Known dispatch limitation: AlertManager stores notifiers in a map keyed by `Type()` (one per type) and dispatches them **sequentially** within a single alert (bounded by the check context `target.Timeout+5s`)
