/* @ds-bundle: {"format":4,"namespace":"HealthMonitorRedesign_502b6f","components":[{"name":"Icon","sourcePath":"components/core/Icon.jsx"},{"name":"IncidentItem","sourcePath":"components/data/IncidentItem.jsx"},{"name":"NotifierRow","sourcePath":"components/data/NotifierRow.jsx"},{"name":"StatCard","sourcePath":"components/data/StatCard.jsx"},{"name":"TargetCard","sourcePath":"components/data/TargetCard.jsx"},{"name":"TargetTable","sourcePath":"components/data/TargetTable.jsx"},{"name":"ConfirmDialog","sourcePath":"components/feedback/ConfirmDialog.jsx"},{"name":"EmptyState","sourcePath":"components/feedback/EmptyState.jsx"},{"name":"Skeleton","sourcePath":"components/feedback/Skeleton.jsx"},{"name":"SlideOver","sourcePath":"components/feedback/SlideOver.jsx"},{"name":"Toast","sourcePath":"components/feedback/Toast.jsx"},{"name":"Button","sourcePath":"components/forms/Button.jsx"},{"name":"DurationField","sourcePath":"components/forms/DurationField.jsx"},{"name":"IconButton","sourcePath":"components/forms/IconButton.jsx"},{"name":"Segmented","sourcePath":"components/forms/Segmented.jsx"},{"name":"SelectField","sourcePath":"components/forms/SelectField.jsx"},{"name":"Switch","sourcePath":"components/forms/Switch.jsx"},{"name":"TagInput","sourcePath":"components/forms/TagInput.jsx"},{"name":"TextField","sourcePath":"components/forms/TextField.jsx"},{"name":"LiveIndicator","sourcePath":"components/status/LiveIndicator.jsx"},{"name":"Sparkline","sourcePath":"components/status/Sparkline.jsx"},{"name":"StatusBadge","sourcePath":"components/status/StatusBadge.jsx"},{"name":"StatusDot","sourcePath":"components/status/StatusDot.jsx"},{"name":"TickBar","sourcePath":"components/status/TickBar.jsx"}],"sourceHashes":{"components/core/Icon.jsx":"086e9a4695ab","components/data/IncidentItem.jsx":"8b314b62f6be","components/data/NotifierRow.jsx":"df25f0bb0cf2","components/data/StatCard.jsx":"47b6c69c3e68","components/data/TargetCard.jsx":"8fc2b45c6c3f","components/data/TargetTable.jsx":"e23dba8ee46e","components/feedback/ConfirmDialog.jsx":"cbadb8a1cca2","components/feedback/EmptyState.jsx":"dd4b8c7bfe52","components/feedback/Skeleton.jsx":"41ba7b88d374","components/feedback/SlideOver.jsx":"ce0747eee175","components/feedback/Toast.jsx":"794040e17fe2","components/forms/Button.jsx":"b9185a287310","components/forms/DurationField.jsx":"09b54557019e","components/forms/IconButton.jsx":"b80254a557b7","components/forms/Segmented.jsx":"5b6de23fab24","components/forms/SelectField.jsx":"25263d4c678c","components/forms/Switch.jsx":"b82c1b0feea7","components/forms/TagInput.jsx":"3083b6982724","components/forms/TextField.jsx":"aa41e2d686ef","components/status/LiveIndicator.jsx":"d15a56941e21","components/status/Sparkline.jsx":"b71025894315","components/status/StatusBadge.jsx":"c68bed00f4dd","components/status/StatusDot.jsx":"521ddf1b233f","components/status/TickBar.jsx":"ed1c71f8b9f8","ui_kits/dashboard/app.jsx":"6b74903fe928","ui_kits/dashboard/data.js":"282c5e4f8cfb","ui_kits/dashboard/detail.jsx":"c95626038932","ui_kits/dashboard/forms.jsx":"d280b443d29e"},"inlinedExternals":[],"unexposedExports":[]} */

