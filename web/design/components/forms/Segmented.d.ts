/** Segmented control / tabs: period switcher (24h/7d/30d), card↔table toggle, status filters. */
export interface SegmentedProps {
  options: Array<string | { value: string; label: string; icon?: React.ReactNode; count?: number }>;
  value: string;
  onChange?: (value: string) => void;
  ariaLabel?: string;
}
