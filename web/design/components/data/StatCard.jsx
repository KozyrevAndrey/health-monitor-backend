import React from "react";

/**
 * Compact metric card. `tone` tints the value (e.g. danger for failing count).
 */
export function StatCard({ label, value, unit, sub, tone, icon, style }) {
  const color = tone === "danger" ? "var(--status-down)"
    : tone === "success" ? "var(--status-up)"
    : tone === "warning" ? "var(--status-degraded)"
    : "var(--text-1)";
  return (
    <div className="hm-card hm-statcard" style={style}>
      <span className="overline">{icon}{label}</span>
      <b style={{ color }}>
        {value}
        {unit && <span style={{ fontSize: "var(--text-lg)", fontWeight: "var(--weight-medium)", color: "var(--text-3)", marginLeft: 3 }}>{unit}</span>}
      </b>
      {sub && <small>{sub}</small>}
    </div>
  );
}
