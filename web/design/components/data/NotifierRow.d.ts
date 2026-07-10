/** Notifier row for the settings area: type icon, name, masked config detail, enable switch. */
export interface NotifierRowProps {
  notifier: {
    name: string;
    type: "telegram" | "email" | "gmail" | "gmail_oauth" | "webhook";
    enabled: boolean;
    /** Masked config summary, e.g. "chat •••4821" or "ops@•••.dev" */
    detail?: string;
  };
  onToggle?: (enabled: boolean) => void;
  onEdit?: () => void;
  onDelete?: () => void;
}
