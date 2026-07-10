Live-data building blocks: StatCard, TargetCard (grid view), TargetTable (dense view), IncidentItem, NotifierRow.

```jsx
<StatCard label="Failing" value={2} tone="danger" sub="of 12 targets" />
<TargetCard target={{ name: "API", type: "http", status: "up", checks, uptime: 99.98, lastMs: 187, endpoint: "https://api.example.dev/health" }} onOpen={…} />
<TargetTable targets={[…]} onOpen={…} onEdit={…} onDelete={…} />
<IncidentItem incident={{ target: "API", status: "ongoing", severity: "critical", startedAt: "Started 09:14", duration: "23m", failureCount: 14, lastError: "dial tcp: i/o timeout" }} />
<NotifierRow notifier={{ name: "Ops channel", type: "telegram", enabled: true, detail: "chat •••4821" }} />
```

DOWN states escalate visually: red card border + pulsing dot + tinted row — never just a recolored pill.
