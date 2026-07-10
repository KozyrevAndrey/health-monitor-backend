/** Tick bar of recent check results (oldest→newest), with hover tooltips. */
export interface TickBarProps {
  /** Recent checks, oldest first. tip renders as a hover tooltip (monospace, supports \n). */
  checks: Array<{ s: "up" | "down" | "degraded" | "unknown"; tip?: string }>;
  /** Fixed slot count; short histories pad left with empty ticks (default 45) */
  slots?: number;
  /** Bar height in px (default 20) */
  height?: number;
  ariaLabel?: string;
}
