import React from "react";

/**
 * Segmented control — doubles as tabs (period switcher, view toggle, filters).
 * options: [{value, label, icon?, count?}] or strings.
 */
export function Segmented({ options = [], value, onChange, ariaLabel, style }) {
  return (
    <div className="hm-seg" role="tablist" aria-label={ariaLabel} style={style}>
      {options.map((o) => {
        const v = typeof o === "string" ? { value: o, label: o } : o;
        const sel = v.value === value;
        return (
          <button
            key={v.value}
            type="button"
            role="tab"
            className="hm-seg-btn"
            aria-selected={sel}
            onClick={() => onChange && onChange(v.value)}
          >
            {v.icon}
            {v.label}
            {v.count != null && (
              <span className="num" style={{ fontSize: "var(--text-xs)", color: sel ? "var(--text-2)" : "var(--text-3)" }}>{v.count}</span>
            )}
          </button>
        );
      })}
    </div>
  );
}
