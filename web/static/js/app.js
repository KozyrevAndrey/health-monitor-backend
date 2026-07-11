/* Health Monitor dashboard — vanilla JS implementation of the web/design system.
   No build step: this file + tokens.css/app.css are the whole frontend. */
'use strict';

/* ═══════════════ Icons (inline SVG, Lucide-style — see web/design/components/core/Icon.jsx) ═══════════════ */

const ICON_PATHS = {
    check: ['M20 6 9 17l-5-5'],
    x: ['M18 6 6 18', 'M6 6l12 12'],
    plus: ['M12 5v14', 'M5 12h14'],
    search: ['M21 21l-4.35-4.35', 'M11 4a7 7 0 1 1 0 14 7 7 0 0 1 0-14z'],
    'chevron-down': ['m6 9 6 6 6-6'],
    'chevron-right': ['m9 6 6 6-6 6'],
    'arrow-left': ['M19 12H5', 'm12 19-7-7 7-7'],
    clock: ['M12 3a9 9 0 1 1 0 18 9 9 0 0 1 0-18z', 'M12 7v5l3 2'],
    trash: ['M4 7h16', 'M9 7V4h6v3', 'M6 7l1 13h10l1-13'],
    pencil: ['M17 3l4 4L8 20l-5 1 1-5L17 3'],
    'alert-triangle': ['M12 3 2 20h20L12 3z', 'M12 10v4', 'M12 17h.01'],
    'check-circle': ['M12 3a9 9 0 1 1 0 18 9 9 0 0 1 0-18z', 'm8.5 12.5 2.5 2.5 4.5-5'],
    'x-circle': ['M12 3a9 9 0 1 1 0 18 9 9 0 0 1 0-18z', 'M15 9l-6 6', 'M9 9l6 6'],
    info: ['M12 3a9 9 0 1 1 0 18 9 9 0 0 1 0-18z', 'M12 16v-5', 'M12 8h.01'],
    activity: ['M22 12h-4l-3 9L9 3l-3 9H2'],
    globe: ['M12 3a9 9 0 1 1 0 18 9 9 0 0 1 0-18z', 'M3 12h18', 'M12 3c2.5 2.5 3.8 5.6 3.8 9S14.5 18.5 12 21c-2.5-2.5-3.8-5.6-3.8-9S9.5 5.5 12 3z'],
    server: ['M4 4h16a2 2 0 0 1 2 2v3a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2z', 'M4 13h16a2 2 0 0 1 2 2v3a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2v-3a2 2 0 0 1 2-2z', 'M6 7.5h.01', 'M6 16.5h.01'],
    'at-sign': ['M12 8a4 4 0 1 0 0 8 4 4 0 0 0 0-8z', 'M16 8v5a3 3 0 0 0 6 0v-1a10 10 0 1 0-3.9 7.9'],
    bell: ['M6 9a6 6 0 0 1 12 0c0 5 2 6 2 6H4s2-1 2-6', 'M10 19a2 2 0 0 0 4 0'],
    zap: ['M13 2 3 14h7l-1 8 11-12h-7l1-8z'],
    pause: ['M9 5v14', 'M15 5v14'],
    send: ['M22 2 11 13', 'M22 2 15 22l-4-9-9-4z'],
    mail: ['M4 5h16a1 1 0 0 1 1 1v12a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1V6a1 1 0 0 1 1-1z', 'm3 7 9 6 9-6'],
    link: ['M10 14a5 5 0 0 1 0-7l1.5-1.5a5 5 0 0 1 7 7L17.5 13.5', 'M14 10a5 5 0 0 1 0 7l-1.5 1.5a5 5 0 0 1-7-7L6.5 10.5'],
    moon: ['M21 13A9 9 0 1 1 11 3a7 7 0 0 0 10 10z'],
    sun: ['M12 8a4 4 0 1 1 0 8 4 4 0 0 1 0-8z', 'M12 2v2', 'M12 20v2', 'M2 12h2', 'M20 12h2', 'm5 5 1.4 1.4', 'm17.6 17.6 1.4 1.4', 'm19 5-1.4 1.4', 'm6.4 17.6-1.4 1.4'],
    grid: ['M4 4h6v6H4z', 'M14 4h6v6h-6z', 'M4 14h6v6H4z', 'M14 14h6v6h-6z'],
    list: ['M8 6h13', 'M8 12h13', 'M8 18h13', 'M3 6h.01', 'M3 12h.01', 'M3 18h.01'],
    external: ['M14 3h7v7', 'M21 3l-9 9', 'M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h6'],
    refresh: ['M21 12a9 9 0 1 1-2.9-6.6L21 8', 'M21 3v5h-5'],
    logout: ['M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4', 'm16 17 5-5-5-5', 'M21 12H9'],
};

function icon(name, size = 16, extra = '') {
    const paths = ICON_PATHS[name] || ICON_PATHS.info;
    return `<svg width="${size}" height="${size}" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true" style="flex:none;${extra}">${paths.map(d => `<path d="${d}"/>`).join('')}</svg>`;
}

/* ═══════════════ Utils ═══════════════ */