(() => {

const __ds_ns = (window.HealthMonitorRedesign_502b6f = window.HealthMonitorRedesign_502b6f || {});

const __ds_scope = {};

(__ds_ns.__errors = __ds_ns.__errors || []);

// components/core/Icon.jsx
try { (() => {
function _extends() { return _extends = Object.assign ? Object.assign.bind() : function (n) { for (var e = 1; e < arguments.length; e++) { var t = arguments[e]; for (var r in t) ({}).hasOwnProperty.call(t, r) && (n[r] = t[r]); } return n; }, _extends.apply(null, arguments); }
const P = {
  check: ["M20 6 9 17l-5-5"],
  x: ["M18 6 6 18", "M6 6l12 12"],
  plus: ["M12 5v14", "M5 12h14"],
  search: ["M21 21l-4.35-4.35", "M11 4a7 7 0 1 1 0 14 7 7 0 0 1 0-14z"],
  "chevron-down": ["m6 9 6 6 6-6"],
  "chevron-right": ["m9 6 6 6-6 6"],
  "chevron-left": ["m15 6-6 6 6 6"],
  "arrow-left": ["M19 12H5", "m12 19-7-7 7-7"],
  clock: ["M12 3a9 9 0 1 1 0 18 9 9 0 0 1 0-18z", "M12 7v5l3 2"],
  trash: ["M4 7h16", "M9 7V4h6v3", "M6 7l1 13h10l1-13"],
  pencil: ["M17 3l4 4L8 20l-5 1 1-5L17 3"],
  "alert-triangle": ["M12 3 2 20h20L12 3z", "M12 10v4", "M12 17h.01"],
  "check-circle": ["M12 3a9 9 0 1 1 0 18 9 9 0 0 1 0-18z", "m8.5 12.5 2.5 2.5 4.5-5"],
  "x-circle": ["M12 3a9 9 0 1 1 0 18 9 9 0 0 1 0-18z", "M15 9l-6 6", "M9 9l6 6"],
  info: ["M12 3a9 9 0 1 1 0 18 9 9 0 0 1 0-18z", "M12 16v-5", "M12 8h.01"],
  activity: ["M22 12h-4l-3 9L9 3l-3 9H2"],
  globe: ["M12 3a9 9 0 1 1 0 18 9 9 0 0 1 0-18z", "M3 12h18", "M12 3c2.5 2.5 3.8 5.6 3.8 9S14.5 18.5 12 21c-2.5-2.5-3.8-5.6-3.8-9S9.5 5.5 12 3z"],
  server: ["M4 4h16a2 2 0 0 1 2 2v3a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2z", "M4 13h16a2 2 0 0 1 2 2v3a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2v-3a2 2 0 0 1 2-2z", "M6 7.5h.01", "M6 16.5h.01"],
  "at-sign": ["M12 8a4 4 0 1 0 0 8 4 4 0 0 0 0-8z", "M16 8v5a3 3 0 0 0 6 0v-1a10 10 0 1 0-3.9 7.9"],
  bell: ["M6 9a6 6 0 0 1 12 0c0 5 2 6 2 6H4s2-1 2-6", "M10 19a2 2 0 0 0 4 0"],
  sliders: ["M4 8h10", "M18 8h2", "M4 16h4", "M12 16h8", "M16 8h.01M14 6a2 2 0 1 0 4 0 2 2 0 0 0-4 0z", "M8 16a2 2 0 1 0 4 0 2 2 0 0 0-4 0z"],
  zap: ["M13 2 3 14h7l-1 8 11-12h-7l1-8z"],
  pause: ["M9 5v14", "M15 5v14"],
  send: ["M22 2 11 13", "M22 2 15 22l-4-9-9-4z"],
  mail: ["M4 5h16a1 1 0 0 1 1 1v12a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1V6a1 1 0 0 1 1-1z", "m3 7 9 6 9-6"],
  link: ["M10 14a5 5 0 0 1 0-7l1.5-1.5a5 5 0 0 1 7 7L17.5 13.5", "M14 10a5 5 0 0 1 0 7l-1.5 1.5a5 5 0 0 1-7-7L6.5 10.5"],
  moon: ["M21 13A9 9 0 1 1 11 3a7 7 0 0 0 10 10z"],
  sun: ["M12 8a4 4 0 1 1 0 8 4 4 0 0 1 0-8z", "M12 2v2", "M12 20v2", "M2 12h2", "M20 12h2", "m5 5 1.4 1.4", "m17.6 17.6 1.4 1.4", "m19 5-1.4 1.4", "m6.4 17.6-1.4 1.4"],
  grid: ["M4 4h6v6H4z", "M14 4h6v6h-6z", "M4 14h6v6H4z", "M14 14h6v6h-6z"],
  list: ["M8 6h13", "M8 12h13", "M8 18h13", "M3 6h.01", "M3 12h.01", "M3 18h.01"],
  "more-h": ["M5 12h.01", "M12 12h.01", "M19 12h.01"],
  external: ["M14 3h7v7", "M21 3l-9 9", "M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h6"],
  refresh: ["M21 12a9 9 0 1 1-2.9-6.6L21 8", "M21 3v5h-5"],
  wifi: ["M5 12a11 11 0 0 1 14 0", "M8.5 15.5a6 6 0 0 1 7 0", "M12 19h.01"],
  filter: ["M3 5h18l-7 8v6l-4-2v-4L3 5z"]
};

/** Inline stroke icon (Lucide-style: 24 grid, 2px stroke, round caps). */
function Icon({
  name,
  size = 16,
  strokeWidth = 2,
  style,
  ...rest
}) {
  const paths = P[name] || P["info"];
  return /*#__PURE__*/React.createElement("svg", _extends({
    width: size,
    height: size,
    viewBox: "0 0 24 24",
    fill: "none",
    stroke: "currentColor",
    strokeWidth: strokeWidth,
    strokeLinecap: "round",
    strokeLinejoin: "round",
    "aria-hidden": rest["aria-label"] ? undefined : true,
    style: {
      flex: "none",
      ...style
    }
  }, rest), paths.map((d, i) => /*#__PURE__*/React.createElement("path", {
    key: i,
    d: d
  })));
}
Icon.names = Object.keys(P);
Object.assign(__ds_scope, { Icon });
})(); } catch (e) { __ds_ns.__errors.push({ path: "components/core/Icon.jsx", error: String((e && e.message) || e) }); }

// components/data/IncidentItem.jsx
try { (() => {
const SEV = {
  critical: {
    color: "var(--status-down)",
    icon: "alert-triangle"
  },
  warning: {
    color: "var(--status-degraded)",
    icon: "zap"
  },
  info: {
    color: "var(--accent)",
    icon: "info"
  }
};

/**
 * Incident row for timeline lists. Ongoing incidents show a pulsing marker
 * and red tint; resolved ones a check + muted rail.
 */
function IncidentItem({
  incident,
  last = false,
  onOpen,
  style
}) {
  const {
    target,
    status = "resolved",
    severity = "critical",
    startedAt,
    duration,
    failureCount,
    lastError
  } = incident;
  const ongoing = status === "ongoing";
  const sev = SEV[severity] || SEV.critical;
  return /*#__PURE__*/React.createElement("div", {
    className: "hm-incident",
    style: {
      background: ongoing ? "var(--status-down-subtle)" : "transparent",
      cursor: onOpen ? "pointer" : "default",
      ...style
    },
    onClick: onOpen
  }, /*#__PURE__*/React.createElement("div", {
    className: "hm-incident-rail",
    "aria-hidden": "true"
  }, /*#__PURE__*/React.createElement("span", {
    className: `hm-dot${ongoing ? " hm-dot--pulse" : ""}`,
    "data-status": ongoing ? "down" : "up"
  }), !last && /*#__PURE__*/React.createElement("span", {
    className: "line"
  })), /*#__PURE__*/React.createElement("div", {
    style: {
      flex: 1,
      minWidth: 0
    }
  }, /*#__PURE__*/React.createElement("div", {
    style: {
      display: "flex",
      alignItems: "center",
      gap: "var(--space-2)",
      flexWrap: "wrap"
    }
  }, /*#__PURE__*/React.createElement("b", {
    style: {
      fontSize: "var(--text-md)",
      fontWeight: "var(--weight-semibold)"
    }
  }, target), /*#__PURE__*/React.createElement("span", {
    style: {
      display: "inline-flex",
      alignItems: "center",
      gap: 4,
      fontSize: "var(--text-xs)",
      fontWeight: "var(--weight-semibold)",
      color: sev.color,
      textTransform: "uppercase",
      letterSpacing: "var(--tracking-caps)"
    }
  }, /*#__PURE__*/React.createElement(__ds_scope.Icon, {
    name: sev.icon,
    size: 11
  }), severity), /*#__PURE__*/React.createElement("span", {
    style: {
      flex: 1
    }
  }), /*#__PURE__*/React.createElement("span", {
    className: "num",
    style: {
      fontSize: "var(--text-sm)",
      color: ongoing ? "var(--status-down)" : "var(--text-3)",
      fontWeight: ongoing ? "var(--weight-semibold)" : "var(--weight-regular)"
    }
  }, ongoing ? `ongoing · ${duration}` : duration)), /*#__PURE__*/React.createElement("div", {
    style: {
      fontSize: "var(--text-sm)",
      color: "var(--text-2)",
      marginTop: 2
    }
  }, startedAt, failureCount != null && /*#__PURE__*/React.createElement(React.Fragment, null, " \xB7 ", /*#__PURE__*/React.createElement("span", {
    className: "num"
  }, failureCount), " failed checks")), lastError && /*#__PURE__*/React.createElement("code", {
    style: {
      display: "block",
      marginTop: 4,
      fontSize: "var(--text-xs)",
      color: "var(--text-3)",
      overflow: "hidden",
      textOverflow: "ellipsis",
      whiteSpace: "nowrap"
    }
  }, lastError)));
}
Object.assign(__ds_scope, { IncidentItem });
})(); } catch (e) { __ds_ns.__errors.push({ path: "components/data/IncidentItem.jsx", error: String((e && e.message) || e) }); }

// components/data/StatCard.jsx
try { (() => {
/**
 * Compact metric card. `tone` tints the value (e.g. danger for failing count).
 */
function StatCard({
  label,
  value,
  unit,
  sub,
  tone,
  icon,
  style
}) {
  const color = tone === "danger" ? "var(--status-down)" : tone === "success" ? "var(--status-up)" : tone === "warning" ? "var(--status-degraded)" : "var(--text-1)";
  return /*#__PURE__*/React.createElement("div", {
    className: "hm-card hm-statcard",
    style: style
  }, /*#__PURE__*/React.createElement("span", {
    className: "overline"
  }, icon, label), /*#__PURE__*/React.createElement("b", {
    style: {
      color
    }
  }, value, unit && /*#__PURE__*/React.createElement("span", {
    style: {
      fontSize: "var(--text-lg)",
      fontWeight: "var(--weight-medium)",
      color: "var(--text-3)",
      marginLeft: 3
    }
  }, unit)), sub && /*#__PURE__*/React.createElement("small", null, sub));
}
Object.assign(__ds_scope, { StatCard });
})(); } catch (e) { __ds_ns.__errors.push({ path: "components/data/StatCard.jsx", error: String((e && e.message) || e) }); }

// components/feedback/EmptyState.jsx
try { (() => {
/** Designed empty state: icon, one-line title, short body, primary action. */
function EmptyState({
  icon = "activity",
  title,
  body,
  action,
  style
}) {
  return /*#__PURE__*/React.createElement("div", {
    className: "hm-empty",
    style: style
  }, /*#__PURE__*/React.createElement(__ds_scope.Icon, {
    name: icon,
    size: 28
  }), /*#__PURE__*/React.createElement("b", null, title), body && /*#__PURE__*/React.createElement("p", null, body), action);
}
Object.assign(__ds_scope, { EmptyState });
})(); } catch (e) { __ds_ns.__errors.push({ path: "components/feedback/EmptyState.jsx", error: String((e && e.message) || e) }); }

// components/feedback/Skeleton.jsx
try { (() => {
/** Shimmering skeleton block. Compose several to sketch a loading layout. */
function Skeleton({
  width = "100%",
  height = 14,
  radius,
  style
}) {
  return /*#__PURE__*/React.createElement("span", {
    className: "hm-skel",
    "aria-hidden": "true",
    style: {
      display: "block",
      width,
      height,
      borderRadius: radius,
      ...style
    }
  });
}
Object.assign(__ds_scope, { Skeleton });
})(); } catch (e) { __ds_ns.__errors.push({ path: "components/feedback/Skeleton.jsx", error: String((e && e.message) || e) }); }

// components/feedback/SlideOver.jsx
try { (() => {
/**
 * Right-hand slide-over panel (target detail, add/edit forms).
 * Esc closes; backdrop click closes; body scrolls; optional footer actions.
 */
function SlideOver({
  title,
  titleExtra,
  width = 560,
  onClose,
  footer,
  children
}) {
  React.useEffect(() => {
    const onKey = e => {
      if (e.key === "Escape" && onClose) onClose();
    };
    document.addEventListener("keydown", onKey);
    return () => document.removeEventListener("keydown", onKey);
  }, []);
  return /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("div", {
    className: "hm-backdrop",
    onClick: onClose
  }), /*#__PURE__*/React.createElement("div", {
    className: "hm-slideover",
    role: "dialog",
    "aria-modal": "true",
    "aria-label": typeof title === "string" ? title : "Panel",
    style: {
      width: `min(${width}px, 100vw)`
    }
  }, /*#__PURE__*/React.createElement("div", {
    className: "hm-slideover-head"
  }, /*#__PURE__*/React.createElement("b", {
    style: {
      fontSize: "var(--text-xl)",
      fontWeight: "var(--weight-semibold)",
      flex: 1,
      minWidth: 0,
      display: "flex",
      alignItems: "center",
      gap: "var(--space-2)"
    }
  }, title), titleExtra, /*#__PURE__*/React.createElement("button", {
    type: "button",
    className: "hm-iconbtn",
    "aria-label": "Close panel",
    onClick: onClose
  }, /*#__PURE__*/React.createElement(__ds_scope.Icon, {
    name: "x",
    size: 16
  }))), /*#__PURE__*/React.createElement("div", {
    className: "hm-slideover-body scroll-thin"
  }, children), footer && /*#__PURE__*/React.createElement("div", {
    className: "hm-slideover-foot"
  }, footer)));
}
Object.assign(__ds_scope, { SlideOver });
})(); } catch (e) { __ds_ns.__errors.push({ path: "components/feedback/SlideOver.jsx", error: String((e && e.message) || e) }); }

// components/feedback/Toast.jsx
try { (() => {
const KIND_ICON = {
  success: "check-circle",
  danger: "alert-triangle",
  warning: "zap",
  info: "info"
};

/** Single toast. Compose inside a `.hm-toaststack` fixed container. */
function Toast({
  kind = "info",
  title,
  message,
  onDismiss,
  style
}) {
  return /*#__PURE__*/React.createElement("div", {
    className: "hm-toast",
    "data-kind": kind,
    role: "status",
    style: style
  }, /*#__PURE__*/React.createElement(__ds_scope.Icon, {
    name: KIND_ICON[kind] || "info",
    size: 16
  }), /*#__PURE__*/React.createElement("div", {
    style: {
      flex: 1,
      minWidth: 0
    }
  }, /*#__PURE__*/React.createElement("b", null, title), message && /*#__PURE__*/React.createElement("p", null, message)), onDismiss && /*#__PURE__*/React.createElement("button", {
    type: "button",
    className: "hm-iconbtn",
    style: {
      width: 24,
      height: 24,
      margin: "-4px -6px 0 0"
    },
    "aria-label": "Dismiss",
    onClick: onDismiss
  }, /*#__PURE__*/React.createElement(__ds_scope.Icon, {
    name: "x",
    size: 13
  })));
}
Object.assign(__ds_scope, { Toast });
})(); } catch (e) { __ds_ns.__errors.push({ path: "components/feedback/Toast.jsx", error: String((e && e.message) || e) }); }

// components/forms/Button.jsx
try { (() => {
function _extends() { return _extends = Object.assign ? Object.assign.bind() : function (n) { for (var e = 1; e < arguments.length; e++) { var t = arguments[e]; for (var r in t) ({}).hasOwnProperty.call(t, r) && (n[r] = t[r]); } return n; }, _extends.apply(null, arguments); }
/**
 * Button. Variants: default (subtle), primary, danger (outline), ghost.
 * `size="lg"` meets the 44px touch target for mobile/forms.
 */
function Button({
  variant = "default",
  size = "md",
  icon,
  children,
  style,
  ...rest
}) {
  const cls = ["hm-btn", variant !== "default" ? `hm-btn--${variant}` : "", size === "lg" ? "hm-btn--lg" : ""].filter(Boolean).join(" ");
  return /*#__PURE__*/React.createElement("button", _extends({
    type: "button",
    className: cls,
    style: style
  }, rest), icon, children);
}
Object.assign(__ds_scope, { Button });
})(); } catch (e) { __ds_ns.__errors.push({ path: "components/forms/Button.jsx", error: String((e && e.message) || e) }); }

// components/feedback/ConfirmDialog.jsx
try { (() => {
/**
 * Styled confirm dialog (replaces window.confirm). Focus lands on Cancel;
 * Esc closes. Render only when open.
 */
function ConfirmDialog({
  title,
  body,
  confirmLabel = "Delete",
  cancelLabel = "Cancel",
  danger = true,
  onConfirm,
  onCancel
}) {
  const cancelRef = React.useRef(null);
  React.useEffect(() => {
    if (cancelRef.current) cancelRef.current.focus();
    const onKey = e => {
      if (e.key === "Escape" && onCancel) onCancel();
    };
    document.addEventListener("keydown", onKey);
    return () => document.removeEventListener("keydown", onKey);
  }, []);
  return /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("div", {
    className: "hm-backdrop",
    onClick: onCancel
  }), /*#__PURE__*/React.createElement("div", {
    className: "hm-dialog",
    role: "dialog",
    "aria-modal": "true",
    "aria-label": title
  }, /*#__PURE__*/React.createElement("b", {
    style: {
      display: "block",
      fontSize: "var(--text-xl)",
      fontWeight: "var(--weight-semibold)"
    }
  }, title), /*#__PURE__*/React.createElement("p", {
    style: {
      margin: "var(--space-2) 0 var(--space-5)",
      color: "var(--text-2)",
      fontSize: "var(--text-md)"
    }
  }, body), /*#__PURE__*/React.createElement("div", {
    style: {
      display: "flex",
      justifyContent: "flex-end",
      gap: "var(--space-2)"
    }
  }, /*#__PURE__*/React.createElement("button", {
    type: "button",
    ref: cancelRef,
    className: "hm-btn",
    onClick: onCancel
  }, cancelLabel), /*#__PURE__*/React.createElement(__ds_scope.Button, {
    variant: danger ? "danger" : "primary",
    onClick: onConfirm
  }, confirmLabel))));
}
Object.assign(__ds_scope, { ConfirmDialog });
})(); } catch (e) { __ds_ns.__errors.push({ path: "components/feedback/ConfirmDialog.jsx", error: String((e && e.message) || e) }); }

// components/forms/DurationField.jsx
try { (() => {
const PRESETS = ["30s", "1m", "5m", "15m", "1h"];

/**
 * Duration input: preset chips + free-form field ("30s", "1m", "2h").
 * Value is a Go-style duration string.
 */
function DurationField({
  label,
  value = "1m",
  onChange,
  presets = PRESETS,
  hint,
  style
}) {
  return /*#__PURE__*/React.createElement("div", {
    className: "hm-field",
    style: style
  }, label && /*#__PURE__*/React.createElement("span", {
    className: "hm-label"
  }, label), /*#__PURE__*/React.createElement("div", {
    style: {
      display: "flex",
      gap: "var(--space-2)",
      alignItems: "center"
    }
  }, /*#__PURE__*/React.createElement("div", {
    className: "hm-seg",
    role: "group",
    "aria-label": `${label || "Duration"} presets`
  }, presets.map(p => /*#__PURE__*/React.createElement("button", {
    key: p,
    type: "button",
    className: "hm-seg-btn num",
    "aria-pressed": value === p,
    onClick: () => onChange && onChange(p)
  }, p))), /*#__PURE__*/React.createElement("input", {
    className: "hm-input hm-input--mono",
    style: {
      width: 72,
      flex: "none"
    },
    value: value,
    "aria-label": `${label || "Duration"} custom value`,
    onChange: e => onChange && onChange(e.target.value)
  })), hint && /*#__PURE__*/React.createElement("span", {
    className: "hm-hint"
  }, hint));
}
Object.assign(__ds_scope, { DurationField });
})(); } catch (e) { __ds_ns.__errors.push({ path: "components/forms/DurationField.jsx", error: String((e && e.message) || e) }); }

// components/forms/IconButton.jsx
try { (() => {
function _extends() { return _extends = Object.assign ? Object.assign.bind() : function (n) { for (var e = 1; e < arguments.length; e++) { var t = arguments[e]; for (var r in t) ({}).hasOwnProperty.call(t, r) && (n[r] = t[r]); } return n; }, _extends.apply(null, arguments); }
/** 32px square icon button. Requires `label` (becomes aria-label + tooltip). */
function IconButton({
  name,
  label,
  danger = false,
  style,
  ...rest
}) {
  return /*#__PURE__*/React.createElement("button", _extends({
    type: "button",
    className: `hm-iconbtn${danger ? " hm-iconbtn--danger" : ""}`,
    "aria-label": label,
    "data-tip": label,
    style: style
  }, rest), /*#__PURE__*/React.createElement(__ds_scope.Icon, {
    name: name,
    size: 16
  }));
}
Object.assign(__ds_scope, { IconButton });
})(); } catch (e) { __ds_ns.__errors.push({ path: "components/forms/IconButton.jsx", error: String((e && e.message) || e) }); }

