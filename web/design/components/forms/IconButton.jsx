import React from "react";
import { Icon } from "../core/Icon.jsx";

/** 32px square icon button. Requires `label` (becomes aria-label + tooltip). */
export function IconButton({ name, label, danger = false, style, ...rest }) {
  return (
    <button
      type="button"
      className={`hm-iconbtn${danger ? " hm-iconbtn--danger" : ""}`}
      aria-label={label}
      data-tip={label}
      style={style}
      {...rest}
    >
      <Icon name={name} size={16} />
    </button>
  );
}
