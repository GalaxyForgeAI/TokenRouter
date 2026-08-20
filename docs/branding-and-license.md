# TokenRouter · Brand & License

## Brand

| Item | Value |
|---|---|
| Product name (EN) | **TokenRouter** |
| Product name (ZH) | **炬枢** |
| Tagline | 炬枢 · Enterprise AI Gateway |
| Brand mark | Rounded hexagon badge + converging beams + amber core (`frontend/public/brand/tokenrouter-logo.svg`) |

TokenRouter and 炬枢 are product-specific brands. Do not reuse the brand mark in derivative projects without permission.

## Open-Source License

This project is released under the **Apache License 2.0** (see `LICENSE` at the repository root).

- ✅ Commercial use, modification, redistribution, and patent grants
- ✅ Closed-source commercial distribution
- ✅ Branded derivative products on top of this codebase
- 📌 Obligations: retain copyright and license notices, state changes, preserve NOTICE (if any)

## Upstream Credits

This project is derived from [astaxie/TokenHub](https://github.com/astaxie/TokenHub) (Apache-2.0), maintained by Asta Xie and the community.

**Scope of the derivation**:

| Layer | Description |
|---|---|
| Brand layer | Product name, logo, mark, and UI copy (en/zh/ja) replaced with TokenRouter 炬枢 |
| Design layer | Design system (`design-system/tokenrouter/MASTER.md`), SaaS-style UI, React Bits motion, interactive chart upgrades |
| Functional layer | Gateway capabilities and API contracts kept in sync with upstream (OpenAI-compatible `/v1`, admin APIs, env var system) |

**Compliance notes**:

- The upstream Apache-2.0 notice is preserved in `LICENSE` and file headers
- Internal contracts (env vars, CLI commands, storage keys) keep upstream naming (`TOKENHUB_*`) so upstream docs, tooling, and migration paths remain compatible
- Upstream GitHub references (release checks, upgrade scripts) keep pointing at the original repository for upstream tracking

## Trademark Statement

TokenRouter / 炬枢 has **no affiliation with, and is not endorsed by**, astaxie/TokenHub or its authors. TokenHub is the project name of its original author and belongs to its respective owners. The brand of this project represents only this project itself.

## Contributing

- Coding style and repository gates follow `AGENTS.md` (English commit messages, trilingual copy sync, gate checks)
- Trilingual documentation sync: Simplified Chinese is the canonical UI copy source; English and Japanese translations align with it
