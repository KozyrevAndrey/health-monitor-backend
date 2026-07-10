import React from "react";

/**
 * UptimeRobot-style tick bar of the last N checks (oldest left, newest right).
 * Each check: { s: "up"|"down"|"degraded"|"unknown", tip?: string }.
 * Pads with empty ticks on the left when fewer than `slots` checks exist.
 */
export function TickBar({ checks = [], slots = 45, height = 20, style, ariaLabel }) {
  const pad = Math.max(0, slots - checks.length);
  const shown = checks.slice(-slots);
  const downs = shown.filter((c) => c.s === "down").length;
  return (
    <div
      className="hm-ticks"
      style={{ height, ...style }}
      role="img"
      aria-label={ariaLabel || `Last ${shown.length} checks, ${downs} failed`}
    >
      {Array.from({ length: pad }).map((_, i) => (
        <span key={`p${i}`} className="hm-tick"></span>
      ))}
      {shown.map((c, i) => (
        <span key={i} className="hm-tick" data-s={c.s} data-tip={c.tip || undefined}></span>
      ))}
    </div>
  );
}
