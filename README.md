# Pantry

Household shopping list and pantry app. React + TypeScript frontend, Go API, Postgres.

## Requirements

- Go 1.27+
- Node 22+ and pnpm
- Docker (Colima locally)

Go tools (sqlc, goose, oapi-codegen, air) are pinned in `backend/go.mod` via the `tool`
directive and run with `go tool <name>`. No global installs needed.

## Setup

```sh
make setup     # copies .env.example to .env, installs frontend deps, downloads Go modules
make db-up     # Postgres: dev on :5432, test on :5433 (tmpfs, wiped on restart)
make migrate   # applies migrations to both dev and test databases
make gen       # regenerates sqlc, oapi-codegen and TypeScript API types
```

## Development

```sh
make dev       # Go API on :8080 (live reload via air) + Vite on :5180
```

Open http://localhost:5180. Vite proxies `/api/*` to the Go server, so there is no CORS
setup in development.

## Commands

| Command | What it does |
| --- | --- |
| `make help` | List all targets |
| `make db-up` / `make db-down` | Start / stop Postgres |
| `make migrate` / `make migrate-down` | Apply migrations (dev + test) / roll back last (dev) |
| `make migrate-status` | Show applied migrations on dev DB |
| `make migrate-new name=<name>` | Create a new goose SQL migration |
| `make gen` | sqlc + oapi-codegen + openapi-typescript |
| `make seed` | Load realistic fake data |
| `make dev` | Backend with reload + Vite dev server |
| `make test` | Go and frontend tests (store tests need `make db-up` and `make migrate` first) |
| `make lint` | gofmt, go vet, tsc, eslint |

## Configuration

The backend reads environment variables once at startup and exits if any are invalid.
`make` loads them from `.env`.

| Variable | Required | Default | Notes |
| --- | --- | --- | --- |
| `DATABASE_URL` | yes | | Postgres connection string |
| `HTTP_ADDR` | no | `:8080` | Listen address |
| `LOG_FORMAT` | no | `json` | `json` or `text` |
| `TEST_DATABASE_URL` | for store tests | | Points at the `db-test` container |

## Testing

Store tests run every query against the `db-test` container. Each test runs inside a
transaction that is rolled back at the end, so tests never see each other's data and the
database stays empty. Go's test cache is off (`-count=1`) because it cannot see database
changes.

## Layout

- `api/openapi.yaml`: HTTP API contract, source of truth
- `backend/`: Go API (`cmd/server`, `internal/*`, `db/migrations`, `db/queries`)
- `frontend/`: Vite + React + TS (`src/api`, `src/routes`, `src/features`, `src/styles`)
- `docs/DECISIONS.md`: architecture and tradeoff decisions