const esc = (s) => String(s ?? '').replace(/[&<>"']/g, (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c]));

const API_STATUS = { success: 'up', failure: 'down', warning: 'degraded' };
const toUiStatus = (s) => API_STATUS[s] || 'unknown';
const STATUS_ORDER = { down: 0, degraded: 1, up: 2, unknown: 3 };
const TYPE_ICON = { http: 'globe', tcp: 'server', dns: 'at-sign' };

function fmtDuration(ns) {
    const s = Math.round(ns / 1e9);
    if (s < 60) return s + 's';
    if (s < 3600) return (s % 60 === 0 ? s / 60 : Math.floor(s / 60)) + 'm';
    return Math.floor(s / 3600) + 'h' + (Math.floor((s % 3600) / 60) ? Math.floor((s % 3600) / 60) + 'm' : '');
}

function parseDuration(str) {
    const m = String(str).trim().match(/^(\d+)([smh])$/);
    if (!m) return 30e9;
    return parseInt(m[1], 10) * { s: 1, m: 60, h: 3600 }[m[2]] * 1e9;
}

function humanMs(ms) {
    if (ms == null) return '—';
    const s = Math.round(ms / 1000);
    if (s < 60) return s + 's';
    if (s < 3600) return Math.floor(s / 60) + 'm' + (s % 60 ? ' ' + String(s % 60).padStart(2, '0') + 's' : '');
    return Math.floor(s / 3600) + 'h ' + String(Math.floor((s % 3600) / 60)).padStart(2, '0') + 'm';
}

function timeAgo(dateStr) {
    if (!dateStr) return 'never';
    const s = Math.max(0, Math.round((Date.now() - new Date(dateStr).getTime()) / 1000));
    if (s < 60) return s + 's ago';
    if (s < 3600) return Math.floor(s / 60) + 'm ago';
    if (s < 86400) return Math.floor(s / 3600) + 'h ago';
    return new Date(dateStr).toLocaleString();
}

function endpointOf(target) {
    const c = target.config || {};
    if (target.type === 'tcp') return `${c.host || ''}:${c.port || ''}`;
    if (target.type === 'dns') return `${c.domain || ''} · ${(c.record_type || 'A').toUpperCase()}`;
    return c.url || '—';
}

function cssVar(name) {
    return getComputedStyle(document.documentElement).getPropertyValue(name).trim();
}

async function fetchAPI(endpoint, options = {}) {
    const res = await fetch(endpoint, {
        ...options,
        headers: { 'Content-Type': 'application/json', ...options.headers },
    });
    if (res.status === 401) {
        window.location.href = '/login';
        throw new Error('Session expired');
    }
    if (!res.ok) {
        const err = await res.json().catch(() => ({ error: res.statusText }));
        throw new Error(err.message || err.error || res.statusText);
    }
    return res.status === 204 ? null : res.json();
}

/* ═══════════════ State ═══════════════ */

const S = {
    targets: [],          // [{ t, results[], status, lastMs, lastCheckAt, uptime }]
    incidents: [],
    notifiers: [],
    tab: 'dashboard',
    query: '',
    filter: 'all',
    sort: 'status',
    view: localStorage.getItem('hm-view') || 'cards',
    resolvedOpen: false,
    loaded: false,
};

let renderDebounce = null;
let incidentsDebounce = null;
let eventStream = null;

/* ═══════════════ Small components ═══════════════ */

function segRender(el, options, value, onChange) {
    el.innerHTML = options.map(o => `
        <button type="button" role="tab" class="hm-seg-btn" aria-selected="${o.value === value}" data-value="${esc(o.value)}">
            ${o.icon ? icon(o.icon, 13) : ''}${esc(o.label)}${o.count != null ? `<span class="num" style="font-size:var(--text-xs);color:var(--text-3);">${o.count}</span>` : ''}
        </button>`).join('');
    el.querySelectorAll('.hm-seg-btn').forEach(btn => {
        btn.addEventListener('click', () => onChange(btn.dataset.value));
    });
}

function statusDot(status, { pulse = false, lg = false } = {}) {
    return `<span class="hm-dot${lg ? ' hm-dot--lg' : ''}${pulse ? ' hm-dot--pulse' : ''}" data-status="${status}"></span>`;
}

const BADGE = {
    up: { icon: 'check', label: 'Up' },
    down: { icon: 'alert-triangle', label: 'Down' },
    degraded: { icon: 'zap', label: 'Degraded' },
    unknown: { icon: 'pause', label: 'Unknown' },
};

function statusBadge(status, label) {
    const m = BADGE[status] || BADGE.unknown;
    return `<span class="hm-badge" data-status="${status}">${icon(m.icon, 12)}${esc(label || m.label)}</span>`;
}

function tickBar(results, { slots = 45, height = 18 } = {}) {
    // results come newest-first from the API; ticks render oldest → newest
    const shown = results.slice(0, slots).reverse();
    const pad = Math.max(0, slots - shown.length);
    const downs = shown.filter(r => r.status === 'failure').length;
    const ticks = shown.map(r => {
        const s = toUiStatus(r.status);
        const time = r.checked_at ? new Date(r.checked_at).toLocaleTimeString() : '';
        const tip = r.status === 'failure'
            ? `${time} · ${r.error || r.message || 'failed'}`
            : `${time} · ${r.response_time_ms} ms${r.status === 'warning' ? ' · slow' : ''}`;
        return `<span class="hm-tick" data-s="${s}" data-tip="${esc(tip)}"></span>`;
    }).join('');
    return `<div class="hm-ticks" style="height:${height}px" role="img" aria-label="Last ${shown.length} checks, ${downs} failed">${'<span class="hm-tick"></span>'.repeat(pad)}${ticks}</div>`;
}

function emptyState({ icon: name = 'activity', title, body, actionHtml = '' }) {
    return `<div class="hm-empty">${icon(name, 28)}<b>${esc(title)}</b>${body ? `<p>${esc(body)}</p>` : ''}${actionHtml}</div>`;
}

function skeletonCards(n = 4) {
    const one = `<div class="hm-card" style="padding:12px 16px;display:flex;flex-direction:column;gap:10px;">
        <div style="display:flex;gap:8px;align-items:center;"><span class="hm-skel" style="display:block;width:8px;height:8px;border-radius:999px;"></span><span class="hm-skel" style="display:block;width:45%;height:14px;"></span></div>
        <span class="hm-skel" style="display:block;width:70%;height:11px;"></span>
        <span class="hm-skel" style="display:block;width:100%;height:18px;"></span>
        <div style="display:flex;gap:12px;"><span class="hm-skel" style="display:block;width:52px;height:12px;"></span><span class="hm-skel" style="display:block;width:52px;height:12px;"></span></div>
    </div>`;
    return `<div class="hm-grid">${one.repeat(n)}</div>`;
}

/* ── Toasts ── */

function toast(kind, title, message) {
    const el = document.createElement('div');
    el.className = 'hm-toast';
    el.dataset.kind = kind;
    el.setAttribute('role', 'status');
    const iconName = { success: 'check-circle', danger: 'alert-triangle', warning: 'zap', info: 'info' }[kind] || 'info';
    el.innerHTML = `${icon(iconName, 16)}<div style="flex:1;min-width:0;"><b>${esc(title)}</b>${message ? `<p>${esc(message)}</p>` : ''}</div>
        <button type="button" class="hm-iconbtn" style="width:24px;height:24px;margin:-4px -6px 0 0;" aria-label="Dismiss">${icon('x', 13)}</button>`;
    el.querySelector('button').addEventListener('click', () => el.remove());
    document.getElementById('toasts').appendChild(el);
    setTimeout(() => el.remove(), 5200);
}

/* ── Overlays: slide-over panel + confirm dialog ── */

const overlayRoot = () => document.getElementById('overlay');
let panelOnClose = null;

function closePanel() {
    overlayRoot().innerHTML = '';
    if (panelOnClose) { const f = panelOnClose; panelOnClose = null; f(); }
}

function openPanel({ titleHtml, titleExtraHtml = '', width = 560, bodyHtml, footerHtml = '' }) {
    const root = overlayRoot();
    root.innerHTML = `
        <div class="hm-backdrop"></div>
        <div class="hm-slideover" role="dialog" aria-modal="true" style="width:min(${width}px,100vw);">
            <div class="hm-slideover-head">
                <b style="font-size:var(--text-xl);font-weight:var(--weight-semibold);flex:1;min-width:0;display:flex;align-items:center;gap:var(--space-2);">${titleHtml}</b>
                ${titleExtraHtml}
                <button type="button" class="hm-iconbtn" aria-label="Close panel" data-close>${icon('x', 16)}</button>
            </div>
            <div class="hm-slideover-body scroll-thin"></div>
            ${footerHtml ? `<div class="hm-slideover-foot">${footerHtml}</div>` : ''}
        </div>`;
    root.querySelector('.hm-slideover-body').innerHTML = bodyHtml;
    root.querySelector('.hm-backdrop').addEventListener('click', closePanel);
    root.querySelector('[data-close]').addEventListener('click', closePanel);
    root.querySelector('[data-close]').focus();
    return root.querySelector('.hm-slideover');
}

function confirmDialog({ title, body, confirmLabel = 'Delete' }, onConfirm) {
    const root = overlayRoot();
    root.innerHTML = `
        <div class="hm-backdrop"></div>
        <div class="hm-dialog" role="dialog" aria-modal="true" aria-label="${esc(title)}">
            <b style="display:block;font-size:var(--text-xl);font-weight:var(--weight-semibold);">${esc(title)}</b>
            <p style="margin:var(--space-2) 0 var(--space-5);color:var(--text-2);font-size:var(--text-md);">${esc(body)}</p>
            <div style="display:flex;justify-content:flex-end;gap:var(--space-2);">
                <button type="button" class="hm-btn" data-cancel>Cancel</button>
                <button type="button" class="hm-btn hm-btn--danger" data-ok>${esc(confirmLabel)}</button>
            </div>
        </div>`;
    root.querySelector('.hm-backdrop').addEventListener('click', closePanel);
    root.querySelector('[data-cancel]').addEventListener('click', closePanel);
    root.querySelector('[data-ok]').addEventListener('click', () => { closePanel(); onConfirm(); });
    root.querySelector('[data-cancel]').focus();
}

document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape' && overlayRoot().childElementCount) closePanel();
});

/* ═══════════════ Data loading ═══════════════ */

async function loadTargets() {
    const raw = await fetchAPI('/api/v1/targets') || [];
    const withResults = await Promise.all(raw.map(async (t) => {
        let results = [];
        try { results = await fetchAPI(`/api/v1/targets/${t.id}/results?limit=45`) || []; } catch (_) { /* keep empty */ }
        return makeEntry(t, results);
    }));
    S.targets = withResults;
}

function makeEntry(t, results) {
    const last = results[0] || null;
    const status = !t.enabled ? 'unknown' : toUiStatus(last?.status);
    const counted = results.length;
    const ok = results.filter(r => r.status === 'success').length;
    return {
        t,
        results,
        status,
        paused: !t.enabled,
        lastMs: last && last.status !== 'failure' ? last.response_time_ms : null,
        lastCheckAt: last?.checked_at || null,
        uptime: counted ? (100 * ok / counted).toFixed(1) : null,
    };
}

async function loadIncidents() {
    S.incidents = await fetchAPI('/api/v1/incidents?limit=20') || [];
}

async function loadNotifiers() {
    S.notifiers = await fetchAPI('/api/v1/notifiers') || [];
}

async function refreshData() {
    try {
        await Promise.all([loadTargets(), loadIncidents(), loadNotifiers()]);
        S.loaded = true;
        renderAll();
        document.getElementById('lastUpdated').textContent = new Date().toLocaleTimeString();
    } catch (err) {
        if (!S.loaded) {
            document.getElementById('targets').innerHTML =
                `<div class="hm-card">${emptyState({ icon: 'alert-triangle', title: 'Failed to load', body: err.message })}</div>`;
        }
        toast('danger', 'Refresh failed', err.message);
    }
}

function scheduleRender() {
    if (renderDebounce) return;
    renderDebounce = setTimeout(() => { renderDebounce = null; renderAll(); }, 250);
}

/* ═══════════════ Rendering ═══════════════ */

function renderAll() {
    renderTabs();
    if (S.tab === 'dashboard') {
        renderHero();
        renderStatRow();
        renderOngoing();
        renderToolbarSegs();
        renderTargets();
        renderResolved();
    } else {
        renderNotifiers();
    }
    document.getElementById('view-dashboard').hidden = S.tab !== 'dashboard';
    document.getElementById('view-notifiers').hidden = S.tab !== 'notifiers';
}

function renderTabs() {
    segRender(document.getElementById('tabSeg'), [
        { value: 'dashboard', label: 'Dashboard' },
        { value: 'notifiers', label: 'Notifiers', count: S.notifiers.length },
    ], S.tab, (v) => { S.tab = v; renderAll(); });
}

function renderHero() {
    const downs = S.targets.filter(e => e.status === 'down');
    const degraded = S.targets.filter(e => e.status === 'degraded');
    const active = S.targets.filter(e => !e.paused);
    const ok = downs.length === 0;
    const el = document.getElementById('hero');
    el.className = 'hm-hero';
    el.dataset.status = ok ? 'up' : 'down';
    el.style.marginTop = 'var(--space-4)';
    const sub = ok
        ? `${active.length} target${active.length === 1 ? '' : 's'} · all checks passing${degraded.length ? ` · ${degraded.length} degraded` : ''}`
        : downs.map(e => e.t.name).join(', ') + (degraded.length ? ` · ${degraded.length} degraded` : '');
    el.innerHTML = `
        <span class="hm-hero-icon">${icon(ok ? 'check' : 'alert-triangle', 22)}</span>
        <div style="flex:1;min-width:0;">
            <h1>${ok ? 'All systems operational' : `${downs.length} target${downs.length > 1 ? 's' : ''} down`}</h1>
            <p>${esc(S.loaded ? sub : 'Loading…')}</p>
        </div>
        ${ok ? '' : `<div style="display:flex;gap:8px;flex-wrap:wrap;">${downs.map(e => statusBadge('down', e.t.name)).join('')}</div>`}`;
}

function renderStatRow() {
    const all = S.targets.length;
    const paused = S.targets.filter(e => e.paused).length;
    const down = S.targets.filter(e => e.status === 'down').length;
    const degraded = S.targets.filter(e => e.status === 'degraded').length;
    const up = S.targets.filter(e => e.status === 'up').length;
    const withMs = S.targets.filter(e => e.lastMs != null);
    const avg = withMs.length ? Math.round(withMs.reduce((a, e) => a + e.lastMs, 0) / withMs.length) : null;
    const card = (label, value, { tone, unit, sub } = {}) => {
        const color = tone === 'danger' ? 'var(--status-down)' : tone === 'success' ? 'var(--status-up)' : tone === 'warning' ? 'var(--status-degraded)' : 'var(--text-1)';
        return `<div class="hm-card hm-statcard"><span class="overline">${esc(label)}</span>
            <b style="color:${color}">${value}${unit ? `<span style="font-size:var(--text-lg);font-weight:var(--weight-medium);color:var(--text-3);margin-left:3px;">${unit}</span>` : ''}</b>
            ${sub ? `<small>${esc(sub)}</small>` : '<small>&nbsp;</small>'}</div>`;
    };
    document.getElementById('statrow').innerHTML =
        card('Targets', all, { sub: `${all - paused} enabled` }) +
        card('Up', up, { tone: 'success', sub: degraded ? `${degraded} degraded` : undefined }) +
        card('Down', down, { tone: down ? 'danger' : undefined, sub: down ? 'see incidents' : undefined }) +
        card('Avg response', avg ?? '—', { unit: avg != null ? 'ms' : '', sub: 'across all targets' });
}

