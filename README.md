# Health Monitor

A lightweight, extensible, self-hosted monitoring system for tracking the availability and performance of web services, API endpoints, and network resources.

## Features

### Currently Available ✅
- **Lightweight** - single Go binary (~15MB)
- **HTTP/HTTPS monitoring** - with SSL validation and expiry checking
- **Flexible scheduling** - independent intervals per target
- **SQLite storage** - automatic migrations and data retention
- **YAML configuration** - with environment variable overrides
- **Structured logging** - JSON and console formats
- **Graceful shutdown** - proper cleanup of all resources
- **Docker support** - multi-stage builds with Alpine
- **Clean Architecture** - fully testable and extensible
- **Alert management** - consecutive failures, response time, SSL expiry detection
- **Incident tracking** - automatic creation and resolution
- **Telegram notifications** - rich Markdown messages with icons and metadata
- **Email notifications** - HTML and plain text with SMTP support

### Coming Soon 🚧
- Multiple check types (TCP, DNS, ICMP)
- Additional notifiers (Webhook)
- HTTP API (REST endpoints)
- Simple web interface
- Real-time updates via Server-Sent Events

## Quick Start

### Prerequisites

- Go 1.21+ (for building from source)
- Docker & Docker Compose (for containerized deployment)

### Build from Source

```bash
# Clone the repository
git clone <repository-url>
cd health-monitor-backend

# Build the application
make build

# Run with example configuration
make run
```

### Run with Docker

```bash
# Build and run with docker-compose
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

### Run Binary Directly

```bash
# After building
./bin/health-monitor --config configs/example.yaml
```

## Configuration

Configuration is done via YAML file. See `configs/example.yaml` for a complete example.

### Basic Configuration

```yaml
server:
  host: "0.0.0.0"
  port: 8080

database:
  type: "sqlite"
  path: "./data/health-monitor.db"

logging:
  level: "info"
  format: "json"

targets:
  - id: "my-api"
    name: "My API"
    type: "http"
    enabled: true
    interval: "30s"
    config:
      url: "https://api.example.com/health"
      expected_status_code: 200
```

### Environment Variables

Configuration can be overridden using environment variables with prefix `HEALTH_MONITOR_`:

```bash
export HEALTH_MONITOR_SERVER_PORT=9090
export HEALTH_MONITOR_LOGGING_LEVEL=debug
```

## Target Types

### HTTP/HTTPS

Monitor HTTP endpoints with configurable checks:

```yaml
- id: "api-check"
  type: "http"
  config:
    url: "https://api.example.com/health"
    method: "GET"
    expected_status_code: 200
    max_response_time_ms: 500
    validate_ssl: true
    check_ssl_expiry: true
```

### TCP Port

Check TCP port availability:

```yaml
- id: "database"
  type: "tcp"
  config:
    host: "db.example.com"
    port: 5432
```

### DNS

Verify DNS resolution:

```yaml
- id: "dns-check"
  type: "dns"
  config:
    domain: "example.com"
    record_type: "A"
    expected_ips:
      - "93.184.216.34"
```

### ICMP (Ping)

Monitor host reachability:

```yaml
- id: "gateway"
  type: "icmp"
  config:
    host: "192.168.1.1"
    packet_count: 3
    max_rtt_ms: 100
```

## Notification Channels

### Webhook

Send alerts to any webhook endpoint:

```yaml
notifiers:
  - id: "slack-webhook"
    type: "webhook"
    config:
      url: "https://hooks.slack.com/services/YOUR/WEBHOOK"
      method: "POST"
```

### Email

Email notifications via SMTP:

```yaml
notifiers:
  - id: "email-team"
    type: "email"
    config:
      smtp_host: "smtp.gmail.com"
      smtp_port: 587
      from: "alerts@example.com"
      to:
        - "team@example.com"
```

### Telegram

Telegram bot notifications:

```yaml
notifiers:
  - id: "telegram"
    type: "telegram"
    config:
      bot_token: "YOUR_BOT_TOKEN"
      chat_ids:
        - "-1001234567890"
