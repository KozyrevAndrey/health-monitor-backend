import React from "react";

const P = {
  check: ["M20 6 9 17l-5-5"],
  x: ["M18 6 6 18", "M6 6l12 12"],
  plus: ["M12 5v14", "M5 12h14"],
  search: ["M21 21l-4.35-4.35", "M11 4a7 7 0 1 1 0 14 7 7 0 0 1 0-14z"],
  "chevron-down": ["m6 9 6 6 6-6"],
  "chevron-right": ["m9 6 6 6-6 6"],
  "chevron-left": ["m15 6-6 6 6 6"],
  "arrow-left": ["M19 12H5", "m12 19-7-7 7-7"],
  clock: ["M12 3a9 9 0 1 1 0 18 9 9 0 0 1 0-18z", "M12 7v5l3 2"],
  trash: ["M4 7h16", "M9 7V4h6v3", "M6 7l1 13h10l1-13"],
  pencil: ["M17 3l4 4L8 20l-5 1 1-5L17 3"],
  "alert-triangle": ["M12 3 2 20h20L12 3z", "M12 10v4", "M12 17h.01"],
  "check-circle": ["M12 3a9 9 0 1 1 0 18 9 9 0 0 1 0-18z", "m8.5 12.5 2.5 2.5 4.5-5"],
  "x-circle": ["M12 3a9 9 0 1 1 0 18 9 9 0 0 1 0-18z", "M15 9l-6 6", "M9 9l6 6"],
  info: ["M12 3a9 9 0 1 1 0 18 9 9 0 0 1 0-18z", "M12 16v-5", "M12 8h.01"],
  activity: ["M22 12h-4l-3 9L9 3l-3 9H2"],
  globe: ["M12 3a9 9 0 1 1 0 18 9 9 0 0 1 0-18z", "M3 12h18", "M12 3c2.5 2.5 3.8 5.6 3.8 9S14.5 18.5 12 21c-2.5-2.5-3.8-5.6-3.8-9S9.5 5.5 12 3z"],
  server: ["M4 4h16a2 2 0 0 1 2 2v3a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2z", "M4 13h16a2 2 0 0 1 2 2v3a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2v-3a2 2 0 0 1 2-2z", "M6 7.5h.01", "M6 16.5h.01"],
  "at-sign": ["M12 8a4 4 0 1 0 0 8 4 4 0 0 0 0-8z", "M16 8v5a3 3 0 0 0 6 0v-1a10 10 0 1 0-3.9 7.9"],
  bell: ["M6 9a6 6 0 0 1 12 0c0 5 2 6 2 6H4s2-1 2-6", "M10 19a2 2 0 0 0 4 0"],
  sliders: ["M4 8h10", "M18 8h2", "M4 16h4", "M12 16h8", "M16 8h.01M14 6a2 2 0 1 0 4 0 2 2 0 0 0-4 0z", "M8 16a2 2 0 1 0 4 0 2 2 0 0 0-4 0z"],
  zap: ["M13 2 3 14h7l-1 8 11-12h-7l1-8z"],
  pause: ["M9 5v14", "M15 5v14"],
  send: ["M22 2 11 13", "M22 2 15 22l-4-9-9-4z"],
  mail: ["M4 5h16a1 1 0 0 1 1 1v12a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1V6a1 1 0 0 1 1-1z", "m3 7 9 6 9-6"],
  link: ["M10 14a5 5 0 0 1 0-7l1.5-1.5a5 5 0 0 1 7 7L17.5 13.5", "M14 10a5 5 0 0 1 0 7l-1.5 1.5a5 5 0 0 1-7-7L6.5 10.5"],
  moon: ["M21 13A9 9 0 1 1 11 3a7 7 0 0 0 10 10z"],
  sun: ["M12 8a4 4 0 1 1 0 8 4 4 0 0 1 0-8z", "M12 2v2", "M12 20v2", "M2 12h2", "M20 12h2", "m5 5 1.4 1.4", "m17.6 17.6 1.4 1.4", "m19 5-1.4 1.4", "m6.4 17.6-1.4 1.4"],
  grid: ["M4 4h6v6H4z", "M14 4h6v6h-6z", "M4 14h6v6H4z", "M14 14h6v6h-6z"],
  list: ["M8 6h13", "M8 12h13", "M8 18h13", "M3 6h.01", "M3 12h.01", "M3 18h.01"],
  "more-h": ["M5 12h.01", "M12 12h.01", "M19 12h.01"],
  external: ["M14 3h7v7", "M21 3l-9 9", "M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h6"],
  refresh: ["M21 12a9 9 0 1 1-2.9-6.6L21 8", "M21 3v5h-5"],
  wifi: ["M5 12a11 11 0 0 1 14 0", "M8.5 15.5a6 6 0 0 1 7 0", "M12 19h.01"],
  filter: ["M3 5h18l-7 8v6l-4-2v-4L3 5z"],
};

/** Inline stroke icon (Lucide-style: 24 grid, 2px stroke, round caps). */
export function Icon({ name, size = 16, strokeWidth = 2, style, ...rest }) {
  const paths = P[name] || P["info"];
  return (
    <svg
      width={size} height={size} viewBox="0 0 24 24" fill="none"
      stroke="currentColor" strokeWidth={strokeWidth} strokeLinecap="round" strokeLinejoin="round"
      aria-hidden={rest["aria-label"] ? undefined : true}
      style={{ flex: "none", ...style }} {...rest}
    >
      {paths.map((d, i) => <path key={i} d={d} />)}
    </svg>
  );
}

Icon.names = Object.keys(P);
