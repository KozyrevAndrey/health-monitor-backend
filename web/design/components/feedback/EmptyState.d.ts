/** Empty state (no targets yet, no incidents, no search results). */
export interface EmptyStateProps {
  /** Icon name (default "activity") */
  icon?: string;
  title: string;
  body?: string;
  /** Primary action, e.g. <Button variant="primary">Add your first target</Button> */
  action?: React.ReactNode;
}
