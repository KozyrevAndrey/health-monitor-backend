# Health Monitor - Project Summary

## Date: 2026-06-27

A snapshot of the project's current state. For the authoritative feature checklist
see [`implementation-progress.md`](implementation-progress.md); for deferred work
see [`backlog.md`](backlog.md).

---

## What's implemented

### Infrastructure & core (Phase 0-1) ✅
- Go module, project structure, Makefile (20+ commands)
- Docker (multi-stage, CGO for SQLite) + docker-compose (dev & prod)
- Domain layer: Target, CheckResult, Alert, Incident, NotifierConfig + interfaces
- Config (Viper, YAML, env override), logging (zerolog), graceful shutdown

### Storage (Phase 2) ✅
- GORM + SQLite, AutoMigrate, connection pooling
- Repositories: Target, CheckResult (history/stats/retention), Incident, NotifierConfig
- Targets and notifiers are managed via the API/UI (DB-backed), not from YAML

### Checkers (Phase 3-5) ✅
- **HTTP/HTTPS** — status code, response time, SSL validation + expiry warning,
  custom headers/methods, redirects, skip-verify
- **TCP** — port availability + connection time
- **DNS** — A/AAAA/CNAME/MX/TXT/NS, custom DNS server, expected-value validation
- Thread-safe checker registry; type-aware target form in the UI
- (ICMP/Ping is deferred — needs CAP_NET_RAW; see backlog)

### Scheduler (Phase 4) ✅
- Ticker-based, independent per-target intervals, concurrent execution
- AddTarget/RemoveTarget/UpdateTarget, graceful shutdown
- (Worker pool / jitter / backoff deferred — see backlog)

### Alerting (Phase 5-6) ✅
- AlertManager: per-target state, consecutive fail/success, rule engine
- Rules: consecutive failures, response-time threshold, SSL expiry, DOWN/UP
- Incident creation/resolution, notifier broadcast, hot-reload

### Notifiers (Phase 7) ✅
- **Telegram** (Bot API, multiple chats, Markdown, **http/https/socks5 proxy**)
- **Email/SMTP** (auth, TLS incl. implicit 465, HTML+plaintext)
- **Gmail** Service Account (domain delegation) and **Gmail OAuth2**
- **Webhook** (generic: url/method/headers/templated payload, retry+backoff)
- CRUD via REST API + hot-reload + secret masking

### HTTP API (Phase 8) ✅
- Chi router, middleware (Request ID, Real IP, logging, recovery, timeout, CORS,
  **Basic Auth**)
- Targets CRUD + results + stats; Incidents (list/get/ongoing + per-target);
  Notifiers CRUD
- **`GET /health`** with real probes (DB + scheduler → healthy/degraded/unhealthy)
- OpenAPI/Swagger via `oapi-codegen`

### Web Dashboard (Phase 9) ✅
- Vanilla-JS SPA, full CRUD of targets and notifiers
- **Real-time via SSE** (`/api/v1/events`) with incremental in-place updates;
  polling only as a fallback
- **Target detail modal**: response-time chart (self-hosted Chart.js),
  period-switchable stats (24h/7d/30d), recent checks, incidents, config

### Data Retention ✅
- Background cleanup of old check results and resolved incidents per `RetentionConfig`

### Deployment & CI ✅
- Production docker-compose + Traefik labels, `.env`
- **CI** (`test.yml`): vet + build + `go test -short -race` (blocking) +
  golangci-lint (advisory)
- **Release** (`release.yml`): on tag `v*` — push Docker image to GHCR + binary
  to the GitHub Release

---

## Project structure

```
health-monitor-backend/
├── cmd/server/main.go              # Entry point, lifecycle wiring
├── internal/
│   ├── domain/                     # Models & interfaces
│   ├── storage/                    # GORM + SQLite repositories + models/
│   ├── checker/                    # http.go, tcp.go, dns.go, registry.go
│   ├── scheduler/                  # Ticker-based scheduler
│   ├── alerting/                   # AlertManager (rules, incidents)
│   ├── notifier/                   # telegram, email, gmail, gmail_oauth, webhook
│   ├── events/                     # In-process pub/sub broker (SSE)
│   ├── retention/                  # Background cleanup job
│   ├── api/                        # Chi server, handlers, SSE, health
│   └── generated/                  # oapi-codegen output
├── pkg/config, pkg/logger
├── web/static/                     # SPA + vendored chart.umd.min.js
├── migrations/
├── configs/example.yaml
├── .github/workflows/              # test.yml, release.yml
├── docker-compose.yml, docker-compose.prod.yml, Dockerfile
└── docs/                           # plan, progress, checkers, backlog, this file
```

~50 Go source files, 10 test files across alerting/checker/events/notifier/retention/api.

---

## Current status

Functionally complete for self-hosted monitoring:
- 3 checkers (HTTP/TCP/DNS), 5 notifiers, smart alerting with incidents
- Real-time dashboard with charts, retention, health check, CI/CD

Remaining (see [`backlog.md`](backlog.md)): ICMP checker, worker pool, Prometheus
`/metrics`, `validate`/`backup` CLI, lint cleanup, multi-arch images.

---

## Commands to run

```bash
# Local development
make build          # Build the project
make test           # Run tests
make run            # Run the application

# Docker (dev)
docker compose up -d --build   # Build & run (rebuild needed after web/ changes)
docker compose logs -f
docker compose down

# Tests
go test ./...                  # All tests
go test -short -race ./...     # What CI runs (skips external-network tests)
```

---

## Important notes

1. **SQLite + CGO**: builds require `CGO_ENABLED=1` (gcc); handled in Dockerfile and CI.
2. **Database path**: `./data/health-monitor.db` (created automatically).
3. **Config**: via YAML or env vars (prefix `HEALTH_MONITOR_`). Targets/notifiers
   live in the DB, managed via the API/UI.
4. **Static assets are baked into the Docker image** — after changing `web/`,
   rebuild with `docker compose up -d --build` and hard-refresh the browser.
5. **Architecture**: Clean Architecture, interface-based DI, graceful shutdown,
   nil-safe optional dependencies (event publisher, DB pinger) wired via setters.

---

## Recent git history

```bash
git log --oneline -10
```

```
docs: translate all docs to English and add GHCR deploy to backlog
ci: add test and release GitHub Actions workflows
feat(api): add real health check for database and scheduler
feat(ui): add target detail view with response-time chart
feat(api): add per-target incidents endpoint and harden stats period
perf(ui): apply SSE events incrementally without full refetch
feat(events): add SSE real-time stream for checks and alerts
feat(retention): add background cleanup job for old results and incidents
feat(checker): add DNS checker with record validation and UI
feat(checker): add TCP checker with type-aware target UI
```

---

**The project is production-ready; continue from the backlog.** 🚀
