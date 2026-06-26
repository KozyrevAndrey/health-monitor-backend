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
