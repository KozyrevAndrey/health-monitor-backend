/**
 * Inline SVG stroke icon (Lucide-style). Always paired with text for status
 * meaning — never the sole carrier of state.
 */
export interface IconProps {
  /** Icon name, e.g. "check", "alert-triangle", "globe", "server", "at-sign", "bell", "trash", "pencil", "search", "plus", "zap", "pause", "send", "mail", "link", "moon", "sun", "grid", "list", "clock", "activity", "sliders", "refresh", "wifi", "filter", "x", "info", "check-circle", "x-circle", "chevron-down", "chevron-right", "chevron-left", "arrow-left", "more-h", "external" */
  name: string;
  /** Pixel size (default 16) */
  size?: number;
  /** Stroke width on the 24px grid (default 2) */
  strokeWidth?: number;
}
