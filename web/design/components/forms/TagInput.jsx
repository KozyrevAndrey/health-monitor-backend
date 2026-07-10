import React from "react";
import { Icon } from "../core/Icon.jsx";

/** Tag chips input: Enter/comma adds, Backspace removes last, × removes one. */
export function TagInput({ value = [], onChange, placeholder = "Add tag…", style }) {
  const [draft, setDraft] = React.useState("");
  const ref = React.useRef(null);
  const commit = () => {
    const t = draft.trim().replace(/,$/, "");
    if (t && !value.includes(t)) onChange && onChange([...value, t]);
    setDraft("");
  };
  return (
    <div className="hm-taginput" style={style} onClick={() => ref.current && ref.current.focus()}>
      {value.map((t) => (
        <span key={t} className="hm-tag">
          {t}
          <button type="button" aria-label={`Remove tag ${t}`} onClick={() => onChange && onChange(value.filter((x) => x !== t))}>
            <Icon name="x" size={11} />
          </button>
        </span>
      ))}
      <input
        ref={ref}
        value={draft}
        placeholder={value.length ? "" : placeholder}
        onChange={(e) => setDraft(e.target.value)}
        onKeyDown={(e) => {
          if (e.key === "Enter" || e.key === ",") { e.preventDefault(); commit(); }
          else if (e.key === "Backspace" && !draft && value.length) onChange && onChange(value.slice(0, -1));
        }}
        onBlur={commit}
      />
    </div>
  );
}
