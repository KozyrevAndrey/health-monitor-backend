# Health Monitor - Implementation Progress

> Документ отражает **фактическое** состояние кода (сверено с git и исходниками).

---

## ✅ Completed

### Фаза 0-1: Инфраструктура и базовая архитектура
- [x] Go module, структура проекта, `.gitignore`/`.dockerignore`
- [x] Makefile (20+ команд), Dockerfile (multi-stage, CGO для SQLite), docker-compose
- [x] `.golangci.yml`, README.md
- [x] Domain Layer (`internal/domain/`): Target, CheckResult, Alert, Incident, NotifierConfig + интерфейсы (Checker, *Repository, Notifier, Scheduler, AlertManager, CheckerRegistry, NotifierRepository)
- [x] Configuration (`pkg/config/`): Viper, YAML, env override (явный binding env vars — `cc21b15`)
- [x] Logging (`pkg/logger/`): zerolog, уровни, json/console, stdout/stderr/file
- [x] Entry point (`cmd/server/main.go`): graceful shutdown, signal handling, lifecycle

### Фаза 2: Storage Layer
- [x] Миграции (`migrations/001_initial_schema.sql`), таблицы targets/check_results/incidents/alerts + 14 индексов
- [x] GORM + SQLite, AutoMigrate, connection pooling, auto-create db dir
- [x] Репозитории: TargetRepository, CheckResultRepository (история, stats, DeleteOlderThan), IncidentRepository
- [x] **NotifierRepository** + модель/миграция (`8155946`) — конфиги нотификаторов в БД

### Фаза 3-5: Checkers
- [x] Checker registry (thread-safe), резолв checker'а по `target.Type`
- [x] **HTTP/HTTPS** checker: status code, response time, SSL validation + expiry warning, custom headers/methods, redirects, skip-verify
- [x] **TCP** checker (`tcp.go`): проверка доступности порта + connection time, метадата host/port/address ✨
- [x] **DNS** checker (`dns.go`): резолвинг A/AAAA/CNAME/MX/TXT/NS, resolution time, кастомный DNS-сервер, валидация против expected values, дедуп записей ✨
- [x] Тесты (`http_test.go`, `tcp_test.go`, `dns_test.go`)
- [x] UI: форма таргета **type-aware** (HTTP url/status · TCP host/port · DNS domain/record/server/expected), endpoint в карточке

### Фаза 4: Scheduler
- [x] Ticker-based, независимые интервалы на target, concurrent execution (goroutine на target)
- [x] AddTarget/RemoveTarget/UpdateTarget, graceful shutdown (WaitGroup)
- [x] Загрузка targets, сохранение результатов, обработка через AlertManager

### Фаза 5-6: Alert System
- [x] AlertManager (`internal/alerting/manager.go`): state per target, consecutive fail/success, rule engine
- [x] Правила: consecutive failures, response-time threshold, SSL expiry, DOWN/UP transitions
- [x] Incident management (создание/резолв, last error)
- [x] Notifier registry + broadcast, интеграция со scheduler
- [x] **ClearNotifiers** для hot-reload (`f9cb8d2`)
- [x] Тесты (`internal/alerting/manager_test.go`)

### Фаза 7: Notification System
- [x] **Telegram** (`telegram.go`): Bot API, multiple chat IDs, Markdown, иконки/severity, **proxy (http/https/socks5)** ✨
- [x] **Email/SMTP** (`email.go`): auth, TLS/SSL (incl. implicit TLS на 465 — `f573485`), HTML+plaintext multipart
- [x] **Gmail Service Account** (`gmail.go`): domain-wide delegation, impersonation
- [x] **Gmail OAuth2** (`gmail_oauth.go`): личный аккаунт через credentials/token файлы
- [x] **Webhook (generic)** (`webhook.go`): url/method/custom headers/payload-шаблон (text/template), retry+backoff на 5xx/429, опц. proxy ✨
- [x] Тесты для telegram/email/webhook
- [x] **Notifier CRUD через REST API** (`1034d6f`) + hot-reload + маскирование секретов; нотификаторы управляются из БД/UI, не из YAML (`8518538`)

