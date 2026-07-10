import React from "react";

/** Shimmering skeleton block. Compose several to sketch a loading layout. */
export function Skeleton({ width = "100%", height = 14, radius, style }) {
  return (
    <span
      className="hm-skel"
      aria-hidden="true"
      style={{ display: "block", width, height, borderRadius: radius, ...style }}
    ></span>
  );
}
