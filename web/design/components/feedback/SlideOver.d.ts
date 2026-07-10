/** Right slide-over panel for target detail & forms (roomier than a modal). */
export interface SlideOverProps {
  title: React.ReactNode;
  /** Extra header content (e.g. status badge, period tabs) */
  titleExtra?: React.ReactNode;
  /** Panel width in px (default 560; use ~680 for detail w/ chart) */
  width?: number;
  onClose?: () => void;
  /** Sticky footer actions (buttons) */
  footer?: React.ReactNode;
  children?: React.ReactNode;
}
