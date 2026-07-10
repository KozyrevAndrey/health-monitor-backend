/** Status pill with icon + text — status is never conveyed by color alone. */
export interface StatusBadgeProps {
  status: "up" | "down" | "degraded" | "unknown";
  /** Override the default text (Up / Down / Degraded / Unknown) */
  label?: string;
}
