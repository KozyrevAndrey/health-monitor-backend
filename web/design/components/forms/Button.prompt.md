Form controls: Button, IconButton, TextField, SelectField, Switch, Segmented, TagInput, DurationField.

```jsx
<Button variant="primary" icon={<Icon name="plus" size={15} />}>Add target</Button>
<IconButton name="trash" label="Delete target" danger />
<TextField label="URL" mono placeholder="https://example.com/health" error="Must start with http(s)://" />
<Switch checked label="Enabled" />
<Segmented options={["24h", "7d", "30d"]} value="24h" />
<DurationField label="Check interval" value="1m" />
<TagInput value={["prod", "api"]} />
```

Forms group into sections with `.hm-formsection` + `.overline` headings (Identity → Check settings → Advanced). Inline errors under fields, never `alert()`.
