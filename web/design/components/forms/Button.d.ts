/** Button. default = subtle surface; primary = accent; danger = destructive outline; ghost = borderless. */
export interface ButtonProps {
  variant?: "default" | "primary" | "danger" | "ghost";
  /** md 32px | lg 44px (touch target) */
  size?: "md" | "lg";
  /** Leading icon element, e.g. <Icon name="plus" size={15} /> */
  icon?: React.ReactNode;
  disabled?: boolean;
  onClick?: () => void;
  children?: React.ReactNode;
}
