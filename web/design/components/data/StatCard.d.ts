/** Compact metric card (label / big tabular number / optional sub-line). */
export interface StatCardProps {
  label: string;
  value: string | number;
  /** Small unit suffix, e.g. "ms", "%" */
  unit?: string;
  /** Secondary line under the value */
  sub?: string;
  /** Tint the value by meaning */
  tone?: "default" | "success" | "warning" | "danger";
  icon?: React.ReactNode;
}
