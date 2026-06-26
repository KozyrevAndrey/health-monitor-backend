# Health Monitor - Development Plan

**Version:** 1.0  
**Date:** December 2024  
**Goal:** A lightweight, extensible system for monitoring domains, hosts, and health endpoints, written in Go

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [Technology Stack](#technology-stack)
3. [Architectural Principles](#architectural-principles)
4. [Development Phases](#development-phases)
5. [Project Structure](#project-structure)
6. [Roadmap and Priorities](#roadmap-and-priorities)
7. [Deployment and Infrastructure](#deployment-and-infrastructure)

---

## Project Overview

### Purpose
A self-hosted monitoring system for tracking the availability and performance of web services, API endpoints, and network resources.

### Key Features
- ✅ Lightweight (a single Go binary)
- ✅ Extensible architecture
- ✅ Support for many types of checks
- ✅ Flexible notification system
- ✅ Simple web interface
- ✅ Real-time updates
- ✅ Low resource footprint

### Target Use Cases
- Monitoring production services
- Checking API availability
- Tracking SSL certificates
- Monitoring internal services
- SLA tracking

---

## Technology Stack

### Backend
- **Language:** Go 1.21+
- **HTTP Router:** chi/gin (to be chosen during implementation)
- **Database:** SQLite (with the option to migrate to PostgreSQL)
- **ORM:** GORM
- **Scheduler:** robfig/cron or a custom ticker-based one
- **Logging:** zerolog/zap

### Frontend
- **Framework:** HTMX + Alpine.js
- **Styles:** Pico.css / Water.css
- **Real-time:** Server-Sent Events (SSE)
- **Charts:** Chart.js (optional)

### Infrastructure
- **Containerization:** Docker
- **Orchestration:** Docker Compose (Kubernetes optional)
- **CI/CD:** GitHub Actions

---

## Architectural Principles

### 1. Clean Architecture

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
- All dependencies are passed through constructors
- Use of interfaces instead of concrete types
- Easy testing via mocks

### 3. Interface Segregation
Key interfaces:

```go
// Examples of interfaces (concept, not actual code)

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
- **Checker Registry:** Registering new types of checks
- **Notifier Registry:** Adding notification channels
- **Storage Adapters:** Support for different databases
- **Plugin System:** (optional) Loading external plugins

---

## Development Phases

## Phase 1: Architecture Design

**Goal:** Design an extensible and maintainable architecture

### 1.1 Defining Interfaces

**Tasks:**
- [ ] Define the `Checker` interface with methods:
  - `Check(ctx, target) -> result`
  - `Type() -> string`
  - `Validate(config) -> error`
- [ ] Define the `Storage` interface:
  - CRUD for targets
  - Saving check results
  - Retrieving history and statistics
- [ ] Define the `Notifier` interface:
  - `Notify(ctx, alert) -> error`
  - `Type() -> string`
- [ ] Define the `Scheduler` interface:
  - Task management
  - Lifecycle methods (Start/Stop)

**Outcome:** A document describing all interfaces

### 1.2 Configuration Structure

**Tasks:**
- [ ] Design the YAML configuration schema
- [ ] Define the Target structure:
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
- [ ] Define global settings
- [ ] Design retention policies

**Outcome:** Example configuration files

### 1.3 Data Models

**Tasks:**
- [ ] Describe the `Target` structure
- [ ] Describe the `CheckResult` structure:
  - Timestamp
  - Success/Failure
  - Response time
  - Status code
  - Error message
  - Metadata (headers, body snippet)
- [ ] Describe the `Alert` structure
- [ ] Describe the `Incident` structure:
  - Grouping of consecutive failures
  - Start/End time
  - Status (ongoing/resolved)

**Outcome:** Domain models as Go structs (conceptually)

---

## Phase 2: Basic Infrastructure

**Goal:** Set up the project and the basic infrastructure

### 2.1 Project Setup

**Tasks:**
- [ ] Initialize the Go module
  - `go mod init github.com/username/health-monitor`
- [ ] Create the directory structure:
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
- [ ] Configure golangci-lint
- [ ] Create a Makefile with commands:
  - `make build`
  - `make test`
  - `make lint`
  - `make run`
  - `make docker-build`

**Outcome:** A working project skeleton

### 2.2 Configuration System

**Tasks:**
- [ ] Choose a library (viper / cleanenv / custom)
- [ ] Implement YAML parsing
- [ ] Add configuration validation
- [ ] Support env variables for sensitive data
- [ ] Hot-reload of the config (optional for v2)

**Outcome:** Working configuration loading

### 2.3 Logging

**Tasks:**
- [ ] Choose a library (zerolog recommended)
- [ ] Set up structured logging
- [ ] Configurable log levels
- [ ] Logging to a file + stdout
- [ ] Log rotation (if logging to a file)

**Outcome:** Centralized logging

---

## Phase 3: Storage Layer

**Goal:** Implement data persistence

### 3.1 Repository Interface

**Tasks:**
- [ ] Define the `TargetRepository` interface:
  - `Create(target) error`
  - `Get(id) (Target, error)`
  - `List() ([]Target, error)`
  - `Update(target) error`
  - `Delete(id) error`
- [ ] Define the `CheckResultRepository` interface:
  - `Save(result) error`
  - `GetHistory(targetID, limit, offset) ([]Result, error)`
  - `GetLatest(targetID) (Result, error)`
  - `GetStats(targetID, from, to) (Stats, error)`
- [ ] Define the `IncidentRepository` interface

**Outcome:** A description of all repository methods

### 3.2 SQLite Implementation

**Tasks:**
- [ ] Design the database schema:
  ```sql
  -- Example structure
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
- [ ] Set up migrations (golang-migrate or goose)
- [ ] Implement all repository methods
- [ ] Add indexes for performance
- [ ] Connection pooling
- [ ] Transactions where necessary

**Outcome:** Working SQLite storage

### 3.3 Preparing for Scaling

**Tasks:**
- [ ] Abstract away SQL-specific code
- [ ] Prepare an interface for PostgreSQL
- [ ] Graceful shutdown with data flush
- [ ] Automatic cleanup of old data (retention)

**Outcome:** The ability to switch databases

---

## Phase 4: Checker System

**Goal:** Implement various types of checks

### 4.1 Base Interface

**Tasks:**
- [ ] Implement the base `Checker` interface
- [ ] Add timeout handling
- [ ] Implement context propagation
- [ ] Optional retry logic with backoff
- [ ] Performance metrics for checkers

**Outcome:** A base mechanism for checks

### 4.2 HTTP Checker (priority 1)

**Tasks:**
- [ ] Check the HTTP status code
- [ ] Measure response time
- [ ] Validate the response body:
  - Check for the presence of a string
  - Regex matching
  - JSON path validation
- [ ] Check headers
- [ ] Check SSL certificates:
  - Validity
  - Expiration date (warning N days in advance)
- [ ] Follow redirects (configurable)
- [ ] Custom HTTP methods (GET, POST, HEAD)
- [ ] Custom headers and authentication

**Outcome:** A fully featured HTTP checker

### 4.3 TCP Checker (priority 2)

**Tasks:**
- [ ] Check port availability
- [ ] Measure connection time
- [ ] Optional sending/receiving of data
- [ ] TLS support

**Outcome:** A TCP port checker

### 4.4 DNS Checker (priority 3)

**Tasks:**
- [ ] Resolving A/AAAA records
- [ ] Checking MX, TXT, CNAME records
- [ ] Measuring DNS resolution time
- [ ] Checking against expected values

**Outcome:** A DNS checker

### 4.5 ICMP/Ping Checker (priority 4)

**Tasks:**
- [ ] ICMP ping
- [ ] Measuring RTT
- [ ] Packet loss detection
- [ ] Requires privileges - document this

**Outcome:** A ping checker

### 4.6 Custom Script Checker (optional)

**Tasks:**
- [ ] Executing external scripts
- [ ] Parsing the exit code
- [ ] Timeout for scripts
- [ ] Security (sandboxing)

**Outcome:** Extension via scripts

### 4.7 Checker Registry

**Tasks:**
- [ ] Factory pattern for creating checkers
- [ ] Registering new types
- [ ] Validating configuration by type

**Outcome:** An easily extensible system

---

## Phase 5: Scheduler

**Goal:** Scheduling and executing checks

### 5.1 Task Manager

**Tasks:**
- [ ] Ticker-based scheduler for each target
- [ ] Worker pool for parallel execution:
  - Configurable number of workers
  - Queue for tasks
  - Graceful worker shutdown
- [ ] Lifecycle management:
  - `Start()` - start the scheduler
  - `Stop()` - graceful stop
  - `AddTarget()` - add a new target
  - `RemoveTarget()` - remove a target
  - `UpdateTarget()` - update the interval
- [ ] Context propagation for cancellation

**Outcome:** A basic scheduler

### 5.2 Optimizations

**Tasks:**
- [ ] Jitter to avoid thundering herd:
  - Random delay at startup
  - Distributing checks over time
- [ ] Backoff on consecutive failures:
  - Increasing the interval on errors
  - Returning to the normal interval on success
- [ ] Check prioritization:
  - Critical targets are checked more frequently
  - Lower priority under high load
- [ ] Health check of the scheduler itself

**Outcome:** An optimized scheduler

### 5.3 Scheduler Metrics

**Tasks:**
- [ ] Number of active tasks
- [ ] Queue depth
- [ ] Worker utilization
- [ ] Avg check duration

**Outcome:** Scheduler observability

---

## Phase 6: Alerting System

**Goal:** Notifications about problems

### 6.1 Alert Manager

**Tasks:**
- [ ] State machine for alerts:
  ```
  OK → WARNING → CRITICAL → RECOVERY → OK
  ```
- [ ] Defining conditions:
  - N consecutive failures
  - Response time exceeds the threshold
  - SSL certificate expires in N days
- [ ] Debouncing:
  - Don't send an alert on every flap
  - Grace period before an alert
  - Recovery notification
- [ ] Grouping into incidents:
  - One incident for a series of failures
  - Tracking downtime
  - Resolution tracking

**Outcome:** A smart alerting system

### 6.2 Notifier Implementations

**Priority 1: Webhook Notifier**
- [ ] POST request to a URL
- [ ] Configurable payload template
- [ ] Retry logic
- [ ] Custom headers

**Priority 2: Email Notifier**
- [ ] SMTP sending
- [ ] HTML templates for emails
- [ ] Attachment support (charts, optional)

**Priority 3: Telegram Notifier**
- [ ] Bot API integration
- [ ] Formatted messages
- [ ] Multiple chat IDs

**Optional:**
- [ ] Slack notifier
- [ ] Discord notifier
- [ ] PagerDuty integration
- [ ] Custom webhook formats (Slack, Discord)

**Outcome:** Multiple notification channels

### 6.3 Alert Routing

**Tasks:**
- [ ] Different notifiers for different targets
- [ ] Severity-based routing:
  - WARNING → email
  - CRITICAL → email + telegram + webhook
- [ ] Quiet hours (do not disturb at night):
  - Configurable time windows
  - Timezone support
- [ ] Escalation policies (optional):
  - Level 1 → email after 5 min
  - Level 2 → telegram after 15 min
  - Level 3 → phone call after 30 min

**Outcome:** A flexible routing system

### 6.4 Alert Templates

**Tasks:**
- [ ] Template engine for messages
- [ ] Available variables:
  - Target name, URL
  - Error message
  - Response time
  - Timestamp
  - Downtime duration
- [ ] Different templates for different events:
  - DOWN alert
  - RECOVERY alert
  - WARNING alert
  - SSL expiration

**Outcome:** Customizable messages

---

## Phase 7: API Layer

**Goal:** A REST API for management and data

### 7.1 REST API Endpoints

**Targets Management:**
- [ ] `GET /api/v1/targets` - list all targets
  - Filtering by type, status
  - Pagination
  - Sorting
- [ ] `GET /api/v1/targets/:id` - target details
- [ ] `GET /api/v1/targets/:id/status` - current status
- [ ] `POST /api/v1/targets` - create (optional for MVP)
- [ ] `PUT /api/v1/targets/:id` - update
- [ ] `DELETE /api/v1/targets/:id` - delete

**Check History:**
- [ ] `GET /api/v1/targets/:id/checks` - check history
  - Pagination
  - Date range filtering
  - Success/failure filtering
- [ ] `GET /api/v1/targets/:id/stats` - statistics
  - Uptime percentage
  - Avg response time
  - Success rate
  - Time range (24h, 7d, 30d)

**Incidents:**
- [ ] `GET /api/v1/incidents` - list incidents
- [ ] `GET /api/v1/incidents/:id` - incident details

**System:**
- [ ] `GET /api/health` - health check of the monitor itself
- [ ] `GET /api/metrics` - system metrics (Prometheus format optional)

**Outcome:** A fully featured API

### 7.2 Real-time Updates

**Tasks:**
- [ ] Server-Sent Events (SSE) endpoint:
  - `GET /api/v1/events` - event stream
- [ ] Events to send:
  - Check completed
  - Status changed
  - Alert triggered
  - Alert resolved
- [ ] Client reconnection handling
- [ ] Event filtering by target ID

**Outcome:** Real-time notifications

### 7.3 Middleware

**Tasks:**
- [ ] Logging middleware:
  - Request/response logging
  - Duration tracking
- [ ] Recovery middleware (panic recovery)
- [ ] CORS middleware (if needed)
- [ ] Rate limiting (optional)
- [ ] Authentication middleware:
  - Basic Auth to start with
  - API keys (optional)
  - JWT (for production)
- [ ] Request ID for tracing

**Outcome:** A production-ready API

### 7.4 API Documentation

**Tasks:**
- [ ] OpenAPI/Swagger specification
- [ ] Request/response examples
- [ ] Swagger UI (optional)

**Outcome:** A documented API

---

## Phase 8: Web UI Integration

**Goal:** A web interface for monitoring

### 8.1 Static Files

**Tasks:**
- [ ] Embed HTML/CSS/JS into the binary (`go:embed`)
- [ ] Serving static files via an HTTP handler
- [ ] Cache headers for static files
- [ ] Compression (gzip)

**Outcome:** Single binary deployment

### 8.2 Template Rendering

**Tasks:**
- [ ] HTML templates (`html/template`)
- [ ] Layout + partials
- [ ] Passing data into templates
- [ ] HTMX integration:
  - Partial updates
  - Auto-refresh of the targets list
  - Inline editing (optional)

**Outcome:** Server-side rendering

### 8.3 Dashboard Components

**Pages:**

**1. Overview Dashboard:**
- [ ] Cards with statistics:
  - Total targets
  - Healthy / Unhealthy
  - Active incidents
- [ ] List of all targets with their current status
- [ ] Filtering and search
- [ ] Sorting by status, name, response time

**2. Target Detail Page:**
- [ ] Current status (large indicator)
- [ ] Last N checks
- [ ] Response time chart (Chart.js)
- [ ] Uptime over different periods (24h, 7d, 30d)
- [ ] Incident history
- [ ] Target configuration

**3. Incidents Page:**
- [ ] List of active incidents
- [ ] History of resolved incidents
- [ ] Timeline visualization

**4. Settings Page (optional):**
- [ ] Managing targets (CRUD)
- [ ] Notifier configuration
- [ ] System settings

**Outcome:** A fully featured UI

### 8.4 Real-time Updates in the UI

**Tasks:**
- [ ] SSE connection from the browser
- [ ] Auto-updating statuses without a reload
- [ ] Browser notifications (optional)
- [ ] Visual indicators for changes

**Outcome:** A live dashboard

### 8.5 Responsive Design

**Tasks:**
- [ ] Mobile-friendly layout
- [ ] Responsive tables
- [ ] Touch-friendly buttons

**Outcome:** Works on all devices

---

## Phase 9: Observability

**Goal:** Monitoring the monitor itself

### 9.1 Internal Metrics

**Tasks:**
- [ ] Metrics about the system itself:
  - Number of active targets
  - Number of checks being executed
  - Queue depth
  - Worker utilization
  - DB connection pool usage
  - Memory usage
  - Goroutines count
- [ ] Prometheus metrics endpoint (optional):
  - `GET /metrics` in Prometheus format
- [ ] Performance metrics:
  - Avg check duration by type
  - P50/P95/P99 response times
  - Error rates

**Outcome:** System observability

### 9.2 Health Check

**Tasks:**
- [ ] Endpoint `GET /api/health`
- [ ] Checks:
  - DB connection alive
  - Scheduler is running
  - All workers are alive
  - Disk space available
- [ ] Health status:
  - `healthy` - everything OK
  - `degraded` - partial problems
  - `unhealthy` - critical problems
- [ ] Detailed problem breakdown in the response

**Outcome:** Self-monitoring

### 9.3 Structured Events

**Tasks:**
- [ ] Event log for important events:
  - System started/stopped
  - Target added/removed
  - Alert triggered/resolved
  - Configuration reloaded
- [ ] Audit log (optional):
  - API calls
  - Configuration changes
  - User actions

**Outcome:** A complete history of events

---

## Phase 10: Production-Ready Features

**Goal:** Preparing for production use

### 10.1 Graceful Shutdown

**Tasks:**
- [ ] Signal handling:
  - SIGTERM
  - SIGINT
- [ ] Shutdown sequence:
  1. Stop accepting new HTTP requests
  2. Complete active HTTP requests (with a timeout)
  3. Stop the scheduler (no new checks)
  4. Complete active checks (with a timeout)
  5. Flush data to the database
  6. Close DB connections
  7. Close notifiers
- [ ] Graceful timeout (30 seconds by default)
- [ ] Logging of the shutdown process

**Outcome:** A safe shutdown

### 10.2 Configuration Validation

**Tasks:**
- [ ] Validation at startup:
  - YAML correctness
  - Required fields
  - Valid intervals (>0)
  - Valid URLs
- [ ] Validate command:
  - `health-monitor validate config.yaml`
- [ ] Detailed error messages

**Outcome:** Prevention of configuration errors

### 10.3 Data Retention

**Tasks:**
- [ ] Automatic cleanup of old data:
  - Configurable retention period
  - Aggregation of old data (optional)
- [ ] Background job for cleanup
- [ ] Logging of deleted records

**Outcome:** Control over database size

### 10.4 Backup & Restore

**Tasks:**
- [ ] Backup command:
  - `health-monitor backup --output backup.db`
  - Backing up SQLite
- [ ] Restore command
- [ ] Process documentation

**Outcome:** Data protection

### 10.5 Security

**Tasks:**
- [ ] Authentication:
  - Basic Auth for the web interface
  - API keys for the API
- [ ] HTTPS support:
  - TLS configuration
  - Let's Encrypt integration (optional)
- [ ] Rate limiting on the API
- [ ] Input validation for all endpoints
- [ ] SQL injection prevention (via the ORM)
- [ ] XSS prevention in templates

**Outcome:** A secure application

### 10.6 Performance

**Tasks:**
- [ ] Connection pooling for the database
- [ ] HTTP client pooling for checkers
- [ ] Caching where possible:
  - Static files
  - Config (in memory)
- [ ] Database indexes
- [ ] Pagination for large lists

**Outcome:** Optimized performance

---

## Project Structure

### Final Directory Structure

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

## Roadmap and Priorities

### MVP (Minimum Viable Product)

**Goal:** A working monitor with basic functionality

**Includes:**
- ✅ HTTP checker (status code, response time)
- ✅ SQLite storage
- ✅ Simple scheduler (without a worker pool)
- ✅ Webhook notifier
- ✅ Basic REST API (read-only)
- ✅ Simple web interface (targets list + details)
- ✅ Configuration via YAML

**Does not include:**
- Additional checkers (TCP, DNS, ping)
- Additional notifiers (email, telegram)
- CRUD operations in the UI
- Real-time updates
- Advanced alerting rules

**Timeframe:** 2-3 weeks of development

---

### Version 1.0

**Goal:** A production-ready solution

**Adds to the MVP:**
- ✅ Worker pool in the scheduler
- ✅ TCP and DNS checkers
- ✅ Email and Telegram notifiers
- ✅ Incident management
- ✅ Real-time updates (SSE)
- ✅ Improved UI with charts
- ✅ Graceful shutdown
- ✅ Authentication
- ✅ Metrics endpoint

**Timeframe:** +2-3 weeks after the MVP

---

### Version 1.1

**Goal:** Extended capabilities

**Adds:**
- ✅ CRUD operations in the UI
- ✅ PostgreSQL support
- ✅ SSL certificate monitoring
- ✅ Advanced alerting rules
- ✅ Alert routing and escalation
- ✅ Data retention policies
- ✅ Prometheus metrics export

**Timeframe:** +2 weeks after v1.0

---

### Version 2.0 (Future)

**Ideas for the future:**
- 🔮 Multi-user support with permissions
- 🔮 Incident timeline visualization
- 🔮 SLA tracking and reporting
- 🔮 Status page generation (a public status page)
- 🔮 Distributed monitoring (multiple instances)
- 🔮 Plugin system for custom checkers
- 🔮 Mobile app
- 🔮 Advanced analytics and ML for anomaly detection

---

## Deployment and Infrastructure

### Docker Deployment

**Dockerfile (multi-stage):**

```dockerfile
# Conceptual example of the structure

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
# Conceptual example

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
# Example systemd unit file

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

### Kubernetes Deployment (optional)

```yaml
# Conceptual example

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

## Testing

### Unit Tests

**Coverage:**
- [ ] All checkers (mocks for the HTTP client)
- [ ] Storage repositories (in-memory or testcontainers)
- [ ] Alert manager state machine
- [ ] Scheduler logic
- [ ] API handlers (HTTP tests)

**Tools:**
- `testing` package
- `testify` for assertions
- `gomock` for mocks
- `httptest` for HTTP tests

### Integration Tests

**Scenarios:**
- [ ] End-to-end flow: config → check → save → alert
- [ ] API integration tests
- [ ] DB migrations
- [ ] Real checkers against test servers

### Load Tests (optional)

**Verify:**
- [ ] Performance with 1000+ targets
- [ ] Memory leaks
- [ ] Goroutine leaks
- [ ] DB performance

---

## Documentation

### Required Documentation

**README.md:**
- Project description
- Quick start guide
- Installation instructions
- Basic usage examples

**docs/configuration.md:**
- Full description of all configuration options
- Examples for different use cases
- Best practices

**docs/api.md:**
- API endpoints documentation
- Request/response examples
- Authentication

**docs/architecture.md:**
- Architecture diagram
- Description of components
- Design decisions

**docs/deployment.md:**
- Docker deployment
- Systemd setup
- Kubernetes deployment
- Backup/restore procedures
- Upgrading guide

**CONTRIBUTING.md:**
- How to add a new checker
- How to add a new notifier
- Code style guide
- PR process

---

## Success Metrics

### Technical Metrics

- ✅ Binary <20MB
- ✅ Memory usage <100MB with 100 targets
- ✅ Startup time <1 second
- ✅ Check latency overhead <10ms
- ✅ Test coverage >70%

### Functional Metrics

- ✅ Support for 1000+ targets simultaneously
- ✅ Check intervals from 10 seconds
- ✅ Alert latency <30 seconds
- ✅ 99.9% uptime of the monitor itself
- ✅ UI response time <100ms

### UX Metrics

- ✅ Setup in <5 minutes
- ✅ Adding a target in <1 minute
- ✅ Intuitive UI (requires no documentation)
- ✅ Mobile-friendly

---

## Risks and Mitigations

### Technical Risks

**Risk 1: Goroutine leaks when scaling**
- *Mitigation:* Regular testing, proper context usage, pprof monitoring

**Risk 2: DB bottleneck with a large number of targets**
- *Mitigation:* Batching writes, async writes, indexes, migration to PostgreSQL

**Risk 3: Memory leaks**
- *Mitigation:* Regular profiling, load testing, bounded queues

**Risk 4: Configuration complexity**
- *Mitigation:* Validation, clear error messages, example configs

### Business Risks

**Risk 1: Competition with existing solutions**
- *Mitigation:* Focus on simplicity and being lightweight, self-hosted without limitations

**Risk 2: Support and maintenance**
- *Mitigation:* High-quality documentation, simple architecture, tests

---

## Conclusion

This plan provides a structured approach to developing a lightweight, extensible monitor in Go. Following the phases and priorities will allow us to:

1. Quickly get a working MVP
2. Iteratively add functionality
3. Maintain a clean architecture
4. Easily extend the system with new capabilities

### Next Steps

1. ✅ Review of this plan
2. ⏳ Creating the GitHub repository
3. ⏳ Setting up the basic project structure
4. ⏳ Starting development from Phase 1

---

**Good luck with the development! 🚀**