function incidentItem(i, last = false) {
    const ongoing = i.status === 'ongoing';
    const sev = { critical: { c: 'var(--status-down)', i: 'alert-triangle' }, warning: { c: 'var(--status-degraded)', i: 'zap' }, info: { c: 'var(--accent)', i: 'info' } }[i.severity] || { c: 'var(--status-down)', i: 'alert-triangle' };
    const dur = ongoing
        ? humanMs(Date.now() - new Date(i.started_at).getTime())
        : (i.duration ? humanMs(i.duration / 1e6) : '—');
    return `
        <div class="hm-incident" data-target-id="${esc(i.target_id || '')}" style="background:${ongoing ? 'var(--status-down-subtle)' : 'transparent'};cursor:pointer;">
            <div class="hm-incident-rail" aria-hidden="true">
                <span class="hm-dot${ongoing ? ' hm-dot--pulse' : ''}" data-status="${ongoing ? 'down' : 'up'}"></span>
                ${last ? '' : '<span class="line"></span>'}
            </div>
            <div style="flex:1;min-width:0;">
                <div style="display:flex;align-items:center;gap:var(--space-2);flex-wrap:wrap;">
                    <b style="font-size:var(--text-md);font-weight:var(--weight-semibold);">${esc(i.target_name || i.target_id)}</b>
                    <span style="display:inline-flex;align-items:center;gap:4px;font-size:var(--text-xs);font-weight:var(--weight-semibold);color:${sev.c};text-transform:uppercase;letter-spacing:var(--tracking-caps);">${icon(sev.i, 11)}${esc(i.severity || 'critical')}</span>
                    <span style="flex:1;"></span>
                    <span class="num" style="font-size:var(--text-sm);color:${ongoing ? 'var(--status-down)' : 'var(--text-3)'};font-weight:${ongoing ? 'var(--weight-semibold)' : 'var(--weight-regular)'};">${ongoing ? `ongoing · ${dur}` : dur}</span>
                </div>
                <div style="font-size:var(--text-sm);color:var(--text-2);margin-top:2px;">
                    ${new Date(i.started_at).toLocaleString()}${i.failure_count != null ? ` · <span class="num">${i.failure_count}</span> failed checks` : ''}
                </div>
                ${i.last_error ? `<code style="display:block;margin-top:4px;font-size:var(--text-xs);color:var(--text-3);overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">${esc(i.last_error)}</code>` : ''}
            </div>
        </div>`;
}

function bindIncidentClicks(container) {
    container.querySelectorAll('.hm-incident[data-target-id]').forEach(el => {
        el.addEventListener('click', () => { if (el.dataset.targetId) openDetail(el.dataset.targetId); });
    });
}

function renderOngoing() {
    const el = document.getElementById('ongoing');
    const ongoing = S.incidents.filter(i => i.status === 'ongoing');
    if (!ongoing.length) { el.innerHTML = ''; return; }
    el.innerHTML = `
        <div class="hm-sectionhead">
            <h2>Ongoing incidents</h2>
            <span class="hm-badge" data-status="down">${icon('alert-triangle', 12)}${ongoing.length} active</span>
        </div>
        <div class="hm-card hm-card--down">${ongoing.map((i, k) => incidentItem(i, k === ongoing.length - 1)).join('')}</div>`;
    bindIncidentClicks(el);
}

function renderResolved() {
    const resolved = S.incidents.filter(i => i.status !== 'ongoing');
    const btn = document.getElementById('resolvedToggle');
    btn.innerHTML = `${icon(S.resolvedOpen ? 'chevron-down' : 'chevron-right', 14)}${resolved.length} resolved`;
    btn.onclick = () => { S.resolvedOpen = !S.resolvedOpen; renderResolved(); };
    const el = document.getElementById('resolved');
    el.hidden = !S.resolvedOpen;
    if (S.resolvedOpen) {
        el.innerHTML = resolved.length
            ? `<div class="hm-card">${resolved.map((i, k) => incidentItem(i, k === resolved.length - 1)).join('')}</div>`
            : `<div class="hm-card">${emptyState({ icon: 'check-circle', title: 'No resolved incidents', body: 'Nothing has failed recently.' })}</div>`;
        bindIncidentClicks(el);
    }
}

function renderToolbarSegs() {
    const counts = {
        all: S.targets.length,
        down: S.targets.filter(e => e.status === 'down').length,
        degraded: S.targets.filter(e => e.status === 'degraded').length,
        paused: S.targets.filter(e => e.paused).length,
    };
    segRender(document.getElementById('filterSeg'), [
        { value: 'all', label: 'All', count: counts.all },
        { value: 'down', label: 'Down', count: counts.down },
        { value: 'degraded', label: 'Degraded', count: counts.degraded },
        { value: 'paused', label: 'Paused', count: counts.paused },
    ], S.filter, (v) => { S.filter = v; renderTargets(); });
    segRender(document.getElementById('viewSeg'), [
        { value: 'cards', label: '', icon: 'grid' },
        { value: 'table', label: '', icon: 'list' },
    ], S.view, (v) => { S.view = v; localStorage.setItem('hm-view', v); renderTargets(); });
}

function shownTargets() {
    const q = S.query.trim().toLowerCase();
    return S.targets
        .filter(e => {
            if (S.filter === 'down' && e.status !== 'down') return false;
            if (S.filter === 'degraded' && e.status !== 'degraded') return false;
            if (S.filter === 'paused' && !e.paused) return false;
            if (!q) return true;
            return e.t.name.toLowerCase().includes(q)
                || endpointOf(e.t).toLowerCase().includes(q)
                || (e.t.tags || []).some(x => String(x).toLowerCase().includes(q))
                || (e.t.description || '').toLowerCase().includes(q);
        })
        .sort((a, b) =>
            S.sort === 'name' ? a.t.name.localeCompare(b.t.name)
            : S.sort === 'uptime' ? (Number(a.uptime ?? 101)) - (Number(b.uptime ?? 101))
            : S.sort === 'response' ? (b.lastMs || 0) - (a.lastMs || 0)
            : (STATUS_ORDER[a.status] - STATUS_ORDER[b.status]) || a.t.name.localeCompare(b.t.name));
}

function targetActions(e) {
    return `
        <button type="button" class="hm-iconbtn" aria-label="Edit ${esc(e.t.name)}" data-act="edit" data-id="${esc(e.t.id)}">${icon('pencil', 14)}</button>
        <button type="button" class="hm-iconbtn hm-iconbtn--danger" aria-label="Delete ${esc(e.t.name)}" data-act="delete" data-id="${esc(e.t.id)}">${icon('trash', 14)}</button>`;
}

