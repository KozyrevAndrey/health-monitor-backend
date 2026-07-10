---
name: health-monitor-design
description: Use this skill to generate well-branded interfaces and assets for Health Monitor (self-hosted uptime monitoring dashboard), either for production or throwaway prototypes/mocks/etc. Contains essential design guidelines, colors, type, fonts, assets, and UI kit components for prototyping.
user-invocable: true
---

Read the README.md file within this skill, and explore the other available files.
If creating visual artifacts (slides, mocks, throwaway prototypes, etc), copy assets out and create static HTML files for the user to view. If working on production code, you can copy assets and read the rules here to become an expert in designing with this brand.
If the user invokes this skill without any other guidance, ask them what they want to build or design, ask some questions, and act as an expert designer who outputs HTML artifacts _or_ production code, depending on the need.

Key facts: dark-first (cool blue-black, `#0d1117` page), semantic status palette (up/down/degraded/unknown) does the talking, blue accent for interaction only, Inter + JetBrains Mono (tabular numerals for all metrics), 4px spacing scale, no emoji, inline SVG stroke icons. The vanilla CSS contract lives in `tokens/components.css` (`hm-*` classes) — the whole UI must remain buildable as a static page with vanilla JS.
