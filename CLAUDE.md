# Pantry

Household shopping list and pantry app. Rewrite of an older version, built agentically as a
learning project. Target stack mirrors a modern product team: React + TypeScript frontend,
Go API, Postgres.

## Who you're working with

- Senior web developer: strong in PHP/Laravel/Rails, JS, SQL, API design, DevOps.
- **New to Go.** Treat Go as something I'm learning, not something I already know.
- Fairly new to large-scale React/TS patterns (TanStack Query, Router, optimistic updates).
- I care more about understanding the code than shipping fast. A smaller diff I understand
  beats a big one I don't.

## Ownership rules (most important section)

- **I own the API contract** (`api/openapi.yaml`) **and the database schema**
  (`backend/db/migrations/`). You may propose changes, but never edit either without
  explicitly asking first and showing the diff.
- No new endpoint, field, or table exists until it's in the contract or a migration.
  Code follows the contract, never the other way around.
- Tradeoff decisions (conflict resolution, caching strategy, indexes, auth approach) are mine.
  Lay out the options with pros/cons and a recommendation, then wait for my call.
  Once decided, append it to `docs/DECISIONS.md`.

## Learning mode

- Before using a Go idiom or pattern for the first time in this project (interfaces,
  `context`, goroutines/channels, embedding, generics, `defer` subtleties, error wrapping),
  explain it in 2-4 sentences, compared to how I'd do it in PHP/Laravel where it helps.
- After finishing a task, add a short "Go concepts used" list at the end of your summary.
- Don't hide complexity behind libraries I haven't agreed to. Ask before adding any dependency.

## Repo layout

```
pantry/
├── api/openapi.yaml           # source of truth for the HTTP API
├── backend/
│   ├── cmd/server/main.go     # entrypoint: config, DB pool, router, graceful shutdown
│   ├── internal/
│   │   ├── api/               # oapi-codegen generated interfaces + handler impls
│   │   ├── store/             # sqlc generated code (do not edit by hand)
│   │   ├── service/           # business logic, only when handlers get non-trivial
│   │   └── config/
│   ├── db/
│   │   ├── migrations/        # goose SQL migrations
│   │   └── queries/           # hand-written SQL for sqlc
│   └── sqlc.yaml
├── frontend/                  # Vite + React + TS
│   └── src/
│       ├── api/               # generated types + thin fetch client
│       ├── routes/            # TanStack Router
│       ├── features/<name>/   # components, hooks, queries per feature
│       └── styles/
├── docs/DECISIONS.md
├── compose.yaml               # Postgres for local dev (runs on Colima)
└── Makefile
```

Start flat. Don't create `service/` or extra packages until a handler actually needs them.

## Backend: Go rules

- Latest stable Go. Standard library `net/http` with method + path routing. No web framework.
- `pgx/v5` with `pgxpool`. `sqlc` for all queries. **No ORM, no string-built SQL.**
- `goose` for migrations, plain SQL files, always with a working `-- +goose Down`.
- `oapi-codegen` generates the server interface from `api/openapi.yaml`; handlers implement it.
- Errors: wrap with context, `fmt.Errorf("load list %d: %w", id, err)`. Never swallow errors.
  Use `errors.Is` / `errors.As`. Map domain errors to HTTP status codes in one place.
- `context.Context` is the first parameter of anything doing I/O. Respect cancellation.
- No global mutable state. Dependencies are passed explicitly (constructor or struct fields).
- Logging with `log/slog`, structured, JSON in production.
- Config from environment variables, parsed once at startup, fail fast on missing values.
- Graceful shutdown on SIGINT/SIGTERM.
- Keep functions small and boring. Prefer clarity over cleverness. `gofmt` and `go vet` clean.

## Database rules

- Postgres. Every table has `id`, `created_at`, `updated_at`. Use `timestamptz`.
- Foreign keys and `NOT NULL` by default; justify any nullable column.
- Add indexes deliberately, and note why in the migration comment.
- For any query that will power a list or report, show the `EXPLAIN ANALYZE` against seeded
  data before calling it done.
- Keep a seed script that generates realistic volume (years of purchase history, several
  households) so performance problems show up locally.

## Frontend rules

- React + TypeScript, `strict: true`. No `any`; use `unknown` and narrow.
- API types are generated from the OpenAPI spec with `openapi-typescript`. Never hand-write
  response types.
- Server state lives in **TanStack Query**. No duplicating server data into `useState` or a
  global store. Local UI state stays local.
- Query keys are defined in one factory per feature (`features/lists/queries.ts`).
- Mutations that affect what the user sees use optimistic updates with rollback on error.
- Routing with TanStack Router, using loaders where it makes sense.
- Styling: plain CSS, BEM components plus a small set of utilities. Components own their
  layout; utilities are never load-bearing. No Tailwind, no CSS-in-JS.
- Accessible by default: semantic HTML, labels, keyboard support, visible focus.

## Testing

- Go: table-driven tests; handler tests with `httptest`; store tests against a real Postgres
  (test DB in compose), not mocks.
- Frontend: Vitest + Testing Library for hooks and components with real logic. Don't test
  trivial markup.
- Tests are part of the task, not a follow-up.

## Commands

Keep these working in the Makefile; update this list if they change.

- `make db-up` / `make db-down`: Postgres via compose
- `make migrate` / `make migrate-down`: goose
- `make gen`: sqlc + oapi-codegen + openapi-typescript
- `make seed`: load realistic fake data
- `make dev`: backend with reload + Vite dev server
- `make test`: Go and frontend tests
- `make lint`: gofmt, go vet, tsc --noEmit, eslint

## Workflow

1. For anything beyond a trivial change: short plan first (files touched, contract/schema
   impact, open questions). Wait for my OK.
2. Work in small, reviewable steps. One concern per commit.
3. Run `make gen`, `make lint`, `make test` before saying a task is done.
4. Summarize: what changed, what to look at closely, Go concepts used, any follow-ups.
5. If you're unsure about intent, ask. Don't guess and build on the guess.

## Roadmap

1. **CRUD foundation:** households, users, auth, lists, items. Learn Go layout and testing.
2. **In-store mode:** optimistic check-off, rollback, offline queue for bad signal.
3. **Shared lists:** two people editing one list live (SSE or polling), cache invalidation,
   explicit conflict rules.
4. **Insights:** purchase history, recurring items, "runs out when" forecast, spending over
   time. Real aggregation queries, one chart that genuinely communicates something.

Current phase: **1**.

## Don't

- Don't edit the OpenAPI spec or migrations without asking.
- Don't add dependencies without asking.
- Don't edit generated code (`internal/store`, `internal/api` generated files, `src/api` types).
- Don't build ahead of the current phase.
