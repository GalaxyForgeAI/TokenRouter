# TokenRouter · Product Overview

> TokenRouter (炬枢) — the private, self-hosted AI gateway for enterprises.
> One governed entry point for every model call: **controlled, traceable, attributable**.

## Why TokenRouter

Teams integrating AI quickly hit four problems:

| Problem | Symptom |
|---|---|
| **Scattered access** | Each team wires OpenAI / Anthropic / Gemini / domestic models separately; keys leak everywhere |
| **No governance** | No unified policy for who can call what, how much a project may spend, or how hard upstream may be hit |
| **Cost blind spots** | Month-end bills cannot be attributed to team, project, model, or cost center |
| **Missing audit trail** | After an incident, no answer to "who, when, with which key, called which model, with what payload" |

TokenRouter solves all four with one self-hosted gateway: **unified access → policy governance → usage attribution → full audit**.

## Core Capabilities

### Unified Access

- OpenAI-compatible APIs: `/v1/chat/completions`, `/v1/responses`, `/v1/embeddings`, `/v1/images/*`
- Anthropic Messages API: `/v1/messages`, `/v1/messages/count_tokens`
- Native adapters: OpenAI, Azure OpenAI, Anthropic, Gemini, DeepSeek, Qwen, Codex subscriptions, local vLLM / Ollama
- **150+ provider templates** plus custom OpenAI-compatible upstreams

### Policy Governance

- Model catalog + routing policies: priority, weight, failover order, route health diagnostics
- Project-scoped API keys: team ownership, member permissions, quotas, concurrency controls
- Approval flows for quota changes and sensitive operations
- Security policies: content safety, proxy egress, identity sources (OAuth/OIDC), RBAC, audit trails

### Usage & Cost

- Attribution across **user, project, team, model, and cost center**
- Request logs with token detail (input / output / cache / reasoning) and cost detail
- Executive dashboard: department token comparison, top models, cost share, trends

### Observability & Security

- Health checks, alert rules, notification channels
- Database backup (SQLite / PostgreSQL) and status monitoring
- Administrative audit events (who changed what configuration)

## Three Roles

| Role | Focus | Primary work |
|---|---|---|
| **User** | Using models | Find available models, create project keys, call the model API, review personal usage |
| **Team Leader** | Managing projects | Project spaces, members, keys, team reports, cost attribution |
| **Administrator** | Running the platform | Providers, model catalog, routing, identity sources, RBAC, audit, cost controls |

## Architecture

```
┌─────────────────────────────┐
│  Admin Console (Next.js 16)  │  role-aware nav, global search, dual theme, API docs
└──────────┬──────────────────┘
           │ /api/admin/* (JWT)
┌──────────▼──────────────────┐
│  Gateway Core (Go 1.26)      │  auth, routing, quotas, audit, usage recording
│  /v1/chat/completions etc.   │
└──────────┬──────────────────┘
           │ routing (priority / weight / failover)
┌──────────▼──────────────────┐
│  Provider Channels          │  OpenAI / Anthropic / Gemini / domestic / local
└─────────────────────────────┘
```

- **Backend**: Go 1.26, SQLite (default) / PostgreSQL
- **Frontend**: Next.js 16 + React 19, trilingual (en/zh/ja)
- **Deployment**: Docker Compose, reverse-proxy ready

## Get Started

- [Quick Start](quickstart.md) — up and running in 5 minutes
- [Brand & License](branding-and-license.md) — open-source compliance and upstream credits
- [Changelog](changelog.md) — release history
