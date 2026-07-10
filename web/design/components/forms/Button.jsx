import React from "react";

/**
 * Button. Variants: default (subtle), primary, danger (outline), ghost.
 * `size="lg"` meets the 44px touch target for mobile/forms.
 */
export function Button({ variant = "default", size = "md", icon, children, style, ...rest }) {
  const cls = [
    "hm-btn",
    variant !== "default" ? `hm-btn--${variant}` : "",
    size === "lg" ? "hm-btn--lg" : "",
  ].filter(Boolean).join(" ");
  return (
    <button type="button" className={cls} style={style} {...rest}>
      {icon}
      {children}
    </button>
  );
}
