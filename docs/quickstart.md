# TokenRouter · Quick Start

Get a governed gateway running in 5 minutes: start → log in → connect a provider → call a model with an API key.

## Option 1: Docker Compose (recommended)

```bash
git clone <your-tokenrouter-repo-url> && cd tokenrouter
docker compose --env-file deploy/.env.example up -d
```

| Service | Address |
|---|---|
| Admin console | http://localhost:3000 |
| Gateway API | http://localhost:8080 |
| Health check | http://localhost:8080/healthz |

## Option 2: Local development

Requirements: Go 1.26+, Node.js 20+.

```bash
# 1. Backend (SQLite, default port 8080)
cd backend
go build -o tokenrouter ./cmd/tokenhub
TOKENHUB_ENV=dev ./tokenrouter

# 2. Frontend (new terminal, port 3000)
cd frontend
npm install
TOKENHUB_API_BASE_URL=http://localhost:8080 npm run dev
```

> The frontend targets `http://localhost:8080` by default; point `TOKENHUB_API_BASE_URL` elsewhere when the backend port changes.

## Initialize and sign in

1. Open http://localhost:3000
2. Sign in with the bootstrap administrator: **admin / admin123456** (change the password immediately after first sign-in)
3. Land on the Overview dashboard

> The bootstrap admin password can be preset with the `TOKENHUB_BOOTSTRAP_ADMIN_PASSWORD` environment variable.

## Connect your first provider

1. Sidebar → **AI Resources → Provider Channels** → **Add Provider**
2. Pick a template from the 150+ catalog, or choose "Custom OpenAI-compatible":
   - Name: e.g. `DeepSeek Production`
   - Base URL: `https://api.deepseek.com/v1`
   - API Key: your upstream credential
3. Save, then **import models** (the upstream model list is fetched automatically; select the ones you need)
4. Go to **Routing Policies** and create a route for the model (select channel + weight/priority), then save

## Create an API key and call the API

1. **Organization → Project Spaces** → open your project → **Key Management → Add Key**
2. Copy the generated API key (shown only once)

```bash
curl http://localhost:8080/v1/chat/completions \
  -H "Authorization: Bearer <YOUR_API_KEY>" \
  -H "content-type: application/json" \
  -d '{
    "model": "<model-enabled-in-routing>",
    "messages": [{"role": "user", "content": "Hello"}]
  }'
```

A successful call returns the OpenAI-shaped response. Afterwards:

- **Request Logs** show the full path of the call (route, upstream, tokens, cost)
- **Usage Analytics** attribute usage to project / model / cost center
- **Cost Dashboard** aggregates spend in real time

## Next steps

| Goal | Doc |
|---|---|
| Full deployment (PostgreSQL, reverse proxy, env vars) | [Deployment Guide](deployment.md) |
| Role-based manuals | [User Guide](user-guide.md) / [Team Leader Guide](team-leader-guide.md) / [Administrator Guide](administrator-guide.md) |
| Connect Codex subscriptions | [Codex Quick Setup](codex-tokenhub-profile-quick-start.md) |
| Architecture and data flow | [Architecture](architecture.md) |