### Фаза 8: HTTP API
- [x] Chi router v5, конфигурируемые таймауты, graceful shutdown
- [x] Middleware: Request ID, Real IP, logging, recovery, timeout, CORS, **Basic Auth** (`fbdbe5a`)
- [x] Targets CRUD + results + stats; Incidents (list/get/ongoing); Notifiers CRUD; `GET /health`
- [x] **OpenAPI/Swagger**: `oapi-codegen` (`internal/generated/api.gen.go`), swagger handlers (`swagger_handlers.go`, `openapi_adapter.go`)

### Фаза 9: Web Dashboard
- [x] SPA (`web/static/index.html`), vanilla JS, адаптивная вёрстка
- [x] Stats cards, target cards со статусами, последние инциденты
- [x] **Полный CRUD targets и notifiers** через UI-модалки (`c274a8d`, `c776017`); per-type поля нотификаторов
- [x] Auto-refresh **polling каждые 10с** (не SSE), manual refresh
- [x] Static file server через Chi

### Deployment
- [x] **Production docker-compose** (`docker-compose.prod.yml`) + Traefik labels + `DOMAIN`/basic-auth env (`a3a9f9f`, `5d4d98a`, `578b240`)
- [x] `.env` / `.env.example`, mount secrets volume (`4bd3014`)

---

## ⏳ Not Done (по плану)

- [ ] **Дополнительные чекеры**: ICMP/Ping (есть http.go + tcp.go + dns.go)
- [ ] **SSE / real-time** (`GET /api/v1/events`) — сейчас UI на polling'е
- [ ] **Target detail page** с графиками (Chart.js): uptime 24h/7d/30d, история чекетов
- [ ] **Worker pool** в scheduler (сейчас goroutine-per-target без bounded pool/очереди)
- [ ] **Data retention job** (метод `DeleteOlderThan` есть, фонового cleanup нет)
- [ ] **Prometheus `/metrics`**
- [ ] **CLI**: `validate`, `backup`/`restore`
- [ ] **CI/CD** (GitHub Actions), load/security тесты
- [ ] PostgreSQL adapter, RBAC/multi-user, maintenance windows, status page, alert routing/escalation, quiet hours

---

## 📊 Current Status

| Компонент | % |
|---|---|
| Infrastructure / Domain / Config / Logging | 100% |
| Storage (SQLite + repos + notifier configs) | 100% |
| Checkers: HTTP + TCP + DNS | 100% |
| Scheduler (без worker pool) | 90% |
| Alert Manager (+ hot-reload) | 100% |
| Notifiers: Telegram(+proxy)/Email/Gmail SA/Gmail OAuth/Webhook | 100% |
| HTTP API (+ OpenAPI/Swagger, Basic Auth) | 100% |
| Web Dashboard (полный CRUD, polling) | 95% |
| Prod deploy (Traefik) | 100% |
| Доп. чекеры (ICMP) | 0% |
| SSE / графики / retention job / metrics | 0% |

**Last Updated:** 2026-06-26

---

## 🚀 Next Steps (приоритет)

1. **ICMP/Ping checker** (требует привилегий — CAP_NET_RAW)
2. **Retention cleanup job** — фоновая очистка по `cleanup_interval` (метод `DeleteOlderThan` готов)
3. **SSE real-time** + замена polling'а в UI
4. **Target detail page** с графиками (Chart.js)
5. **Worker pool** в scheduler
6. **Prometheus `/metrics`**, CLI `validate`/`backup`, CI/CD

---

## 📝 Notes
- Clean Architecture, interface-based DI, graceful shutdown везде
- Targets и notifiers управляются **только через API/UI** (убраны из YAML)
- Новый тип нотификатора подключается в 3 слоях: `internal/notifier/<type>.go` → два switch'а (`cmd/server/main.go`, `internal/api/notifier_handlers.go`) → UI (`web/static/index.html`)
- Конфиги нотификаторов валидируются через конструктор при create/update (битый url/шаблон/прокси → 400)
- Известные ограничения маскирования: значения внутри webhook `headers` (напр. `Authorization`) и креды в `proxy_url` (`socks5://user:pass@...`) **не маскируются** в ответах API — маскируются только top-level sensitive-ключи (`bot_token`, `password`, `token`, `secret`, `smtp_password`)
- Известное ограничение dispatch: AlertManager хранит нотификаторы в map по `Type()` (один на тип) и рассылает **последовательно** в рамках одного алерта (ограничено контекстом проверки `target.Timeout+5s`)
