# Checkers — how to use each check

Health Monitor supports several types of checks (checkers). The type is set by the
`type` field of a target, and the specific parameters — in the `config` object. Targets are created
**through the Web UI** (Targets → Add section) or **through the REST API** (`POST /api/v1/targets`).

Supported types: **`http`**, **`tcp`**, **`dns`** (ICMP — planned).

---

## Common target fields

| Field | Type | Description |
|---|---|---|
| `id` | string | Unique ID (`[a-z0-9-]+`), e.g. `api-prod` |
| `name` | string | Human-readable name |
| `type` | string | `http` / `tcp` / `dns` |
| `enabled` | bool | Whether the target is enabled |
| `interval` | int (nanoseconds) | Check period. In the API — **nanoseconds** (`30s` = `30000000000`) |
| `timeout` | int (nanoseconds) | Timeout of a single check (`10s` = `10000000000`) |
| `description` | string | Optional description |
| `tags` | []string | Optional tags |
| `config` | object | Parameters of the specific checker (see below) |

> ⚠️ **Important for the API:** `interval` and `timeout` are serialized as nanoseconds
> (Go `time.Duration`). In the UI you can enter `30s`/`1m`/`5m` — conversion is automatic.

Each check returns a status: `success` (all good), `warning` (working,
but with a concern — e.g. a slow response), `failure` (unavailable/did not match),
`unknown` (not yet checked).

---

## HTTP / HTTPS (`type: "http"`)

Checks an HTTP endpoint: response code, response time, and optionally the SSL certificate expiry.

### `config` fields

| Field | Type | Default | Description |
|---|---|---|---|
| `url` | string | — (**required**) | Full URL, e.g. `https://api.example.com/health` |
| `method` | string | `GET` | HTTP method |
| `headers` | object | — | Custom headers (string→string) |
| `body` | string | — | Request body |
| `expected_status_code` | int | `200` | Expected response code |
| `max_response_time_ms` | int | — | Response time threshold; exceeding it → `warning` status |
| `follow_redirects` | bool | `true` | Whether to follow redirects |
| `validate_ssl` | bool | `true` | Validate the TLS certificate (`false` — for self-signed) |
| `check_ssl_expiry` | bool | `false` | Count the days until certificate expiry |
| `ssl_expiry_days` | int | — | If ≤ N days until expiry — add a warning to the metadata |

### Result logic
- `failure` — request error, or code ≠ `expected_status_code`.
- `warning` — code matched, but response time > `max_response_time_ms`.
- `success` — code matched and time is within bounds.
- Metadata: `url`, `method`, `status_text`, and with `check_ssl_expiry` — `ssl_expiry_days`, `ssl_expires_at`.

### Example (API)
```bash
curl -X POST http://localhost:8080/api/v1/targets -H 'Content-Type: application/json' -d '{
  "id": "api-health",
  "name": "Production API",
  "type": "http",
  "enabled": true,
  "interval": 30000000000,
  "timeout": 10000000000,
  "config": {
    "url": "https://api.example.com/health",
    "method": "GET",
    "expected_status_code": 200,
    "max_response_time_ms": 500,
    "validate_ssl": true,
    "check_ssl_expiry": true,
    "ssl_expiry_days": 30
  }
}'
```

### In the UI
Type **HTTP/HTTPS** → **URL** and **Expected Status Code** fields. The remaining parameters
(method GET, validate_ssl, check_ssl_expiry) are set to sensible default
values.

---

## TCP Port (`type: "tcp"`)

Checks the availability of a TCP port and measures the connection establishment time.
Suitable for databases (5432, 3306), brokers, SMTP, and any socket services.

### `config` fields

| Field | Type | Default | Description |
|---|---|---|---|
| `host` | string | — (**required**) | Host or IP |
| `port` | int | — (**required**) | Port, 1–65535 |

### Result logic
- `success` — connection established (port is open). Time in `response_time_ms`.
- `failure` — refused/timeout (e.g. `connection refused`). Text in `error`.
- Metadata: `host`, `port`, `address`.

### Example (API)
```bash
curl -X POST http://localhost:8080/api/v1/targets -H 'Content-Type: application/json' -d '{
  "id": "postgres-prod",
  "name": "Postgres",
  "type": "tcp",
  "enabled": true,
  "interval": 30000000000,
  "timeout": 5000000000,
  "config": { "host": "db.example.com", "port": 5432 }
}'
```

### In the UI
Type **TCP Port** → **Host** and **Port** fields.

---

## DNS (`type: "dns"`)

Resolves the DNS records of a domain, measures the resolution time and (optionally) compares
the result with expected values. You can specify a custom DNS server.

### `config` fields

| Field | Type | Default | Description |
|---|---|---|---|
| `domain` | string | — (**required**) | Domain to resolve |
| `record_type` | string | `A` | `A`, `AAAA`, `CNAME`, `MX`, `TXT`, `NS` |
| `dns_server` | string | system | Custom DNS server (`8.8.8.8` or `8.8.8.8:53`) |
| `expected_ips` | []string | — | Expected IPs (for A/AAAA) |
| `expect_values` | []string | — | Expected values (MX/NS/CNAME hosts, TXT strings) |

`expected_ips` and `expect_values` are combined: if set, **every** value
must be present in the response (case-insensitive, tolerant of a trailing dot;
for MX a host match is enough).

### Result logic
- `success` — resolution succeeded and (if set) all expected values were found.
- `failure` — resolution error (`no such host`), no records, or an expected value is missing.
- Metadata: `domain`, `record_type`, `records` (deduplicated, sorted), `record_count`, and if present — `dns_server`.

### Example (API)
```bash
curl -X POST http://localhost:8080/api/v1/targets -H 'Content-Type: application/json' -d '{
  "id": "dns-example",
  "name": "example.com A",
  "type": "dns",
  "enabled": true,
  "interval": 60000000000,
  "timeout": 5000000000,
  "config": {
    "domain": "example.com",
    "record_type": "A",
    "dns_server": "8.8.8.8"
  }
}'
```

To check specific values, add `expect_values` (or `expected_ips`)
with the records relevant to your domain, e.g. `"expect_values": ["mx1.example.com"]`
for `record_type: "MX"`.

### In the UI
Type **DNS** → **Domain**, **Record Type**, **DNS Server** (opt.),
**Expected records** (comma-separated, opt.) fields.

---

## Data Retention

Old data is automatically deleted by a background task so that the database does not grow unbounded.
Configured in the `retention` section of the configuration:

```yaml
retention:
  check_results: 720h    # keep check results for 30 days
  incidents: 2160h       # keep resolved incidents for 90 days
  cleanup_interval: 24h  # how often to run cleanup
```

- Cleanup runs once immediately at startup, then every `cleanup_interval`.
- A value of `0` for `check_results`/`incidents` = "keep forever" (cleanup is skipped).
- `cleanup_interval` ≤ `0` fully disables the background task.
- Only **resolved** incidents are deleted; active ones remain.