// components/forms/Segmented.jsx
try { (() => {
/**
 * Segmented control — doubles as tabs (period switcher, view toggle, filters).
 * options: [{value, label, icon?, count?}] or strings.
 */
function Segmented({
  options = [],
  value,
  onChange,
  ariaLabel,
  style
}) {
  return /*#__PURE__*/React.createElement("div", {
    className: "hm-seg",
    role: "tablist",
    "aria-label": ariaLabel,
    style: style
  }, options.map(o => {
    const v = typeof o === "string" ? {
      value: o,
      label: o
    } : o;
    const sel = v.value === value;
    return /*#__PURE__*/React.createElement("button", {
      key: v.value,
      type: "button",
      role: "tab",
      className: "hm-seg-btn",
      "aria-selected": sel,
      onClick: () => onChange && onChange(v.value)
    }, v.icon, v.label, v.count != null && /*#__PURE__*/React.createElement("span", {
      className: "num",
      style: {
        fontSize: "var(--text-xs)",
        color: sel ? "var(--text-2)" : "var(--text-3)"
      }
    }, v.count));
  }));
}
Object.assign(__ds_scope, { Segmented });
})(); } catch (e) { __ds_ns.__errors.push({ path: "components/forms/Segmented.jsx", error: String((e && e.message) || e) }); }

// components/forms/SelectField.jsx
try { (() => {
function _extends() { return _extends = Object.assign ? Object.assign.bind() : function (n) { for (var e = 1; e < arguments.length; e++) { var t = arguments[e]; for (var r in t) ({}).hasOwnProperty.call(t, r) && (n[r] = t[r]); } return n; }, _extends.apply(null, arguments); }
let selSeq = 0;

/** Labeled native select. options: [{value, label}] or plain strings. */
function SelectField({
  label,
  hint,
  options = [],
  id,
  style,
  ...rest
}) {
  const [uid] = React.useState(() => id || `hm-s${++selSeq}`);
  return /*#__PURE__*/React.createElement("div", {
    className: "hm-field",
    style: style
  }, label && /*#__PURE__*/React.createElement("label", {
    className: "hm-label",
    htmlFor: uid
  }, label), /*#__PURE__*/React.createElement("select", _extends({
    id: uid,
    className: "hm-select"
  }, rest), options.map(o => {
    const v = typeof o === "string" ? {
      value: o,
      label: o
    } : o;
    return /*#__PURE__*/React.createElement("option", {
      key: v.value,
      value: v.value
    }, v.label);
  })), hint && /*#__PURE__*/React.createElement("span", {
    className: "hm-hint"
  }, hint));
}
Object.assign(__ds_scope, { SelectField });
})(); } catch (e) { __ds_ns.__errors.push({ path: "components/forms/SelectField.jsx", error: String((e && e.message) || e) }); }

// components/forms/Switch.jsx
try { (() => {
/** Toggle switch (green = enabled). Renders an optional side label. */
function Switch({
  checked = false,
  onChange,
  label,
  disabled = false,
  style
}) {
  return /*#__PURE__*/React.createElement("label", {
    style: {
      display: "inline-flex",
      alignItems: "center",
      gap: "var(--space-2)",
      cursor: disabled ? "not-allowed" : "pointer",
      opacity: disabled ? 0.5 : 1,
      ...style
    }
  }, /*#__PURE__*/React.createElement("span", {
    className: "hm-switch"
  }, /*#__PURE__*/React.createElement("input", {
    type: "checkbox",
    role: "switch",
    checked: checked,
    disabled: disabled,
    onChange: e => onChange && onChange(e.target.checked)
  }), /*#__PURE__*/React.createElement("i", null)), label && /*#__PURE__*/React.createElement("span", {
    style: {
      fontSize: "var(--text-md)",
      color: "var(--text-1)"
    }
  }, label));
}
Object.assign(__ds_scope, { Switch });
})(); } catch (e) { __ds_ns.__errors.push({ path: "components/forms/Switch.jsx", error: String((e && e.message) || e) }); }

// components/data/NotifierRow.jsx
try { (() => {
const TYPE_META = {
  telegram: {
    icon: "send",
    label: "Telegram"
  },
  email: {
    icon: "mail",
    label: "Email (SMTP)"
  },
  gmail: {
    icon: "mail",
    label: "Gmail"
  },
  gmail_oauth: {
    icon: "mail",
    label: "Gmail (OAuth)"
  },
  webhook: {
    icon: "link",
    label: "Webhook"
  }
};

/** Notifier settings row: type icon + name + masked detail + enable switch. */
function NotifierRow({
  notifier,
  onToggle,
  onEdit,
  onDelete,
  style
}) {
  const {
    name,
    type = "webhook",
    enabled = true,
    detail
  } = notifier;
  const meta = TYPE_META[type] || TYPE_META.webhook;
  return /*#__PURE__*/React.createElement("div", {
    style: {
      display: "flex",
      alignItems: "center",
      gap: "var(--space-3)",
      padding: "var(--space-3) var(--space-4)",
      ...style
    }
  }, /*#__PURE__*/React.createElement("span", {
    style: {
      display: "flex",
      alignItems: "center",
      justifyContent: "center",
      width: 32,
      height: 32,
      flex: "none",
      borderRadius: "var(--radius-md)",
      background: "var(--bg-2)",
      border: "1px solid var(--border-1)",
      color: enabled ? "var(--text-2)" : "var(--text-3)"
    }
  }, /*#__PURE__*/React.createElement(__ds_scope.Icon, {
    name: meta.icon,
    size: 15
  })), /*#__PURE__*/React.createElement("div", {
    style: {
      flex: 1,
      minWidth: 0,
      opacity: enabled ? 1 : 0.55
    }
  }, /*#__PURE__*/React.createElement("b", {
    style: {
      display: "block",
      fontSize: "var(--text-md)",
      fontWeight: "var(--weight-semibold)"
    }
  }, name), /*#__PURE__*/React.createElement("span", {
    style: {
      fontSize: "var(--text-sm)",
      color: "var(--text-3)"
    }
  }, meta.label, detail && /*#__PURE__*/React.createElement(React.Fragment, null, " \xB7 ", /*#__PURE__*/React.createElement("code", {
    style: {
      fontSize: "var(--text-xs)"
    }
  }, detail)))), /*#__PURE__*/React.createElement(__ds_scope.Switch, {
    checked: enabled,
    onChange: onToggle,
    label: ""
  }), onEdit && /*#__PURE__*/React.createElement("button", {
    type: "button",
    className: "hm-iconbtn",
    "aria-label": `Edit ${name}`,
    onClick: onEdit
  }, /*#__PURE__*/React.createElement(__ds_scope.Icon, {
    name: "pencil",
    size: 14
  })), onDelete && /*#__PURE__*/React.createElement("button", {
    type: "button",
    className: "hm-iconbtn hm-iconbtn--danger",
    "aria-label": `Delete ${name}`,
    onClick: onDelete
  }, /*#__PURE__*/React.createElement(__ds_scope.Icon, {
    name: "trash",
    size: 14
  })));
}
Object.assign(__ds_scope, { NotifierRow });
})(); } catch (e) { __ds_ns.__errors.push({ path: "components/data/NotifierRow.jsx", error: String((e && e.message) || e) }); }

// components/forms/TagInput.jsx
try { (() => {
/** Tag chips input: Enter/comma adds, Backspace removes last, × removes one. */
function TagInput({
  value = [],
  onChange,
  placeholder = "Add tag…",
  style
}) {
  const [draft, setDraft] = React.useState("");
  const ref = React.useRef(null);
  const commit = () => {
    const t = draft.trim().replace(/,$/, "");
    if (t && !value.includes(t)) onChange && onChange([...value, t]);
    setDraft("");
  };
  return /*#__PURE__*/React.createElement("div", {
    className: "hm-taginput",
    style: style,
    onClick: () => ref.current && ref.current.focus()
  }, value.map(t => /*#__PURE__*/React.createElement("span", {
    key: t,
    className: "hm-tag"
  }, t, /*#__PURE__*/React.createElement("button", {
    type: "button",
    "aria-label": `Remove tag ${t}`,
    onClick: () => onChange && onChange(value.filter(x => x !== t))
  }, /*#__PURE__*/React.createElement(__ds_scope.Icon, {
    name: "x",
    size: 11
  })))), /*#__PURE__*/React.createElement("input", {
    ref: ref,
    value: draft,
    placeholder: value.length ? "" : placeholder,
    onChange: e => setDraft(e.target.value),
    onKeyDown: e => {
      if (e.key === "Enter" || e.key === ",") {
        e.preventDefault();
        commit();
      } else if (e.key === "Backspace" && !draft && value.length) onChange && onChange(value.slice(0, -1));
    },
    onBlur: commit
  }));
}
Object.assign(__ds_scope, { TagInput });
})(); } catch (e) { __ds_ns.__errors.push({ path: "components/forms/TagInput.jsx", error: String((e && e.message) || e) }); }

// components/forms/TextField.jsx
try { (() => {
function _extends() { return _extends = Object.assign ? Object.assign.bind() : function (n) { for (var e = 1; e < arguments.length; e++) { var t = arguments[e]; for (var r in t) ({}).hasOwnProperty.call(t, r) && (n[r] = t[r]); } return n; }, _extends.apply(null, arguments); }
let fieldSeq = 0;

/**
 * Labeled text input with hint / inline error. `mono` for URLs & host:port.
 */
function TextField({
  label,
  hint,
  error,
  mono = false,
  id,
  style,
  ...rest
}) {
  const [uid] = React.useState(() => id || `hm-f${++fieldSeq}`);
  return /*#__PURE__*/React.createElement("div", {
    className: "hm-field",
    style: style
  }, label && /*#__PURE__*/React.createElement("label", {
    className: "hm-label",
    htmlFor: uid
  }, label), /*#__PURE__*/React.createElement("input", _extends({
    id: uid,
    className: `hm-input${error ? " hm-input--invalid" : ""}${mono ? " hm-input--mono" : ""}`,
    "aria-invalid": error ? true : undefined
  }, rest)), error ? /*#__PURE__*/React.createElement("span", {
    className: "hm-error-text"
  }, /*#__PURE__*/React.createElement(__ds_scope.Icon, {
    name: "alert-triangle",
    size: 13
  }), error) : hint ? /*#__PURE__*/React.createElement("span", {
    className: "hm-hint"
  }, hint) : null);
}
Object.assign(__ds_scope, { TextField });
})(); } catch (e) { __ds_ns.__errors.push({ path: "components/forms/TextField.jsx", error: String((e && e.message) || e) }); }

// components/status/LiveIndicator.jsx
try { (() => {
const LABEL = {
  live: "Live",
  reconnecting: "Reconnecting…",
  polling: "Polling (30s)"
};

/** SSE connection indicator: live / reconnecting / polling fallback. */
function LiveIndicator({
  state = "live",
  label,
  style
}) {
  return /*#__PURE__*/React.createElement("span", {
    className: "hm-live",
    "data-state": state,
    style: style,
    role: "status"
  }, /*#__PURE__*/React.createElement("span", {
    className: "hm-dot"
  }), label || LABEL[state] || state);
}
Object.assign(__ds_scope, { LiveIndicator });
})(); } catch (e) { __ds_ns.__errors.push({ path: "components/status/LiveIndicator.jsx", error: String((e && e.message) || e) }); }

// components/status/Sparkline.jsx
try { (() => {
/**
 * Minimal response-time sparkline (SVG line + soft area fill).
 * `points` is an array of numbers (ms); scales to its own min/max.
 */
function Sparkline({
  points = [],
  width = 120,
  height = 28,
  stroke = "var(--chart-line)",
  fill = "var(--chart-fill)",
  style
}) {
  if (points.length < 2) return /*#__PURE__*/React.createElement("svg", {
    width: width,
    height: height,
    style: style,
    "aria-hidden": "true"
  });
  const min = Math.min(...points);
  const max = Math.max(...points);
  const span = max - min || 1;
  const step = width / (points.length - 1);
  const y = v => 2 + (height - 4) * (1 - (v - min) / span);
  const pts = points.map((v, i) => `${(i * step).toFixed(1)},${y(v).toFixed(1)}`);
  return /*#__PURE__*/React.createElement("svg", {
    width: width,
    height: height,
    style: style,
    role: "img",
    "aria-label": `Response time trend, ${min}–${max} ms`
  }, /*#__PURE__*/React.createElement("polygon", {
    points: `0,${height} ${pts.join(" ")} ${width},${height}`,
    fill: fill,
    stroke: "none"
  }), /*#__PURE__*/React.createElement("polyline", {
    points: pts.join(" "),
    fill: "none",
    stroke: stroke,
    strokeWidth: "1.5",
    strokeLinejoin: "round"
  }));
}
Object.assign(__ds_scope, { Sparkline });
})(); } catch (e) { __ds_ns.__errors.push({ path: "components/status/Sparkline.jsx", error: String((e && e.message) || e) }); }

