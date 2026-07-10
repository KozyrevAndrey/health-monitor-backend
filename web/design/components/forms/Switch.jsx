import React from "react";

/** Toggle switch (green = enabled). Renders an optional side label. */
export function Switch({ checked = false, onChange, label, disabled = false, style }) {
  return (
    <label style={{ display: "inline-flex", alignItems: "center", gap: "var(--space-2)", cursor: disabled ? "not-allowed" : "pointer", opacity: disabled ? 0.5 : 1, ...style }}>
      <span className="hm-switch">
        <input
          type="checkbox"
          role="switch"
          checked={checked}
          disabled={disabled}
          onChange={(e) => onChange && onChange(e.target.checked)}
        />
        <i></i>
      </span>
      {label && <span style={{ fontSize: "var(--text-md)", color: "var(--text-1)" }}>{label}</span>}
    </label>
  );
}