function renderTargets() {
    const el = document.getElementById('targets');
    if (!S.loaded) { el.innerHTML = skeletonCards(); return; }
    if (!S.targets.length) {
        el.innerHTML = `<div class="hm-card">${emptyState({
            icon: 'activity', title: 'No targets yet',
            body: 'Add your first HTTP, TCP or DNS check to start monitoring.',
            actionHtml: `<button type="button" class="hm-btn hm-btn--primary" data-act="add">${icon('plus', 15)}Add target</button>`,
        })}</div>`;
        el.querySelector('[data-act="add"]')?.addEventListener('click', () => openTargetForm(null));
        return;
    }
    const shown = shownTargets();
    if (!shown.length) {
        el.innerHTML = `<div class="hm-card">${emptyState({
            icon: 'search', title: 'No matching targets',
            body: S.query ? `Nothing matches “${S.query}”. Try a different search or clear the filter.` : 'Nothing matches this filter.',
            actionHtml: `<button type="button" class="hm-btn" data-act="clear">Clear filters</button>`,
        })}</div>`;
        el.querySelector('[data-act="clear"]')?.addEventListener('click', () => {
            S.query = ''; S.filter = 'all';
            document.getElementById('searchInput').value = '';
            renderToolbarSegs(); renderTargets();
        });
        return;
    }

    if (S.view === 'cards') {
        el.innerHTML = `<div class="hm-grid">${shown.map(e => {
            const down = e.status === 'down';
            return `
            <div class="hm-card hm-card--hover${down ? ' hm-card--down' : ''}" data-target-id="${esc(e.t.id)}" role="button" tabindex="0"
                 aria-label="${esc(e.t.name)}, ${e.paused ? 'paused' : e.status}"
                 style="padding:var(--space-3) var(--space-4);display:flex;flex-direction:column;gap:var(--space-2);">
                <div style="display:flex;align-items:center;gap:var(--space-2);">
                    ${statusDot(e.paused ? 'unknown' : e.status, { pulse: down })}
                    <b style="font-size:var(--text-lg);font-weight:var(--weight-semibold);overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">${esc(e.t.name)}</b>
                    <span style="color:var(--text-3);" data-tip="${esc(e.t.type.toUpperCase())}">${icon(TYPE_ICON[e.t.type] || 'globe', 13)}</span>
                    <span style="flex:1;"></span>
                    <div style="display:flex;gap:2px;" data-stop>${targetActions(e)}</div>
                </div>
                <span class="num" style="font-size:var(--text-sm);color:var(--text-3);overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">${esc(endpointOf(e.t))}</span>
                ${tickBar(e.results)}
                <div class="tcard-meta">
                    <span class="num" style="color:${down ? 'var(--status-down)' : 'var(--status-up)'};font-weight:var(--weight-semibold);">${e.paused || e.uptime == null ? '—' : e.uptime + '%'}</span>
                    <span class="num">${e.lastMs != null ? e.lastMs + ' ms' : '—'}</span>
                    <span style="color:var(--text-3);">${esc(timeAgo(e.lastCheckAt))} · every ${esc(fmtDuration(e.t.interval))}</span>
                    ${e.paused ? `<span style="color:var(--text-3);display:inline-flex;align-items:center;gap:4px;">${icon('pause', 11)}paused</span>` : ''}
                </div>
            </div>`;
        }).join('')}</div>`;
    } else {
        el.innerHTML = `<div class="hm-card hm-tablewrap scroll-thin"><table class="hm-table">
            <thead><tr>
                <th style="width:24px;"><span class="sr-only">Status</span></th>
                <th>Name</th><th>Type</th><th style="width:200px;">Last 45 checks</th>
                <th style="text-align:right;">Uptime</th><th style="text-align:right;">Resp.</th><th style="width:76px;"></th>
            </tr></thead>
            <tbody>${shown.map(e => `
                <tr data-target-id="${esc(e.t.id)}" data-status="${e.paused ? 'unknown' : e.status}" tabindex="0">
                    <td>${statusDot(e.paused ? 'unknown' : e.status, { pulse: e.status === 'down' })}</td>
                    <td><b style="font-weight:var(--weight-semibold);">${esc(e.t.name)}</b>${e.paused ? `<span style="margin-left:8px;font-size:var(--text-xs);color:var(--text-3);">paused</span>` : ''}</td>
                    <td><span style="display:inline-flex;align-items:center;gap:5px;color:var(--text-2);font-size:var(--text-sm);">${icon(TYPE_ICON[e.t.type] || 'globe', 13)}${esc(e.t.type.toUpperCase())}</span></td>
                    <td>${tickBar(e.results, { height: 14 })}</td>
                    <td class="num" style="text-align:right;color:${e.status === 'down' ? 'var(--status-down)' : 'var(--status-up)'};font-weight:var(--weight-medium);">${e.paused || e.uptime == null ? '—' : e.uptime + '%'}</td>
                    <td class="num" style="text-align:right;color:var(--text-2);">${e.lastMs != null ? e.lastMs + ' ms' : '—'}</td>
                    <td data-stop><div style="display:flex;gap:2px;justify-content:flex-end;">${targetActions(e)}</div></td>
                </tr>`).join('')}</tbody>
        </table></div>`;
    }

    // open detail on card/row click; buttons stop propagation
    el.querySelectorAll('[data-target-id]').forEach(node => {
        node.addEventListener('click', (ev) => {
            if (ev.target.closest('[data-stop]')) return;
            openDetail(node.dataset.targetId);
        });
        node.addEventListener('keydown', (ev) => {
            if (ev.key === 'Enter' && !ev.target.closest('[data-stop]')) openDetail(node.dataset.targetId);
        });
    });
    el.querySelectorAll('[data-act="edit"]').forEach(b => b.addEventListener('click', () => {
        const e = S.targets.find(x => x.t.id === b.dataset.id);
        if (e) openTargetForm(e.t);
    }));
    el.querySelectorAll('[data-act="delete"]').forEach(b => b.addEventListener('click', () => {
        const e = S.targets.find(x => x.t.id === b.dataset.id);
        if (e) confirmDialog({
            title: 'Delete target?',
            body: `“${e.t.name}” and its check history will be permanently removed. Notifiers are not affected.`,
        }, async () => {
            try {
                await fetchAPI(`/api/v1/targets/${e.t.id}`, { method: 'DELETE' });
                toast('info', 'Target deleted', `${e.t.name} is no longer monitored`);
                refreshData();
            } catch (err) { toast('danger', 'Delete failed', err.message); }
        });
    }));
}

/* ── Notifiers ── */

const NOTIFIER_META = {
    telegram: { icon: 'send', label: 'Telegram' },
    email: { icon: 'mail', label: 'Email (SMTP)' },
    gmail: { icon: 'mail', label: 'Gmail' },
    gmail_oauth: { icon: 'mail', label: 'Gmail (OAuth)' },
    webhook: { icon: 'link', label: 'Webhook' },
};

function notifierDetail(n) {
    const c = n.config || {};
    if (n.type === 'gmail_oauth' || n.type === 'gmail') return `${c.from || '—'} → ${(c.to || []).join(', ')}`;
    if (n.type === 'telegram') return `chat ${(c.chat_ids || []).join(', ')}${c.proxy_url ? ' · proxy' : ''}`;
    if (n.type === 'email') return `${c.smtp_host || ''}:${c.smtp_port || ''} · ${c.from || '—'}`;
    if (n.type === 'webhook') return `${c.method || 'POST'} ${c.url || '—'}`;
    return '';
}

function renderNotifiers() {
    const el = document.getElementById('notifiers');
    if (!S.notifiers.length) {
        el.innerHTML = `<div class="hm-card">${emptyState({
            icon: 'bell', title: 'No notifiers yet',
            body: 'Add a Telegram, email or webhook channel to receive alerts.',
            actionHtml: `<button type="button" class="hm-btn hm-btn--primary" data-act="add">${icon('plus', 15)}Add notifier</button>`,
        })}</div>`;
        el.querySelector('[data-act="add"]')?.addEventListener('click', () => openNotifierForm(null));
        return;
    }
    el.innerHTML = `<div class="hm-card">${S.notifiers.map((n, k) => {
        const meta = NOTIFIER_META[n.type] || NOTIFIER_META.webhook;
        return `
        <div style="display:flex;align-items:center;gap:var(--space-3);padding:var(--space-3) var(--space-4);${k ? 'border-top:1px solid var(--border-1);' : ''}">
            <span style="display:flex;align-items:center;justify-content:center;width:32px;height:32px;flex:none;border-radius:var(--radius-md);background:var(--bg-2);border:1px solid var(--border-1);color:${n.enabled ? 'var(--text-2)' : 'var(--text-3)'};">${icon(meta.icon, 15)}</span>
            <div style="flex:1;min-width:0;opacity:${n.enabled ? 1 : 0.55};">
                <b style="display:block;font-size:var(--text-md);font-weight:var(--weight-semibold);">${esc(n.id)}</b>
                <span style="font-size:var(--text-sm);color:var(--text-3);">${esc(meta.label)} · <code style="font-size:var(--text-xs);">${esc(notifierDetail(n))}</code></span>
            </div>
            <label class="hm-switch" aria-label="Toggle ${esc(n.id)}">
                <input type="checkbox" role="switch" ${n.enabled ? 'checked' : ''} data-act="toggle" data-id="${esc(n.id)}"><i></i>
            </label>
            <button type="button" class="hm-iconbtn" aria-label="Edit ${esc(n.id)}" data-act="edit" data-id="${esc(n.id)}">${icon('pencil', 14)}</button>
            <button type="button" class="hm-iconbtn hm-iconbtn--danger" aria-label="Delete ${esc(n.id)}" data-act="del" data-id="${esc(n.id)}">${icon('trash', 14)}</button>
        </div>`;
    }).join('')}</div>`;

    el.querySelectorAll('[data-act="toggle"]').forEach(cb => cb.addEventListener('change', async () => {
        try {
            const cfg = await fetchAPI(`/api/v1/notifiers/${cb.dataset.id}`);
            cfg.enabled = cb.checked;
            await fetchAPI(`/api/v1/notifiers/${cb.dataset.id}`, { method: 'PUT', body: JSON.stringify(cfg) });
            toast(cb.checked ? 'success' : 'info', cb.checked ? 'Notifier enabled' : 'Notifier disabled', cb.dataset.id);
            await loadNotifiers(); renderNotifiers();
        } catch (err) { toast('danger', 'Toggle failed', err.message); cb.checked = !cb.checked; }
    }));
    el.querySelectorAll('[data-act="edit"]').forEach(b => b.addEventListener('click', () => {
        const n = S.notifiers.find(x => x.id === b.dataset.id);
        if (n) openNotifierForm(n);
    }));
    el.querySelectorAll('[data-act="del"]').forEach(b => b.addEventListener('click', () => {
        confirmDialog({ title: 'Delete notifier?', body: `“${b.dataset.id}” will stop receiving alerts.` }, async () => {
            try {
                await fetchAPI(`/api/v1/notifiers/${b.dataset.id}`, { method: 'DELETE' });
                toast('info', 'Notifier deleted', b.dataset.id);
                await loadNotifiers(); renderAll();
            } catch (err) { toast('danger', 'Delete failed', err.message); }
        });
    }));
}

/* ═══════════════ Target detail (slide-over) ═══════════════ */

const PERIODS = { '24h': '24h', '7d': '168h', '30d': '720h' };
const PERIOD_MS = { '24h': 864e5, '7d': 6048e5, '30d': 2592e6 };
let detail = null; // { target, period, results, chart }

function ensureChartJs() {
    if (window.Chart) return Promise.resolve();
    return new Promise((resolve, reject) => {
        const s = document.createElement('script');
        s.src = '/static/js/chart.umd.min.js';
        s.onload = resolve;
        s.onerror = () => reject(new Error('Failed to load chart library'));
        document.head.appendChild(s);
    });
}

