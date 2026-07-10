import React from "react";
import { Icon } from "../core/Icon.jsx";

/** Designed empty state: icon, one-line title, short body, primary action. */
export function EmptyState({ icon = "activity", title, body, action, style }) {
  return (
    <div className="hm-empty" style={style}>
      <Icon name={icon} size={28} />
      <b>{title}</b>
      {body && <p>{body}</p>}
      {action}
    </div>
  );
}
