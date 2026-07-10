/** Dense table view of targets — same data as TargetCard, one row each. */
export interface TargetTableProps {
  targets: Array<{
    id?: string;
    name: string;
    type: "http" | "tcp" | "dns";
    status: "up" | "down" | "degraded" | "unknown";
    checks: Array<{ s: string; tip?: string }>;
    uptime?: number | string;
    lastMs?: number;
    paused?: boolean;
  }>;
  onOpen?: (target: any) => void;
  onEdit?: (target: any) => void;
  onDelete?: (target: any) => void;
}
