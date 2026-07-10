/** Monitored-target card: dot + name + type icon, 45-check tick bar, uptime % + last ms. */
export interface TargetCardProps {
  target: {
    name: string;
    type: "http" | "tcp" | "dns";
    status: "up" | "down" | "degraded" | "unknown";
    /** Recent checks for the tick bar, oldest first */
    checks: Array<{ s: string; tip?: string }>;
    /** Uptime percentage for the active period, e.g. 99.98 */
    uptime?: number | string;
    /** Last response time in ms */
    lastMs?: number;
    /** URL / host:port / domain, shown in monospace */
    endpoint?: string;
    paused?: boolean;
  };
  /** Open detail view (whole card is clickable) */
  onOpen?: () => void;
  onEdit?: (e: any) => void;
  onDelete?: (e: any) => void;
}