async function openDetail(targetId) {
    let target;
    try { target = await fetchAPI(`/api/v1/targets/${targetId}`); }
    catch (err) { toast('danger', 'Failed to open target', err.message); return; }

    const entry = S.targets.find(e => e.t.id === targetId);
    const status = entry ? entry.status : 'unknown';
    detail = { target, period: '24h', results: [], chart: null };

    const panel = openPanel({
        width: 720,
        titleHtml: `${statusDot(status, { lg: true, pulse: status === 'down' })}<span style="overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">${esc(target.name)}</span><span style="color:var(--text-3);">${icon(TYPE_ICON[target.type] || 'globe', 14)}</span>`,
        titleExtraHtml: statusBadge(status, !target.enabled ? 'Paused' : undefined),
        bodyHtml: `
            <div style="display:flex;align-items:center;gap:8px;flex-wrap:wrap;">
                <div class="hm-seg" role="tablist" aria-label="Stats period" id="dPeriodSeg"></div>
                <span style="flex:1;"></span>
                <span class="hm-live" data-state="live"><span class="hm-dot"></span>Checked ${esc(timeAgo(entry?.lastCheckAt))} · every ${esc(fmtDuration(target.interval))}</span>
            </div>
            <div class="detail-statrow" id="dStats"><span class="hm-skel" style="display:block;width:100%;height:58px;"></span></div>
            <div class="detail-section">
                <div class="detail-section-head"><span class="overline" id="dChartTitle">Response time — 24h</span></div>
                <div class="chart-box"><canvas id="dChart"></canvas></div>
            </div>
            <div class="detail-section">
                <div class="detail-section-head"><span class="overline">Recent checks</span></div>
                <div id="dTicks" style="margin-bottom:8px;"></div>
                <div class="hm-card" style="overflow:hidden;border-radius:var(--radius-md);" id="dChecks"></div>
            </div>
            <div class="detail-section">
                <div class="detail-section-head"><span class="overline">Incidents</span></div>
                <div id="dIncidents"><span class="hm-skel" style="display:block;width:100%;height:40px;"></span></div>
            </div>
            <div class="detail-section">
                <div class="detail-section-head">
                    <span class="overline">Configuration</span>
                    <span style="flex:1;"></span>
                    <button type="button" class="hm-btn hm-btn--ghost" id="dEdit">${icon('pencil', 13)}Edit</button>
                </div>
                <div class="config-box" id="dConfig"></div>
            </div>`,
    });
    panelOnClose = () => { if (detail?.chart) detail.chart.destroy(); detail = null; };

    panel.querySelector('#dEdit').addEventListener('click', () => { const t = detail.target; closePanel(); openTargetForm(t); });
    renderDetailPeriodSeg(panel);
    renderDetailConfig(panel);
    loadDetailStats(panel);
    loadDetailHistory(panel);
    loadDetailIncidents(panel);
}

function renderDetailPeriodSeg(panel) {
    segRender(panel.querySelector('#dPeriodSeg'), [
        { value: '24h', label: '24h' }, { value: '7d', label: '7d' }, { value: '30d', label: '30d' },
    ], detail.period, (v) => {
        detail.period = v;
        renderDetailPeriodSeg(panel);
        panel.querySelector('#dChartTitle').textContent = `Response time — ${v}`;
        loadDetailStats(panel);
        renderDetailChart(panel); // same fetched series, re-filtered by period
    });
}

async function loadDetailStats(panel) {
    const el = panel.querySelector('#dStats');
    try {
        const s = await fetchAPI(`/api/v1/targets/${detail.target.id}/stats?period=${PERIODS[detail.period]}`);
        const uptime = (s.uptime_percentage ?? 0).toFixed(2);
        const upColor = uptime >= 99 ? 'var(--status-up)' : uptime >= 90 ? 'var(--status-degraded)' : 'var(--status-down)';
        const stat = (label, value, unit = '', color = 'var(--text-1)') => `
            <div class="detail-stat"><span class="overline">${label}</span>
            <b class="num" style="color:${color}">${value}${unit ? `<span style="font-size:var(--text-sm);color:var(--text-3);margin-left:2px;">${unit}</span>` : ''}</b></div>`;
        el.innerHTML =
            stat('Uptime', uptime, '%', upColor) +
            stat('Avg', Math.round(s.avg_response_time_ms ?? 0), 'ms') +
            stat('Min / Max', `${s.min_response_time_ms ?? 0}–${s.max_response_time_ms ?? 0}`, 'ms') +
            stat('Failed', s.failed_checks ?? 0, `/ ${s.total_checks ?? 0}`, s.failed_checks ? 'var(--status-down)' : 'var(--text-1)') +
            (s.consecutive_failures ? stat('Consecutive', s.consecutive_failures, 'fails', 'var(--status-down)') : '');
    } catch (err) {
        el.innerHTML = `<span style="color:var(--status-down);font-size:var(--text-sm);">Failed to load stats: ${esc(err.message)}</span>`;
    }
}

async function loadDetailHistory(panel) {
    try {
        detail.results = await fetchAPI(`/api/v1/targets/${detail.target.id}/results?limit=200`) || [];
    } catch (err) {
        panel.querySelector('#dChecks').innerHTML = emptyState({ icon: 'alert-triangle', title: 'Failed to load checks', body: err.message });
        return;
    }
    panel.querySelector('#dTicks').innerHTML = tickBar(detail.results, { height: 16 });
    renderDetailChart(panel);
    renderDetailChecks(panel);
}

async function renderDetailChart(panel) {
    if (!detail) return;
    const cutoff = Date.now() - PERIOD_MS[detail.period];
    const series = detail.results.filter(r => new Date(r.checked_at).getTime() >= cutoff).reverse();
    try { await ensureChartJs(); } catch (err) { toast('warning', 'Chart unavailable', err.message); return; }
    if (!detail) return; // panel closed while loading
    if (detail.chart) { detail.chart.destroy(); detail.chart = null; }
    const canvas = panel.querySelector('#dChart');
    if (!canvas) return;
    const colors = { up: cssVar('--status-up'), down: cssVar('--status-down'), degraded: cssVar('--status-degraded'), unknown: cssVar('--status-unknown') };
    detail.chart = new Chart(canvas.getContext('2d'), {
        type: 'line',
        data: {
            labels: series.map(r => new Date(r.checked_at).toLocaleTimeString()),
            datasets: [{
                label: 'Response time (ms)',
                data: series.map(r => r.response_time_ms),
                borderColor: cssVar('--accent'),
                backgroundColor: cssVar('--chart-fill') || 'rgba(59,130,246,0.1)',
                fill: true, tension: 0.2, borderWidth: 1.5,
                pointRadius: series.map(r => r.status === 'success' ? 1.5 : 3.5),
                pointBackgroundColor: series.map(r => colors[toUiStatus(r.status)]),
                pointBorderColor: series.map(r => colors[toUiStatus(r.status)]),
            }],
        },
        options: {
            responsive: true, maintainAspectRatio: false, animation: false,
            scales: {
                y: { beginAtZero: true, grid: { color: cssVar('--chart-grid') || 'rgba(128,128,128,0.1)' }, ticks: { color: cssVar('--text-3'), font: { size: 10 } } },
                x: { ticks: { maxTicksLimit: 8, autoSkip: true, color: cssVar('--text-3'), font: { size: 10 } }, grid: { display: false } },
            },
            plugins: { legend: { display: false } },
        },
    });
}

function renderDetailChecks(panel) {
    const el = panel.querySelector('#dChecks');
    if (!detail.results.length) {
        el.innerHTML = emptyState({ icon: 'clock', title: 'No checks yet', body: 'The first check runs on the next scheduler tick.' });
        return;
    }
    const rows = detail.results.slice(0, 12).map(r => {
        const s = toUiStatus(r.status);
        const c = s === 'down' ? 'var(--status-down)' : s === 'degraded' ? 'var(--status-degraded)' : s === 'up' ? 'var(--status-up)' : 'var(--status-unknown)';
        const ic = s === 'down' ? 'x-circle' : s === 'degraded' ? 'zap' : 'check-circle';
        const lbl = s === 'down' ? 'Failure' : s === 'degraded' ? 'Slow' : 'Success';
        const info = s === 'down' ? (r.error || r.message || 'failed') : `${r.response_time_ms} ms${r.status_code ? ` · ${r.status_code}` : ''}`;
        return `<tr style="cursor:default;">
            <td class="num" style="color:var(--text-2);">${new Date(r.checked_at).toLocaleTimeString()}</td>
            <td><span style="display:inline-flex;align-items:center;gap:6px;color:${c};font-size:var(--text-sm);font-weight:500;">${icon(ic, 13)}${lbl}</span></td>
            <td class="num" style="text-align:right;color:var(--text-2);font-size:var(--text-sm);overflow:hidden;text-overflow:ellipsis;max-width:280px;white-space:nowrap;">${esc(info)}</td>
        </tr>`;
    }).join('');
    el.innerHTML = `<table class="hm-table"><thead><tr><th>Time</th><th>Result</th><th style="text-align:right;">Detail</th></tr></thead><tbody>${rows}</tbody></table>`;
}

async function loadDetailIncidents(panel) {
    const el = panel.querySelector('#dIncidents');
    try {
        const incidents = await fetchAPI(`/api/v1/targets/${detail.target.id}/incidents?limit=20`) || [];
        el.innerHTML = incidents.length
            ? `<div class="hm-card" style="border-radius:var(--radius-md);">${incidents.map((i, k) => incidentItem({ ...i, target_id: '' }, k === incidents.length - 1)).join('')}</div>`
            : `<p style="margin:0;font-size:var(--text-md);color:var(--text-3);">No incidents in recent history.</p>`;
    } catch (err) {
        el.innerHTML = `<span style="color:var(--status-down);font-size:var(--text-sm);">Failed to load incidents: ${esc(err.message)}</span>`;
    }
}

function renderDetailConfig(panel) {
    const t = detail.target;
    const c = t.config || {};
    const row = (k, v, mono = false) => v == null || v === '' ? '' :
        `<div class="config-row"><span>${esc(k)}</span><span class="${mono ? 'num' : ''}">${esc(v)}</span></div>`;
    let typeRows = '';
    if (t.type === 'http') {
        typeRows = row('URL', c.url, true) + row('Expected status', c.expected_status_code, true);
    } else if (t.type === 'tcp') {
        typeRows = row('Host : port', `${c.host || ''}:${c.port || ''}`, true);
    } else if (t.type === 'dns') {
        typeRows = row('Domain · record', `${c.domain || ''} · ${(c.record_type || 'A').toUpperCase()}`, true)
            + row('DNS server', c.dns_server, true)
            + row('Expected values', [...(c.expected_ips || []), ...(c.expect_values || [])].join(', '), true);
    }
    panel.querySelector('#dConfig').innerHTML =
        row('Type', t.type.toUpperCase()) + typeRows +
        row('Interval', fmtDuration(t.interval), true) +
        row('Timeout', fmtDuration(t.timeout), true) +
        row('Tags', (t.tags || []).join(', ')) +
        row('Description', t.description) +
        `<div class="config-row"><span>Enabled</span><span style="color:${t.enabled ? 'var(--status-up)' : 'var(--text-3)'};">${t.enabled ? 'Yes' : 'No — paused'}</span></div>`;
}

