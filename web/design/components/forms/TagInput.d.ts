/** Tag chips input (Enter/comma to add, Backspace to remove). */
export interface TagInputProps {
  value: string[];
  onChange?: (tags: string[]) => void;
  placeholder?: string;
}