// components/status/StatusBadge.jsx
try { (() => {
const MAP = {
  up: {
    icon: "check",
    label: "Up"
  },
  down: {
    icon: "alert-triangle",
    label: "Down"
  },
  degraded: {
    icon: "zap",
    label: "Degraded"
  },
  unknown: {
    icon: "pause",
    label: "Unknown"
  }
};

/**
 * Status pill: icon + text, never color alone. Pass `label` to override the
 * default text (e.g. "Operational", "Paused", "3 failing").
 */
function StatusBadge({
  status = "unknown",
  label,
  style
}) {
  const m = MAP[status] || MAP.unknown;
  return /*#__PURE__*/React.createElement("span", {
    className: "hm-badge",
    "data-status": status,
    style: style
  }, /*#__PURE__*/React.createElement(__ds_scope.Icon, {
    name: m.icon,
    size: 12
  }), label || m.label);
}
Object.assign(__ds_scope, { StatusBadge });
})(); } catch (e) { __ds_ns.__errors.push({ path: "components/status/StatusBadge.jsx", error: String((e && e.message) || e) }); }

// components/status/StatusDot.jsx
try { (() => {
/**
 * Small colored status dot. `pulse` adds a radiating ring — reserved for
 * ongoing-down states so motion means "needs attention".
 */
function StatusDot({
  status = "unknown",
  size = "md",
  pulse = false,
  style
}) {
  return /*#__PURE__*/React.createElement("span", {
    className: `hm-dot${size === "lg" ? " hm-dot--lg" : ""}${pulse ? " hm-dot--pulse" : ""}`,
    "data-status": status,
    style: style
  });
}
Object.assign(__ds_scope, { StatusDot });
})(); } catch (e) { __ds_ns.__errors.push({ path: "components/status/StatusDot.jsx", error: String((e && e.message) || e) }); }

// components/status/TickBar.jsx
try { (() => {
/**
 * UptimeRobot-style tick bar of the last N checks (oldest left, newest right).
 * Each check: { s: "up"|"down"|"degraded"|"unknown", tip?: string }.
 * Pads with empty ticks on the left when fewer than `slots` checks exist.
 */
function TickBar({
  checks = [],
  slots = 45,
  height = 20,
  style,
  ariaLabel
}) {
  const pad = Math.max(0, slots - checks.length);
  const shown = checks.slice(-slots);
  const downs = shown.filter(c => c.s === "down").length;
  return /*#__PURE__*/React.createElement("div", {
    className: "hm-ticks",
    style: {
      height,
      ...style
    },
    role: "img",
    "aria-label": ariaLabel || `Last ${shown.length} checks, ${downs} failed`
  }, Array.from({
    length: pad
  }).map((_, i) => /*#__PURE__*/React.createElement("span", {
    key: `p${i}`,
    className: "hm-tick"
  })), shown.map((c, i) => /*#__PURE__*/React.createElement("span", {
    key: i,
    className: "hm-tick",
    "data-s": c.s,
    "data-tip": c.tip || undefined
  })));
}
Object.assign(__ds_scope, { TickBar });
})(); } catch (e) { __ds_ns.__errors.push({ path: "components/status/TickBar.jsx", error: String((e && e.message) || e) }); }

// components/data/TargetCard.jsx
try { (() => {
const TYPE_ICON = {
  http: "globe",
  tcp: "server",
  dns: "at-sign"
};

/**
 * Target card (comfortable grid view): status dot + name + type icon,
 * tick bar of recent checks, uptime % + last response time.
 * A DOWN target gets a red border + pulsing dot.
 */
function TargetCard({
  target,
  onOpen,
  onEdit,
  onDelete,
  style
}) {
  const {
    name,
    type = "http",
    status = "unknown",
    checks = [],
    uptime,
    lastMs,
    endpoint,
    paused
  } = target;
  const down = status === "down";
  return /*#__PURE__*/React.createElement("div", {
    className: `hm-card hm-card--hover${down ? " hm-card--down" : ""}`,
    style: {
      padding: "var(--space-3) var(--space-4)",
      display: "flex",
      flexDirection: "column",
      gap: "var(--space-2)",
      ...style
    },
    onClick: onOpen,
    role: "button",
    tabIndex: 0,
    onKeyDown: e => {
      if (e.key === "Enter" && onOpen) onOpen();
    },
    "aria-label": `${name}, ${paused ? "paused" : status}`
  }, /*#__PURE__*/React.createElement("div", {
    style: {
      display: "flex",
      alignItems: "center",
      gap: "var(--space-2)"
    }
  }, /*#__PURE__*/React.createElement(__ds_scope.StatusDot, {
    status: paused ? "unknown" : status,
    pulse: down
  }), /*#__PURE__*/React.createElement("b", {
    style: {
      fontSize: "var(--text-lg)",
      fontWeight: "var(--weight-semibold)",
      overflow: "hidden",
      textOverflow: "ellipsis",
      whiteSpace: "nowrap"
    }
  }, name), /*#__PURE__*/React.createElement(__ds_scope.Icon, {
    name: TYPE_ICON[type] || "globe",
    size: 13,
    style: {
      color: "var(--text-3)"
    },
    "data-tip": type.toUpperCase()
  }), /*#__PURE__*/React.createElement("span", {
    style: {
      flex: 1
    }
  }), /*#__PURE__*/React.createElement("div", {
    style: {
      display: "flex",
      gap: 2
    },
    onClick: e => e.stopPropagation()
  }, onEdit && /*#__PURE__*/React.createElement("button", {
    type: "button",
    className: "hm-iconbtn",
    "aria-label": `Edit ${name}`,
    onClick: onEdit
  }, /*#__PURE__*/React.createElement(__ds_scope.Icon, {
    name: "pencil",
    size: 14
  })), onDelete && /*#__PURE__*/React.createElement("button", {
    type: "button",
    className: "hm-iconbtn hm-iconbtn--danger",
    "aria-label": `Delete ${name}`,
    onClick: onDelete
  }, /*#__PURE__*/React.createElement(__ds_scope.Icon, {
    name: "trash",
    size: 14
  })))), endpoint && /*#__PURE__*/React.createElement("span", {
    className: "num",
    style: {
      fontSize: "var(--text-sm)",
      color: "var(--text-3)",
      overflow: "hidden",
      textOverflow: "ellipsis",
      whiteSpace: "nowrap"
    }
  }, endpoint), /*#__PURE__*/React.createElement(__ds_scope.TickBar, {
    checks: checks,
    slots: 45,
    height: 18
  }), /*#__PURE__*/React.createElement("div", {
    style: {
      display: "flex",
      alignItems: "baseline",
      gap: "var(--space-3)",
      fontSize: "var(--text-sm)",
      color: "var(--text-2)"
    }
  }, /*#__PURE__*/React.createElement("span", {
    className: "num",
    style: {
      color: down ? "var(--status-down)" : "var(--status-up)",
      fontWeight: "var(--weight-semibold)"
    }
  }, paused ? "—" : `${uptime}%`), /*#__PURE__*/React.createElement("span", {
    className: "num"
  }, lastMs != null ? `${lastMs} ms` : "—"), paused && /*#__PURE__*/React.createElement("span", {
    style: {
      color: "var(--text-3)",
      display: "inline-flex",
      alignItems: "center",
      gap: 4
    }
  }, /*#__PURE__*/React.createElement(__ds_scope.Icon, {
    name: "pause",
    size: 11
  }), "paused")));
}
Object.assign(__ds_scope, { TargetCard });
})(); } catch (e) { __ds_ns.__errors.push({ path: "components/data/TargetCard.jsx", error: String((e && e.message) || e) }); }

// components/data/TargetTable.jsx
try { (() => {
const TYPE_ICON = {
  http: "globe",
  tcp: "server",
  dns: "at-sign"
};

/**
 * Compact table of targets (dense alternative to the card grid).
 * Columns: status, name, type, tick bar, uptime %, last ms, actions.
 */
function TargetTable({
  targets = [],
  onOpen,
  onEdit,
  onDelete,
  style
}) {
  return /*#__PURE__*/React.createElement("table", {
    className: "hm-table",
    style: style
  }, /*#__PURE__*/React.createElement("thead", null, /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("th", {
    style: {
      width: 24
    }
  }, /*#__PURE__*/React.createElement("span", {
    style: {
      position: "absolute",
      width: 1,
      height: 1,
      overflow: "hidden",
      clip: "rect(0 0 0 0)"
    }
  }, "Status")), /*#__PURE__*/React.createElement("th", null, "Name"), /*#__PURE__*/React.createElement("th", null, "Type"), /*#__PURE__*/React.createElement("th", {
    style: {
      width: 200
    }
  }, "Last 45 checks"), /*#__PURE__*/React.createElement("th", {
    style: {
      textAlign: "right"
    }
  }, "Uptime"), /*#__PURE__*/React.createElement("th", {
    style: {
      textAlign: "right"
    }
  }, "Resp."), /*#__PURE__*/React.createElement("th", {
    style: {
      width: 76
    }
  }))), /*#__PURE__*/React.createElement("tbody", null, targets.map((t, i) => {
    const down = t.status === "down";
    return /*#__PURE__*/React.createElement("tr", {
      key: t.id || i,
      "data-status": t.paused ? "unknown" : t.status,
      onClick: () => onOpen && onOpen(t),
      tabIndex: 0,
      onKeyDown: e => {
        if (e.key === "Enter" && onOpen) onOpen(t);
      }
    }, /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement(__ds_scope.StatusDot, {
      status: t.paused ? "unknown" : t.status,
      pulse: down
    })), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("b", {
      style: {
        fontWeight: "var(--weight-semibold)"
      }
    }, t.name), t.paused && /*#__PURE__*/React.createElement("span", {
      style: {
        marginLeft: 8,
        fontSize: "var(--text-xs)",
        color: "var(--text-3)"
      }
    }, "paused")), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("span", {
      style: {
        display: "inline-flex",
        alignItems: "center",
        gap: 5,
        color: "var(--text-2)",
        fontSize: "var(--text-sm)"
      }
    }, /*#__PURE__*/React.createElement(__ds_scope.Icon, {
      name: TYPE_ICON[t.type] || "globe",
      size: 13
    }), (t.type || "http").toUpperCase())), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement(__ds_scope.TickBar, {
      checks: t.checks || [],
      slots: 45,
      height: 14
    })), /*#__PURE__*/React.createElement("td", {
      className: "num",
      style: {
        textAlign: "right",
        color: down ? "var(--status-down)" : "var(--status-up)",
        fontWeight: "var(--weight-medium)"
      }
    }, t.paused ? "—" : `${t.uptime}%`), /*#__PURE__*/React.createElement("td", {
      className: "num",
      style: {
        textAlign: "right",
        color: "var(--text-2)"
      }
    }, t.lastMs != null ? `${t.lastMs} ms` : "—"), /*#__PURE__*/React.createElement("td", {
      onClick: e => e.stopPropagation()
    }, /*#__PURE__*/React.createElement("div", {
      style: {
        display: "flex",
        gap: 2,
        justifyContent: "flex-end"
      }
    }, onEdit && /*#__PURE__*/React.createElement("button", {
      type: "button",
      className: "hm-iconbtn",
      "aria-label": `Edit ${t.name}`,
      onClick: () => onEdit(t)
    }, /*#__PURE__*/React.createElement(__ds_scope.Icon, {
      name: "pencil",
      size: 14
    })), onDelete && /*#__PURE__*/React.createElement("button", {
      type: "button",
      className: "hm-iconbtn hm-iconbtn--danger",
      "aria-label": `Delete ${t.name}`,
      onClick: () => onDelete(t)
    }, /*#__PURE__*/React.createElement(__ds_scope.Icon, {
      name: "trash",
      size: 14
    })))));
  })));
}
Object.assign(__ds_scope, { TargetTable });
})(); } catch (e) { __ds_ns.__errors.push({ path: "components/data/TargetTable.jsx", error: String((e && e.message) || e) }); }

