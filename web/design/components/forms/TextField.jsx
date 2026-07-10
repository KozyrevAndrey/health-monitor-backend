import React from "react";
import { Icon } from "../core/Icon.jsx";

let fieldSeq = 0;

/**
 * Labeled text input with hint / inline error. `mono` for URLs & host:port.
 */
export function TextField({ label, hint, error, mono = false, id, style, ...rest }) {
  const [uid] = React.useState(() => id || `hm-f${++fieldSeq}`);
  return (
    <div className="hm-field" style={style}>
      {label && <label className="hm-label" htmlFor={uid}>{label}</label>}
      <input
        id={uid}
        className={`hm-input${error ? " hm-input--invalid" : ""}${mono ? " hm-input--mono" : ""}`}
        aria-invalid={error ? true : undefined}
        {...rest}
      />
      {error ? (
        <span className="hm-error-text"><Icon name="alert-triangle" size={13} />{error}</span>
      ) : hint ? (
        <span className="hm-hint">{hint}</span>
      ) : null}
    </div>
  );
}
