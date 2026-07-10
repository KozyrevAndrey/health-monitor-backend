/** Confirm dialog replacing window.confirm(); Esc closes, focus starts on Cancel. */
export interface ConfirmDialogProps {
  title: string;
  body: string;
  confirmLabel?: string;
  cancelLabel?: string;
  /** Destructive styling on the confirm button (default true) */
  danger?: boolean;
  onConfirm?: () => void;
  onCancel?: () => void;
}
