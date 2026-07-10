/** Incident row with timeline rail; ongoing = pulsing red, resolved = muted. */
export interface IncidentItemProps {
  incident: {
    target: string;
    status: "ongoing" | "resolved";
    severity: "critical" | "warning" | "info";
    /** Human string, e.g. "Started 12 Jun, 09:14" */
    startedAt: string;
    /** Human duration, e.g. "23m" or "1h 05m" */
    duration: string;
    failureCount?: number;
    lastError?: string;
  };
  /** Hide the connecting rail line (last item in list) */
  last?: boolean;
  onOpen?: () => void;
}
