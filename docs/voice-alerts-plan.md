# Voice alerts (phone calls) — implementation plan

Status: v1 implemented. Provider: SMSPILOT voice calls
(`https://smspilot.ru/api.php`, `from=GOLOS`).

Where it lives: `internal/domain/policy.go` (delivery policy),
`internal/alerting/dispatch.go` (incident dispatcher),
`internal/storage/notification_send_repository.go` (send journal),
`internal/notifier/smspilot.go` (provider), `internal/notifier/factory.go`
(one place to register a provider), voice fields in `web/static/js/app.js`.

Two things changed versus the plan below: the endpoint is overridden through
the server-side `SMSPILOT_API_URL` variable and a package-internal constructor
rather than a config key, and the API key is env-only, so no masking of an
`api_key` field was needed.

**v1 scope, agreed with the stakeholder (Nick Mitin, 2026-09-23): one call per
incident, no repeats.** The call fires only after N consecutive real failures,
because a single failed check is often a false alarm. Repeat/escalation
(call every 30 minutes until acknowledged) is deliberately deferred — the
schema and the policy struct leave room for it, the scheduling code is not
written yet.

---

## Design in one line

A generic per-notifier delivery policy (how many failures before notifying, how
many times, how often) plus one asynchronous provider implementing the existing
`domain.Notifier` interface. The provider is configured from the existing
notifier form in the dashboard and stored in the existing `notifier_configs`
table.

```text
check result
  → incident + FailureCount (DB)
  → per-notifier policy check
  → claim row in notification_sends (unique)
  → notifier.Notify
  → store delivery receipt
```

No separate OS process and no background worker: dispatch happens on the
scheduler tick that already delivers every check result to
`alerting.Manager.ProcessCheckResult`.

---

## What the current code does today (verified)

These are the constraints that shape the plan:

1. A `down` alert fires on the **first** failure — `internal/alerting/manager.go:180-182`
   (`result.Status == Failure && previousStatus == Success`). Wiring a voice
   notifier into the existing alert path would call on the first blip.
2. Alert rules are hardcoded: consecutive-failure threshold is 3
   (`getDefaultAlertRules`, `manager.go:389`), with a 5-minute in-memory
   cooldown per alert type (`shouldSendAlert`, `manager.go:199-210`).
3. `evaluateAlertRules` runs **before** `manageIncident`, so on the first
   failure no incident exists yet and there is no `incident_id` to key
   deliveries on.
4. `updateState` (`manager.go:116-128`) increments `consecutiveFailures` for
   anything that is not `Success` — including `CheckStatusWarning` (slow
   response). Three slow-but-healthy checks currently trigger the
   consecutive-failure alert.
5. Alert state lives in memory (`Manager.states`). After a restart
   `currentIncidentID` is `nil`, so a failure creates a **duplicate** ongoing
   incident and a success never resolves the old one — it stays `ongoing`
   forever.
6. `Manager.notifiers` is keyed by `notifier.Type()` (`manager.go:59-63`), so
   two configs of the same type overwrite each other. Production callers of
   `RegisterNotifier` are only `cmd/server/main.go:237` and
   `internal/api/notifier_handlers.go:168`.
7. `maskNotifierConfig` (`notifier_handlers.go:180-186`) masks
   `smtp_password`, `bot_token`, `password`, `token`, `secret` — an `api_key`
   field would be returned in clear text by the API.
8. Schema is managed by `Database.AutoMigrate` (`internal/storage/database.go:88-97`).
   `migrations/001_initial_schema.sql` is not applied at runtime.

---

## Stage 0 — restart recovery (required before any call is enabled)

In `ProcessCheckResult`, the `!exists` branch (`manager.go:84-91`) must hydrate
state from the database: look up `incidentRepo.GetOngoing(targetID)` and restore
`currentIncidentID` and the failure count from it.

This is a prerequisite, not a nice-to-have: without it a restart during an
incident produces a duplicate incident and a permanently unresolved old one —
which, once repeats exist, means calling about a service that already recovered.
It also fixes duplicate incidents that happen today.

## Stage 1 — key the notifier registry by config ID

Change `Manager.notifiers` from `map[string]Notifier` keyed by type to a map
keyed by the notifier config ID, holding the parsed policy next to the notifier.
Touches `domain.AlertManager` (`internal/domain/interfaces.go:149`) and the two
call sites listed above. Without this, per-notifier policies have nothing stable
to attach to, and two Telegram channels still clobber each other.

## Stage 2 — delivery policy and the deliveries journal

**Policy lives in the notifier's own config JSON** (`notifier_configs.Config`),
not in a separate policies table. That table already has CRUD, masking and
automatic reload on save (`notifier_handlers.go:140-176`), so this needs no
migration and no new API surface:

```json
{
  "min_failures": 3,
  "max_sends": 1,
  "repeat_interval": "",
  "target_ids": []
}
```

`repeat_interval` empty or `0` means "once" — the v1 default. `target_ids`
empty means all targets.

**One new table, `notification_sends`**, added to `AutoMigrate` with a composite
unique index on `(incident_id, notifier_id, sequence)`:

| column | purpose |
| --- | --- |
| `incident_id`, `notifier_id`, `sequence` | identity + dedupe (unique index) |
| `status` | `sent`, `failed`, `unknown` |
| `external_id` | SMSPILOT `server_id` |
| `external_status`, `cost`, `error` | receipt details |
| `created_at`, `completed_at` | timing, drives `repeat_interval` later |

