import React from "react";

const LABEL = { live: "Live", reconnecting: "Reconnecting…", polling: "Polling (30s)" };

/** SSE connection indicator: live / reconnecting / polling fallback. */
export function LiveIndicator({ state = "live", label, style }) {
  return (
    <span className="hm-live" data-state={state} style={style} role="status">
      <span className="hm-dot"></span>
      {label || LABEL[state] || state}
    </span>
  );
}
