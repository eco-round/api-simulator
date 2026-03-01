# EcoRound — API Simulator

Go/Gin REST API that simulates three esports data sources (PandaScore, VLR, Liquipedia) from a single PostgreSQL database. Used by the Chainlink CRE oracle for match consensus.

## Role in the System

```
CRE Workflow (oracle)
  ├── GET /api/v1/admin/matches          → list all open/locked matches
  ├── GET /api/v1/pandascore/matches/:id → source 1 result
  ├── GET /api/v1/vlr/matches/:id        → source 2 result
  └── GET /api/v1/liquipedia/matches/:id → source 3 result

Admin Panel / Panel-v2
  ├── POST   /api/v1/admin/matches        → create match
  ├── PATCH  /api/v1/admin/matches/:id   → update status
  ├── PATCH  /api/v1/admin/matches/:id/vault → save vault address
  └── POST   /api/v1/admin/matches/:id/result → set source result (simulate reporting)
```

## API Endpoints

### Admin (no auth)

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/v1/admin/matches` | Create a match (with optional `vault_address`, `on_chain_match_id`) |
| `GET` | `/api/v1/admin/matches` | List all matches (filter: `?status=open`) |
| `GET` | `/api/v1/admin/matches/:id` | Get match with results |
| `PATCH` | `/api/v1/admin/matches/:id` | Update match status |
| `PATCH` | `/api/v1/admin/matches/:id/vault` | Save `vault_address` + `on_chain_match_id` |
| `POST` | `/api/v1/admin/matches/:id/result` | Set a source result (simulate data source reporting) |

### Data Sources (protected by `X-Api-Key` header if keys are set)

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/v1/pandascore/matches/:id` | PandaScore result for match |
| `GET` | `/api/v1/vlr/matches/:id` | VLR result for match |
| `GET` | `/api/v1/liquipedia/matches/:id` | Liquipedia result for match |

> If `PANDASCORE_API_KEY` env var is empty, auth is skipped (open dev mode).
> For CRE confidential HTTP, set keys matching `secrets.yaml` in the `cre/` directory.

## Setup

```bash
cd api-simulator

# Configure environment
cp .env.example .env
# Edit .env:
#   DATABASE_URL=postgresql://...
#   PANDASCORE_API_KEY=pandascore123ECO
#   VLR_API_KEY=vlr123ECO
#   LIQUIPEDIA_API_KEY=liquipedia123ECO

# Run
go run .
# Server starts on :8080
```

## Environment Variables

| Variable | Description |
|---|---|
| `DATABASE_URL` | Neon PostgreSQL connection string |
| `PANDASCORE_API_KEY` | API key for PandaScore endpoint (empty = open access) |
| `VLR_API_KEY` | API key for VLR endpoint |
| `LIQUIPEDIA_API_KEY` | API key for Liquipedia endpoint |

## Database Models

**Match**: `id`, `on_chain_match_id`, `vault_address`, `team_a_name`, `team_a_tag`, `team_b_name`, `team_b_tag`, `status`, `best_of`, `event`, `start_time`

**MatchResult**: `match_id`, `source`, `match_status` (upcoming/started/ended), `winner`, `score_a`, `score_b`, `map_count`

## Simulating a Match Flow

```bash
# 1. Create a match (done via panel-v2, which also deploys the vault)

# 2. Simulate sources reporting "started" (triggers lockMatch via oracle)
curl -X POST http://localhost:8080/api/v1/admin/matches/1/result \
  -H "Content-Type: application/json" \
  -d '{"source":"pandascore","match_status":"started"}'

# 3. Simulate sources reporting "ended" with winner
curl -X POST http://localhost:8080/api/v1/admin/matches/1/result \
  -H "Content-Type: application/json" \
  -d '{"source":"pandascore","match_status":"ended","winner":"TeamA","score_a":2,"score_b":0}'
```
