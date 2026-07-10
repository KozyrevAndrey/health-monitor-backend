import React from "react";
import { Button } from "../forms/Button.jsx";

/**
 * Styled confirm dialog (replaces window.confirm). Focus lands on Cancel;
 * Esc closes. Render only when open.
 */
export function ConfirmDialog({ title, body, confirmLabel = "Delete", cancelLabel = "Cancel", danger = true, onConfirm, onCancel }) {
  const cancelRef = React.useRef(null);
  React.useEffect(() => {
    if (cancelRef.current) cancelRef.current.focus();
    const onKey = (e) => { if (e.key === "Escape" && onCancel) onCancel(); };
    document.addEventListener("keydown", onKey);
    return () => document.removeEventListener("keydown", onKey);
  }, []);
  return (
    <>
      <div className="hm-backdrop" onClick={onCancel}></div>
      <div className="hm-dialog" role="dialog" aria-modal="true" aria-label={title}>
        <b style={{ display: "block", fontSize: "var(--text-xl)", fontWeight: "var(--weight-semibold)" }}>{title}</b>
        <p style={{ margin: "var(--space-2) 0 var(--space-5)", color: "var(--text-2)", fontSize: "var(--text-md)" }}>{body}</p>
        <div style={{ display: "flex", justifyContent: "flex-end", gap: "var(--space-2)" }}>
          <button type="button" ref={cancelRef} className="hm-btn" onClick={onCancel}>{cancelLabel}</button>
          <Button variant={danger ? "danger" : "primary"} onClick={onConfirm}>{confirmLabel}</Button>
        </div>
      </div>
    </>
  );
}