// ui_kits/dashboard/app.jsx
try { (() => {
// Dashboard app shell. Reads flags: HM_SCENARIO ("ok"|"alarm"), HM_TAB,
// HM_OPEN ("detail"|"target-form"|"notifier-form"), HM_VIEW ("cards"|"table").
const ADS = window.HealthMonitorRedesign_502b6f;
const {
  Icon: AIcon,
  StatusDot: ADot,
  StatusBadge: ABadge,
  LiveIndicator,
  StatCard: AStat,
  TargetCard: ACard,
  TargetTable: ATable,
  IncidentItem: AInc,
  NotifierRow: ANotif,
  Toast: AToast,
  ConfirmDialog: AConfirm,
  SlideOver: APanel,
  EmptyState: AEmpty,
  Button: ABtn,
  Segmented: ASeg,
  SelectField: ASelect
} = ADS;
function Hero({
  data
}) {
  const downs = data.targets.filter(t => t.status === "down");
  const degraded = data.targets.filter(t => t.status === "degraded");
  const ok = downs.length === 0;
  return /*#__PURE__*/React.createElement("div", {
    className: "hm-hero",
    "data-status": ok ? "up" : "down",
    role: "status",
    "aria-live": "polite",
    style: {
      marginTop: "var(--space-4)"
    }
  }, /*#__PURE__*/React.createElement("span", {
    className: "hm-hero-icon"
  }, /*#__PURE__*/React.createElement(AIcon, {
    name: ok ? "check" : "alert-triangle",
    size: 22
  })), /*#__PURE__*/React.createElement("div", {
    style: {
      flex: 1,
      minWidth: 0
    }
  }, /*#__PURE__*/React.createElement("h1", null, ok ? "All systems operational" : `${downs.length} target${downs.length > 1 ? "s" : ""} down`), /*#__PURE__*/React.createElement("p", null, ok ? `${data.targets.filter(t => !t.paused).length} targets · all checks passing` : downs.map(t => t.name).join(", ") + (degraded.length ? ` · ${degraded.length} degraded` : ""))), !ok && /*#__PURE__*/React.createElement("div", {
    style: {
      display: "flex",
      gap: 8,
      flexWrap: "wrap"
    }
  }, downs.map(t => /*#__PURE__*/React.createElement(ABadge, {
    key: t.id,
    status: "down",
    label: t.name
  }))));
}
function App() {
  const scenario = window.HM_SCENARIO || "alarm";
  const data = React.useMemo(() => window.HMData(scenario), [scenario]);
  const [tab, setTab] = React.useState(window.HM_TAB || "dashboard");
  const [view, setView] = React.useState(window.HM_VIEW || "cards");
  const [query, setQuery] = React.useState("");
  const [filter, setFilter] = React.useState("all");
  const [sort, setSort] = React.useState("status");
  const [theme, setTheme] = React.useState(document.documentElement.getAttribute("data-theme") || "dark");
  const [showResolved, setShowResolved] = React.useState(false);
  const [detail, setDetail] = React.useState(window.HM_OPEN === "detail" ? data.targets[2] : null);
  const [targetForm, setTargetForm] = React.useState(window.HM_OPEN === "target-form" ? {
    mode: "add"
  } : null);
  const [notifierForm, setNotifierForm] = React.useState(window.HM_OPEN === "notifier-form" ? {
    mode: "add"
  } : null);
  const [confirm, setConfirm] = React.useState(null);
  const [toasts, setToasts] = React.useState([]);
  const toast = (kind, title, message) => {
    const id = Date.now() + Math.random();
    setToasts(ts => [...ts, {
      id,
      kind,
      title,
      message
    }]);
    setTimeout(() => setToasts(ts => ts.filter(t => t.id !== id)), 5200);
  };
  const setThemeAttr = next => {
    setTheme(next);
    document.documentElement.setAttribute("data-theme", next);
  };
  const counts = {
    all: data.targets.length,
    down: data.targets.filter(t => t.status === "down").length,
    degraded: data.targets.filter(t => t.status === "degraded").length,
    paused: data.targets.filter(t => t.paused).length
  };
  const order = {
    down: 0,
    degraded: 1,
    up: 2,
    unknown: 3
  };
  const shown = data.targets.filter(t => {
    if (filter === "down" && t.status !== "down") return false;
    if (filter === "degraded" && t.status !== "degraded") return false;
    if (filter === "paused" && !t.paused) return false;
    const q = query.trim().toLowerCase();
    return !q || t.name.toLowerCase().includes(q) || t.endpoint.toLowerCase().includes(q) || t.tags.some(x => x.includes(q));
  }).sort((a, b) => sort === "name" ? a.name.localeCompare(b.name) : sort === "uptime" ? Number(a.uptime) - Number(b.uptime) : sort === "response" ? (b.lastMs || 0) - (a.lastMs || 0) : order[a.status] - order[b.status] || a.name.localeCompare(b.name));
  const askDelete = t => setConfirm({
    title: "Delete target?",
    body: `“${t.name}” and its check history will be permanently removed. Notifiers are not affected.`,
    onConfirm: () => {
      setConfirm(null);
      toast("info", "Target deleted", `${t.name} is no longer monitored`);
    }
  });
  const avgMs = Math.round(data.targets.filter(t => t.lastMs).reduce((a, t) => a + t.lastMs, 0) / data.targets.filter(t => t.lastMs).length);
  return /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("header", {
    className: "hm-topbar"
  }, /*#__PURE__*/React.createElement("span", {
    className: "hm-brand"
  }, /*#__PURE__*/React.createElement("span", {
    className: "hm-brand-mark"
  }, /*#__PURE__*/React.createElement(AIcon, {
    name: "activity",
    size: 15
  })), "Health Monitor"), /*#__PURE__*/React.createElement(ASeg, {
    ariaLabel: "Section",
    options: [{
      value: "dashboard",
      label: "Dashboard"
    }, {
      value: "notifiers",
      label: "Notifiers",
      count: data.notifiers.length
    }],
    value: tab,
    onChange: setTab
  }), /*#__PURE__*/React.createElement("span", {
    style: {
      flex: 1
    }
  }), /*#__PURE__*/React.createElement(LiveIndicator, {
    state: "live"
  }), /*#__PURE__*/React.createElement("button", {
    type: "button",
    className: "hm-iconbtn",
    "aria-label": `Switch to ${theme === "dark" ? "light" : "dark"} theme`,
    onClick: () => setThemeAttr(theme === "dark" ? "light" : "dark")
  }, /*#__PURE__*/React.createElement(AIcon, {
    name: theme === "dark" ? "sun" : "moon",
    size: 15
  })), /*#__PURE__*/React.createElement(ABtn, {
    variant: "primary",
    icon: /*#__PURE__*/React.createElement(AIcon, {
      name: "plus",
      size: 15
    }),
    onClick: () => setTargetForm({
      mode: "add"
    })
  }, /*#__PURE__*/React.createElement("span", {
    className: "hm-topbar-add-label"
  }, "Add target"))), /*#__PURE__*/React.createElement("main", {
    className: "hm-shell"
  }, tab === "dashboard" && /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement(Hero, {
    data: data
  }), /*#__PURE__*/React.createElement("div", {
    className: "hm-statrow"
  }, /*#__PURE__*/React.createElement(AStat, {
    label: "Targets",
    value: counts.all,
    sub: `${counts.all - counts.paused} enabled`
  }), /*#__PURE__*/React.createElement(AStat, {
    label: "Up",
    value: counts.all - counts.down - counts.degraded - counts.paused,
    tone: "success"
  }), /*#__PURE__*/React.createElement(AStat, {
    label: "Down",
    value: counts.down,
    tone: counts.down ? "danger" : "default",
    sub: counts.down ? "see incidents" : "\u00a0"
  }), /*#__PURE__*/React.createElement(AStat, {
    label: "Avg response",
    value: avgMs,
    unit: "ms",
    sub: "across all targets"
  })), data.incidents.ongoing.length > 0 && /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("div", {
    className: "hm-sectionhead"
  }, /*#__PURE__*/React.createElement("h2", null, "Ongoing incidents"), /*#__PURE__*/React.createElement("span", {
    className: "hm-badge",
    "data-status": "down"
  }, /*#__PURE__*/React.createElement(AIcon, {
    name: "alert-triangle",
    size: 12
  }), data.incidents.ongoing.length, " active")), /*#__PURE__*/React.createElement("div", {
    className: "hm-card hm-card--down"
  }, data.incidents.ongoing.map((inc, i) => /*#__PURE__*/React.createElement(AInc, {
    key: inc.id,
    incident: inc,
    last: i === data.incidents.ongoing.length - 1,
    onOpen: () => setDetail(data.targets.find(t => t.id === inc.targetId))
  })))), /*#__PURE__*/React.createElement("div", {
    className: "hm-toolbar"
  }, /*#__PURE__*/React.createElement("div", {
    className: "hm-search"
  }, /*#__PURE__*/React.createElement(AIcon, {
    name: "search",
    size: 14
  }), /*#__PURE__*/React.createElement("input", {
    className: "hm-input",
    placeholder: "Search targets, tags\u2026",
    "aria-label": "Search targets",
    value: query,
    onChange: e => setQuery(e.target.value)
  })), /*#__PURE__*/React.createElement(ASeg, {
    ariaLabel: "Status filter",
    options: [{
      value: "all",
      label: "All",
      count: counts.all
    }, {
      value: "down",
      label: "Down",
      count: counts.down
    }, {
      value: "degraded",
      label: "Degraded",
      count: counts.degraded
    }, {
      value: "paused",
      label: "Paused",
      count: counts.paused
    }],
    value: filter,
    onChange: setFilter
  }), /*#__PURE__*/React.createElement("span", {
    style: {
      flex: 1
    }
  }), /*#__PURE__*/React.createElement("select", {
    className: "hm-select",
    style: {
      width: "auto",
      minHeight: 30
    },
    "aria-label": "Sort targets",
    value: sort,
    onChange: e => setSort(e.target.value)
  }, /*#__PURE__*/React.createElement("option", {
    value: "status"
  }, "Sort: status"), /*#__PURE__*/React.createElement("option", {
    value: "name"
  }, "Sort: name"), /*#__PURE__*/React.createElement("option", {
    value: "uptime"
  }, "Sort: uptime"), /*#__PURE__*/React.createElement("option", {
    value: "response"
  }, "Sort: response time")), /*#__PURE__*/React.createElement(ASeg, {
    ariaLabel: "View",
    options: [{
      value: "cards",
      label: "",
      icon: /*#__PURE__*/React.createElement(AIcon, {
        name: "grid",
        size: 13,
        "aria-label": "Card view"
      })
    }, {
      value: "table",
      label: "",
      icon: /*#__PURE__*/React.createElement(AIcon, {
        name: "list",
        size: 13,
        "aria-label": "Table view"
      })
    }],
    value: view,
    onChange: setView
  })), shown.length === 0 ? /*#__PURE__*/React.createElement("div", {
    className: "hm-card"
  }, /*#__PURE__*/React.createElement(AEmpty, {
    icon: "search",
    title: "No matching targets",
    body: query ? `Nothing matches “${query}”. Try a different search or clear the filter.` : "Nothing matches this filter.",
    action: /*#__PURE__*/React.createElement(ABtn, {
      onClick: () => {
        setQuery("");
        setFilter("all");
      }
    }, "Clear filters")
  })) : view === "cards" ? /*#__PURE__*/React.createElement("div", {
    className: "hm-grid"
  }, shown.map(t => /*#__PURE__*/React.createElement(ACard, {
    key: t.id,
    target: t,
    onOpen: () => setDetail(t),
    onEdit: () => setTargetForm({
      mode: "edit",
      target: t
    }),
    onDelete: () => askDelete(t)
  }))) : /*#__PURE__*/React.createElement("div", {
    className: "hm-card hm-tablewrap scroll-thin"
  }, /*#__PURE__*/React.createElement(ATable, {
    targets: shown,
    onOpen: t => setDetail(t),
    onEdit: t => setTargetForm({
      mode: "edit",
      target: t
    }),
    onDelete: t => askDelete(t)
  })), /*#__PURE__*/React.createElement("div", {
    className: "hm-sectionhead"
  }, /*#__PURE__*/React.createElement("h2", null, "Incident history"), /*#__PURE__*/React.createElement("span", {
    style: {
      flex: 1
    }
  }), /*#__PURE__*/React.createElement(ABtn, {
    variant: "ghost",
    icon: /*#__PURE__*/React.createElement(AIcon, {
      name: showResolved ? "chevron-down" : "chevron-right",
      size: 14
    }),
    onClick: () => setShowResolved(!showResolved)
  }, data.incidents.resolved.length, " resolved")), showResolved && /*#__PURE__*/React.createElement("div", {
    className: "hm-card"
  }, data.incidents.resolved.map((inc, i) => /*#__PURE__*/React.createElement(AInc, {
    key: inc.id,
    incident: inc,
    last: i === data.incidents.resolved.length - 1,
    onOpen: () => setDetail(data.targets.find(t => t.id === inc.targetId))
  }))), /*#__PURE__*/React.createElement("div", {
    className: "hm-foot"
  }, /*#__PURE__*/React.createElement(LiveIndicator, {
    state: "live"
  }), /*#__PURE__*/React.createElement("span", null, "Last updated ", /*#__PURE__*/React.createElement("span", {
    className: "num"
  }, "12s"), " ago"))), tab === "notifiers" && /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("div", {
    className: "hm-sectionhead",
    style: {
      marginTop: "var(--space-6)"
    }
  }, /*#__PURE__*/React.createElement("h2", null, "Notifiers"), /*#__PURE__*/React.createElement("span", {
    style: {
      flex: 1
    }
  }), /*#__PURE__*/React.createElement(ABtn, {
    variant: "primary",
    icon: /*#__PURE__*/React.createElement(AIcon, {
      name: "plus",
      size: 15
    }),
    onClick: () => setNotifierForm({
      mode: "add"
    })
  }, "Add notifier")), /*#__PURE__*/React.createElement("p", {
    style: {
      margin: "0 0 var(--space-3)",
      color: "var(--text-2)",
      fontSize: "var(--text-md)",
      maxWidth: 560
    }
  }, "Alert channels for down / recovered / slow-response events. Rarely touched \u2014 they live here, out of the way of live status."), /*#__PURE__*/React.createElement("div", {
    className: "hm-card"
  }, data.notifiers.map((n, i) => /*#__PURE__*/React.createElement("div", {
    key: n.id,
    style: {
      borderTop: i ? "1px solid var(--border-1)" : "none"
    }
  }, /*#__PURE__*/React.createElement(ANotif, {
    notifier: n,
    onToggle: on => toast(on ? "success" : "info", on ? "Notifier enabled" : "Notifier disabled", n.name),
    onEdit: () => setNotifierForm({
      mode: "edit",
      notifier: n
    }),
    onDelete: () => setConfirm({
      title: "Delete notifier?",
      body: `“${n.name}” will stop receiving alerts.`,
      onConfirm: () => {
        setConfirm(null);
        toast("info", "Notifier deleted", n.name);
      }
    })
  })))))), detail && /*#__PURE__*/React.createElement(APanel, {
    width: 720,
    title: /*#__PURE__*/React.createElement(window.DetailTitle, {
      target: detail
    }),
    titleExtra: /*#__PURE__*/React.createElement(ABadge, {
      status: detail.status,
      label: detail.paused ? "Paused" : undefined
    }),
    onClose: () => setDetail(null)
  }, /*#__PURE__*/React.createElement(window.TargetDetailBody, {
    target: detail,
    incidents: data.incidents,
    onEdit: () => {
      setDetail(null);
      setTargetForm({
        mode: "edit",
        target: detail
      });
    }
  })), targetForm && /*#__PURE__*/React.createElement(APanel, {
    width: 520,
    title: targetForm.mode === "add" ? "Add target" : `Edit ${targetForm.target.name}`,
    onClose: () => setTargetForm(null),
    footer: /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement(ABtn, {
      onClick: () => setTargetForm(null)
    }, "Cancel"), /*#__PURE__*/React.createElement(ABtn, {
      variant: "primary",
      onClick: () => {
        setTargetForm(null);
        toast("success", targetForm.mode === "add" ? "Target added" : "Target saved", "Checks start on the next tick");
      }
    }, targetForm.mode === "add" ? "Add target" : "Save changes"))
  }, /*#__PURE__*/React.createElement(window.TargetForm, {
    target: targetForm.target
  })), notifierForm && /*#__PURE__*/React.createElement(APanel, {
    width: 520,
    title: notifierForm.mode === "add" ? "Add notifier" : `Edit ${notifierForm.notifier.name}`,
    onClose: () => setNotifierForm(null),
    footer: /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement(ABtn, {
      onClick: () => setNotifierForm(null)
    }, "Cancel"), /*#__PURE__*/React.createElement(ABtn, {
      variant: "primary",
      onClick: () => {
        setNotifierForm(null);
        toast("success", "Notifier saved", "A test message was sent");
      }
    }, "Save notifier"))
  }, /*#__PURE__*/React.createElement(window.NotifierForm, {
    notifier: notifierForm.notifier
  })), confirm && /*#__PURE__*/React.createElement(AConfirm, {
    title: confirm.title,
    body: confirm.body,
    onConfirm: confirm.onConfirm,
    onCancel: () => setConfirm(null)
  }), /*#__PURE__*/React.createElement("div", {
    className: "hm-toaststack",
    "aria-live": "polite"
  }, toasts.map(t => /*#__PURE__*/React.createElement(AToast, {
    key: t.id,
    kind: t.kind,
    title: t.title,
    message: t.message,
    onDismiss: () => setToasts(ts => ts.filter(x => x.id !== t.id))
  }))));
}
window.HMApp = App;
})(); } catch (e) { __ds_ns.__errors.push({ path: "ui_kits/dashboard/app.jsx", error: String((e && e.message) || e) }); }

