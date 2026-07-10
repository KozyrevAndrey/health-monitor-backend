/** SSE connection indicator (live / reconnecting / polling fallback). */
export interface LiveIndicatorProps {
  state: "live" | "reconnecting" | "polling";
  label?: string;
}
