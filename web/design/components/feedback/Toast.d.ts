/** Toast notification (SSE alerts, CRUD results). Stack in .hm-toaststack, bottom-right. */
export interface ToastProps {
  kind: "success" | "danger" | "warning" | "info";
  title: string;
  message?: string;
  onDismiss?: () => void;
}
