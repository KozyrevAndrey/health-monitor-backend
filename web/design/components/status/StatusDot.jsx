import React from "react";

/**
 * Small colored status dot. `pulse` adds a radiating ring — reserved for
 * ongoing-down states so motion means "needs attention".
 */
export function StatusDot({ status = "unknown", size = "md", pulse = false, style }) {
  return (
    <span
      className={`hm-dot${size === "lg" ? " hm-dot--lg" : ""}${pulse ? " hm-dot--pulse" : ""}`}
      data-status={status}
      style={style}
    ></span>
  );
}
