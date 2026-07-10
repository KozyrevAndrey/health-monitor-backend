import React from "react";
import { Icon } from "../core/Icon.jsx";
import { StatusDot } from "../status/StatusDot.jsx";
import { TickBar } from "../status/TickBar.jsx";

const TYPE_ICON = { http: "globe", tcp: "server", dns: "at-sign" };

/**
 * Target card (comfortable grid view): status dot + name + type icon,
 * tick bar of recent checks, uptime % + last response time.
 * A DOWN target gets a red border + pulsing dot.
 */
export function TargetCard({ target, onOpen, onEdit, onDelete, style }) {
  const { name, type = "http", status = "unknown", checks = [], uptime, lastMs, endpoint, paused } = target;
  const down = status === "down";
  return (
    <div
      className={`hm-card hm-card--hover${down ? " hm-card--down" : ""}`}
      style={{ padding: "var(--space-3) var(--space-4)", display: "flex", flexDirection: "column", gap: "var(--space-2)", ...style }}
      onClick={onOpen}
      role="button"
      tabIndex={0}
      onKeyDown={(e) => { if (e.key === "Enter" && onOpen) onOpen(); }}
      aria-label={`${name}, ${paused ? "paused" : status}`}
    >
      <div style={{ display: "flex", alignItems: "center", gap: "var(--space-2)" }}>
        <StatusDot status={paused ? "unknown" : status} pulse={down} />
        <b style={{ fontSize: "var(--text-lg)", fontWeight: "var(--weight-semibold)", overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>{name}</b>
        <Icon name={TYPE_ICON[type] || "globe"} size={13} style={{ color: "var(--text-3)" }} data-tip={type.toUpperCase()} />
        <span style={{ flex: 1 }}></span>
        <div style={{ display: "flex", gap: 2 }} onClick={(e) => e.stopPropagation()}>
          {onEdit && (
            <button type="button" className="hm-iconbtn" aria-label={`Edit ${name}`} onClick={onEdit}><Icon name="pencil" size={14} /></button>
          )}
          {onDelete && (
            <button type="button" className="hm-iconbtn hm-iconbtn--danger" aria-label={`Delete ${name}`} onClick={onDelete}><Icon name="trash" size={14} /></button>
          )}
        </div>
      </div>
      {endpoint && (
        <span className="num" style={{ fontSize: "var(--text-sm)", color: "var(--text-3)", overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>{endpoint}</span>
      )}
      <TickBar checks={checks} slots={45} height={18} />
      <div style={{ display: "flex", alignItems: "baseline", gap: "var(--space-3)", fontSize: "var(--text-sm)", color: "var(--text-2)" }}>
        <span className="num" style={{ color: down ? "var(--status-down)" : "var(--status-up)", fontWeight: "var(--weight-semibold)" }}>
          {paused ? "—" : `${uptime}%`}
        </span>
        <span className="num">{lastMs != null ? `${lastMs} ms` : "—"}</span>
        {paused && <span style={{ color: "var(--text-3)", display: "inline-flex", alignItems: "center", gap: 4 }}><Icon name="pause" size={11} />paused</span>}
      </div>
    </div>
  );
}
