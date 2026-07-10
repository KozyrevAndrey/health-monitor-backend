import React from "react";
import { Icon } from "../core/Icon.jsx";
import { StatusDot } from "../status/StatusDot.jsx";
import { TickBar } from "../status/TickBar.jsx";

const TYPE_ICON = { http: "globe", tcp: "server", dns: "at-sign" };

/**
 * Compact table of targets (dense alternative to the card grid).
 * Columns: status, name, type, tick bar, uptime %, last ms, actions.
 */
export function TargetTable({ targets = [], onOpen, onEdit, onDelete, style }) {
  return (
    <table className="hm-table" style={style}>
      <thead>
        <tr>
          <th style={{ width: 24 }}><span style={{ position: "absolute", width: 1, height: 1, overflow: "hidden", clip: "rect(0 0 0 0)" }}>Status</span></th>
          <th>Name</th>
          <th>Type</th>
          <th style={{ width: 200 }}>Last 45 checks</th>
          <th style={{ textAlign: "right" }}>Uptime</th>
          <th style={{ textAlign: "right" }}>Resp.</th>
          <th style={{ width: 76 }}></th>
        </tr>
      </thead>
      <tbody>
        {targets.map((t, i) => {
          const down = t.status === "down";
          return (
            <tr key={t.id || i} data-status={t.paused ? "unknown" : t.status} onClick={() => onOpen && onOpen(t)} tabIndex={0}
                onKeyDown={(e) => { if (e.key === "Enter" && onOpen) onOpen(t); }}>
              <td><StatusDot status={t.paused ? "unknown" : t.status} pulse={down} /></td>
              <td>
                <b style={{ fontWeight: "var(--weight-semibold)" }}>{t.name}</b>
                {t.paused && <span style={{ marginLeft: 8, fontSize: "var(--text-xs)", color: "var(--text-3)" }}>paused</span>}
              </td>
              <td>
                <span style={{ display: "inline-flex", alignItems: "center", gap: 5, color: "var(--text-2)", fontSize: "var(--text-sm)" }}>
                  <Icon name={TYPE_ICON[t.type] || "globe"} size={13} />{(t.type || "http").toUpperCase()}
                </span>
              </td>
              <td><TickBar checks={t.checks || []} slots={45} height={14} /></td>
              <td className="num" style={{ textAlign: "right", color: down ? "var(--status-down)" : "var(--status-up)", fontWeight: "var(--weight-medium)" }}>
                {t.paused ? "—" : `${t.uptime}%`}
              </td>
              <td className="num" style={{ textAlign: "right", color: "var(--text-2)" }}>{t.lastMs != null ? `${t.lastMs} ms` : "—"}</td>
              <td onClick={(e) => e.stopPropagation()}>
                <div style={{ display: "flex", gap: 2, justifyContent: "flex-end" }}>
                  {onEdit && <button type="button" className="hm-iconbtn" aria-label={`Edit ${t.name}`} onClick={() => onEdit(t)}><Icon name="pencil" size={14} /></button>}
                  {onDelete && <button type="button" className="hm-iconbtn hm-iconbtn--danger" aria-label={`Delete ${t.name}`} onClick={() => onDelete(t)}><Icon name="trash" size={14} /></button>}
                </div>
              </td>
            </tr>
          );
        })}
      </tbody>
    </table>
  );
}