/* ═══════════════ Forms ═══════════════ */

function fieldHtml({ id, label, value = '', placeholder = '', hint = '', type = 'text', mono = false, required = false, pattern = '', readOnly = false }) {
    return `<div class="hm-field">
        <label class="hm-label" for="${id}">${esc(label)}</label>
        <input class="hm-input${mono ? ' hm-input--mono' : ''}" id="${id}" type="${type}" value="${esc(value)}"
            placeholder="${esc(placeholder)}" ${required ? 'required' : ''} ${pattern ? `pattern="${pattern}"` : ''} ${readOnly ? 'readonly' : ''}>
        <span class="hm-hint" data-hint-for="${id}">${esc(hint)}</span>
    </div>`;
}

function setFieldError(panel, id, msg) {
    const input = panel.querySelector('#' + id);
    const hint = panel.querySelector(`[data-hint-for="${id}"]`);
    if (!input) return;
    if (msg) {
        input.classList.add('hm-input--invalid');
        input.setAttribute('aria-invalid', 'true');
        if (hint) { hint.className = 'hm-error-text'; hint.innerHTML = `${icon('alert-triangle', 13)}${esc(msg)}`; }
        input.focus();
    } else {
        input.classList.remove('hm-input--invalid');
        input.removeAttribute('aria-invalid');
    }
}

function durationFieldHtml(id, label, value, presets) {
    return `<div class="hm-field">
        <span class="hm-label">${esc(label)}</span>
        <div style="display:flex;gap:var(--space-2);align-items:center;flex-wrap:wrap;">
            <div class="hm-seg" role="group" aria-label="${esc(label)} presets" data-dur-presets="${id}">
                ${presets.map(p => `<button type="button" class="hm-seg-btn num" aria-pressed="${p === value}" data-v="${p}">${p}</button>`).join('')}
            </div>
            <input class="hm-input hm-input--mono" style="width:72px;flex:none;" id="${id}" value="${esc(value)}" aria-label="${esc(label)} custom value">
        </div>
    </div>`;
}

function bindDurationField(panel, id) {
    const wrap = panel.querySelector(`[data-dur-presets="${id}"]`);
    const input = panel.querySelector('#' + id);
    const sync = () => wrap.querySelectorAll('.hm-seg-btn').forEach(b => b.setAttribute('aria-pressed', b.dataset.v === input.value.trim()));
    wrap.querySelectorAll('.hm-seg-btn').forEach(b => b.addEventListener('click', () => { input.value = b.dataset.v; sync(); }));
    input.addEventListener('input', sync);
}

function tagInputHtml(id, values) {
    return `<div class="hm-taginput" data-taginput="${id}">
        ${values.map(v => `<span class="hm-tag" data-value="${esc(v)}">${esc(v)}<button type="button" aria-label="Remove ${esc(v)}">${icon('x', 11)}</button></span>`).join('')}
        <input placeholder="${values.length ? '' : 'Add and press Enter…'}" aria-label="Add value">
    </div>`;
}

function bindTagInput(panel, id) {
    const wrap = panel.querySelector(`[data-taginput="${id}"]`);
    const input = wrap.querySelector('input');
    const commit = () => {
        const v = input.value.trim().replace(/,$/, '');
        if (v && !tagValues(panel, id).includes(v)) {
            const tag = document.createElement('span');
            tag.className = 'hm-tag';
            tag.dataset.value = v;
            tag.innerHTML = `${esc(v)}<button type="button" aria-label="Remove ${esc(v)}">${icon('x', 11)}</button>`;
            tag.querySelector('button').addEventListener('click', () => tag.remove());
            wrap.insertBefore(tag, input);
        }
        input.value = '';
    };
    wrap.addEventListener('click', () => input.focus());
    wrap.querySelectorAll('.hm-tag button').forEach(b => b.addEventListener('click', () => b.parentElement.remove()));
    input.addEventListener('keydown', (e) => {
        if (e.key === 'Enter' || e.key === ',') { e.preventDefault(); commit(); }
        else if (e.key === 'Backspace' && !input.value) wrap.querySelector('.hm-tag:last-of-type')?.remove();
    });
    input.addEventListener('blur', commit);
}

function tagValues(panel, id) {
    return [...panel.querySelectorAll(`[data-taginput="${id}"] .hm-tag`)].map(t => t.dataset.value);
}

/* ── Target form ── */

function openTargetForm(target) {
    const isEdit = !!target;
    const t = target || { type: 'http', enabled: true, config: {}, tags: [] };
    const c = t.config || {};
    let curType = t.type || 'http';

    const panel = openPanel({
        width: 520,
        titleHtml: esc(isEdit ? `Edit ${t.name}` : 'Add target'),
        bodyHtml: `
            <form id="targetForm" novalidate>
            <div class="hm-formsection">
                <span class="overline">Identity</span>
                ${fieldHtml({ id: 'fName', label: 'Name', value: t.name || '', placeholder: 'API /health', hint: 'Shown on the dashboard and in alerts', required: true })}
                ${fieldHtml({ id: 'fId', label: 'Target ID', value: t.id || '', placeholder: 'my-api-service', mono: true, hint: isEdit ? 'IDs are permanent' : 'Lowercase letters, numbers and hyphens', pattern: '[a-z0-9-]+', readOnly: isEdit, required: true })}
                ${fieldHtml({ id: 'fDesc', label: 'Description', value: t.description || '', placeholder: 'Optional note for your future self' })}
                <div class="hm-field"><span class="hm-label">Tags</span>${tagInputHtml('fTags', t.tags || [])}<span class="hm-hint">Used for search; Enter to add</span></div>
            </div>
            <div class="hm-formsection">
                <span class="overline">Check</span>
                <div class="hm-field"><span class="hm-label">Type</span><div class="hm-seg" role="tablist" aria-label="Target type" id="fTypeSeg"></div></div>
                <div id="fTypeFields"></div>
                ${durationFieldHtml('fInterval', 'Check interval', isEdit ? fmtDuration(t.interval) : '30s', ['30s', '1m', '5m', '15m', '1h'])}
                ${durationFieldHtml('fTimeout', 'Timeout', isEdit ? fmtDuration(t.timeout) : '10s', ['3s', '5s', '10s', '30s'])}
            </div>
            <div class="hm-formsection">
                <span class="overline">Advanced</span>
                <label style="display:inline-flex;align-items:center;gap:var(--space-2);cursor:pointer;">
                    <span class="hm-switch"><input type="checkbox" role="switch" id="fEnabled" ${t.enabled ? 'checked' : ''}><i></i></span>
                    <span style="font-size:var(--text-md);">Enabled — run checks on schedule</span>
                </label>
            </div>
            </form>`,
        footerHtml: `
            <button type="button" class="hm-btn" data-cancel>Cancel</button>
            <button type="button" class="hm-btn hm-btn--primary" data-save>${isEdit ? 'Save changes' : 'Add target'}</button>`,
    });

    const typeFields = () => {
        const el = panel.querySelector('#fTypeFields');
        if (curType === 'http') {
            el.innerHTML = fieldHtml({ id: 'fUrl', label: 'URL', value: c.url || '', placeholder: 'https://api.example.com/health', mono: true, required: true })
                + `<div class="field-cols">${fieldHtml({ id: 'fStatus', label: 'Expected status', value: c.expected_status_code || 200, mono: true, type: 'number' })}<div></div></div>`;
        } else if (curType === 'tcp') {
            el.innerHTML = `<div class="field-cols">
                ${fieldHtml({ id: 'fHost', label: 'Host', value: c.host || '', placeholder: 'db.example.com', mono: true, required: true })}
                ${fieldHtml({ id: 'fPort', label: 'Port', value: c.port || '', placeholder: '5432', mono: true, type: 'number', required: true })}
            </div>`;
        } else {
            el.innerHTML = `<div class="field-cols">
                ${fieldHtml({ id: 'fDomain', label: 'Domain', value: c.domain || '', placeholder: 'example.com', mono: true, required: true })}
                <div class="hm-field"><label class="hm-label" for="fRecord">Record</label>
                    <select class="hm-select" id="fRecord">${['A', 'AAAA', 'CNAME', 'MX', 'TXT', 'NS'].map(r => `<option ${((c.record_type || 'A').toUpperCase() === r) ? 'selected' : ''}>${r}</option>`).join('')}</select>
                <span class="hm-hint"></span></div>
            </div>`
                + fieldHtml({ id: 'fDnsServer', label: 'DNS server (optional)', value: c.dns_server || '', placeholder: '8.8.8.8 — default: system resolver', mono: true })
                + fieldHtml({ id: 'fDnsExpect', label: 'Expected values (optional)', value: [...(c.expected_ips || []), ...(c.expect_values || [])].join(', '), placeholder: '1.2.3.4, 5.6.7.8', mono: true, hint: 'Comma-separated; resolution must contain these' });
        }
    };
    const renderTypeSeg = () => {
        segRender(panel.querySelector('#fTypeSeg'), [
            { value: 'http', label: 'HTTP', icon: 'globe' },
            { value: 'tcp', label: 'TCP', icon: 'server' },
            { value: 'dns', label: 'DNS', icon: 'at-sign' },
        ], curType, (v) => { curType = v; renderTypeSeg(); typeFields(); });
    };
    renderTypeSeg();
    typeFields();
    bindDurationField(panel, 'fInterval');
    bindDurationField(panel, 'fTimeout');
    bindTagInput(panel, 'fTags');

    panel.querySelector('[data-cancel]').addEventListener('click', closePanel);
    panel.querySelector('[data-save]').addEventListener('click', async () => {
        const val = (id) => panel.querySelector('#' + id)?.value.trim() ?? '';
        const name = val('fName');
        if (!name) { setFieldError(panel, 'fName', 'Name is required'); return; }
        const id = val('fId');
        if (!isEdit && !/^[a-z0-9-]+$/.test(id)) { setFieldError(panel, 'fId', 'Lowercase letters, numbers and hyphens only'); return; }

        let config;
        if (curType === 'http') {
            const url = val('fUrl');
            if (!/^https?:\/\/.+/.test(url)) { setFieldError(panel, 'fUrl', 'Must start with http:// or https://'); return; }
            config = { url, method: 'GET', expected_status_code: parseInt(val('fStatus'), 10) || 200, validate_ssl: true, check_ssl_expiry: true };
        } else if (curType === 'tcp') {
            const host = val('fHost');
            const port = parseInt(val('fPort'), 10);
            if (!host) { setFieldError(panel, 'fHost', 'Host is required'); return; }
            if (!port || port < 1 || port > 65535) { setFieldError(panel, 'fPort', 'Port must be 1–65535'); return; }
            config = { host, port };
        } else {
            const domain = val('fDomain');
            if (!domain) { setFieldError(panel, 'fDomain', 'Domain is required'); return; }
            config = { domain, record_type: panel.querySelector('#fRecord').value };
            if (val('fDnsServer')) config.dns_server = val('fDnsServer');
            const expect = val('fDnsExpect').split(',').map(s => s.trim()).filter(Boolean);
            if (expect.length) config.expect_values = expect;
        }
        if (!/^\d+[smh]$/.test(val('fInterval'))) { setFieldError(panel, 'fInterval', 'Use formats like 30s, 1m, 5m'); return; }
        if (!/^\d+[smh]$/.test(val('fTimeout'))) { setFieldError(panel, 'fTimeout', 'Use formats like 10s'); return; }

        const payload = {
            id, name, type: curType,
            enabled: panel.querySelector('#fEnabled').checked,
            interval: parseDuration(val('fInterval')),
            timeout: parseDuration(val('fTimeout')),
            description: val('fDesc'),
            tags: tagValues(panel, 'fTags'),
            config,
        };
        try {
            if (isEdit) await fetchAPI(`/api/v1/targets/${t.id}`, { method: 'PUT', body: JSON.stringify(payload) });
            else await fetchAPI('/api/v1/targets', { method: 'POST', body: JSON.stringify(payload) });
            closePanel();
            toast('success', isEdit ? 'Target saved' : 'Target added', 'Checks start on the next tick');
            refreshData();
        } catch (err) { toast('danger', 'Save failed', err.message); }
    });
}

