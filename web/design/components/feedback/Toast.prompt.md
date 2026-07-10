Feedback & overlays: Toast, ConfirmDialog, SlideOver, EmptyState, Skeleton.

```jsx
<div className="hm-toaststack">
  <Toast kind="danger" title="API is down" message="HTTP 503 · 3 consecutive failures" onDismiss={…} />
</div>
<ConfirmDialog title="Delete target?" body="“API /health” and its check history will be removed." onConfirm={…} onCancel={…} />
<SlideOver title="API /health" width={680} onClose={…} footer={<Button variant="primary">Save</Button>}>…</SlideOver>
<EmptyState icon="activity" title="No targets yet" body="Add your first HTTP, TCP or DNS check." action={<Button variant="primary">Add target</Button>} />
<Skeleton height={18} width="60%" />
```

Rules: never `window.confirm`/`alert`; toasts bottom-right; slide-over for detail + forms; Esc always closes overlays.
