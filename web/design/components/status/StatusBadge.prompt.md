Status primitives — dot, badge, tick bar, sparkline, live indicator. Status is always icon/shape + text, never color alone.

```jsx
<StatusDot status="down" pulse />
<StatusBadge status="up" label="Operational" />
<TickBar checks={[{ s: "up", tip: "12:04 · 187 ms" }, { s: "down", tip: "12:05 · timeout" }]} slots={45} />
<Sparkline points={[120, 180, 140, 460, 130]} />
<LiveIndicator state="live" />
```

Canonical statuses: `up | down | degraded | unknown` (map API success/failure/warning/unknown to these).
