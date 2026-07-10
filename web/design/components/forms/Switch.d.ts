/** Toggle switch; checked track uses --status-up (enabled = healthy green). */
export interface SwitchProps {
  checked: boolean;
  onChange?: (checked: boolean) => void;
  label?: string;
  disabled?: boolean;
}
