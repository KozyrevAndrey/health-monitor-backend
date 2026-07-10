// Health Monitor UI kit — deterministic sample data.
// Scenario: "ok" (all up) or "alarm" (2 down, 1 degraded). window.HM_SCENARIO picks it.
window.HMData = (() => {
  function rng(seed) {
    let s = seed;
    return () => ((s = (s * 1103515245 + 12345) % 2147483648) / 2147483648);
  }

  const defs = [
    { id: "t1", name: "API /health", type: "http", endpoint: "https://api.example.dev/health", interval: "30s", timeout: "5s", tags: ["prod", "api"], desc: "Main REST API liveness probe", base: 180 },
    { id: "t2", name: "Marketing site", type: "http", endpoint: "https://example.dev", interval: "1m", timeout: "10s", tags: ["prod", "web"], base: 320, alarm: "degraded" },
    { id: "t3", name: "Postgres primary", type: "tcp", endpoint: "10.0.4.2:5432", interval: "30s", timeout: "3s", tags: ["prod", "db"], base: 12, alarm: "down", err: "dial tcp 10.0.4.2:5432: i/o timeout" },
    { id: "t4", name: "Redis cache", type: "tcp", endpoint: "10.0.4.3:6379", interval: "30s", timeout: "3s", tags: ["prod", "cache"], base: 8 },
    { id: "t5", name: "Backup worker", type: "http", endpoint: "https://backup.example.dev/ping", interval: "5m", timeout: "10s", tags: ["infra"], base: 260, alarm: "down", err: "HTTP 503 Service Unavailable" },
    { id: "t6", name: "example.dev DNS", type: "dns", endpoint: "example.dev · A", interval: "5m", timeout: "5s", tags: ["infra", "dns"], base: 45 },
    { id: "t7", name: "Staging API", type: "http", endpoint: "https://staging.example.dev/health", interval: "1m", timeout: "5s", tags: ["staging"], base: 210 },
    { id: "t8", name: "SMTP relay", type: "tcp", endpoint: "10.0.4.9:25", interval: "5m", timeout: "5s", tags: ["infra", "mail"], base: 30 },
    { id: "t9", name: "Status page", type: "http", endpoint: "https://status.example.dev", interval: "5m", timeout: "10s", tags: ["web"], base: 190, paused: true },
    { id: "t10", name: "Grafana", type: "http", endpoint: "https://grafana.internal.example.dev", interval: "1m", timeout: "10s", tags: ["infra", "observability"], base: 150 },
    { id: "t11", name: "MX records", type: "dns", endpoint: "example.dev · MX", interval: "15m", timeout: "5s", tags: ["infra", "dns", "mail"], base: 52 },
    { id: "t12", name: "CDN edge", type: "http", endpoint: "https://cdn.example.dev/health", interval: "1m", timeout: "5s", tags: ["prod", "web"], base: 95 },
  ];

  function pad(n) { return String(n).padStart(2, "0"); }

  function build(def, i, scenario) {
    const r = rng(i * 7919 + 13);
    const state = scenario === "alarm" ? def.alarm || "up" : "up";
    const ms = () => Math.round(def.base * (0.7 + r() * 0.9));

    // 45 recent ticks (~last 45 checks), newest right
    const checks = [];
    for (let k = 0; k < 45; k++) {
      const t = `12:${pad(k + 10)}`;
      let s = r() > 0.985 ? "degraded" : "up";
      if (state === "down" && k >= 41) s = "down";
      if (state === "degraded" && (k === 38 || k >= 43)) s = "degraded";
      checks.push({ s, tip: s === "down" ? `${t} · ${def.err || "failed"}` : `${t} · ${ms()} ms${s === "degraded" ? " · slow" : ""}` });
    }

    // per-period stats + chart series
    const periods = {};
    [["24h", 48, (k) => `${pad(Math.floor(k / 2))}:${k % 2 ? "30" : "00"}`],
     ["7d", 56, (k) => ["Thu","Fri","Sat","Sun","Mon","Tue","Wed"][Math.floor(k / 8)]],
     ["30d", 60, (k) => `${pad(9 + Math.floor(k / 2) - 30 + 21)} ${Math.floor(k/2) < 21 ? "Jun" : "Jul"}`],
    ].forEach(([p, n, lbl]) => {
      const rr = rng(i * 104729 + p.length * 31);
      const series = [];
      let failed = 0;
      for (let k = 0; k < n; k++) {
        let s = rr() > 0.99 ? "degraded" : "up";
        let v = Math.round(def.base * (0.7 + rr() * 0.9));
        if (state === "down" && k >= n - 3) { s = "down"; v = null; failed++; }
        if (state === "degraded" && k >= n - 4) { s = "degraded"; v = Math.round(def.base * 3.2); }
        series.push({ v, s, label: lbl(k) });
      }
      const vals = series.filter((x) => x.v != null).map((x) => x.v);
      const secs = { "24h": 86400, "7d": 604800, "30d": 2592000 }[p];
      const total = Math.round(secs / ({ "30s": 30, "1m": 60, "5m": 300, "15m": 900 }[def.interval] || 60));
      let failTotal = state === "down" ? (p === "24h" ? 14 : p === "7d" ? 14 : 22)
        : state === "degraded" ? (p === "24h" ? 2 : 3)
        : p === "24h" ? (rr() < 0.4 ? 1 : 0)
        : Math.round(rr() * (p === "7d" ? 4 : 7));
      failTotal = Math.min(failTotal, Math.max(total - 1, 0));
      periods[p] = {
        series,
        uptime: (100 * (1 - failTotal / total)).toFixed(p === "24h" ? 2 : 3),
        avg: Math.round(vals.reduce((a, b) => a + b, 0) / vals.length),
        min: Math.min.apply(null, vals),
        max: Math.max.apply(null, vals),
        totalChecks: total,
        failedChecks: failTotal,
      };
    });

    const status = def.paused ? "unknown" : state;
    return Object.assign({}, def, {
      status,
      paused: !!def.paused,
      checks,
      periods,
      uptime: periods["24h"].uptime,
      lastMs: state === "down" ? null : checks[44].tip.match(/(\d+) ms/) ? Number(checks[44].tip.match(/(\d+) ms/)[1]) : def.base,
      consecutiveFailures: state === "down" ? 14 : 0,
      lastCheck: "12s ago",
    });
  }

  const notifiers = [
    { id: "n1", name: "Ops channel", type: "telegram", enabled: true, detail: "chat •••4821" },
    { id: "n2", name: "On-call email", type: "email", enabled: true, detail: "oncall@•••.dev via smtp.example.dev:587" },
    { id: "n3", name: "Founders Gmail", type: "gmail_oauth", enabled: false, detail: "alerts@•••.dev" },
    { id: "n4", name: "PagerDuty bridge", type: "webhook", enabled: true, detail: "https://events.pagerduty.com/•••" },
  ];

  function incidents(scenario) {
    const ongoing = scenario === "alarm" ? [
      { id: "i1", targetId: "t3", target: "Postgres primary", status: "ongoing", severity: "critical", startedAt: "Today, 09:14", duration: "23m", failureCount: 14, lastError: "dial tcp 10.0.4.2:5432: i/o timeout" },
      { id: "i2", targetId: "t5", target: "Backup worker", status: "ongoing", severity: "warning", startedAt: "Today, 09:29", duration: "8m", failureCount: 2, lastError: "HTTP 503 Service Unavailable" },
    ] : [];
    const resolved = [
      { id: "i3", targetId: "t2", target: "Marketing site", status: "resolved", severity: "warning", startedAt: "8 Jul, 22:03", duration: "4m", failureCount: 3, lastError: "HTTP 502 Bad Gateway" },
      { id: "i4", targetId: "t12", target: "CDN edge", status: "resolved", severity: "critical", startedAt: "6 Jul, 03:41", duration: "18m", failureCount: 11, lastError: "context deadline exceeded" },
      { id: "i5", targetId: "t6", target: "example.dev DNS", status: "resolved", severity: "info", startedAt: "1 Jul, 11:02", duration: "2m", failureCount: 1, lastError: "SERVFAIL" },
    ];
    return { ongoing, resolved };
  }

  return function make(scenario) {
    return {
      scenario,
      targets: defs.map((d, i) => build(d, i, scenario)),
      notifiers,
      incidents: incidents(scenario),
    };
  };
})();
