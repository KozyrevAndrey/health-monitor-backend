/** Response-time sparkline (line + soft area fill), self-scaling. */
export interface SparklineProps {
  /** Values in ms, oldest first */
  points: number[];
  width?: number;
  height?: number;
  stroke?: string;
  fill?: string;
}
