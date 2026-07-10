// Target detail: period-driven stats + chart + incidents + checks + config.
const DS = window.HealthMonitorRedesign_502b6f;
const { Icon, StatusDot, StatusBadge, Segmented, Button, TickBar } = DS;

const detailTypeIcon = { http: "globe", tcp: "server", dns: "at-sign" };

/* ---- SVG response-time chart, points colored by status ---- */
function ResponseChart({ series, height = 170 }) {
  const w = 640, padL = 44, padB = 22, padT = 10;
  const vals = series.filter((p) => p.v != null).map((p) => p.v);
  const max = Math.max(...vals) * 1.15, min = 0;
  const iw = w - padL - 8, ih = height - padT - padB;
  const x = (i) => padL + (i / (series.length - 1)) * iw;
  const y = (v) => padT + ih * (1 - (v - min) / (max - min));
  const segs = [];
  let cur = [];
  series.forEach((p, i) => {
    if (p.v == null) { if (cur.length > 1) segs.push(cur); cur = []; }
    else cur.push(`${x(i).toFixed(1)},${y(p.v).toFixed(1)}`);
  });
  if (cur.length > 1) segs.push(cur);
  const gridVals = [0.25, 0.5, 0.75, 1].map((f) => Math.round(max * f));
  const labelEvery = Math.ceil(series.length / 6);
  return (
    <svg width="100%" viewBox={`0 0 ${w} ${height}`} role="img" aria-label="Response time chart" style={{ display: "block" }}>
      {gridVals.map((v) => (
        <g key={v}>
          <line x1={padL} x2={w - 8} y1={y(v)} y2={y(v)} stroke="var(--chart-grid)" strokeWidth="1" />
          <text x={padL - 6} y={y(v) + 3} textAnchor="end" fontSize="9" fill="var(--text-3)" fontFamily="var(--font-mono)">{v}</text>
        </g>
      ))}
      {segs.map((s, i) => (
        <polyline key={i} points={s.join(" ")} fill="none" stroke="var(--chart-line)" strokeWidth="1.5" strokeLinejoin="round" />
      ))}
      {series.map((p, i) => {
        if (p.v == null) {
          return <rect key={i} x={x(i) - 2} y={padT} width={4} height={ih} fill="var(--status-down-subtle)" stroke="var(--status-down)" strokeWidth="0.5" rx={1} />;
        }
        const bad = p.s !== "up";
        if (!bad && i % 2) return null;
        return <circle key={i} cx={x(i)} cy={y(p.v)} r={bad ? 3 : 1.6} fill={bad ? "var(--status-degraded)" : "var(--chart-line)"} />;
      })}
      {series.map((p, i) => (i % labelEvery === 0 ? (
        <text key={`l${i}`} x={x(i)} y={height - 6} textAnchor="middle" fontSize="9" fill="var(--text-3)" fontFamily="var(--font-mono)">{p.label}</text>
      ) : null))}
    </svg>
  );
}

function DetailStat({ label, value, unit, tone }) {
  return (
    <div style={{ flex: 1, minWidth: 0, background: "var(--bg-inset)", border: "1px solid var(--border-1)", borderRadius: "var(--radius-md)", padding: "8px 12px" }}>
      <span className="overline" style={{ display: "block" }}>{label}</span>
      <b className="num" style={{ fontSize: "var(--text-2xl)", fontWeight: 600, color: tone || "var(--text-1)" }}>
        {value}{unit && <span style={{ fontSize: "var(--text-sm)", color: "var(--text-3)", marginLeft: 2 }}>{unit}</span>}
      </b>
    </div>
  );
}

function ConfigRow({ k, v, mono }) {
  return (
    <div style={{ display: "grid", gridTemplateColumns: "140px 1fr", gap: 8, padding: "6px 0", borderBottom: "1px solid var(--border-1)", fontSize: "var(--text-md)" }}>
      <span style={{ color: "var(--text-3)" }}>{k}</span>
      <span className={mono ? "num" : ""} style={{ color: "var(--text-1)", overflowWrap: "anywhere" }}>{v}</span>
    </div>
  );
}

