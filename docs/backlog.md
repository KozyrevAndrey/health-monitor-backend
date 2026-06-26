# Backlog

Отложенные задачи (не в активной работе). Приоритетный список «что дальше» — в
[`implementation-progress.md`](implementation-progress.md) → Next Steps.

---

## ICMP / Ping checker

**Суть:** проверка живости хоста через ICMP Echo (как `ping`): доступность,
RTT (round-trip time), packet loss. Самая низкоуровневая проверка — без портов и
приложений. Полезно для роутеров/свитчей/серверов без открытого HTTP/TCP-сервиса.

**Почему отложено — требует привилегий.** Raw ICMP-сокеты нужны root / Linux
capability `CAP_NET_RAW`, либо unprivileged-режим через UDP-сокет (`ip4:icmp`) с
разрешением в sysctl `net.ipv4.ping_group_range`. В Docker — `cap_add: [NET_RAW]`
или `sysctls`. То есть, в отличие от TCP/DNS, «просто заработает» не везде —
нужны правки деплоя.

**Как реализовать (тот же паттерн, что TCP/DNS):**
- `internal/checker/icmp.go`, зарегистрировать в `NewDefaultRegistry()`.
- `ICMPConfig` уже есть в domain: `host`, `packet_count`, `max_rtt_ms`.
- Библиотека: `github.com/prometheus-community/pro-bing` (умеет privileged и
  unprivileged режимы).
- Результат: `success` если потерь нет и RTT в норме; `warning` при потерях или
  превышении `max_rtt_ms`; `failure` если хост не ответил.
- Добавить ICMP-поля в type-aware форму targets (`web/static/index.html`).
- Документировать в `docs/checkers.md` требование `NET_RAW`/sysctl + пример Docker.

---

## Lint cleanup → сделать lint блокирующим

**Сейчас:** `lint` job в CI помечен `continue-on-error: true` (advisory) — репортит,
но не валит CI. Причина: строгий `.golangci.yml` выдаёт **~201 finding** на код,
который раньше не линтился.

**Разбивка findings (v1.64.8):** revive 79 (в основном «exported … should have
comment»), errcheck 34, gocritic 27, errorlint 7, unused 6, govet 6, unparam 5,
gosec 4, goconst 4, gosimple 3, misspell 2, gocyclo 2, contextcheck 2.

**Как закрыть (на выбор):**
- *Постепенно по категориям:* сначала потенциально-реальные (errcheck, unused,
  gosec, govet, contextcheck), отложить чистую стилистику (revive `exported`,
  gocritic `opinionated/experimental`). По мере зачистки — снять `continue-on-error`.
- *Смягчить конфиг:* убрать самые педантичные правила (revive `exported`,
  gocritic-теги `opinionated`/`experimental`) → останется ~два десятка, их добить и
  сделать lint блокирующим.

**Файлы:** `.golangci.yml`, `.github/workflows/test.yml` (флаг `continue-on-error`).

---

## Worker pool в scheduler (+ jitter / backoff)

**Сейчас:** `internal/scheduler/scheduler.go` заводит **горутину-на-таргет** с
собственным тикером; глобального ограничения на число одновременных проверок нет.
Для десятков-сотен таргетов это нормально.

**Зачем:** при сотнях-тысячах таргетов или совпадении тиков может стартовать
много проверок разом → всплеск по сети/CPU/FD и нагрузке на запись в SQLite.
Пул отвязывает «сколько таргетов» от «сколько проверок бежит одновременно».

**Как реализовать:**
- Тикеры кладут задачу «проверь X» в канал-очередь; N воркеров разбирают и
  выполняют (конфигурируемое число воркеров).
- Добавить **jitter** при старте (размазать первые проверки во времени — против
  thundering herd) и **backoff** при сериях ошибок (увеличивать интервал, возврат
  к норме при success). См. план, Фаза 5.1–5.2.
- Опционально метрики: queue depth, worker utilization (см. Prometheus ниже).

---

## Prometheus `/metrics`

**Зачем:** интеграция с существующим стеком Prometheus + Grafana + Alertmanager.
Брать **только если такой стек есть/планируется** — иначе встроенный дашборд (SSE,
графики), статистика (uptime/response) и нотификаторы уже закрывают потребности.

**Два среза:**
- *Про сам монитор:* активные таргеты, queue depth, длительность проверок
  P50/P95/P99, error rate, goroutines, пул соединений к БД (Фаза 9.1).
- *Данные мониторинга как метрики:* `health_target_up{target=…}`,
  `health_target_response_ms{…}` → дашборды/алерты в Grafana/Alertmanager.

**Как:** `github.com/prometheus/client_golang`, отдать `/metrics` (обычно вне
basic-auth), инкрементить в тех же точках scheduler/alert-manager, куда уже
вешаются SSE-события.

---

## CLI: `validate` и `backup`/`restore`

**Сейчас:** бинарь принимает только `--config` и `--version` (flag-based).

- **`validate`** (Фаза 10.2) — проверить конфиг без запуска: валидный YAML,
  обязательные поля, корректные интервалы/таймауты, путь к БД. `validate --config
  config.yaml` → exit 0/1. Дёшево, удобно в CI/перед деплоем.
- **`backup` / `restore`** (Фаза 10.4) — корректный бэкап SQLite «на лету»
  (`VACUUM INTO` / Online Backup API, не `cp`). Польза средняя — часто хватает
  снапшота тома/файла.

Для подкоманд — либо `cobra`/`urfave/cli` (рефактор `main.go`), либо ручной разбор
`os.Args[1]`.

---

## CI/CD — расширения

- **Multi-arch Docker** (arm64): сейчас `Dockerfile` хардкодит `GOARCH=amd64` +
  CGO; для arm64 нужен cross-toolchain или buildx с QEMU и правки Dockerfile.
- **Security-скан:** CodeQL и/или `gosec` отдельным job в CI.
- **`release.yml`:** при желании — доп. бинари (darwin/arm64) — но CGO+SQLite
  усложняет кросс-сборку (нужны cross-CC per target).

---

## Долгосрочное (план v1.1 / 2.0)

- PostgreSQL adapter (абстракция storage уже на интерфейсах).
- Alert routing / severity-based / quiet hours / escalation (Фаза 6.3).
- Multi-user + RBAC, API tokens / JWT.
- Maintenance windows, публичная status page, distributed monitoring.
- Data retention: агрегация старых данных (rollup) вместо простого удаления.
