import React from "react";

let selSeq = 0;

/** Labeled native select. options: [{value, label}] or plain strings. */
export function SelectField({ label, hint, options = [], id, style, ...rest }) {
  const [uid] = React.useState(() => id || `hm-s${++selSeq}`);
  return (
    <div className="hm-field" style={style}>
      {label && <label className="hm-label" htmlFor={uid}>{label}</label>}
      <select id={uid} className="hm-select" {...rest}>
        {options.map((o) => {
          const v = typeof o === "string" ? { value: o, label: o } : o;
          return <option key={v.value} value={v.value}>{v.label}</option>;
        })}
      </select>
      {hint && <span className="hm-hint">{hint}</span>}
    </div>
  );
}