function ChecksTable({ target }) {
  const rows = target.checks.slice(-8).reverse().map((c, i) => {
    const m = c.tip.match(/^(\S+) · (.*)$/);
    return { time: m ? m[1] : "", info: m ? m[2] : c.tip, s: c.s };
  });
  return (
    <table className="hm-table">
      <thead><tr><th>Time</th><th>Result</th><th style={{ textAlign: "right" }}>Detail</th></tr></thead>
      <tbody>
        {rows.map((r, i) => (
          <tr key={i} style={{ cursor: "default" }}>
            <td className="num" style={{ color: "var(--text-2)" }}>{r.time}</td>
            <td>
              <span style={{ display: "inline-flex", alignItems: "center", gap: 6, color: r.s === "down" ? "var(--status-down)" : r.s === "degraded" ? "var(--status-degraded)" : "var(--status-up)", fontSize: "var(--text-sm)", fontWeight: 500 }}>
                <Icon name={r.s === "down" ? "x-circle" : r.s === "degraded" ? "zap" : "check-circle"} size={13} />
                {r.s === "down" ? "Failure" : r.s === "degraded" ? "Slow" : "Success"}
              </span>
            </td>
            <td className="num" style={{ textAlign: "right", color: "var(--text-2)", fontSize: "var(--text-sm)" }}>{r.info}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

function Section({ title, children, extra }) {
  return (
    <section style={{ marginTop: "var(--space-5)" }}>
      <div style={{ display: "flex", alignItems: "center", gap: 8, marginBottom: "var(--space-2)" }}>
        <span className="overline">{title}</span>
        <span style={{ flex: 1 }}></span>
        {extra}
      </div>
      {children}
    </section>
  );
}

/* ---- shared detail body (used by slide-over AND full-page variant) ---- */
function TargetDetailBody({ target, incidents, onEdit }) {
  const [period, setPeriod] = React.useState("24h");
  const p = target.periods[period];
  const down = target.status === "down";
  const tIncidents = incidents.ongoing.concat(incidents.resolved).filter((i) => i.targetId === target.id);
  return (
    <div>
      <div style={{ display: "flex", alignItems: "center", gap: 8, flexWrap: "wrap" }}>
        <Segmented options={["24h", "7d", "30d"]} value={period} onChange={setPeriod} ariaLabel="Stats period" />
        <span style={{ flex: 1 }}></span>
        <span className="hm-live" data-state="live"><span className="hm-dot"></span>Checked {target.lastCheck} · every {target.interval}</span>
      </div>

      <div style={{ display: "flex", gap: 8, marginTop: "var(--space-3)", flexWrap: "wrap" }}>
        <DetailStat label="Uptime" value={p.uptime} unit="%" tone={down ? "var(--status-down)" : "var(--status-up)"} />
        <DetailStat label="Avg" value={p.avg} unit="ms" />
        <DetailStat label="Min / Max" value={`${p.min}–${p.max}`} unit="ms" />
        <DetailStat label="Failed" value={p.failedChecks} unit={`/ ${p.totalChecks}`} tone={p.failedChecks ? "var(--status-down)" : undefined} />
      </div>

      <Section title={`Response time — ${period}`}>
        <div style={{ background: "var(--bg-inset)", border: "1px solid var(--border-1)", borderRadius: "var(--radius-md)", padding: "10px 8px 4px" }}>
          <ResponseChart series={p.series} />
        </div>
      </Section>

      <Section title="Recent checks">
        <TickBar checks={target.checks} slots={45} height={16} style={{ marginBottom: 8 }} />
        <div className="hm-card" style={{ overflow: "hidden", borderRadius: "var(--radius-md)" }}>
          <ChecksTable target={target} />
        </div>
      </Section>

      <Section title="Incidents">
        {tIncidents.length ? (
          <div className="hm-card" style={{ borderRadius: "var(--radius-md)" }}>
            {tIncidents.map((inc, i) => <DS.IncidentItem key={inc.id} incident={inc} last={i === tIncidents.length - 1} />)}
          </div>
        ) : (
          <p style={{ margin: 0, fontSize: "var(--text-md)", color: "var(--text-3)" }}>No incidents in the last 30 days.</p>
        )}
      </Section>

      <Section title="Configuration" extra={<Button variant="ghost" icon={<Icon name="pencil" size={13} />} onClick={onEdit}>Edit</Button>}>
        <div style={{ background: "var(--bg-inset)", border: "1px solid var(--border-1)", borderRadius: "var(--radius-md)", padding: "4px 14px" }}>
          <ConfigRow k="Type" v={target.type.toUpperCase()} />
          <ConfigRow k={target.type === "http" ? "URL" : target.type === "tcp" ? "Host : port" : "Domain · record"} v={target.endpoint} mono />
          <ConfigRow k="Interval" v={target.interval} mono />
          <ConfigRow k="Timeout" v={target.timeout} mono />
          <ConfigRow k="Tags" v={target.tags.join(", ")} />
          {target.desc && <ConfigRow k="Description" v={target.desc} />}
          <div style={{ display: "grid", gridTemplateColumns: "140px 1fr", gap: 8, padding: "6px 0", fontSize: "var(--text-md)" }}>
            <span style={{ color: "var(--text-3)" }}>Enabled</span>
            <span style={{ color: target.paused ? "var(--text-3)" : "var(--status-up)" }}>{target.paused ? "No — paused" : "Yes"}</span>
          </div>
        </div>
      </Section>
    </div>
  );
}

function DetailTitle({ target }) {
  return (
    <>
      <StatusDot status={target.status} pulse={target.status === "down"} size="lg" />
      <span style={{ overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>{target.name}</span>
      <Icon name={detailTypeIcon[target.type]} size={14} style={{ color: "var(--text-3)" }} />
    </>
  );
}

Object.assign(window, { TargetDetailBody, DetailTitle, ResponseChart });
