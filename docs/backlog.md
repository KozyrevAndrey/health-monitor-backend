# Backlog

Deferred tasks (not in active work). The prioritized "what's next" list lives in
[`implementation-progress.md`](implementation-progress.md) → Next Steps.

---

## ICMP / Ping checker

**What it is:** host liveness via ICMP Echo (like `ping`): reachability, RTT
(round-trip time), packet loss. The lowest-level check — no ports, no
applications. Useful for routers/switches/servers that expose no HTTP/TCP service.

**Why deferred — requires privileges.** Raw ICMP sockets need root / the Linux
`CAP_NET_RAW` capability, or unprivileged mode via a UDP socket (`ip4:icmp`) with
the `net.ipv4.ping_group_range` sysctl allowed. In Docker that means
`cap_add: [NET_RAW]` or `sysctls`. So unlike TCP/DNS it won't "just work"
everywhere — it needs deployment changes.

**How to implement (same pattern as TCP/DNS):**
- `internal/checker/icmp.go`, register it in `NewDefaultRegistry()`.
- `ICMPConfig` already exists in domain: `host`, `packet_count`, `max_rtt_ms`.
- Library: `github.com/prometheus-community/pro-bing` (supports both privileged
  and unprivileged modes).
- Result: `success` if no loss and RTT within bounds; `warning` on loss or
  exceeding `max_rtt_ms`; `failure` if the host doesn't respond.
- Add ICMP fields to the type-aware target form (`web/static/index.html`).
- Document the `NET_RAW`/sysctl requirement + a Docker example in `docs/checkers.md`.

---

## Lint cleanup → make lint blocking

**Current state:** the `lint` job in CI is marked `continue-on-error: true`
(advisory) — it reports but does not fail CI. Reason: the strict `.golangci.yml`
surfaces **~201 findings** against a codebase that was never linted before.

