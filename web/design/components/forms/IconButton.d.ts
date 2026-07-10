/** Square icon-only button; label is mandatory (aria-label + hover tooltip). */
export interface IconButtonProps {
  /** Icon name from the Icon set */
  name: string;
  /** Accessible label, e.g. "Edit target" */
  label: string;
  /** Red hover treatment for destructive actions */
  danger?: boolean;
  onClick?: () => void;
}