// ui_kits/dashboard/data.js
try { (() => {
// Health Monitor UI kit — deterministic sample data.
// Scenario: "ok" (all up) or "alarm" (2 down, 1 degraded). window.HM_SCENARIO picks it.
window.HMData = (() => {
  function rng(seed) {
    let s = seed;
    return () => (s = (s * 1103515245 + 12345) % 2147483648) / 2147483648;
  }
  const defs = [{
    id: "t1",
    name: "API /health",
    type: "http",
    endpoint: "https://api.example.dev/health",
    interval: "30s",
    timeout: "5s",
    tags: ["prod", "api"],
    desc: "Main REST API liveness probe",
    base: 180
  }, {
    id: "t2",
    name: "Marketing site",
    type: "http",
    endpoint: "https://example.dev",
    interval: "1m",
    timeout: "10s",
    tags: ["prod", "web"],
    base: 320,
    alarm: "degraded"
  }, {
    id: "t3",
    name: "Postgres primary",
    type: "tcp",
    endpoint: "10.0.4.2:5432",
    interval: "30s",
    timeout: "3s",
    tags: ["prod", "db"],
    base: 12,
    alarm: "down",
    err: "dial tcp 10.0.4.2:5432: i/o timeout"
  }, {
    id: "t4",
    name: "Redis cache",
    type: "tcp",
    endpoint: "10.0.4.3:6379",
    interval: "30s",
    timeout: "3s",
    tags: ["prod", "cache"],
    base: 8
  }, {
    id: "t5",
    name: "Backup worker",
    type: "http",
    endpoint: "https://backup.example.dev/ping",
    interval: "5m",
    timeout: "10s",
    tags: ["infra"],
    base: 260,
    alarm: "down",
    err: "HTTP 503 Service Unavailable"
  }, {
    id: "t6",
    name: "example.dev DNS",
    type: "dns",
    endpoint: "example.dev · A",
    interval: "5m",
    timeout: "5s",
    tags: ["infra", "dns"],
    base: 45
  }, {
    id: "t7",
    name: "Staging API",
    type: "http",
    endpoint: "https://staging.example.dev/health",
    interval: "1m",
    timeout: "5s",
    tags: ["staging"],
    base: 210
  }, {
    id: "t8",
    name: "SMTP relay",
    type: "tcp",
    endpoint: "10.0.4.9:25",
    interval: "5m",
    timeout: "5s",
    tags: ["infra", "mail"],
    base: 30
  }, {
    id: "t9",
    name: "Status page",
    type: "http",
    endpoint: "https://status.example.dev",
    interval: "5m",
    timeout: "10s",
    tags: ["web"],
    base: 190,
    paused: true
  }, {
    id: "t10",
    name: "Grafana",
    type: "http",
    endpoint: "https://grafana.internal.example.dev",
    interval: "1m",
    timeout: "10s",
    tags: ["infra", "observability"],
    base: 150
  }, {
    id: "t11",
    name: "MX records",
    type: "dns",
    endpoint: "example.dev · MX",
    interval: "15m",
    timeout: "5s",
    tags: ["infra", "dns", "mail"],
    base: 52
  }, {
    id: "t12",
    name: "CDN edge",
    type: "http",
    endpoint: "https://cdn.example.dev/health",
    interval: "1m",
    timeout: "5s",
    tags: ["prod", "web"],
    base: 95
  }];
  function pad(n) {
    return String(n).padStart(2, "0");
  }
  function build(def, i, scenario) {
    const r = rng(i * 7919 + 13);
    const state = scenario === "alarm" ? def.alarm || "up" : "up";
    const ms = () => Math.round(def.base * (0.7 + r() * 0.9));

    // 45 recent ticks (~last 45 checks), newest right
    const checks = [];
    for (let k = 0; k < 45; k++) {
      const t = `12:${pad(k + 10)}`;
      let s = r() > 0.985 ? "degraded" : "up";
      if (state === "down" && k >= 41) s = "down";
      if (state === "degraded" && (k === 38 || k >= 43)) s = "degraded";
      checks.push({
        s,
        tip: s === "down" ? `${t} · ${def.err || "failed"}` : `${t} · ${ms()} ms${s === "degraded" ? " · slow" : ""}`
      });
    }

    // per-period stats + chart series
    const periods = {};
    [["24h", 48, k => `${pad(Math.floor(k / 2))}:${k % 2 ? "30" : "00"}`], ["7d", 56, k => ["Thu", "Fri", "Sat", "Sun", "Mon", "Tue", "Wed"][Math.floor(k / 8)]], ["30d", 60, k => `${pad(9 + Math.floor(k / 2) - 30 + 21)} ${Math.floor(k / 2) < 21 ? "Jun" : "Jul"}`]].forEach(([p, n, lbl]) => {
      const rr = rng(i * 104729 + p.length * 31);
      const series = [];
      let failed = 0;
      for (let k = 0; k < n; k++) {
        let s = rr() > 0.99 ? "degraded" : "up";
        let v = Math.round(def.base * (0.7 + rr() * 0.9));
        if (state === "down" && k >= n - 3) {
          s = "down";
          v = null;
          failed++;
        }
        if (state === "degraded" && k >= n - 4) {
          s = "degraded";
          v = Math.round(def.base * 3.2);
        }
        series.push({
          v,
          s,
          label: lbl(k)
        });
      }
      const vals = series.filter(x => x.v != null).map(x => x.v);
      const secs = {
        "24h": 86400,
        "7d": 604800,
        "30d": 2592000
      }[p];
      const total = Math.round(secs / ({
        "30s": 30,
        "1m": 60,
        "5m": 300,
        "15m": 900
      }[def.interval] || 60));
      let failTotal = state === "down" ? p === "24h" ? 14 : p === "7d" ? 14 : 22 : state === "degraded" ? p === "24h" ? 2 : 3 : p === "24h" ? rr() < 0.4 ? 1 : 0 : Math.round(rr() * (p === "7d" ? 4 : 7));
      failTotal = Math.min(failTotal, Math.max(total - 1, 0));
      periods[p] = {
        series,
        uptime: (100 * (1 - failTotal / total)).toFixed(p === "24h" ? 2 : 3),
        avg: Math.round(vals.reduce((a, b) => a + b, 0) / vals.length),
        min: Math.min.apply(null, vals),
        max: Math.max.apply(null, vals),
        totalChecks: total,
        failedChecks: failTotal
      };
    });
    const status = def.paused ? "unknown" : state;
    return Object.assign({}, def, {
      status,
      paused: !!def.paused,
      checks,
      periods,
      uptime: periods["24h"].uptime,
      lastMs: state === "down" ? null : checks[44].tip.match(/(\d+) ms/) ? Number(checks[44].tip.match(/(\d+) ms/)[1]) : def.base,
      consecutiveFailures: state === "down" ? 14 : 0,
      lastCheck: "12s ago"
    });
  }
  const notifiers = [{
    id: "n1",
    name: "Ops channel",
    type: "telegram",
    enabled: true,
    detail: "chat •••4821"
  }, {
    id: "n2",
    name: "On-call email",
    type: "email",
    enabled: true,
    detail: "oncall@•••.dev via smtp.example.dev:587"
  }, {
    id: "n3",
    name: "Founders Gmail",
    type: "gmail_oauth",
    enabled: false,
    detail: "alerts@•••.dev"
  }, {
    id: "n4",
    name: "PagerDuty bridge",
    type: "webhook",
    enabled: true,
    detail: "https://events.pagerduty.com/•••"
  }];
  function incidents(scenario) {
    const ongoing = scenario === "alarm" ? [{
      id: "i1",
      targetId: "t3",
      target: "Postgres primary",
      status: "ongoing",
      severity: "critical",
      startedAt: "Today, 09:14",
      duration: "23m",
      failureCount: 14,
      lastError: "dial tcp 10.0.4.2:5432: i/o timeout"
    }, {
      id: "i2",
      targetId: "t5",
      target: "Backup worker",
      status: "ongoing",
      severity: "warning",
      startedAt: "Today, 09:29",
      duration: "8m",
      failureCount: 2,
      lastError: "HTTP 503 Service Unavailable"
    }] : [];
    const resolved = [{
      id: "i3",
      targetId: "t2",
      target: "Marketing site",
      status: "resolved",
      severity: "warning",
      startedAt: "8 Jul, 22:03",
      duration: "4m",
      failureCount: 3,
      lastError: "HTTP 502 Bad Gateway"
    }, {
      id: "i4",
      targetId: "t12",
      target: "CDN edge",
      status: "resolved",
      severity: "critical",
      startedAt: "6 Jul, 03:41",
      duration: "18m",
      failureCount: 11,
      lastError: "context deadline exceeded"
    }, {
      id: "i5",
      targetId: "t6",
      target: "example.dev DNS",
      status: "resolved",
      severity: "info",
      startedAt: "1 Jul, 11:02",
      duration: "2m",
      failureCount: 1,
      lastError: "SERVFAIL"
    }];
    return {
      ongoing,
      resolved
    };
  }
  return function make(scenario) {
    return {
      scenario,
      targets: defs.map((d, i) => build(d, i, scenario)),
      notifiers,
      incidents: incidents(scenario)
    };
  };
})();
})(); } catch (e) { __ds_ns.__errors.push({ path: "ui_kits/dashboard/data.js", error: String((e && e.message) || e) }); }

