import React from "react";

const PRESETS = ["30s", "1m", "5m", "15m", "1h"];

/**
 * Duration input: preset chips + free-form field ("30s", "1m", "2h").
 * Value is a Go-style duration string.
 */
export function DurationField({ label, value = "1m", onChange, presets = PRESETS, hint, style }) {
  return (
    <div className="hm-field" style={style}>
      {label && <span className="hm-label">{label}</span>}
      <div style={{ display: "flex", gap: "var(--space-2)", alignItems: "center" }}>
        <div className="hm-seg" role="group" aria-label={`${label || "Duration"} presets`}>
          {presets.map((p) => (
            <button key={p} type="button" className="hm-seg-btn num" aria-pressed={value === p} onClick={() => onChange && onChange(p)}>
              {p}
            </button>
          ))}
        </div>
        <input
          className="hm-input hm-input--mono"
          style={{ width: 72, flex: "none" }}
          value={value}
          aria-label={`${label || "Duration"} custom value`}
          onChange={(e) => onChange && onChange(e.target.value)}
        />
      </div>
      {hint && <span className="hm-hint">{hint}</span>}
    </div>
  );
}
