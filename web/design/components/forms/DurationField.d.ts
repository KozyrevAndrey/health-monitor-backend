/** Duration input: preset chips (30s/1m/5m/…) + free-form Go-duration field. */
export interface DurationFieldProps {
  label?: string;
  /** Go-style duration string, e.g. "30s", "5m" */
  value: string;
  onChange?: (value: string) => void;
  presets?: string[];
  hint?: string;
}