// ui_kits/dashboard/detail.jsx
try { (() => {
// Target detail: period-driven stats + chart + incidents + checks + config.
const DS = window.HealthMonitorRedesign_502b6f;
const {
  Icon,
  StatusDot,
  StatusBadge,
  Segmented,
  Button,
  TickBar
} = DS;
const detailTypeIcon = {
  http: "globe",
  tcp: "server",
  dns: "at-sign"
};

/* ---- SVG response-time chart, points colored by status ---- */
function ResponseChart({
  series,
  height = 170
}) {
  const w = 640,
    padL = 44,
    padB = 22,
    padT = 10;
  const vals = series.filter(p => p.v != null).map(p => p.v);
  const max = Math.max(...vals) * 1.15,
    min = 0;
  const iw = w - padL - 8,
    ih = height - padT - padB;
  const x = i => padL + i / (series.length - 1) * iw;
  const y = v => padT + ih * (1 - (v - min) / (max - min));
  const segs = [];
  let cur = [];
  series.forEach((p, i) => {
    if (p.v == null) {
      if (cur.length > 1) segs.push(cur);
      cur = [];
    } else cur.push(`${x(i).toFixed(1)},${y(p.v).toFixed(1)}`);
  });
  if (cur.length > 1) segs.push(cur);
  const gridVals = [0.25, 0.5, 0.75, 1].map(f => Math.round(max * f));
  const labelEvery = Math.ceil(series.length / 6);
  return /*#__PURE__*/React.createElement("svg", {
    width: "100%",
    viewBox: `0 0 ${w} ${height}`,
    role: "img",
    "aria-label": "Response time chart",
    style: {
      display: "block"
    }
  }, gridVals.map(v => /*#__PURE__*/React.createElement("g", {
    key: v
  }, /*#__PURE__*/React.createElement("line", {
    x1: padL,
    x2: w - 8,
    y1: y(v),
    y2: y(v),
    stroke: "var(--chart-grid)",
    strokeWidth: "1"
  }), /*#__PURE__*/React.createElement("text", {
    x: padL - 6,
    y: y(v) + 3,
    textAnchor: "end",
    fontSize: "9",
    fill: "var(--text-3)",
    fontFamily: "var(--font-mono)"
  }, v))), segs.map((s, i) => /*#__PURE__*/React.createElement("polyline", {
    key: i,
    points: s.join(" "),
    fill: "none",
    stroke: "var(--chart-line)",
    strokeWidth: "1.5",
    strokeLinejoin: "round"
  })), series.map((p, i) => {
    if (p.v == null) {
      return /*#__PURE__*/React.createElement("rect", {
        key: i,
        x: x(i) - 2,
        y: padT,
        width: 4,
        height: ih,
        fill: "var(--status-down-subtle)",
        stroke: "var(--status-down)",
        strokeWidth: "0.5",
        rx: 1
      });
    }
    const bad = p.s !== "up";
    if (!bad && i % 2) return null;
    return /*#__PURE__*/React.createElement("circle", {
      key: i,
      cx: x(i),
      cy: y(p.v),
      r: bad ? 3 : 1.6,
      fill: bad ? "var(--status-degraded)" : "var(--chart-line)"
    });
  }), series.map((p, i) => i % labelEvery === 0 ? /*#__PURE__*/React.createElement("text", {
    key: `l${i}`,
    x: x(i),
    y: height - 6,
    textAnchor: "middle",
    fontSize: "9",
    fill: "var(--text-3)",
    fontFamily: "var(--font-mono)"
  }, p.label) : null));
}
function DetailStat({
  label,
  value,
  unit,
  tone
}) {
  return /*#__PURE__*/React.createElement("div", {
    style: {
      flex: 1,
      minWidth: 0,
      background: "var(--bg-inset)",
      border: "1px solid var(--border-1)",
      borderRadius: "var(--radius-md)",
      padding: "8px 12px"
    }
  }, /*#__PURE__*/React.createElement("span", {
    className: "overline",
    style: {
      display: "block"
    }
  }, label), /*#__PURE__*/React.createElement("b", {
    className: "num",
    style: {
      fontSize: "var(--text-2xl)",
      fontWeight: 600,
      color: tone || "var(--text-1)"
    }
  }, value, unit && /*#__PURE__*/React.createElement("span", {
    style: {
      fontSize: "var(--text-sm)",
      color: "var(--text-3)",
      marginLeft: 2
    }
  }, unit)));
}
function ConfigRow({
  k,
  v,
  mono
}) {
  return /*#__PURE__*/React.createElement("div", {
    style: {
      display: "grid",
      gridTemplateColumns: "140px 1fr",
      gap: 8,
      padding: "6px 0",
      borderBottom: "1px solid var(--border-1)",
      fontSize: "var(--text-md)"
    }
  }, /*#__PURE__*/React.createElement("span", {
    style: {
      color: "var(--text-3)"
    }
  }, k), /*#__PURE__*/React.createElement("span", {
    className: mono ? "num" : "",
    style: {
      color: "var(--text-1)",
      overflowWrap: "anywhere"
    }
  }, v));
}
function ChecksTable({
  target
}) {
  const rows = target.checks.slice(-8).reverse().map((c, i) => {
    const m = c.tip.match(/^(\S+) · (.*)$/);
    return {
      time: m ? m[1] : "",
      info: m ? m[2] : c.tip,
      s: c.s
    };
  });
  return /*#__PURE__*/React.createElement("table", {
    className: "hm-table"
  }, /*#__PURE__*/React.createElement("thead", null, /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("th", null, "Time"), /*#__PURE__*/React.createElement("th", null, "Result"), /*#__PURE__*/React.createElement("th", {
    style: {
      textAlign: "right"
    }
  }, "Detail"))), /*#__PURE__*/React.createElement("tbody", null, rows.map((r, i) => /*#__PURE__*/React.createElement("tr", {
    key: i,
    style: {
      cursor: "default"
    }
  }, /*#__PURE__*/React.createElement("td", {
    className: "num",
    style: {
      color: "var(--text-2)"
    }
  }, r.time), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("span", {
    style: {
      display: "inline-flex",
      alignItems: "center",
      gap: 6,
      color: r.s === "down" ? "var(--status-down)" : r.s === "degraded" ? "var(--status-degraded)" : "var(--status-up)",
      fontSize: "var(--text-sm)",
      fontWeight: 500
    }
  }, /*#__PURE__*/React.createElement(Icon, {
    name: r.s === "down" ? "x-circle" : r.s === "degraded" ? "zap" : "check-circle",
    size: 13
  }), r.s === "down" ? "Failure" : r.s === "degraded" ? "Slow" : "Success")), /*#__PURE__*/React.createElement("td", {
    className: "num",
    style: {
      textAlign: "right",
      color: "var(--text-2)",
      fontSize: "var(--text-sm)"
    }
  }, r.info)))));
}
function Section({
  title,
  children,
  extra
}) {
  return /*#__PURE__*/React.createElement("section", {
    style: {
      marginTop: "var(--space-5)"
    }
  }, /*#__PURE__*/React.createElement("div", {
    style: {
      display: "flex",
      alignItems: "center",
      gap: 8,
      marginBottom: "var(--space-2)"
    }
  }, /*#__PURE__*/React.createElement("span", {
    className: "overline"
  }, title), /*#__PURE__*/React.createElement("span", {
    style: {
      flex: 1
    }
  }), extra), children);
}

/* ---- shared detail body (used by slide-over AND full-page variant) ---- */
function TargetDetailBody({
  target,
  incidents,
  onEdit
}) {
  const [period, setPeriod] = React.useState("24h");
  const p = target.periods[period];
  const down = target.status === "down";
  const tIncidents = incidents.ongoing.concat(incidents.resolved).filter(i => i.targetId === target.id);
  return /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement("div", {
    style: {
      display: "flex",
      alignItems: "center",
      gap: 8,
      flexWrap: "wrap"
    }
  }, /*#__PURE__*/React.createElement(Segmented, {
    options: ["24h", "7d", "30d"],
    value: period,
    onChange: setPeriod,
    ariaLabel: "Stats period"
  }), /*#__PURE__*/React.createElement("span", {
    style: {
      flex: 1
    }
  }), /*#__PURE__*/React.createElement("span", {
    className: "hm-live",
    "data-state": "live"
  }, /*#__PURE__*/React.createElement("span", {
    className: "hm-dot"
  }), "Checked ", target.lastCheck, " \xB7 every ", target.interval)), /*#__PURE__*/React.createElement("div", {
    style: {
      display: "flex",
      gap: 8,
      marginTop: "var(--space-3)",
      flexWrap: "wrap"
    }
  }, /*#__PURE__*/React.createElement(DetailStat, {
    label: "Uptime",
    value: p.uptime,
    unit: "%",
    tone: down ? "var(--status-down)" : "var(--status-up)"
  }), /*#__PURE__*/React.createElement(DetailStat, {
    label: "Avg",
    value: p.avg,
    unit: "ms"
  }), /*#__PURE__*/React.createElement(DetailStat, {
    label: "Min / Max",
    value: `${p.min}–${p.max}`,
    unit: "ms"
  }), /*#__PURE__*/React.createElement(DetailStat, {
    label: "Failed",
    value: p.failedChecks,
    unit: `/ ${p.totalChecks}`,
    tone: p.failedChecks ? "var(--status-down)" : undefined
  })), /*#__PURE__*/React.createElement(Section, {
    title: `Response time — ${period}`
  }, /*#__PURE__*/React.createElement("div", {
    style: {
      background: "var(--bg-inset)",
      border: "1px solid var(--border-1)",
      borderRadius: "var(--radius-md)",
      padding: "10px 8px 4px"
    }
  }, /*#__PURE__*/React.createElement(ResponseChart, {
    series: p.series
  }))), /*#__PURE__*/React.createElement(Section, {
    title: "Recent checks"
  }, /*#__PURE__*/React.createElement(TickBar, {
    checks: target.checks,
    slots: 45,
    height: 16,
    style: {
      marginBottom: 8
    }
  }), /*#__PURE__*/React.createElement("div", {
    className: "hm-card",
    style: {
      overflow: "hidden",
      borderRadius: "var(--radius-md)"
    }
  }, /*#__PURE__*/React.createElement(ChecksTable, {
    target: target
  }))), /*#__PURE__*/React.createElement(Section, {
    title: "Incidents"
  }, tIncidents.length ? /*#__PURE__*/React.createElement("div", {
    className: "hm-card",
    style: {
      borderRadius: "var(--radius-md)"
    }
  }, tIncidents.map((inc, i) => /*#__PURE__*/React.createElement(DS.IncidentItem, {
    key: inc.id,
    incident: inc,
    last: i === tIncidents.length - 1
  }))) : /*#__PURE__*/React.createElement("p", {
    style: {
      margin: 0,
      fontSize: "var(--text-md)",
      color: "var(--text-3)"
    }
  }, "No incidents in the last 30 days.")), /*#__PURE__*/React.createElement(Section, {
    title: "Configuration",
    extra: /*#__PURE__*/React.createElement(Button, {
      variant: "ghost",
      icon: /*#__PURE__*/React.createElement(Icon, {
        name: "pencil",
        size: 13
      }),
      onClick: onEdit
    }, "Edit")
  }, /*#__PURE__*/React.createElement("div", {
    style: {
      background: "var(--bg-inset)",
      border: "1px solid var(--border-1)",
      borderRadius: "var(--radius-md)",
      padding: "4px 14px"
    }
  }, /*#__PURE__*/React.createElement(ConfigRow, {
    k: "Type",
    v: target.type.toUpperCase()
  }), /*#__PURE__*/React.createElement(ConfigRow, {
    k: target.type === "http" ? "URL" : target.type === "tcp" ? "Host : port" : "Domain · record",
    v: target.endpoint,
    mono: true
  }), /*#__PURE__*/React.createElement(ConfigRow, {
    k: "Interval",
    v: target.interval,
    mono: true
  }), /*#__PURE__*/React.createElement(ConfigRow, {
    k: "Timeout",
    v: target.timeout,
    mono: true
  }), /*#__PURE__*/React.createElement(ConfigRow, {
    k: "Tags",
    v: target.tags.join(", ")
  }), target.desc && /*#__PURE__*/React.createElement(ConfigRow, {
    k: "Description",
    v: target.desc
  }), /*#__PURE__*/React.createElement("div", {
    style: {
      display: "grid",
      gridTemplateColumns: "140px 1fr",
      gap: 8,
      padding: "6px 0",
      fontSize: "var(--text-md)"
    }
  }, /*#__PURE__*/React.createElement("span", {
    style: {
      color: "var(--text-3)"
    }
  }, "Enabled"), /*#__PURE__*/React.createElement("span", {
    style: {
      color: target.paused ? "var(--text-3)" : "var(--status-up)"
    }
  }, target.paused ? "No — paused" : "Yes")))));
}
function DetailTitle({
  target
}) {
  return /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement(StatusDot, {
    status: target.status,
    pulse: target.status === "down",
    size: "lg"
  }), /*#__PURE__*/React.createElement("span", {
    style: {
      overflow: "hidden",
      textOverflow: "ellipsis",
      whiteSpace: "nowrap"
    }
  }, target.name), /*#__PURE__*/React.createElement(Icon, {
    name: detailTypeIcon[target.type],
    size: 14,
    style: {
      color: "var(--text-3)"
    }
  }));
}
Object.assign(window, {
  TargetDetailBody,
  DetailTitle,
  ResponseChart
});
})(); } catch (e) { __ds_ns.__errors.push({ path: "ui_kits/dashboard/detail.jsx", error: String((e && e.message) || e) }); }