v1 always writes `sequence = 1`; the column exists so repeats can be added
without reshaping the unique index (awkward in SQLite).

The row is inserted **before** the provider call (claim), then updated with the
result. A crash between request and response leaves the row in `unknown`, which
is never retried — a call may well have been placed.

**Dispatch point:** a new step in `ProcessCheckResult`, **after**
`manageIncident`, driven by the ongoing incident — not inside `CreateAlert`,
which sits behind the hardcoded threshold and the 5-minute cooldown (constraints
2 and 3 above).

**Threshold source:** `incident.FailureCount`, not `state.consecutiveFailures`
(constraint 4) — warnings must not place phone calls, and the DB counter
survives restarts.

## Stage 3 — SMSPILOT notifier

New `internal/notifier/smspilot.go`, modelled on `telegram.go`, reusing
`buildHTTPClient` from `http_client.go`.

- Config: `phones`, `voice` (`GOLOS` / `GOLOSM`), `message_template`, `ttl`,
  `api_key` or `api_key_env`.
- Add `api_key` to `sensitiveKeys` in `maskNotifierConfig`.
- Register the type in **both** factory switches: `cmd/server/main.go:214-228`
  and `buildNotifier` in `notifier_handlers.go:206-221`.
- `domain.Notifier.Notify` returns only `error`, which cannot carry the
  `server_id` or cost. Add an optional interface checked by type assertion, so
  existing notifiers are untouched:

  ```go
  type ReceiptNotifier interface {
      NotifyWithReceipt(ctx context.Context, alert *Alert) (Receipt, error)
  }
  ```

- Error classification: an explicit provider rejection is `failed`; a timeout,
  non-200 response or unparseable body is a typed `ErrUnknownResult` and is
  **never** auto-retried.
- Known provider facts: statuses `0` queued, `1` handed to operator,
  `2` delivered, `3` deferred, `-1` not delivered, `-2` error; error `232` means
  the number must be whitelisted or a business tariff enabled. Take the request
  contract (`to`, `send`, `from`, `apikey`, `format`) from the SMSPILOT docs —
  the Python snippet circulated in chat does not run (positional argument after
  keyword arguments in `_request`).
- The API key stays server-side: prefer `api_key_env`; never log the request
  body.
- **The endpoint URL must be an injectable struct field** (default
  `https://smspilot.ru/api.php`), not a package const. `telegram.go:17` hardcodes
  its URL as a const, which is exactly why `telegram_test.go` never exercises the
  send path; `webhook.go` takes the URL from config and is testable. Here a
  hardcoded const would mean any send-path test places a real phone call.

## Stage 4 — dashboard

In `web/static/js/app.js`: add the type to `NOTIFIER_TYPES` (lines 1067-1073), a
branch in `typeFields`, and a branch in the payload builder (lines 1180-1225).

Below the per-type fields, a shared "Delivery rules" block for every notifier
type: notify after N failures, how many times, which targets. Defaults: 3
failures, 1 call, no repeats. The repeat-interval field ships disabled/hidden
until repeats are implemented.

## Stage 5 — tests

**A test that reaches the real provider rings someone's phone and costs money.
Nothing in `go test ./...` may ever touch `smspilot.ru`.** Rules:

- The SMSPILOT notifier is constructed in tests with `api_url` pointed at an
  `httptest.NewServer` — the same pattern as `webhook_test.go:130`. This is only
  possible because of the injectable URL field in Stage 3.
- Tests use literal fake credentials (`"test-key"`, `"79000000000"`). No test
  reads `SMSPILOT_API_KEY` or `ALERT_PHONE` from the environment — a developer
  machine with a real key exported must not turn a unit test into a call.
- The dispatch tests below (threshold, dedupe, restart) use an in-memory fake
  notifier that records calls, not the SMSPILOT implementation at all — no HTTP
  involved.
- No live smoke test in CI. If a one-off real call is ever needed to verify the
  account, it goes in a separate file behind `//go:build manual`, takes the key
  and phone from the environment, and is run by hand.
- Do not add a "Send test notification" button for the voice type in the
  dashboard (none exists for any type today). If one is ever added, it must
  require an explicit typed confirmation for voice.

Cases:

- No call on the 1st and 2nd failure; call on the 3rd.
- Recovery before the threshold cancels the call entirely.
- A `Warning` (slow) result does not count towards the threshold.
- No duplicate call for the same incident across repeated ticks.
- Restart mid-incident: no second call, old incident resolves on recovery.
- `ErrUnknownResult` is not re-sent.
- `api_key` never appears in API responses or logs.

---

## Compatibility and known ceilings

- A notifier config without policy keys behaves exactly as today; Telegram and
  email are not routed through the new dispatch path.
- Event-type filtering (`down` / `recovery` / `slow` / SSL) is out of v1 — a
  voice call only makes sense on incident escalation.
- Dispatch is tied to scheduler ticks, so repeat granularity (once implemented)
  equals the target's check interval, and a target disabled or deleted mid-
  incident stops producing calls.
- Delivery-status polling (SMSPILOT `check` ~5 minutes after the call) is out of
  v1: it needs its own timer, and status `2` means "delivered", not "a human
  picked it up".
- Acknowledgement ("Accepted" button, `acknowledged_at/by`) is only needed once
  repeats exist — nothing to stop while v1 calls once.