**Findings breakdown (v1.64.8):** revive 79 (mostly "exported … should have
comment"), errcheck 34, gocritic 27, errorlint 7, unused 6, govet 6, unparam 5,
gosec 4, goconst 4, gosimple 3, misspell 2, gocyclo 2, contextcheck 2.

**How to close (pick one):**
- *Incrementally, by category:* first the potentially-real ones (errcheck,
  unused, gosec, govet, contextcheck), defer pure style (revive `exported`,
  gocritic `opinionated/experimental`). Drop `continue-on-error` as you clean up.
- *Soften the config:* remove the most pedantic rules (revive `exported`,
  gocritic `opinionated`/`experimental` tags) → ~two dozen left; fix those and
  make lint blocking.

**Files:** `.golangci.yml`, `.github/workflows/test.yml` (the `continue-on-error`
flag).

---

## Worker pool in the scheduler (+ jitter / backoff)

**Current state:** `internal/scheduler/scheduler.go` spawns a **goroutine per
target** with its own ticker; there is no global cap on concurrent checks. Fine
for tens–hundreds of targets.

**Why:** with hundreds–thousands of targets, or when ticks coincide, many checks
can start at once → spikes in network/CPU/FD usage and SQLite write load. A pool
decouples "how many targets" from "how many checks run at once".

**How to implement:**
- Tickers enqueue a "check X" job onto a channel; N workers pull and execute
  (configurable worker count).
- Add **jitter** at startup (spread the first checks over time — against the
  thundering herd) and **backoff** on error streaks (grow the interval, return to
  normal on success). See the plan, Phase 5.1–5.2.
- Optional metrics: queue depth, worker utilization (see Prometheus below).

---

## Prometheus `/metrics`

**Why:** integration with an existing Prometheus + Grafana + Alertmanager stack.
Adopt it **only if you have/plan such a stack** — otherwise the built-in
dashboard (SSE, charts), stats (uptime/response) and notifiers already cover the
need.

**Two angles:**
- *About the monitor itself:* active targets, queue depth, check duration
  P50/P95/P99, error rate, goroutines, DB connection pool (Phase 9.1).
- *Monitoring data as metrics:* `health_target_up{target=…}`,
  `health_target_response_ms{…}` → dashboards/alerts in Grafana/Alertmanager.

**How:** `github.com/prometheus/client_golang`, serve `/metrics` (usually outside
the login auth), increment at the same scheduler/alert-manager points that already
emit SSE events.

---

## CLI: `validate` and `backup`/`restore`

**Current state:** the binary accepts only `--config` and `--version`
(flag-based).

- **`validate`** (Phase 10.2) — validate config without starting: valid YAML,
  required fields, sane intervals/timeouts, DB path. `validate --config
  config.yaml` → exit 0/1. Cheap, handy in CI / before deploy.
- **`backup` / `restore`** (Phase 10.4) — a correct live SQLite backup
  (`VACUUM INTO` / the Online Backup API, not `cp`). Medium value — a volume/file
  snapshot is often enough.

For subcommands — either `cobra`/`urfave/cli` (a `main.go` refactor) or hand-parse
`os.Args[1]`.

---

## Deploying to a server via GHCR

**Goal:** the last mile — tag → CI pushes the image to GHCR → run it on the prod
server. Today `docker-compose.prod.yml` **builds the image locally** (`build:`);
switch it to **pull** from the registry.

**1. One-time: edit `docker-compose.prod.yml`** — replace the build block with an
image reference:
```yaml
  health-monitor:
    image: ghcr.io/kozyrevandrey/health-monitor:${TAG:-latest}
    container_name: health-monitor-prod
    # build: ...  ← remove entirely
```
Add `TAG=v1.0.0` to `.env`. Prefer pinning a concrete version in prod over
`latest` (reproducible, explicit about what's deployed).

**2. One-time: log in to GHCR on the server**
- If the package is **public** (toggle in the package settings on GitHub) — no
  login needed, `pull` works anonymously.
- If **private**:
  ```bash
  echo $GHCR_PAT | docker login ghcr.io -u KozyrevAndrey --password-stdin
  ```
  where `GHCR_PAT` is a Personal Access Token (classic) with the **`read:packages`**
  scope (read is enough for pull).

**3. Deploy (each release)** on the server, in the compose dir:
```bash
# (if pinning a version) update TAG in .env, e.g. TAG=v1.0.1
docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml up -d
docker image prune -f   # clean up old layers
```
The `health-monitor-data` volume (DB) and `./secrets` are preserved.

**Full release cycle:**
```
git tag v1.0.1 && git push origin v1.0.1
   → release.yml builds & pushes ghcr.io/.../health-monitor:1.0.1 + :latest
   → on the server: TAG=v1.0.1, docker compose pull && up -d
```

**Optional automation:**
- **Watchtower** — a container that watches GHCR and auto-updates the service on
  a new image (give it GHCR creds for a private package). Simple, but "magic" and
  tied to `latest`.
- **CI auto-deploy over SSH** — a `release.yml` job that SSHes into the server and
  runs `pull && up -d` (needs secrets: SSH key, host). Controlled, but opens CI
  access to the server.
- **Manual pull** (the flow above) — simplest and safest; recommended to start.

---

## CI/CD — extensions

- **Multi-arch Docker** (arm64): the `Dockerfile` currently hardcodes
  `GOARCH=amd64` + CGO; arm64 needs a cross-toolchain or buildx with QEMU plus
  Dockerfile changes.
- **Security scan:** CodeQL and/or `gosec` as a separate CI job.
- **`release.yml`:** optionally extra binaries (darwin/arm64) — but CGO+SQLite
  complicates cross-compilation (cross-CC per target).

---

## Long-term (plan v1.1 / 2.0)

- PostgreSQL adapter (the storage abstraction is already interface-based).
- Alert routing / severity-based / quiet hours / escalation (Phase 6.3).
- Multi-user + RBAC, API tokens / JWT.
- Maintenance windows, public status page, distributed monitoring.
- Data retention: rolling up old data (aggregation) instead of plain deletion.
