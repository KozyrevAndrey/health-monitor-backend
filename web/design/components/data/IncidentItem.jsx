import React from "react";
import { Icon } from "../core/Icon.jsx";

const SEV = {
  critical: { color: "var(--status-down)", icon: "alert-triangle" },
  warning: { color: "var(--status-degraded)", icon: "zap" },
  info: { color: "var(--accent)", icon: "info" },
};

/**
 * Incident row for timeline lists. Ongoing incidents show a pulsing marker
 * and red tint; resolved ones a check + muted rail.
 */
export function IncidentItem({ incident, last = false, onOpen, style }) {
  const { target, status = "resolved", severity = "critical", startedAt, duration, failureCount, lastError } = incident;
  const ongoing = status === "ongoing";
  const sev = SEV[severity] || SEV.critical;
  return (
    <div className="hm-incident" style={{ background: ongoing ? "var(--status-down-subtle)" : "transparent", cursor: onOpen ? "pointer" : "default", ...style }} onClick={onOpen}>
      <div className="hm-incident-rail" aria-hidden="true">
        <span className={`hm-dot${ongoing ? " hm-dot--pulse" : ""}`} data-status={ongoing ? "down" : "up"}></span>
        {!last && <span className="line"></span>}
      </div>
      <div style={{ flex: 1, minWidth: 0 }}>
        <div style={{ display: "flex", alignItems: "center", gap: "var(--space-2)", flexWrap: "wrap" }}>
          <b style={{ fontSize: "var(--text-md)", fontWeight: "var(--weight-semibold)" }}>{target}</b>
          <span style={{ display: "inline-flex", alignItems: "center", gap: 4, fontSize: "var(--text-xs)", fontWeight: "var(--weight-semibold)", color: sev.color, textTransform: "uppercase", letterSpacing: "var(--tracking-caps)" }}>
            <Icon name={sev.icon} size={11} />{severity}
          </span>
          <span style={{ flex: 1 }}></span>
          <span className="num" style={{ fontSize: "var(--text-sm)", color: ongoing ? "var(--status-down)" : "var(--text-3)", fontWeight: ongoing ? "var(--weight-semibold)" : "var(--weight-regular)" }}>
            {ongoing ? `ongoing · ${duration}` : duration}
          </span>
        </div>
        <div style={{ fontSize: "var(--text-sm)", color: "var(--text-2)", marginTop: 2 }}>
          {startedAt}{failureCount != null && <> · <span className="num">{failureCount}</span> failed checks</>}
        </div>
        {lastError && (
          <code style={{ display: "block", marginTop: 4, fontSize: "var(--text-xs)", color: "var(--text-3)", overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>{lastError}</code>
        )}
      </div>
    </div>
  );
}
