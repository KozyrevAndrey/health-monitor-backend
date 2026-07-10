// Dashboard app shell. Reads flags: HM_SCENARIO ("ok"|"alarm"), HM_TAB,
// HM_OPEN ("detail"|"target-form"|"notifier-form"), HM_VIEW ("cards"|"table").
const ADS = window.HealthMonitorRedesign_502b6f;
const {
  Icon: AIcon, StatusDot: ADot, StatusBadge: ABadge, LiveIndicator, StatCard: AStat,
  TargetCard: ACard, TargetTable: ATable, IncidentItem: AInc, NotifierRow: ANotif,
  Toast: AToast, ConfirmDialog: AConfirm, SlideOver: APanel, EmptyState: AEmpty,
  Button: ABtn, Segmented: ASeg, SelectField: ASelect,
} = ADS;

function Hero({ data }) {
  const downs = data.targets.filter((t) => t.status === "down");
  const degraded = data.targets.filter((t) => t.status === "degraded");
  const ok = downs.length === 0;
  return (
    <div className="hm-hero" data-status={ok ? "up" : "down"} role="status" aria-live="polite" style={{ marginTop: "var(--space-4)" }}>
      <span className="hm-hero-icon">
        <AIcon name={ok ? "check" : "alert-triangle"} size={22} />
      </span>
      <div style={{ flex: 1, minWidth: 0 }}>
        <h1>{ok ? "All systems operational" : `${downs.length} target${downs.length > 1 ? "s" : ""} down`}</h1>
        <p>
          {ok
            ? `${data.targets.filter((t) => !t.paused).length} targets · all checks passing`
            : downs.map((t) => t.name).join(", ") + (degraded.length ? ` · ${degraded.length} degraded` : "")}
        </p>
      </div>
      {!ok && (
        <div style={{ display: "flex", gap: 8, flexWrap: "wrap" }}>
          {downs.map((t) => <ABadge key={t.id} status="down" label={t.name} />)}
        </div>
      )}
    </div>
  );
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
  const [targetForm, setTargetForm] = React.useState(window.HM_OPEN === "target-form" ? { mode: "add" } : null);
  const [notifierForm, setNotifierForm] = React.useState(window.HM_OPEN === "notifier-form" ? { mode: "add" } : null);
  const [confirm, setConfirm] = React.useState(null);
  const [toasts, setToasts] = React.useState([]);

  const toast = (kind, title, message) => {
    const id = Date.now() + Math.random();
    setToasts((ts) => [...ts, { id, kind, title, message }]);
    setTimeout(() => setToasts((ts) => ts.filter((t) => t.id !== id)), 5200);
  };

  const setThemeAttr = (next) => {
    setTheme(next);
    document.documentElement.setAttribute("data-theme", next);
  };

  const counts = {
    all: data.targets.length,
    down: data.targets.filter((t) => t.status === "down").length,
    degraded: data.targets.filter((t) => t.status === "degraded").length,
    paused: data.targets.filter((t) => t.paused).length,
  };

  const order = { down: 0, degraded: 1, up: 2, unknown: 3 };
  const shown = data.targets
    .filter((t) => {
      if (filter === "down" && t.status !== "down") return false;
      if (filter === "degraded" && t.status !== "degraded") return false;
      if (filter === "paused" && !t.paused) return false;
      const q = query.trim().toLowerCase();
      return !q || t.name.toLowerCase().includes(q) || t.endpoint.toLowerCase().includes(q) || t.tags.some((x) => x.includes(q));
    })
    .sort((a, b) =>
      sort === "name" ? a.name.localeCompare(b.name)
      : sort === "uptime" ? Number(a.uptime) - Number(b.uptime)
      : sort === "response" ? (b.lastMs || 0) - (a.lastMs || 0)
      : order[a.status] - order[b.status] || a.name.localeCompare(b.name));

  const askDelete = (t) => setConfirm({
    title: "Delete target?",
    body: `“${t.name}” and its check history will be permanently removed. Notifiers are not affected.`,
    onConfirm: () => { setConfirm(null); toast("info", "Target deleted", `${t.name} is no longer monitored`); },
  });

  const avgMs = Math.round(data.targets.filter((t) => t.lastMs).reduce((a, t) => a + t.lastMs, 0) / data.targets.filter((t) => t.lastMs).length);

  return (
    <>
      <header className="hm-topbar">
        <span className="hm-brand">
          <span className="hm-brand-mark"><AIcon name="activity" size={15} /></span>
          Health Monitor
        </span>
        <ASeg
          ariaLabel="Section"
          options={[{ value: "dashboard", label: "Dashboard" }, { value: "notifiers", label: "Notifiers", count: data.notifiers.length }]}
          value={tab} onChange={setTab}
        />
        <span style={{ flex: 1 }}></span>
        <LiveIndicator state="live" />
        <button type="button" className="hm-iconbtn" aria-label={`Switch to ${theme === "dark" ? "light" : "dark"} theme`}
          onClick={() => setThemeAttr(theme === "dark" ? "light" : "dark")}>
          <AIcon name={theme === "dark" ? "sun" : "moon"} size={15} />
        </button>
        <ABtn variant="primary" icon={<AIcon name="plus" size={15} />} onClick={() => setTargetForm({ mode: "add" })}>
          <span className="hm-topbar-add-label">Add target</span>
        </ABtn>
      </header>

      <main className="hm-shell">
        {tab === "dashboard" && (
          <>
            <Hero data={data} />

            <div className="hm-statrow">
              <AStat label="Targets" value={counts.all} sub={`${counts.all - counts.paused} enabled`} />
              <AStat label="Up" value={counts.all - counts.down - counts.degraded - counts.paused} tone="success" />
              <AStat label="Down" value={counts.down} tone={counts.down ? "danger" : "default"} sub={counts.down ? "see incidents" : " "} />
              <AStat label="Avg response" value={avgMs} unit="ms" sub="across all targets" />
            </div>

            {data.incidents.ongoing.length > 0 && (
              <>
                <div className="hm-sectionhead">
                  <h2>Ongoing incidents</h2>
                  <span className="hm-badge" data-status="down"><AIcon name="alert-triangle" size={12} />{data.incidents.ongoing.length} active</span>
                </div>
                <div className="hm-card hm-card--down">
                  {data.incidents.ongoing.map((inc, i) => (
                    <AInc key={inc.id} incident={inc} last={i === data.incidents.ongoing.length - 1}
                      onOpen={() => setDetail(data.targets.find((t) => t.id === inc.targetId))} />
                  ))}
                </div>
              </>
            )}

            <div className="hm-toolbar">
              <div className="hm-search">
                <AIcon name="search" size={14} />
                <input className="hm-input" placeholder="Search targets, tags…" aria-label="Search targets"
                  value={query} onChange={(e) => setQuery(e.target.value)} />
              </div>
              <ASeg
                ariaLabel="Status filter"
                options={[
                  { value: "all", label: "All", count: counts.all },
                  { value: "down", label: "Down", count: counts.down },
                  { value: "degraded", label: "Degraded", count: counts.degraded },
                  { value: "paused", label: "Paused", count: counts.paused },
                ]}
                value={filter} onChange={setFilter}
              />
              <span style={{ flex: 1 }}></span>
              <select className="hm-select" style={{ width: "auto", minHeight: 30 }} aria-label="Sort targets" value={sort} onChange={(e) => setSort(e.target.value)}>
                <option value="status">Sort: status</option>
                <option value="name">Sort: name</option>
                <option value="uptime">Sort: uptime</option>
                <option value="response">Sort: response time</option>
              </select>
              <ASeg
                ariaLabel="View"
                options={[
                  { value: "cards", label: "", icon: <AIcon name="grid" size={13} aria-label="Card view" /> },
                  { value: "table", label: "", icon: <AIcon name="list" size={13} aria-label="Table view" /> },
                ]}
                value={view} onChange={setView}
              />
            </div>

            {shown.length === 0 ? (
              <div className="hm-card">
                <AEmpty icon="search" title="No matching targets"
                  body={query ? `Nothing matches “${query}”. Try a different search or clear the filter.` : "Nothing matches this filter."}
                  action={<ABtn onClick={() => { setQuery(""); setFilter("all"); }}>Clear filters</ABtn>} />
              </div>
            ) : view === "cards" ? (
              <div className="hm-grid">
                {shown.map((t) => (
                  <ACard key={t.id} target={t}
                    onOpen={() => setDetail(t)}
                    onEdit={() => setTargetForm({ mode: "edit", target: t })}
                    onDelete={() => askDelete(t)} />
                ))}
              </div>
            ) : (
              <div className="hm-card hm-tablewrap scroll-thin">
                <ATable targets={shown}
                  onOpen={(t) => setDetail(t)}
                  onEdit={(t) => setTargetForm({ mode: "edit", target: t })}
                  onDelete={(t) => askDelete(t)} />
              </div>
            )}

            <div className="hm-sectionhead">
              <h2>Incident history</h2>
              <span style={{ flex: 1 }}></span>
              <ABtn variant="ghost" icon={<AIcon name={showResolved ? "chevron-down" : "chevron-right"} size={14} />}
                onClick={() => setShowResolved(!showResolved)}>
                {data.incidents.resolved.length} resolved
              </ABtn>
            </div>
            {showResolved && (
              <div className="hm-card">
                {data.incidents.resolved.map((inc, i) => (
                  <AInc key={inc.id} incident={inc} last={i === data.incidents.resolved.length - 1}
                    onOpen={() => setDetail(data.targets.find((t) => t.id === inc.targetId))} />
                ))}
              </div>
            )}

            <div className="hm-foot">
              <LiveIndicator state="live" />
              <span>Last updated <span className="num">12s</span> ago</span>
            </div>
          </>
        )}

        {tab === "notifiers" && (
          <>
            <div className="hm-sectionhead" style={{ marginTop: "var(--space-6)" }}>
              <h2>Notifiers</h2>
              <span style={{ flex: 1 }}></span>
              <ABtn variant="primary" icon={<AIcon name="plus" size={15} />} onClick={() => setNotifierForm({ mode: "add" })}>Add notifier</ABtn>
            </div>
            <p style={{ margin: "0 0 var(--space-3)", color: "var(--text-2)", fontSize: "var(--text-md)", maxWidth: 560 }}>
              Alert channels for down / recovered / slow-response events. Rarely touched — they live here, out of the way of live status.
            </p>
            <div className="hm-card">
              {data.notifiers.map((n, i) => (
                <div key={n.id} style={{ borderTop: i ? "1px solid var(--border-1)" : "none" }}>
                  <ANotif notifier={n}
                    onToggle={(on) => toast(on ? "success" : "info", on ? "Notifier enabled" : "Notifier disabled", n.name)}
                    onEdit={() => setNotifierForm({ mode: "edit", notifier: n })}
                    onDelete={() => setConfirm({
                      title: "Delete notifier?",
                      body: `“${n.name}” will stop receiving alerts.`,
                      onConfirm: () => { setConfirm(null); toast("info", "Notifier deleted", n.name); },
                    })} />
                </div>
              ))}
            </div>
          </>
        )}
      </main>

      {detail && (
        <APanel width={720} title={<window.DetailTitle target={detail} />}
          titleExtra={<ABadge status={detail.status} label={detail.paused ? "Paused" : undefined} />}
          onClose={() => setDetail(null)}>
          <window.TargetDetailBody target={detail} incidents={data.incidents}
            onEdit={() => { setDetail(null); setTargetForm({ mode: "edit", target: detail }); }} />
        </APanel>
      )}

      {targetForm && (
        <APanel width={520}
          title={targetForm.mode === "add" ? "Add target" : `Edit ${targetForm.target.name}`}
          onClose={() => setTargetForm(null)}
          footer={<>
            <ABtn onClick={() => setTargetForm(null)}>Cancel</ABtn>
            <ABtn variant="primary" onClick={() => { setTargetForm(null); toast("success", targetForm.mode === "add" ? "Target added" : "Target saved", "Checks start on the next tick"); }}>
              {targetForm.mode === "add" ? "Add target" : "Save changes"}
            </ABtn>
          </>}>
          <window.TargetForm target={targetForm.target} />
        </APanel>
      )}

      {notifierForm && (
        <APanel width={520}
          title={notifierForm.mode === "add" ? "Add notifier" : `Edit ${notifierForm.notifier.name}`}
          onClose={() => setNotifierForm(null)}
          footer={<>
            <ABtn onClick={() => setNotifierForm(null)}>Cancel</ABtn>
            <ABtn variant="primary" onClick={() => { setNotifierForm(null); toast("success", "Notifier saved", "A test message was sent"); }}>Save notifier</ABtn>
          </>}>
          <window.NotifierForm notifier={notifierForm.notifier} />
        </APanel>
      )}

      {confirm && <AConfirm title={confirm.title} body={confirm.body} onConfirm={confirm.onConfirm} onCancel={() => setConfirm(null)} />}

      <div className="hm-toaststack" aria-live="polite">
        {toasts.map((t) => (
          <AToast key={t.id} kind={t.kind} title={t.title} message={t.message}
            onDismiss={() => setToasts((ts) => ts.filter((x) => x.id !== t.id))} />
        ))}
      </div>
    </>
  );
}

window.HMApp = App;