```

## Development

### Project Structure

```
health-monitor/
cmd/server/      # Application entry point
internal/        # Internal packages
domain/          # Domain models & interfaces
checker/         # Health check implementations
storage/         # Database layer
scheduler/       # Task scheduling
alerting/        # Alert management
notifier/        # Notification implementations
api/             # HTTP API
pkg/             # Public packages
config/          # Configuration management
logger/          # Logging utilities
web/             # Frontend assets
configs/         # Configuration examples
docs/            # Documentation
```

### Make Commands

```bash
make help           # Show all available commands
make build          # Build the application
make test           # Run tests
make lint           # Run linter
make fmt            # Format code
make run            # Build and run
make dev            # Run without building
make clean          # Clean build artifacts
make docker-build   # Build Docker image
```

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific package tests
go test -v ./pkg/config/...
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run vet
make vet
```

## Architecture

The project follows Clean Architecture principles with clear separation of concerns:

- **Domain Layer**: Core business logic and interfaces
- **Application Layer**: Use cases and orchestration
- **Infrastructure Layer**: External dependencies (DB, HTTP, etc.)
- **Presentation Layer**: API handlers and UI

### Key Interfaces

- `Checker`: Health check implementations
- `Storage`: Data persistence
- `Notifier`: Alert delivery
- `Scheduler`: Task scheduling

## Deployment

### Docker

Single container deployment:

```bash
docker build -t health-monitor .
docker run -p 8080:8080 \
  -v $(pwd)/configs:/configs \
  -v $(pwd)/data:/data \
  health-monitor
```

### Docker Compose

Full stack with PostgreSQL (optional):

```bash
docker-compose up -d
```

### Systemd

Create a systemd service:

```bash
sudo cp health-monitor /usr/local/bin/
sudo cp configs/example.yaml /etc/health-monitor/config.yaml

# Create systemd service file
sudo nano /etc/systemd/system/health-monitor.service

# Enable and start
sudo systemctl enable health-monitor
sudo systemctl start health-monitor
```

## API Documentation

The application exposes a REST API for programmatic access.

### Endpoints

- `GET /api/v1/targets` - List all targets
- `GET /api/v1/targets/:id` - Get target details
- `GET /api/v1/targets/:id/checks` - Get check history
- `GET /api/v1/targets/:id/stats` - Get statistics
- `GET /api/v1/incidents` - List incidents
- `GET /api/health` - Health check
- `GET /api/v1/events` - Server-Sent Events stream

## Roadmap

**Progress: 70%**

### Core Features (Phase 1-7) - Completed
- [x] Basic infrastructure setup (Makefile, Docker, CI/CD)
- [x] Configuration system (Viper with YAML and env variables)
- [x] Logging system (zerolog with structured logging)
- [x] Domain models and interfaces (Clean Architecture)
- [x] SQLite storage layer (GORM with auto-migration)
  - Targets, check results, incidents repositories
  - Statistics and history queries
- [x] HTTP/HTTPS checker implementation
  - Status code validation
  - Response time measurement
  - SSL certificate validation and expiry checking
- [x] Scheduler system
  - Ticker-based periodic checks
  - Concurrent execution
  - Graceful shutdown
- [x] Alert manager
  - Consecutive failures detection (default: 3)
  - Response time threshold alerts (default: 5s)
  - SSL certificate expiry warnings (default: 30 days)
  - DOWN/UP transition alerts
  - Incident tracking and automatic resolution
  - Alert deduplication with cooldown periods
- [x] Telegram notifier
  - Bot API integration with HTTP client
  - Multiple chat IDs support
  - Rich Markdown formatting with icons
  - Alert type-specific icons (🔴 DOWN, 🟢 UP, 🐌 slow, 🔐 SSL)
  - Metadata and timestamp display
- [x] Email notifier
  - SMTP support with TLS/SSL
  - HTML and plain text email formats
  - Multiple recipients support
  - Rich HTML templates with severity colors
  - Multipart/alternative MIME format

### In Progress (Phase 8+)
- [ ] Additional notifiers (Webhook)

### Planned (Phase 7+)
- [ ] Additional checkers (TCP, DNS, ICMP)
- [ ] HTTP API (REST endpoints)
- [ ] Simple web UI (dashboard)
- [ ] Real-time updates (Server-Sent Events)
- [ ] Additional notifiers (Email, Telegram)
- [ ] PostgreSQL support
- [ ] Advanced alerting rules
- [ ] Multi-user support and RBAC

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[Your License Here]

## Support

For issues and questions, please use the GitHub issue tracker.
