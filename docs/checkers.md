# Checkers — как использовать каждую проверку

Health Monitor поддерживает несколько типов проверок (checkers). Тип задаётся полем
`type` у таргета, а специфичные параметры — в объекте `config`. Таргеты создаются
**через Web UI** (раздел Targets → Add) или **через REST API** (`POST /api/v1/targets`).

Поддерживаемые типы: **`http`**, **`tcp`**, **`dns`** (ICMP — в планах).

---

## Общие поля таргета

| Поле | Тип | Описание |
|---|---|---|
| `id` | string | Уникальный ID (`[a-z0-9-]+`), напр. `api-prod` |
| `name` | string | Человекочитаемое имя |
| `type` | string | `http` / `tcp` / `dns` |
| `enabled` | bool | Включён ли таргет |
| `interval` | int (наносекунды) | Период проверки. В API — **наносекунды** (`30s` = `30000000000`) |
| `timeout` | int (наносекунды) | Таймаут одной проверки (`10s` = `10000000000`) |
| `description` | string | Опциональное описание |
| `tags` | []string | Опциональные теги |
| `config` | object | Параметры конкретного чекера (см. ниже) |

> ⚠️ **Важно для API:** `interval` и `timeout` сериализуются как наносекунды
> (Go `time.Duration`). В UI можно вводить `30s`/`1m`/`5m` — конвертация автоматическая.

Каждая проверка возвращает статус: `success` (всё хорошо), `warning` (работает,
но с нареканием — напр. медленный ответ), `failure` (недоступно/не совпало),
`unknown` (ещё не проверялось).

---

## HTTP / HTTPS (`type: "http"`)

Проверяет HTTP-эндпоинт: код ответа, время ответа, опционально срок SSL-сертификата.

### Поля `config`

| Поле | Тип | По умолчанию | Описание |
|---|---|---|---|
| `url` | string | — (**обязательно**) | Полный URL, напр. `https://api.example.com/health` |
| `method` | string | `GET` | HTTP-метод |
| `headers` | object | — | Кастомные заголовки (string→string) |
| `body` | string | — | Тело запроса |
| `expected_status_code` | int | `200` | Ожидаемый код ответа |
| `max_response_time_ms` | int | — | Порог времени ответа; превышение → статус `warning` |
| `follow_redirects` | bool | `true` | Следовать ли редиректам |
| `validate_ssl` | bool | `true` | Проверять валидность TLS-сертификата (`false` — для self-signed) |
| `check_ssl_expiry` | bool | `false` | Считать дни до истечения сертификата |
| `ssl_expiry_days` | int | — | Если до истечения ≤ N дней — добавить предупреждение в метадату |

### Логика результата
- `failure` — ошибка запроса, либо код ≠ `expected_status_code`.
- `warning` — код совпал, но время ответа > `max_response_time_ms`.
- `success` — код совпал и время в норме.
- Метадата: `url`, `method`, `status_text`, и при `check_ssl_expiry` — `ssl_expiry_days`, `ssl_expires_at`.

### Пример (API)
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

### В UI
Тип **HTTP/HTTPS** → поля **URL** и **Expected Status Code**. Остальные параметры
(метод GET, validate_ssl, check_ssl_expiry) выставляются разумными значениями по
умолчанию.

---

## TCP Port (`type: "tcp"`)

Проверяет доступность TCP-порта и измеряет время установления соединения.
Подходит для БД (5432, 3306), брокеров, SMTP, любых сокет-сервисов.

### Поля `config`

| Поле | Тип | По умолчанию | Описание |
|---|---|---|---|
| `host` | string | — (**обязательно**) | Хост или IP |
| `port` | int | — (**обязательно**) | Порт, 1–65535 |

### Логика результата
- `success` — соединение установлено (порт открыт). Время в `response_time_ms`.
- `failure` — отказ/таймаут (напр. `connection refused`). Текст в `error`.
- Метадата: `host`, `port`, `address`.

### Пример (API)
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

### В UI
Тип **TCP Port** → поля **Host** и **Port**.

---

## DNS (`type: "dns"`)

Резолвит DNS-записи домена, измеряет время резолва и (опционально) сверяет
результат с ожидаемыми значениями. Можно указать кастомный DNS-сервер.

### Поля `config`

| Поле | Тип | По умолчанию | Описание |
|---|---|---|---|
| `domain` | string | — (**обязательно**) | Домен для резолва |
| `record_type` | string | `A` | `A`, `AAAA`, `CNAME`, `MX`, `TXT`, `NS` |
| `dns_server` | string | системный | Кастомный DNS-сервер (`8.8.8.8` или `8.8.8.8:53`) |
| `expected_ips` | []string | — | Ожидаемые IP (для A/AAAA) |
| `expect_values` | []string | — | Ожидаемые значения (хосты MX/NS/CNAME, строки TXT) |

`expected_ips` и `expect_values` объединяются: если заданы, **каждое** значение
должно присутствовать в ответе (без учёта регистра, толерантно к завершающей точке;
для MX достаточно совпадения хоста).

### Логика результата
- `success` — резолв успешен и (если заданы) все ожидаемые значения найдены.
- `failure` — ошибка резолва (`no such host`), нет записей, либо ожидаемое значение отсутствует.
- Метадата: `domain`, `record_type`, `records` (дедуплицированы, отсортированы), `record_count`, при наличии — `dns_server`.

### Пример (API)
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

Чтобы проверять конкретные значения, добавьте `expect_values` (или `expected_ips`)
с актуальными для вашего домена записями, напр. `"expect_values": ["mx1.example.com"]`
для `record_type: "MX"`.

### В UI
Тип **DNS** → поля **Domain**, **Record Type**, **DNS Server** (опц.),
**Expected records** (через запятую, опц.).

---

## Data Retention

Старые данные автоматически удаляются фоновой задачей, чтобы БД не разрасталась.
Настраивается в секции `retention` конфигурации:

```yaml
retention:
  check_results: 720h    # хранить результаты проверок 30 дней
  incidents: 2160h       # хранить resolved-инциденты 90 дней
  cleanup_interval: 24h  # как часто запускать очистку
```

- Очистка запускается один раз сразу при старте, затем каждые `cleanup_interval`.
- Значение `0` для `check_results`/`incidents` = «хранить вечно» (очистка пропускается).
- `cleanup_interval` ≤ `0` полностью отключает фоновую задачу.
- Удаляются только **resolved**-инциденты; активные остаются.
