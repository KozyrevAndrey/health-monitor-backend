import React from "react";
import { Icon } from "../core/Icon.jsx";

/**
 * Right-hand slide-over panel (target detail, add/edit forms).
 * Esc closes; backdrop click closes; body scrolls; optional footer actions.
 */
export function SlideOver({ title, titleExtra, width = 560, onClose, footer, children }) {
  React.useEffect(() => {
    const onKey = (e) => { if (e.key === "Escape" && onClose) onClose(); };
    document.addEventListener("keydown", onKey);
    return () => document.removeEventListener("keydown", onKey);
  }, []);
  return (
    <>
      <div className="hm-backdrop" onClick={onClose}></div>
      <div className="hm-slideover" role="dialog" aria-modal="true" aria-label={typeof title === "string" ? title : "Panel"} style={{ width: `min(${width}px, 100vw)` }}>
        <div className="hm-slideover-head">
          <b style={{ fontSize: "var(--text-xl)", fontWeight: "var(--weight-semibold)", flex: 1, minWidth: 0, display: "flex", alignItems: "center", gap: "var(--space-2)" }}>{title}</b>
          {titleExtra}
          <button type="button" className="hm-iconbtn" aria-label="Close panel" onClick={onClose}><Icon name="x" size={16} /></button>
        </div>
        <div className="hm-slideover-body scroll-thin">{children}</div>
        {footer && <div className="hm-slideover-foot">{footer}</div>}
      </div>
    </>
  );
}
