# TokenRouter · Changelog

## v0.6.0-tokenrouter.1 (2026-08-20)

First branded release of TokenRouter. Based on upstream TokenHub v0.6.0 (Apache-2.0).

### Branding

- Renamed **TokenHub → TokenRouter 炬枢** across all UI copy (en/zh/ja), page titles, sidebar, sign-in page, breadcrumbs
- New brand mark: rounded hexagon badge + converging beams + amber core (SVG/PNG/ICO), replacing the old logo and favicon
- Brand lockup redesign: Chinese「炬枢」as the hero (21px sidebar primary, 36px sign-in lockup), English TokenRouter as a letter-spaced subtitle
- Version status moved to the sidebar footer, visible only to administrators (with upgrade access)

### Design & UX

- **Design system**: `design-system/tokenrouter/MASTER.md` (ui-ux-pro-max) — primary `#1E40AF`, background `#F8FAFC`, Fira Sans + Fira Code, data-dense enterprise dashboard spec
- **SaaS-style global upgrade**: ambient brand glow, sidebar gradient, gradient primary buttons, 16px rounded panels, slim scrollbars, keyboard focus rings, table row hover, themeable status pills
- **React Bits motion**: KPI count-up (CountUp), spotlight cards (SpotlightCard), gradient headings (GradientText), shimmer buttons (ShinyText), view transitions, animated nav active pill (motion layoutId)
- **Premium charts**: smooth Catmull-Rom trend line with gradient stroke, hover crosshair + tooltip, pulsing end-point halo; SVG donut (12 o'clock start, hover-linked legend); department token bar chart with hover detail; gradient progress bars for top models
- **Request log redesign**: single-line row meta (· separators), brand-blue selection state; detail panel scrolls within the viewport, 1px-divider metadata grid, side-by-side payload blocks
- Fixed stale-content navigation bug (AnimatePresence mode="wait" → keyed motion.div)
- Accessibility: `prefers-reduced-motion` fallbacks, visible focus

### Engineering & Demo

- `tools/mock-llm/server.js`: zero-dependency OpenAI-compatible mock upstream for lighting up the observability pages locally
- End-to-end data path verified: Provider → model import → routing → API key → `/v1` calls (incl. streaming) → usage / cost / audit
- Repository gates green: `next build` / `tsc` / `eslint` / 126 node tests / trilingual translation check / env contract / line-count baselines

## Upstream Baseline

- **TokenHub v0.6.0** (astaxie/TokenHub, Apache-2.0): enterprise AI gateway core — OpenAI/Anthropic-compatible APIs, provider channels and template catalog, model routing policies, project-scoped keys, usage/cost attribution, identity sources and RBAC, audit and alerts, SQLite/PostgreSQL.
