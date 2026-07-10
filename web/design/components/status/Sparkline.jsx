import React from "react";

/**
 * Minimal response-time sparkline (SVG line + soft area fill).
 * `points` is an array of numbers (ms); scales to its own min/max.
 */
export function Sparkline({ points = [], width = 120, height = 28, stroke = "var(--chart-line)", fill = "var(--chart-fill)", style }) {
  if (points.length < 2) return <svg width={width} height={height} style={style} aria-hidden="true"></svg>;
  const min = Math.min(...points);
  const max = Math.max(...points);
  const span = max - min || 1;
  const step = width / (points.length - 1);
  const y = (v) => 2 + (height - 4) * (1 - (v - min) / span);
  const pts = points.map((v, i) => `${(i * step).toFixed(1)},${y(v).toFixed(1)}`);
  return (
    <svg width={width} height={height} style={style} role="img" aria-label={`Response time trend, ${min}–${max} ms`}>
      <polygon points={`0,${height} ${pts.join(" ")} ${width},${height}`} fill={fill} stroke="none" />
      <polyline points={pts.join(" ")} fill="none" stroke={stroke} strokeWidth="1.5" strokeLinejoin="round" />
    </svg>
  );
}
