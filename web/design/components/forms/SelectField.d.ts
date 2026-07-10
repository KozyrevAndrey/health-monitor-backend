/** Labeled native select with custom chevron. */
export interface SelectFieldProps {
  label?: string;
  hint?: string;
  options: Array<string | { value: string; label: string }>;
  value?: string;
  onChange?: (e: any) => void;
}
