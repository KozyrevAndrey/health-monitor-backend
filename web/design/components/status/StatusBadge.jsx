import React from "react";
import { Icon } from "../core/Icon.jsx";

const MAP = {
  up: { icon: "check", label: "Up" },
  down: { icon: "alert-triangle", label: "Down" },
  degraded: { icon: "zap", label: "Degraded" },
  unknown: { icon: "pause", label: "Unknown" },
};

/**
 * Status pill: icon + text, never color alone. Pass `label` to override the
 * default text (e.g. "Operational", "Paused", "3 failing").
 */
export function StatusBadge({ status = "unknown", label, style }) {
  const m = MAP[status] || MAP.unknown;
  return (
    <span className="hm-badge" data-status={status} style={style}>
      <Icon name={m.icon} size={12} />
      {label || m.label}
    </span>
  );
}
