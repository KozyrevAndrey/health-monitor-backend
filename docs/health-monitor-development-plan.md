# Health Monitor - План разработки

**Версия:** 1.0  
**Дата:** Декабрь 2024  
**Цель:** Легковесная, расширяемая система мониторинга доменов, хостов и health endpoints на Go

---

## Оглавление

1. [Обзор проекта](#обзор-проекта)
2. [Технологический стек](#технологический-стек)
3. [Архитектурные принципы](#архитектурные-принципы)
4. [Фазы разработки](#фазы-разработки)
5. [Структура проекта](#структура-проекта)
6. [Roadmap и приоритеты](#roadmap-и-приоритеты)
7. [Deployment и инфраструктура](#deployment-и-инфраструктура)

---

## Обзор проекта

### Назначение
Self-hosted система мониторинга для отслеживания доступности и производительности веб-сервисов, API endpoints и сетевых ресурсов.

### Ключевые особенности
- ✅ Легковесность (один Go бинарник)
- ✅ Расширяемая архитектура
- ✅ Поддержка множества типов проверок
- ✅ Гибкая система уведомлений
- ✅ Простой веб-интерфейс
- ✅ Real-time обновления
- ✅ Low resource footprint

### Целевое использование
- Мониторинг production сервисов
- Проверка доступности API
- Отслеживание SSL сертификатов
- Мониторинг внутренних сервисов
- SLA tracking

---

## Технологический стек

### Backend
- **Язык:** Go 1.21+
- **HTTP Router:** chi/gin (выбрать на этапе реализации)
- **БД:** SQLite (с возможностью миграции на PostgreSQL)
- **ORM:** GORM
- **Scheduler:** robfig/cron или custom ticker-based
- **Логирование:** zerolog/zap

### Frontend
- **Framework:** HTMX + Alpine.js
- **Стили:** Pico.css / Water.css
- **Real-time:** Server-Sent Events (SSE)
- **Графики:** Chart.js (опционально)

### Инфраструктура
- **Контейнеризация:** Docker
- **Оркестрация:** Docker Compose (Kubernetes опционально)
- **CI/CD:** GitHub Actions

---

## Архитектурные принципы

### 1. Чистая архитектура (Clean Architecture)

```
┌─────────────────────────────────────┐
│         Presentation Layer          │
│     (HTTP Handlers, Templates)      │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│        Application Layer            │
│    (Use Cases, Business Logic)      │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│          Domain Layer               │
│   (Entities, Interfaces, Rules)     │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│      Infrastructure Layer           │
│  (DB, External Services, Checkers)  │
└─────────────────────────────────────┘
```

### 2. Dependency Injection
- Все зависимости передаются через конструкторы
- Использование интерфейсов вместо конкретных типов
- Легкое тестирование через моки

### 3. Interface Segregation
Ключевые интерфейсы:

```go
// Примеры интерфейсов (не код, а концепция)

type Checker interface {
    Check(ctx context.Context, target Target) CheckResult
    Type() string
}

type Storage interface {
    SaveTarget(target Target) error
    GetTarget(id string) (Target, error)
    SaveCheckResult(result CheckResult) error
    GetHistory(targetID string, limit int) ([]CheckResult, error)
}

type Notifier interface {
    Notify(ctx context.Context, alert Alert) error
    Type() string
}

type Scheduler interface {
    AddTarget(target Target, interval time.Duration)
    RemoveTarget(targetID string)
    Start() error
    Stop() error
}
```

### 4. Extensibility Points
- **Checker Registry:** Регистрация новых типов проверок
- **Notifier Registry:** Добавление каналов уведомлений
- **Storage Adapters:** Поддержка разных БД
- **Plugin System:** (опционально) Загрузка внешних плагинов

---

## Фазы разработки

## Фаза 1: Проектирование архитектуры

**Цель:** Спроектировать расширяемую и поддерживаемую архитектуру

### 1.1 Определение интерфейсов

**Задачи:**
- [ ] Определить интерфейс `Checker` с методами:
  - `Check(ctx, target) -> result`
  - `Type() -> string`
  - `Validate(config) -> error`
- [ ] Определить интерфейс `Storage`:
  - CRUD для targets
  - Сохранение результатов проверок
  - Получение истории и статистики
- [ ] Определить интерфейс `Notifier`:
  - `Notify(ctx, alert) -> error`
  - `Type() -> string`
- [ ] Определить интерфейс `Scheduler`:
  - Управление задачами
  - Lifecycle методы (Start/Stop)

**Результат:** Документ с описанием всех интерфейсов

### 1.2 Структура конфигурации

**Задачи:**
- [ ] Спроектировать YAML схему конфигурации
- [ ] Определить структуру Target:
  ```yaml
  targets:
    - id: "api-production"
      name: "Production API"
      type: "http"
      url: "https://api.example.com/health"
      interval: "30s"
      timeout: "10s"
      checks:
        - type: "status_code"
          expected: 200
        - type: "response_time"
          max_ms: 500
      alerts:
        - type: "webhook"
          url: "https://hooks.slack.com/..."
  ```
- [ ] Определить глобальные настройки
- [ ] Спроектировать retention policies

**Результат:** Примеры конфигурационных файлов

### 1.3 Модели данных

**Задачи:**
- [ ] Описать структуру `Target`
- [ ] Описать структуру `CheckResult`:
  - Timestamp
  - Success/Failure
  - Response time
  - Status code
  - Error message
  - Metadata (headers, body snippet)
- [ ] Описать структуру `Alert`
- [ ] Описать структуру `Incident`:
  - Группировка последовательных failures
  - Start/End time
  - Status (ongoing/resolved)

**Результат:** Domain модели в виде Go structs (концептуально)

---

## Фаза 2: Базовая инфраструктура

**Цель:** Настроить проект и базовую инфраструктуру

### 2.1 Настройка проекта

**Задачи:**
- [ ] Инициализировать Go модуль
  - `go mod init github.com/username/health-monitor`
- [ ] Создать структуру директорий:
  ```
  health-monitor/
  ├── cmd/
  │   └── server/
  │       └── main.go
  ├── internal/
  │   ├── domain/
  │   ├── checker/
  │   ├── storage/
  │   ├── notifier/
  │   ├── scheduler/
  │   ├── alerting/
  │   └── api/
  ├── pkg/
  │   └── config/
  ├── web/
  │   ├── templates/
  │   └── static/
  ├── migrations/
  ├── configs/
  │   └── example.yaml
  ├── scripts/
  ├── Makefile
  ├── Dockerfile
  ├── docker-compose.yml
  └── README.md
  ```
- [ ] Настроить golangci-lint
- [ ] Создать Makefile с командами:
  - `make build`
  - `make test`
  - `make lint`
  - `make run`
  - `make docker-build`

**Результат:** Работающий скелет проекта

### 2.2 Система конфигурации

**Задачи:**
- [ ] Выбрать библиотеку (viper / cleanenv / custom)
- [ ] Реализовать парсинг YAML
- [ ] Добавить валидацию конфигурации
- [ ] Поддержка env переменных для чувствительных данных
- [ ] Hot-reload конфига (опционально для v2)

**Результат:** Рабочая загрузка конфигурации

### 2.3 Логирование

**Задачи:**
- [ ] Выбрать библиотеку (zerolog рекомендуется)
- [ ] Настроить structured logging
- [ ] Конфигурируемые уровни логов
- [ ] Логирование в файл + stdout
- [ ] Ротация логов (если в файл)

**Результат:** Централизованное логирование

---

## Фаза 3: Storage слой

**Цель:** Реализовать персистентность данных

### 3.1 Интерфейс репозитория

**Задачи:**
- [ ] Определить интерфейс `TargetRepository`:
  - `Create(target) error`
  - `Get(id) (Target, error)`
  - `List() ([]Target, error)`
  - `Update(target) error`
  - `Delete(id) error`
- [ ] Определить интерфейс `CheckResultRepository`:
  - `Save(result) error`
  - `GetHistory(targetID, limit, offset) ([]Result, error)`
  - `GetLatest(targetID) (Result, error)`
  - `GetStats(targetID, from, to) (Stats, error)`
- [ ] Определить интерфейс `IncidentRepository`

**Результат:** Описание всех методов репозиториев

### 3.2 SQLite имплементация

**Задачи:**
- [ ] Спроектировать схему БД:
  ```sql
  -- Пример структуры
  CREATE TABLE targets (
      id TEXT PRIMARY KEY,
      name TEXT NOT NULL,
      type TEXT NOT NULL,
      config JSON NOT NULL,
      interval INTEGER NOT NULL,
      enabled BOOLEAN DEFAULT true,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );
  
  CREATE TABLE check_results (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      target_id TEXT NOT NULL,
      success BOOLEAN NOT NULL,
      response_time_ms INTEGER,
      status_code INTEGER,
      error TEXT,
      checked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      FOREIGN KEY (target_id) REFERENCES targets(id)
  );
  
  CREATE INDEX idx_check_results_target_time 
      ON check_results(target_id, checked_at DESC);
  ```
- [ ] Настроить миграции (golang-migrate или goose)
- [ ] Реализовать все методы репозитория
- [ ] Добавить индексы для производительности
- [ ] Connection pooling
- [ ] Транзакции где необходимо

**Результат:** Рабочий SQLite storage

### 3.3 Подготовка к масштабированию

**Задачи:**
- [ ] Абстрагировать SQL-специфичный код
- [ ] Подготовить интерфейс для PostgreSQL
- [ ] Graceful shutdown с flush данных
- [ ] Автоматическая очистка старых данных (retention)

**Результат:** Возможность переключения БД

---

## Фаза 4: Checker система

**Цель:** Реализовать различные типы проверок

### 4.1 Базовый интерфейс

**Задачи:**
- [ ] Реализовать базовый интерфейс `Checker`
- [ ] Добавить timeout handling
- [ ] Реализовать context propagation
- [ ] Опциональная retry логика с backoff
- [ ] Метрики производительности checkers

**Результат:** Базовый механизм проверок

### 4.2 HTTP Checker (приоритет 1)

**Задачи:**
- [ ] Проверка HTTP status code
- [ ] Измерение response time
- [ ] Валидация response body:
  - Проверка наличия строки
  - Regex matching
  - JSON path validation
- [ ] Проверка headers
- [ ] Проверка SSL сертификатов:
  - Валидность
  - Срок истечения (предупреждение за N дней)
- [ ] Follow redirects (конфигурируемо)
- [ ] Custom HTTP методы (GET, POST, HEAD)
- [ ] Custom headers и authentication

**Результат:** Полнофункциональный HTTP checker

### 4.3 TCP Checker (приоритет 2)

**Задачи:**
- [ ] Проверка доступности порта
- [ ] Измерение connection time
- [ ] Опциональная отправка/получение данных
- [ ] Поддержка TLS

**Результат:** TCP port checker

### 4.4 DNS Checker (приоритет 3)

**Задачи:**
- [ ] Резолвинг A/AAAA записей
- [ ] Проверка MX, TXT, CNAME записей
- [ ] Измерение DNS resolution time
- [ ] Проверка против ожидаемых значений

**Результат:** DNS checker

### 4.5 ICMP/Ping Checker (приоритет 4)

**Задачи:**
- [ ] ICMP ping
- [ ] Измерение RTT
- [ ] Packet loss detection
- [ ] Требует привилегий - документировать

**Результат:** Ping checker

### 4.6 Custom Script Checker (опционально)

**Задачи:**
- [ ] Выполнение внешних скриптов
- [ ] Парсинг exit code
- [ ] Timeout для скриптов
- [ ] Безопасность (sandboxing)

**Результат:** Расширение через скрипты

### 4.7 Checker Registry

**Задачи:**
- [ ] Factory pattern для создания checkers
- [ ] Регистрация новых типов
- [ ] Валидация конфигурации по типу

**Результат:** Легко расширяемая система

---

## Фаза 5: Scheduler

**Цель:** Планирование и выполнение проверок

### 5.1 Менеджер задач

**Задачи:**
- [ ] Ticker-based scheduler для каждого target
- [ ] Worker pool для параллельного выполнения:
  - Конфигурируемое количество workers
  - Queue для задач
  - Graceful worker shutdown
- [ ] Управление lifecycle:
  - `Start()` - запуск планировщика
  - `Stop()` - graceful остановка
  - `AddTarget()` - добавление нового target
  - `RemoveTarget()` - удаление target
  - `UpdateTarget()` - обновление интервала
- [ ] Context propagation для отмены

**Результат:** Базовый scheduler

### 5.2 Оптимизации

**Задачи:**
- [ ] Jitter для избежания thundering herd:
  - Случайная задержка при старте
  - Распределение проверок во времени
- [ ] Backoff при последовательных failures:
  - Увеличение интервала при ошибках
  - Возврат к нормальному интервалу при success
- [ ] Приоритизация проверок:
  - Critical targets проверяются чаще
  - Lower priority при высокой нагрузке
- [ ] Health check самого scheduler

**Результат:** Оптимизированный scheduler

### 5.3 Метрики scheduler

**Задачи:**
- [ ] Количество активных задач
- [ ] Queue depth
- [ ] Worker utilization
- [ ] Avg check duration

**Результат:** Наблюдаемость scheduler

---

## Фаза 6: Alerting система

**Цель:** Уведомления о проблемах

### 6.1 Alert Manager

**Задачи:**
- [ ] State machine для алертов:
  ```
  OK → WARNING → CRITICAL → RECOVERY → OK
  ```
- [ ] Определение условий:
  - N последовательных failures
  - Response time превышает threshold
  - SSL сертификат истекает через N дней
- [ ] Debouncing:
  - Не отправлять алерт при каждом flap
  - Grace period перед алертом
  - Recovery notification
- [ ] Группировка в incidents:
  - Один incident для серии failures
  - Tracking времени downtime
  - Resolution tracking

**Результат:** Умная система алертинга

### 6.2 Notifier реализации

**Приоритет 1: Webhook Notifier**
- [ ] POST запрос на URL
- [ ] Конфигурируемый payload template
- [ ] Retry логика
- [ ] Custom headers

**Приоритет 2: Email Notifier**
- [ ] SMTP отправка
- [ ] HTML templates для писем
- [ ] Attachment support (графики, опционально)

**Приоритет 3: Telegram Notifier**
- [ ] Bot API интеграция
- [ ] Форматированные сообщения
- [ ] Multiple chat IDs

**Опционально:**
- [ ] Slack notifier
- [ ] Discord notifier
- [ ] PagerDuty integration
- [ ] Custom webhook formats (Slack, Discord)

**Результат:** Множественные каналы уведомлений

### 6.3 Alert Routing

**Задачи:**
- [ ] Разные notifiers для разных targets
- [ ] Severity-based routing:
  - WARNING → email
  - CRITICAL → email + telegram + webhook
- [ ] Quiet hours (не беспокоить ночью):
  - Конфигурируемые временные окна
  - Timezone support
- [ ] Escalation policies (опционально):
  - Level 1 → email после 5 мин
  - Level 2 → telegram после 15 мин
  - Level 3 → phone call после 30 мин

**Результат:** Гибкая система роутинга

### 6.4 Alert Templates

**Задачи:**
- [ ] Template engine для сообщений
- [ ] Доступные переменные:
  - Target name, URL
  - Error message
  - Response time
  - Timestamp
  - Downtime duration
- [ ] Разные templates для разных событий:
  - DOWN alert
  - RECOVERY alert
  - WARNING alert
  - SSL expiration

**Результат:** Настраиваемые сообщения

---

## Фаза 7: API Layer

**Цель:** REST API для управления и данных

### 7.1 REST API Endpoints

**Targets Management:**
- [ ] `GET /api/v1/targets` - список всех targets
  - Фильтрация по типу, статусу
  - Pagination
  - Sorting
- [ ] `GET /api/v1/targets/:id` - детали target
- [ ] `GET /api/v1/targets/:id/status` - текущий статус
- [ ] `POST /api/v1/targets` - создание (опционально для MVP)
- [ ] `PUT /api/v1/targets/:id` - обновление
- [ ] `DELETE /api/v1/targets/:id` - удаление

**Check History:**
- [ ] `GET /api/v1/targets/:id/checks` - история проверок
  - Pagination
  - Date range filtering
  - Success/failure filtering
- [ ] `GET /api/v1/targets/:id/stats` - статистика
  - Uptime percentage
  - Avg response time
  - Success rate
  - Time range (24h, 7d, 30d)

**Incidents:**
- [ ] `GET /api/v1/incidents` - список incidents
- [ ] `GET /api/v1/incidents/:id` - детали incident

**System:**
- [ ] `GET /api/health` - health check самого монитора
- [ ] `GET /api/metrics` - метрики системы (Prometheus format опционально)

**Результат:** Полнофункциональный API

### 7.2 Real-time Updates

**Задачи:**
- [ ] Server-Sent Events (SSE) endpoint:
  - `GET /api/v1/events` - stream событий
- [ ] События для отправки:
  - Check completed
  - Status changed
  - Alert triggered
  - Alert resolved
- [ ] Client reconnection handling
- [ ] Event filtering по target ID

**Результат:** Real-time уведомления

### 7.3 Middleware

**Задачи:**
- [ ] Logging middleware:
  - Request/response logging
  - Duration tracking
- [ ] Recovery middleware (panic recovery)
- [ ] CORS middleware (если нужно)
- [ ] Rate limiting (опционально)
- [ ] Authentication middleware:
  - Basic Auth для начала
  - API keys (опционально)
  - JWT (для production)
- [ ] Request ID для трейсинга

**Результат:** Production-ready API

### 7.4 API Documentation

**Задачи:**
- [ ] OpenAPI/Swagger спецификация
- [ ] Примеры запросов/ответов
- [ ] Swagger UI (опционально)

**Результат:** Документированный API

---

## Фаза 8: Web UI Integration

**Цель:** Веб-интерфейс для мониторинга

### 8.1 Статические файлы

**Задачи:**
- [ ] Embed HTML/CSS/JS в бинарник (`go:embed`)
- [ ] Serving статики через HTTP handler
- [ ] Cache headers для статики
- [ ] Compression (gzip)

**Результат:** Single binary deployment

### 8.2 Template Rendering

**Задачи:**
- [ ] HTML templates (`html/template`)
- [ ] Layout + partials
- [ ] Передача данных в templates
- [ ] HTMX integration:
  - Partial updates
  - Auto-refresh списка targets
  - Inline editing (опционально)

**Результат:** Server-side rendering

### 8.3 Dashboard Components

**Страницы:**

**1. Overview Dashboard:**
- [ ] Карточки со статистикой:
  - Всего targets
  - Healthy / Unhealthy
  - Active incidents
- [ ] Список всех targets с текущим статусом
- [ ] Фильтрация и поиск
- [ ] Сортировка по статусу, имени, response time

**2. Target Detail Page:**
- [ ] Текущий статус (большой индикатор)
- [ ] Последние N проверок
- [ ] График response time (Chart.js)
- [ ] Uptime за разные периоды (24h, 7d, 30d)
- [ ] История incidents
- [ ] Конфигурация target

**3. Incidents Page:**
- [ ] Список активных incidents
- [ ] История resolved incidents
- [ ] Timeline визуализация

**4. Settings Page (опционально):**
- [ ] Управление targets (CRUD)
- [ ] Конфигурация notifiers
- [ ] Системные настройки

**Результат:** Полнофункциональный UI

### 8.4 Real-time Updates в UI

**Задачи:**
- [ ] SSE connection из браузера
- [ ] Автообновление статусов без перезагрузки
- [ ] Уведомления в браузере (опционально)
- [ ] Visual indicators для изменений

**Результат:** Live dashboard

### 8.5 Responsive Design

**Задачи:**
- [ ] Mobile-friendly layout
- [ ] Адаптивные таблицы
- [ ] Touch-friendly кнопки

**Результат:** Работает на всех устройствах

---

## Фаза 9: Observability

**Цель:** Мониторинг самого монитора

### 9.1 Internal Metrics

**Задачи:**
- [ ] Метрики о самой системе:
  - Количество активных targets
  - Количество выполняемых проверок
  - Queue depth
  - Worker utilization
  - БД connection pool usage
  - Memory usage
  - Goroutines count
- [ ] Prometheus metrics endpoint (опционально):
  - `GET /metrics` в Prometheus формате
- [ ] Метрики производительности:
  - Avg check duration по типам
  - P50/P95/P99 response times
  - Error rates

**Результат:** Наблюдаемость системы

### 9.2 Health Check

**Задачи:**
- [ ] Endpoint `GET /api/health`
- [ ] Проверки:
  - БД connection alive
  - Scheduler работает
  - Все workers живы
  - Disk space available
- [ ] Health status:
  - `healthy` - все OK
  - `degraded` - частичные проблемы
  - `unhealthy` - критические проблемы
- [ ] Детализация проблем в response

**Результат:** Self-monitoring

### 9.3 Structured Events

**Задачи:**
- [ ] Event log для важных событий:
  - System started/stopped
  - Target added/removed
  - Alert triggered/resolved
  - Configuration reloaded
- [ ] Audit log (опционально):
  - API calls
  - Configuration changes
  - User actions

**Результат:** Полная история событий

---

## Фаза 10: Production-Ready Features

**Цель:** Подготовка к production использованию

### 10.1 Graceful Shutdown

**Задачи:**
- [ ] Обработка сигналов:
  - SIGTERM
  - SIGINT
- [ ] Shutdown sequence:
  1. Остановка приема новых HTTP запросов
  2. Завершение активных HTTP запросов (с timeout)
  3. Остановка scheduler (no new checks)
  4. Завершение активных проверок (с timeout)
  5. Flush данных в БД
  6. Закрытие БД connections
  7. Закрытие notifiers
- [ ] Graceful timeout (30 секунд по умолчанию)
- [ ] Логирование shutdown process

**Результат:** Безопасная остановка

### 10.2 Configuration Validation

**Задачи:**
- [ ] Валидация при старте:
  - Корректность YAML
  - Обязательные поля
  - Валидные интервалы (>0)
  - Валидные URLs
- [ ] Validate команда:
  - `health-monitor validate config.yaml`
- [ ] Подробные сообщения об ошибках

**Результат:** Предотвращение ошибок конфигурации

### 10.3 Data Retention

**Задачи:**
- [ ] Автоматическая очистка старых данных:
  - Конфигурируемый retention period
  - Агрегация старых данных (опционально)
- [ ] Background job для cleanup
- [ ] Логирование удаленных записей

**Результат:** Контроль размера БД

### 10.4 Backup & Restore

**Задачи:**
- [ ] Backup команда:
  - `health-monitor backup --output backup.db`
  - Резервное копирование SQLite
- [ ] Restore команда
- [ ] Документация процесса

**Результат:** Защита данных

### 10.5 Security

**Задачи:**
- [ ] Authentication:
  - Basic Auth для веб-интерфейса
  - API keys для API
- [ ] HTTPS support:
  - Конфигурация TLS
  - Let's Encrypt integration (опционально)
- [ ] Rate limiting на API
- [ ] Input validation для всех endpoints
- [ ] SQL injection prevention (через ORM)
- [ ] XSS prevention в templates

**Результат:** Безопасное приложение

### 10.6 Performance

**Задачи:**
- [ ] Connection pooling для БД
- [ ] HTTP client pooling для checkers
- [ ] Кэширование где возможно:
  - Static files
  - Config (в памяти)
- [ ] Индексы БД
- [ ] Pagination для больших списков

**Результат:** Оптимизированная производительность

---

## Структура проекта

### Финальная структура директорий

```
health-monitor/
├── cmd/
│   └── server/
│       └── main.go                 # Entry point
│
├── internal/
│   ├── domain/                     # Domain models & interfaces
│   │   ├── target.go
│   │   ├── check_result.go
│   │   ├── alert.go
│   │   ├── incident.go
│   │   └── interfaces.go           # All interfaces
│   │
│   ├── checker/                    # Checker implementations
│   │   ├── checker.go              # Base interface
│   │   ├── http.go                 # HTTP checker
│   │   ├── tcp.go                  # TCP checker
│   │   ├── dns.go                  # DNS checker
│   │   ├── ping.go                 # ICMP checker
│   │   └── registry.go             # Checker registry
│   │
│   ├── storage/                    # Storage implementations
│   │   ├── storage.go              # Interfaces
│   │   ├── sqlite/
│   │   │   ├── sqlite.go
│   │   │   ├── target_repo.go
│   │   │   ├── result_repo.go
│   │   │   └── incident_repo.go
│   │   └── postgres/               # Future PostgreSQL
│   │
│   ├── scheduler/                  # Scheduling logic
│   │   ├── scheduler.go
│   │   ├── worker_pool.go
│   │   └── task.go
│   │
│   ├── alerting/                   # Alert management
│   │   ├── manager.go
│   │   ├── state_machine.go
│   │   └── rules.go
│   │
│   ├── notifier/                   # Notifier implementations
│   │   ├── notifier.go             # Interface
│   │   ├── webhook.go
│   │   ├── email.go
│   │   ├── telegram.go
│   │   └── registry.go
│   │
│   ├── api/                        # HTTP API
│   │   ├── server.go               # HTTP server setup
│   │   ├── middleware/
│   │   │   ├── logging.go
│   │   │   ├── recovery.go
│   │   │   └── auth.go
│   │   ├── handlers/
│   │   │   ├── targets.go
│   │   │   ├── checks.go
│   │   │   ├── incidents.go
│   │   │   ├── health.go
│   │   │   └── events.go           # SSE
│   │   └── routes.go
│   │
│   └── app/                        # Application orchestration
│       ├── app.go                  # Main app struct
│       └── lifecycle.go            # Start/Stop logic
│
├── pkg/                            # Public packages
│   ├── config/
│   │   ├── config.go               # Config structs
│   │   ├── loader.go               # YAML loading
│   │   └── validator.go            # Validation
│   │
│   └── logger/
│       └── logger.go               # Logger setup
│
├── web/                            # Frontend
│   ├── templates/
│   │   ├── layouts/
│   │   │   └── base.html
│   │   ├── pages/
│   │   │   ├── dashboard.html
│   │   │   ├── target_detail.html
│   │   │   └── incidents.html
│   │   └── partials/
│   │       ├── target_card.html
│   │       └── status_badge.html
│   │
│   └── static/
│       ├── css/
│       │   └── style.css
│       └── js/
│           ├── htmx.min.js
│           ├── alpine.min.js
│           └── app.js
│
├── migrations/                     # Database migrations
│   ├── 001_initial.up.sql
│   ├── 001_initial.down.sql
│   ├── 002_add_incidents.up.sql
│   └── 002_add_incidents.down.sql
│
├── configs/                        # Example configs
│   ├── example.yaml
│   ├── production.yaml.example
│   └── development.yaml
│
├── scripts/                        # Utility scripts
│   ├── build.sh
│   └── migrate.sh
│
├── docs/                           # Documentation
│   ├── architecture.md
│   ├── api.md
│   ├── configuration.md
│   └── deployment.md
│
├── .github/
│   └── workflows/
│       ├── test.yml
│       └── release.yml
│
├── Makefile
├── Dockerfile
├── docker-compose.yml
├── .dockerignore
├── .gitignore
├── .golangci.yml
├── go.mod
├── go.sum
└── README.md
```

---

## Roadmap и приоритеты

### MVP (Минимально жизнеспособный продукт)

**Цель:** Работающий монитор с базовым функционалом

**Включает:**
- ✅ HTTP checker (status code, response time)
- ✅ SQLite storage
- ✅ Простой scheduler (без worker pool)
- ✅ Webhook notifier
- ✅ Базовый REST API (read-only)
- ✅ Простой веб-интерфейс (список targets + детали)
- ✅ Конфигурация через YAML

**Не включает:**
- Дополнительные checkers (TCP, DNS, ping)
- Дополнительные notifiers (email, telegram)
- CRUD операции в UI
- Real-time updates
- Advanced alerting rules

**Срок:** 2-3 недели разработки

---

### Version 1.0

**Цель:** Production-ready решение

**Добавляет к MVP:**
- ✅ Worker pool в scheduler
- ✅ TCP и DNS checkers
- ✅ Email и Telegram notifiers
- ✅ Incident management
- ✅ Real-time updates (SSE)
- ✅ Улучшенный UI с графиками
- ✅ Graceful shutdown
- ✅ Authentication
- ✅ Metrics endpoint

**Срок:** +2-3 недели после MVP

---

### Version 1.1

**Цель:** Расширенные возможности

**Добавляет:**
- ✅ CRUD операции в UI
- ✅ PostgreSQL support
- ✅ SSL certificate monitoring
- ✅ Advanced alerting rules
- ✅ Alert routing и escalation
- ✅ Data retention policies
- ✅ Prometheus metrics export

**Срок:** +2 недели после v1.0

---

### Version 2.0 (Future)

**Идеи для будущего:**
- 🔮 Multi-user support с permissions
- 🔮 Incident timeline visualization
- 🔮 SLA tracking и reporting
- 🔮 Status page generation (публичная страница статуса)
- 🔮 Distributed monitoring (multiple instances)
- 🔮 Plugin system для custom checkers
- 🔮 Mobile app
- 🔮 Advanced analytics и ML для anomaly detection

---

## Deployment и инфраструктура

### Docker Deployment

**Dockerfile (multi-stage):**

```dockerfile
# Концептуальный пример структуры

# Stage 1: Build
FROM golang:1.21-alpine AS builder
# ... build steps

# Stage 2: Runtime
FROM alpine:latest
# ... copy binary
# ... setup non-root user
ENTRYPOINT ["/app/health-monitor"]
```

**docker-compose.yml:**

```yaml
# Концептуальный пример

version: '3.8'
services:
  health-monitor:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./config.yaml:/etc/health-monitor/config.yaml
      - ./data:/var/lib/health-monitor
    environment:
      - LOG_LEVEL=info
    restart: unless-stopped
```

### Systemd Service

```ini
# Пример systemd unit файла

[Unit]
Description=Health Monitor Service
After=network.target

[Service]
Type=simple
User=health-monitor
ExecStart=/usr/local/bin/health-monitor --config /etc/health-monitor/config.yaml
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

### Kubernetes Deployment (опционально)

```yaml
# Концептуальный пример

apiVersion: apps/v1
kind: Deployment
metadata:
  name: health-monitor
spec:
  replicas: 1
  template:
    spec:
      containers:
      - name: health-monitor
        image: health-monitor:latest
        ports:
        - containerPort: 8080
        volumeMounts:
        - name: config
          mountPath: /etc/health-monitor
        - name: data
          mountPath: /var/lib/health-monitor
```

---

## Тестирование

### Unit тесты

**Покрытие:**
- [ ] Все checkers (моки для HTTP клиента)
- [ ] Storage repositories (in-memory или testcontainers)
- [ ] Alert manager state machine
- [ ] Scheduler логика
- [ ] API handlers (HTTP тесты)

**Инструменты:**
- `testing` package
- `testify` для assertions
- `gomock` для моков
- `httptest` для HTTP тестов

### Integration тесты

**Сценарии:**
- [ ] End-to-end flow: конфиг → проверка → сохранение → алерт
- [ ] API integration тесты
- [ ] БД миграции
- [ ] Real checkers против test servers

### Load тесты (опционально)

**Проверить:**
- [ ] Производительность при 1000+ targets
- [ ] Memory leaks
- [ ] Goroutine leaks
- [ ] БД performance

---

## Документация

### Необходимая документация

**README.md:**
- Описание проекта
- Quick start guide
- Installation instructions
- Basic usage examples

**docs/configuration.md:**
- Полное описание всех конфигурационных опций
- Примеры для разных use cases
- Best practices

**docs/api.md:**
- API endpoints документация
- Request/response примеры
- Authentication

**docs/architecture.md:**
- Архитектурная диаграмма
- Описание компонентов
- Design decisions

**docs/deployment.md:**
- Docker deployment
- Systemd setup
- Kubernetes deployment
- Backup/restore процедуры
- Upgrading guide

**CONTRIBUTING.md:**
- Как добавить новый checker
- Как добавить новый notifier
- Code style guide
- PR process

---

## Метрики успеха

### Технические метрики

- ✅ Бинарник <20MB
- ✅ Memory usage <100MB при 100 targets
- ✅ Startup time <1 секунды
- ✅ Check latency overhead <10ms
- ✅ Test coverage >70%

### Функциональные метрики

- ✅ Поддержка 1000+ targets одновременно
- ✅ Check intervals от 10 секунд
- ✅ Alert latency <30 секунд
- ✅ 99.9% uptime самого монитора
- ✅ UI response time <100ms

### UX метрики

- ✅ Setup за <5 минут
- ✅ Добавление target за <1 минуту
- ✅ Интуитивный UI (не требует документации)
- ✅ Mobile-friendly

---

## Риски и митигации

### Технические риски

**Риск 1: Goroutine leaks при масштабировании**
- *Митигация:* Регулярное тестирование, proper context usage, pprof monitoring

**Риск 2: БД bottleneck при большом количестве targets**
- *Митигация:* Батчинг записей, async writes, индексы, миграция на PostgreSQL

**Риск 3: Memory leaks**
- *Митигация:* Regular profiling, load testing, bounded queues

**Риск 4: Сложность конфигурации**
- *Митигация:* Валидация, понятные error messages, примеры конфигов

### Бизнес риски

**Риск 1: Конкуренция с существующими решениями**
- *Митигация:* Фокус на простоту и легковесность, self-hosted без ограничений

**Риск 2: Поддержка и maintenance**
- *Митигация:* Качественная документация, простая архитектура, тесты

---

## Заключение

Этот план предоставляет структурированный подход к разработке легковесного, расширяемого монитора на Go. Следование фазам и приоритетам позволит:

1. Быстро получить работающий MVP
2. Итеративно добавлять функциональность
3. Поддерживать чистую архитектуру
4. Легко расширять систему новыми возможностями

### Следующие шаги

1. ✅ Ревью данного плана
2. ⏳ Создание GitHub репозитория
3. ⏳ Setup базовой структуры проекта
4. ⏳ Начало разработки с Фазы 1

---

**Удачи в разработке! 🚀**
