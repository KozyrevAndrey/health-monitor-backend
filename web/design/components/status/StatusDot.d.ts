/** Colored status dot; pulse reserved for ongoing DOWN states. */
export interface StatusDotProps {
  status: "up" | "down" | "degraded" | "unknown";
  /** "md" 8px (default) | "lg" 12px */
  size?: "md" | "lg";
  /** Radiating ring animation — only for live-down attention */
  pulse?: boolean;
}
