import React from "react";
import { Icon } from "../core/Icon.jsx";

const KIND_ICON = { success: "check-circle", danger: "alert-triangle", warning: "zap", info: "info" };

/** Single toast. Compose inside a `.hm-toaststack` fixed container. */
export function Toast({ kind = "info", title, message, onDismiss, style }) {
  return (
    <div className="hm-toast" data-kind={kind} role="status" style={style}>
      <Icon name={KIND_ICON[kind] || "info"} size={16} />
      <div style={{ flex: 1, minWidth: 0 }}>
        <b>{title}</b>
        {message && <p>{message}</p>}
      </div>
      {onDismiss && (
        <button type="button" className="hm-iconbtn" style={{ width: 24, height: 24, margin: "-4px -6px 0 0" }} aria-label="Dismiss" onClick={onDismiss}>
          <Icon name="x" size={13} />
        </button>
      )}
    </div>
  );
}