/* ── Notifier form ── */

const NOTIFIER_TYPES = [
    { value: 'telegram', label: 'Telegram', icon: 'send', hint: 'Bot message to a chat' },
    { value: 'email', label: 'Email (SMTP)', icon: 'mail', hint: 'Any SMTP server' },
    { value: 'gmail_oauth', label: 'Gmail OAuth', icon: 'mail', hint: 'Personal account token' },
    { value: 'gmail', label: 'Gmail SA', icon: 'mail', hint: 'Service account + delegation' },
    { value: 'webhook', label: 'Webhook', icon: 'link', hint: 'HTTP request anywhere' },
];

function openNotifierForm(notifier) {
    const isEdit = !!notifier;
    const n = notifier || { type: 'telegram', enabled: true, config: {} };
    const c = n.config || {};
    let curType = n.type;

    const panel = openPanel({
        width: 520,
        titleHtml: esc(isEdit ? `Edit ${n.id}` : 'Add notifier'),
        bodyHtml: `
            <form id="notifierForm" novalidate>
            <div class="hm-formsection">
                <span class="overline">Channel</span>
                <div class="type-picker" id="nTypePicker"></div>
            </div>
            <div class="hm-formsection">
                <span class="overline">Settings</span>
                ${fieldHtml({ id: 'nId', label: 'ID', value: n.id || '', placeholder: 'ops-telegram', mono: true, pattern: '[a-z0-9-]+', readOnly: isEdit, required: true, hint: isEdit ? 'IDs are permanent' : 'Lowercase letters, numbers and hyphens' })}
                <div id="nTypeFields"></div>
            </div>
            <div class="hm-formsection">
                <span class="overline">Advanced</span>
                <label style="display:inline-flex;align-items:center;gap:var(--space-2);cursor:pointer;">
                    <span class="hm-switch"><input type="checkbox" role="switch" id="nEnabled" ${n.enabled ? 'checked' : ''}><i></i></span>
                    <span style="font-size:var(--text-md);">Enabled — send alerts through this channel</span>
                </label>
            </div>
            </form>`,
        footerHtml: `
            <button type="button" class="hm-btn" data-cancel>Cancel</button>
            <button type="button" class="hm-btn hm-btn--primary" data-save>${isEdit ? 'Save notifier' : 'Add notifier'}</button>`,
    });

    const renderPicker = () => {
        panel.querySelector('#nTypePicker').innerHTML = NOTIFIER_TYPES.map(t => `
            <button type="button" class="type-card" aria-pressed="${t.value === curType}" data-v="${t.value}" ${isEdit ? 'disabled style="cursor:not-allowed;opacity:' + (t.value === curType ? 1 : 0.45) + ';"' : ''}>
                <span class="t-title">${icon(t.icon, 14)}${esc(t.label)}</span>
                <span class="t-hint">${esc(t.hint)}</span>
            </button>`).join('');
        if (!isEdit) panel.querySelectorAll('.type-card').forEach(b => b.addEventListener('click', () => {
            curType = b.dataset.v; renderPicker(); typeFields();
        }));
    };

    const typeFields = () => {
        const el = panel.querySelector('#nTypeFields');
        const cc = curType === n.type ? c : {};
        if (curType === 'telegram') {
            el.innerHTML = fieldHtml({ id: 'nTgToken', label: 'Bot token', value: cc.bot_token || '', placeholder: '123456:ABC-DEF…', mono: true, type: 'password', hint: isEdit ? 'Leave “***” to keep the current token' : '' })
                + `<div class="hm-field"><span class="hm-label">Chat IDs</span>${tagInputHtml('nTgChats', cc.chat_ids || [])}<span class="hm-hint">Enter to add</span></div>`
                + fieldHtml({ id: 'nTgProxy', label: 'Proxy URL (optional)', value: cc.proxy_url || '', placeholder: 'socks5://user:pass@host:1080', mono: true });
        } else if (curType === 'email') {
            el.innerHTML = `<div class="field-cols">
                    ${fieldHtml({ id: 'nSmtpHost', label: 'SMTP host', value: cc.smtp_host || '', placeholder: 'smtp.gmail.com', mono: true, required: true })}
                    ${fieldHtml({ id: 'nSmtpPort', label: 'Port', value: cc.smtp_port || 587, mono: true, type: 'number' })}
                </div>`
                + fieldHtml({ id: 'nSmtpUser', label: 'SMTP user', value: cc.smtp_user || '', placeholder: 'you@gmail.com', mono: true })
                + fieldHtml({ id: 'nSmtpPass', label: 'SMTP password', value: '', type: 'password', hint: isEdit ? 'Leave blank to keep the existing password' : '' })
                + fieldHtml({ id: 'nEmailFrom', label: 'From', value: cc.from || '', placeholder: 'you@gmail.com', mono: true, required: true })
                + `<div class="hm-field"><span class="hm-label">To</span>${tagInputHtml('nEmailTo', cc.to || [])}<span class="hm-hint">Enter to add</span></div>`
                + `<label style="display:inline-flex;align-items:center;gap:var(--space-2);cursor:pointer;">
                    <span class="hm-switch"><input type="checkbox" role="switch" id="nSmtpTLS" ${cc.use_tls !== false ? 'checked' : ''}><i></i></span>
                    <span style="font-size:var(--text-md);">Use TLS</span></label>`;
        } else if (curType === 'gmail_oauth') {
            el.innerHTML = fieldHtml({ id: 'nOauthCred', label: 'Credentials file', value: cc.credentials_file || 'secrets/gmail-token.json', mono: true, hint: 'OAuth2 client JSON from Google Cloud Console' })
                + fieldHtml({ id: 'nOauthToken', label: 'Token file', value: cc.token_file || 'secrets/gmail-oauth-token.json', mono: true, hint: 'Generated by the gmail_auth helper' })
                + fieldHtml({ id: 'nOauthFrom', label: 'From', value: cc.from || '', placeholder: 'you@gmail.com', mono: true, required: true })
                + `<div class="hm-field"><span class="hm-label">To</span>${tagInputHtml('nOauthTo', cc.to || [])}<span class="hm-hint">Enter to add</span></div>`;
        } else if (curType === 'gmail') {
            el.innerHTML = fieldHtml({ id: 'nGmailSA', label: 'Service account file', value: cc.service_account_file || '', placeholder: 'secrets/google-service-account.json', mono: true, hint: 'Path inside the container' })
                + fieldHtml({ id: 'nGmailFrom', label: 'From', value: cc.from || '', placeholder: 'alerts@yourdomain.com', mono: true, required: true })
                + `<div class="hm-field"><span class="hm-label">To</span>${tagInputHtml('nGmailTo', cc.to || [])}<span class="hm-hint">Enter to add</span></div>`
                + fieldHtml({ id: 'nGmailImp', label: 'Impersonate user (optional)', value: cc.impersonate_user || '', placeholder: 'user@yourworkspace.com', mono: true, hint: 'Google Workspace only' });
        } else { // webhook
            el.innerHTML = fieldHtml({ id: 'nWhUrl', label: 'URL', value: cc.url || '', placeholder: 'https://hooks.slack.com/services/…', mono: true, required: true })
                + `<div class="field-cols">
                    <div class="hm-field"><label class="hm-label" for="nWhMethod">Method</label>
                        <select class="hm-select" id="nWhMethod">${['POST', 'PUT', 'GET', 'PATCH'].map(m => `<option ${(cc.method || 'POST') === m ? 'selected' : ''}>${m}</option>`).join('')}</select>
                    <span class="hm-hint"></span></div>
                    ${fieldHtml({ id: 'nWhTimeout', label: 'Timeout', value: cc.timeout || '', placeholder: '10s', mono: true })}
                </div>`
                + `<div class="hm-field"><label class="hm-label" for="nWhHeaders">Headers (JSON, optional)</label>
                    <textarea class="hm-textarea hm-input--mono" id="nWhHeaders" rows="2" placeholder='{"Authorization": "Bearer …"}'>${esc(cc.headers ? JSON.stringify(cc.headers) : '')}</textarea>
                    <span class="hm-hint" data-hint-for="nWhHeaders">String → string object; Content-Type defaults to application/json</span></div>`
                + `<div class="hm-field"><label class="hm-label" for="nWhPayload">Payload template (optional)</label>
                    <textarea class="hm-textarea hm-input--mono" id="nWhPayload" rows="3" placeholder='{"text": "{{.Message}} ({{.TargetName}})"}'>${esc(cc.payload || '')}</textarea>
                    <span class="hm-hint">Go text/template over alert fields; empty = raw alert JSON</span></div>`
                + `<div class="field-cols">
                    ${fieldHtml({ id: 'nWhRetries', label: 'Max retries', value: cc.max_retries != null ? cc.max_retries : 3, mono: true, type: 'number' })}
                    ${fieldHtml({ id: 'nWhProxy', label: 'Proxy URL (optional)', value: cc.proxy_url || '', placeholder: 'socks5://host:1080', mono: true })}
                </div>`;
        }
        el.querySelectorAll('[data-taginput]').forEach(w => bindTagInput(panel, w.dataset.taginput));
    };

    renderPicker();
    typeFields();

    panel.querySelector('[data-cancel]').addEventListener('click', closePanel);
    panel.querySelector('[data-save]').addEventListener('click', async () => {
        const val = (id) => panel.querySelector('#' + id)?.value.trim() ?? '';
        const id = val('nId');
        if (!isEdit && !/^[a-z0-9-]+$/.test(id)) { setFieldError(panel, 'nId', 'Lowercase letters, numbers and hyphens only'); return; }

        let config = {};
        if (curType === 'telegram') {
            config = { chat_ids: tagValues(panel, 'nTgChats') };
            if (!config.chat_ids.length) { toast('warning', 'Chat IDs required', 'Add at least one chat ID'); return; }
            const token = val('nTgToken');
            if (token && token !== '***') config.bot_token = token;
            if (val('nTgProxy')) config.proxy_url = val('nTgProxy');
        } else if (curType === 'email') {
            if (!val('nSmtpHost')) { setFieldError(panel, 'nSmtpHost', 'SMTP host is required'); return; }
            config = {
                smtp_host: val('nSmtpHost'),
                smtp_port: parseInt(val('nSmtpPort'), 10) || 587,
                smtp_user: val('nSmtpUser'),
                from: val('nEmailFrom'),
                to: tagValues(panel, 'nEmailTo'),
                use_tls: panel.querySelector('#nSmtpTLS').checked,
            };
            if (val('nSmtpPass')) config.smtp_password = val('nSmtpPass');
        } else if (curType === 'gmail_oauth') {
            config = {
                credentials_file: val('nOauthCred') || 'secrets/gmail-token.json',
                token_file: val('nOauthToken') || 'secrets/gmail-oauth-token.json',
                from: val('nOauthFrom'),
                to: tagValues(panel, 'nOauthTo'),
            };
        } else if (curType === 'gmail') {
            config = {
                service_account_file: val('nGmailSA'),
                from: val('nGmailFrom'),
                to: tagValues(panel, 'nGmailTo'),
                impersonate_user: val('nGmailImp'),
            };
        } else {
            if (!val('nWhUrl')) { setFieldError(panel, 'nWhUrl', 'URL is required'); return; }
            config = { url: val('nWhUrl'), method: panel.querySelector('#nWhMethod').value };
            const headersRaw = val('nWhHeaders');
            if (headersRaw) {
                try { config.headers = JSON.parse(headersRaw); }
                catch (_) { setFieldError(panel, 'nWhHeaders', 'Headers must be valid JSON'); return; }
            }
            if (panel.querySelector('#nWhPayload').value.trim()) config.payload = panel.querySelector('#nWhPayload').value;
            if (val('nWhTimeout')) config.timeout = val('nWhTimeout');
            if (val('nWhRetries') !== '') config.max_retries = parseInt(val('nWhRetries'), 10) || 0;
            if (val('nWhProxy')) config.proxy_url = val('nWhProxy');
        }

        const payload = { id, type: curType, enabled: panel.querySelector('#nEnabled').checked, config };
        try {
            if (isEdit) await fetchAPI(`/api/v1/notifiers/${n.id}`, { method: 'PUT', body: JSON.stringify(payload) });
            else await fetchAPI('/api/v1/notifiers', { method: 'POST', body: JSON.stringify(payload) });
            closePanel();
            toast('success', 'Notifier saved', id);
            await loadNotifiers();
            renderAll();
        } catch (err) { toast('danger', 'Save failed', err.message); }
    });
}

