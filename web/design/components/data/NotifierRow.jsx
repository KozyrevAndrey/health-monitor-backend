import React from "react";
import { Icon } from "../core/Icon.jsx";
import { Switch } from "../forms/Switch.jsx";

const TYPE_META = {
  telegram: { icon: "send", label: "Telegram" },
  email: { icon: "mail", label: "Email (SMTP)" },
  gmail: { icon: "mail", label: "Gmail" },
  gmail_oauth: { icon: "mail", label: "Gmail (OAuth)" },
  webhook: { icon: "link", label: "Webhook" },
};

/** Notifier settings row: type icon + name + masked detail + enable switch. */
export function NotifierRow({ notifier, onToggle, onEdit, onDelete, style }) {
  const { name, type = "webhook", enabled = true, detail } = notifier;
  const meta = TYPE_META[type] || TYPE_META.webhook;
  return (
    <div style={{ display: "flex", alignItems: "center", gap: "var(--space-3)", padding: "var(--space-3) var(--space-4)", ...style }}>
      <span style={{ display: "flex", alignItems: "center", justifyContent: "center", width: 32, height: 32, flex: "none", borderRadius: "var(--radius-md)", background: "var(--bg-2)", border: "1px solid var(--border-1)", color: enabled ? "var(--text-2)" : "var(--text-3)" }}>
        <Icon name={meta.icon} size={15} />
      </span>
      <div style={{ flex: 1, minWidth: 0, opacity: enabled ? 1 : 0.55 }}>
        <b style={{ display: "block", fontSize: "var(--text-md)", fontWeight: "var(--weight-semibold)" }}>{name}</b>
        <span style={{ fontSize: "var(--text-sm)", color: "var(--text-3)" }}>
          {meta.label}{detail && <> · <code style={{ fontSize: "var(--text-xs)" }}>{detail}</code></>}
        </span>
      </div>
      <Switch checked={enabled} onChange={onToggle} label="" />
      {onEdit && <button type="button" className="hm-iconbtn" aria-label={`Edit ${name}`} onClick={onEdit}><Icon name="pencil" size={14} /></button>}
      {onDelete && <button type="button" className="hm-iconbtn hm-iconbtn--danger" aria-label={`Delete ${name}`} onClick={onDelete}><Icon name="trash" size={14} /></button>}
    </div>
  );
}