// ui_kits/dashboard/forms.jsx
try { (() => {
// Add/Edit Target form + Notifier form (slide-over bodies).
const FDS = window.HealthMonitorRedesign_502b6f;
const {
  Icon: FIcon,
  TextField,
  SelectField,
  Switch: FSwitch,
  Segmented: FSeg,
  TagInput,
  DurationField,
  Button: FBtn
} = FDS;
function FormSection({
  title,
  children
}) {
  return /*#__PURE__*/React.createElement("div", {
    className: "hm-formsection"
  }, /*#__PURE__*/React.createElement("span", {
    className: "overline"
  }, title), children);
}

/* ---------------- Target form ---------------- */
function TargetForm({
  target
}) {
  const t = target || {};
  const [type, setType] = React.useState(t.type || "http");
  const [url, setUrl] = React.useState(t.type === "http" ? t.endpoint || "" : "https://");
  const [urlErr, setUrlErr] = React.useState("");
  const [tags, setTags] = React.useState(t.tags || []);
  const [interval, setInterval_] = React.useState(t.interval || "1m");
  const [timeout_, setTimeout_] = React.useState(t.timeout || "5s");
  const [enabled, setEnabled] = React.useState(t.paused ? false : true);
  const checkUrl = v => setUrlErr(/^https?:\/\/.+/.test(v) ? "" : "Must start with http:// or https://");
  return /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement(FormSection, {
    title: "Identity"
  }, /*#__PURE__*/React.createElement(TextField, {
    label: "Name",
    defaultValue: t.name,
    placeholder: "API /health",
    hint: "Shown on the dashboard and in alerts"
  }), /*#__PURE__*/React.createElement(TextField, {
    label: "Description",
    defaultValue: t.desc,
    placeholder: "Optional note for your future self"
  }), /*#__PURE__*/React.createElement("div", {
    className: "hm-field"
  }, /*#__PURE__*/React.createElement("span", {
    className: "hm-label"
  }, "Tags"), /*#__PURE__*/React.createElement(TagInput, {
    value: tags,
    onChange: setTags
  }), /*#__PURE__*/React.createElement("span", {
    className: "hm-hint"
  }, "Used for filtering; comma or Enter to add"))), /*#__PURE__*/React.createElement(FormSection, {
    title: "Check"
  }, /*#__PURE__*/React.createElement("div", {
    className: "hm-field"
  }, /*#__PURE__*/React.createElement("span", {
    className: "hm-label"
  }, "Type"), /*#__PURE__*/React.createElement(FSeg, {
    ariaLabel: "Target type",
    options: [{
      value: "http",
      label: "HTTP",
      icon: /*#__PURE__*/React.createElement(FIcon, {
        name: "globe",
        size: 13
      })
    }, {
      value: "tcp",
      label: "TCP",
      icon: /*#__PURE__*/React.createElement(FIcon, {
        name: "server",
        size: 13
      })
    }, {
      value: "dns",
      label: "DNS",
      icon: /*#__PURE__*/React.createElement(FIcon, {
        name: "at-sign",
        size: 13
      })
    }],
    value: type,
    onChange: setType
  })), type === "http" && /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement(TextField, {
    label: "URL",
    mono: true,
    value: url,
    error: urlErr,
    onChange: e => {
      setUrl(e.target.value);
      checkUrl(e.target.value);
    }
  }), /*#__PURE__*/React.createElement("div", {
    style: {
      display: "grid",
      gridTemplateColumns: "1fr 1fr",
      gap: 12
    }
  }, /*#__PURE__*/React.createElement(SelectField, {
    label: "Method",
    options: ["GET", "HEAD", "POST"]
  }), /*#__PURE__*/React.createElement(TextField, {
    label: "Expected status",
    mono: true,
    defaultValue: "200",
    hint: "Comma-separate multiple codes"
  }))), type === "tcp" && /*#__PURE__*/React.createElement(TextField, {
    label: "Host : port",
    mono: true,
    defaultValue: t.type === "tcp" ? t.endpoint : "",
    placeholder: "10.0.4.2:5432"
  }), type === "dns" && /*#__PURE__*/React.createElement("div", {
    style: {
      display: "grid",
      gridTemplateColumns: "2fr 1fr",
      gap: 12
    }
  }, /*#__PURE__*/React.createElement(TextField, {
    label: "Domain",
    mono: true,
    placeholder: "example.dev"
  }), /*#__PURE__*/React.createElement(SelectField, {
    label: "Record",
    options: ["A", "AAAA", "CNAME", "MX", "TXT"]
  })), /*#__PURE__*/React.createElement(DurationField, {
    label: "Check interval",
    value: interval,
    onChange: setInterval_
  }), /*#__PURE__*/React.createElement(DurationField, {
    label: "Timeout",
    value: timeout_,
    onChange: setTimeout_,
    presets: ["3s", "5s", "10s", "30s"]
  })), /*#__PURE__*/React.createElement(FormSection, {
    title: "Advanced"
  }, /*#__PURE__*/React.createElement(FSwitch, {
    checked: enabled,
    onChange: setEnabled,
    label: "Enabled \u2014 run checks on schedule"
  })));
}

/* ---------------- Notifier form ---------------- */
const NOTIFIER_TYPES = [{
  value: "telegram",
  label: "Telegram",
  icon: "send",
  hint: "Bot message to a chat"
}, {
  value: "email",
  label: "Email (SMTP)",
  icon: "mail",
  hint: "Any SMTP server"
}, {
  value: "gmail",
  label: "Gmail",
  icon: "mail",
  hint: "App password"
}, {
  value: "gmail_oauth",
  label: "Gmail OAuth",
  icon: "mail",
  hint: "Google sign-in"
}, {
  value: "webhook",
  label: "Webhook",
  icon: "link",
  hint: "POST JSON anywhere"
}];
function TypeCard({
  t,
  selected,
  onSelect
}) {
  return /*#__PURE__*/React.createElement("button", {
    type: "button",
    onClick: onSelect,
    "aria-pressed": selected,
    style: {
      display: "flex",
      flexDirection: "column",
      gap: 4,
      alignItems: "flex-start",
      cursor: "pointer",
      padding: "10px 12px",
      textAlign: "left",
      minHeight: 64,
      background: selected ? "var(--accent-subtle)" : "var(--bg-inset)",
      border: `1px solid ${selected ? "var(--accent-border)" : "var(--border-1)"}`,
      borderRadius: "var(--radius-md)",
      color: "var(--text-1)",
      font: "inherit"
    }
  }, /*#__PURE__*/React.createElement("span", {
    style: {
      display: "inline-flex",
      alignItems: "center",
      gap: 6,
      fontWeight: 600,
      fontSize: "var(--text-md)",
      color: selected ? "var(--accent)" : "var(--text-1)"
    }
  }, /*#__PURE__*/React.createElement(FIcon, {
    name: t.icon,
    size: 14
  }), t.label), /*#__PURE__*/React.createElement("span", {
    style: {
      fontSize: "var(--text-sm)",
      color: "var(--text-3)"
    }
  }, t.hint));
}
function NotifierForm({
  notifier
}) {
  const n = notifier || {};
  const [type, setType] = React.useState(n.type || "telegram");
  const [enabled, setEnabled] = React.useState(n.enabled != null ? n.enabled : true);
  return /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement(FormSection, {
    title: "Channel"
  }, /*#__PURE__*/React.createElement("div", {
    style: {
      display: "grid",
      gridTemplateColumns: "1fr 1fr",
      gap: 8
    }
  }, NOTIFIER_TYPES.map(t => /*#__PURE__*/React.createElement(TypeCard, {
    key: t.value,
    t: t,
    selected: type === t.value,
    onSelect: () => setType(t.value)
  })))), /*#__PURE__*/React.createElement(FormSection, {
    title: "Settings"
  }, /*#__PURE__*/React.createElement(TextField, {
    label: "Name",
    defaultValue: n.name,
    placeholder: "Ops channel"
  }), type === "telegram" && /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement(TextField, {
    label: "Bot token",
    mono: true,
    type: "password",
    defaultValue: n.id ? "••••••••••••" : ""
  }), /*#__PURE__*/React.createElement(TextField, {
    label: "Chat ID",
    mono: true,
    placeholder: "-1001234567890"
  })), type === "email" && /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("div", {
    style: {
      display: "grid",
      gridTemplateColumns: "2fr 1fr",
      gap: 12
    }
  }, /*#__PURE__*/React.createElement(TextField, {
    label: "SMTP host",
    mono: true,
    placeholder: "smtp.example.dev"
  }), /*#__PURE__*/React.createElement(TextField, {
    label: "Port",
    mono: true,
    defaultValue: "587"
  })), /*#__PURE__*/React.createElement(TextField, {
    label: "From",
    mono: true,
    placeholder: "alerts@example.dev"
  }), /*#__PURE__*/React.createElement(TextField, {
    label: "To",
    mono: true,
    placeholder: "oncall@example.dev",
    hint: "Comma-separate multiple recipients"
  })), type === "gmail" && /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement(TextField, {
    label: "Gmail address",
    mono: true,
    placeholder: "alerts@gmail.com"
  }), /*#__PURE__*/React.createElement(TextField, {
    label: "App password",
    mono: true,
    type: "password"
  }), /*#__PURE__*/React.createElement(TextField, {
    label: "To",
    mono: true,
    placeholder: "oncall@example.dev"
  })), type === "gmail_oauth" && /*#__PURE__*/React.createElement("div", {
    className: "hm-field"
  }, /*#__PURE__*/React.createElement("span", {
    className: "hm-label"
  }, "Google account"), /*#__PURE__*/React.createElement(FBtn, {
    icon: /*#__PURE__*/React.createElement(FIcon, {
      name: "external",
      size: 14
    })
  }, "Sign in with Google"), /*#__PURE__*/React.createElement("span", {
    className: "hm-hint"
  }, "Opens Google consent; token is stored locally")), type === "webhook" && /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement(TextField, {
    label: "URL",
    mono: true,
    placeholder: "https://hooks.example.dev/uptime"
  }), /*#__PURE__*/React.createElement(SelectField, {
    label: "Method",
    options: ["POST", "PUT"]
  }))), /*#__PURE__*/React.createElement(FormSection, {
    title: "Advanced"
  }, /*#__PURE__*/React.createElement(FSwitch, {
    checked: enabled,
    onChange: setEnabled,
    label: "Enabled \u2014 send alerts through this channel"
  })));
}
Object.assign(window, {
  TargetForm,
  NotifierForm
});
})(); } catch (e) { __ds_ns.__errors.push({ path: "ui_kits/dashboard/forms.jsx", error: String((e && e.message) || e) }); }

__ds_ns.Icon = __ds_scope.Icon;

__ds_ns.IncidentItem = __ds_scope.IncidentItem;

__ds_ns.NotifierRow = __ds_scope.NotifierRow;

__ds_ns.StatCard = __ds_scope.StatCard;

__ds_ns.TargetCard = __ds_scope.TargetCard;

__ds_ns.TargetTable = __ds_scope.TargetTable;

__ds_ns.ConfirmDialog = __ds_scope.ConfirmDialog;

__ds_ns.EmptyState = __ds_scope.EmptyState;

__ds_ns.Skeleton = __ds_scope.Skeleton;

__ds_ns.SlideOver = __ds_scope.SlideOver;

__ds_ns.Toast = __ds_scope.Toast;

__ds_ns.Button = __ds_scope.Button;

__ds_ns.DurationField = __ds_scope.DurationField;

__ds_ns.IconButton = __ds_scope.IconButton;

__ds_ns.Segmented = __ds_scope.Segmented;

__ds_ns.SelectField = __ds_scope.SelectField;

__ds_ns.Switch = __ds_scope.Switch;

__ds_ns.TagInput = __ds_scope.TagInput;

__ds_ns.TextField = __ds_scope.TextField;

__ds_ns.LiveIndicator = __ds_scope.LiveIndicator;

__ds_ns.Sparkline = __ds_scope.Sparkline;

__ds_ns.StatusBadge = __ds_scope.StatusBadge;

__ds_ns.StatusDot = __ds_scope.StatusDot;

__ds_ns.TickBar = __ds_scope.TickBar;

})();