/* ═══════════════ Live updates (SSE) + polling fallback ═══════════════ */

function setLive(state, label) {
    const el = document.getElementById('liveIndicator');
    el.dataset.state = state;
    document.getElementById('liveLabel').textContent = label;
}

function applyCheckEvent(d) {
    const entry = S.targets.find(e => e.t.id === d.target_id);
    if (!entry) { scheduleFullRefresh(); return; }
    entry.results.unshift({
        status: d.status,
        response_time_ms: d.response_time_ms,
        message: d.message,
        checked_at: d.checked_at || new Date().toISOString(),
    });
    if (entry.results.length > 45) entry.results.length = 45;
    Object.assign(entry, makeEntry(entry.t, entry.results));
    document.getElementById('lastUpdated').textContent = new Date().toLocaleTimeString();
    if (S.tab === 'dashboard') scheduleRender();
}

let fullRefreshDebounce = null;
function scheduleFullRefresh() {
    if (fullRefreshDebounce) return;
    fullRefreshDebounce = setTimeout(() => { fullRefreshDebounce = null; refreshData(); }, 800);
}

function scheduleIncidentsRefresh() {
    if (incidentsDebounce) return;
    incidentsDebounce = setTimeout(async () => {
        incidentsDebounce = null;
        try { await loadIncidents(); renderOngoing(); renderResolved(); renderHero(); } catch (_) {}
    }, 800);
}

function connectEventStream() {
    if (!window.EventSource) { setLive('polling', 'Polling (30s)'); return; }
    try {
        eventStream = new EventSource('/api/v1/events');
        eventStream.onopen = () => setLive('live', 'Live');
        eventStream.onerror = () => setLive('reconnecting', 'Reconnecting…');
        eventStream.addEventListener('check', (e) => {
            try { applyCheckEvent(JSON.parse(e.data)); } catch (_) { scheduleFullRefresh(); }
        });
        eventStream.addEventListener('alert', (e) => {
            try {
                const a = JSON.parse(e.data);
                const kind = a.type === 'up' ? 'success' : a.type === 'down' ? 'danger' : 'warning';
                toast(kind, a.target_name || a.target_id || 'Alert', a.message || a.type);
            } catch (_) {}
            scheduleIncidentsRefresh();
        });
    } catch (_) {
        setLive('polling', 'Polling (30s)');
    }
}

/* ═══════════════ Theme ═══════════════ */

const THEMES = ['dark', 'light', 'makima'];
const THEME_ICON = { dark: 'moon', light: 'sun', makima: 'activity' };

function applyTheme(theme) {
    document.documentElement.setAttribute('data-theme', theme);
    try { localStorage.setItem('hm-theme', theme); } catch (_) {}
    const next = THEMES[(THEMES.indexOf(theme) + 1) % THEMES.length];
    const btn = document.getElementById('themeBtn');
    btn.innerHTML = icon(THEME_ICON[theme], 15);
    btn.setAttribute('aria-label', `Theme: ${theme}. Switch to ${next}`);
    btn.title = `Theme: ${theme} → ${next}`;
    // re-tint the open chart, if any
    if (detail?.chart) renderDetailChart(overlayRoot().querySelector('.hm-slideover'));
}

/* ═══════════════ Init ═══════════════ */

function init() {
    document.getElementById('brandMark').innerHTML = icon('activity', 15);
    document.getElementById('addTargetIcon').innerHTML = icon('plus', 15);
    document.getElementById('searchIcon').innerHTML = icon('search', 14);
    document.getElementById('addNotifierBtn').innerHTML = `${icon('plus', 15)}Add notifier`;

    document.getElementById('addTargetBtn').addEventListener('click', () => openTargetForm(null));
    document.getElementById('addNotifierBtn').addEventListener('click', () => openNotifierForm(null));
    document.getElementById('searchInput').addEventListener('input', (e) => { S.query = e.target.value; renderTargets(); });
    document.getElementById('sortSelect').addEventListener('change', (e) => { S.sort = e.target.value; renderTargets(); });
    document.getElementById('themeBtn').addEventListener('click', () => {
        const cur = document.documentElement.getAttribute('data-theme') || 'dark';
        applyTheme(THEMES[(THEMES.indexOf(cur) + 1) % THEMES.length]);
    });

    // Show the logout button only when the server actually requires login.
    fetchAPI('/api/v1/auth/me').then((auth) => {
        if (!auth?.auth_enabled) return;
        const btn = document.getElementById('logoutBtn');
        btn.hidden = false;
        btn.innerHTML = icon('logout', 15);
        btn.addEventListener('click', async () => {
            await fetchAPI('/api/v1/auth/logout', { method: 'POST' }).catch(() => {});
            window.location.href = '/login';
        });
    }).catch(() => {});

    applyTheme(document.documentElement.getAttribute('data-theme') || 'dark');
    renderAll();
    refreshData();
    connectEventStream();
    setInterval(refreshData, 30000);
}

window.addEventListener('DOMContentLoaded', init);
window.addEventListener('beforeunload', () => { if (eventStream) eventStream.close(); });
