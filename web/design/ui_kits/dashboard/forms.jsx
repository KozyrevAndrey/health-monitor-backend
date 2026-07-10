// Add/Edit Target form + Notifier form (slide-over bodies).
const FDS = window.HealthMonitorRedesign_502b6f;
const { Icon: FIcon, TextField, SelectField, Switch: FSwitch, Segmented: FSeg, TagInput, DurationField, Button: FBtn } = FDS;

function FormSection({ title, children }) {
  return (
    <div className="hm-formsection">
      <span className="overline">{title}</span>
      {children}
    </div>
  );
}

/* ---------------- Target form ---------------- */
function TargetForm({ target }) {
  const t = target || {};
  const [type, setType] = React.useState(t.type || "http");
  const [url, setUrl] = React.useState(t.type === "http" ? t.endpoint || "" : "https://");
  const [urlErr, setUrlErr] = React.useState("");
  const [tags, setTags] = React.useState(t.tags || []);
  const [interval, setInterval_] = React.useState(t.interval || "1m");
  const [timeout_, setTimeout_] = React.useState(t.timeout || "5s");
  const [enabled, setEnabled] = React.useState(t.paused ? false : true);

  const checkUrl = (v) => setUrlErr(/^https?:\/\/.+/.test(v) ? "" : "Must start with http:// or https://");

  return (
    <div>
      <FormSection title="Identity">
        <TextField label="Name" defaultValue={t.name} placeholder="API /health" hint="Shown on the dashboard and in alerts" />
        <TextField label="Description" defaultValue={t.desc} placeholder="Optional note for your future self" />
        <div className="hm-field">
          <span className="hm-label">Tags</span>
          <TagInput value={tags} onChange={setTags} />
          <span className="hm-hint">Used for filtering; comma or Enter to add</span>
        </div>
      </FormSection>

      <FormSection title="Check">
        <div className="hm-field">
          <span className="hm-label">Type</span>
          <FSeg
            ariaLabel="Target type"
            options={[
              { value: "http", label: "HTTP", icon: <FIcon name="globe" size={13} /> },
              { value: "tcp", label: "TCP", icon: <FIcon name="server" size={13} /> },
              { value: "dns", label: "DNS", icon: <FIcon name="at-sign" size={13} /> },
            ]}
            value={type} onChange={setType}
          />
        </div>
        {type === "http" && (
          <>
            <TextField label="URL" mono value={url} error={urlErr}
              onChange={(e) => { setUrl(e.target.value); checkUrl(e.target.value); }} />
            <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 12 }}>
              <SelectField label="Method" options={["GET", "HEAD", "POST"]} />
              <TextField label="Expected status" mono defaultValue="200" hint="Comma-separate multiple codes" />
            </div>
          </>
        )}
        {type === "tcp" && (
          <TextField label="Host : port" mono defaultValue={t.type === "tcp" ? t.endpoint : ""} placeholder="10.0.4.2:5432" />
        )}
        {type === "dns" && (
          <div style={{ display: "grid", gridTemplateColumns: "2fr 1fr", gap: 12 }}>
            <TextField label="Domain" mono placeholder="example.dev" />
            <SelectField label="Record" options={["A", "AAAA", "CNAME", "MX", "TXT"]} />
          </div>
        )}
        <DurationField label="Check interval" value={interval} onChange={setInterval_} />
        <DurationField label="Timeout" value={timeout_} onChange={setTimeout_} presets={["3s", "5s", "10s", "30s"]} />
      </FormSection>

      <FormSection title="Advanced">
        <FSwitch checked={enabled} onChange={setEnabled} label="Enabled — run checks on schedule" />
      </FormSection>
    </div>
  );
}

/* ---------------- Notifier form ---------------- */
const NOTIFIER_TYPES = [
  { value: "telegram", label: "Telegram", icon: "send", hint: "Bot message to a chat" },
  { value: "email", label: "Email (SMTP)", icon: "mail", hint: "Any SMTP server" },
  { value: "gmail", label: "Gmail", icon: "mail", hint: "App password" },
  { value: "gmail_oauth", label: "Gmail OAuth", icon: "mail", hint: "Google sign-in" },
  { value: "webhook", label: "Webhook", icon: "link", hint: "POST JSON anywhere" },
];

function TypeCard({ t, selected, onSelect }) {
  return (
    <button type="button" onClick={onSelect} aria-pressed={selected}
      style={{
        display: "flex", flexDirection: "column", gap: 4, alignItems: "flex-start", cursor: "pointer",
        padding: "10px 12px", textAlign: "left", minHeight: 64,
        background: selected ? "var(--accent-subtle)" : "var(--bg-inset)",
        border: `1px solid ${selected ? "var(--accent-border)" : "var(--border-1)"}`,
        borderRadius: "var(--radius-md)", color: "var(--text-1)", font: "inherit",
      }}>
      <span style={{ display: "inline-flex", alignItems: "center", gap: 6, fontWeight: 600, fontSize: "var(--text-md)", color: selected ? "var(--accent)" : "var(--text-1)" }}>
        <FIcon name={t.icon} size={14} />{t.label}
      </span>
      <span style={{ fontSize: "var(--text-sm)", color: "var(--text-3)" }}>{t.hint}</span>
    </button>
  );
}

function NotifierForm({ notifier }) {
  const n = notifier || {};
  const [type, setType] = React.useState(n.type || "telegram");
  const [enabled, setEnabled] = React.useState(n.enabled != null ? n.enabled : true);
  return (
    <div>
      <FormSection title="Channel">
        <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 8 }}>
          {NOTIFIER_TYPES.map((t) => (
            <TypeCard key={t.value} t={t} selected={type === t.value} onSelect={() => setType(t.value)} />
          ))}
        </div>
      </FormSection>
      <FormSection title="Settings">
        <TextField label="Name" defaultValue={n.name} placeholder="Ops channel" />
        {type === "telegram" && (
          <>
            <TextField label="Bot token" mono type="password" defaultValue={n.id ? "••••••••••••" : ""} />
            <TextField label="Chat ID" mono placeholder="-1001234567890" />
          </>
        )}
        {(type === "email") && (
          <>
            <div style={{ display: "grid", gridTemplateColumns: "2fr 1fr", gap: 12 }}>
              <TextField label="SMTP host" mono placeholder="smtp.example.dev" />
              <TextField label="Port" mono defaultValue="587" />
            </div>
            <TextField label="From" mono placeholder="alerts@example.dev" />
            <TextField label="To" mono placeholder="oncall@example.dev" hint="Comma-separate multiple recipients" />
          </>
        )}
        {(type === "gmail") && (
          <>
            <TextField label="Gmail address" mono placeholder="alerts@gmail.com" />
            <TextField label="App password" mono type="password" />
            <TextField label="To" mono placeholder="oncall@example.dev" />
          </>
        )}
        {type === "gmail_oauth" && (
          <div className="hm-field">
            <span className="hm-label">Google account</span>
            <FBtn icon={<FIcon name="external" size={14} />}>Sign in with Google</FBtn>
            <span className="hm-hint">Opens Google consent; token is stored locally</span>
          </div>
        )}
        {type === "webhook" && (
          <>
            <TextField label="URL" mono placeholder="https://hooks.example.dev/uptime" />
            <SelectField label="Method" options={["POST", "PUT"]} />
          </>
        )}
      </FormSection>
      <FormSection title="Advanced">
        <FSwitch checked={enabled} onChange={setEnabled} label="Enabled — send alerts through this channel" />
      </FormSection>
    </div>
  );
}

Object.assign(window, { TargetForm, NotifierForm });
