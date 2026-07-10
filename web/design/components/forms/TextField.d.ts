/** Labeled input with hint & inline error states. */
export interface TextFieldProps {
  label?: string;
  hint?: string;
  /** Error message; presence switches to invalid styling */
  error?: string;
  /** Monospace — use for URLs, host:port, domains */
  mono?: boolean;
  placeholder?: string;
  value?: string;
  onChange?: (e: any) => void;
  type?: string;
}
